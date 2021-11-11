package ctk

import (
	"fmt"

	"github.com/go-curses/cdk"
)

type AccelKey struct {
	Key   cdk.Key
	Mods  cdk.ModMask
	Flags AccelFlags
}

func MakeAccelKey(key cdk.Key, mods cdk.ModMask, flags AccelFlags) (accelKey AccelKey) {
	accelKey = AccelKey{
		Key:   key,
		Mods:  mods,
		Flags: flags,
	}
	return
}

func (a AccelKey) String() (key string) {
	return fmt.Sprintf("%v %v", a.Mods.String(), cdk.LookupKeyName(a.Key))
}
