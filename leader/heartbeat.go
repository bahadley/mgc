package leader

import (
	"sync"
	"time"

	"github.com/bahadley/mgc/config"
	"github.com/bahadley/mgc/log"
)

type (
	Heartbeat struct {
		dst          string
		seqNo        uint16
		transmitTime time.Time
	}
)

var (
	heartbeatChan chan *Heartbeat
	outputChan     chan *Heartbeat

	wg sync.WaitGroup
)

func PushHeartbeats() {
	dsts := config.DstAddrs()

	// Counting semaphore set to the number of addrs.
	wg.Add(len(dsts))

	// Launch all threads.  Each thread has a different follower.
	for _, dst := range dsts {
		go pushToFollower(dst)
	}

	// Wait for the threads to finish.
	wg.Wait()
}

func Output() {
	for {
		hb := <-outputChan
		log.Info.Printf("Sent heartbeat: time (ns): %d, dst: %s, seqno: %d",
			hb.transmitTime.UnixNano(), hb.dst, hb.seqNo)
	}
}

func pushToFollower(dst string) {
	defer wg.Done()

	go egress(dst, heartbeatChan, outputChan)

	var seqNo uint16 = 0

	ticker := time.NewTicker(config.DurationOfHeartbeatInterval())
	timer := time.NewTimer(config.DurationToRegimeStart())
	<-timer.C

	for range ticker.C {
		heartbeatChan <- &Heartbeat{dst: dst, seqNo: seqNo}
		seqNo++
	}
}

func init() {
	heartbeatChan = make(chan *Heartbeat)
	outputChan = make(chan *Heartbeat, config.ChannelBufSz())
}
