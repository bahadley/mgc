package leader

import (
	"net"
	"sync"
	"time"

	"github.com/bahadley/mgc/config"
	"github.com/bahadley/mgc/log"
)

var (
	wg sync.WaitGroup
)

func Transmit() {
	dsts := config.DstAddrs()

	// Counting semaphore set to the number of addrs.
	wg.Add(len(dsts))

	// Launch all threads.  Each thread has a different destination.
	for _, dst := range dsts {
		go egress(dst)
	}

	// Wait for the threads to finish.
	wg.Wait()
}

func egress(dst string) {
	dstAddr, err := net.ResolveUDPAddr("udp",
		dst+":"+config.DstPort())
	if err != nil {
		log.Error.Fatal(err.Error())
	}

	srcAddr, err := net.ResolveUDPAddr("udp",
		config.Addr()+":0")
	if err != nil {
		log.Error.Fatal(err.Error())
	}

	conn, err := net.DialUDP("udp", srcAddr, dstAddr)
	if err != nil {
		log.Error.Fatal(err.Error())
	}

	defer conn.Close()
	defer wg.Done()

	hbts := config.NumHeartbeats()
	delayInt := config.DelayInterval()
	msg := []byte("alive")

	for i := 0; i < hbts; i++ {
		log.Trace.Printf("Tx(%s): %s", dstAddr, msg)

		_, err = conn.Write(msg)
		if err != nil {
			log.Warning.Println(err.Error())
		}

		time.Sleep(delayInt * time.Millisecond)
	}
}
