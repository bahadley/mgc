package follower

import (
	"bytes"
	"encoding/binary"
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
	var seqNo uint16
	for {
		n, caddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Warning.Println(err.Error())
			continue
		}

		seqNoBuf := bytes.NewReader(buf)
		err = binary.Read(seqNoBuf, binary.LittleEndian, &seqNo)
		if err != nil {
			log.Error.Fatal(err.Error())
		}

		heartbeatChan <- &heartbeat{
			src:         caddr.String(),
			seqNo:       seqNo,
			arrivalTime: time.Now()}

		log.Trace.Printf("Rx(%s): % x", caddr, buf[0:n])
	}
}
