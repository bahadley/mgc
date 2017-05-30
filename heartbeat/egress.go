package heartbeat

import (
	"net"
	"sync"
	"time"

	"github.com/bahadley/mgc/log"
)

var (
	send bool

	wg sync.WaitGroup
)

func Transmit(addr string, dsts []string) {
	//addrs := DstAddr()

	// Counting semaphore set to the number of addrs.
	wg.Add(len(dsts))

	// Launch all threads.  Each thread has a different destination.
	for _, dst := range dsts {
		go egress(addr, dst)
	}

	// Wait for the threads to finish.
	wg.Wait()
}

func egress(addr string, dst string) {
	dstAddr, err := net.ResolveUDPAddr("udp",
		dst+":"+DstPort())
	if err != nil {
		log.Error.Fatal(err.Error())
	}

	srcAddr, err := net.ResolveUDPAddr("udp",
		addr+":0")
	if err != nil {
		log.Error.Fatal(err.Error())
	}

	conn, err := net.DialUDP("udp", srcAddr, dstAddr)
	if err != nil {
		log.Error.Fatal(err.Error())
	}

	defer conn.Close()
	defer wg.Done()

	hbs := NumHeartbeats()
	delayInt := DelayInterval()
	msg := []byte("alive")

	for i := 0; i < hbs; i++ {
		log.Trace.Printf("Tx(%s): %s", dstAddr, msg)

		if send {
			_, err = conn.Write(msg)
			if err != nil {
				log.Warning.Println(err.Error())
			}
		}

		time.Sleep(delayInt * time.Millisecond)
	}
}

func init() {
	send = Send()
}
