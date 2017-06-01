package follower

import (
	"fmt"
	"net"
	"time"

	"github.com/bahadley/mgc/config"
	"github.com/bahadley/mgc/log"
)

var (
	outputChan chan string
)

func Ingress() {
	srcAddr, err := net.ResolveUDPAddr("udp",
		config.Addr()+":"+config.DstPort())
	if err != nil {
		log.Error.Fatal(err.Error())
	}

	conn, err := net.ListenUDP("udp", srcAddr)
	if err != nil {
		log.Error.Fatal(err.Error())
	}

	defer conn.Close()

	log.Info.Printf("Listening for heartbeats (%s UDP) ...",
		srcAddr.String())

	buf := make([]byte, 128, 1024)
	for {
		n, caddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Warning.Println(err.Error())
			continue
		}

		t := time.Now().UnixNano()
		log.Trace.Printf("Rx(%s): %s", caddr, buf[0:n])
		outputChan <- fmt.Sprintf("time (ns): %d, msg: %s", t, buf[0:n])
	}
}

func Output() {
	outputChan = make(chan string, 10)

	for {
		hb := <-outputChan
		fmt.Println(hb)
	}
}
