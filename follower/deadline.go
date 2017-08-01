package follower

import (
	"time"

	"github.com/bahadley/mgc/config"
)

type deadline interface {
	nextDeadline(t time.Time) time.Time
}

type last struct{}

func (n *last) nextDeadline(t time.Time) time.Time {
	var predictor time.Time
	var lastDelay time.Duration

	// Search the observation window for the most recent observation
	// with a heartbeat recorded.
	for _, hb := range hbWindow {
		if !(hb == nil || hb.ArrivalTime.IsZero()) {
			lastDelay = hb.TransDelay
			break
		}
	}

	if lastDelay > 0 {
		predictor = t.Add(lastDelay)
	} else {
		predictor = t.Add(config.DefaultDeadline() * time.Millisecond)
	}

	// Add a constant safety margin.
	return predictor.Add(config.DefaultSafetyMargin() * time.Millisecond)
}
