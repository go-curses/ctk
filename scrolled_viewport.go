package ctk

import (
	"github.com/go-curses/cdk"
	"github.com/go-curses/cdk/lib/enums"
	cmath "github.com/go-curses/cdk/lib/math"
	"github.com/go-curses/cdk/lib/paint"
	"github.com/go-curses/cdk/lib/ptypes"
	"github.com/go-curses/cdk/memphis"
)

const (
	TypeScrolledViewport cdk.CTypeTag = "ctk-scrolled-viewport"
)

func init() {
	_ = cdk.TypesManager.AddType(TypeScrolledViewport, func() interface{} { return MakeScrolledViewport() }, "ctk-scrolled-window")
	ctkBuilderTranslators[TypeScrolledViewport] = func(builder Builder, widget Widget, name, value string) error {
		switch name {
		case "hscrollbar-policy", "h-scrollbar-policy":
			if err := widget.SetPropertyFromString(PropertyHScrollbarPolicy, value); err != nil {
				return err
			}
			return nil
		case "vscrollbar-policy", "v-scrollbar-policy":
			if err := widget.SetPropertyFromString(PropertyVScrollbarPolicy, value); err != nil {
				return err
			}
			return nil
		}
		return ErrFallthrough
	}
}

// ScrolledViewport Hierarchy:
//      Object
//        +- Widget
//          +- Container
//            +- Bin
//              +- ScrolledViewport
type ScrolledViewport interface {
	Viewport

	Init() (already bool)
	Build(builder Builder, element *CBuilderElement) error
	GetHAdjustment() (value *CAdjustment)
	GetVAdjustment() (value *CAdjustment)
	SetPolicy(hScrollbarPolicy PolicyType, vScrollbarPolicy PolicyType)
	AddWithViewport(child Widget)
	SetPlacement(windowPlacement CornerType)
	UnsetPlacement()
	SetShadowType(t ShadowType)
	SetHAdjustment(hAdjustment *CAdjustment)
	SetVAdjustment(vAdjustment *CAdjustment)
	GetPlacement() (value CornerType)
	GetPolicy() (hScrollbarPolicy PolicyType, vScrollbarPolicy PolicyType)
	GetShadowType() (value ShadowType)
	VerticalShowByPolicy() (show bool)
	HorizontalShowByPolicy() (show bool)
	Add(w Widget)
	Remove(w Widget)
	GetChild() Widget
	GetHScrollbar() *CHScrollbar
	GetVScrollbar() *CVScrollbar
	Show()
	Hide()
	GetWidgetAt(p *ptypes.Point2I) Widget
	CancelEvent()
	GetRegions() (c, h, v ptypes.Region)
	GrabFocus()
	GrabEventFocus()
}

type CScrolledViewport struct {
	CViewport

	scrollbarsWithinBevel bool
}

func MakeScrolledViewport() *CScrolledViewport {
	return NewScrolledViewport()
}

func NewScrolledViewport() *CScrolledViewport {
	s := new(CScrolledViewport)
	s.Init()
	return s
}

func (s *CScrolledViewport) Init() (already bool) {
	if s.InitTypeItem(TypeScrolledViewport, s) {
		return true
	}
	s.CViewport.Init()
	s.flags = NULL_WIDGET_FLAG
	s.SetFlags(SENSITIVE | CAN_FOCUS | APP_PAINTABLE)
	_ = s.InstallProperty(PropertyHScrollbarPolicy, cdk.StructProperty, true, PolicyAlways)
	_ = s.InstallProperty(PropertyScrolledViewportShadowType, cdk.StructProperty, true, SHADOW_NONE)
	_ = s.InstallProperty(PropertyVScrollbarPolicy, cdk.StructProperty, true, PolicyAlways)
	_ = s.InstallProperty(PropertyWindowPlacement, cdk.StructProperty, true, GravityNorthWest)
	_ = s.InstallProperty(PropertyWindowPlacementSet, cdk.BoolProperty, true, false)
	s.SetTheme(paint.DefaultColorTheme)
	s.SetPolicy(PolicyAlways, PolicyAlways)
	// hScrollbar
	s.CContainer.Add(NewHScrollbar())
	if hs := s.GetHScrollbar(); hs != nil {
		hs.SetParent(s)
		hs.SetWindow(s.GetWindow())
		s.SetHAdjustment(hs.GetAdjustment())
		hs.UnsetFlags(CAN_FOCUS)
	}
	// vScrollbar
	s.CContainer.Add(NewVScrollbar())
	if vs := s.GetVScrollbar(); vs != nil {
		vs.SetParent(s)
		vs.SetWindow(s.GetWindow())
		s.SetVAdjustment(vs.GetAdjustment())
		vs.UnsetFlags(CAN_FOCUS)
	}
	s.Connect(SignalCdkEvent, ScrolledViewportEventHandle, s.event)
	s.Connect(SignalInvalidate, ScrolledViewportDrawHandle, s.invalidate)
	s.Connect(SignalResize, ScrolledViewportDrawHandle, s.resize)
	s.Connect(SignalDraw, ScrolledViewportDrawHandle, s.draw)
	s.Invalidate()
	return false
}

func (s *CScrolledViewport) Build(builder Builder, element *CBuilderElement) error {
	s.Freeze()
	defer s.Thaw()
	if err := s.CObject.Build(builder, element); err != nil {
		return err
	}
	if len(element.Children) > 0 {
		topChild := element.Children[0]
		if topClass, ok := topChild.Attributes["class"]; ok {
			switch topClass {
			case "GtkViewport":
				// GtkScrolledWindow -> GtkViewport -> Thing
				if len(topChild.Children) > 0 {
					grandchild := topChild.Children[0]
					if newChild := builder.Build(grandchild); newChild != nil {
						if grandWidget, ok := grandchild.Instance.(Widget); ok {
							s.Add(grandWidget)
						} else {
							s.LogError("viewport grandchild is not of Widget type: %v (%T)", grandchild, grandchild)
						}
					}
				} else {
					s.LogError("viewport child has no descendants")
				}
			default:
				// GtkScrolledWindow -> ScrollableThing
				if newChild := builder.Build(topChild); newChild != nil {
					if newWidget, ok := newChild.(Widget); ok {
						s.Add(newWidget)
					}
				}
			}
		}
	}
	return nil
}

// Returns the horizontal scrollbar's adjustment, used to connect the
// horizontal scrollbar to the child widget's horizontal scroll
// functionality.
// Returns:
//      the horizontal Adjustment.
func (s *CScrolledViewport) GetHAdjustment() (value Adjustment) {
	var ok bool
	if v, err := s.GetStructProperty(PropertyHAdjustment); err != nil {
		s.LogErr(err)
	} else if value, ok = v.(Adjustment); !ok {
		s.LogError("value stored in struct property is not of Adjustment type: %v (%T)", v, v)
	}
	return
}

// Returns the vertical scrollbar's adjustment, used to connect the vertical
// scrollbar to the child widget's vertical scroll functionality.
// Returns:
//      the vertical Adjustment.
func (s *CScrolledViewport) GetVAdjustment() (value Adjustment) {
	var ok bool
	if v, err := s.GetStructProperty(PropertyVAdjustment); err != nil {
		s.LogErr(err)
	} else if value, ok = v.(Adjustment); !ok {
		s.LogError("value stored in struct property is not an Adjustment: %v (%T)", v, v)
	}
	return
}

// Sets the scrollbar policy for the horizontal and vertical scrollbars. The
// policy determines when the scrollbar should appear; it is a value from the
// PolicyType enumeration. If GTK_POLICY_ALWAYS, the scrollbar is always
// present; if GTK_POLICY_NEVER, the scrollbar is never present; if
// GTK_POLICY_AUTOMATIC, the scrollbar is present only if needed (that is, if
// the slider part of the bar would be smaller than the trough - the display
// is larger than the page size).
// Parameters:
//      hScrollbarPolicy        policy for horizontal bar
//      vScrollbarPolicy        policy for vertical bar
func (s *CScrolledViewport) SetPolicy(hScrollbarPolicy PolicyType, vScrollbarPolicy PolicyType) {
	if err := s.SetStructProperty(PropertyHScrollbarPolicy, hScrollbarPolicy); err != nil {
		s.LogErr(err)
	}
	if err := s.SetStructProperty(PropertyVScrollbarPolicy, vScrollbarPolicy); err != nil {
		s.LogErr(err)
	}
	return
}

// Used to add children without native scrolling capabilities. This is simply
// a convenience function; it is equivalent to adding the unscrollable child
// to a viewport, then adding the viewport to the scrolled window. If a child
// has native scrolling, use ContainerAdd instead of this function.
// The viewport scrolls the child by moving its Window, and takes the size
// of the child to be the size of its toplevel Window. This will be very
// wrong for most widgets that support native scrolling; for example, if you
// add a widget such as TreeView with a viewport, the whole widget will
// scroll, including the column headings. Thus, widgets with native scrolling
// support should not be used with the Viewport proxy. A widget supports
// scrolling natively if the set_scroll_adjustments_signal field in
// WidgetClass is non-zero, i.e. has been filled in with a valid signal
// identifier.
// Parameters:
//      child   the widget you want to scroll
func (s *CScrolledViewport) AddWithViewport(child Widget) {}

// Sets the placement of the contents with respect to the scrollbars for the
// scrolled window. The default is GTK_CORNER_TOP_LEFT, meaning the child is
// in the top left, with the scrollbars underneath and to the right. Other
// values in CornerType are GTK_CORNER_TOP_RIGHT, GTK_CORNER_BOTTOM_LEFT,
// and GTK_CORNER_BOTTOM_RIGHT. See also GetPlacement
// and UnsetPlacement.
// Parameters:
//      windowPlacement position of the child window
func (s *CScrolledViewport) SetPlacement(windowPlacement CornerType) {}

// Unsets the placement of the contents with respect to the scrollbars for
// the scrolled window. If no window placement is set for a scrolled window,
// it obeys the "gtk-scrolled-window-placement" XSETTING. See also
// SetPlacement and
// GetPlacement.
func (s *CScrolledViewport) UnsetPlacement() {}

// Changes the type of shadow drawn around the contents of scrolled_window .
// Parameters:
//      type    kind of shadow to draw around scrolled window contents
func (s *CScrolledViewport) SetShadowType(t ShadowType) {
	if err := s.SetStructProperty(PropertyScrolledViewportShadowType, t); err != nil {
		s.LogErr(err)
	}
}

// Sets the Adjustment for the horizontal scrollbar.
// Parameters:
//      hAdjustment     horizontal scroll adjustment
func (s *CScrolledViewport) SetHAdjustment(hAdjustment *CAdjustment) {
	if err := s.SetStructProperty(PropertyHAdjustment, hAdjustment); err != nil {
		s.LogErr(err)
	}
}

// Sets the Adjustment for the vertical scrollbar.
// Parameters:
//      vAdjustment     vertical scroll adjustment
func (s *CScrolledViewport) SetVAdjustment(vAdjustment *CAdjustment) {
	if err := s.SetStructProperty(PropertyVAdjustment, vAdjustment); err != nil {
		s.LogErr(err)
	}
}

// Gets the placement of the contents with respect to the scrollbars for the
// scrolled window. See SetPlacement.
// Returns:
//      the current placement value.
//      See also SetPlacement and
//      UnsetPlacement.
func (s *CScrolledViewport) GetPlacement() (value CornerType) {
	return
}

// Retrieves the current policy values for the horizontal and vertical
// scrollbars. See SetPolicy.
// Parameters:
//      hScrollbarPolicy        location to store the policy
// for the horizontal scrollbar, or NULL.
//      vScrollbarPolicy        location to store the policy
// for the vertical scrollbar, or NULL.
func (s *CScrolledViewport) GetPolicy() (hScrollbarPolicy PolicyType, vScrollbarPolicy PolicyType) {
	var ok bool
	if v, err := s.GetStructProperty(PropertyHScrollbarPolicy); err != nil {
		s.LogErr(err)
	} else if hScrollbarPolicy, ok = v.(PolicyType); !ok {
		s.LogError("value stored in struct property is not of PolicyType: %v (%T)", v, v)
	}
	ok = false
	if v, err := s.GetStructProperty(PropertyVScrollbarPolicy); err != nil {
		s.LogErr(err)
	} else if vScrollbarPolicy, ok = v.(PolicyType); !ok {
		s.LogError("value stored in struct property is not of PolicyType: %v (%T)", v, v)
	}
	return
}

// Gets the shadow type of the scrolled window. See
// SetShadowType.
// Returns:
//      the current shadow type
func (s *CScrolledViewport) GetShadowType() (value ShadowType) {
	var ok bool
	if v, err := s.GetStructProperty(PropertyScrolledViewportShadowType); err != nil {
		s.LogErr(err)
	} else if value, ok = v.(ShadowType); !ok {
		s.LogError("value stored in struct property is not of ShadowType: %v (%T)", v, v)
	}
	return
}

func (s *CScrolledViewport) VerticalShowByPolicy() (show bool) {
	vPolicy, _ := s.GetPolicy()
	if vertical := s.GetVAdjustment(); vertical != nil {
		show = vertical.ShowByPolicy(vPolicy)
		if !show && vPolicy == PolicyAutomatic && vertical.Moot() {
			if child := s.GetChild(); child != nil {
				childSize := ptypes.NewRectangle(child.GetSizeRequest())
				if childSize.H > 0 {
					alloc := s.GetAllocation()
					if childSize.H > alloc.H {
						show = true
					}
				}
			}
		}
	} else {
		s.LogError("missing vertical adjustment")
	}
	return
}

func (s *CScrolledViewport) HorizontalShowByPolicy() (show bool) {
	_, hPolicy := s.GetPolicy()
	if horizontal := s.GetHAdjustment(); horizontal != nil {
		show = horizontal.ShowByPolicy(hPolicy)
		if !show && hPolicy == PolicyAutomatic && horizontal.Moot() {
			if child := s.GetChild(); child != nil {
				childSize := ptypes.NewRectangle(child.GetSizeRequest())
				if childSize.W > 0 {
					alloc := s.GetAllocation()
					if childSize.W > alloc.W {
						show = true
					}
				}
			}
		}
	} else {
		s.LogError("missing horizontal adjustment")
	}
	return
}

func (s *CScrolledViewport) Add(w Widget) {
	if len(s.children) < 3 {
		s.CContainer.Add(w)
		s.Invalidate()
	} else {
		s.LogError("too many children for scrolled viewport")
	}
}

func (s *CScrolledViewport) Remove(w Widget) {
	s.CContainer.Remove(w)
	s.Invalidate()
}

func (s *CScrolledViewport) GetChild() Widget {
	for _, child := range s.GetChildren() {
		if _, ok := child.(Scrollbar); !ok {
			return child
		}
	}
	return nil
}

func (s *CScrolledViewport) GetHScrollbar() *CHScrollbar {
	for _, child := range s.GetChildren() {
		if v, ok := child.(*CHScrollbar); ok {
			return v
		}
	}
	return nil
}

func (s *CScrolledViewport) GetVScrollbar() *CVScrollbar {
	for _, child := range s.GetChildren() {
		if v, ok := child.(*CVScrollbar); ok {
			return v
		}
	}
	return nil
}

func (s *CScrolledViewport) Show() {
	s.CViewport.Show()
	if child := s.GetChild(); child != nil {
		child.Show()
	}
	if hs := s.GetHScrollbar(); hs != nil {
		hs.Show()
	}
	if vs := s.GetVScrollbar(); vs != nil {
		vs.Show()
	}
	s.Invalidate()
}

func (s *CScrolledViewport) Hide() {
	s.CViewport.Hide()
	if child := s.GetChild(); child != nil {
		child.Hide()
	}
	if hs := s.GetHScrollbar(); hs != nil {
		hs.Hide()
	}
	if vs := s.GetVScrollbar(); vs != nil {
		vs.Hide()
	}
	s.Invalidate()
}

func (s *CScrolledViewport) GetWidgetAt(p *ptypes.Point2I) Widget {
	if s.HasPoint(p) && s.IsVisible() {
		return s
	}
	return nil
}

func (s *CScrolledViewport) internalGetWidgetAt(p *ptypes.Point2I) Widget {
	if s.HasPoint(p) {
		if vs := s.GetVScrollbar(); vs != nil {
			if vs.HasPoint(p) {
				return vs
			}
		}
		if hs := s.GetHScrollbar(); hs != nil {
			if hs.HasPoint(p) {
				return hs
			}
		}
		if child := s.GetChild(); child != nil {
			if child.HasPoint(p) {
				return child
			}
		}
		return s
	}
	return nil
}

func (s *CScrolledViewport) CancelEvent() {
	if child := s.GetChild(); child != nil {
		if cs, ok := child.(Sensitive); ok {
			cs.CancelEvent()
		}
	}
	if vs := s.GetVScrollbar(); vs != nil {
		vs.CancelEvent()
	}
	if hs := s.GetHScrollbar(); hs != nil {
		hs.CancelEvent()
	}
	s.Invalidate()
}

// GrabFocus will take the focus of the associated Window if the Widget instance
// CanFocus(). Any previously focused Widget will emit a lost-focus signal and
// the newly focused Widget will emit a gained-focus signal. This method emits a
// grab-focus signal initially and if the listeners return EVENT_PASS, the
// changes are applied.
//
// Note that this method needs to be implemented within each Drawable that can
// be focused because of the golang interface system losing the concrete struct
// when a Widget interface reference is passed as a generic interface{}
// argument.
func (s *CScrolledViewport) GrabFocus() {
	if s.CanFocus() && s.IsVisible() && s.IsSensitive() {
		if r := s.Emit(SignalGrabFocus, s); r == enums.EVENT_PASS {
			if tl := s.GetWindow(); tl != nil {
				if focused := tl.GetFocus(); focused != nil {
					if fw, ok := focused.(Widget); ok && fw.ObjectID() != s.ObjectID() {
						fw.Emit(SignalLostFocus)
						fw.UnsetState(StateSelected)
						fw.LogDebug("has lost focus")
					}
				}
				tl.SetFocus(s)
				s.Emit(SignalGainedFocus)
				s.SetState(StateSelected)
				s.LogDebug("has taken focus")
			}
		}
	} else {
		s.LogError("cannot grab focus: can't focus, invisible or insensitive")
	}
}

// GrabEventFocus will emit a grab-event-focus signal and if all signal handlers
// return enums.EVENT_PASS will set the Button instance as the Window event
// focus handler.
//
// Note that this method needs to be implemented within each Drawable that can
// be focused because of the golang interface system losing the concrete struct
// when a Widget interface reference is passed as a generic interface{}
// argument.
func (b *CScrolledViewport) GrabEventFocus() {
	if window := b.GetWindow(); window != nil {
		if f := b.Emit(SignalGrabEventFocus, b, window); f == enums.EVENT_PASS {
			window.SetEventFocus(b)
		}
	} else {
		b.LogError("cannot grab focus: can't focus, invisible or insensitive")
	}
}

func (s *CScrolledViewport) event(data []interface{}, argv ...interface{}) enums.EventFlag {
	if evt, ok := argv[1].(cdk.Event); ok {
		s.Lock()
		defer s.Unlock()
		switch e := evt.(type) {
		case *cdk.EventMouse:
			if e.IsWheelImpulse() {
				s.GrabFocus()
				switch e.WheelImpulse() {
				case cdk.WheelUp:
					if vs := s.GetVScrollbar(); vs != nil {
						if f := vs.ForwardStep(); f == enums.EVENT_STOP {
							s.Invalidate()
							return enums.EVENT_STOP
						}
						return enums.EVENT_PASS
					}
				case cdk.WheelLeft:
					if hs := s.GetHScrollbar(); hs != nil {
						if f := hs.BackwardStep(); f == enums.EVENT_STOP {
							s.Invalidate()
							return enums.EVENT_STOP
						}
						return enums.EVENT_PASS
					}
				case cdk.WheelDown:
					if vs := s.GetVScrollbar(); vs != nil {
						if f := vs.BackwardStep(); f == enums.EVENT_STOP {
							s.Invalidate()
							return enums.EVENT_STOP
						}
					}
				case cdk.WheelRight:
					if hs := s.GetHScrollbar(); hs != nil {
						if f := hs.ForwardStep(); f == enums.EVENT_STOP {
							s.Invalidate()
							return enums.EVENT_STOP
						}
						return enums.EVENT_PASS
					}
				}
			}
			point := ptypes.NewPoint2I(e.Position())
			if f := s.processEventAtPoint(point, e); f == enums.EVENT_STOP {
				s.GrabFocus()
				return enums.EVENT_STOP
			}
		case *cdk.EventKey:
			if vs := s.GetVScrollbar(); vs != nil {
				if f := vs.ProcessEvent(evt); f == enums.EVENT_STOP {
					s.Invalidate()
					return enums.EVENT_STOP
				}
			}
			if hs := s.GetHScrollbar(); hs != nil {
				if f := hs.ProcessEvent(evt); f == enums.EVENT_STOP {
					s.Invalidate()
					return enums.EVENT_STOP
				}
			}
		}
	}
	return enums.EVENT_PASS
}

func (s *CScrolledViewport) processEventAtPoint(p *ptypes.Point2I, evt *cdk.EventMouse) enums.EventFlag {
	if w := s.internalGetWidgetAt(p); w != nil {
		if w.ObjectID() != s.ObjectID() {
			if ws, ok := w.(Sensitive); ok {
				if f := ws.ProcessEvent(evt); f == enums.EVENT_STOP {
					s.Invalidate()
					return enums.EVENT_STOP
				}
			}
		}
	}
	return enums.EVENT_PASS
}

// Returns a CDK Region for each of the viewport child space, horizontal and vertical
// scrollbar spaces.
func (s *CScrolledViewport) GetRegions() (c, h, v ptypes.Region) {
	if child := s.GetChild(); child != nil {
		o := child.GetOrigin()
		a := child.GetAllocation()
		c = ptypes.MakeRegion(o.X, o.Y, a.W, a.H)
	}
	if hs := s.GetHScrollbar(); hs != nil && s.HorizontalShowByPolicy() {
		o := hs.GetOrigin()
		a := hs.GetAllocation()
		h = ptypes.MakeRegion(o.X, o.Y, a.W, a.H)
	}
	if vs := s.GetVScrollbar(); vs != nil && s.VerticalShowByPolicy() {
		o := vs.GetOrigin()
		a := vs.GetAllocation()
		v = ptypes.MakeRegion(o.X, o.Y, a.W, a.H)
	}
	return
}

func (s *CScrolledViewport) resize(data []interface{}, argv ...interface{}) enums.EventFlag {
	// s.resizeViewport()
	// s.resizeScrollbars()
	s.Invalidate()
	return enums.EVENT_STOP
}

func (s *CScrolledViewport) draw(data []interface{}, argv ...interface{}) enums.EventFlag {
	if surface, ok := argv[1].(*memphis.CSurface); ok {
		s.Lock()
		defer s.Unlock()
		alloc := s.GetAllocation()
		if !s.IsVisible() || alloc.W <= 0 || alloc.H <= 0 {
			s.LogTrace("not visible, zero width or zero height")
			return enums.EVENT_PASS
		}
		child := s.GetChild()
		if child != nil {
			surface.BoxWithTheme(
				s.GetOrigin(),
				s.GetAllocation(),
				false,
				true,
				child.GetTheme(),
			)
			if f := child.Draw(); f == enums.EVENT_STOP {
				if err := surface.Composite(child.ObjectID()); err != nil {
					s.LogError("child composite error: %v", err)
				}
			}
		}
		if vs := s.GetVScrollbar(); child != nil && vs != nil && s.VerticalShowByPolicy() {
			if f := vs.Draw(); f == enums.EVENT_STOP {
				if err := surface.Composite(vs.ObjectID()); err != nil {
					s.LogError("vertical scrollbar composite error: %v", err)
				}
			}
		}
		if hs := s.GetHScrollbar(); child != nil && hs != nil && s.HorizontalShowByPolicy() {
			if f := hs.Draw(); f == enums.EVENT_STOP {
				if err := surface.Composite(hs.ObjectID()); err != nil {
					s.LogError("horizontal scrollbar composite error: %v", err)
				}
			}
		}
		if child != nil && s.VerticalShowByPolicy() && s.HorizontalShowByPolicy() {
			// fill in the corner gap between scrollbars
			_ = surface.SetRune(alloc.W-1, alloc.H-1, s.GetTheme().Content.FillRune, s.GetTheme().Content.Normal)
		}

		if debug, _ := s.GetBoolProperty(cdk.PropertyDebug); debug {
			surface.DebugBox(paint.ColorSilver, s.ObjectInfo())
		}
		return enums.EVENT_STOP
	}
	return enums.EVENT_PASS
}

func (s *CScrolledViewport) invalidate(data []interface{}, argv ...interface{}) enums.EventFlag {
	s.resizeViewport()
	s.resizeScrollbars()
	origin := s.GetOrigin()
	// alloc := s.GetAllocation()
	if child := s.GetChild(); child != nil {
		local := child.GetOrigin()
		local.SubPoint(origin)
		size := child.GetAllocation() // set by resizeViewport() call
		style := child.GetTheme().Content.Normal
		if err := memphis.ConfigureSurface(child.ObjectID(), local, size, style); err != nil {
			child.LogErr(err)
		}
	}
	if vs := s.GetVScrollbar(); vs != nil && s.VerticalShowByPolicy() {
		local := vs.GetOrigin()
		local.SubPoint(origin)
		if err := memphis.ConfigureSurface(vs.ObjectID(), local, vs.GetAllocation(), s.GetThemeRequest().Content.Normal); err != nil {
			vs.LogErr(err)
		}
		vs.Show()
	}
	if hs := s.GetHScrollbar(); hs != nil && s.HorizontalShowByPolicy() {
		local := hs.GetOrigin()
		local.SubPoint(origin)
		if err := memphis.ConfigureSurface(hs.ObjectID(), local, hs.GetAllocation(), s.GetThemeRequest().Content.Normal); err != nil {
			hs.LogErr(err)
		}
		hs.Show()
	}
	return enums.EVENT_STOP
}

func (s *CScrolledViewport) makeAdjustments() (region ptypes.Region, changed bool) {
	changed = false
	region = ptypes.MakeRegion(0, 0, 0, 0)
	origin := s.GetOrigin()
	alloc := s.GetAllocation()
	horizontal, vertical := s.GetHAdjustment(), s.GetVAdjustment()
	if alloc.W == 0 || alloc.H == 0 {
		if horizontal != nil {
			ohValue, ohLower, ohUpper, ohStepIncrement, ohPageIncrement, ohPageSize := horizontal.Settings()
			ah := []int{ohValue, ohLower, ohUpper, ohStepIncrement, ohPageIncrement, ohPageSize}
			bh := []int{0, 0, 0, 0, 0, 0}
			changed = cmath.EqInts(ah, bh)
			horizontal.Configure(0, 0, 0, 0, 0, 0)
		}
		if vertical != nil {
			ovValue, ovLower, ovUpper, ovStepIncrement, ovPageIncrement, ovPageSize := vertical.Settings()
			av := []int{ovValue, ovLower, ovUpper, ovStepIncrement, ovPageIncrement, ovPageSize}
			bv := []int{0, 0, 0, 0, 0, 0}
			changed = changed || cmath.EqInts(av, bv)
			vertical.Configure(0, 0, 0, 0, 0, 0)
		}
		return
	}
	hValue, hLower, hUpper, hStepIncrement, hPageIncrement, hPageSize := 0, 0, 0, 0, 0, 0
	vValue, vLower, vUpper, vStepIncrement, vPageIncrement, vPageSize := 0, 0, 0, 0, 0, 0
	if child := s.GetChild(); child != nil {
		size := ptypes.NewRectangle(child.GetSizeRequest())
		if size.W <= -1 { // auto
			size.W = alloc.W
			// if s.VerticalShowByPolicy() {
			// 	size.W -= 1
			// }
		}
		if size.H <= -1 { // auto
			size.H = alloc.H
			// if s.HorizontalShowByPolicy() {
			// 	size.H -= 1
			// }
		}
		if size.W >= alloc.W {
			hStepIncrement, hPageIncrement, hPageSize = 1, alloc.W/2, alloc.W
			if size.W >= alloc.W {
				overflow := size.W - alloc.W
				hLower, hUpper = 0, overflow
				if s.VerticalShowByPolicy() {
					hUpper += 1
				}
			} else {
				hLower, hUpper, hValue = 0, 0, 0
			}
			if horizontal != nil {
				hValue = cmath.ClampI(horizontal.GetValue(), hLower, hUpper)
			}
		}
		region.X = origin.X - hValue
		region.W = size.W
		if size.H >= alloc.H {
			vStepIncrement, vPageIncrement, vPageSize = 1, alloc.H/2, alloc.H
			if size.H >= alloc.H {
				overflow := size.H - alloc.H
				vLower, vUpper = 0, overflow
				if s.HorizontalShowByPolicy() {
					vUpper += 1
				}
			} else {
				vLower, vUpper, vValue = 0, 0, 0
			}
			if vertical != nil {
				vValue = cmath.ClampI(vertical.GetValue(), vLower, vUpper)
			}
		}
		region.Y = origin.Y - vValue
		region.H = size.H
	}
	// horizontal
	if horizontal != nil {
		ohValue, ohLower, ohUpper, ohStepIncrement, ohPageIncrement, ohPageSize := horizontal.Settings()
		ah := []int{ohValue, ohLower, ohUpper, ohStepIncrement, ohPageIncrement, ohPageSize}
		bh := []int{hValue, hLower, hUpper, hStepIncrement, hPageIncrement, hPageSize}
		if !cmath.EqInts(ah, bh) {
			changed = true
			horizontal.Configure(hValue, hLower, hUpper, hStepIncrement, hPageIncrement, hPageSize)
		}
	}
	// vertical
	if vertical != nil {
		ovValue, ovLower, ovUpper, ovStepIncrement, ovPageIncrement, ovPageSize := vertical.Settings()
		av := []int{ovValue, ovLower, ovUpper, ovStepIncrement, ovPageIncrement, ovPageSize}
		bv := []int{vValue, vLower, vUpper, vStepIncrement, vPageIncrement, vPageSize}
		if !cmath.EqInts(av, bv) {
			changed = true
			vertical.Configure(vValue, vLower, vUpper, vStepIncrement, vPageIncrement, vPageSize)
		}
	}
	return
}

func (s *CScrolledViewport) resizeViewport() enums.EventFlag {
	region, _ := s.makeAdjustments()
	if child := s.GetChild(); child != nil {
		child.SetOrigin(region.X, region.Y)
		child.SetAllocation(region.Size())
		return child.Resize()
	}
	return enums.EVENT_STOP
}

func (s *CScrolledViewport) resizeScrollbars() enums.EventFlag {
	origin := s.GetOrigin()
	alloc := s.GetAllocation()
	if hs := s.GetHScrollbar(); hs != nil {
		o := ptypes.MakePoint2I(origin.X, origin.Y+alloc.H-1)
		a := ptypes.MakeRectangle(alloc.W, 1)
		if s.VerticalShowByPolicy() {
			a.W -= 1
		}
		hs.SetOrigin(o.X, o.Y)
		hs.SetAllocation(a)
		if s.IsFocused() {
			hs.SetState(StateSelected)
		} else {
			hs.UnsetState(StateSelected)
		}
		hs.Resize()
	}
	if vs := s.GetVScrollbar(); vs != nil {
		o := ptypes.MakePoint2I(origin.X+alloc.W-1, origin.Y)
		a := ptypes.MakeRectangle(1, alloc.H)
		if s.HorizontalShowByPolicy() {
			a.H -= 1
		}
		vs.SetOrigin(o.X, o.Y)
		vs.SetAllocation(a)
		if s.IsFocused() {
			vs.SetState(StateSelected)
		} else {
			vs.UnsetState(StateSelected)
		}
		vs.Resize()
	}
	return enums.EVENT_STOP
}

func (s *CScrolledViewport) lostFocus([]interface{}, ...interface{}) enums.EventFlag {
	s.Invalidate()
	return enums.EVENT_PASS
}

func (s *CScrolledViewport) gainedFocus([]interface{}, ...interface{}) enums.EventFlag {
	s.Invalidate()
	return enums.EVENT_PASS
}

// The Adjustment for the horizontal position.
// Flags: Read / Write / Construct
const PropertyHAdjustment cdk.Property = "h-adjustment"

// When the horizontal scrollbar is displayed.
// Flags: Read / Write
// Default value: GTK_POLICY_ALWAYS
const PropertyHScrollbarPolicy cdk.Property = "h-scrollbar-policy"

// Style of bevel around the contents.
// Flags: Read / Write
// Default value: GTK_SHADOW_NONE
const PropertyScrolledViewportShadowType cdk.Property = "viewport-shadow-type"

// The Adjustment for the vertical position.
// Flags: Read / Write / Construct
const PropertyVAdjustment cdk.Property = "v-adjustment"

// When the vertical scrollbar is displayed.
// Flags: Read / Write
// Default value: GTK_POLICY_ALWAYS
const PropertyVScrollbarPolicy cdk.Property = "vscrollbar-policy"

// Where the contents are located with respect to the scrollbars. This
// property only takes effect if "window-placement-set" is TRUE.
// Flags: Read / Write
// Default value: GTK_CORNER_TOP_LEFT
const PropertyWindowPlacement cdk.Property = "window-placement"

// Whether "window-placement" should be used to determine the location of the
// contents with respect to the scrollbars. Otherwise, the
// "gtk-scrolled-window-placement" setting is used.
// Flags: Read / Write
// Default value: FALSE
const PropertyWindowPlacementSet cdk.Property = "window-placement-set"

// Listener function arguments:
//      arg1 DirectionType
const SignalMoveFocusOut cdk.Signal = "move-focus-out"

// The ::scroll-child signal is a which gets emitted when a keybinding that
// scrolls is pressed. The horizontal or vertical adjustment is updated which
// triggers a signal that the scrolled windows child may listen to and scroll
// itself.
const SignalScrollChild cdk.Signal = "scroll-child"

const ScrolledViewportLostFocusHandle = "scrolled-viewport-lost-focus-handler"
const ScrolledViewportGainedFocusHandle = "scrolled-viewport-gained-focus-handler"
const ScrolledViewportEventHandle = "scrolled-viewport-event-handler"
const ScrolledViewportInvalidateHandle = "scrolled-viewport-invalidate-handler"
const ScrolledViewportResizeHandle = "scrolled-viewport-resize-handler"
const ScrolledViewportDrawHandle = "scrolled-viewport-draw-handler"
