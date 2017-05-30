package main

import (
	"flag"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bahadley/mgc/heartbeat"
	"github.com/bahadley/mgc/log"
)

func main() {
	addr := flag.String("addr", "localhost", "host address (IPv4)")
	dst := flag.String("dst", "localhost", "comma delimited list of destination addresses (IPv4)")
	flag.Parse()
	log.Info.Println("Starting up ...")

	// Allow the node to be shut down gracefully.
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		// Block waiting for signal.
		<-c
		log.Info.Println("Shutting down ...")
		os.Exit(0)
	}()

	heartbeat.Transmit(*addr, strings.Split(*dst, ","))

	log.Info.Println("Shutting down ...")
}
