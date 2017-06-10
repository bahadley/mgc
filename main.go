package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/bahadley/mgc/config"
	"github.com/bahadley/mgc/follower"
	"github.com/bahadley/mgc/leader"
	"github.com/bahadley/mgc/log"
)

func main() {
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

	if config.IsLeader() {
		leader.RunHeartbeats()
	} else if config.IsFollower() {
		follower.RunFailureDetector()
	} else {
		log.Error.Println("Role not implemented")
	}

	log.Info.Println("Shutting down ...")
}
