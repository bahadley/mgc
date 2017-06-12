package follower

import (
	"sync"
	"time"

	"github.com/bahadley/mgc/config"
	"github.com/bahadley/mgc/log"
)

const (
	// Types for events passed in the eventChan
	heartbeatEvent = "H"
	freshnessEvent = "F"
	queryEvent     = "Q"
)

type (
	event struct {
		eventTime time.Time
		eventType string
		src       string
		seqNo     uint16
	}

	report struct {
		suspect        bool
		freshnessPoint time.Time
	}
)

var (
	// Leader is suspect if true, trusted if false
	leaderSuspect bool

	eventChan  chan *event
	reportChan chan *report
	outputChan chan *event

	wg sync.WaitGroup
)

func RunFailureDetector() {
	// Counting semaphore set to the number of threads.
	wg.Add(4)

	go runOutput()
	go runObservations()
	go runControlLoop()
	go runIngress()

	// Wait for the threads to finish.
	wg.Wait()
}

func runControlLoop() {
	timerStart := time.NewTimer(config.DurationToRegimeStart())
	<-timerStart.C

	var seqNo uint16 = 0
	ticker := time.NewTicker(config.DurationOfHeartbeatInterval())
	for t := range ticker.C {
		eventChan <- &event{
			eventTime: t,
			eventType: queryEvent,
			seqNo:     seqNo}
		seqNo++

		rptF := <-reportChan
		freshnessPoint := time.NewTimer(rptF.freshnessPoint.Sub(time.Now()))
		<-freshnessPoint.C

		eventChan <- &event{
			eventTime: time.Now(),
			eventType: freshnessEvent}

		rptL := <-reportChan
		leaderSuspect = rptL.suspect
	}
}

func runOutput() {
	for {
		switch event := <-outputChan; event.eventType {
		case heartbeatEvent:
			log.Info.Printf("Rcvd heartbeat: time (ns): %d, seqno: %d",
				event.eventTime.UnixNano(), event.seqNo)
		case freshnessEvent:
			log.Info.Printf("Freshness point: time (ns) %d",
				event.eventTime.UnixNano())
		default:
			log.Error.Println("Invalid event type encountered")
		}
	}
}

func init() {
	eventChan = make(chan *event, config.ChannelBufSz())
	reportChan = make(chan *report)
	outputChan = make(chan *event, config.ChannelBufSz())
}
