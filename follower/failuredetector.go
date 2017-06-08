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

	n := noop{}
	ticker := time.NewTicker(config.DurationOfHeartbeatInterval())
	for range ticker.C {
		timerFreshnessPoint := time.NewTimer(durationToNextFreshnessPoint(n))
		<-timerFreshnessPoint.C

		latestHeartbeat := ingestHeartbeats()
		if latestHeartbeat == nil {
			log.Info.Println("Leader is suspect")
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

	for hb := range heartbeatChan {
		outputChan <- hb
		if latestHeartbeat == nil ||
			(hb.arrivalTime).After(latestHeartbeat.arrivalTime) {
			latestHeartbeat = hb
		}
	}

	return latestHeartbeat
}

func init() {
	heartbeatChan = make(chan *heartbeat, config.ChannelBufSz())
	outputChan = make(chan *heartbeat, config.ChannelBufSz())
}
