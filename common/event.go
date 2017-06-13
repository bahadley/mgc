package common

import (
	"time"
)

const (
	// Types for Event.EventType
	HeartbeatEvent = "H"
	FreshnessEvent = "F"
	QueryEvent     = "Q"
)

type (
	Event struct {
		EventTime time.Time
		EventType string
		Src       string
		SeqNo     uint16
	}
)
