package ctk

import (
	"github.com/go-curses/cdk/lib/ptypes"
)

type AccelGroupEntry struct {
	Accelerator AccelKey
	Closure     GClosure
	Quark       ptypes.QuarkID
}

func NewAccelGroupEntry(key AccelKey, closure GClosure, quark ptypes.QuarkID) (age *AccelGroupEntry) {
	age = &AccelGroupEntry{
		Accelerator: key,
		Closure:     closure,
		Quark:       quark,
	}
	return
}
