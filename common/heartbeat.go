package common

import (
	"time"
)

type (
	SeqNoType uint16

	Heartbeat struct {
		Dst         string
		SeqNo       SeqNoType
		SendTime    time.Time
		ArrivalTime time.Time
		TransDelay  time.Duration
	}
)
