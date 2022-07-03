package ctk

import (
	"github.com/go-curses/cdk"
	cenums "github.com/go-curses/cdk/lib/enums"
	cmath "github.com/go-curses/cdk/lib/math"
	"github.com/go-curses/cdk/lib/paint"
	"github.com/go-curses/cdk/lib/ptypes"
	"github.com/go-curses/cdk/memphis"

	"github.com/go-curses/ctk/lib/enums"
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

	GetHAdjustment() (value Adjustment)
	GetVAdjustment() (value Adjustment)
	SetPolicy(hScrollbarPolicy enums.PolicyType, vScrollbarPolicy enums.PolicyType)
	AddWithViewport(child Widget)
	SetPlacement(windowPlacement enums.CornerType)
	UnsetPlacement()
	SetShadowType(t enums.ShadowType)
	SetHAdjustment(hAdjustment Adjustment)
	SetVAdjustment(vAdjustment Adjustment)
	GetPlacement() (value enums.CornerType)
	GetPolicy() (hScrollbarPolicy enums.PolicyType, vScrollbarPolicy enums.PolicyType)
	GetShadowType() (value enums.ShadowType)
	VerticalShowByPolicy() (show bool)
	HorizontalShowByPolicy() (show bool)
	GetHScrollbar() HScrollbar
	GetVScrollbar() VScrollbar
	CancelEvent()
	GetRegions() (c, h, v ptypes.Region)
	ScrollTo(child Widget)
}

var _ ScrolledViewport = (*CScrolledViewport)(nil)

type CScrolledViewport struct {
	CViewport

	scrollbarsWithinBevel bool
}

func MakeScrolledViewport() ScrolledViewport {
	return NewScrolledViewport()
}

func NewScrolledViewport() ScrolledViewport {
	s := new(CScrolledViewport)
	s.Init()
	return s
}

func (s *CScrolledViewport) Init() (already bool) {
	if s.InitTypeItem(TypeScrolledViewport, s) {
		return true
	}
	s.CViewport.Init()
	s.flags = enums.NULL_WIDGET_FLAG
	s.SetFlags(enums.SENSITIVE | enums.CAN_FOCUS | enums.APP_PAINTABLE)
	s.SetTheme(paint.DefaultColorTheme)

	_ = s.InstallProperty(PropertyHScrollbarPolicy, cdk.StructProperty, true, enums.PolicyAlways)
	_ = s.InstallProperty(PropertyScrolledViewportShadowType, cdk.StructProperty, true, enums.SHADOW_NONE)
	_ = s.InstallProperty(PropertyVScrollbarPolicy, cdk.StructProperty, true, enums.PolicyAlways)
	_ = s.InstallProperty(PropertyWindowPlacement, cdk.StructProperty, true, enums.GravityNorthWest)
	_ = s.InstallProperty(PropertyWindowPlacementSet, cdk.BoolProperty, true, false)

	s.SetPolicy(enums.PolicyAlways, enums.PolicyAlways)

	// hScrollbar
	s.CContainer.Add(NewHScrollbar())
	if hs := s.GetHScrollbar(); hs != nil {
		hs.Show()
		hs.SetParent(s)
		hs.SetWindow(s.GetWindow())
		ha := hs.GetAdjustment()
		s.SetHAdjustment(ha)
		hs.UnsetFlags(enums.CAN_FOCUS)
		hs.Connect(SignalValueChanged, "value-changed-handler", func(data []interface{}, argv ...interface{}) cenums.EventFlag {
			s.LogDebug("hs value changed: %v", ha.GetValue())
			s.Resize()
			s.RequestDrawAndShow()
			return cenums.EVENT_PASS
		})
	}

	// vScrollbar
	s.CContainer.Add(NewVScrollbar())
	if vs := s.GetVScrollbar(); vs != nil {
		vs.Show()
		vs.SetParent(s)
		vs.SetWindow(s.GetWindow())
		va := vs.GetAdjustment()
		s.SetVAdjustment(va)
		vs.UnsetFlags(enums.CAN_FOCUS)
		vs.Connect(SignalValueChanged, "value-changed-handler", func(data []interface{}, argv ...interface{}) cenums.EventFlag {
			s.LogDebug("vs value changed: %v", va.GetValue())
			s.Resize()
			s.RequestDrawAndShow()
			return cenums.EVENT_PASS
		})
	}

	s.Connect(SignalCdkEvent, ScrolledViewportEventHandle, s.event)
	s.Connect(SignalLostFocus, ScrolledViewportLostFocusHandle, s.lostFocus)
	s.Connect(SignalGainedFocus, ScrolledViewportGainedFocusHandle, s.gainedFocus)
	s.Connect(SignalInvalidate, ScrolledViewportDrawHandle, s.invalidate)
	s.Connect(SignalResize, ScrolledViewportDrawHandle, s.resize)
	s.Connect(SignalDraw, ScrolledViewportDrawHandle, s.draw)
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
	} else if value, ok = v.(*CAdjustment); !ok {
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
	} else if value, ok = v.(*CAdjustment); !ok {
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
func (s *CScrolledViewport) SetPolicy(hScrollbarPolicy enums.PolicyType, vScrollbarPolicy enums.PolicyType) {
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
func (s *CScrolledViewport) SetPlacement(windowPlacement enums.CornerType) {}

// Unsets the placement of the contents with respect to the scrollbars for
// the scrolled window. If no window placement is set for a scrolled window,
// it obeys the "gtk-scrolled-window-placement" XSETTING. See also
// SetPlacement and
// GetPlacement.
func (s *CScrolledViewport) UnsetPlacement() {}

// Changes the type of shadow drawn around the contents of scrolled_window .
// Parameters:
//      type    kind of shadow to draw around scrolled window contents
func (s *CScrolledViewport) SetShadowType(t enums.ShadowType) {
	if err := s.SetStructProperty(PropertyScrolledViewportShadowType, t); err != nil {
		s.LogErr(err)
	}
}

// Sets the Adjustment for the horizontal scrollbar.
// Parameters:
//      hAdjustment     horizontal scroll adjustment
func (s *CScrolledViewport) SetHAdjustment(hAdjustment Adjustment) {
	if err := s.SetStructProperty(PropertyHAdjustment, hAdjustment); err != nil {
		s.LogErr(err)
	} else {
		if vs := s.GetVScrollbar(); vs != nil {
			vs.SetScrollAdjustments(hAdjustment, s.GetVAdjustment())
		}
		if hs := s.GetHScrollbar(); hs != nil {
			hs.SetScrollAdjustments(hAdjustment, s.GetVAdjustment())
		}
	}
}

// Sets the Adjustment for the vertical scrollbar.
// Parameters:
//      vAdjustment     vertical scroll adjustment
func (s *CScrolledViewport) SetVAdjustment(vAdjustment Adjustment) {
	if err := s.SetStructProperty(PropertyVAdjustment, vAdjustment); err != nil {
		s.LogErr(err)
	} else {
		if vs := s.GetVScrollbar(); vs != nil {
			vs.SetScrollAdjustments(s.GetHAdjustment(), vAdjustment)
		}
		if hs := s.GetHScrollbar(); hs != nil {
			hs.SetScrollAdjustments(s.GetHAdjustment(), vAdjustment)
		}
	}
}

// Gets the placement of the contents with respect to the scrollbars for the
// scrolled window. See SetPlacement.
// Returns:
//      the current placement value.
//      See also SetPlacement and
//      UnsetPlacement.
func (s *CScrolledViewport) GetPlacement() (value enums.CornerType) {
	return
}

// Retrieves the current policy values for the horizontal and vertical
// scrollbars. See SetPolicy.
// Parameters:
//      hScrollbarPolicy        location to store the policy
// for the horizontal scrollbar, or NULL.
//      vScrollbarPolicy        location to store the policy
// for the vertical scrollbar, or NULL.
func (s *CScrolledViewport) GetPolicy() (hScrollbarPolicy enums.PolicyType, vScrollbarPolicy enums.PolicyType) {
	var ok bool
	if v, err := s.GetStructProperty(PropertyHScrollbarPolicy); err != nil {
		s.LogErr(err)
	} else if hScrollbarPolicy, ok = v.(enums.PolicyType); !ok {
		s.LogError("value stored in struct property is not of PolicyType: %v (%T)", v, v)
	}
	ok = false
	if v, err := s.GetStructProperty(PropertyVScrollbarPolicy); err != nil {
		s.LogErr(err)
	} else if vScrollbarPolicy, ok = v.(enums.PolicyType); !ok {
		s.LogError("value stored in struct property is not of PolicyType: %v (%T)", v, v)
	}
	return
}

// Gets the shadow type of the scrolled window. See
// SetShadowType.
// Returns:
//      the current shadow type
func (s *CScrolledViewport) GetShadowType() (value enums.ShadowType) {
	var ok bool
	if v, err := s.GetStructProperty(PropertyScrolledViewportShadowType); err != nil {
		s.LogErr(err)
	} else if value, ok = v.(enums.ShadowType); !ok {
		s.LogError("value stored in struct property is not of ShadowType: %v (%T)", v, v)
	}
	return
}

func (s *CScrolledViewport) VerticalShowByPolicy() (show bool) {
	vPolicy, _ := s.GetPolicy()
	vertical := s.GetVAdjustment()
	child := s.GetChild()
	alloc := s.GetAllocation()
	s.RLock()
	if vertical != nil {
		show = vertical.ShowByPolicy(vPolicy)
		if !show && vPolicy == enums.PolicyAutomatic && vertical.Moot() {
			if child != nil {
				childSize := ptypes.NewRectangle(child.GetSizeRequest())
				if childSize.H > 0 {
					if childSize.H > alloc.H {
						show = true
					}
				}
			}
		}
	} else {
		s.LogError("missing vertical adjustment")
	}
	s.RUnlock()
	return
}

func (s *CScrolledViewport) HorizontalShowByPolicy() (show bool) {
	_, hPolicy := s.GetPolicy()
	horizontal := s.GetHAdjustment()
	child := s.GetChild()
	alloc := s.GetAllocation()
	s.RLock()
	if horizontal != nil {
		show = horizontal.ShowByPolicy(hPolicy)
		if !show && hPolicy == enums.PolicyAutomatic && horizontal.Moot() {
			if child != nil {
				childSize := ptypes.NewRectangle(child.GetSizeRequest())
				if childSize.W > 0 {
					if childSize.W > alloc.W {
						show = true
					}
				}
			}
		}
	} else {
		s.LogError("missing horizontal adjustment")
	}
	s.RUnlock()
	return
}

func (s *CScrolledViewport) Add(w Widget) {
	if _, ok := w.Self().(Scrollbar); ok {
		s.LogError("cannot Add a scrollbar as the Viewport content: %v", w)
		return
	}
	for _, child := range s.GetChildren() {
		if _, ok := child.(Scrollbar); !ok {
			s.CContainer.Remove(child)
		}
	}
	s.CContainer.Add(w)
	s.Invalidate()
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

func (s *CScrolledViewport) GetHScrollbar() HScrollbar {
	for _, child := range s.GetChildren() {
		if v, ok := child.Self().(*CHScrollbar); ok {
			return v
		}
	}
	return nil
}

func (s *CScrolledViewport) GetVScrollbar() VScrollbar {
	for _, child := range s.GetChildren() {
		if v, ok := child.Self().(*CVScrollbar); ok {
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
		if s.VerticalShowByPolicy() {
			if vs := s.GetVScrollbar(); vs != nil {
				if vs.HasPoint(p) {
					return vs
				}
			}
		}

		if s.HorizontalShowByPolicy() {
			if hs := s.GetHScrollbar(); hs != nil {
				if hs.HasPoint(p) {
					return hs
				}
			}
		}

		if child := s.GetChild(); child != nil {
			if child.HasPoint(p) {
				if found := child.GetWidgetAt(p); found != nil {
					return found
				}
			}
		}
		return s
	}
	return nil
}

func (s *CScrolledViewport) processEventAtPoint(p *ptypes.Point2I, evt *cdk.EventMouse) cenums.EventFlag {
	if w := s.internalGetWidgetAt(p); w != nil {
		if w.ObjectID() != s.ObjectID() {
			if ws, ok := w.Self().(Sensitive); ok {
				if f := ws.ProcessEvent(evt); f == cenums.EVENT_STOP {
					s.Invalidate()
					return cenums.EVENT_STOP
				}
			}
		}
	}
	return cenums.EVENT_PASS
}

func (s *CScrolledViewport) CancelEvent() {
	if child := s.GetChild(); child != nil {
		if cs, ok := child.Self().(Sensitive); ok {
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

func (s *CScrolledViewport) SetWindow(w Window) {
	if oldWindow := s.GetWindow(); oldWindow != nil {
		_ = oldWindow.Disconnect(SignalFocusChanged, ScrolledViewportWindowFocusHandle)
	}
	s.CBin.SetWindow(w)
	if w != nil {
		w.Connect(SignalFocusChanged, ScrolledViewportWindowFocusHandle, s.windowFocusSet)
	}
}

func (s *CScrolledViewport) ScrollTo(target Widget) {
	if !s.HasChild(target) {
		return
	}

	origin := s.GetOrigin()
	alloc := s.GetAllocation()
	region := ptypes.MakeRegion(origin.X, origin.Y, alloc.W, alloc.H)
	farPoint := region.FarPoint()
	childOrigin := target.GetOrigin()
	childAlloc := target.GetAllocation()
	childRegion := ptypes.MakeRegion(childOrigin.X, childOrigin.Y, childAlloc.W, childAlloc.H)
	childFarPoint := childRegion.FarPoint()

	if region.ContainsRegion(childRegion) {
		s.LogDebug("already visible: %v", target.ObjectInfo())
		return
	}

	// if region.OccludesRegion(childRegion) {
	// 	s.LogDebug("partially visible: %v", target.ObjectInfo())
	// 	return
	// }

	s.LogDebug("scrolling move position to include: %v within %v for %v", childRegion, region, target)

	wantResize := false

	if s.HorizontalShowByPolicy() {
		// check for left aligning, x+w visible
		horizontal := s.GetHAdjustment()
		upper := horizontal.GetUpper()
		delta := childFarPoint.Y - farPoint.Y + horizontal.GetValue()
		if delta > upper {
			delta = upper
		}
		horizontal.SetValue(delta)
		s.LogDebug("set horizontal delta: %v", delta)
		wantResize = true
	}

	if s.VerticalShowByPolicy() {
		// child for top aligning, y+h visible
		vertical := s.GetVAdjustment()
		upper := vertical.GetUpper()
		delta := childFarPoint.Y - farPoint.Y + vertical.GetValue()
		if delta > upper {
			delta = upper
		}
		vertical.SetValue(delta)
		s.LogDebug("set vertical delta: %v", delta)
		wantResize = true
	}

	if wantResize {
		s.Resize()
	}
}

func (s *CScrolledViewport) windowFocusSet(data []interface{}, argv ...interface{}) cenums.EventFlag {
	if len(argv) >= 2 {
		if focus, ok := argv[1].(Widget); ok {
			s.ScrollTo(focus)
		}
	}
	return cenums.EVENT_PASS
}

func (s *CScrolledViewport) event(data []interface{}, argv ...interface{}) cenums.EventFlag {
	if evt, ok := argv[1].(cdk.Event); ok {
		switch e := evt.(type) {
		case *cdk.EventMouse:
			if e.IsWheelImpulse() {
				vs := s.GetVScrollbar()
				hs := s.GetHScrollbar()
				switch e.WheelImpulse() {
				case cdk.WheelUp:
					if vs != nil {
						if f := vs.ForwardPage(); f == cenums.EVENT_STOP {
							s.Invalidate()
							return cenums.EVENT_STOP
						}
					}
				case cdk.WheelLeft:
					if hs != nil {
						if f := hs.BackwardStep(); f == cenums.EVENT_STOP {
							s.Invalidate()
							return cenums.EVENT_STOP
						}
					}
				case cdk.WheelDown:
					if vs != nil {
						if f := vs.BackwardPage(); f == cenums.EVENT_STOP {
							s.Invalidate()
							return cenums.EVENT_STOP
						}
					}
				case cdk.WheelRight:
					if hs != nil {
						if f := hs.ForwardStep(); f == cenums.EVENT_STOP {
							s.Invalidate()
							return cenums.EVENT_STOP
						}
					}
				}
			}
			point := ptypes.NewPoint2I(e.Position())
			if f := s.processEventAtPoint(point, e); f == cenums.EVENT_STOP {
				s.Invalidate()
				return cenums.EVENT_STOP
			}
		case *cdk.EventKey:
			if vs := s.GetVScrollbar(); vs != nil {
				if f := vs.ProcessEvent(evt); f == cenums.EVENT_STOP {
					s.Invalidate()
					return cenums.EVENT_STOP
				}
			}
			if hs := s.GetHScrollbar(); hs != nil {
				if f := hs.ProcessEvent(evt); f == cenums.EVENT_STOP {
					s.Invalidate()
					return cenums.EVENT_STOP
				}
			}
		}
	}
	return cenums.EVENT_PASS
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

func (s *CScrolledViewport) invalidate(data []interface{}, argv ...interface{}) cenums.EventFlag {
	child := s.GetChild()
	vs := s.GetVScrollbar()
	hs := s.GetHScrollbar()
	horizontalShow := s.HorizontalShowByPolicy()
	verticalShow := s.VerticalShowByPolicy()

	if child != nil {
		WidgetRecurseInvalidate(child)
	}

	if vs != nil {
		if verticalShow {
			vs.SetFlags(enums.VISIBLE)
			WidgetRecurseInvalidate(vs)
		} else if vs.IsVisible() {
			vs.UnsetFlags(enums.VISIBLE)
		}
	}

	if hs != nil {
		if horizontalShow {
			hs.SetFlags(enums.VISIBLE)
			WidgetRecurseInvalidate(hs)
		} else if hs.IsVisible() {
			hs.UnsetFlags(enums.VISIBLE)
		}
	}
	return cenums.EVENT_PASS
}

func (s *CScrolledViewport) resize(data []interface{}, argv ...interface{}) cenums.EventFlag {

	s.resizeViewport()
	s.resizeScrollbars()

	s.Invalidate()
	return cenums.EVENT_STOP
}

func (s *CScrolledViewport) draw(data []interface{}, argv ...interface{}) cenums.EventFlag {

	if surface, ok := argv[1].(*memphis.CSurface); ok {
		alloc := s.GetAllocation()
		if !s.IsVisible() || alloc.W <= 0 || alloc.H <= 0 {
			s.LogTrace("not visible, zero width or zero height")
			return cenums.EVENT_PASS
		}

		child := s.GetChild()
		vs := s.GetVScrollbar()
		hs := s.GetHScrollbar()
		verticalShow := s.VerticalShowByPolicy()
		horizontalShow := s.HorizontalShowByPolicy()

		theme := s.GetThemeRequest()
		surface.Fill(theme)

		if child != nil {
			child.Draw()
			if err := surface.Composite(child.ObjectID()); err != nil {
				s.LogError("child composite error: %v", err)
			}

			if vs != nil && vs.IsVisible() && verticalShow {
				vs.Draw()
				if err := surface.Composite(vs.ObjectID()); err != nil {
					s.LogError("vertical scrollbar composite error: %v", err)
				}
			}

			if hs != nil && hs.IsVisible() && horizontalShow {
				hs.Draw()
				if err := surface.Composite(hs.ObjectID()); err != nil {
					s.LogError("horizontal scrollbar composite error: %v", err)
				}
			}

			if verticalShow && horizontalShow {
				// fill in the corner gap between scrollbars
				_ = surface.SetRune(alloc.W-1, alloc.H-1, theme.Content.FillRune, theme.Content.Normal)
			}
		}

		if debug, _ := s.GetBoolProperty(cdk.PropertyDebug); debug {
			surface.DebugBox(paint.ColorSilver, s.ObjectInfo())
		}

		return cenums.EVENT_STOP
	}
	return cenums.EVENT_PASS
}

func (s *CScrolledViewport) makeAdjustments() (region ptypes.Region, changed bool) {
	origin := s.GetOrigin()
	alloc := s.GetAllocation()
	child := s.GetChild()
	verticalShow := s.VerticalShowByPolicy()
	horizontalShow := s.HorizontalShowByPolicy()
	horizontal, vertical := s.GetHAdjustment(), s.GetVAdjustment()

	s.Lock()
	defer s.Unlock()

	changed = false
	region = ptypes.MakeRegion(0, 0, 0, 0)

	if alloc.W == 0 || alloc.H == 0 {
		if horizontal != nil {
			ohValue, ohLower, ohUpper, ohStepIncrement, ohPageIncrement, ohPageSize := horizontal.Settings()
			ah := []int{ohValue, ohLower, ohUpper}
			bh := []int{0, 0, 0}
			changed = cmath.EqInts(ah, bh)
			horizontal.Configure(0, 0, 0, ohStepIncrement, ohPageIncrement, ohPageSize)
		}
		if vertical != nil {
			ovValue, ovLower, ovUpper, ovStepIncrement, ovPageIncrement, ovPageSize := vertical.Settings()
			av := []int{ovValue, ovLower, ovUpper}
			bv := []int{0, 0, 0}
			changed = changed || cmath.EqInts(av, bv)
			vertical.Configure(0, 0, 0, ovStepIncrement, ovPageIncrement, ovPageSize)
		}
		return
	}

	hValue, hLower, hUpper := 0, 0, 0
	vValue, vLower, vUpper := 0, 0, 0
	if child != nil {
		size := ptypes.NewRectangle(child.GetSizeRequest())
		if size.W <= -1 { // auto
			size.W = alloc.W
		}
		if size.H <= -1 { // auto
			size.H = alloc.H
		}
		if size.W >= alloc.W {
			if size.W >= alloc.W {
				overflow := size.W - alloc.W
				hLower, hUpper = 0, overflow
				if verticalShow {
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
			if size.H >= alloc.H {
				overflow := size.H - alloc.H
				vLower, vUpper = 0, overflow
				if horizontalShow {
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

	if horizontal != nil {
		ohValue, ohLower, ohUpper, ohStepIncrement, ohPageIncrement, ohPageSize := horizontal.Settings()
		ah := []int{ohValue, ohLower, ohUpper}
		bh := []int{hValue, hLower, hUpper}
		if !cmath.EqInts(ah, bh) {
			changed = true
			horizontal.Configure(hValue, hLower, hUpper, ohStepIncrement, ohPageIncrement, ohPageSize)
		}
	}

	if vertical != nil {
		ovValue, ovLower, ovUpper, ovStepIncrement, ovPageIncrement, ovPageSize := vertical.Settings()
		av := []int{ovValue, ovLower, ovUpper}
		bv := []int{vValue, vLower, vUpper}
		if !cmath.EqInts(av, bv) {
			changed = true
			vertical.Configure(vValue, vLower, vUpper, ovStepIncrement, ovPageIncrement, ovPageSize)
		}
	}

	return
}

func (s *CScrolledViewport) resizeViewport() cenums.EventFlag {
	region, _ := s.makeAdjustments()
	child := s.GetChild()
	if child != nil {
		child.SetOrigin(region.X, region.Y)
		child.SetAllocation(region.Size())
		child.Resize()
	}
	return cenums.EVENT_STOP
}

func (s *CScrolledViewport) resizeScrollbars() cenums.EventFlag {
	origin := s.GetOrigin()
	alloc := s.GetAllocation()
	hs := s.GetHScrollbar()
	vs := s.GetVScrollbar()
	verticalShow := s.VerticalShowByPolicy()
	horizontalShow := s.HorizontalShowByPolicy()
	state := s.GetState()
	if hs != nil {
		o := ptypes.MakePoint2I(origin.X, origin.Y+alloc.H-1)
		a := ptypes.MakeRectangle(alloc.W, 1)
		if verticalShow {
			a.W -= 1
		}
		hs.SetOrigin(o.X, o.Y)
		hs.SetAllocation(a)
		hs.SetState(enums.StateNone)
		hs.SetState(state)
		hs.Resize()
	}
	if vs != nil {
		o := ptypes.MakePoint2I(origin.X+alloc.W-1, origin.Y)
		a := ptypes.MakeRectangle(1, alloc.H)
		if horizontalShow {
			a.H -= 1
		}
		vs.SetOrigin(o.X, o.Y)
		vs.SetAllocation(a)
		vs.SetState(enums.StateNone)
		vs.SetState(state)
		vs.Resize()
	}
	return cenums.EVENT_STOP
}

func (s *CScrolledViewport) lostFocus([]interface{}, ...interface{}) cenums.EventFlag {
	s.UnsetState(enums.StateSelected)
	if hs := s.GetHScrollbar(); hs != nil {
		hs.UnsetState(enums.StateSelected)
		WidgetRecurseInvalidate(hs)
	}
	if vs := s.GetVScrollbar(); vs != nil {
		vs.UnsetState(enums.StateSelected)
		WidgetRecurseInvalidate(vs)
	}
	s.Invalidate()
	return cenums.EVENT_PASS
}

func (s *CScrolledViewport) gainedFocus([]interface{}, ...interface{}) cenums.EventFlag {
	s.SetState(enums.StateSelected)
	if hs := s.GetHScrollbar(); hs != nil {
		hs.SetState(enums.StateSelected)
		WidgetRecurseInvalidate(hs)
	}
	if vs := s.GetVScrollbar(); vs != nil {
		vs.SetState(enums.StateSelected)
		WidgetRecurseInvalidate(vs)
	}
	s.Invalidate()
	return cenums.EVENT_PASS
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
const ScrolledViewportWindowFocusHandle = "scrolled-viewport-window-focus-handler"