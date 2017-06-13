package follower

import (
	"time"

	"github.com/bahadley/mgc/common"
	"github.com/bahadley/mgc/log"
)

const (
	// Length of window.
	bufSz uint32 = 4
)

var (
	// Invariant:  Heartbeats are in descending order by event.seqNo.
	hbWindow []*common.Heartbeat

	// Used to calculate next freshness points.  Defined in freshnessPoint.go
	fpCalc deadline
)

func runObservations() {
	for {
		switch event := <-eventChan; event.EventType {
		case common.HeartbeatEvent:
			if !update(event.SeqNo, event.EventTime) {
				log.Warning.Printf("Heartbeat from %s with seqNo %d not registered",
					event.Src, event.SeqNo)
			}
			outputChan <- event
		case common.QueryEvent:
			if !insert(&common.Heartbeat{SeqNo: event.SeqNo, SendTime: event.EventTime}) {
				log.Warning.Printf("Heartbeat initialization with seqNo %d not inserted",
					event.SeqNo)
			}
			reportChan <- &report{freshnessPoint: event.EventTime.Add(time.Millisecond * 500)}
		case common.FreshnessEvent:
			reportChan <- &report{suspect: false}
		default:
			log.Error.Println("Invalid event type encountered")
		}
	}
}

func insert(tmp *common.Heartbeat) bool {
	inserted := false

	for idx, hb := range hbWindow {
		if inserted ||
			(!inserted && hb != nil && tmp.SeqNo > hb.SeqNo) {
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
		if hbWindow[idx] != nil && hbWindow[idx].SeqNo == seqNo {
			hbWindow[idx].ArrivalTime = arrivalTime
			hbWindow[idx].TransDelay = arrivalTime.Sub(hbWindow[idx].SendTime)
			updated = true
			break
		}
	}

	return updated
}

func init() {
	hbWindow = make([]*common.Heartbeat, bufSz)

	fpCalc = &last{}
}
