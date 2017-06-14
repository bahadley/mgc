package follower

type verdict interface {
	check(seqNo uint16) bool
}

type basic struct{}

func (b *basic) check(seqNo uint16) bool {
	suspect := true

	// Search for heartbeat with matching sequence number.
	for _, hb := range hbWindow {
		if hb != nil && hb.SeqNo == seqNo {
			// Leader is not suspected if ArrivalTime is set.
			suspect = hb.ArrivalTime.IsZero()
		}
	}

	return suspect
}
