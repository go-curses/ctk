package ctk

import (
	"github.com/go-curses/cdk"
)

type CAccelGroupEntry struct {
	Handle   string
	Closure  GClosure
	AccelKey AccelKey
}

func NewAccelGroupEntry(accelerator AccelKey, handle string, closure GClosure) (age *CAccelGroupEntry) {
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
