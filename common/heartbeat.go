package common

import (
	"time"
)

type (
	Heartbeat struct {
		Dst         string
		SeqNo       uint16
		SendTime    time.Time
		ArrivalTime time.Time
		TransDelay  time.Duration
	}
)
