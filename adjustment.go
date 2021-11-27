package ctk

import (
	"github.com/go-curses/cdk"
	cenums "github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/ctk/lib/enums"
)

const TypeAdjustment cdk.CTypeTag = "ctk-adjustment"

func init() {
	_ = cdk.TypesManager.AddType(TypeAdjustment, func() interface{} { return MakeAdjustment() })
}

// Adjustment Hierarchy:
//	Object
//	  +- Adjustment
//
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
	Changed() cenums.EventFlag
	ValueChanged() cenums.EventFlag
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
	ShowByPolicy(policy enums.PolicyType) bool
}

// The CAdjustment structure implements the Adjustment interface and is
// exported to facilitate type embedding with custom implementations. No member
// variables are exported as the interface methods are the only intended means
// of interacting with Adjustment objects.
type CAdjustment struct {
	CObject
}

// MakeAdjustment is used by the Buildable system to construct a new Adjustment.
func MakeAdjustment() Adjustment {
	return NewAdjustment(0, 0, 0, 0, 0, 0)
}

// NewAdjustment is the constructor for new Adjustment instances.
func NewAdjustment(value, lower, upper, stepIncrement, pageIncrement, pageSize int) Adjustment {
	a := new(CAdjustment)
	a.Init()
	a.Configure(value, lower, upper, stepIncrement, pageIncrement, pageSize)
	return a
}

// Init initializes an Adjustment object. This must be called at least once to
// set up the necessary defaults and allocate any memory structures. Calling
// this more than once is safe though unnecessary. Only the first call will
// result in any effect upon the Adjustment instance. Init is used in the
// NewAdjustment constructor and only necessary when implementing a derivative
// Adjustment type.
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

// GetValue returns the current value of the adjustment.
// See: SetValue()
//
// Locking: read
func (a *CAdjustment) GetValue() (value int) {
	var err error
	if value, err = a.GetIntProperty(PropertyValue); err != nil {
		a.LogErr(err)
	}
	return
}

// SetValue updates the current value of the Adjustment. This method emits a
// set-value signal initially and if the listeners return an EVENT_PASS then the
// new value is applied and a call to ValueChanged() is made.
//
// Locking: write
func (a *CAdjustment) SetValue(value int) {
	if f := a.Emit(SignalSetValue, value); f == cenums.EVENT_PASS {
		if err := a.SetIntProperty(PropertyValue, value); err != nil {
			a.LogErr(err)
		} else {
			a.ValueChanged()
		}
	}
}

// ClampPage is a convenience method to set both the upper and lower bounds
// without emitting multiple changed signals. This method emits a clamp-page
// signal initially and if the listeners return an EVENT_PASS then the new upper
// and lower bounds are applied and a call to Changed() is made.
//
// Locking: write
func (a *CAdjustment) ClampPage(upper, lower int) {
	var up, lo int
	up, _ = a.GetIntProperty(PropertyUpper)
	lo, _ = a.GetIntProperty(PropertyLower)
	if f := a.Emit(SignalClampPage, upper, lower); f == cenums.EVENT_PASS {
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
}

// Changed emits a changed signal. The changed signal reflects that one or more
// of the configurable aspects have changed, excluding the actual value of the
// Adjustment.
// See: ValueChanged()
func (a *CAdjustment) Changed() cenums.EventFlag {
	return a.Emit(SignalChanged, a)
}

// ValueChanged emits a value-changed signal. The value-changed signal reflects
// that the actual value of the Adjustment has changed.
func (a *CAdjustment) ValueChanged() cenums.EventFlag {
	return a.Emit(SignalValueChanged, a)
}

// Settings is a convenience method to retrieve all the configurable Adjustment
// values in one statement.
//
// Locking: read
func (a *CAdjustment) Settings() (value, lower, upper, stepIncrement, pageIncrement, pageSize int) {
	value, lower, upper = a.GetValue(), a.GetLower(), a.GetUpper()
	stepIncrement, pageIncrement, pageSize = a.GetStepIncrement(), a.GetPageIncrement(), a.GetPageSize()
	return
}

// Configure updates all the configurable aspects of an Adjustment while
// emitting only a single changed and/or value-changed signal. The method emits
// a configure signal initially and if the listeners return an EVENT_PASS then
// applies all the changes and calls Changed() and/or ValueChanged()
// accordingly. The same effect of this method can be achieved by passing the
// changed and value-changed signal emissions for the Adjustment instance,
// making all the individual calls to the setter methods, resuming the changed
// and value-changed signals and finally calling the Changed() and/or
// ValueChanged() methods accordingly.
//
// Parameters:
// 	value	the new value
// 	lower	the new minimum value
// 	upper	the new maximum value
// 	stepIncrement	the new step increment
// 	pageIncrement	the new page increment
// 	pageSize	the new page size
//
// Locking: write
func (a *CAdjustment) Configure(value, lower, upper, stepIncrement, pageIncrement, pageSize int) {
	if f := a.Emit(SignalConfigure, value, lower, upper, stepIncrement, pageIncrement, pageSize); f == cenums.EVENT_PASS {
		a.Freeze()
		aValue, aLower, aUpper, aStepIncrement, aPageIncrement, aPageSize := a.Settings()
		valueChanged := aValue != value
		changed := aLower != lower ||
			aUpper != upper ||
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
}

// GetLower returns the current lower bounds of the Adjustment.
//
// Locking: read
func (a *CAdjustment) GetLower() (value int) {
	var err error
	if value, err = a.GetIntProperty(PropertyLower); err != nil {
		a.LogErr(err)
	}
	return
}

// SetLower updates the lower bounds of the Adjustment. This method emits a
// set-lower signal initially and if the listeners return an EVENT_PASS then the
// value is applied and a call to Changed() is made.
//
// Locking: write
func (a *CAdjustment) SetLower(lower int) {
	if f := a.Emit(SignalSetLower, lower); f == cenums.EVENT_PASS {
		if err := a.SetIntProperty(PropertyLower, lower); err != nil {
			a.LogErr(err)
		} else {
			a.Changed()
		}
	}
}

// GetUpper returns the current lower bounds of the Adjustment.
//
// Locking: read
func (a *CAdjustment) GetUpper() (upper int) {
	var err error
	if upper, err = a.GetIntProperty(PropertyUpper); err != nil {
		a.LogErr(err)
	}
	return
}

// SetUpper updates the upper bounds of the Adjustment. This method emits a
// set-upper signal initially and if the listeners return an EVENT_PASS then the
// value is applied and a call to Changed() is made.
//
// Locking: write
func (a *CAdjustment) SetUpper(upper int) {
	if f := a.Emit(SignalSetUpper, upper); f == cenums.EVENT_PASS {
		if err := a.SetIntProperty(PropertyUpper, upper); err != nil {
			a.LogErr(err)
		} else {
			a.Changed()
		}
	}
}

// GetStepIncrement returns the current step increment of the Adjustment.
// Adjustment values are intended to be increased or decreased by either a step
// or page amount. The step increment is the shorter movement such as moving up
// or down a line of text.
//
// Locking: read
func (a *CAdjustment) GetStepIncrement() (stepIncrement int) {
	var err error
	if stepIncrement, err = a.GetIntProperty(PropertyStepIncrement); err != nil {
		a.LogErr(err)
	}
	return
}

// SetStepIncrement updates the step increment of the Adjustment. This method
// emits a set-step-increment signal initially and if the listeners return an
// EVENT_PASS then the value is applied and a call to Changed() is made.
//
// Locking: write
func (a *CAdjustment) SetStepIncrement(stepIncrement int) {
	if f := a.Emit(SignalSetStepIncrement, stepIncrement); f == cenums.EVENT_PASS {
		if err := a.SetIntProperty(PropertyStepIncrement, stepIncrement); err != nil {
			a.LogErr(err)
		} else {
			a.Changed()
		}
	}
}

// GetPageIncrement returns the current page increment of the Adjustment.
// Adjustment values are intended to be increased or decreased by either a step
// or page amount. The page increment is the longer movement such as moving up
// or down half a page of text.
//
// Locking: read
func (a *CAdjustment) GetPageIncrement() (pageIncrement int) {
	var err error
	if pageIncrement, err = a.GetIntProperty(PropertyPageIncrement); err != nil {
		a.LogErr(err)
	}
	return
}

// SetPageIncrement updates the page increment of the Adjustment. This method
// emits a set-page-increment signal initially and if the listeners return an
// EVENT_PASS then the value is applied and a call to Changed() is made.
//
// Locking: write
func (a *CAdjustment) SetPageIncrement(pageIncrement int) {
	if f := a.Emit(SignalSetPageIncrement, pageIncrement); f == cenums.EVENT_PASS {
		if err := a.SetIntProperty(PropertyPageIncrement, pageIncrement); err != nil {
			a.LogErr(err)
		} else {
			a.Changed()
		}
	}
}

// GetPageSize returns the page size of the Adjustment. Adjustment values are
// intended to be increased or decreased by either a step or page amount and
// having a separate Adjustment variable to track the end-user facing page size
// is beneficial. This value does not have to be the same as the page increment,
// however the page increment is in effect clamped to the page size of the
// Adjustment.
//
// Locking: read
func (a *CAdjustment) GetPageSize() (pageSize int) {
	var err error
	if pageSize, err = a.GetIntProperty(PropertyPageSize); err != nil {
		a.LogErr(err)
	}
	return
}

// SetPageSize updates the page size of the Adjustment. This method emits a
// set-page-size signal initially and if the listeners return an EVENT_PASS then the value is
// applied and a call to Changed() is made.
//
// Locking: write
func (a *CAdjustment) SetPageSize(pageSize int) {
	if f := a.Emit(SignalSetPageSize, pageSize); f == cenums.EVENT_PASS {
		if err := a.SetIntProperty(PropertyPageSize, pageSize); err != nil {
			a.LogErr(err)
		} else {
			a.Changed()
		}
	}
}

// Moot is a convenience method to return TRUE if the Adjustment object is in a
// "moot" state in which the values are irrelevant because the upper bounds is
// set to zero. Given an upper bounds of zero, the value and lower bounds are
// also clamped to zero and so on. Thus, this convenience method enables a Human
// readable understanding of why the upper bounds is being checked for a zero
// value (eliminates a useful magic value).
//
// Locking: read
func (a *CAdjustment) Moot() bool {
	return a.GetUpper() == 0 || a.GetLower() == a.GetUpper()
}

// ShowByPolicy is a convenience method that given a PolicyType, determines if
// the Widget using the Adjustment should be rendered or otherwise used to an
// end-user-facing effect. This method is primarily used by other Widgets as a
// convenient way to determine if they should show or hide their presence.
//
// Locking: read
func (a *CAdjustment) ShowByPolicy(policy enums.PolicyType) bool {
	switch policy {
	case enums.PolicyNever:
		return false
	case enums.PolicyAutomatic:
		return !a.Moot()
	case enums.PolicyAlways:
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

const SignalSetValue cdk.Signal = "set-value"

const SignalClampPage cdk.Signal = "clamp-page"

const SignalConfigure cdk.Signal = "configure"

const SignalSetLower cdk.Signal = "set-lower"

const SignalSetUpper cdk.Signal = "set-upper"

const SignalSetStepIncrement cdk.Signal = "set-step-increment"

const SignalSetPageIncrement cdk.Signal = "set-page-increment"

const SignalSetPageSize cdk.Signal = "set-page-size"

const SignalValueChanged cdk.Signal = "value-changed"
