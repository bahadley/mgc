package follower

import (
	"time"
)

type deadline interface {
	nextDeadline(t time.Time) time.Time
}

type last struct{}

func (n *last) nextDeadline(t time.Time) time.Time {
	var dl time.Time
	var lastDelay time.Duration

	for _, hb := range hbWindow {
		if !(hb == nil || hb.ArrivalTime.IsZero()) {
			lastDelay = hb.TransDelay
		}
	}

	if lastDelay > 0 {
		dl = t.Add(lastDelay)
	} else {
		dl = t.Add(time.Millisecond * 500)
	}

	return dl
}
