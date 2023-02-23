package ctk

import (
	"github.com/go-curses/cdk"
	cenums "github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/cdk/lib/paint"
	"github.com/go-curses/cdk/lib/ptypes"
	"github.com/go-curses/cdk/memphis"

	"github.com/go-curses/ctk/lib/enums"
)

const TypeFrame cdk.CTypeTag = "ctk-frame"

func init() {
	_ = cdk.TypesManager.AddType(TypeFrame, func() interface{} { return MakeFrame() })
}

// Frame Hierarchy:
//
//	Object
//	  +- Widget
//	    +- Container
//	      +- Bin
//	        +- Frame
//	          +- AspectFrame
//
// The Frame Widget wraps other Widgets with a border and optional title label.
type Frame interface {
	Bin
	Buildable

	GetLabel() (value string)
	SetLabel(label string)
	GetLabelWidget() (value Widget)
	SetLabelWidget(labelWidget Widget)
	GetLabelAlign() (xAlign float64, yAlign float64)
	SetLabelAlign(xAlign float64, yAlign float64)
	GetShadowType() (value enums.ShadowType)
	SetShadowType(shadowType enums.ShadowType)
	GetFocusWithChild() (focusWithChild bool)
	SetFocusWithChild(focusWithChild bool)
}

var _ Frame = (*CFrame)(nil)

// The CFrame structure implements the Frame interface and is exported
// to facilitate type embedding with custom implementations. No member variables
// are exported as the interface methods are the only intended means of
// interacting with Frame objects.
type CFrame struct {
	CBin

	focusWithChild bool
}

// MakeFrame is used by the Buildable system to construct a new Frame.
func MakeFrame() Frame {
	return NewFrame("")
}

// NewFrame is the constructor for new Frame instances.
func NewFrame(text string) Frame {
	f := new(CFrame)
	f.Init()
	label := NewLabel(text)
	label.SetSingleLineMode(true)
	label.SetLineWrap(false)
	label.SetLineWrapMode(cenums.WRAP_NONE)
	label.SetJustify(cenums.JUSTIFY_LEFT)
	label.Show()
	theme := label.GetTheme()
	theme.Content.FillRune = paint.DefaultNilRune
	theme.Border.FillRune = paint.DefaultNilRune
	label.SetTheme(theme)
	f.SetLabelWidget(label)
	return f
}

// NewFrameWithWidget will construct a new Frame with the given widget instead
// of the default Label.
func NewFrameWithWidget(w Widget) Frame {
	f := new(CFrame)
	f.Init()
	f.SetLabelWidget(w)
	return f
}

// Init initializes a Frame object. This must be called at least once to
// set up the necessary defaults and allocate any memory structures. Calling
// this more than once is safe though unnecessary. Only the first call will
// result in any effect upon the Frame instance. Init is used in the
// NewFrame constructor and only necessary when implementing a derivative
// Frame type.
func (f *CFrame) Init() (already bool) {
	if f.InitTypeItem(TypeFrame, f) {
		return true
	}
	f.CBin.Init()
	f.flags = enums.NULL_WIDGET_FLAG
	f.SetFlags(enums.PARENT_SENSITIVE | enums.APP_PAINTABLE)
	f.focusWithChild = false

	_ = f.InstallProperty(PropertyLabel, cdk.StringProperty, true, nil)
	_ = f.InstallProperty(PropertyLabelWidget, cdk.StructProperty, true, nil)
	_ = f.InstallProperty(PropertyLabelXAlign, cdk.FloatProperty, true, 0.0)
	_ = f.InstallProperty(PropertyLabelYAlign, cdk.FloatProperty, true, 0.5)
	_ = f.InstallProperty(PropertyShadow, cdk.StructProperty, true, nil)
	_ = f.InstallProperty(PropertyShadowType, cdk.StructProperty, true, nil)

	f.Connect(SignalResize, FrameResizeHandle, f.resize)
	f.Connect(SignalDraw, FrameDrawHandle, f.draw)
	return false
}

func (f *CFrame) SetWindow(w Window) {
	if widget := f.GetLabelWidget(); widget != nil {
		WidgetRecurseSetWindow(widget, w)
	}
	f.CBin.SetWindow(w)
}

// GetLabel returns the text in the label Widget, if the Widget is in
// fact of Label Widget type. If the label Widget is not an actual Label, the
// value of the Frame label property is returned.
//
// Returns:
//
//	the text in the label, or NULL if there was no label widget or
//	the label widget was not a Label. This string is owned by
//
// Locking: read
func (f *CFrame) GetLabel() (value string) {
	var err error
	if w := f.GetLabelWidget(); w != nil {
		if lw, ok := w.Self().(Label); ok {
			f.RLock()
			value = lw.GetLabel()
			f.RUnlock()
			return
		}
	}
	f.RLock()
	if value, err = f.GetStringProperty(PropertyLabel); err != nil {
		f.LogErr(err)
	}
	f.RUnlock()
	return
}

// SetLabel updates the text of the Label.
//
// Parameters:
//
//	label	the text to use as the label of the frame.
//
// Locking: write
func (f *CFrame) SetLabel(label string) {
	if err := f.SetStringProperty(PropertyLabel, label); err != nil {
		f.LogErr(err)
	} else {
		if w := f.GetLabelWidget(); w != nil {
			if lw, ok := w.Self().(Label); ok {
				lw.SetText(label)
				f.Invalidate()
			}
		}
	}
}

// GetLabelWidget retrieves the label widget for the Frame.
// See: SetLabelWidget()
//
// Locking: read
func (f *CFrame) GetLabelWidget() (value Widget) {
	// f.RLock()
	// defer f.RUnlock()
	if v, err := f.GetStructProperty(PropertyLabelWidget); err != nil {
		f.LogErr(err)
	} else if v != nil {
		var ok bool
		if value, ok = v.(Widget); !ok {
			f.LogError("value stored in %v is not a Widget: %v (%T)", PropertyLabelWidget, v, v)
		}
	}
	return
}

// SetLabelWidget removes any existing Widget and replaces it with the given
// one. This is the widget that will appear embedded in the top edge of the
// frame as a title.
//
// Parameters:
//
//	labelWidget	the new label widget
//
// Locking: write
func (f *CFrame) SetLabelWidget(widget Widget) {
	var previousWidget Widget
	if found, err := f.GetStructProperty(PropertyLabelWidget); err == nil {
		if fw, ok := found.(Widget); ok {
			previousWidget = fw
		}
	}
	if err := f.SetStructProperty(PropertyLabelWidget, widget); err != nil {
		f.LogErr(err)
	} else {
		if previousWidget != nil {
			f.PopCompositeChild(previousWidget)
		}
		f.PushCompositeChild(widget)
		f.Invalidate()
	}
}

// GetLabelAlign retrieves the X and Y alignment of the frame's label. If the
// label Widget is not of Label Widget type, then the values of the
// label-x-align and label-y-align properties are returned.
// See: SetLabelAlign()
//
// Parameters:
//
//	xAlign	X alignment of frame's label
//	yAlign	Y alignment of frame's label
//
// Locking: read
func (f *CFrame) GetLabelAlign() (xAlign float64, yAlign float64) {
	var err error
	if w := f.GetLabelWidget(); w != nil {
		if lw, ok := w.Self().(Label); ok {
			xAlign, yAlign = lw.GetAlignment()
			return
		}
	}
	if xAlign, err = f.GetFloatProperty(PropertyLabelXAlign); err != nil {
		f.LogErr(err)
	}
	if yAlign, err = f.GetFloatProperty(PropertyLabelYAlign); err != nil {
		f.LogErr(err)
	}
	return
}

// SetLabelAlign is a convenience method for setting the label-x-align and
// label-y-align properties of the Frame. The default values for a newly created
// frame are 0.0 and 0.5. If the label Widget is in fact of Label Widget type,
// SetAlignment with the given x and y alignment values.
//
// Parameters:
//
//	xAlign	The position of the label along the top edge of the widget. A value
//	        of 0.0 represents left alignment; 1.0 represents right alignment.
//	yAlign	The y alignment of the label. A value of 0.0 aligns under the frame;
//	        1.0 aligns above the frame. If the values are exactly 0.0 or 1.0 the
//	        gap in the frame won't be painted because the label will be
//	        completely above or below the frame.
//
// Locking: write
func (f *CFrame) SetLabelAlign(xAlign float64, yAlign float64) {
	f.Lock()
	if err := f.SetFloatProperty(PropertyLabelXAlign, xAlign); err != nil {
		f.LogErr(err)
	}
	if err := f.SetFloatProperty(PropertyLabelYAlign, yAlign); err != nil {
		f.LogErr(err)
	}
	f.Unlock()
	if w := f.GetLabelWidget(); w != nil {
		if lw, ok := w.Self().(Label); ok {
			f.Lock()
			lw.SetAlignment(xAlign, yAlign)
			f.Unlock()
		}
	}
}

// GetShadowType returns the shadow type of the frame.
// See: SetShadowType()
//
// Locking: read
func (f *CFrame) GetShadowType() (value enums.ShadowType) {
	f.RLock()
	defer f.RUnlock()
	if v, err := f.GetStructProperty(PropertyShadowType); err != nil {
		f.LogErr(err)
	} else {
		var ok bool
		if value, ok = v.(enums.ShadowType); !ok {
			f.LogError("value stored in %v is not a ShadowType: %v (%T)", PropertyShadowType, v, v)
		}
	}
	return
}

// SetShadowType updates the shadow-type property for the Frame.
//
// Parameters:
//
//	type	the new ShadowType
//
// Note that usage of this within CTK is unimplemented at this time
func (f *CFrame) SetShadowType(shadowType enums.ShadowType) {
	f.Lock()
	defer f.Unlock()
	if err := f.SetStructProperty(PropertyShadowType, shadowType); err != nil {
		f.LogErr(err)
	}
}

// Add will add the given Widget to the Frame. As the Frame Widget is of Bin
// type, any previous child Widget is removed first.
//
// Locking: write
func (f *CFrame) Add(w Widget) {
	f.CBin.Add(w)
	w.Connect(SignalLostFocus, FrameChildLostFocusHandle, f.childLostFocus)
	w.Connect(SignalGainedFocus, FrameChildGainedFocusHandle, f.childGainedFocus)
	f.Invalidate()
}

// Remove will remove the given Widget from the Frame.
//
// Locking: write
func (f *CFrame) Remove(w Widget) {
	_ = w.Disconnect(SignalLostFocus, FrameChildLostFocusHandle)
	_ = w.Disconnect(SignalGainedFocus, FrameChildGainedFocusHandle)
	f.CBin.Remove(w)
	f.Invalidate()
}

// IsFocus is a convenience method for returning whether the child Widget is the
// focused Widget. If no child Widget exists, or the child Widget cannot be
// focused itself, then the return value is whether the Frame itself is the
// focused Widget.
//
// Locking: read
func (f *CFrame) IsFocus() bool {
	if f.GetFocusWithChild() {
		if child := f.GetChild(); child != nil && child.CanFocus() {
			return child.IsFocus()
		}
	}
	return f.CBin.IsFocus()
}

// GetFocusWithChild returns true if the Frame is supposed to follow the focus
// of its child Widget or if it should follow its own focus.
// See: SetFocusWithChild()
//
// Locking: read
func (f *CFrame) GetFocusWithChild() (focusWithChild bool) {
	f.RLock()
	focusWithChild = f.focusWithChild
	f.RUnlock()
	return
}

// SetFocusWithChild updates whether the Frame's theme will reflect the focused
// state of the Frame's child Widget.
//
// Locking: write
func (f *CFrame) SetFocusWithChild(focusWithChild bool) {
	f.Lock()
	f.focusWithChild = focusWithChild
	f.Unlock()
}

// GetSizeRequest returns the requested size of the Frame, taking into account
// any children and their size requests.
func (f *CFrame) GetSizeRequest() (width, height int) {
	_, yAlign := f.GetLabelAlign()
	size := ptypes.NewRectangle(f.CWidget.GetSizeRequest())
	if child := f.GetChild(); child != nil {
		childSize := ptypes.NewRectangle(child.GetSizeRequest())
		if size.W <= -1 {
			if childSize.W > -1 {
				size.W = 1 + childSize.W + 1
			}
		}
		if size.H <= -1 {
			if childSize.H > -1 {
				size.H = 1 + childSize.H + 1
				if yAlign == 0.0 {
					size.H += 1
				}
			}
		}
	}
	return size.W, size.H
}

func (f *CFrame) GetWidgetAt(p *ptypes.Point2I) Widget {
	if f.HasPoint(p) && f.IsVisible() {
		if widget := f.GetLabelWidget(); widget != nil && widget.HasPoint(p) {
			return widget.GetWidgetAt(p)
		}
		if child := f.GetChild(); child != nil {
			return child.GetWidgetAt(p)
		}
		return f
	}
	return nil
}

func (f *CFrame) childLostFocus(_ []interface{}, _ ...interface{}) cenums.EventFlag {
	f.UnsetState(enums.StateSelected)
	if child := f.GetChild(); child != nil {
		child.UnsetState(enums.StateSelected)
		child.Invalidate()
	}
	f.Invalidate()
	return cenums.EVENT_PASS
}

func (f *CFrame) childGainedFocus(_ []interface{}, _ ...interface{}) cenums.EventFlag {
	f.SetState(enums.StateSelected)
	if child := f.GetChild(); child != nil {
		child.SetState(enums.StateSelected)
		child.Invalidate()
	}
	f.Invalidate()
	return cenums.EVENT_PASS
}

func (f *CFrame) resize(data []interface{}, argv ...interface{}) cenums.EventFlag {

	// our allocation has been set prior to Resize() being called
	alloc := f.GetAllocation()
	widget := f.GetLabelWidget()
	origin := f.GetOrigin()
	xAlign, yAlign := f.GetLabelAlign()
	child := f.GetChild()

	if widget == nil {
		return cenums.EVENT_PASS
	}

	label, _ := widget.(Label)
	if alloc.W <= 0 && alloc.H <= 0 {
		if label != nil {
			label.SetAllocation(ptypes.MakeRectangle(0, 0))
			label.Resize()
		} else if widget != nil {
			widget.SetAllocation(ptypes.MakeRectangle(0, 0))
			widget.Resize()
		}
		return cenums.EVENT_PASS
	}

	alloc.Floor(0, 0)
	childOrigin := ptypes.MakePoint2I(origin.X+1, origin.Y+1)
	childAlloc := ptypes.MakeRectangle(alloc.W-2, alloc.H-2)
	labelOrigin := ptypes.MakePoint2I(origin.X+2, origin.Y)
	labelAlloc := ptypes.MakeRectangle(alloc.W-4, 1)

	if yAlign <= 0.0 {
		yAlign = 0.0
		childAlloc.H -= 1
		labelOrigin.X -= 1
		childOrigin.Y += 1
	} else if yAlign >= 1.0 {
		yAlign = 1.0
		labelOrigin.Y += 1
		childAlloc.H -= 1
		childOrigin.Y += 1
	} else {
		yAlign = 0.5
	}

	if label != nil {
		label.SetAlignment(xAlign, yAlign)
		label.SetMaxWidthChars(labelAlloc.W)
		label.SetOrigin(labelOrigin.X, labelOrigin.Y)
		label.SetAllocation(labelAlloc)
		label.Resize()
	} else if widget != nil {
		widget.SetOrigin(labelOrigin.X, labelOrigin.Y)
		widget.SetAllocation(labelAlloc)
		widget.Resize()
	}

	if child != nil {
		child.SetOrigin(childOrigin.X, childOrigin.Y)
		child.SetAllocation(childAlloc)
		child.Resize()
	}

	f.Invalidate()
	return cenums.EVENT_STOP
}

func (f *CFrame) draw(data []interface{}, argv ...interface{}) cenums.EventFlag {

	if surface, ok := argv[1].(*memphis.CSurface); ok {
		alloc := f.GetAllocation()
		if !f.IsVisible() || alloc.W <= 0 || alloc.H <= 0 {
			f.LogTrace("not visible, zero width or zero height")
			return cenums.EVENT_PASS
		}

		// render the box and border, with widget
		_, yAlign := f.GetLabelAlign()
		widget := f.GetLabelWidget()
		child := f.GetChild()
		theme := f.GetThemeRequest()

		boxOrigin := ptypes.MakePoint2I(0, 0)
		boxSize := alloc
		if yAlign == 0.0 { // top
			boxOrigin.Y += 1
			boxSize.H -= 1
		}

		surface.BoxWithTheme(boxOrigin, boxSize, true, true, theme)

		if widget != nil {
			if label, ok := widget.Self().(Label); ok {
				label.Draw()
				label.LockDraw()
				if err := surface.Composite(label.ObjectID()); err != nil {
					f.LogError("composite error: %v", err)
				}
				label.UnlockDraw()
			} else {
				widget.Draw()
				widget.LockDraw()
				if err := surface.Composite(widget.ObjectID()); err != nil {
					f.LogError("composite error: %v", err)
				}
				widget.UnlockDraw()
			}
		}

		if child != nil {
			child.Draw()
			child.LockDraw()
			if err := surface.Composite(child.ObjectID()); err != nil {
				f.LogError("composite error: %v", err)
			}
			child.UnlockDraw()
		}

		if debug, _ := f.GetBoolProperty(cdk.PropertyDebug); debug {
			surface.DebugBox(paint.ColorSilver, f.ObjectInfo())
		}

		return cenums.EVENT_STOP
	}
	return cenums.EVENT_PASS
}

// Text of the frame's label.
// Flags: Read / Write
// Default value: NULL
// const PropertyLabel cdk.Property = "label"

// A widget to display in place of the usual frame label.
// Flags: Read / Write
const PropertyLabelWidget cdk.Property = "label-widget"

// The horizontal alignment of the label.
// Flags: Read / Write
// Allowed values: [0,1]
// Default value: 0
const PropertyLabelXAlign cdk.Property = "label-x-align"

// The vertical alignment of the label.
// Flags: Read / Write
// Allowed values: [0,1]
// Default value: 0.5
const PropertyLabelYAlign cdk.Property = "label-y-align"

// Deprecated property, use shadow_type instead.
// Flags: Read / Write
// Default value: GTK_SHADOW_ETCHED_IN
const PropertyShadow cdk.Property = "shadow"

// Appearance of the frame border.
// Flags: Read / Write
// Default value: GTK_SHADOW_ETCHED_IN
const PropertyShadowType cdk.Property = "shadow-type"

const FrameChildLostFocusHandle = "frame-child-lost-focus-handler"

const FrameChildGainedFocusHandle = "frame-child-gained-focus-handler"

const FrameInvalidateHandle = "frame-invalidate-handler"

const FrameResizeHandle = "frame-resize-handler"

const FrameDrawHandle = "frame-draw-handler"