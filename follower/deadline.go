package follower

import (
	"time"
)

type deadline interface {
	nextFreshnessPoint(t time.Time) time.Time
	recordObservation(t time.Time, hb *event)
}

type last struct {
	lastDelayObs time.Duration
}

func (n *last) nextFreshnessPoint(t time.Time) time.Time {
	var fp time.Time

	if n.lastDelayObs != 0 {
		fp = t.Add(n.lastDelayObs)
	} else {
		fp = t.Add(time.Millisecond * 500)
	}

	return fp
}

func (n *last) recordObservation(t time.Time, hb *event) {
	n.lastDelayObs = (hb.eventTime).Sub(t)
}

func nextFreshnessPoint(t time.Time, d deadline) time.Duration {
	freshnessPoint := d.nextFreshnessPoint(t)
	return freshnessPoint.Sub(t)
}
