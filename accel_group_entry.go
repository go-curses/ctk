package ctk

import (
	"github.com/go-curses/cdk"
	"github.com/go-curses/ctk/lib/enums"
)

type CAccelGroupEntry struct {
	Handle   string
	Closure  enums.GClosure
	AccelKey AccelKey
}

func NewCAccelGroupEntry(accelerator AccelKey, handle string, closure enums.GClosure) (age *CAccelGroupEntry) {
	age = &CAccelGroupEntry{
		Handle:   handle,
		Closure:  closure,
		AccelKey: accelerator,
	}
	return
}

func (a *CAccelGroupEntry) Match(key cdk.Key, modifier cdk.ModMask) (match bool) {
	return a.AccelKey.Match(key, modifier)
}
