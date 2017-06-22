package net

import (
	"bytes"
	"encoding/binary"
	"net"
	"time"

	"github.com/bahadley/mgc/common"
	"github.com/bahadley/mgc/config"
	"github.com/bahadley/mgc/log"
)

func Ingress(output chan<- *common.Event) {
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
	var seqNo common.SeqNoType
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

		output <- &common.Event{
			EventTime: time.Now(),
			EventType: common.HeartbeatEvent,
			Src:       caddr.String(),
			SeqNo:     seqNo}

		log.Trace.Printf("Rx(%s): % x", caddr, buf[0:n])
	}
}
