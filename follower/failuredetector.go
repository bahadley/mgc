package follower

import (
	"sync"
	"time"

	"github.com/bahadley/mgc/common"
	"github.com/bahadley/mgc/config"
	"github.com/bahadley/mgc/log"
	"github.com/bahadley/mgc/net"
)

type (
	report struct {
		suspect        bool
		freshnessPoint time.Time
	}
)

var (
	// Leader is suspect if true, trusted if false
	leaderSuspect bool

	eventChan  chan *common.Event
	reportChan chan *report
	outputChan chan *common.Event

	wg sync.WaitGroup
)

func RunFailureDetector() {
	// Counting semaphore set to the number of threads.
	wg.Add(4)

	go output()
	go runObservations()
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
		// Leader is sending, get a deadline for that heartbeat.
		eventChan <- &common.Event{
			EventTime: t,
			EventType: common.QueryEvent,
			SeqNo:     seqNo}
		seqNo++

		// Block waiting for deadline calc from observations.
		rptF := <-reportChan
		deadline := time.NewTimer(rptF.freshnessPoint.Sub(time.Now()))
		<-deadline.C

		// Determine if heartbeat arrived.
		eventChan <- &common.Event{
			EventTime: time.Now(),
			EventType: common.FreshnessEvent}

		// Block waiting for trust/suspect verdict.
		rptL := <-reportChan
		leaderSuspect = rptL.suspect
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
	reportChan = make(chan *report)
	outputChan = make(chan *common.Event, config.ChannelBufSz())
}
