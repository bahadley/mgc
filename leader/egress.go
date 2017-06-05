package leader

import (
	"bytes"
	"encoding/binary"
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
		dst+":"+config.Port())
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
	var seqNo uint16 = 0

	timer := time.NewTimer((config.Start()).Sub(time.Now()))
	<-timer.C

	for i := 0; i < hbts; i++ {
		buf := new(bytes.Buffer)
		err := binary.Write(buf, binary.LittleEndian, seqNo)
		if err != nil {
			log.Error.Fatal(err.Error())
		}

		log.Trace.Printf("Tx(%s): % x", dstAddr, buf.Bytes())

		_, err = conn.Write(buf.Bytes())
		if err != nil {
			log.Warning.Println(err.Error())
		}

		seqNo++
		time.Sleep(delayInt * time.Millisecond)
	}
}
