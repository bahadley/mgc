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

	go runOutput()
	go runObservations()
	go runControlLoop()
	go net.Ingress(eventChan)

	// Wait for the threads to finish.
	wg.Wait()
}

func runControlLoop() {
	timerStart := time.NewTimer(config.DurationToRegimeStart())
	<-timerStart.C

	var seqNo uint16 = 0
	ticker := time.NewTicker(config.DurationOfHeartbeatInterval())
	for t := range ticker.C {
		eventChan <- &common.Event{
			EventTime: t,
			EventType: common.QueryEvent,
			SeqNo:     seqNo}
		seqNo++

		rptF := <-reportChan
		freshnessPoint := time.NewTimer(rptF.freshnessPoint.Sub(time.Now()))
		<-freshnessPoint.C

		eventChan <- &common.Event{
			EventTime: time.Now(),
			EventType: common.FreshnessEvent}

		rptL := <-reportChan
		leaderSuspect = rptL.suspect
	}
}

func runOutput() {
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
