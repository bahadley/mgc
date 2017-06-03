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
	printChan     chan *heartbeat
)

func IngestHeartbeats() {
	for {
		hb := <-heartbeatChan
		printChan <- hb
	}
}

func Print() {
	printChan = make(chan *heartbeat, config.ChannelBufSz())

	for {
		hb := <-printChan
		log.Info.Printf("Rcvd heartbeat: time (ns): %d, seqno: %d",
			hb.arrivalTime.UnixNano(), hb.seqNo)
	}
}

func init() {
	heartbeatChan = make(chan *heartbeat, config.ChannelBufSz())
}
