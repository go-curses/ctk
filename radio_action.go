package ctk

import (
	"github.com/go-curses/cdk"
)

const TypeRadioAction cdk.CTypeTag = "ctk-radio-action"

func init() {
	_ = cdk.TypesManager.AddType(TypeRadioAction, func() interface{} { return MakeRadioAction() })
}

// RadioAction Hierarchy:
//	Object
//	  +- Action
//	    +- ToggleAction
//	      +- RadioAction
type RadioAction interface {
	ToggleAction

	Init() (already bool)
	GetGroup() (value ActionGroup)
	SetGroup(group ActionGroup)
	GetCurrentValue() (value int)
	SetCurrentValue(currentValue int)
}

// The CRadioAction structure implements the RadioAction interface and is
// exported to facilitate type embedding with custom implementations. No member
// variables are exported as the interface methods are the only intended means
// of interacting with RadioAction objects.
type CRadioAction struct {
	CToggleAction
}

// MakeRadioAction is used by the Buildable system to construct a new RadioAction.
func MakeRadioAction() RadioAction {
	return NewRadioAction("", "", "", "", 0)
}

// NewRadioAction is the constructor for new RadioAction instances.
func NewRadioAction(name string, label string, tooltip string, stockId string, value int) (r RadioAction) {
	r = new(CRadioAction)
	r.Init()
	return r
}

// Init initializes an RadioAction object. This must be called at least once to
// set up the necessary defaults and allocate any memory structures. Calling
// this more than once is safe though unnecessary. Only the first call will
// result in any effect upon the RadioAction instance. Init is used in the
// NewRadioAction constructor and only necessary when implementing a derivative
// RadioAction type.
func (r *CRadioAction) Init() (already bool) {
	if r.InitTypeItem(TypeRadioAction, r) {
		return true
	}
	r.CToggleAction.Init()
	_ = r.InstallProperty(PropertyCurrentValue, cdk.IntProperty, true, 0)
	_ = r.InstallProperty(PropertyGroup, cdk.StructProperty, true, nil)
	_ = r.InstallProperty(PropertyValue, cdk.IntProperty, true, 0)
	return false
}

// GetGroup returns the list representing the radio group for this object. Note
// that the returned list is only valid until the next change to the group.
//
// Parameters:
// 	action	the action object
func (r *CRadioAction) GetGroup() (value ActionGroup) {
	var ok bool
	if v, err := r.GetStructProperty(PropertyGroup); err != nil {
		r.LogErr(err)
	} else if value, ok = v.(ActionGroup); !ok {
		value = nil
		r.LogError("value stored in %v property is not of ActionGroup type: %v (%T)", PropertyGroup, v, v)
	}
	return
}

// SetGroup updates the radio group for the radio action object.
//
// Parameters:
// 	group	a list representing a radio group
func (r *CRadioAction) SetGroup(group ActionGroup) {
	if err := r.SetStructProperty(PropertyGroup, group); err != nil {
		r.LogErr(err)
	}
}

// GetCurrentValue returns the value property of the currently active member of
// the group to which action belongs.
//
// Returns:
// 	The value of the currently active group member
func (r *CRadioAction) GetCurrentValue() (value int) {
	var err error
	if value, err = r.GetIntProperty(PropertyCurrentValue); err != nil {
		r.LogErr(err)
	}
	return
}

// SetCurrentValue updates the currently active group member to the member with
// value property current_value.
//
// Parameters:
// 	currentValue	the new value
func (r *CRadioAction) SetCurrentValue(currentValue int) {
	if err := r.SetIntProperty(PropertyCurrentValue, currentValue); err != nil {
		r.LogErr(err)
	}
}

// The value property of the currently active member of the group to which
// this action belongs.
// Flags: Read / Write
// Default value: 0
const PropertyCurrentValue cdk.Property = "current-value"

// Sets a new group for a radio action.
// Flags: Write
const PropertyGroup cdk.Property = "group"

// The value is an arbitrary integer which can be used as a convenient way to
// determine which action in the group is currently active in an ::activate
// or ::changed signal handler. See GetCurrentValue and
// RadioActionEntry for convenient ways to get and set this property.
// Flags: Read / Write
// Default value: 0
const PropertyRadioActionValue cdk.Property = "value"

// The ::changed signal is emitted on every member of a radio group when the
// active member is changed. The signal gets emitted after the ::activate
// signals for the previous and current active members.
const SignalRadioActionChanged cdk.Signal = "changed"
