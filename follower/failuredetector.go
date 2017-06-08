package follower

import (
	"time"

	"github.com/bahadley/mgc/config"
	"github.com/bahadley/mgc/log"
	"github.com/bahadley/mgc/util"
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
	timer := time.NewTimer(util.DurationToRegimeStart())
	<-timer.C

	freshnessInterval, err := time.ParseDuration("500ms")
	if err != nil {
		log.Error.Fatal(err.Error())
	}
	ticker := time.NewTicker(freshnessInterval)

	for range ticker.C {
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
