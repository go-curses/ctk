package ctk

import (
	"fmt"

	"github.com/go-curses/cdk"
	cmath "github.com/go-curses/cdk/lib/math"
	"github.com/go-curses/ctk/lib/enums"
	"github.com/gofrs/uuid"
)

const TypeAccelGroup cdk.CTypeTag = "ctk-accel-group"

func init() {
	_ = cdk.TypesManager.AddType(TypeAccelGroup, func() interface{} { return MakeAccelGroup() })
}

// AccelGroup Hierarchy:
//	Object
//	  +- AccelGroup
//
// An AccelGroup represents a group of keyboard accelerators, typically
// attached to a toplevel Window (with Window.AddAccelGroup). Usually
// you won't need to create a AccelGroup directly; instead CTK automatically
// sets up the accelerators for your menus in the item factory's AccelGroup.
// Note that accelerators are different from mnemonics. Accelerators are
// shortcuts for activating a menu item; they appear alongside the menu item
// they're a shortcut for. For example "Ctrl+Q" might appear alongside the
// "Quit" menu item. Mnemonics are shortcuts for GUI elements such as text
// entries or buttons; they appear as underlined characters. Menu items can have
// both accelerators and mnemonics, of course.
// See: Label.NewWithMnemonic()
type AccelGroup interface {
	Object

	Init() (already bool)
	AccelConnect(accelKey cdk.Key, accelMods cdk.ModMask, accelFlags enums.AccelFlags, handle string, closure enums.GClosure) (id uuid.UUID)
	ConnectByPath(accelPath string, handle string, closure enums.GClosure)
	AccelGroupActivate(keyval cdk.Key, modifier cdk.ModMask) (activated bool)
	AccelDisconnect(id uuid.UUID) (removed bool)
	DisconnectKey(accelKey cdk.Key, accelMods cdk.ModMask) (removed bool)
	Query(accelKey cdk.Key, accelMods cdk.ModMask) (entries []*CAccelGroupEntry)
	Activate(accelKey cdk.Key, accelMods cdk.ModMask) (value bool)
	LockGroup()
	UnlockGroup()
	GetIsLocked() (locked bool)
	GetModifierMask() (value cdk.ModMask)
	Find(findFunc AccelGroupFindFunc, data interface{}) (key *AccelKey)
	AcceleratorValid(keyval cdk.Key, modifiers cdk.ModMask) (valid bool)
	AcceleratorParse(accelerator string) (acceleratorKey cdk.Key, acceleratorMods cdk.ModMask)
	AcceleratorName(acceleratorKey cdk.Key, acceleratorMods cdk.ModMask) (value string)
	AcceleratorGetLabel(acceleratorKey cdk.Key, acceleratorMods cdk.ModMask) (value string)
	AcceleratorSetDefaultModMask(defaultModMask cdk.ModMask)
	AcceleratorGetDefaultModMask() (value int)
}

// The CAccelGroup structure implements the AccelGroup interface and is
// exported to facilitate type embedding with custom implementations. No member
// variables are exported as the interface methods are the only intended means
// of interacting with AccelGroup objects.
type CAccelGroup struct {
	CObject

	entries map[uuid.UUID]*CAccelGroupEntry
	locking int
}

// MakeAccelGroup is used by the Buildable system to construct a new AccelGroup.
func MakeAccelGroup() AccelGroup {
	return NewAccelGroup()
}

// NewAccelGroup is the constructor for new AccelGroup instances.
func NewAccelGroup() (value AccelGroup) {
	a := new(CAccelGroup)
	a.Init()
	return a
}

// Init initializes an AccelGroup object. This must be called at least once to
// set up the necessary defaults and allocate any memory structures. Calling
// this more than once is safe though unnecessary. Only the first call will
// result in any effect upon the AccelGroup instance. Init is used in the
// NewAccelGroup constructor and only necessary when implementing a derivative
// AccelGroup type.
func (a *CAccelGroup) Init() (already bool) {
	if a.InitTypeItem(TypeAccelGroup, a) {
		return true
	}
	a.CObject.Init()
	a.entries = make(map[uuid.UUID]*CAccelGroupEntry, 0)
	a.locking = 0
	_ = a.InstallProperty(PropertyIsLocked, cdk.BoolProperty, false, false)
	_ = a.InstallProperty(PropertyModifierMask, cdk.StructProperty, false, nil)
	return false
}

// AccelConnect installs an accelerator in this group. When accel_group is being
// activated in response to a call to AccelGroupsActivate, closure will be
// invoked if the accel_key and accel_mods from AccelGroupsActivate match those
// of this connection. The signature used for the closure is that of
// AccelGroupActivate. Note that, due to implementation details, a single
// closure can only be connected to one accelerator group.
//
// Parameters:
// 	accelKey	key value of the accelerator
// 	accelMods	modifier combination of the accelerator
// 	accelFlags	a flag mask to configure this accelerator
//	handle	string to tag the closure for later use
// 	closure	code to be executed upon accelerator activation
func (a *CAccelGroup) AccelConnect(accelKey cdk.Key, accelMods cdk.ModMask, accelFlags enums.AccelFlags, handle string, closure enums.GClosure) (id uuid.UUID) {
	a.CObject.Lock()
	key := NewAccelKey(accelKey, accelMods, accelFlags)
	age := NewCAccelGroupEntry(key, handle, closure)
	id, _ = uuid.NewV4()
	a.entries[id] = age
	a.CObject.Unlock()
	return
}

// ConnectByPath installs an accelerator in this group, using an accelerator
// path to look up the appropriate key and modifiers (see AccelMapAddEntry).
// When accel_group is being activated in response to a call to
// AccelGroupsActivate, closure will be invoked if the accel_key and accel_mods
// from AccelGroupsActivate match the key and modifiers for the path.
//
// Parameters:
// 	accelPath	path used for determining key and modifiers.
//	handle 	string to tag the closure for later use
// 	closure	code to be executed upon accelerator activation
func (a *CAccelGroup) ConnectByPath(accelPath string, handle string, closure enums.GClosure) {
	accelMap := GetAccelMap()
	if accelerator, ok := accelMap.LookupEntry(accelPath); ok {
		a.AccelConnect(accelerator.Key(), accelerator.Mods(), enums.ACCEL_VISIBLE, handle, closure)
	} else {
		a.LogError("accelerator path not found: %v", accelPath)
	}
}

// AccelGroupActivate queries for entries matching the given keyval and
// modifier, then calling Closure functions for each entry found until one of
// them returns TRUE or the list of entries is exhausted, returning FALSE.
//
// Parameters:
// 	accelKey	key value of the accelerator
// 	accelMods	modifier combination of the accelerator
func (a *CAccelGroup) AccelGroupActivate(keyval cdk.Key, modifier cdk.ModMask) (activated bool) {
	entries := a.Query(keyval, modifier)
	for _, entry := range entries {
		if entry.Closure(keyval, modifier, entry.AccelKey.GetFlags()) {
			activated = true
			break
		}
	}
	return false
}

// AccelDisconnect removes an accelerator previously installed through Connect.
//
// Parameters:
// 	accelGroup	the accelerator group to remove an accelerator from
// 	closure	handle for the closure code to remove
func (a *CAccelGroup) AccelDisconnect(id uuid.UUID) (removed bool) {
	a.Lock()
	defer a.Unlock()
	if _, ok := a.entries[id]; ok {
		delete(a.entries, id)
		return true
	}
	return false
}

// DisconnectKey removes an accelerator previously installed through Connect.
//
// Parameters:
// 	accelKey	key value of the accelerator
// 	accelMods	modifier combination of the accelerator
func (a *CAccelGroup) DisconnectKey(accelKey cdk.Key, accelMods cdk.ModMask) (removed bool) {
	a.Lock()
	defer a.Unlock()
	for id, entry := range a.entries {
		if entry.AccelKey.GetKey() == accelKey {
			if entry.AccelKey.GetMods() == accelMods {
				delete(a.entries, id)
				return true
			}
		}
	}
	return false
}

// Query searches an accelerator group for all entries matching accel_key and
// accel_mods.
//
// Parameters:
// 	accelGroup	the accelerator group to query
// 	accelKey	key value of the accelerator
// 	accelMods	modifier combination of the accelerator
func (a *CAccelGroup) Query(accelKey cdk.Key, accelMods cdk.ModMask) (entries []*CAccelGroupEntry) {
	a.RLock()
	defer a.RUnlock()
	for _, entry := range a.entries {
		if entry.Match(accelKey, accelMods) {
			entries = append(entries, entry)
		}
	}
	return
}

// Activate finds the first accelerator in accel_group that matches accel_key
// and accel_mods, and activates it.
//
// Parameters:
// 	accelKey	accelerator keyval from a key event
// 	accelMods	keyboard state mask from a key event
func (a *CAccelGroup) Activate(accelKey cdk.Key, accelMods cdk.ModMask) (value bool) {
	a.RLock()
	defer a.RUnlock()
	for _, entry := range a.entries {
		if entry.AccelKey.GetKey() == accelKey && entry.AccelKey.GetMods() == accelMods {
			return entry.Closure(accelKey, accelMods, entry.AccelKey.GetFlags())
		}
	}
	return false
}

// LockGroup locks the given accelerator group. Locking an accelerator group
// prevents the accelerators contained within it to be changed during runtime.
// Refer o AccelMap.ChangeEntry() about runtime accelerator changes. If called
// more than once, accel_group remains locked until UnlockGroup has been called
// an equivalent number of times.
func (a *CAccelGroup) LockGroup() {
	a.Lock()
	defer a.Unlock()
	a.locking += 1
	if a.locking == 1 {
		if err := a.SetBoolProperty(PropertyIsLocked, true); err != nil {
			a.LogErr(err)
		}
	}
}

// UnlockGroup releases the last call to LockGroup.
func (a *CAccelGroup) UnlockGroup() {
	a.Lock()
	defer a.Unlock()
	a.locking = cmath.FloorI(a.locking-1, 0)
	if a.locking == 0 {
		if err := a.SetBoolProperty(PropertyIsLocked, true); err != nil {
			a.LogErr(err)
		}
	}
}

// GetIsLocked checks if the group is locked or not. Locks are added and removed
// using LockGroup and UnlockGroup.
func (a *CAccelGroup) GetIsLocked() (locked bool) {
	var err error
	if locked, err = a.GetBoolProperty(PropertyIsLocked); err != nil {
		a.LogErr(err)
	}
	return
}

// GetModifierMask returns a cdk.ModMask representing the mask for this
// accel_group. For example, CONTROL_MASK, SHIFT_MASK, etc.
func (a *CAccelGroup) GetModifierMask() (value cdk.ModMask) {
	var ok bool
	if v, err := a.GetStructProperty(PropertyModifierMask); err != nil {
		a.LogErr(err)
	} else if value, ok = v.(cdk.ModMask); !ok {
		a.LogError("value stored in PropertyModifierMask is not of cdk.ModMask: %v (%T)", v, v)
	}
	return
}

// Find finds the first entry in an accelerator group for which find_func
// returns TRUE and returns its AccelKey.
//
// Parameters:
// 	findFunc	a function to filter the entries of accel_group
// 	data	arbitrary data to pass to find_func
func (a *CAccelGroup) Find(findFunc AccelGroupFindFunc, data interface{}) (key *AccelKey) {
	a.LogError("method not implemented")
	return
}

// AcceleratorValid determines whether a given keyval and modifier mask
// constitute a valid keyboard accelerator.
//
// Parameters:
// 	keyval	a GDK keyval
// 	modifiers	modifier mask
func (a *CAccelGroup) AcceleratorValid(keyval cdk.Key, modifiers cdk.ModMask) (valid bool) {
	a.LogError("method unimplemented")
	return true
}

// AcceleratorParse parses a string representing an accelerator. The format
// looks like "<Control>a" or "<Shift><Alt>F1". The parser is fairly liberal and
// allows lower or upper case, and also abbreviations such as "<Ctl>" and
// "<Ctrl>". Key names are parsed using KeyvalFromName. For character keys the
// name is not the symbol, but the lowercase name, e.g. one would use
// "<Ctrl>minus" instead of "<Ctrl>-". If the parse fails, accelerator_key and
// accelerator_mods will be set to 0 (zero).
//
// Parameters:
// 	accelerator	string representing an accelerator
func (a *CAccelGroup) AcceleratorParse(accelerator string) (acceleratorKey cdk.Key, acceleratorMods cdk.ModMask) {
	a.LogError("method unimplemented")
	return
}

// AcceleratorName converts an accelerator keyval and modifier mask into a
// string parseable by AcceleratorParse. For example, if you pass in
// cdk.KeySmallQ and CONTROL_MASK, this function returns "<Control>q". If you
// need to display accelerators in the user interface, see AcceleratorGetLabel.
//
// Parameters:
// 	acceleratorKey	accelerator keyval
// 	acceleratorMods	accelerator modifier mask
func (a *CAccelGroup) AcceleratorName(acceleratorKey cdk.Key, acceleratorMods cdk.ModMask) (value string) {
	return fmt.Sprintf("%v%v", acceleratorMods.String(), cdk.LookupKeyName(acceleratorKey))
}

// AcceleratorGetLable converts an accelerator keyval and modifier mask into a
// string which can be used to represent the accelerator to the user.
//
// Parameters:
// 	acceleratorKey	accelerator keyval
// 	acceleratorMods	accelerator modifier mask
func (a *CAccelGroup) AcceleratorGetLabel(acceleratorKey cdk.Key, acceleratorMods cdk.ModMask) (value string) {
	return fmt.Sprintf("%v%v", acceleratorMods.String(), cdk.LookupKeyName(acceleratorKey))
}

// AcceleratorSetDefaultMask updates the modifiers that will be considered
// significant for keyboard accelerators. The default mod mask is CONTROL_MASK |
// SHIFT_MASK | MOD1_MASK | SUPER_MASK | HYPER_MASK | META_MASK, that is,
// Control, Shift, Alt, Super, Hyper and Meta. Other modifiers will by default
// be ignored by AccelGroup. You must include at least the three modifiers
// Control, Shift and Alt in any value you pass to this function. The default
// mod mask should be changed on application startup, before using any
// accelerator groups.
//
// Parameters:
// 	defaultModMask	accelerator modifier mask
func (a *CAccelGroup) AcceleratorSetDefaultModMask(defaultModMask cdk.ModMask) {}

// AcceleratorGetDefaultModMask returns the value set by
// AcceleratorSetDefaultModMask.
//
// Parameters:
// 	returns	the default accelerator modifier mask
func (a *CAccelGroup) AcceleratorGetDefaultModMask() (value int) {
	return 0
}

// Is the accel group locked.
// Flags: Read
// Default value: FALSE
const PropertyIsLocked cdk.Property = "is-locked"

// Modifier Mask.
// Flags: Read
// Default value: GDK_SHIFT_MASK | GDK_CONTROL_MASK | GDK_MOD1_MASK | GDK_SUPER_MASK | GDK_HYPER_MASK | GDK_META_MASK
const PropertyModifierMask cdk.Property = "modifier-mask"

// The accel-activate signal is an implementation detail of AccelGroup and
// not meant to be used by applications.
const SignalAccelActivate cdk.Signal = "accel-activate"

// The accel-changed signal is emitted when a CAccelGroupEntry is added to
// or removed from the accel group. Widgets like AccelLabel which display
// an associated accelerator should connect to this signal, and rebuild their
// visual representation if the accel_closure is theirs.
// Listener function arguments:
// 	keyval int	the accelerator keyval
// 	modifier cdk.ModMask	the modifier combination of the accelerator
// 	accelClosure GClosure	the GClosure of the accelerator
const SignalAccelChanged cdk.Signal = "accel-changed"

type AccelGroupFindFunc = func(key AccelKey, closure enums.GClosure, data []interface{}) bool
