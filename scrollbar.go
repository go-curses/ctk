package ctk

import (
	"fmt"

	"github.com/go-curses/cdk"
	"github.com/go-curses/cdk/lib/enums"
	cmath "github.com/go-curses/cdk/lib/math"
	"github.com/go-curses/cdk/lib/paint"
	"github.com/go-curses/cdk/lib/ptypes"
	"github.com/go-curses/cdk/memphis"
)

// CDK type-tag for Scrollbar objects
const TypeScrollbar cdk.CTypeTag = "ctk-scrollbar"

func init() {
	_ = cdk.TypesManager.AddType(TypeScrollbar, nil)
}

var (
	DefaultMonoScrollbarTheme = paint.Theme{
		Content: paint.ThemeAspect{
			Normal:      paint.DefaultMonoStyle,
			Focused:     paint.DefaultMonoStyle.Dim(false),
			Active:      paint.DefaultMonoStyle.Dim(false).Bold(true),
			FillRune:    paint.DefaultFillRune,
			BorderRunes: paint.DefaultBorderRune,
			ArrowRunes:  paint.DefaultArrowRune,
			Overlay:     false,
		},
		Border: paint.ThemeAspect{
			Normal:      paint.DefaultMonoStyle,
			Focused:     paint.DefaultMonoStyle.Dim(false),
			Active:      paint.DefaultMonoStyle.Dim(false).Bold(true),
			FillRune:    paint.DefaultFillRune,
			BorderRunes: paint.DefaultBorderRune,
			ArrowRunes:  paint.DefaultArrowRune,
			Overlay:     false,
		},
	}
	DefaultColorScrollbarTheme = paint.Theme{
		// slider
		Content: paint.ThemeAspect{
			Normal:      paint.DefaultColorStyle.Foreground(paint.ColorDarkGray).Background(paint.ColorSilver).Dim(true).Bold(false),
			Focused:     paint.DefaultColorStyle.Foreground(paint.ColorBlack).Background(paint.ColorWhite).Dim(false).Bold(true),
			Active:      paint.DefaultColorStyle.Foreground(paint.ColorBlack).Background(paint.ColorWhite).Dim(false).Bold(true),
			FillRune:    paint.DefaultFillRune,
			BorderRunes: paint.DefaultBorderRune,
			ArrowRunes:  paint.DefaultArrowRune,
			Overlay:     false,
		},
		// trough
		Border: paint.ThemeAspect{
			Normal:      paint.DefaultColorStyle.Foreground(paint.ColorBlack).Background(paint.ColorGray).Dim(true).Bold(false),
			Focused:     paint.DefaultColorStyle.Foreground(paint.ColorWhite).Background(paint.ColorDarkGray).Dim(false).Bold(true),
			Active:      paint.DefaultColorStyle.Foreground(paint.ColorBlack).Background(paint.ColorSilver).Dim(false).Bold(true),
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
	Forward(step int) enums.EventFlag
	ForwardStep() enums.EventFlag
	ForwardPage() enums.EventFlag
	Backward(step int) enums.EventFlag
	BackwardStep() enums.EventFlag
	BackwardPage() enums.EventFlag
	GetSizeRequest() (width, height int)
	Resize() enums.EventFlag
	GetWidgetAt(p *ptypes.Point2I) Widget
	ValueChanged()
	Changed()
	GrabFocus()
	CancelEvent()
	ProcessEvent(evt cdk.Event) enums.EventFlag
	ProcessEventAtPoint(p *ptypes.Point2I, e *cdk.EventMouse) enums.EventFlag
	Invalidate() enums.EventFlag
	GetAllStepperRegions() (fwd, bwd, sFwd, sBwd ptypes.Region)
	GetStepperRegions() (start, end ptypes.Region)
	GetTroughRegion() (region ptypes.Region)
	GetSliderRegion() (region ptypes.Region)
}

// The CScrollbar structure implements the Scrollbar interface and is
// exported to facilitate type embedding with custom implementations. No member
// variables are exported as the interface methods are the only intended means
// of interacting with Scrollbar objects
type CScrollbar struct {
	CRange

	orientation     enums.Orientation
	minSliderLength int
	sliderMoving    bool
	prevSliderPos   *ptypes.Point2I
	focusedButton   *CButton

	hasBackwardStepper          bool
	hasForwardStepper           bool
	hasSecondaryBackwardStepper bool
	hasSecondaryForwardStepper  bool

	slider                   *CButton
	backwardStepper          *CButton
	forwardStepper           *CButton
	secondaryBackwardStepper *CButton
	secondaryForwardStepper  *CButton
}

// Scrollbar object initialization. This must be called at least once to setup
// the necessary defaults and allocate any memory structures. Calling this more
// than once is safe though unnecessary. Only the first call will result in any
// effect upon the Scrollbar instance
func (s *CScrollbar) Init() (already bool) {
	if s.InitTypeItem(TypeScrollbar, s) {
		return true
	}
	s.CRange.Init()
	s.flags = NULL_WIDGET_FLAG
	s.SetFlags(SENSITIVE | PARENT_SENSITIVE | CAN_FOCUS | APP_PAINTABLE)
	if s.orientation == enums.ORIENTATION_NONE {
		s.orientation = enums.ORIENTATION_VERTICAL
	}
	s.focusedButton = nil
	s.hasBackwardStepper = true
	s.hasForwardStepper = true
	s.hasSecondaryBackwardStepper = false
	s.hasSecondaryForwardStepper = false
	s.SetTheme(DefaultColorScrollbarTheme)
	s.Resize()
	s.Connect(SignalDraw, ScrollbarDrawHandle, s.draw)
	return false
}

// Display the standard backward arrow button.
// Flags: Read
// Default value: TRUE
func (s *CScrollbar) GetHasBackwardStepper() (hasBackwardStepper bool) {
	return s.hasBackwardStepper
}

// Display the standard backward arrow button.
// Flags: Read
// Default value: TRUE
func (s *CScrollbar) SetHasBackwardStepper(hasBackwardStepper bool) {
	s.hasBackwardStepper = hasBackwardStepper
	s.backwardStepper.Destroy()
	s.backwardStepper = nil
}

// Display the standard forward arrow button.
// Flags: Read
// Default value: TRUE
func (s *CScrollbar) GetHasForwardStepper() (hasForwardStepper bool) {
	return s.hasForwardStepper
}

// Display the standard forward arrow button.
// Flags: Read
// Default value: TRUE
func (s *CScrollbar) SetHasForwardStepper(hasForwardStepper bool) {
	s.hasForwardStepper = hasForwardStepper
	if !hasForwardStepper {
		s.forwardStepper.Destroy()
		s.forwardStepper = nil
	}
}

// Display a second backward arrow button on the opposite end of the scrollbar.
// Flags: Read
// Default value: FALSE
func (s *CScrollbar) GetHasSecondaryBackwardStepper() (hasSecondaryBackwardStepper bool) {
	return s.hasSecondaryBackwardStepper
}

// Display a second backward arrow button on the opposite end of the scrollbar.
// Flags: Read
// Default value: FALSE
func (s *CScrollbar) SetHasSecondaryBackwardStepper(hasSecondaryBackwardStepper bool) {
	s.hasSecondaryBackwardStepper = hasSecondaryBackwardStepper
	if !hasSecondaryBackwardStepper {
		s.secondaryBackwardStepper.Destroy()
		s.secondaryBackwardStepper = nil
	}
}

// Display a second forward arrow button on the opposite end of the scrollbar.
// Flags: Read
// Default value: FALSE
func (s *CScrollbar) GetHasSecondaryForwardStepper() (hasSecondaryForwardStepper bool) {
	return s.hasSecondaryForwardStepper
}

// Display a second forward arrow button on the opposite end of the scrollbar.
// Flags: Read
// Default value: FALSE
func (s *CScrollbar) SetHasSecondaryForwardStepper(hasSecondaryForwardStepper bool) {
	s.hasSecondaryForwardStepper = hasSecondaryForwardStepper
	if !hasSecondaryForwardStepper {
		s.secondaryForwardStepper.Destroy()
		s.secondaryForwardStepper = nil
	}
}

func (s *CScrollbar) Forward(step int) enums.EventFlag {
	min, max := s.GetRange()
	value := s.GetValue()
	want := value + step
	s.SetValue(want)
	got := s.GetValue()
	s.LogDebug("Forward: (step: %v, wants: %d, got:%d, range: %d-%d)", step, want, got, min, max)
	if value != got {
		s.Invalidate()
		return enums.EVENT_STOP
	}
	return enums.EVENT_PASS
}

func (s *CScrollbar) ForwardStep() enums.EventFlag {
	step, _ := s.GetIncrements()
	return s.Forward(step)
}

func (s *CScrollbar) ForwardPage() enums.EventFlag {
	_, page := s.GetIncrements()
	return s.Forward(page)
}

func (s *CScrollbar) Backward(step int) enums.EventFlag {
	min, max := s.GetRange()
	value := s.GetValue()
	want := value - step
	s.SetValue(want)
	got := s.GetValue()
	s.LogDebug("Backward: (step: %v, wants: %d, got:%d, range: %d-%d)", step, want, got, min, max)
	if value != got {
		s.Invalidate()
		return enums.EVENT_STOP
	}
	return enums.EVENT_PASS
}

func (s *CScrollbar) BackwardStep() enums.EventFlag {
	step, _ := s.GetIncrements()
	return s.Backward(step)
}

func (s *CScrollbar) BackwardPage() enums.EventFlag {
	_, page := s.GetIncrements()
	return s.Backward(page)
}

func (s *CScrollbar) GetSizeRequest() (width, height int) {
	size := ptypes.NewRectangle(s.CWidget.GetSizeRequest())
	switch s.orientation {
	case enums.ORIENTATION_HORIZONTAL:
		size.H = 1
	case enums.ORIENTATION_VERTICAL:
		fallthrough
	default:
		size.W = 1
	}
	return size.W, size.H
}

func (s *CScrollbar) Resize() enums.EventFlag {
	s.resizeSteppers()
	s.resizeSlider()
	s.Invalidate()
	return s.Emit(SignalResize, s)
}

func (s *CScrollbar) GetWidgetAt(p *ptypes.Point2I) Widget {
	if s.HasPoint(p) && s.IsVisible() {
		fwd, bwd, sFwd, sBwd := s.GetAllStepperRegions()
		if s.hasForwardStepper && s.forwardStepper != nil && fwd.HasPoint(*p) {
			return s.forwardStepper
		}
		if s.hasBackwardStepper && s.backwardStepper != nil && bwd.HasPoint(*p) {
			return s.backwardStepper
		}
		if s.hasSecondaryForwardStepper && s.secondaryForwardStepper != nil && sFwd.HasPoint(*p) {
			return s.secondaryForwardStepper
		}
		if s.hasSecondaryBackwardStepper && s.secondaryBackwardStepper != nil && sBwd.HasPoint(*p) {
			return s.secondaryBackwardStepper
		}
		if s.slider != nil && s.slider.HasPoint(p) {
			return s.slider
		}
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

// If the Widget instance CanFocus() then take the focus of the associated
// Window. Any previously focused Widget will emit a lost-focus signal and the
// newly focused Widget will emit a gained-focus signal. This method emits a
// grab-focus signal initially and if the listeners return EVENT_PASS, the
// changes are applied
//
// Emits: SignalGrabFocus, Argv=[Widget instance]
// Emits: SignalLostFocus, Argv=[Previous focus Widget instance], From=Previous focus Widget instance
// Emits: SignalGainedFocus, Argv=[Widget instance, previous focus Widget instance]
func (s *CScrollbar) GrabFocus() {
	if s.CanFocus() {
		if r := s.Emit(SignalGrabFocus, s); r == enums.EVENT_PASS {
			tl := s.GetWindow()
			if tl != nil {
				var fw Widget
				focused := tl.GetFocus()
				tl.SetFocus(s)
				if focused != nil {
					var ok bool
					if fw, ok = focused.(Widget); ok && fw.ObjectID() != s.ObjectID() {
						if f := fw.Emit(SignalLostFocus, fw); f == enums.EVENT_STOP {
							fw = nil
						}
					}
				}
				if f := s.Emit(SignalGainedFocus, s, fw); f == enums.EVENT_STOP {
					if fw != nil {
						tl.SetFocus(fw)
					}
				}
				s.LogDebug("has taken focus")
			}
		}
	}
}

func (s *CScrollbar) CancelEvent() {
	if s.HasEventFocus() {
		s.ReleaseEventFocus()
	}
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
	s.Invalidate()
}

func (s *CScrollbar) ProcessEvent(evt cdk.Event) enums.EventFlag {
	s.Lock()
	defer s.Unlock()
	switch e := evt.(type) {
	case *cdk.EventMouse:
		point := ptypes.NewPoint2I(e.Position())
		return s.ProcessEventAtPoint(point, e)
	case *cdk.EventKey:
		if s.HasEventFocus() {
			s.CancelEvent()
			return enums.EVENT_STOP
		}
		step, page := 0, 0
		if adjustment := s.GetAdjustment(); adjustment != nil {
			step, page = adjustment.GetStepIncrement(), adjustment.GetPageIncrement()
		} else {
			s.LogError("missing adjustment")
		}
		switch s.orientation {
		case enums.ORIENTATION_HORIZONTAL:
			switch e.Key() {
			case cdk.KeyLeft:
				if e.Modifiers().Has(cdk.ModShift) {
					return s.Backward(page)
				}
				return s.Backward(step)
			case cdk.KeyRight:
				if e.Modifiers().Has(cdk.ModShift) {
					return s.Forward(page)
				}
				return s.Forward(step)
			}
		case enums.ORIENTATION_VERTICAL:
			fallthrough
		default:
			switch e.Key() {
			case cdk.KeyUp:
				return s.Backward(step)
			case cdk.KeyDown:
				return s.Forward(step)
			case cdk.KeyPgUp:
				return s.Backward(page)
			case cdk.KeyPgDn:
				return s.Forward(page)
			}
		}
	}
	return enums.EVENT_PASS
}

func (s *CScrollbar) ProcessEventAtPoint(p *ptypes.Point2I, e *cdk.EventMouse) enums.EventFlag {
	// me := NewMouseEvent(e)
	slider := s.GetSliderRegion()
	switch e.State() {
	case cdk.BUTTON_PRESS:
		if w := s.GetWidgetAt(p); w != nil && w.IsVisible() {
			if w.ObjectID() != s.ObjectID() {
				if wb, ok := w.(*CButton); ok {
					// s.GrabFocus()
					// s.GrabEventFocus()
					wb.SetPressed(true)
					s.focusedButton = wb
					return enums.EVENT_STOP
				}
			}
			if slider.HasPoint(*p) {
				s.prevSliderPos = p
				s.sliderMoving = true
				s.focusedButton = nil
				// s.GrabFocus()
				// s.GrabEventFocus()
				return enums.EVENT_STOP
			}
		}
	case cdk.DRAG_START:
		if !s.sliderMoving {
			s.focusedButton = nil
			s.sliderMoving = true
			// s.GrabFocus()
			// s.GrabEventFocus()
		}
		fallthrough
	case cdk.DRAG_MOVE:
		if s.sliderMoving {
			if s.prevSliderPos != nil {
				if s.prevSliderPos.X != p.X && s.orientation == enums.ORIENTATION_HORIZONTAL {
					// moved horizontally
					if s.textDirection == TextDirRtl {
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
					return enums.EVENT_STOP
				}
				if s.prevSliderPos.Y != p.Y && s.orientation == enums.ORIENTATION_VERTICAL {
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
					return enums.EVENT_STOP
				}
			}
			s.prevSliderPos = p
		}
	case cdk.DRAG_STOP:
		if s.HasEventFocus() {
			s.ReleaseEventFocus()
		}
		s.focusedButton = nil
		s.sliderMoving = false
		s.prevSliderPos = nil
		return enums.EVENT_STOP
	case cdk.BUTTON_RELEASE:
		if s.HasEventFocus() {
			s.ReleaseEventFocus()
		}
		if s.focusedButton != nil {
			if s.focusedButton.HasPoint(p) {
				s.focusedButton.SetPressed(false)
				s.focusedButton.Activate()
				s.focusedButton = nil
				s.sliderMoving = false
				s.prevSliderPos = nil
				return enums.EVENT_STOP
			}
		}
		s.focusedButton = nil
		s.sliderMoving = false
		s.prevSliderPos = nil
		slider := s.GetSliderRegion()
		if s.orientation == enums.ORIENTATION_HORIZONTAL {
			if s.textDirection == TextDirRtl {
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
	return enums.EVENT_PASS
}

func (s *CScrollbar) Invalidate() enums.EventFlag {
	origin := s.GetOrigin()
	size := ptypes.MakeRectangle(1, 1)
	theme := s.GetThemeRequest()
	style := theme.Content.Normal
	doStepper := func(b *CButton) {
		if b != nil {
			local := b.GetOrigin()
			local.SubPoint(origin)
			if err := memphis.ConfigureSurface(b.ObjectID(), local, size, style); err != nil {
				b.LogErr(err)
			}
			b.Invalidate()
		}
	}
	doStepper(s.forwardStepper)
	doStepper(s.backwardStepper)
	doStepper(s.secondaryForwardStepper)
	doStepper(s.secondaryBackwardStepper)
	return enums.EVENT_PASS
}

func (s *CScrollbar) draw(data []interface{}, argv ...interface{}) enums.EventFlag {
	if surface, ok := argv[1].(*memphis.CSurface); ok {
		s.Lock()
		defer s.Unlock()
		alloc := s.GetAllocation()
		if !s.IsVisible() || alloc.W <= 0 || alloc.H <= 0 {
			s.LogTrace("not visible, zero width or zero height")
			return enums.EVENT_PASS
		}
		theme := s.GetThemeRequest()
		origin := s.GetOrigin()
		// draw the trough
		trough := s.GetTroughRegion()
		trough.X -= origin.X
		trough.Y -= origin.Y
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
				b.Draw()
				return surface.Composite(b.ObjectID())
			}
			return nil
		}
		fwd, bwd, sFwd, sBwd := s.GetAllStepperRegions()
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
		return enums.EVENT_STOP
	}
	return enums.EVENT_PASS
}

func (s *CScrollbar) GetAllStepperRegions() (fwd, bwd, sFwd, sBwd ptypes.Region) {
	start, end := s.GetStepperRegions()
	fwd, bwd, sFwd, sBwd = end, start, start, end
	switch s.orientation {
	case enums.ORIENTATION_HORIZONTAL:
		if fwd.W == 2 {
			fwd.X += 1
			fwd.W = 1
		}
		if sFwd.W == 2 {
			sFwd.X += 1
			sFwd.W = 1
		}
	case enums.ORIENTATION_VERTICAL:
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
	return
}

func (s *CScrollbar) GetStepperRegions() (start, end ptypes.Region) {
	alloc := s.GetAllocation()
	origin := s.GetOrigin()
	switch s.orientation {
	case enums.ORIENTATION_HORIZONTAL:
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
	case enums.ORIENTATION_VERTICAL:
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
	return
}

func (s *CScrollbar) GetTroughRegion() (region ptypes.Region) {
	alloc := s.GetAllocation()
	start, end := s.GetStepperRegions()
	region = ptypes.MakeRegion(start.X, start.Y, 1, 1)
	switch s.orientation {
	case enums.ORIENTATION_HORIZONTAL:
		region.X += start.W
		region.W = alloc.W - start.W - end.W
	case enums.ORIENTATION_VERTICAL:
		fallthrough
	default:
		region.Y += start.H
		region.H = alloc.H - start.H - end.H
	}
	region.Floor(0, 0)
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
	region = ptypes.MakeRegion(trough.X, trough.Y, 1, 1)
	switch s.orientation {
	case enums.ORIENTATION_HORIZONTAL:
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
	case enums.ORIENTATION_VERTICAL:
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
	return
}

func (s *CScrollbar) resizeSteppers() {
	fwd, bwd, sFwd, sBwd := s.GetAllStepperRegions()
	aFwd, aBwd := ArrowDown, ArrowUp
	if s.orientation == enums.ORIENTATION_HORIZONTAL {
		aFwd, aBwd = ArrowRight, ArrowLeft
	}
	s.forwardStepper = s.resizeStepper(
		aFwd, aBwd,
		s.hasForwardStepper, s.forwardStepper,
		true,
		fwd.X, fwd.Y, fwd.W, fwd.H,
	)
	s.backwardStepper = s.resizeStepper(
		aFwd, aBwd,
		s.hasBackwardStepper, s.backwardStepper,
		false,
		bwd.X, bwd.Y, bwd.W, bwd.H,
	)
	s.secondaryForwardStepper = s.resizeStepper(
		aFwd, aBwd,
		s.hasSecondaryForwardStepper, s.secondaryForwardStepper,
		true,
		sFwd.X, sFwd.Y, sFwd.W, sFwd.H,
	)
	s.secondaryBackwardStepper = s.resizeStepper(
		aFwd, aBwd,
		s.hasSecondaryBackwardStepper, s.secondaryBackwardStepper,
		false,
		sBwd.X, bwd.Y, sBwd.W, sBwd.H,
	)
}

func (s *CScrollbar) resizeStepper(fArrow, bArrow ArrowType, has bool, b *CButton, forward bool, x, y, w, h int) *CButton {
	if has {
		if b == nil {
			if forward {
				fa := NewArrow(fArrow)
				fa.SetOrigin(0, 0)
				// fa.SetSizeRequest(1, 1)
				fa.SetAllocation(ptypes.MakeRectangle(1, 1))
				fa.SetTheme(DefaultColorButtonTheme)
				fa.UnsetFlags(CAN_FOCUS)
				fa.Show()
				b = NewButtonWithWidget(fa)
			} else {
				ba := NewArrow(bArrow)
				ba.SetOrigin(0, 0)
				// ba.SetSizeRequest(1, 1)
				ba.SetAllocation(ptypes.MakeRectangle(1, 1))
				ba.SetTheme(DefaultColorButtonTheme)
				ba.UnsetFlags(CAN_FOCUS)
				ba.Show()
				b = NewButtonWithWidget(ba)
			}
			b.ShowAll()
			b.SetFocusOnClick(false)
			b.Connect(
				SignalActivate,
				fmt.Sprintf("%v.activate", s.ObjectName()),
				func(data []interface{}, argv ...interface{}) enums.EventFlag {
					if adjustment := s.GetAdjustment(); adjustment != nil {
						if forward {
							s.Forward(adjustment.GetStepIncrement())
						} else {
							s.Backward(adjustment.GetStepIncrement())
						}
					} else {
						s.LogError("missing adjustment")
					}
					return enums.EVENT_STOP
				},
			)
		}
		if bc := b.GetChild(); bc != nil {
			if ba, ok := bc.(Arrow); ok {
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
		b.SetParent(s.GetParent())
		b.SetWindow(s.GetWindow())
		b.SetOrigin(x, y)
		// b.SetSizeRequest(w, h)
		b.SetAllocation(ptypes.MakeRectangle(w, h))
		b.ShowAll()
		b.Resize()
	}
	return b
}

func (s *CScrollbar) resizeSlider() {
	if s.slider == nil {
		l := NewLabel("*")
		l.SetSingleLineMode(true)
		l.SetMaxWidthChars(1)
		// l.SetSizeRequest(1, 1)
		l.Show()
		s.slider = NewButtonWithWidget(l)
		s.slider.Show()
		s.slider.SetTheme(s.GetTheme())
	}
	sr := s.GetSliderRegion()
	s.slider.SetOrigin(sr.X, sr.Y)
	// s.slider.SetSizeRequest(sr.W, sr.H)
	s.slider.SetAllocation(sr.Size())
	s.slider.Resize()
}

const ScrollbarDrawHandle = "scrollbar-draw-handler"
