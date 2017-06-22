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

	// Used to calculate next freshness points.
	dlCalc deadline
	// Used to determine the trust/suspect verdict for the leader.
	vCalc verdict
)

func manageObservations() {
	for {
		switch event := <-eventChan; event.EventType {
		// Heartbeat received from network interface.  Set the arrival time in
		// matching shell record in the observation window
		case common.HeartbeatEvent:
			if !update(event.SeqNo, event.EventTime) {
				log.Warning.Printf("Heartbeat from %s with seqNo %d not registered",
					event.Src, event.SeqNo)
			}
			//outputChan <- event
		// Leader will be sending a heartbeat now.  Calculate a deadline and
		// create a shell record in the observation window.
		case common.Query:
			reportChan <- &common.Event{
				EventType:      common.Query,
				EventTime:      event.EventTime,
				FreshnessPoint: dlCalc.nextDeadline(event.EventTime)}
			if !insert(&common.Heartbeat{
				SeqNo:    event.SeqNo,
				SendTime: event.EventTime}) {
				log.Warning.Printf("Heartbeat initialization with seqNo %d not inserted",
					event.SeqNo)
			}
		// Deadline has expired.  Determine if a heartbeat has arrived for
		// this period.
		case common.FreshnessEvent:
			reportChan <- &common.Event{
				EventType: common.Verdict,
				EventTime: time.Now(),
				SeqNo:     event.SeqNo,
				Suspect:   vCalc.check(event.SeqNo)}
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

func update(seqNo common.SeqNoType, arrivalTime time.Time) bool {
	updated := false

	// Search for heartbeat with same sequence number.
	for _, hb := range hbWindow {
		if hb != nil && hb.SeqNo == seqNo {
			hb.ArrivalTime = arrivalTime
			hb.TransDelay = arrivalTime.Sub(hb.SendTime)
			updated = true
			break
		}
	}

	return updated
}

func init() {
	hbWindow = make([]*common.Heartbeat, bufSz)

	dlCalc = &last{}
	vCalc = &basic{}
}
