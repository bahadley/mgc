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

	eventChan   chan *common.Event
	outputChan  chan *common.Event
	reportChan  chan *common.Event
	verdictChan chan *common.Event

	wg sync.WaitGroup
)

func RunFailureDetector() {
	// Counting semaphore set to the number of threads.
	wg.Add(5)

	go output()
	go manageObservations()
	go stateControl()
	go controlLoop()
	go net.Ingress(eventChan)

	// Wait for the threads to finish.
	wg.Wait()
}

func controlLoop() {
	// Block waiting for coordinated regime start.
	timerStart := time.NewTimer(config.DurationToRegimeStart())
	<-timerStart.C

	var seqNo common.SeqNoType = 0
	// Tick when leader sends a heartbeat.
	ticker := time.NewTicker(config.DurationOfHeartbeatInterval())
	for t := range ticker.C {
		// Leader is scheduled to send, so get a deadline for the heartbeat.
		eventChan <- &common.Event{
			EventType: common.Query,
			EventTime: t,
			SeqNo:     seqNo}

		// Block waiting for deadline calc from observations.
		rptF := <-reportChan
		outputChan <- rptF
		deadline := time.NewTimer(rptF.FreshnessPoint.Sub(time.Now()))
		<-deadline.C

		// Determine if heartbeat arrived.
		eventChan <- &common.Event{
			EventType: common.FreshnessEvent,
			EventTime: time.Now(),
			SeqNo:     seqNo}

		seqNo++
	}
}

func stateControl() {
	for {
		verdict := <-verdictChan
		if leaderSuspect != verdict.Suspect {
			leaderSuspect = verdict.Suspect
			outputChan <- verdict
		}
	}
}

func output() {
	for {
		switch event := <-outputChan; event.EventType {
		case common.HeartbeatEvent:
			log.Info.Printf("%s|%d|||%d|||%d", event.EventType,
				event.EventTime.UnixNano(), event.SeqNo,
				event.HeartbeatDelay)
		case common.Query:
			log.Info.Printf("%s|%d|||%d||%d|", event.EventType,
				event.EventTime.UnixNano(), event.SeqNo,
				event.FreshnessPoint.Sub(event.EventTime))
		case common.Verdict:
			log.Info.Printf("%s|%d|||%d|%t||", event.EventType,
				event.EventTime.UnixNano(), event.SeqNo,
				event.Suspect)
		default:
			log.Error.Println("Invalid event type encountered")
		}
	}
}

func init() {
	eventChan = make(chan *common.Event, config.ChannelBufSz())
	outputChan = make(chan *common.Event, config.ChannelBufSz())
	reportChan = make(chan *common.Event)
	verdictChan = make(chan *common.Event)
}
