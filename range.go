package ctk

import (
	"github.com/go-curses/cdk"
	cmath "github.com/go-curses/cdk/lib/math"
	"github.com/go-curses/cdk/lib/ptypes"
)

// CDK type-tag for Range objects
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
	GetRangeRect(rangeRect ptypes.Rectangle)
	GetSliderRange(sliderStart int, sliderEnd int)
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

// The CRange structure implements the Range interface and is
// exported to facilitate type embedding with custom implementations. No member
// variables are exported as the interface methods are the only intended means
// of interacting with Range objects
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

// Range object initialization. This must be called at least once to setup
// the necessary defaults and allocate any memory structures. Calling this more
// than once is safe though unnecessary. Only the first call will result in any
// effect upon the Range instance
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

// Gets the current position of the fill level indicator.
// Returns:
// 	The current fill level
func (r *CRange) GetFillLevel() (value float64) {
	var err error
	if value, err = r.GetFloat64Property(PropertyFillLevel); err != nil {
		r.LogErr(err)
	}
	return
}

// Gets whether the range is restricted to the fill level.
// Returns:
// 	TRUE if range is restricted to the fill level.
func (r *CRange) GetRestrictToFillLevel() (value bool) {
	var err error
	if value, err = r.GetBoolProperty(PropertyRestrictToFillLevel); err != nil {
		r.LogErr(err)
	}
	return
}

// Gets whether the range displays the fill level graphically.
// Returns:
// 	TRUE if range shows the fill level.
func (r *CRange) GetShowFillLevel() (value bool) {
	var err error
	if value, err = r.GetBoolProperty(PropertyShowFillLevel); err != nil {
		r.LogErr(err)
	}
	return
}

// Set the new position of the fill level indicator. The "fill level" is
// probably best described by its most prominent use case, which is an
// indicator for the amount of pre-buffering in a streaming media player. In
// that use case, the value of the range would indicate the current play
// position, and the fill level would be the position up to which the
// file/stream has been downloaded. This amount of prebuffering can be
// displayed on the range's trough and is themeable separately from the
// trough. To enable fill level display, use SetShowFillLevel.
// The range defaults to not showing the fill level. Additionally, it's
// possible to restrict the range's slider position to values which are
// smaller than the fill level. This is controller by
// SetRestrictToFillLevel and is by default enabled.
// Parameters:
// 	fillLevel	the new position of the fill level indicator
func (r *CRange) SetFillLevel(fillLevel float64) {
	if err := r.SetFloatProperty(PropertyFillLevel, cmath.ClampF(fillLevel, 0.0, 1.0)); err != nil {
		r.LogErr(err)
	}
}

// Sets whether the slider is restricted to the fill level. See
// SetFillLevel for a general description of the fill level
// concept.
// Parameters:
// 	restrictToFillLevel	Whether the fill level restricts slider movement.
func (r *CRange) SetRestrictToFillLevel(restrictToFillLevel bool) {
	if err := r.SetBoolProperty(PropertyRestrictToFillLevel, restrictToFillLevel); err != nil {
		r.LogErr(err)
	}
}

// Sets whether a graphical fill level is show on the trough. See
// SetFillLevel for a general description of the fill level
// concept.
// Parameters:
// 	showFillLevel	Whether a fill level indicator graphics is shown.
func (r *CRange) SetShowFillLevel(showFillLevel bool) {
	if err := r.SetBoolProperty(PropertyShowFillLevel, showFillLevel); err != nil {
		r.LogErr(err)
	}
}

// Get the Adjustment which is the "model" object for Range. See
// SetAdjustment for details. The return value does not have a
// reference added, so should not be unreferenced.
// Returns:
// 	a Adjustment.
// 	[transfer none]
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

// Gets the value set by SetInverted.
// Returns:
// 	TRUE if the range is inverted
func (r *CRange) GetInverted() (value bool) {
	var err error
	if value, err = r.GetBoolProperty(PropertyInverted); err != nil {
		r.LogErr(err)
	}
	return
}

// Ranges normally move from lower to higher values as the slider moves from
// top to bottom or left to right. Inverted ranges have higher values at the
// top or on the right rather than on the bottom or left.
// Parameters:
// 	setting	TRUE to invert the range
func (r *CRange) SetInverted(setting bool) {
	if err := r.SetBoolProperty(PropertyInverted, setting); err != nil {
		r.LogErr(err)
	}
}

// Sets the step and page sizes for the range. The step size is used when the
// user clicks the Scrollbar arrows or moves Scale via arrow keys. The
// page size is used for example when moving via Page Up or Page Down keys.
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

// Sets the allowable values in the Range, and clamps the range value to
// be between min and max . (If the range has a non-zero page size, it is
// clamped between min and max - page-size.)
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

// Gets the current value of the range.
// Returns:
// 	current value of the range.
func (r *CRange) GetValue() (value int) {
	value = -1
	if adjustment := r.GetAdjustment(); adjustment != nil {
		value = adjustment.GetValue()
	} else {
		r.LogError("missing adjustment")
	}
	return
}

// Sets the current value of the range; if the value is outside the minimum
// or maximum range values, it will be clamped to fit inside them. The range
// emits the “value-changed” signal if the value changes.
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

// Gets the number of digits to round the value to when it changes. See
// “change-value”.
// Returns:
// 	the number of digits to round to
func (r *CRange) GetRoundDigits() (value int) {
	var err error
	if value, err = r.GetIntProperty(PropertyRoundDigits); err != nil {
		r.LogErr(err)
	}
	return
}

// Sets the number of digits to round the value to when it changes. See
// “change-value”.
// Parameters:
// 	roundDigits	the precision in digits, or -1
func (r *CRange) SetRoundDigits(roundDigits int) {
	if err := r.SetIntProperty(PropertyRoundDigits, roundDigits); err != nil {
		r.LogErr(err)
	}
}

// Sets the sensitivity policy for the stepper that points to the 'lower' end
// of the Range's adjustment.
// Parameters:
// 	sensitivity	the lower stepper's sensitivity policy.
func (r *CRange) SetLowerStepperSensitivity(sensitivity SensitivityType) {
	if err := r.SetStructProperty(PropertyLowerStepperSensitivity, sensitivity); err != nil {
		r.LogErr(err)
	}
}

// Gets the sensitivity policy for the stepper that points to the 'lower' end
// of the Range's adjustment.
// Returns:
// 	The lower stepper's sensitivity policy.
func (r *CRange) GetLowerStepperSensitivity() (value SensitivityType) {
	var ok bool
	if v, err := r.GetStructProperty(PropertyLowerStepperSensitivity); err != nil {
		r.LogErr(err)
	} else if value, ok = v.(SensitivityType); !ok {
		r.LogError("value stored in %v property is not of SensitivityType type: %v (%T)", PropertyLowerStepperSensitivity, v, v)
	}
	return
}

// Sets the sensitivity policy for the stepper that points to the 'upper' end
// of the Range's adjustment.
// Parameters:
// 	sensitivity	the upper stepper's sensitivity policy.
func (r *CRange) SetUpperStepperSensitivity(sensitivity SensitivityType) {
	if err := r.SetStructProperty(PropertyUpperStepperSensitivity, sensitivity); err != nil {
		r.LogErr(err)
	}
}

// Gets the sensitivity policy for the stepper that points to the 'upper' end
// of the Range's adjustment.
// Returns:
// 	The upper stepper's sensitivity policy.
func (r *CRange) GetUpperStepperSensitivity() (value SensitivityType) {
	var ok bool
	if v, err := r.GetStructProperty(PropertyUpperStepperSensitivity); err != nil {
		r.LogErr(err)
	} else if value, ok = v.(SensitivityType); !ok {
		r.LogError("value stored in %v property is not of SensitivityType type: %v (%T)", PropertyUpperStepperSensitivity, v, v)
	}
	return
}

// Gets the value set by SetFlippable.
// Returns:
// 	TRUE if the range is flippable
func (r *CRange) GetFlippable() (value bool) {
	return r.flippable
}

// If a range is flippable, it will switch its direction if it is horizontal
// and its direction is GTK_TEXT_DIR_RTL. See WidgetGetDirection.
// Parameters:
// 	flippable	TRUE to make the range flippable
func (r *CRange) SetFlippable(flippable bool) {
	r.flippable = flippable
}

// This function is useful mainly for Range subclasses. See
// SetMinSliderSize.
// Returns:
// 	The minimum size of the range's slider.
func (r *CRange) GetMinSliderSize() (value int) {
	return 0
}

// This function returns the area that contains the range's trough and its
// steppers, in widget->window coordinates. This function is useful mainly
// for Range subclasses.
// Parameters:
// 	rangeRect	return location for the range rectangle.
func (r *CRange) GetRangeRect(rangeRect ptypes.Rectangle) {}

// This function returns sliders range along the long dimension, in
// widget->window coordinates. This function is useful mainly for Range
// subclasses.
// Parameters:
// 	sliderStart	return location for the slider's
// start, or NULL.
// 	sliderEnd	return location for the slider's
// end, or NULL.
func (r *CRange) GetSliderRange(sliderStart int, sliderEnd int) {}

// This function is useful mainly for Range subclasses. See
// SetSliderSizeFixed.
// Returns:
// 	whether the range's slider has a fixed size.
func (r *CRange) GetSliderSizeFixed() (value bool) {
	return r.sliderSizeFixed
}

// Sets the minimum size of the range's slider. This function is useful
// mainly for Range subclasses.
// Parameters:
// 	minSize	The slider's minimum size
func (r *CRange) SetMinSliderSize(minSize bool) {}

// Sets whether the range's slider has a fixed size, or a size that depends
// on it's adjustment's page size. This function is useful mainly for
// Range subclasses.
// Parameters:
// 	sizeFixed	TRUE to make the slider size constant
func (r *CRange) SetSliderSizeFixed(sizeFixed bool) {
	r.sliderSizeFixed = sizeFixed
}

// Gets the step and page sizes for the range.
//   The step size is used when the user clicks the Scrollbar
//   arrows or moves Scale via arrow keys. The page size
//   is used for example when moving via Page Up or Page Down keys.
// Returns:
// 	step	step size
// 	page	page size
func (r *CRange) GetIncrements() (step int, page int) {
	step, page = -1, -1
	if adjustment := r.GetAdjustment(); adjustment != nil {
		step, page = adjustment.GetStepIncrement(), adjustment.GetPageIncrement()
	} else {
		r.LogError("missing adjustment")
	}
	return
}

// Sets the allowable values in the Range, and clamps the range
//   value to be between min
//    and max
//   . (If the range has a non-zero
//   page size, it is clamped between min
//    and max
//    - page-size.)
// Parameters:
// 	min	minimum range value
// 	max	maximum range value
func (r *CRange) GetRange() (min, max int) {
	min, max = -1, -1
	if adjustment := r.GetAdjustment(); adjustment != nil {
		min, max = adjustment.GetLower(), adjustment.GetUpper()
	} else {
		r.LogError("missing adjustment")
	}
	return
}

// This function is useful mainly for Range subclasses.
// See SetMinSliderLength().
// Returns:
// 	The minimum size of the range's slider.
func (r *CRange) GetMinSliderLength() (length int) {
	length = r.minSliderLength
	if r.minSliderLength <= -1 {
		length = 1
	}
	return
}

// Sets the minimum size of the range's slider.
// This function is useful mainly for Range subclasses.
// Parameters:
// 	minSize	The slider's minimum size
func (r *CRange) SetMinSliderLength(length int) {
	r.minSliderLength = length
}

// Returns the length of the scrollbar or scale thumb.
// Flags: Read
func (r *CRange) GetSliderLength() (length int) {
	length = r.sliderLength
	return
}

// Set the length of the scrollbar or scale thumb. Sets fixed slider length to
//  true. Set to -1 for variable slider length.
// Flags: Read
// Allowed values: >= 0
// Default value: 14
func (r *CRange) SetSliderLength(length int) {
	if length <= -1 {
		r.sliderLength = -1
		r.SetSliderSizeFixed(false)
	} else if r.sliderLength != length {
		r.SetSliderSizeFixed(true)
		r.sliderLength = length
	}
}

// Length of step buttons at ends.
// Flags: Read
// Allowed values: >= 0
// Default value: 1
func (r *CRange) GetStepperSize() (size int) {
	return r.stepperSize
}

// Length of step buttons at ends.
// Flags: Read
// Allowed values: >= 0
// Default value: 1
func (r *CRange) SetStepperSize(size int) {
	r.stepperSize = size
}

// The spacing between the stepper buttons and thumb. Note that
// setting this value to anything > 0 will automatically set the
// trough-under-steppers style property to TRUE as well. Also,
// stepper-spacing won't have any effect if there are no steppers.
// Flags: Read
// Allowed values: >= 0
// Default value: 0
func (r *CRange) GetStepperSpacing() (spacing int) {
	return r.stepperSpacing
}

// The spacing between the stepper buttons and thumb. Note that
// setting this value to anything > 0 will automatically set the
// trough-under-steppers style property to TRUE as well. Also,
// stepper-spacing won't have any effect if there are no steppers.
// Flags: Read
// Allowed values: >= 0
// Default value: 0
func (r *CRange) SetStepperSpacing(spacing int) {
	r.stepperSpacing = spacing
}

// Whether to draw the trough across the full length of the range or
// to exclude the steppers and their spacing. Note that setting the
// stepper-spacing style property to any value > 0 will
// automatically enable trough-under-steppers too.
// Flags: Read
// Default value: TRUE
// Flags: Run Last
func (r *CRange) GetTroughUnderSteppers() (underSteppers bool) {
	return r.troughUnderSteppers
}

// Whether to draw the trough across the full length of the range or
// to exclude the steppers and their spacing. Note that setting the
// stepper-spacing style property to any value > 0 will
// automatically enable trough-under-steppers too.
// Flags: Read
// Default value: TRUE
// Flags: Run Last
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
