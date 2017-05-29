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

func Transmit() {
	addrs := DstAddr()

	// Counting semaphore set to the number of addrs.
	wg.Add(len(addrs))

	// Launch all threads.  Each thread has a different destination.
	for _, addr := range addrs {
		go egress(addr)
	}

	// Wait for the threads to finish.
	wg.Wait()
}

func egress(addr string) {
	dstAddr, err := net.ResolveUDPAddr("udp",
		addr+":"+DstPort())
	if err != nil {
		log.Error.Fatal(err.Error())
	}

	srcAddr, err := net.ResolveUDPAddr("udp",
		Addr()+":0")
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
