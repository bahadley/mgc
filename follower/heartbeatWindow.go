package follower

import (
	"time"

	"github.com/bahadley/mgc/log"
)

const (
	// Length of window.
	bufSz uint32 = 4
)

type (
	heartbeat struct {
		seqNo       uint16
		sendTime    time.Time
		arrivalTime time.Time
		transDelay  time.Duration
	}
)

var (
	// Invariant:  Heartbeats are in descending order by event.seqNo.
	hbWindow []*heartbeat

	// Used to calculate next freshness points.  Defined in freshnessPoint.go
	fpCalc deadline
)

func runObservations() {
	for {
		switch event := <-eventChan; event.eventType {
		case heartbeatEvent:
			if !update(event.seqNo, event.eventTime) {
				log.Warning.Printf("Heartbeat from %s with seqNo %d not registered",
					event.src, event.seqNo)
			}
			outputChan <- event
		case queryEvent:
			if !insert(&heartbeat{seqNo: event.seqNo, sendTime: event.eventTime}) {
				log.Warning.Printf("Heartbeat initialization with seqNo %d not inserted",
					event.seqNo)
			}
			reportChan <- &report{freshnessPoint: event.eventTime.Add(time.Millisecond * 500)}
		case freshnessEvent:
			reportChan <- &report{suspect: false}
		default:
			log.Error.Println("Invalid event type encountered")
		}
	}
}

func insert(tmp *heartbeat) bool {
	inserted := false

	for idx, hb := range hbWindow {
		if inserted ||
			(!inserted && hb != nil && tmp.seqNo > hb.seqNo) {
			// Insert the new heartbeat and shift the subsequent heartbeats towards
			// the back of the window.  The last heartbeat will fall off if the
			// window is full.
			hbWindow[idx] = tmp
			tmp = hb
			inserted = true
		} else if !inserted && hb == nil {
			// Window is currently empty and this is the first arriving heartbeat.
			hbWindow[idx] = tmp
			inserted = true
			break
		}
	}

	return inserted
}

func update(seqNo uint16, arrivalTime time.Time) bool {
	updated := false

	for idx := bufSz - 1; idx >= 0; idx-- {
		if hbWindow[idx] != nil && hbWindow[idx].seqNo == seqNo {
			hbWindow[idx].arrivalTime = arrivalTime
			hbWindow[idx].transDelay = arrivalTime.Sub(hbWindow[idx].sendTime)
			updated = true
			break
		}
	}

	return updated
}

func init() {
	hbWindow = make([]*heartbeat, bufSz)

	fpCalc = &last{}
}
