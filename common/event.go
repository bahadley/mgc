package common

import (
	"time"
)

const (
	// Types for Event.EventType
	HeartbeatEvent = "H"
	FreshnessEvent = "F"
	Query          = "Q"
	Verdict        = "V"
	SendHeartbeat  = "S"
)

type (
	Event struct {
		EventType      string
		EventTime      time.Time
		Src            string
		Dst            string
		SeqNo          SeqNoType
		Suspect        bool
		FreshnessPoint time.Time
		HeartbeatDelay time.Duration
	}
)
