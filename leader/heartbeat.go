package leader

import (
	"fmt"
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
	heartbeatChan = make(chan *Heartbeat)

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
	printChan = make(chan *Heartbeat, config.ChannelBufSz())

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

	timer := time.NewTimer(durationToRegimeStart())
	<-timer.C

	ticker := time.NewTicker(durationOfHeartbeatInterval())
	for range ticker.C {
		heartbeatChan <- &Heartbeat{dst: dst, seqNo: seqNo}
		seqNo++
	}
}

func durationToRegimeStart() time.Duration {
	return (config.Start()).Sub(time.Now())
}

func durationOfHeartbeatInterval() time.Duration {
	d, err := time.ParseDuration(fmt.Sprintf("%dms", config.DelayInterval()))
	if err != nil {
		log.Error.Fatal(err.Error())
	}
	return d
}
