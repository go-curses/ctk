package ctk

import (
	"fmt"

	"github.com/go-curses/cdk"
)

type AccelKey interface {
	GetKey() cdk.Key
	GetMods() cdk.ModMask
	GetFlags() AccelFlags
	Match(key cdk.Key, mods cdk.ModMask) (match bool)
	String() string
}

type CAccelKey struct {
	Key   cdk.Key
	Mods  cdk.ModMask
	Flags AccelFlags
}

func MakeAccelKey(key cdk.Key, mods cdk.ModMask, flags AccelFlags) (accelKey AccelKey) {
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

func (a *CAccelKey) GetFlags() AccelFlags {
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
