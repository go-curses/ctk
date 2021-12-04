package ctk

import (
	"fmt"

	"github.com/go-curses/cdk"
	cenums "github.com/go-curses/cdk/lib/enums"
	cmath "github.com/go-curses/cdk/lib/math"
	"github.com/go-curses/cdk/lib/paint"
	"github.com/go-curses/cdk/lib/ptypes"
	"github.com/go-curses/cdk/memphis"
	"github.com/go-curses/ctk/lib/enums"
)

const TypeScrollbar cdk.CTypeTag = "ctk-scrollbar"

func init() {
	_ = cdk.TypesManager.AddType(TypeScrollbar, nil)
}

var (
	DefaultScrollbarTheme = paint.Theme{
		// slider
		Content: paint.ThemeAspect{
			Normal:      paint.DefaultColorStyle.Foreground(paint.ColorDarkGray).Background(paint.ColorSilver).Dim(true).Bold(false),
			Selected:    paint.DefaultColorStyle.Foreground(paint.ColorBlack).Background(paint.ColorWhite).Dim(false).Bold(true),
			Active:      paint.DefaultColorStyle.Foreground(paint.ColorBlack).Background(paint.ColorWhite).Dim(false).Bold(true),
			Prelight:    paint.DefaultColorStyle.Foreground(paint.ColorBlack).Background(paint.ColorGray).Dim(false),
			Insensitive: paint.DefaultColorStyle.Foreground(paint.ColorBlack).Background(paint.ColorGray).Dim(true),
			FillRune:    paint.DefaultFillRune,
			BorderRunes: paint.DefaultBorderRune,
			ArrowRunes:  paint.DefaultArrowRune,
			Overlay:     false,
		},
		// trough
		Border: paint.ThemeAspect{
			Normal:      paint.DefaultColorStyle.Foreground(paint.ColorBlack).Background(paint.ColorGray).Dim(true).Bold(false),
			Selected:    paint.DefaultColorStyle.Foreground(paint.ColorWhite).Background(paint.ColorDarkGray).Dim(false).Bold(true),
			Active:      paint.DefaultColorStyle.Foreground(paint.ColorBlack).Background(paint.ColorSilver).Dim(false).Bold(true),
			Prelight:    paint.DefaultColorStyle.Foreground(paint.ColorBlack).Background(paint.ColorDarkGray).Dim(false),
			Insensitive: paint.DefaultColorStyle.Foreground(paint.ColorBlack).Background(paint.ColorDarkGray).Dim(true),
			FillRune:    paint.DefaultFillRune,
			BorderRunes: paint.DefaultBorderRune,
			ArrowRunes:  paint.DefaultArrowRune,
			Overlay:     false,
		},
	}
)

// Scrollbar Hierarchy:
//	Object
//	  +- Widget
//	    +- Range
//	      +- Scrollbar
//	        +- HScrollbar
//	        +- VScrollbar
//
// The Scrollbar Widget is a Range Widget that draws steppers and sliders.
type Scrollbar interface {
	Range

	Init() (already bool)
	GetHasBackwardStepper() (hasBackwardStepper bool)
	SetHasBackwardStepper(hasBackwardStepper bool)
	GetHasForwardStepper() (hasForwardStepper bool)
	SetHasForwardStepper(hasForwardStepper bool)
	GetHasSecondaryBackwardStepper() (hasSecondaryBackwardStepper bool)
	SetHasSecondaryBackwardStepper(hasSecondaryBackwardStepper bool)
	GetHasSecondaryForwardStepper() (hasSecondaryForwardStepper bool)
	SetHasSecondaryForwardStepper(hasSecondaryForwardStepper bool)
	Forward(step int) cenums.EventFlag
	ForwardStep() cenums.EventFlag
	ForwardPage() cenums.EventFlag
	Backward(step int) cenums.EventFlag
	BackwardStep() cenums.EventFlag
	BackwardPage() cenums.EventFlag
	GetSizeRequest() (width, height int)
	GetWidgetAt(p *ptypes.Point2I) Widget
	ValueChanged()
	Changed()
	CancelEvent()
	GetAllStepperRegions() (fwd, bwd, sFwd, sBwd ptypes.Region)
	GetStepperRegions() (start, end ptypes.Region)
	GetTroughRegion() (region ptypes.Region)
	GetSliderRegion() (region ptypes.Region)
}

// The CScrollbar structure implements the Scrollbar interface and is exported
// to facilitate type embedding with custom implementations. No member variables
// are exported as the interface methods are the only intended means of
// interacting with Scrollbar objects.
type CScrollbar struct {
	CRange

	orientation     cenums.Orientation
	minSliderLength int
	sliderMoving    bool
	prevSliderPos   *ptypes.Point2I
	focusedButton   Button

	hasBackwardStepper          bool
	hasForwardStepper           bool
	hasSecondaryBackwardStepper bool
	hasSecondaryForwardStepper  bool

	slider                   Button
	backwardStepper          Button
	forwardStepper           Button
	secondaryBackwardStepper Button
	secondaryForwardStepper  Button
}

// Init initializes a Scrollbar object. This must be called at least once to
// set up the necessary defaults and allocate any memory structures. Calling
// this more than once is safe though unnecessary. Only the first call will
// result in any effect upon the Scrollbar instance. Init is used in the
// NewScrollbar constructor and only necessary when implementing a derivative
// Scrollbar type.
func (s *CScrollbar) Init() (already bool) {
	if s.InitTypeItem(TypeScrollbar, s) {
		return true
	}
	s.CRange.Init()
	s.flags = enums.NULL_WIDGET_FLAG
	s.SetFlags(enums.SENSITIVE | enums.PARENT_SENSITIVE | enums.CAN_FOCUS | enums.APP_PAINTABLE)
	if s.orientation == cenums.ORIENTATION_NONE {
		s.orientation = cenums.ORIENTATION_VERTICAL
	}
	s.SetTheme(DefaultScrollbarTheme)
	s.focusedButton = nil
	s.hasBackwardStepper = true
	s.hasForwardStepper = true
	s.hasSecondaryBackwardStepper = false
	s.hasSecondaryForwardStepper = false
	s.backwardStepper = s.makeStepperButton(enums.ArrowUp, false)
	s.forwardStepper = s.makeStepperButton(enums.ArrowDown, true)
	s.secondaryForwardStepper = s.makeStepperButton(enums.ArrowDown, true)
	s.secondaryBackwardStepper = s.makeStepperButton(enums.ArrowUp, false)
	l := NewLabel("*")
	l.Show()
	l.SetSingleLineMode(true)
	l.SetMaxWidthChars(1)
	s.slider = NewButtonWithWidget(l)
	s.Connect(SignalCdkEvent, ScrollbarEventHandle, s.event)
	s.Connect(SignalInvalidate, ScrollbarInvalidateHandle, s.invalidate)
	s.Connect(SignalResize, ScrollbarResizeHandle, s.resize)
	s.Connect(SignalDraw, ScrollbarDrawHandle, s.draw)
	return false
}

// GetHasBackwardStepper returns whether to display the standard backward arrow
// button.
// See: SetHasBackwardStepper()
//
// Locking: read
func (s *CScrollbar) GetHasBackwardStepper() (hasBackwardStepper bool) {
	s.RLock()
	hasBackwardStepper = s.hasBackwardStepper
	s.RUnlock()
	return
}

// SetHasBackwardStepper updates whether to display the standard backward arrow
// button.
//
// Locking: write
func (s *CScrollbar) SetHasBackwardStepper(hasBackwardStepper bool) {
	s.Lock()
	s.hasBackwardStepper = hasBackwardStepper
	s.Unlock()
}

// GetHasForwardStepper returns whether to display the standard forward arrow
// button.
//
// Locking: read
func (s *CScrollbar) GetHasForwardStepper() (hasForwardStepper bool) {
	s.RLock()
	hasForwardStepper = s.hasForwardStepper
	s.RUnlock()
	return
}

// SetHasForwardSettper updates whether to display the standard forward arrow
// button.
//
// Locking: write
func (s *CScrollbar) SetHasForwardStepper(hasForwardStepper bool) {
	s.Lock()
	s.hasForwardStepper = hasForwardStepper
	s.Unlock()
}

// GetHasSecondaryBackwardStepper returns whether to display a second backward
// arrow button on the opposite end of the scrollbar.
//
// Locking: read
func (s *CScrollbar) GetHasSecondaryBackwardStepper() (hasSecondaryBackwardStepper bool) {
	s.RLock()
	hasSecondaryBackwardStepper = s.hasSecondaryBackwardStepper
	s.RUnlock()
	return
}

// SetHasSecondaryBackwardStepper updates whether to display a second backward
// arrow button on the opposite end of the scrollbar.
//
// Locking: write
func (s *CScrollbar) SetHasSecondaryBackwardStepper(hasSecondaryBackwardStepper bool) {
	s.Lock()
	s.hasSecondaryBackwardStepper = hasSecondaryBackwardStepper
	s.Unlock()
}

// GetHasSecondaryForwardStepper returns whether to display a second backward
// arrow button on the opposite end of the scrollbar.
//
// Locking: read
func (s *CScrollbar) GetHasSecondaryForwardStepper() (hasSecondaryForwardStepper bool) {
	s.RLock()
	hasSecondaryForwardStepper = s.hasSecondaryForwardStepper
	s.RUnlock()
	return
}

// SetHasSecondaryForwardStepper updates whether to display a second backward
// arrow button on the opposite end of the scrollbar.
//
// Locking: write
func (s *CScrollbar) SetHasSecondaryForwardStepper(hasSecondaryForwardStepper bool) {
	s.Lock()
	s.hasSecondaryForwardStepper = hasSecondaryForwardStepper
	s.Unlock()
}

// Forward updates the scrollbar in a forward direction by the given step count.
// Returns EVENT_STOP if changes were made, EVENT_PASS otherwise.
//
// Locking: write
func (s *CScrollbar) Forward(step int) cenums.EventFlag {
	min, max := s.GetRange()
	value := s.GetValue()
	want := value + step
	s.SetValue(want)
	got := s.GetValue()
	s.Invalidate()
	s.LogDebug("Forward: (step: %v, wants: %d, got:%d, range: %d-%d)", step, want, got, min, max)
	if value != got {
		return cenums.EVENT_STOP
	}
	return cenums.EVENT_PASS
}

// ForwardStep updates the scrollbar in a forward direction by the configured
// step increment amount. Returns EVENT_STOP if changes were made, EVENT_PASS
// otherwise.
//
// Locking: write
func (s *CScrollbar) ForwardStep() cenums.EventFlag {
	step, _ := s.GetIncrements()
	return s.Forward(step)
}

// ForwardPage updates the scrollbar in a forward direction by the configured
// page increment amount. Returns EVENT_STOP if changes were made, EVENT_PASS
// otherwise.
//
// Locking: write
func (s *CScrollbar) ForwardPage() cenums.EventFlag {
	_, page := s.GetIncrements()
	return s.Forward(page)
}

// Backward updates the scrollbar in a backward direction by the given step count.
// Returns EVENT_STOP if changes were made, EVENT_PASS otherwise.
//
// Locking: write
func (s *CScrollbar) Backward(step int) cenums.EventFlag {
	min, max := s.GetRange()
	value := s.GetValue()
	want := value - step
	s.SetValue(want)
	got := s.GetValue()
	s.Invalidate()
	s.LogDebug("Backward: (step: %v, wants: %d, got:%d, range: %d-%d)", step, want, got, min, max)
	if value != got {
		return cenums.EVENT_STOP
	}
	return cenums.EVENT_PASS
}

// BackwardStep updates the scrollbar in a backward direction by the configured
// step increment amount. Returns EVENT_STOP if changes were made, EVENT_PASS
// otherwise.
//
// Locking: write
func (s *CScrollbar) BackwardStep() cenums.EventFlag {
	step, _ := s.GetIncrements()
	return s.Backward(step)
}

// BackwardPage updates the scrollbar in a backward direction by the configured
// page increment amount. Returns EVENT_STOP if changes were made, EVENT_PASS
// otherwise.
//
// Locking: write
func (s *CScrollbar) BackwardPage() cenums.EventFlag {
	_, page := s.GetIncrements()
	return s.Backward(page)
}

func (s *CScrollbar) GetSizeRequest() (width, height int) {
	size := ptypes.NewRectangle(s.CWidget.GetSizeRequest())
	s.RLock()
	switch s.orientation {
	case cenums.ORIENTATION_HORIZONTAL:
		size.H = 1
	case cenums.ORIENTATION_VERTICAL:
		fallthrough
	default:
		size.W = 1
	}
	width, height = size.W, size.H
	s.RUnlock()
	return
}

func (s *CScrollbar) GetWidgetAt(p *ptypes.Point2I) Widget {
	if s.HasPoint(p) && s.IsVisible() {
		return s
	}
	return nil
}

func (s *CScrollbar) ValueChanged() {
	s.Invalidate()
	s.Emit(SignalValueChanged, s)
}

func (s *CScrollbar) Changed() {
	s.Invalidate()
	s.Emit(SignalChanged, s)
}

func (s *CScrollbar) CancelEvent() {
	if s.HasEventFocus() {
		s.ReleaseEventFocus()
	}
	s.Lock()
	if s.slider != nil {
		s.slider.CancelEvent()
	}
	if s.forwardStepper != nil {
		s.forwardStepper.CancelEvent()
	}
	if s.backwardStepper != nil {
		s.backwardStepper.CancelEvent()
	}
	if s.secondaryForwardStepper != nil {
		s.secondaryForwardStepper.CancelEvent()
	}
	if s.secondaryBackwardStepper != nil {
		s.secondaryBackwardStepper.CancelEvent()
	}
	s.Unlock()
	s.Invalidate()
}

func (s *CScrollbar) GetAllStepperRegions() (fwd, bwd, sFwd, sBwd ptypes.Region) {
	start, end := s.GetStepperRegions()
	s.RLock()
	fwd, bwd, sFwd, sBwd = end, start, start, end
	switch s.orientation {
	case cenums.ORIENTATION_HORIZONTAL:
		if fwd.W == 2 {
			fwd.X += 1
			fwd.W = 1
		}
		if sFwd.W == 2 {
			sFwd.X += 1
			sFwd.W = 1
		}
	case cenums.ORIENTATION_VERTICAL:
	default:
		if fwd.H == 2 {
			fwd.Y += 1
			fwd.H = 1
		}
		if sFwd.H == 2 {
			sFwd.Y += 1
			sFwd.H = 1
		}
	}
	s.RUnlock()
	return
}

func (s *CScrollbar) GetStepperRegions() (start, end ptypes.Region) {
	alloc := s.GetAllocation()
	origin := s.GetOrigin()
	s.RLock()
	switch s.orientation {
	case cenums.ORIENTATION_HORIZONTAL:
		start.X, start.Y, start.W, start.H = origin.X, origin.Y, 0, 1
		if s.hasForwardStepper {
			start.W += 1
		}
		if s.hasSecondaryBackwardStepper {
			start.W += 1
		}
		end.X, end.Y, end.W, end.H = origin.X+alloc.W, origin.Y, 0, 1
		if s.hasBackwardStepper {
			end.W += 1
			end.X -= 1
		}
		if s.hasSecondaryForwardStepper {
			end.W += 1
			end.X -= 1
		}
	case cenums.ORIENTATION_VERTICAL:
		fallthrough
	default:
		start.X, start.Y, start.W, start.H = origin.X, origin.Y, 1, 0
		if s.hasBackwardStepper {
			start.H += 1
		}
		if s.hasSecondaryForwardStepper {
			start.H += 1
		}
		end.X, end.Y, end.W, end.H = origin.X, origin.Y+alloc.H, 1, 0
		if s.hasForwardStepper {
			end.Y -= 1
			end.H += 1
		}
		if s.hasSecondaryBackwardStepper {
			end.Y -= 1
			end.H += 1
		}
	}
	s.RUnlock()
	return
}

func (s *CScrollbar) GetTroughRegion() (region ptypes.Region) {
	alloc := s.GetAllocation()
	start, end := s.GetStepperRegions()
	s.RLock()
	region = ptypes.MakeRegion(start.X, start.Y, 1, 1)
	switch s.orientation {
	case cenums.ORIENTATION_HORIZONTAL:
		region.X += start.W
		region.W = alloc.W - start.W - end.W
	case cenums.ORIENTATION_VERTICAL:
		fallthrough
	default:
		region.Y += start.H
		region.H = alloc.H - start.H - end.H
	}
	region.Floor(0, 0)
	s.RUnlock()
	return
}

func (s *CScrollbar) GetSliderRegion() (region ptypes.Region) {
	trough := s.GetTroughRegion()
	upper, page, value := 0, 0, 0
	if adjustment := s.GetAdjustment(); adjustment != nil {
		upper = adjustment.GetUpper()
		page = adjustment.GetPageIncrement()
		value = adjustment.GetValue()
	} else {
		s.LogError("missing adjustment")
	}
	s.RLock()
	region = ptypes.MakeRegion(trough.X, trough.Y, 1, 1)
	switch s.orientation {
	case cenums.ORIENTATION_HORIZONTAL:
		if upper == 0 {
			region.W = trough.W
		} else {
			size, fullSize := 1, cmath.FloorI(trough.W-2, 1)
			if s.sliderSizeFixed {
				size = cmath.ClampI(s.sliderLength, s.minSliderLength, trough.W)
			} else {
				if fullSize > 1 {
					size = int((float64(page) / float64(upper)) * float64(fullSize))
				} else if s.minSliderLength > 0 {
					size = s.minSliderLength
				}
			}
			region.W = cmath.ClampI(size, 1, trough.W-1)
			inc := int((float64(value) / float64(upper)) * float64(trough.W-region.W))
			if inc == 0 && value > 0 {
				inc = 1
			} else if inc == upper && value < upper {
				inc -= 1
			}
			region.X += inc
			region.X = cmath.ClampI(region.X, 0, trough.X+trough.W-1)
		}
	case cenums.ORIENTATION_VERTICAL:
		fallthrough
	default:
		if upper == 0 {
			region.H = trough.H
		} else {
			size, fullSize := 1, cmath.FloorI(trough.H-2, 1)
			if s.sliderSizeFixed {
				size = cmath.ClampI(s.sliderLength, s.minSliderLength, trough.H)
			} else {
				if fullSize > 1 {
					size = int((float64(page) / float64(upper)) * float64(fullSize))
				} else if s.minSliderLength > 0 {
					size = s.minSliderLength
				}
			}
			region.H = cmath.ClampI(size, 1, trough.H-1)
			inc := int((float64(value) / float64(upper)) * float64(trough.H-region.H))
			if inc == 0 && value > 0 {
				inc = 1
			} else if inc == upper && value < upper {
				inc -= 1
			}
			region.Y += inc
			region.Y = cmath.ClampI(region.Y, 0, trough.Y+trough.H-1)
		}
	}
	region.Floor(0, 0)
	s.RUnlock()
	return
}

func (s *CScrollbar) processEventAtPoint(p *ptypes.Point2I, e *cdk.EventMouse) cenums.EventFlag {
	// me := NewMouseEvent(e)
	slider := s.GetSliderRegion()
	w := s.GetWidgetAt(p)
	switch e.State() {
	case cdk.BUTTON_PRESS:
		if w != nil && w.IsVisible() {
			if w.ObjectID() != s.ObjectID() {
				if wb, ok := w.Self().(*CButton); ok {
					wb.SetPressed(true)
					s.Lock()
					s.focusedButton = wb
					s.Unlock()
					return cenums.EVENT_STOP
				}
			}
			if slider.HasPoint(*p) {
				s.Lock()
				s.prevSliderPos = p.NewClone()
				s.sliderMoving = true
				s.focusedButton = nil
				s.Unlock()
				return cenums.EVENT_STOP
			}
		}
	case cdk.DRAG_START:
		s.Lock()
		if !s.sliderMoving {
			s.focusedButton = nil
			s.sliderMoving = true
		}
		s.Unlock()
		fallthrough
	case cdk.DRAG_MOVE:
		if s.sliderMoving {
			if s.prevSliderPos != nil {
				if s.prevSliderPos.X != p.X && s.orientation == cenums.ORIENTATION_HORIZONTAL {
					// moved horizontally
					if s.textDirection == enums.TextDirRtl {
						// left=forward, right=backward
						if p.X > s.prevSliderPos.X {
							// right=backward
							s.BackwardPage()
						} else if p.X < s.prevSliderPos.X {
							// left=forward
							s.ForwardPage()
						}
					} else {
						// left=backward, right=forward
						if p.X > s.prevSliderPos.X {
							// right=forward
							s.ForwardPage()
						} else if p.X < s.prevSliderPos.X {
							// left=backward
							s.BackwardPage()
						}
					}
					return cenums.EVENT_STOP
				}
				if s.prevSliderPos.Y != p.Y && s.orientation == cenums.ORIENTATION_VERTICAL {
					// moved vertically
					// down=forward, up=backward
					if p.Y > s.prevSliderPos.Y {
						// down=forward
						s.ForwardPage()
					} else if p.Y < s.prevSliderPos.Y {
						// up=backward
						s.BackwardPage()
					} else {
						// neither
					}
					return cenums.EVENT_STOP
				}
			}
			s.Lock()
			s.prevSliderPos = p.NewClone()
			s.Unlock()
		}
	case cdk.DRAG_STOP:
		if s.HasEventFocus() {
			s.ReleaseEventFocus()
		}
		s.Lock()
		s.focusedButton = nil
		s.sliderMoving = false
		s.prevSliderPos = nil
		s.Unlock()
		return cenums.EVENT_STOP
	case cdk.BUTTON_RELEASE:
		if s.HasEventFocus() {
			s.ReleaseEventFocus()
		}
		if s.focusedButton != nil {
			if s.focusedButton.HasPoint(p) {
				s.focusedButton.SetPressed(false)
				s.focusedButton.Activate()
				s.Lock()
				s.focusedButton = nil
				s.sliderMoving = false
				s.prevSliderPos = nil
				s.Unlock()
				return cenums.EVENT_STOP
			}
		}
		s.Lock()
		s.focusedButton = nil
		s.sliderMoving = false
		s.prevSliderPos = nil
		s.Unlock()
		slider := s.GetSliderRegion()
		if s.orientation == cenums.ORIENTATION_HORIZONTAL {
			if s.textDirection == enums.TextDirRtl {
				if p.X < slider.X {
					return s.ForwardPage()
				} else if p.X >= slider.X+slider.W {
					return s.BackwardPage()
				}
			} else {
				if p.X < slider.X {
					return s.BackwardPage()
				} else if p.X >= slider.X+slider.W {
					return s.ForwardPage()
				}
			}
		} else {
			if p.Y < slider.Y {
				return s.BackwardPage()
			} else if p.Y >= slider.Y+slider.H {
				return s.ForwardPage()
			}
		}
	}
	return cenums.EVENT_PASS
}

func (s *CScrollbar) event(data []interface{}, argv ...interface{}) cenums.EventFlag {
	if evt, ok := argv[1].(cdk.Event); ok {
		switch e := evt.(type) {
		case *cdk.EventMouse:
			point := ptypes.NewPoint2I(e.Position())
			return s.processEventAtPoint(point, e)
		case *cdk.EventKey:
			if s.HasEventFocus() {
				s.CancelEvent()
				return cenums.EVENT_STOP
			}
			adjustment := s.GetAdjustment()
			s.Lock()
			step, page := 0, 0
			if adjustment != nil {
				step, page = adjustment.GetStepIncrement(), adjustment.GetPageIncrement()
			} else {
				s.LogError("missing adjustment")
			}
			switch s.orientation {
			case cenums.ORIENTATION_HORIZONTAL:
				switch e.Key() {
				case cdk.KeyLeft:
					if e.Modifiers().Has(cdk.ModShift) {
						s.Unlock()
						return s.Backward(page)
					}
					s.Unlock()
					return s.Backward(step)
				case cdk.KeyRight:
					if e.Modifiers().Has(cdk.ModShift) {
						s.Unlock()
						return s.Forward(page)
					}
					s.Unlock()
					return s.Forward(step)
				}
			case cenums.ORIENTATION_VERTICAL:
				fallthrough
			default:
				switch e.Key() {
				case cdk.KeyUp:
					s.Unlock()
					return s.Backward(step)
				case cdk.KeyDown:
					s.Unlock()
					return s.Forward(step)
				case cdk.KeyPgUp:
					s.Unlock()
					return s.Backward(page)
				case cdk.KeyPgDn:
					s.Unlock()
					return s.Forward(page)
				}
			}
			s.Unlock()
		}
	}
	return cenums.EVENT_PASS
}

func (s *CScrollbar) resize(data []interface{}, argv ...interface{}) cenums.EventFlag {
	s.resizeSteppers()
	s.resizeSlider()
	s.Invalidate()
	return cenums.EVENT_STOP
}

func (s *CScrollbar) resizeSteppers() {
	fwd, bwd, sFwd, sBwd := s.GetAllStepperRegions()
	aFwd, aBwd := enums.ArrowDown, enums.ArrowUp
	if s.orientation == cenums.ORIENTATION_HORIZONTAL {
		aFwd, aBwd = enums.ArrowRight, enums.ArrowLeft
	}
	s.resizeStepper(
		aFwd, aBwd,
		s.hasForwardStepper, s.forwardStepper,
		true,
		fwd.X, fwd.Y, fwd.W, fwd.H,
	)
	s.resizeStepper(
		aFwd, aBwd,
		s.hasBackwardStepper, s.backwardStepper,
		false,
		bwd.X, bwd.Y, bwd.W, bwd.H,
	)
	s.resizeStepper(
		aFwd, aBwd,
		s.hasSecondaryForwardStepper, s.secondaryForwardStepper,
		true,
		sFwd.X, sFwd.Y, sFwd.W, sFwd.H,
	)
	s.resizeStepper(
		aFwd, aBwd,
		s.hasSecondaryBackwardStepper, s.secondaryBackwardStepper,
		false,
		sBwd.X, bwd.Y, sBwd.W, sBwd.H,
	)
}

func (s *CScrollbar) makeStepperButton(arrow enums.ArrowType, forward bool) Button {
	a := NewArrow(arrow)
	a.Show()
	a.SetOrigin(0, 0)
	// a.SetSizeRequest(1, 1)
	a.SetAllocation(ptypes.MakeRectangle(1, 1))
	a.SetTheme(DefaultButtonTheme)
	a.UnsetFlags(enums.CAN_FOCUS)
	b := NewButtonWithWidget(a)
	b.ShowAll()
	b.SetFocusOnClick(false)
	b.SetParent(s.GetParent())
	b.SetWindow(s.GetWindow())
	b.Connect(
		SignalActivate,
		fmt.Sprintf("%v.activate", s.ObjectName()),
		func(data []interface{}, argv ...interface{}) cenums.EventFlag {
			if adjustment := s.GetAdjustment(); adjustment != nil {
				if forward {
					s.Forward(adjustment.GetStepIncrement())
				} else {
					s.Backward(adjustment.GetStepIncrement())
				}
			} else {
				s.LogError("missing adjustment")
			}
			return cenums.EVENT_STOP
		},
	)
	return b
}

func (s *CScrollbar) resizeStepper(fArrow, bArrow enums.ArrowType, has bool, b Button, forward bool, x, y, w, h int) {
	if has {
		if bc := b.GetChild(); bc != nil {
			if ba, ok := bc.Self().(*CArrow); ok {
				if forward {
					if ba.GetArrowType() != fArrow {
						ba.SetArrowType(fArrow)
					}
				} else {
					if ba.GetArrowType() != bArrow {
						ba.SetArrowType(bArrow)
					}
				}
			}
		}
		b.SetOrigin(x, y)
		// b.SetSizeRequest(w, h)
		b.SetAllocation(ptypes.MakeRectangle(w, h))
		b.ShowAll()
		b.Resize()
	} else {
		b.Hide()
	}
	return
}

func (s *CScrollbar) resizeSlider() {
	sr := s.GetSliderRegion()
	s.Lock()
	s.slider.SetOrigin(sr.X, sr.Y)
	s.slider.SetSizeRequest(sr.W, sr.H)
	s.slider.SetAllocation(sr.Size())
	s.slider.ShowAll()
	s.slider.Resize()
	s.Unlock()
}

func (s *CScrollbar) invalidate(data []interface{}, argv ...interface{}) cenums.EventFlag {
	// s.resizeSteppers()
	// s.resizeSlider()
	origin := s.GetOrigin()
	isSelected := s.HasState(enums.StateSelected)
	size := ptypes.MakeRectangle(1, 1)
	theme := s.GetThemeRequest()
	style := theme.Content.Normal
	doConfigure := func(b Button) {
		if b != nil {
			bid := b.ObjectID()
			if isSelected {
				b.SetState(enums.StateSelected)
			} else {
				b.UnsetState(enums.StateSelected)
			}
			local := b.GetOrigin()
			local.SubPoint(origin)
			if err := memphis.ConfigureSurface(bid, local, size, style); err != nil {
				b.LogErr(err)
			}
			b.Invalidate()
		}
	}
	doConfigure(s.slider)
	doConfigure(s.forwardStepper)
	doConfigure(s.backwardStepper)
	doConfigure(s.secondaryForwardStepper)
	doConfigure(s.secondaryBackwardStepper)
	return cenums.EVENT_STOP
}

func (s *CScrollbar) draw(data []interface{}, argv ...interface{}) cenums.EventFlag {
	if surface, ok := argv[1].(*memphis.CSurface); ok {
		alloc := s.GetAllocation()
		if !s.IsVisible() || alloc.W <= 0 || alloc.H <= 0 {
			s.LogTrace("not visible, zero width or zero height")
			return cenums.EVENT_PASS
		}

		s.LockDraw()
		defer s.UnlockDraw()

		theme := s.GetThemeRequest()
		origin := s.GetOrigin()
		trough := s.GetTroughRegion()
		trough.X -= origin.X
		trough.Y -= origin.Y

		fwd, bwd, sFwd, sBwd := s.GetAllStepperRegions()

		surface.Box(
			trough.Origin(), trough.Size(),
			false, true,
			theme.Border.Overlay,
			theme.Content.FillRune,
			theme.Border.Normal,
			theme.Border.Normal,
			theme.Border.BorderRunes,
		)

		// draw the slider
		if slider := s.slider; slider != nil {
			sliderOrigin := slider.GetOrigin()
			sliderOrigin.SubPoint(origin)
			sliderSize := slider.GetAllocation()
			surface.Box(
				sliderOrigin, sliderSize,
				false, true,
				theme.Content.Overlay,
				theme.Content.FillRune,
				theme.Content.Normal,
				theme.Border.Normal,
				theme.Border.BorderRunes,
			)
		}

		// draw the stepper buttons
		drawStepper := func(has bool, b Button, r ptypes.Region) error {
			if has && b != nil {
				if f := b.Draw(); f == cenums.EVENT_STOP {
					return surface.Composite(b.ObjectID())
				}
			}
			return nil
		}

		if err := drawStepper(s.hasBackwardStepper, s.backwardStepper, bwd); err != nil {
			s.LogError("error compositing backward stepper: %v", err)
		}
		if err := drawStepper(s.hasForwardStepper, s.forwardStepper, fwd); err != nil {
			s.LogError("error compositing forward stepper: %v", err)
		}
		if err := drawStepper(s.hasSecondaryBackwardStepper, s.secondaryBackwardStepper, sBwd); err != nil {
			s.LogError("error compositing secondary backward stepper: %v", err)
		}
		if err := drawStepper(s.hasSecondaryForwardStepper, s.secondaryForwardStepper, sFwd); err != nil {
			s.LogError("error compositing secondary forward stepper: %v", err)
		}

		return cenums.EVENT_STOP
	}
	return cenums.EVENT_PASS
}

const ScrollbarEventHandle = "scrollbar-event-handler"

const ScrollbarInvalidateHandle = "scrollbar-invalidate-handler"

const ScrollbarResizeHandle = "scrollbar-resize-handler"

const ScrollbarDrawHandle = "scrollbar-draw-handler"
