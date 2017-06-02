package follower

import (
	"net"
	"time"

	"github.com/bahadley/mgc/config"
	"github.com/bahadley/mgc/log"
)

func Ingress() {
	srcAddr, err := net.ResolveUDPAddr("udp",
		config.Addr()+":"+config.Port())
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

	buf := make([]byte, config.TupleBufLen(), config.TupleBufCap())
	for {
		n, caddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Warning.Println(err.Error())
			continue
		}

		heartbeatChan <- &heartbeat{
			src:         caddr.String(),
			seqNo:       string(buf[0:n]),
			arrivalTime: time.Now().UnixNano()}

		log.Trace.Printf("Rx(%s): %s", caddr, buf[0:n])
	}
}
