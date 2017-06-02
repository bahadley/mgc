package follower

import (
	"github.com/bahadley/mgc/config"
	"github.com/bahadley/mgc/log"
)

type (
	heartbeat struct {
		src         string
		seqNo       string
		arrivalTime int64
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
		log.Info.Printf("Rcvd heartbeat: time (ns): %d, seqno: %s",
			hb.arrivalTime, hb.seqNo)
	}
}

func init() {
	heartbeatChan = make(chan *heartbeat, config.ChannelBufSz())
}
