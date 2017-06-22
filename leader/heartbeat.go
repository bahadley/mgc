package leader

import (
	"sync"
	"time"

	"github.com/bahadley/mgc/common"
	"github.com/bahadley/mgc/config"
	"github.com/bahadley/mgc/log"
	"github.com/bahadley/mgc/net"
)

var (
	heartbeatChan chan *common.Heartbeat
	outputChan    chan *common.Event

	wg sync.WaitGroup
)

func RunHeartbeats() {
	go runOutput()

	dsts := config.DstAddrs()
	// Counting semaphore set to the number of addrs.
	wg.Add(len(dsts))

	// Launch all threads.  Each thread has a different follower.
	for _, dst := range dsts {
		go runPushToFollower(dst)
	}

	// Wait for the threads to finish.
	wg.Wait()
}

func runPushToFollower(dst string) {
	defer wg.Done()
	go net.Egress(dst, heartbeatChan, outputChan)

	timer := time.NewTimer(config.DurationToRegimeStart())
	<-timer.C

	var seqNo common.SeqNoType = 0
	ticker := time.NewTicker(config.DurationOfHeartbeatInterval())
	for range ticker.C {
		heartbeatChan <- &common.Heartbeat{Dst: dst, SeqNo: seqNo}
		seqNo++
	}
}

func runOutput() {
	for {
		event := <-outputChan
        log.Info.Printf("%s|%d||%s|%d||", event.EventType,
            event.EventTime.UnixNano(), event.Dst, event.SeqNo)
	}
}

func init() {
	heartbeatChan = make(chan *common.Heartbeat)
	outputChan = make(chan *common.Event, config.ChannelBufSz())
}
