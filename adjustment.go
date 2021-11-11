package ctk

import (
	"github.com/go-curses/cdk"
	"github.com/go-curses/cdk/lib/enums"
)

// CDK type-tag for Adjustment objects
const TypeAdjustment cdk.CTypeTag = "ctk-adjustment"

func init() {
	_ = cdk.TypesManager.AddType(TypeAdjustment, func() interface{} { return MakeAdjustment() })
}

// Adjustment Hierarchy:
//	Object
//	  +- Adjustment
// The Adjustment CTK object is a means of managing the state of multiple
// widgets concurrently. By sharing the same Adjustment instance, one or more
// widgets can ensure that all related User Interface elements are reflecting
// the same values and constraints. The Adjustment consists of an integer value,
// with an upper and lower bounds, and pagination rendering features
type Adjustment interface {
	Object

	Init() bool
	GetValue() (value int)
	SetValue(value int)
	ClampPage(upper, lower int)
	Changed() enums.EventFlag
	ValueChanged() enums.EventFlag
	Settings() (value, lower, upper, stepIncrement, pageIncrement, pageSize int)
	Configure(value, lower, upper, stepIncrement, pageIncrement, pageSize int)
	GetLower() (value int)
	SetLower(lower int)
	GetUpper() (upper int)
	SetUpper(upper int)
	GetStepIncrement() (stepIncrement int)
	SetStepIncrement(stepIncrement int)
	GetPageIncrement() (pageIncrement int)
	SetPageIncrement(pageIncrement int)
	GetPageSize() (pageSize int)
	SetPageSize(pageSize int)
	Moot() bool
	ShowByPolicy(policy PolicyType) bool
}

// The CAdjustment structure implements the Adjustment interface and is
// exported to facilitate type embedding with custom implementations. No member
// variables are exported as the interface methods are the only intended means
// of interacting with Adjustment objects
type CAdjustment struct {
	CObject
}

func MakeAdjustment() *CAdjustment {
	return NewAdjustment(0, 0, 0, 0, 0, 0)
}

// Factory method for convenient construction of new Adjustment objects
func NewAdjustment(value, lower, upper, stepIncrement, pageIncrement, pageSize int) *CAdjustment {
	a := new(CAdjustment)
	a.Init()
	a.Configure(value, lower, upper, stepIncrement, pageIncrement, pageSize)
	return a
}

// Adjustment object initialization. This must be called at least once to setup
// the necessary defaults and allocate any memory structures. Calling this more
// than once is safe though unnecessary. Only the first call will result in any
// effect upon the Adjustment instance
func (a *CAdjustment) Init() bool {
	if a.InitTypeItem(TypeAdjustment, a) {
		return true
	}
	a.CObject.Init()
	_ = a.InstallBuildableProperty(PropertyLower, cdk.IntProperty, true, 0)
	_ = a.InstallBuildableProperty(PropertyPageIncrement, cdk.IntProperty, true, 0)
	_ = a.InstallBuildableProperty(PropertyPageSize, cdk.IntProperty, true, 0)
	_ = a.InstallBuildableProperty(PropertyStepIncrement, cdk.IntProperty, true, 0)
	_ = a.InstallBuildableProperty(PropertyUpper, cdk.IntProperty, true, 0)
	_ = a.InstallBuildableProperty(PropertyValue, cdk.IntProperty, true, 0)
	return false
}

// Gets the current value of the adjustment. See SetValue.
// Returns:
// 	The current value of the adjustment.
func (a *CAdjustment) GetValue() (value int) {
	var err error
	if value, err = a.GetIntProperty(PropertyValue); err != nil {
		a.LogErr(err)
	}
	return
}

// Sets the current value of the Adjustment. This method emits a set-value
// signal initially and if the listeners return an EVENT_PASS then the new value
// is applied and a call to ValueChanged() is made
//
// Emits: SignalSetValue, Argv=[Adjustment instance, current value, new value]
func (a *CAdjustment) SetValue(value int) {
	if err := a.SetIntProperty(PropertyValue, value); err != nil {
		a.LogErr(err)
	} else {
		a.ValueChanged()
	}
}

// Convenience method to set both the upper and lower bounds without emitting
// multiple changed signals. This method emits a clamp-page signal initially and
// if the listeners return an EVENT_PASS then the new upper and lower bounds are
// applied and a call to Changed() is made
//
// Emits: SignalClampPage, Argv=[Adjustment instance, given upper, given lower]
func (a *CAdjustment) ClampPage(upper, lower int) {
	var up, lo int
	up, _ = a.GetIntProperty(PropertyUpper)
	lo, _ = a.GetIntProperty(PropertyLower)
	if up != upper || lo != lower {
		a.Freeze()
		if up != upper {
			a.SetUpper(upper)
		}
		if lo != lower {
			a.SetLower(lower)
		}
		a.Thaw()
		a.Changed()
	}
}

// Cause a changed signal to be emitted. The changed signal reflects that one
// or more of the configurable aspects have changed, excluding the actual value
// of the Adjustment. See ValueChanged()
//
// Emits: SignalChanged, Argv=[Adjustment instance]
// Returns: emission result EventFlag
func (a *CAdjustment) Changed() enums.EventFlag {
	return a.Emit(SignalChanged, a)
}

// Cause a value-changed signal to be emitted. The value-changed signal reflects
// that the actual value of the Adjustment has changed
//
// Emits: SignalValueChanged, Argv=[Adjustment instance]
// Returns: emission result EventFlag
func (a *CAdjustment) ValueChanged() enums.EventFlag {
	return a.Emit(SignalValueChanged, a)
}

// Convenience method to retrieve all the configurable Adjustment values in one
// statement
func (a *CAdjustment) Settings() (value, lower, upper, stepIncrement, pageIncrement, pageSize int) {
	value, lower, upper = a.GetValue(), a.GetLower(), a.GetUpper()
	stepIncrement, pageIncrement, pageSize = a.GetStepIncrement(), a.GetPageIncrement(), a.GetPageSize()
	return
}

// Set all of the configurable aspects of an Adjustment while emitting only a
// single changed and/or value-changed signal. The method emits a configure
// signal initially and if the listeners return an EVENT_PASS then applies all
// of the changes and calls Changed() and/or ValueChanged() accordingly. The
// same effect of this method can be achieved by passing the changed and
// value-changed signal emissions for the Adjustment instance, making all the
// individual calls to the setter methods, resuming the changed and
// value-changed signals and finally calling the Changed() and/or ValueChanged()
// methods accordingly
//
// Parameters:
// 	value	the new value
// 	lower	the new minimum value
// 	upper	the new maximum value
// 	stepIncrement	the new step increment
// 	pageIncrement	the new page increment
// 	pageSize	the new page size
//
// Emits: SignalConfigure, Argv=[Adjustment instance, value, lower, upper, stepIncrement, pageIncrement, pageSize]
func (a *CAdjustment) Configure(value, lower, upper, stepIncrement, pageIncrement, pageSize int) {
	a.Freeze()
	aValue, aLower, aUpper, aStepIncrement, aPageIncrement, aPageSize := a.Settings()
	valueChanged := aValue != value
	changed := aLower != lower || aUpper != upper ||
		aStepIncrement != stepIncrement ||
		aPageIncrement != pageIncrement ||
		aPageSize != pageSize
	a.SetValue(value)
	a.SetLower(lower)
	a.SetUpper(upper)
	a.SetStepIncrement(stepIncrement)
	a.SetPageIncrement(pageIncrement)
	a.SetPageSize(pageSize)
	a.Thaw()
	if changed {
		a.Changed()
	}
	if valueChanged {
		a.ValueChanged()
	}
}

// Returns the current lower bounds of the Adjustment
func (a *CAdjustment) GetLower() (value int) {
	var err error
	if value, err = a.GetIntProperty(PropertyLower); err != nil {
		a.LogErr(err)
	}
	return
}

// Sets the lower bounds of the Adjustment. This method emits a set-lower signal
// initially and if the listeners return an EVENT_PASS then the value is applied
// and a call to Changed() is made
func (a *CAdjustment) SetLower(lower int) {
	if err := a.SetIntProperty(PropertyLower, lower); err != nil {
		a.LogErr(err)
	} else {
		a.Changed()
	}
}

// Returns the current lower bounds of the Adjustment
func (a *CAdjustment) GetUpper() (upper int) {
	var err error
	if upper, err = a.GetIntProperty(PropertyUpper); err != nil {
		a.LogErr(err)
	}
	return
}

// Sets the upper bounds of the Adjustment. This method emits a set-upper signal
// initially and if the listeners return an EVENT_PASS then the value is applied
// and a call to Changed() is made
func (a *CAdjustment) SetUpper(upper int) {
	if err := a.SetIntProperty(PropertyUpper, upper); err != nil {
		a.LogErr(err)
	} else {
		a.Changed()
	}
}

// Returns the current step increment of the Adjustment. Adjustment values are
// intended to be increased or decreased by either a step or page amount. The
// step increment is the shorter movement such as moving up or down a line of
// text
func (a *CAdjustment) GetStepIncrement() (stepIncrement int) {
	var err error
	if stepIncrement, err = a.GetIntProperty(PropertyStepIncrement); err != nil {
		a.LogErr(err)
	}
	return
}

// Sets the step increment of the Adjustment. This method emits a set-step
// signal initially and if the listeners return an EVENT_PASS then the value is
// applied and a call to Changed() is made
func (a *CAdjustment) SetStepIncrement(stepIncrement int) {
	if err := a.SetIntProperty(PropertyStepIncrement, stepIncrement); err != nil {
		a.LogErr(err)
	} else {
		a.Changed()
	}
}

// Returns the current page increment of the Adjustment. Adjustment values are
// intended to be increased or decreased by either a step or page amount. The
// page increment is the longer movement such as moving up or down half a page
// of text
func (a *CAdjustment) GetPageIncrement() (pageIncrement int) {
	var err error
	if pageIncrement, err = a.GetIntProperty(PropertyPageIncrement); err != nil {
		a.LogErr(err)
	}
	return
}

// Sets the page increment of the Adjustment. This method emits a set-step
// signal initially and if the listeners return an EVENT_PASS then the value is
// applied and a call to Changed() is made
func (a *CAdjustment) SetPageIncrement(pageIncrement int) {
	if err := a.SetIntProperty(PropertyPageIncrement, pageIncrement); err != nil {
		a.LogErr(err)
	} else {
		a.Changed()
	}
}

// Gets the page size of the Adjustment. Adjustment values are intended to be
// increased or decreased by either a step or page amount and having a separate
// Adjustment variable to track the end-user facing page size is beneficial.
// This value does not have to be the same as the page increment, however the
// page increment is in effect clamped to the page size of the Adjustment
func (a *CAdjustment) GetPageSize() (pageSize int) {
	var err error
	if pageSize, err = a.GetIntProperty(PropertyPageSize); err != nil {
		a.LogErr(err)
	}
	return
}

// Sets the page size of the Adjustment. This method emits a set-page-size
// signal initially and if the listeners return an EVENT_PASS then the value is
// applied and a call to Changed() is made
func (a *CAdjustment) SetPageSize(pageSize int) {
	if err := a.SetIntProperty(PropertyPageSize, pageSize); err != nil {
		a.LogErr(err)
	} else {
		a.Changed()
	}
}

// This method determines if the Adjustment object is in a "moot" state in which
// the values are irrelevant because the upper bounds is set to zero. Given an
// upper bounds of zero, the value and lower bounds are also clamped to zero and
// so on. Thus, this convenience method enables a human readable understanding
// of why the upper bounds is being checked for a zero value (eliminates a
// useful magic value)
func (a *CAdjustment) Moot() bool {
	return a.GetUpper() == 0 || a.GetLower() == a.GetUpper()
}

// Convenience method that given a Policy, determine if the Adjustment should
// be rendered or otherwise used to an end-user-facing effect. This method is
// primarily used by other Widgets as a convenient way to determine if they
// should show or hide their presence
func (a *CAdjustment) ShowByPolicy(policy PolicyType) bool {
	switch policy {
	case PolicyNever:
		return false
	case PolicyAutomatic:
		return !a.Moot()
	case PolicyAlways:
		return true
	}
	a.LogError("unknown policy given: %v", policy)
	return false
}

// The minimum value of the adjustment.
// Flags: Read / Write
// Default value: 0
const PropertyLower cdk.Property = "lower"

// The page increment of the adjustment.
// Flags: Read / Write
// Default value: 0
const PropertyPageIncrement cdk.Property = "page-increment"

// The page size of the adjustment. Note that the page-size is irrelevant and
// should be set to zero if the adjustment is used for a simple scalar value,
// e.g. in a SpinButton.
// Flags: Read / Write
// Default value: 0
const PropertyPageSize cdk.Property = "page-size"

// The step increment of the adjustment.
// Flags: Read / Write
// Default value: 0
const PropertyStepIncrement cdk.Property = "step-increment"

// The maximum value of the adjustment. Note that values will be restricted
// by upper - page-size if the page-size property is nonzero.
// Flags: Read / Write
// Default value: 0
const PropertyUpper cdk.Property = "upper"

// The value of the adjustment.
// Flags: Read / Write
// Default value: 0
const PropertyValue cdk.Property = "value"

const SignalChanged cdk.Signal = "changed"

const SignalValueChanged cdk.Signal = "value-changed"
