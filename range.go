package ctk

import (
	"github.com/go-curses/cdk"
	cmath "github.com/go-curses/cdk/lib/math"
	"github.com/go-curses/cdk/lib/ptypes"
)

const TypeRange cdk.CTypeTag = "ctk-range"

func init() {
	_ = cdk.TypesManager.AddType(TypeRange, nil)
}

// Range Hierarchy:
//	Object
//	  +- Widget
//	    +- Range
//	      +- Scale
//	      +- Scrollbar
//
// The Range Widget is used to manage the position of things within some range
// of values. Scrollbar and Scale are two examples.
type Range interface {
	Widget

	Init() (already bool)
	GetFillLevel() (value float64)
	GetRestrictToFillLevel() (value bool)
	GetShowFillLevel() (value bool)
	SetFillLevel(fillLevel float64)
	SetRestrictToFillLevel(restrictToFillLevel bool)
	SetShowFillLevel(showFillLevel bool)
	GetAdjustment() (adjustment *CAdjustment)
	GetInverted() (value bool)
	SetInverted(setting bool)
	SetIncrements(step int, page int)
	SetRange(min, max int)
	GetValue() (value int)
	SetValue(value int)
	GetRoundDigits() (value int)
	SetRoundDigits(roundDigits int)
	SetLowerStepperSensitivity(sensitivity SensitivityType)
	GetLowerStepperSensitivity() (value SensitivityType)
	SetUpperStepperSensitivity(sensitivity SensitivityType)
	GetUpperStepperSensitivity() (value SensitivityType)
	GetFlippable() (value bool)
	SetFlippable(flippable bool)
	GetMinSliderSize() (value int)
	GetRangeRect() (rangeRect ptypes.Rectangle)
	GetSliderRange() (sliderStart int, sliderEnd int)
	GetSliderSizeFixed() (value bool)
	SetMinSliderSize(minSize bool)
	SetSliderSizeFixed(sizeFixed bool)
	GetIncrements() (step int, page int)
	GetRange() (min, max int)
	GetMinSliderLength() (length int)
	SetMinSliderLength(length int)
	GetSliderLength() (length int)
	SetSliderLength(length int)
	GetStepperSize() (size int)
	SetStepperSize(size int)
	GetStepperSpacing() (spacing int)
	SetStepperSpacing(spacing int)
	GetTroughUnderSteppers() (underSteppers bool)
	SetTroughUnderSteppers(underSteppers bool)
}

// The CRange structure implements the Range interface and is exported
// to facilitate type embedding with custom implementations. No member variables
// are exported as the interface methods are the only intended means of
// interacting with Range objects.
type CRange struct {
	CWidget

	restrictToFillLevel bool
	flippable           bool
	minSliderLength     int
	sliderSizeFixed     bool
	sliderLength        int
	stepperSize         int
	stepperSpacing      int
	troughUnderSteppers bool
}

// Init initializes a Range object. This must be called at least once to
// set up the necessary defaults and allocate any memory structures. Calling
// this more than once is safe though unnecessary. Only the first call will
// result in any effect upon the Range instance. Init is used in the
// NewRange constructor and only necessary when implementing a derivative
// Range type.
func (r *CRange) Init() (already bool) {
	if r.InitTypeItem(TypeRange, r) {
		return true
	}
	r.CWidget.Init()
	r.flags = NULL_WIDGET_FLAG
	r.SetFlags(SENSITIVE | PARENT_SENSITIVE | APP_PAINTABLE)
	_ = r.InstallProperty(PropertyAdjustment, cdk.StructProperty, true, NewAdjustment(0, 0, 0, 0, 0, 0))
	_ = r.InstallProperty(PropertyFillLevel, cdk.FloatProperty, true, 1.0)
	_ = r.InstallProperty(PropertyInverted, cdk.BoolProperty, true, false)
	_ = r.InstallProperty(PropertyLowerStepperSensitivity, cdk.StructProperty, true, SensitivityAuto)
	_ = r.InstallProperty(PropertyRestrictToFillLevel, cdk.BoolProperty, true, false)
	_ = r.InstallProperty(PropertyRoundDigits, cdk.IntProperty, true, -1)
	_ = r.InstallProperty(PropertyShowFillLevel, cdk.BoolProperty, true, false)
	_ = r.InstallProperty(PropertyUpdatePolicy, cdk.StructProperty, true, UpdateContinuous)
	_ = r.InstallProperty(PropertyUpperStepperSensitivity, cdk.StructProperty, true, SensitivityAuto)
	r.restrictToFillLevel = false
	r.troughUnderSteppers = false
	r.flippable = false
	r.minSliderLength = 1
	r.sliderLength = -1
	r.stepperSize = -1
	r.stepperSpacing = 0
	return false
}

// GetFillLevel returns the current position of the fill level indicator.
func (r *CRange) GetFillLevel() (value float64) {
	var err error
	if value, err = r.GetFloat64Property(PropertyFillLevel); err != nil {
		r.LogErr(err)
	}
	return
}

// GetRestricttoFillLevel returns whether the range is restricted to the fill
// level.
func (r *CRange) GetRestrictToFillLevel() (value bool) {
	var err error
	if value, err = r.GetBoolProperty(PropertyRestrictToFillLevel); err != nil {
		r.LogErr(err)
	}
	return
}

// GetShowFillLevel returns whether the range displays the fill level
// graphically.
func (r *CRange) GetShowFillLevel() (value bool) {
	var err error
	if value, err = r.GetBoolProperty(PropertyShowFillLevel); err != nil {
		r.LogErr(err)
	}
	return
}

// SetFillLevel updates the new position of the fill level indicator. The
// "fill level" is probably best described by its most prominent use case, which
// is an indicator for the amount of pre-buffering in a streaming media player.
// In that use case, the value of the range would indicate the current play
// position, and the fill level would be the position up to which the
// file/stream has been downloaded. This amount of pre-buffering can be
// displayed on the range's trough and is themeable separately from the trough.
// To enable fill level display, use SetShowFillLevel. The range defaults to not
// showing the fill level. Additionally, it's possible to restrict the range's
// slider position to values which area smaller than the fill level. This is
// controller by SetRestrictToFillLevel and is by default enabled.
//
// Parameters:
// 	fillLevel	the new position of the fill level indicator
func (r *CRange) SetFillLevel(fillLevel float64) {
	if err := r.SetFloatProperty(PropertyFillLevel, cmath.ClampF(fillLevel, 0.0, 1.0)); err != nil {
		r.LogErr(err)
	}
}

// SetRestrictToFillLevel updates whether the slider is restricted to the fill
// level.
// See: SetFillLevel()
//
// Parameters:
// 	restrictToFillLevel	Whether the fill level restricts slider movement.
func (r *CRange) SetRestrictToFillLevel(restrictToFillLevel bool) {
	if err := r.SetBoolProperty(PropertyRestrictToFillLevel, restrictToFillLevel); err != nil {
		r.LogErr(err)
	}
}

// SetShowFillLevel updates whether a graphical fill level is show on the
// trough.
// See: SetFillLevel()
//
// Parameters:
// 	showFillLevel	Whether a fill level indicator graphics is shown.
func (r *CRange) SetShowFillLevel(showFillLevel bool) {
	if err := r.SetBoolProperty(PropertyShowFillLevel, showFillLevel); err != nil {
		r.LogErr(err)
	}
}

// GetAdjustment returns the Adjustment which is the "model" object for Range.
// See: SetAdjustment()
func (r *CRange) GetAdjustment() (adjustment *CAdjustment) {
	if v, err := r.GetStructProperty(PropertyAdjustment); err != nil {
		r.LogErr(err)
	} else {
		var ok bool
		if adjustment, ok = v.(*CAdjustment); !ok {
			r.LogError("value stored in %v property is not of *CAdjustment type: %v (%T)", PropertyAdjustment, v, v)
		}
	}
	return
}

// GetInverted returns the value set by SetInverted.
func (r *CRange) GetInverted() (value bool) {
	var err error
	if value, err = r.GetBoolProperty(PropertyInverted); err != nil {
		r.LogErr(err)
	}
	return
}

// SetInverted updates the inverted property value. Ranges normally move from
// lower to higher values as the slider moves from top to bottom or left to
// right. Inverted ranges have higher values at the top or on the right rather
// than on the bottom or left.
//
// Parameters:
// 	setting	TRUE to invert the range
func (r *CRange) SetInverted(setting bool) {
	if err := r.SetBoolProperty(PropertyInverted, setting); err != nil {
		r.LogErr(err)
	}
}

// SetIncrements updates the step and page sizes for the range. The step size is
// used when the user clicks the Scrollbar arrows or moves Scale via arrow keys.
// The page size is used for example when moving via Page Up or Page Down keys.
//
// Parameters:
// 	step	step size
// 	page	page size
func (r *CRange) SetIncrements(step int, page int) {
	if adjustment := r.GetAdjustment(); adjustment != nil {
		adjustment.SetStepIncrement(step)
		adjustment.SetPageIncrement(page)
	} else {
		r.LogError("missing adjustment")
	}
}

// SetRange updates the allowable values in the Range, and clamps the range
// value to be between min and max. (If the range has a non-zero page size, it
// is clamped between min and max - page-size.)
//
// Parameters:
// 	min	minimum range value
// 	max	maximum range value
func (r *CRange) SetRange(min, max int) {
	if adjustment := r.GetAdjustment(); adjustment != nil {
		adjustment.SetLower(min)
		adjustment.SetUpper(max)
	} else {
		r.LogError("missing adjustment")
	}
}

// GetValue returns the current value of the range.
func (r *CRange) GetValue() (value int) {
	value = -1
	if adjustment := r.GetAdjustment(); adjustment != nil {
		value = adjustment.GetValue()
	} else {
		r.LogError("missing adjustment")
	}
	return
}

// SetValue updates the current value of the range; if the value is outside the
// minimum or maximum range values, it will be clamped to fit inside them. The
// range emits the value-changed signal if the value changes.
//
// Parameters:
// 	value	new value of the range
func (r *CRange) SetValue(value int) {
	if adjustment := r.GetAdjustment(); adjustment != nil {
		if r.GetRestrictToFillLevel() {
			// 0.0 == lower, 1.0 == upper
			max := int(r.GetFillLevel() * float64(adjustment.GetUpper()))
			value = cmath.ClampI(value, 0, max)
		}
		value = cmath.ClampI(value, adjustment.GetLower(), adjustment.GetUpper())
		previousValue := adjustment.GetValue()
		if previousValue != value {
			adjustment.SetValue(value)
			r.Emit(SignalRangeValueChanged, r, previousValue, value)
		}
	}
}

// GetRoundDigits returns the number of digits to round the value to when it
// changes.
func (r *CRange) GetRoundDigits() (value int) {
	var err error
	if value, err = r.GetIntProperty(PropertyRoundDigits); err != nil {
		r.LogErr(err)
	}
	return
}

// SetRoundDigits updates the number of digits to round the value to when it
// changes.
//
// Parameters:
// 	roundDigits	the precision in digits, or -1
func (r *CRange) SetRoundDigits(roundDigits int) {
	if err := r.SetIntProperty(PropertyRoundDigits, roundDigits); err != nil {
		r.LogErr(err)
	}
}

// GetLowerStepperSensitivity updates the sensitivity policy for the stepper
// that points to the 'lower' end of the Range's adjustment.
func (r *CRange) GetLowerStepperSensitivity() (value SensitivityType) {
	var ok bool
	if v, err := r.GetStructProperty(PropertyLowerStepperSensitivity); err != nil {
		r.LogErr(err)
	} else if value, ok = v.(SensitivityType); !ok {
		r.LogError("value stored in %v property is not of SensitivityType type: %v (%T)", PropertyLowerStepperSensitivity, v, v)
	}
	return
}

// SetLowerStepperSensitivity updates the sensitivity policy for the stepper
// that points to the 'lower' end of the Range's adjustment.
//
// Parameters:
// 	sensitivity	the lower stepper's sensitivity policy.
func (r *CRange) SetLowerStepperSensitivity(sensitivity SensitivityType) {
	if err := r.SetStructProperty(PropertyLowerStepperSensitivity, sensitivity); err != nil {
		r.LogErr(err)
	}
}

// GetUpperStepperSensitivity updates the sensitivity policy for the stepper
// that points to the 'upper' end of the Range's adjustment.
func (r *CRange) GetUpperStepperSensitivity() (value SensitivityType) {
	var ok bool
	if v, err := r.GetStructProperty(PropertyUpperStepperSensitivity); err != nil {
		r.LogErr(err)
	} else if value, ok = v.(SensitivityType); !ok {
		r.LogError("value stored in %v property is not of SensitivityType type: %v (%T)", PropertyUpperStepperSensitivity, v, v)
	}
	return
}

// SetUpperStepperSensitivity updates the sensitivity policy for the stepper
// that points to the 'upper' end of the Range's adjustment.
//
// Parameters:
// 	sensitivity	the upper stepper's sensitivity policy.
func (r *CRange) SetUpperStepperSensitivity(sensitivity SensitivityType) {
	if err := r.SetStructProperty(PropertyUpperStepperSensitivity, sensitivity); err != nil {
		r.LogErr(err)
	}
}

// GetFlippable returns the value set by SetFlippable.
func (r *CRange) GetFlippable() (value bool) {
	return r.flippable
}

// SetFlippable updates whether a range is flippable. If the range is flippable,
// it will switch its direction if it is horizontal and its direction is
// TEXT_DIR_RTL.
// See: Widget.GetTextDirection()
//
// Parameters:
// 	flippable	TRUE to make the range flippable
//
// Note that usage of this within CTK is unimplemented at this time
func (r *CRange) SetFlippable(flippable bool) {
	r.flippable = flippable
}

// GetMinSliderSize returns the minimum slider size. This method is useful
// mainly for Range subclasses.
// See: SetMinSliderSize()
//
// Note that usage of this within CTK is unimplemented at this time
func (r *CRange) GetMinSliderSize() (value int) {
	return 0
}

// GetRangeRect returns the area that contains the range's trough and its
// steppers, in widget->window coordinates. This function is useful mainly
// for Range subclasses.
//
// Note that usage of this within CTK is unimplemented at this time
func (r *CRange) GetRangeRect() (rangeRect ptypes.Rectangle) {
	return
}

// GetSliderRange returns sliders range along the long dimension, in
// widget->window coordinates. This function is useful mainly for Range
// subclasses.
//
// Note that usage of this within CTK is unimplemented at this time
func (r *CRange) GetSliderRange() (sliderStart int, sliderEnd int) {
	return
}

// GetSliderSizeFixed returns whether the slider size is fixed or not. This
// method is useful mainly for Range subclasses.
// See: SetSliderSizeFixed()
//
// Note that usage of this within CTK is unimplemented at this time
func (r *CRange) GetSliderSizeFixed() (value bool) {
	return r.sliderSizeFixed
}

// SetMinSliderSize updates the minimum size of the range's slider. This
// method is useful mainly for Range subclasses.
//
// Parameters:
// 	minSize	The slider's minimum size
//
// Note that usage of this within CTK is unimplemented at this time
func (r *CRange) SetMinSliderSize(minSize bool) {}

// SetSliderSizeFixed updates whether the range's slider has a fixed size, or a
// size that depends on it's adjustment's page size. This function is useful
// mainly for Range subclasses.
//
// Parameters:
// 	sizeFixed	TRUE to make the slider size constant
//
// Note that usage of this within CTK is unimplemented at this time
func (r *CRange) SetSliderSizeFixed(sizeFixed bool) {
	r.sliderSizeFixed = sizeFixed
}

// GetIncrements returns the step and page sizes for the range. The step size is
// used when the user clicks the Scrollbar arrows or moves Scale via arrow keys.
// The page size is used for example when moving via Page Up or Page Down keys.
func (r *CRange) GetIncrements() (step int, page int) {
	step, page = -1, -1
	if adjustment := r.GetAdjustment(); adjustment != nil {
		step, page = adjustment.GetStepIncrement(), adjustment.GetPageIncrement()
	} else {
		r.LogError("missing adjustment")
	}
	return
}

// GetRange returns the allowable values in the Range.
func (r *CRange) GetRange() (min, max int) {
	min, max = -1, -1
	if adjustment := r.GetAdjustment(); adjustment != nil {
		min, max = adjustment.GetLower(), adjustment.GetUpper()
	} else {
		r.LogError("missing adjustment")
	}
	return
}

// GetMinSliderLength returns the minimum slider length. This method is useful
// mainly for Range subclasses.
// See: SetMinSliderLength()
func (r *CRange) GetMinSliderLength() (length int) {
	length = r.minSliderLength
	if r.minSliderLength <= -1 {
		length = 1
	}
	return
}

// SetMinSliderLength updates the minimum size of the range's slider. This
// method is useful mainly for Range subclasses.
//
// Parameters:
// 	minSize	The slider's minimum size
func (r *CRange) SetMinSliderLength(length int) {
	r.minSliderLength = length
}

// GetSliderLength returns the length of the scrollbar or scale thumb.
func (r *CRange) GetSliderLength() (length int) {
	length = r.sliderLength
	return
}

// SetSliderLength updates the length of the scrollbar or scale thumb. Sets
// fixed slider length to true. Set to -1 for variable slider length.
func (r *CRange) SetSliderLength(length int) {
	if length <= -1 {
		r.sliderLength = -1
		r.SetSliderSizeFixed(false)
	} else if r.sliderLength != length {
		r.SetSliderSizeFixed(true)
		r.sliderLength = length
	}
}

// GetStepperSize returns the length of step buttons at ends.
func (r *CRange) GetStepperSize() (size int) {
	return r.stepperSize
}

// SetStepperSize updates the length of step buttons at ends.
func (r *CRange) SetStepperSize(size int) {
	r.stepperSize = size
}

// GetStepperSpacing returns the spacing between the stepper buttons and thumb.
// Note that setting this value to anything > 0 will automatically set the
// trough-under-steppers style property to TRUE as well. Also, stepper-spacing
// won't have any effect if there are no steppers.
func (r *CRange) GetStepperSpacing() (spacing int) {
	return r.stepperSpacing
}

// SetStepperSpacing updates the spacing between the stepper buttons and thumb.
// Note that setting this value to anything > 0 will automatically set the
// trough-under-steppers style property to TRUE as well. Also, stepper-spacing
// won't have any effect if there are no steppers.
func (r *CRange) SetStepperSpacing(spacing int) {
	r.stepperSpacing = spacing
}

// GetTroughUnderSteppers returns whether to draw the trough across the full
// length of the range or to exclude the steppers and their spacing. Note that
// setting the stepper-spacing style property to any value > 0 will
// automatically enable trough-under-steppers too.
func (r *CRange) GetTroughUnderSteppers() (underSteppers bool) {
	return r.troughUnderSteppers
}

// SetTroughUnderSteppers updates whether to draw the trough across the full
// length of the range or to exclude the steppers and their spacing. Note that
// setting the stepper-spacing style property to any value > 0 will
// automatically enable trough-under-steppers too.
func (r *CRange) SetTroughUnderSteppers(underSteppers bool) {
	r.troughUnderSteppers = underSteppers
}

// The Adjustment that contains the current value of this range object.
// Flags: Read / Write / Construct
const PropertyAdjustment cdk.Property = "adjustment"

// The fill level (e.g. prebuffering of a network stream). See
// SetFillLevel.
// Flags: Read / Write
// Default value: 1.79769e+308
const PropertyFillLevel cdk.Property = "fill-level"

// Invert direction slider moves to increase range value.
// Flags: Read / Write
// Default value: FALSE
const PropertyInverted cdk.Property = "inverted"

// The sensitivity policy for the stepper that points to the adjustment's
// lower side.
// Flags: Read / Write
// Default value: GTK_SENSITIVITY_AUTO
const PropertyLowerStepperSensitivity cdk.Property = "lower-stepper-sensitivity"

// The restrict-to-fill-level property controls whether slider movement is
// restricted to an upper boundary set by the fill level. See
// SetRestrictToFillLevel.
// Flags: Read / Write
// Default value: TRUE
const PropertyRestrictToFillLevel cdk.Property = "restrict-to-fill-level"

// The number of digits to round the value to when it changes, or -1. See
// “change-value”.
// Flags: Read / Write
// Allowed values: >= -1
// Default value: -1
const PropertyRoundDigits cdk.Property = "round-digits"

// The show-fill-level property controls whether fill level indicator
// graphics are displayed on the trough. See SetShowFillLevel.
// Flags: Read / Write
// Default value: FALSE
const PropertyShowFillLevel cdk.Property = "show-fill-level"

// How the range should be updated on the screen.
// Flags: Read / Write
// Default value: GTK_UPDATE_CONTINUOUS
const PropertyUpdatePolicy cdk.Property = "update-policy"

// The sensitivity policy for the stepper that points to the adjustment's
// upper side.
// Flags: Read / Write
// Default value: GTK_SENSITIVITY_AUTO
const PropertyUpperStepperSensitivity cdk.Property = "upper-stepper-sensitivity"

// Listener function arguments:
// 	arg1 float64
const SignalAdjustBounds cdk.Signal = "adjust-bounds"

// The ::change-value signal is emitted when a scroll action is performed on
// a range. It allows an application to determine the type of scroll event
// that occurred and the resultant new value. The application can handle the
// event itself and return TRUE to prevent further processing. Or, by
// returning FALSE, it can pass the event to other handlers until the default
// CTK handler is reached. The value parameter is unrounded. An application
// that overrides the ::change-value signal is responsible for clamping the
// value to the desired number of decimal digits; the default CTK handler
// clamps the value based on “round_digits”. It is not possible to use
// delayed update policies in an overridden ::change-value handler.
const SignalChangeValue cdk.Signal = "change-value"

// Virtual function that moves the slider. Used for keybindings.
// Listener function arguments:
// 	step ScrollType	how to move the slider
const SignalMoveSlider cdk.Signal = "move-slider"

// Emitted when the range value changes.
const SignalRangeValueChanged cdk.Signal = "value-changed"
