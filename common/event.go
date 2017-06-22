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
)

type (
	Event struct {
		EventTime      time.Time
		EventType      string
		Src            string
		SeqNo          SeqNoType
		Suspect        bool
		FreshnessPoint time.Time
	}
)
