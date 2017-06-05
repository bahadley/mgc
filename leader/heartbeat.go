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
	printChan     chan *Heartbeat

	wg sync.WaitGroup
)

func PushHeartbeats() {
	dsts := config.DstAddrs()

	// Counting semaphore set to the number of addrs.
	wg.Add(len(dsts))

	// Launch all threads.  Each thread has a different destination.
	for _, dst := range dsts {
		go pushToFollower(dst)
	}

	// Wait for the threads to finish.
	wg.Wait()
}

func pushToFollower(dst string) {
	defer wg.Done()

	go Egress(dst, heartbeatChan, printChan)

	//hbts := config.NumHeartbeats()
	delayInt := config.DelayInterval()
	var seqNo uint16 = 0

	timer := time.NewTimer((config.Start()).Sub(time.Now()))
	<-timer.C

	//for i := 0; i < hbts; i++ {
	for {
		heartbeatChan <- &Heartbeat{dst: dst, seqNo: seqNo}
		seqNo++
		time.Sleep(delayInt * time.Millisecond)
	}

}

func Print() {
	for {
		hb := <-printChan
		log.Info.Printf("Sent heartbeat: time (ns): %d, dst: %s, seqno: %d",
			hb.transmitTime.UnixNano(), hb.dst, hb.seqNo)
	}
}

func init() {
	printChan = make(chan *Heartbeat, config.ChannelBufSz())
	heartbeatChan = make(chan *Heartbeat)
}
