package ctk

import (
	"github.com/go-curses/cdk"
	"github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/cdk/lib/paint"
	"github.com/go-curses/cdk/lib/ptypes"
	"github.com/go-curses/cdk/memphis"
)

const TypeFrame cdk.CTypeTag = "ctk-frame"

func init() {
	_ = cdk.TypesManager.AddType(TypeFrame, func() interface{} { return MakeFrame() })
}

// Frame Hierarchy:
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

	Init() (already bool)
	GetLabel() (value string)
	SetLabel(label string)
	GetLabelWidget() (value Widget)
	SetLabelWidget(labelWidget Widget)
	GetLabelAlign() (xAlign float64, yAlign float64)
	SetLabelAlign(xAlign float64, yAlign float64)
	GetShadowType() (value ShadowType)
	SetShadowType(shadowType ShadowType)
	Add(w Widget)
	Remove(w Widget)
	IsFocus() bool
	GetFocusWithChild() (focusWithChild bool)
	SetFocusWithChild(focusWithChild bool)
	GetThemeRequest() (theme paint.Theme)
	GetSizeRequest() (width, height int)
	GetWidgetAt(p *ptypes.Point2I) Widget
}

// The CFrame structure implements the Frame interface and is exported
// to facilitate type embedding with custom implementations. No member variables
// are exported as the interface methods are the only intended means of
// interacting with Frame objects.
type CFrame struct {
	CBin

	focusWithChild bool
}

// MakeFrame is used by the Buildable system to construct a new Frame.
func MakeFrame() *CFrame {
	return NewFrame("")
}

// NewFrame is the constructor for new Frame instances.
func NewFrame(text string) *CFrame {
	f := new(CFrame)
	f.Init()
	label := NewLabel(text)
	label.SetSingleLineMode(true)
	label.SetLineWrap(false)
	label.SetLineWrapMode(enums.WRAP_NONE)
	label.SetJustify(enums.JUSTIFY_LEFT)
	label.SetParent(f)
	label.SetWindow(f.GetWindow())
	label.Show()
	f.SetLabelWidget(label)
	return f
}

// NewFrameWithWidget will construct a new Frame with the given widget instead
// of the default Label.
func NewFrameWithWidget(w Widget) *CFrame {
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
	f.flags = NULL_WIDGET_FLAG
	f.SetFlags(PARENT_SENSITIVE | APP_PAINTABLE)
	_ = f.InstallProperty(PropertyLabel, cdk.StringProperty, true, nil)
	_ = f.InstallProperty(PropertyLabelWidget, cdk.StructProperty, true, nil)
	_ = f.InstallProperty(PropertyLabelXAlign, cdk.FloatProperty, true, 0.0)
	_ = f.InstallProperty(PropertyLabelYAlign, cdk.FloatProperty, true, 0.5)
	_ = f.InstallProperty(PropertyShadow, cdk.StructProperty, true, nil)
	_ = f.InstallProperty(PropertyShadowType, cdk.StructProperty, true, nil)
	f.focusWithChild = false
	f.Connect(SignalInvalidate, FrameInvalidateHandle, f.invalidate)
	f.Connect(SignalResize, FrameResizeHandle, f.resize)
	f.Connect(SignalDraw, FrameDrawHandle, f.draw)
	return false
}

// GetLabel returns the text in the label Widget, if the Widget is in
// fact of Label Widget type. If the label Widget is not an actual Label, the
// value of the Frame label property is returned.
//
// Returns:
// 	the text in the label, or NULL if there was no label widget or
// 	the label widget was not a Label. This string is owned by
//
// Locking: read
func (f *CFrame) GetLabel() (value string) {
	var err error
	if w := f.GetLabelWidget(); w != nil {
		if lw, ok := w.(Label); ok {
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
// 	label	the text to use as the label of the frame.
//
// Locking: write
func (f *CFrame) SetLabel(label string) {
	if err := f.SetStringProperty(PropertyLabel, label); err != nil {
		f.LogErr(err)
	} else {
		if w := f.GetLabelWidget(); w != nil {
			if lw, ok := w.(Label); ok {
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
// 	labelWidget	the new label widget
//
// Locking: write
func (f *CFrame) SetLabelWidget(labelWidget Widget) {
	window := f.GetWindow()
	if err := f.SetStructProperty(PropertyLabelWidget, labelWidget); err != nil {
		f.LogErr(err)
	} else {
		labelWidget.SetParent(f)
		labelWidget.SetWindow(window)
		labelWidget.Show()
		f.Invalidate()
	}
}

// GetLabelAlign retrieves the X and Y alignment of the frame's label. If the
// label Widget is not of Label Widget type, then the values of the
// label-x-align and label-y-align properties are returned.
// See: SetLabelAlign()
//
// Parameters:
// 	xAlign	X alignment of frame's label
// 	yAlign	Y alignment of frame's label
//
// Locking: read
func (f *CFrame) GetLabelAlign() (xAlign float64, yAlign float64) {
	var err error
	if w := f.GetLabelWidget(); w != nil {
		if lw, ok := w.(Label); ok {
			xAlign, yAlign = lw.GetAlignment()
			return
		}
	}
	f.RLock()
	if xAlign, err = f.GetFloatProperty(PropertyLabelXAlign); err != nil {
		f.LogErr(err)
	}
	if yAlign, err = f.GetFloatProperty(PropertyLabelYAlign); err != nil {
		f.LogErr(err)
	}
	f.RUnlock()
	return
}

// SetLabelAlign is a convenience method for setting the label-x-align and
// label-y-align properties of the Frame. The default values for a newly created
// frame are 0.0 and 0.5. If the label Widget is in fact of Label Widget type,
// SetAlignment with the given x and y alignment values.
//
// Parameters:
// 	xAlign	The position of the label along the top edge of the widget. A value
// 	        of 0.0 represents left alignment; 1.0 represents right alignment.
// 	yAlign	The y alignment of the label. A value of 0.0 aligns under the frame;
// 	        1.0 aligns above the frame. If the values are exactly 0.0 or 1.0 the
// 	        gap in the frame won't be painted because the label will be
// 	        completely above or below the frame.
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
		if lw, ok := w.(Label); ok {
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
func (f *CFrame) GetShadowType() (value ShadowType) {
	f.RLock()
	defer f.RUnlock()
	if v, err := f.GetStructProperty(PropertyShadowType); err != nil {
		f.LogErr(err)
	} else {
		var ok bool
		if value, ok = v.(ShadowType); !ok {
			f.LogError("value stored in %v is not a ShadowType: %v (%T)", PropertyShadowType, v, v)
		}
	}
	return
}

// SetShadowType updates the shadow-type property for the Frame.
//
// Parameters:
// 	type	the new ShadowType
//
// Note that usage of this within CTK is unimplemented at this time
func (f *CFrame) SetShadowType(shadowType ShadowType) {
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
	f.Invalidate()
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
		if child := f.GetChild(); child != nil {
			return child
		}
		return f
	}
	return nil
}

func (f *CFrame) childLostFocus(_ []interface{}, _ ...interface{}) enums.EventFlag {
	f.UnsetState(StateSelected)
	if child := f.GetChild(); child != nil {
		child.UnsetState(StateSelected)
		child.Invalidate()
	}
	f.Invalidate()
	return enums.EVENT_PASS
}

func (f *CFrame) childGainedFocus(_ []interface{}, _ ...interface{}) enums.EventFlag {
	f.SetState(StateSelected)
	if child := f.GetChild(); child != nil {
		child.SetState(StateSelected)
		child.Invalidate()
	}
	f.Invalidate()
	return enums.EVENT_PASS
}

func (f *CFrame) invalidate(data []interface{}, argv ...interface{}) enums.EventFlag {
	wantStop := false
	origin := f.GetOrigin()
	theme := f.GetThemeRequest()
	labelWidget := f.GetLabelWidget()
	child := f.GetChild()
	f.Lock()
	defer f.Unlock()
	if labelChild, ok := labelWidget.(Label); ok && labelChild != nil {
		local := labelChild.GetOrigin()
		local.SubPoint(origin)
		alloc := labelChild.GetAllocation()
		if surface, err := memphis.GetSurface(labelChild.ObjectID()); err != nil {
			labelChild.LogErr(err)
		} else {
			surface.SetOrigin(local)
			surface.Resize(alloc, theme.Content.Normal)
			theme.Content.FillRune = rune(0)
			surface.Fill(theme)
		}
		labelChild.SetTheme(theme)
		labelChild.Invalidate()
		wantStop = true
	}
	if child != nil {
		local := child.GetOrigin()
		local.SubPoint(origin)
		alloc := child.GetAllocation()
		if err := memphis.ConfigureSurface(child.ObjectID(), local, alloc, theme.Content.Normal); err != nil {
			child.LogErr(err)
		}
		wantStop = true
		child.Invalidate()
	}
	if wantStop {
		return enums.EVENT_STOP
	}
	return enums.EVENT_PASS
}

func (f *CFrame) resize(data []interface{}, argv ...interface{}) enums.EventFlag {
	// our allocation has been set prior to Resize() being called
	alloc := f.GetAllocation()
	widget := f.GetLabelWidget()
	origin := f.GetOrigin()
	xAlign, yAlign := f.GetLabelAlign()
	child := f.GetChild()

	f.Lock()

	if widget == nil {
		f.Unlock()
		return enums.EVENT_PASS
	}

	label, _ := widget.(Label)
	if alloc.W <= 0 && alloc.H <= 0 {
		if label != nil {
			label.SetAllocation(ptypes.MakeRectangle(0, 0))
			label.Resize()
		}
		f.Unlock()
		return enums.EVENT_PASS
	}

	alloc.Floor(0, 0)
	childOrigin := ptypes.MakePoint2I(origin.X+1, origin.Y+1)
	childAlloc := ptypes.MakeRectangle(alloc.W-2, alloc.H-2)
	labelOrigin := ptypes.MakePoint2I(origin.X+2, origin.Y)
	labelAlloc := ptypes.MakeRectangle(alloc.W-4, 1)
	if yAlign <= 0.0 {
		yAlign = 0.0
		childAlloc.H -= 1
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
	}
	if child != nil {
		child.SetOrigin(childOrigin.X, childOrigin.Y)
		child.SetAllocation(childAlloc)
		child.Resize()
	}

	f.Unlock()

	f.Invalidate()
	return enums.EVENT_STOP
}

func (f *CFrame) draw(data []interface{}, argv ...interface{}) enums.EventFlag {
	if surface, ok := argv[1].(*memphis.CSurface); ok {
		alloc := f.GetAllocation()
		if !f.IsVisible() || alloc.W <= 0 || alloc.H <= 0 {
			f.LogTrace("not visible, zero width or zero height")
			return enums.EVENT_PASS
		}

		// render the box and border, with widget
		_, yAlign := f.GetLabelAlign()
		widget := f.GetLabelWidget()
		child := f.GetChild()
		theme := f.GetThemeRequest()
		if child != nil {
			theme = child.GetThemeRequest()
		}

		f.Lock()
		defer f.Unlock()

		boxOrigin := ptypes.MakePoint2I(0, 0)
		boxSize := alloc
		if yAlign == 0.0 { // top
			boxOrigin.Y += 1
			boxSize.H -= 1
		}

		surface.BoxWithTheme(boxOrigin, boxSize, true, true, theme)

		if widget != nil {
			if label, ok := widget.(Label); ok {
				if rv := label.Draw(); rv == enums.EVENT_STOP {
					if err := surface.Composite(label.ObjectID()); err != nil {
						f.LogError("composite error: %v", err)
					}
				}
			}
		}

		if child != nil {
			if rv := child.Draw(); rv == enums.EVENT_STOP {
				if err := surface.Composite(child.ObjectID()); err != nil {
					f.LogError("composite error: %v", err)
				}
			}
		}

		if debug, _ := f.GetBoolProperty(cdk.PropertyDebug); debug {
			surface.DebugBox(paint.ColorSilver, f.ObjectInfo())
		}
		return enums.EVENT_STOP
	}
	return enums.EVENT_PASS
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
