package follower

import (
	"time"

	"github.com/bahadley/mgc/config"
	"github.com/bahadley/mgc/log"
)

const (
	heartbeatId      = "H"
	freshnessEventId = "F"
)

type (
	event struct {
		eventTime time.Time
		eventType string
		src       string
		seqNo     uint16
	}
)

var (
	eventChan  chan *event
	outputChan chan *event
)

func RunFailureDetector() {
	timerStart := time.NewTimer(config.DurationToRegimeStart())
	<-timerStart.C

	ticker := time.NewTicker(config.DurationOfHeartbeatInterval())
	for range ticker.C {
		freshnessPoint := time.NewTimer(time.Millisecond * 500)
		<-freshnessPoint.C

		ingestHeartbeats()
	}
}

func Output() {
	for {
		hb := <-outputChan
		log.Info.Printf("Rcvd heartbeat: time (ns): %d, seqno: %d",
			hb.eventTime.UnixNano(), hb.seqNo)
	}
}

func ingestHeartbeats() {
	hb := <-eventChan
	outputChan <- hb
}

func init() {
	eventChan = make(chan *event, config.ChannelBufSz())
	outputChan = make(chan *event, config.ChannelBufSz())
}
