package follower

import (
	"sync"
	"time"

	"github.com/bahadley/mgc/common"
	"github.com/bahadley/mgc/config"
	"github.com/bahadley/mgc/log"
	"github.com/bahadley/mgc/net"
)

var (
	// Leader is suspect if true, trusted if false
	leaderSuspect bool

	eventChan  chan *common.Event
	reportChan chan *common.Event
	outputChan chan *common.Event

	wg sync.WaitGroup
)

func RunFailureDetector() {
	// Counting semaphore set to the number of threads.
	wg.Add(4)

	go output()
	go manageObservations()
	go controlLoop()
	go net.Ingress(eventChan)

	// Wait for the threads to finish.
	wg.Wait()
}

func controlLoop() {
	// Block waiting for coordinated regime start.
	timerStart := time.NewTimer(config.DurationToRegimeStart())
	<-timerStart.C

	var seqNo uint16 = 0
	// Tick when leader sends a heartbeat.
	ticker := time.NewTicker(config.DurationOfHeartbeatInterval())
	for t := range ticker.C {
		// Leader is scheduled to send, so get a deadline for the heartbeat.
		eventChan <- &common.Event{
			EventTime: t,
			EventType: common.QueryEvent,
			SeqNo:     seqNo}

		// Block waiting for deadline calc from observations.
		rptF := <-reportChan
		deadline := time.NewTimer(rptF.FreshnessPoint.Sub(time.Now()))
		<-deadline.C

		// Determine if heartbeat arrived.
		eventChan <- &common.Event{
			EventTime: time.Now(),
			EventType: common.FreshnessEvent,
			SeqNo: seqNo}

		// Block waiting for trust/suspect verdict.
		rptL := <-reportChan
		leaderSuspect = rptL.Suspect

		seqNo++
	}
}

func output() {
	for {
		switch event := <-outputChan; event.EventType {
		case common.HeartbeatEvent:
			log.Info.Printf("Rcvd heartbeat: time (ns): %d, seqno: %d",
				event.EventTime.UnixNano(), event.SeqNo)
		case common.FreshnessEvent:
			log.Info.Printf("Freshness point: time (ns) %d",
				event.EventTime.UnixNano())
		default:
			log.Error.Println("Invalid event type encountered")
		}
	}
}

func init() {
	eventChan = make(chan *common.Event, config.ChannelBufSz())
	reportChan = make(chan *common.Event)
	outputChan = make(chan *common.Event, config.ChannelBufSz())
}
