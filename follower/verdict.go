package follower

import (
	"github.com/bahadley/mgc/common"
)

type verdict interface {
	check(seqNo common.SeqNoType) bool
}

type basic struct{}

func (b *basic) check(seqNo common.SeqNoType) bool {
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
