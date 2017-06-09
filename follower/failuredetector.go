package follower

import (
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
)

func RunFailureDetector() {
	go runObservations()

	timerStart := time.NewTimer(config.DurationToRegimeStart())
	<-timerStart.C

	ticker := time.NewTicker(config.DurationOfHeartbeatInterval())
	for t := range ticker.C {
		eventChan <- &event{
			eventTime: t,
			eventType: queryEvent}

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

func Output() {
	for {
		event := <-outputChan
		if event.eventType == heartbeatEvent {
			log.Info.Printf("Rcvd heartbeat: time (ns): %d, seqno: %d",
				event.eventTime.UnixNano(), event.seqNo)
		} else if event.eventType == freshnessEvent {
			log.Info.Printf("Freshness point: time (ns) %d",
				event.eventTime.UnixNano())
		}
	}
}

func init() {
	eventChan = make(chan *event, config.ChannelBufSz())
	reportChan = make(chan *report)
	outputChan = make(chan *event, config.ChannelBufSz())
}
