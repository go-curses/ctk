package ctk

import (
	"github.com/go-curses/cdk"
)

const TypeAccelerator cdk.CTypeTag = "ctk-accelerator"

func init() {
	_ = cdk.TypesManager.AddType(TypeAccelerator, nil)
}

type Accelerator interface {
	Object

	Init() (already bool)
	LockAccel()
	UnlockAccel()
	IsLocked() (locked bool)
	Path() (path string)
	Key() (key cdk.Key)
	Mods() (mods cdk.ModMask)
	Match(key cdk.Key, mods cdk.ModMask) (match bool)
	Settings() (path string, key cdk.Key, mods cdk.ModMask)
	Configure(key cdk.Key, mods cdk.ModMask)
	UnsetKeyMods()
}

type CAccelerator struct {
	CObject

	accelLocking int
}

func NewDefaultAccelerator(path string) Accelerator {
	return NewAccelerator(path, cdk.KeyNUL, cdk.ModNone)
}

func NewAccelerator(path string, key cdk.Key, mods cdk.ModMask) Accelerator {
	a := &CAccelerator{}
	a.Init()
	a.Configure(key, mods)
	return a
}

func (a *CAccelerator) Init() (already bool) {
	if a.InitTypeItem(TypeAccelMap, a) {
		return true
	}
	a.CObject.Init()
	a.accelLocking = 0
	_ = a.InstallProperty(PropertyAccelPath, cdk.StringProperty, true, "")
	_ = a.InstallProperty(PropertyAccelKey, cdk.StructProperty, true, cdk.KeyNUL)
	_ = a.InstallProperty(PropertyAccelMods, cdk.StructProperty, true, cdk.ModNone)
	_ = a.InstallProperty(PropertyAccelLocked, cdk.BoolProperty, true, false)
	return false
}

func (a *CAccelerator) LockAccel() {
	a.Lock()
	a.accelLocking += 1
	a.Unlock()
	if a.accelLocking == 1 {
		if err := a.SetBoolProperty(PropertyAccelLocked, true); err != nil {
			a.LogErr(err)
		}
	}
}

func (a *CAccelerator) UnlockAccel() {
	a.Lock()
	a.accelLocking -= 1
	a.Unlock()
	if a.accelLocking == 0 {
		if err := a.SetBoolProperty(PropertyAccelLocked, false); err != nil {
			a.LogErr(err)
		}
	}
}

func (a *CAccelerator) IsLocked() (locked bool) {
	a.RLock()
	locked = a.accelLocking > 0
	a.RUnlock()
	return
}

func (a *CAccelerator) Path() (path string) {
	var err error
	if path, err = a.GetStringProperty(PropertyAccelPath); err != nil {
		a.LogErr(err)
	}
	return
}

func (a *CAccelerator) Key() (key cdk.Key) {
	var ok bool
	if v, err := a.GetStructProperty(PropertyAccelKey); err != nil {
		a.LogErr(err)
	} else if key, ok = v.(cdk.Key); !ok {
		key = cdk.KeyNUL
		a.LogError("value stored in %v is not of cdk.Key type: %v (%T)", PropertyAccelKey, v, v)
	}
	return
}

func (a *CAccelerator) Mods() (mods cdk.ModMask) {
	var ok bool
	if v, err := a.GetStructProperty(PropertyAccelMods); err != nil {
		a.LogErr(err)
	} else if mods, ok = v.(cdk.ModMask); !ok {
		mods = cdk.ModNone
		a.LogError("value stored in %v is not of cdk.ModMask type: %v (%T)", PropertyAccelMods, v, v)
	}
	return
}

func (a *CAccelerator) Match(key cdk.Key, mods cdk.ModMask) (match bool) {
	k, m := a.Key(), a.Mods()
	return key == k && mods == m
}

func (a *CAccelerator) Settings() (path string, key cdk.Key, mods cdk.ModMask) {
	path = a.Path()
	key = a.Key()
	mods = a.Mods()
	return
}

func (a *CAccelerator) Configure(key cdk.Key, mods cdk.ModMask) {
	if a.IsLocked() {
		a.LogError("accelerator is locked, cannot Configure: %v", a.Path())
		return
	}
	a.Freeze()
	if err := a.SetStructProperty(PropertyAccelKey, key); err != nil {
		a.LogErr(err)
	}
	if err := a.SetStructProperty(PropertyAccelMods, mods); err != nil {
		a.LogErr(err)
	}
	a.Thaw()
}

func (a *CAccelerator) UnsetKeyMods() {
	a.Configure(cdk.KeyNUL, cdk.ModNone)
}

const PropertyAccelPath cdk.Property = "accel-path"

const PropertyAccelKey cdk.Property = "accel-key"

const PropertyAccelMods cdk.Property = "accel-mods"

const PropertyAccelLocked cdk.Property = "accel-locked"
