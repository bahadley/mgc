package follower

import (
	"github.com/bahadley/mgc/config"
	"github.com/bahadley/mgc/log"
)

type (
	Heartbeat struct {
		Src         string
		SeqNo       int32
		ArrivalTime int64
	}
)

var (
	HeartbeatChan chan *Heartbeat
    printChan chan *Heartbeat 
)

func IngestHeartbeats() {
	for {
		hb := <-HeartbeatChan
        printChan <- hb 
	}
}

func Print() {
    printChan = make(chan *Heartbeat, config.ChannelBufSz())

    for {
        hb := <-printChan
        log.Info.Printf("Rcvd heartbeat: time (ns): %d, seqno: %s", 
            hb.ArrivalTime, hb.SeqNo)
    }
}

func init() {
	HeartbeatChan = make(chan *Heartbeat, config.ChannelBufSz())
}
