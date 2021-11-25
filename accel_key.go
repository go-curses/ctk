package ctk

import (
	"fmt"

	"github.com/go-curses/cdk"
	"github.com/go-curses/ctk/lib/enums"
)

type AccelKey interface {
	GetKey() cdk.Key
	GetMods() cdk.ModMask
	GetFlags() enums.AccelFlags
	Match(key cdk.Key, mods cdk.ModMask) (match bool)
	String() string
}

type CAccelKey struct {
	Key   cdk.Key
	Mods  cdk.ModMask
	Flags enums.AccelFlags
}

func MakeAccelKey(key cdk.Key, mods cdk.ModMask, flags enums.AccelFlags) (accelKey AccelKey) {
	accelKey = &CAccelKey{
		Key:   key,
		Mods:  mods,
		Flags: flags,
	}
	return
}

func (a *CAccelKey) GetKey() cdk.Key {
	return a.Key
}

func (a *CAccelKey) GetMods() cdk.ModMask {
	return a.Mods
}

func (a *CAccelKey) GetFlags() enums.AccelFlags {
	return a.Flags
}

func (a *CAccelKey) Match(key cdk.Key, mods cdk.ModMask) (match bool) {
	k, m := a.GetKey(), a.GetMods()
	return key == k && mods == m
}

func (a *CAccelKey) String() (key string) {
	mods := a.Mods.String()
	name := cdk.LookupKeyName(a.Key)
	if len(mods) > 0 {
		return fmt.Sprintf("%v%v", mods, name)
	}
	return fmt.Sprintf("%v", name)
}
