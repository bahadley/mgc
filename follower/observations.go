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
	window []*event
)

func runObservations() {
	for {
		event := <-eventChan

		if event.eventType == heartbeatEvent {
			if !insert(event) {
				log.Warning.Printf("Heartbeat from %s with seqNo %d not inserted",
					event.src, event.seqNo)
			}
		} else if event.eventType == queryEvent {
			reportChan <- &report{
				freshnessPoint: event.eventTime.Add(time.Millisecond * 500)}
		} else if event.eventType == freshnessEvent {
			reportChan <- &report{
				suspect: false}
		}

		outputChan <- event
	}
}

func insert(tmp *event) bool {
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
	window = make([]*event, bufSz)
}
