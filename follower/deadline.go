package follower

import (
	"time"
)

type deadline interface {
	predictor() int
	safetyMargin() int
}

type noop struct{}

func (n noop) predictor() int {
	return 0
}

func (n noop) safetyMargin() int {
	return 0
}

func durationToNextFreshnessPoint(d deadline) time.Duration {
	//t := d.predictor() + d.safetyMargin()
	return time.Millisecond * 500
}
