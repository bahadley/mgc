package follower

import (
	"time"

	"github.com/bahadley/mgc/config"
	"github.com/bahadley/mgc/log"
)

type (
	heartbeat struct {
		src         string
		seqNo       uint16
		arrivalTime time.Time
	}
)

var (
	heartbeatChan chan *heartbeat
	outputChan    chan *heartbeat
)

func RunFailureDetector() {
	timerStart := time.NewTimer(config.DurationToRegimeStart())
	<-timerStart.C

	var d deadline = &last{}
	ticker := time.NewTicker(config.DurationOfHeartbeatInterval())
	for t := range ticker.C {
		freshnessPoint := time.NewTimer(nextFreshnessPoint(t, d))
		<-freshnessPoint.C

		latestHeartbeat := ingestHeartbeats()
		if latestHeartbeat == nil {
			log.Info.Println("Leader is suspect")
		} else {
			d.recordObservation(t, latestHeartbeat)
		}
	}
}

func Output() {
	for {
		hb := <-outputChan
		log.Info.Printf("Rcvd heartbeat: time (ns): %d, seqno: %d",
			hb.arrivalTime.UnixNano(), hb.seqNo)
	}
}

func ingestHeartbeats() *heartbeat {
	var latestHeartbeat *heartbeat

	select {
	case hb := <-heartbeatChan:
		outputChan <- hb
		if latestHeartbeat == nil ||
			(hb.arrivalTime).After(latestHeartbeat.arrivalTime) {
			latestHeartbeat = hb
		}
	default:
		latestHeartbeat = nil
	}

	return latestHeartbeat
}

func init() {
	heartbeatChan = make(chan *heartbeat)
	outputChan = make(chan *heartbeat, config.ChannelBufSz())
}
