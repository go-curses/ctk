package ctk

import (
	"fmt"

	"github.com/go-curses/cdk"
	"github.com/go-curses/cdk/lib/ptypes"
)

// CDK type-tag for AccelGroup objects
const TypeAccelGroup cdk.CTypeTag = "ctk-accel-group"

func init() {
	_ = cdk.TypesManager.AddType(TypeAccelGroup, func() interface{} { return MakeAccelGroup() })
}

// AccelGroup Hierarchy:
//	Object
//	  +- AccelGroup
// A AccelGroup represents a group of keyboard accelerators, typically
// attached to a toplevel Window (with WindowAddAccelGroup). Usually
// you won't need to create a AccelGroup directly; instead, when using
// ItemFactory, CTK automatically sets up the accelerators for your menus in
// the item factory's AccelGroup. Note that accelerators are different from
// mnemonics. Accelerators are shortcuts for activating a menu item; they
// appear alongside the menu item they're a shortcut for. For example
// "Ctrl+Q" might appear alongside the "Quit" menu item. Mnemonics are
// shortcuts for GUI elements such as text entries or buttons; they appear as
// underlined characters. See LabelNewWithMnemonic. Menu items can
// have both accelerators and mnemonics, of course.
type AccelGroup interface {
	Object

	Init() (already bool)
	AccelConnect(accelKey cdk.Key, accelMods cdk.ModMask, accelFlags AccelFlags, closure GClosure) (id int)
	ConnectByPath(accelPath string, closure GClosure)
	AccelGroupActivate(acceleratable Object, keyval cdk.Key, modifier cdk.ModMask) (activated bool)
	AccelDisconnect(id int) (removed bool)
	DisconnectKey(accelKey cdk.Key, accelMods cdk.ModMask) (removed bool)
	Query(accelKey cdk.Key, accelMods cdk.ModMask) (entries []*AccelGroupEntry)
	Activate(accelQuark ptypes.QuarkID, acceleratable Object, accelKey cdk.Key, accelMods cdk.ModMask) (value bool)
	Lock()
	Unlock()
	GetIsLocked() (locked bool)
	FromAccelClosure(closure GClosure) (value AccelGroup)
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
// of interacting with AccelGroup objects
type CAccelGroup struct {
	CObject

	entries map[int]*AccelGroupEntry
	locking int
}

// Default constructor for AccelGroup objects
func MakeAccelGroup() *CAccelGroup {
	return NewAccelGroup()
}

// Constructor for AccelGroup objects
func NewAccelGroup() (value *CAccelGroup) {
	a := new(CAccelGroup)
	a.Init()
	return a
}

// AccelGroup object initialization. This must be called at least once to setup
// the necessary defaults and allocate any memory structures. Calling this more
// than once is safe though unnecessary. Only the first call will result in any
// effect upon the AccelGroup instance
func (a *CAccelGroup) Init() (already bool) {
	if a.InitTypeItem(TypeAccelGroup, a) {
		return true
	}
	a.CObject.Init()
	a.entries = make(map[int]*AccelGroupEntry, 0)
	a.locking = 0
	_ = a.InstallProperty(PropertyIsLocked, cdk.BoolProperty, false, false)
	_ = a.InstallProperty(PropertyModifierMask, cdk.StructProperty, false, nil)
	return false
}

// Installs an accelerator in this group. When accel_group is being activated
// in response to a call to AccelGroupsActivate, closure will be
// invoked if the accel_key and accel_mods from AccelGroupsActivate
// match those of this connection. The signature used for the closure is that
// of AccelGroupActivate. Note that, due to implementation details, a
// single closure can only be connected to one accelerator group.
// Parameters:
// 	accelGroup	the accelerator group to install an accelerator in
// 	accelKey	key value of the accelerator
// 	accelMods	modifier combination of the accelerator
// 	accelFlags	a flag mask to configure this accelerator
// 	closure	closure to be executed upon accelerator activation
func (a *CAccelGroup) AccelConnect(accelKey cdk.Key, accelMods cdk.ModMask, accelFlags AccelFlags, closure GClosure) (id int) {
	key := MakeAccelKey(accelKey, accelMods, accelFlags)
	age := NewAccelGroupEntry(key, closure, ptypes.QuarkFromString(accelMods.String()))
	next := 0
	for idx, _ := range a.entries {
		if next <= idx {
			next = idx + 1
		}
	}
	a.entries[next] = age
	id = next
	return
}

// Installs an accelerator in this group, using an accelerator path to look
// up the appropriate key and modifiers (see AccelMapAddEntry). When
// accel_group is being activated in response to a call to
// AccelGroupsActivate, closure will be invoked if the accel_key and
// accel_mods from AccelGroupsActivate match the key and modifiers
// for the path. The signature used for the closure is that of
// AccelGroupActivate. Note that accel_path string will be stored in a
// GQuark. Therefore, if you pass a static string, you can save some memory
// by interning it first with g_intern_static_string.
// Parameters:
// 	accelGroup	the accelerator group to install an accelerator in
// 	accelPath	path used for determining key and modifiers.
// 	closure	closure to be executed upon accelerator activation
func (a *CAccelGroup) ConnectByPath(accelPath string, closure GClosure) {}

func (a *CAccelGroup) AccelGroupActivate(acceleratable Object, keyval cdk.Key, modifier cdk.ModMask) (activated bool) {
	return false
}

// Removes an accelerator previously installed through
// Connect. Since 2.20 closure can be NULL.
// Parameters:
// 	accelGroup	the accelerator group to remove an accelerator from
// 	closure	the closure to remove from this accelerator group, or NULL
// to remove all closures.
// 	returns	TRUE if the closure was found and got disconnected
func (a *CAccelGroup) AccelDisconnect(id int) (removed bool) {
	if _, ok := a.entries[id]; ok {
		delete(a.entries, id)
		return true
	}
	return false
}

// Removes an accelerator previously installed through
// Connect.
// Parameters:
// 	accelGroup	the accelerator group to install an accelerator in
// 	accelKey	key value of the accelerator
// 	accelMods	modifier combination of the accelerator
// 	returns	TRUE if there was an accelerator which could be
// removed, FALSE otherwise
func (a *CAccelGroup) DisconnectKey(accelKey cdk.Key, accelMods cdk.ModMask) (removed bool) {
	for id, entry := range a.entries {
		if entry.Accelerator.Key == accelKey {
			if entry.Accelerator.Mods == accelMods {
				delete(a.entries, id)
				return true
			}
		}
	}
	return false
}

// Queries an accelerator group for all entries matching accel_key and
// accel_mods .
// Parameters:
// 	accelGroup	the accelerator group to query
// 	accelKey	key value of the accelerator
// 	accelMods	modifier combination of the accelerator
// 	nEntries	location to return the number of entries found, or NULL.
// 	returns	an array of n_entries
// AccelGroupEntry elements, or NULL. The array is owned by CTK and must not be freed.
func (a *CAccelGroup) Query(accelKey cdk.Key, accelMods cdk.ModMask) (entries []*AccelGroupEntry) {
	for _, entry := range a.entries {
		if entry.Accelerator.Key == accelKey && entry.Accelerator.Mods == accelMods {
			entries = append(entries, entry)
		}
	}
	return
}

// Finds the first accelerator in accel_group that matches accel_key and
// accel_mods , and activates it.
// Parameters:
// 	accelQuark	the quark for the accelerator name
// 	acceleratable	the GObject, usually a Window, on which
// to activate the accelerator.
// 	accelKey	accelerator keyval from a key event
// 	accelMods	keyboard state mask from a key event
// Returns:
// 	TRUE if an accelerator was activated and handled this keypress
func (a *CAccelGroup) Activate(accelQuark ptypes.QuarkID, acceleratable Object, accelKey cdk.Key, accelMods cdk.ModMask) (value bool) {
	for _, entry := range a.entries {
		if entry.Accelerator.Key == accelKey && entry.Accelerator.Mods == accelMods {
			return entry.Closure(acceleratable, accelKey, accelMods, entry.Accelerator.Flags)
		}
	}
	return false
}

// Locks the given accelerator group. Locking an acelerator group prevents
// the accelerators contained within it to be changed during runtime. Refer
// to AccelMapChangeEntry about runtime accelerator changes. If
// called more than once, accel_group remains locked until
// Unlock has been called an equivalent number of times.
func (a *CAccelGroup) Lock() {
	a.locking += 1
	_ = a.SetBoolProperty(PropertyIsLocked, true)
}

// Undoes the last call to Lock on this accel_group .
func (a *CAccelGroup) Unlock() {
	a.locking -= 1
	if a.locking <= 0 {
		_ = a.SetBoolProperty(PropertyIsLocked, false)
	}
}

// Locks are added and removed using Lock and
// Unlock.
// Returns:
// 	TRUE if there are 1 or more locks on the accel_group , FALSE
// 	otherwise.
func (a *CAccelGroup) GetIsLocked() (locked bool) {
	var err error
	if locked, err = a.GetBoolProperty(PropertyIsLocked); err != nil {
		a.LogErr(err)
	}
	return
}

// Finds the AccelGroup to which closure is connected; see
// Connect.
// Parameters:
// 	closure	a GClosure
func (a *CAccelGroup) FromAccelClosure(closure GClosure) (value AccelGroup) {
	return nil
}

// Gets a cdk.ModMask representing the mask for this accel_group . For
// example, GDK_CONTROL_MASK, GDK_SHIFT_MASK, etc.
// Returns:
// 	the modifier mask for this accel group.
func (a *CAccelGroup) GetModifierMask() (value cdk.ModMask) {
	var ok bool
	if v, err := a.GetStructProperty(PropertyModifierMask); err != nil {
		a.LogErr(err)
	} else if value, ok = v.(cdk.ModMask); !ok {
		a.LogError("value stored in PropertyModifierMask is not of cdk.ModMask: %v (%T)", v, v)
	}
	return
}

// Finds the first entry in an accelerator group for which find_func returns
// TRUE and returns its AccelKey.
// Parameters:
// 	findFunc	a function to filter the entries of accel_group
// with
// 	data	data to pass to find_func
//
// 	returns	the key of the first entry passing
// find_func
// . The key is owned by CTK and must not be freed.
func (a *CAccelGroup) Find(findFunc AccelGroupFindFunc, data interface{}) (key *AccelKey) {
	a.LogError("method not implemented")
	return
}

// Determines whether a given keyval and modifier mask constitute a valid
// keyboard accelerator. For example, the GDK_a keyval plus GDK_CONTROL_MASK
// is valid - this is a "Ctrl+a" accelerator. But, you can't, for instance,
// use the GDK_Control_L keyval as an accelerator.
// Parameters:
// 	keyval	a GDK keyval
// 	modifiers	modifier mask
// 	returns	TRUE if the accelerator is valid
func (a *CAccelGroup) AcceleratorValid(keyval cdk.Key, modifiers cdk.ModMask) (valid bool) {
	a.LogError("method unimplemented")
	return true
}

// Parses a string representing an accelerator. The format looks like
// "<Control>a" or "<Shift><Alt>F1" or "<Release>z" (the last one is for key
// release). The parser is fairly liberal and allows lower or upper case, and
// also abbreviations such as "<Ctl>" and "<Ctrl>". Key names are parsed
// using KeyvalFromName. For character keys the name is not the
// symbol, but the lowercase name, e.g. one would use "<Ctrl>minus" instead
// of "<Ctrl>-". If the parse fails, accelerator_key and accelerator_mods
// will be set to 0 (zero).
// Parameters:
// 	accelerator	string representing an accelerator
// Returns:
// 	acceleratorKey	    keyval.
// 	acceleratorMods	    modifier mask.
func (a *CAccelGroup) AcceleratorParse(accelerator string) (acceleratorKey cdk.Key, acceleratorMods cdk.ModMask) {
	a.LogError("method unimplemented")
	return
}

// Converts an accelerator keyval and modifier mask into a string parseable
// by AcceleratorParse. For example, if you pass in GDK_q and
// GDK_CONTROL_MASK, this function returns "<Control>q". If you need to
// display accelerators in the user interface, see
// AcceleratorGetLabel.
// Parameters:
// 	acceleratorKey	accelerator keyval
// 	acceleratorMods	accelerator modifier mask
// Returns:
// 	a newly-allocated accelerator name
func (a *CAccelGroup) AcceleratorName(acceleratorKey cdk.Key, acceleratorMods cdk.ModMask) (value string) {
	return fmt.Sprintf("%v %v", acceleratorMods.String(), cdk.LookupKeyName(acceleratorKey))
}

// Converts an accelerator keyval and modifier mask into a string which can
// be used to represent the accelerator to the user.
// Parameters:
// 	acceleratorKey	accelerator keyval
// 	acceleratorMods	accelerator modifier mask
// Returns:
// 	a newly-allocated string representing the accelerator.
func (a *CAccelGroup) AcceleratorGetLabel(acceleratorKey cdk.Key, acceleratorMods cdk.ModMask) (value string) {
	return fmt.Sprintf("%v %v", acceleratorMods.String(), cdk.LookupKeyName(acceleratorKey))
}

// Sets the modifiers that will be considered significant for keyboard
// accelerators. The default mod mask is GDK_CONTROL_MASK | GDK_SHIFT_MASK |
// GDK_MOD1_MASK | GDK_SUPER_MASK | GDK_HYPER_MASK | GDK_META_MASK, that is,
// Control, Shift, Alt, Super, Hyper and Meta. Other modifiers will by
// default be ignored by AccelGroup. You must include at least the three
// modifiers Control, Shift and Alt in any value you pass to this function.
// The default mod mask should be changed on application startup, before
// using any accelerator groups.
// Parameters:
// 	defaultModMask	accelerator modifier mask
func (a *CAccelGroup) AcceleratorSetDefaultModMask(defaultModMask cdk.ModMask) {}

// Gets the value set by AcceleratorSetDefaultModMask.
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

// The accel-changed signal is emitted when a AccelGroupEntry is added to
// or removed from the accel group. Widgets like AccelLabel which display
// an associated accelerator should connect to this signal, and rebuild their
// visual representation if the accel_closure is theirs.
// Listener function arguments:
// 	keyval int	the accelerator keyval
// 	modifier cdk.ModMask	the modifier combination of the accelerator
// 	accelClosure GClosure	the GClosure of the accelerator
const SignalAccelChanged cdk.Signal = "accel-changed"

type GClosure = func(argv ...interface{}) (handled bool)

type AccelGroupFindFunc = func(key AccelKey, closure GClosure, data []interface{}) bool

// func (a *CAccelGroup) GtkAccelGroupFindFunc(key AccelKey, closure GClosure, data interface{}) (value bool) {
// 	return false
// }
