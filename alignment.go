package ctk

import (
	"github.com/go-curses/cdk"
	"github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/cdk/lib/paint"
	"github.com/go-curses/cdk/lib/ptypes"
	"github.com/go-curses/cdk/memphis"
)

const TypeAlignment cdk.CTypeTag = "ctk-alignment"

func init() {
	_ = cdk.TypesManager.AddType(TypeAlignment, func() interface{} { return MakeAlignment() })
}

// Alignment Hierarchy:
//	Object
//	  +- Widget
//	    +- Container
//	      +- Bin
//	        +- Alignment
//
// The Alignment widget controls the alignment and size of its child widget.
// It has four settings: xScale, yScale, xAlign, and yAlign. The scale
// settings are used to specify how much the child widget should expand to
// fill the space allocated to the Alignment. The values can range from 0
// (meaning the child doesn't expand at all) to 1 (meaning the child expands
// to fill all of the available space). The alignment settings are used to place
// the child widget within the available area. The values range from 0 (top
// or left) to 1 (bottom or right). Of course, if the scale settings are both
// set to 1, the alignment settings have no effect. New Alignment instances can
// be created using NewAlignment.
type Alignment interface {
	Bin
	Buildable

	Init() (already bool)
	Get() (xAlign, yAlign, xScale, yScale float64)
	Set(xAlign, yAlign, xScale, yScale float64)
	GetPadding() (paddingTop, paddingBottom, paddingLeft, paddingRight int)
	SetPadding(paddingTop, paddingBottom, paddingLeft, paddingRight int)
	Add(w Widget)
	Remove(w Widget)
}

// The CAlignment structure implements the Alignment interface and is
// exported to facilitate type embedding with custom implementations. No member
// variables are exported as the interface methods are the only intended means
// of interacting with Alignment objects.
type CAlignment struct {
	CBin
}

// MakeAlignment is used by the Buildable system to construct a new Alignment
// with default settings of: xAlign=0.5, yAlign=1.0, xScale=0.5, yScale=1.0
func MakeAlignment() *CAlignment {
	return NewAlignment(0.5, 1, 0.5, 1)
}

// NewAlignment is the constructor for new Alignment instances.
//
// Parameters:
// 	xAlign	the horizontal alignment of the child widget, from 0 (left) to 1 (right)
// 	yAlign	the vertical alignment of the child widget, from 0 (top) to 1 (bottom)
// 	xScale	the amount that the child widget expands horizontally to fill up unused space, from 0 to 1. A value of 0 indicates that the child widget should never expand. A value of 1 indicates that the child widget will expand to fill all of the space allocated for the Alignment
// 	yScale	the amount that the child widget expands vertically to fill up unused space, from 0 to 1. The values are similar to xScale
func NewAlignment(xAlign float64, yAlign float64, xScale float64, yScale float64) *CAlignment {
	a := new(CAlignment)
	a.Init()
	a.Set(xAlign, yAlign, xScale, yScale)
	return a
}

// Init initializes an Alignment object. This must be called at least once to
// set up the necessary defaults and allocate any memory structures. Calling
// this more than once is safe though unnecessary. Only the first call will
// result in any effect upon the Alignment instance. Init is used in the
// NewAlignment constructor and only necessary when implementing a derivative
// Alignment type.
func (a *CAlignment) Init() (already bool) {
	if a.InitTypeItem(TypeAlignment, a) {
		return true
	}
	a.CBin.Init()
	a.flags = NULL_WIDGET_FLAG
	a.SetFlags(PARENT_SENSITIVE | APP_PAINTABLE)
	_ = a.InstallBuildableProperty(PropertyBottomPadding, cdk.IntProperty, true, 0)
	_ = a.InstallBuildableProperty(PropertyLeftPadding, cdk.IntProperty, true, 0)
	_ = a.InstallBuildableProperty(PropertyRightPadding, cdk.IntProperty, true, 0)
	_ = a.InstallBuildableProperty(PropertyTopPadding, cdk.IntProperty, true, 0)
	_ = a.InstallBuildableProperty(PropertyXAlign, cdk.FloatProperty, true, 0.5)
	_ = a.InstallBuildableProperty(PropertyXScale, cdk.FloatProperty, true, 1)
	_ = a.InstallBuildableProperty(PropertyYAlign, cdk.FloatProperty, true, 0.5)
	_ = a.InstallBuildableProperty(PropertyYScale, cdk.FloatProperty, true, 1)
	a.Connect(SignalLostFocus, AlignmentLostFocusHandle, a.childLostFocus)
	a.Connect(SignalGainedFocus, AlignmentGainedFocusHandle, a.childGainedFocus)
	a.Connect(SignalResize, AlignmentEventResizeHandle, a.resize)
	a.Connect(SignalInvalidate, AlignmentInvalidateHandler, a.invalidate)
	a.Connect(SignalDraw, AlignmentDrawHandle, a.draw)
	return false
}

// Get is a convenience method to return the four main Alignment property values
// See: Set()
func (a *CAlignment) Get() (xAlign, yAlign, xScale, yScale float64) {
	xAlign, _ = a.GetFloatProperty(PropertyXAlign)
	yAlign, _ = a.GetFloatProperty(PropertyYAlign)
	xScale, _ = a.GetFloatProperty(PropertyXScale)
	yScale, _ = a.GetFloatProperty(PropertyYScale)
	return
}

// Set is a convenience method to update the four main Alignment property values
//
// Parameters:
// 	xAlign	the horizontal alignment of the child widget, from 0 (left) to 1 (right)
// 	yAlign	the vertical alignment of the child widget, from 0 (top) to 1 (bottom)
// 	xScale	the amount that the child widget expands horizontally to fill up unused space, from 0 to 1. A value of 0 indicates that the child widget should never expand. A value of 1 indicates that the child widget will expand to fill all of the space allocated for the Alignment
// 	yScale	the amount that the child widget expands vertically to fill up unused space, from 0 to 1. The values are similar to xScale
func (a *CAlignment) Set(xAlign, yAlign, xScale, yScale float64) {
	xa, ya, xs, ys := a.Get()
	if xa != xAlign || ya != yAlign || xs != xScale || ys != yScale {
		a.Freeze()
		if err := a.SetFloatProperty(PropertyXAlign, xAlign); err != nil {
			a.LogErr(err)
		}
		if err := a.SetFloatProperty(PropertyYAlign, yAlign); err != nil {
			a.LogErr(err)
		}
		if err := a.SetFloatProperty(PropertyXScale, xScale); err != nil {
			a.LogErr(err)
		}
		if err := a.SetFloatProperty(PropertyYScale, yScale); err != nil {
			a.LogErr(err)
		}
		a.Thaw()
		a.Emit(SignalChanged)
	}
}

// GetPadding is a convenience method to return the four padding property values
// See: SetPadding()
func (a *CAlignment) GetPadding() (paddingTop, paddingBottom, paddingLeft, paddingRight int) {
	paddingTop, _ = a.GetIntProperty(PropertyTopPadding)
	paddingBottom, _ = a.GetIntProperty(PropertyBottomPadding)
	paddingLeft, _ = a.GetIntProperty(PropertyLeftPadding)
	paddingRight, _ = a.GetIntProperty(PropertyRightPadding)
	return
}

// SetPadding is a convenience method to update the padding for the different
// sides of the widget. The padding adds blank space to the sides of the widget.
// For instance, this can be used to indent the child widget towards the right
// by adding padding on the left.
//
// Parameters:
// 	paddingTop	    the padding at the top of the widget
// 	paddingBottom	the padding at the bottom of the widget
// 	paddingLeft	    the padding at the left of the widget
// 	paddingRight	the padding at the right of the widget.
func (a *CAlignment) SetPadding(paddingTop, paddingBottom, paddingLeft, paddingRight int) {
	t, b, l, r := a.GetPadding()
	if t != paddingTop || b != paddingBottom || l != paddingLeft || r != paddingRight {
		a.Freeze()
		if err := a.SetIntProperty(PropertyXAlign, paddingTop); err != nil {
			a.LogErr(err)
		}
		if err := a.SetIntProperty(PropertyYAlign, paddingBottom); err != nil {
			a.LogErr(err)
		}
		if err := a.SetIntProperty(PropertyXScale, paddingLeft); err != nil {
			a.LogErr(err)
		}
		if err := a.SetIntProperty(PropertyYScale, paddingRight); err != nil {
			a.LogErr(err)
		}
		a.Thaw()
		a.Emit(SignalChanged)
	}
}

// Add will set the current child to the Widget instance given, connect two
// signal handlers for losing and gaining focus and finally resize the
// Alignment instance to accommodate the new child Widget.
func (a *CAlignment) Add(w Widget) {
	a.CBin.Add(w)
	w.Connect(SignalLostFocus, AlignmentLostFocusHandle, a.childLostFocus)
	w.Connect(SignalGainedFocus, AlignmentGainedFocusHandle, a.childGainedFocus)
	a.Resize()
}

// Remove will remove the given Widget from the Alignment instance,
// disconnecting any connected focus signal handlers and finally resize the
// Alignment instance to accommodate the lack of content.
func (a *CAlignment) Remove(w Widget) {
	_ = w.Disconnect(SignalLostFocus, AlignmentLostFocusHandle)
	_ = w.Disconnect(SignalGainedFocus, AlignmentGainedFocusHandle)
	a.CBin.Remove(w)
	a.Resize()
}

func (a *CAlignment) childLostFocus(_ []interface{}, _ ...interface{}) enums.EventFlag {
	theme := a.GetTheme()
	a.SetThemeRequest(theme)
	a.Invalidate()
	return enums.EVENT_PASS
}

func (a *CAlignment) childGainedFocus(_ []interface{}, _ ...interface{}) enums.EventFlag {
	theme := a.GetTheme()
	theme.Content.Normal = theme.Content.Focused
	theme.Border.Normal = theme.Border.Focused
	a.SetThemeRequest(theme)
	a.Invalidate()
	return enums.EVENT_PASS
}

func (a *CAlignment) resize(data []interface{}, argv ...interface{}) enums.EventFlag {
	a.Lock()
	defer a.Unlock()
	alloc := a.GetAllocation()
	if alloc.W <= 0 && alloc.H <= 0 {
		if child := a.GetChild(); child != nil {
			child.SetAllocation(ptypes.MakeRectangle(0, 0))
			child.Resize()
		}
		return enums.EVENT_PASS
	}
	if child := a.GetChild(); child != nil {
		xAlign, yAlign, xScale, yScale := a.Get()
		origin := a.GetOrigin()
		size := ptypes.NewRectangle(child.GetSizeRequest())
		if size.W <= -1 || size.W > alloc.W {
			size.W = alloc.W
		}
		if size.H <= -1 || size.H > alloc.H {
			size.H = alloc.H
		}
		if size.W < alloc.W && size.H < alloc.H {
			// available space
			xDelta := alloc.W - size.W
			yDelta := alloc.H - size.H
			// xScale, yScale
			xSize := int(xScale * float64(xDelta))
			ySize := int(yScale * float64(yDelta))
			xDelta -= xSize
			yDelta -= ySize
			size.W += xSize
			size.H += ySize
			// xAlign, yAlign
			xDeltaValue := xAlign * float64(xDelta)
			origin.X += int(xDeltaValue)
			yDeltaValue := yAlign * float64(yDelta)
			origin.Y += int(yDeltaValue)
		}
		child.SetOrigin(origin.X, origin.Y)
		child.SetAllocation(*size)
		child.Resize()
	}
	a.Invalidate()
	return enums.EVENT_PASS
}

func (a *CAlignment) invalidate(data []interface{}, argv ...interface{}) enums.EventFlag {
	theme := a.GetThemeRequest()
	style := theme.Content.Normal
	origin := a.GetOrigin()
	if child := a.GetChild(); child != nil {
		childOrigin := child.GetOrigin()
		childOrigin.SubPoint(origin)
		childSize := child.GetAllocation()
		if err := memphis.ConfigureSurface(child.ObjectID(), childOrigin, childSize, style); err != nil {
			child.LogErr(err)
		}
	}
	return enums.EVENT_PASS
}

func (a *CAlignment) draw(data []interface{}, argv ...interface{}) enums.EventFlag {
	if surface, ok := argv[1].(*memphis.CSurface); ok {
		a.Lock()
		defer a.Unlock()
		alloc := a.GetAllocation()
		if !a.IsVisible() || alloc.W <= 0 || alloc.H <= 0 {
			a.LogTrace("not visible, zero width or zero height")
			return enums.EVENT_PASS
		}

		// render the box and border, with widget
		theme := a.GetThemeRequest()
		boxOrigin := ptypes.MakePoint2I(0, 0)
		boxSize := alloc

		surface.BoxWithTheme(boxOrigin, boxSize, false, true, theme)

		if child := a.GetChild(); child != nil {
			if f := child.Draw(); f == enums.EVENT_STOP {
				if err := surface.Composite(child.ObjectID()); err != nil {
					a.LogError("composite error: %v", err)
				}
			}
		}

		if debug, _ := a.GetBoolProperty(cdk.PropertyDebug); debug {
			surface.DebugBox(paint.ColorSilver, a.ObjectInfo())
		}
		return enums.EVENT_STOP
	}
	return enums.EVENT_PASS
}

// The padding to insert at the bottom of the widget.
// Flags: Read / Write
// Allowed values: <= G_MAXINT
// Default value: 0
const PropertyBottomPadding cdk.Property = "bottom-padding"

// The padding to insert at the left of the widget.
// Flags: Read / Write
// Allowed values: <= G_MAXINT
// Default value: 0
const PropertyLeftPadding cdk.Property = "left-padding"

// The padding to insert at the right of the widget.
// Flags: Read / Write
// Allowed values: <= G_MAXINT
// Default value: 0
const PropertyRightPadding cdk.Property = "right-padding"

// The padding to insert at the top of the widget.
// Flags: Read / Write
// Allowed values: <= G_MAXINT
// Default value: 0
const PropertyTopPadding cdk.Property = "top-padding"

// Horizontal position of child in available space. 0.0 is left aligned, 1.0
// is right aligned.
// Flags: Read / Write
// Allowed values: [0,1]
// Default value: 0.5
const PropertyXAlign cdk.Property = "x-align"

// If available horizontal space is bigger than needed for the child, how
// much of it to use for the child. 0.0 means none, 1.0 means all.
// Flags: Read / Write
// Allowed values: [0,1]
// Default value: 1
const PropertyXScale cdk.Property = "x-scale"

// Vertical position of child in available space. 0.0 is top aligned, 1.0 is
// bottom aligned.
// Flags: Read / Write
// Allowed values: [0,1]
// Default value: 0.5
const PropertyYAlign cdk.Property = "y-align"

// If available vertical space is bigger than needed for the child, how much
// of it to use for the child. 0.0 means none, 1.0 means all.
// Flags: Read / Write
// Allowed values: [0,1]
// Default value: 1
const PropertyYScale cdk.Property = "y-scale"

const AlignmentLostFocusHandle = "alignment-lost-focus-handler"

const AlignmentGainedFocusHandle = "alignment-gained-focus-handler"

const AlignmentDrawHandle = "alignment-draw-handler"

const AlignmentEventResizeHandle = "alignment-event-resize-handler"

const AlignmentInvalidateHandler = "alignment-invalidate-handler"
