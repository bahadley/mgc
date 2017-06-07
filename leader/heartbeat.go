package leader

import (
	"sync"
	"time"

	"github.com/bahadley/mgc/config"
	"github.com/bahadley/mgc/log"
	"github.com/bahadley/mgc/util"
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
	printChan     chan *Heartbeat

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

func Print() {
	for {
		hb := <-printChan
		log.Info.Printf("Sent heartbeat: time (ns): %d, dst: %s, seqno: %d",
			hb.transmitTime.UnixNano(), hb.dst, hb.seqNo)
	}
}

func pushToFollower(dst string) {
	defer wg.Done()

	go egress(dst, heartbeatChan, printChan)

	var seqNo uint16 = 0

	ticker := time.NewTicker(util.DurationOfHeartbeatInterval())
	timer := time.NewTimer(util.DurationToRegimeStart())
	<-timer.C

	for range ticker.C {
		heartbeatChan <- &Heartbeat{dst: dst, seqNo: seqNo}
		seqNo++
	}
}

func init() {
	heartbeatChan = make(chan *Heartbeat)
	printChan = make(chan *Heartbeat, config.ChannelBufSz())
}
