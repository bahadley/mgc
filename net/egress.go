package net

import (
	"bytes"
	"encoding/binary"
	"net"
	"time"

	"github.com/bahadley/mgc/config"
	"github.com/bahadley/mgc/common"
	"github.com/bahadley/mgc/log"
)

func Egress(dst string, input <-chan *common.Heartbeat, output chan<- *common.Heartbeat) {
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

	for {
		hb := <-input

		buf := new(bytes.Buffer)
		err := binary.Write(buf, binary.LittleEndian, hb.SeqNo)
		if err != nil {
			log.Error.Fatal(err.Error())
		}

		log.Trace.Printf("Tx(%s): % x", dstAddr, buf.Bytes())
		hb.SendTime = time.Now()
		output <- hb

		_, err = conn.Write(buf.Bytes())
		if err != nil {
			log.Warning.Println(err.Error())
		}
	}
}
