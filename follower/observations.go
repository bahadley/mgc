package follower

import (
	"time"

	"github.com/bahadley/mgc/log"
)

const (
	// Length of window.
	bufSz uint32 = 4
)

var (
	// Invariant:  Heartbeats are in descending order by event.seqNo.
	hbWindow []*event

	// Used to calculate next freshness points.  Defined in freshnessPoint.go
	fpCalc deadline
)

func runObservations() {
	for {
		switch event := <-eventChan; event.eventType {
		case heartbeatEvent:
			if !insert(event, hbWindow) {
				log.Warning.Printf("Heartbeat from %s with seqNo %d not inserted",
					event.src, event.seqNo)
			}
			outputChan <- event
		case queryEvent:
			reportChan <- &report{freshnessPoint: event.eventTime.Add(time.Millisecond * 500)}
		case freshnessEvent:
			reportChan <- &report{suspect: false}
		default:
			log.Error.Println("Invalid event type encountered")
		}
	}
}

func insert(tmp *event, window []*event) bool {
	inserted := false

	for idx, hb := range window {
		if inserted ||
			(!inserted && hb != nil && tmp.seqNo > hb.seqNo) {
			// Insert the new heartbeat and shift the subsequent heartbeats towards
			// the back of the window.  The last heartbeat will fall off if the
			// window is full.
			window[idx] = tmp
			tmp = hb
			inserted = true
		} else if !inserted && hb == nil {
			// Window is currently empty and this is the first arriving heartbeat, or ...
			// Out of order arrival and there is room at the back of the window.
			window[idx] = tmp
			inserted = true
			break
		}
	}

	return inserted
}

func init() {
	hbWindow = make([]*event, bufSz)

	fpCalc = &last{}
}
