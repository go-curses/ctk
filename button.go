package ctk

// TODO: mnemonics support, Accel

import (
	"strings"

	"github.com/go-curses/cdk"
	"github.com/go-curses/cdk/lib/enums"
	cmath "github.com/go-curses/cdk/lib/math"
	"github.com/go-curses/cdk/lib/paint"
	"github.com/go-curses/cdk/lib/ptypes"
	cstrings "github.com/go-curses/cdk/lib/strings"
	"github.com/go-curses/cdk/memphis"
)

const TypeButton cdk.CTypeTag = "ctk-button"

var (
	// DefaultButtonTheme enables customized theming of default stock buttons.
	DefaultButtonTheme = paint.Theme{
		Content: paint.ThemeAspect{
			Normal:      paint.DefaultColorStyle.Foreground(paint.ColorWhite).Background(paint.ColorFireBrick).Dim(true).Bold(false),
			Focused:     paint.DefaultColorStyle.Foreground(paint.ColorWhite).Background(paint.ColorDarkRed).Dim(false).Bold(true),
			Active:      paint.DefaultColorStyle.Foreground(paint.ColorWhite).Background(paint.ColorDarkRed).Dim(false).Bold(true).Reverse(true),
			FillRune:    paint.DefaultFillRune,
			BorderRunes: paint.DefaultBorderRune,
			ArrowRunes:  paint.DefaultArrowRune,
			Overlay:     false,
		},
		Border: paint.ThemeAspect{
			Normal:      paint.DefaultColorStyle.Foreground(paint.ColorWhite).Background(paint.ColorFireBrick).Dim(true).Bold(false),
			Focused:     paint.DefaultColorStyle.Foreground(paint.ColorWhite).Background(paint.ColorDarkRed).Dim(false).Bold(true),
			Active:      paint.DefaultColorStyle.Foreground(paint.ColorWhite).Background(paint.ColorDarkRed).Dim(false).Bold(true).Reverse(true),
			FillRune:    paint.DefaultFillRune,
			BorderRunes: paint.DefaultBorderRune,
			ArrowRunes:  paint.DefaultArrowRune,
			Overlay:     false,
		},
	}
)

func init() {
	_ = cdk.TypesManager.AddType(TypeButton, func() interface{} { return MakeButton() })
}

// Button Hierarchy:
//	Object
//	  +- Widget
//	    +- Container
//	      +- Bin
//	        +- Button
//	          +- ToggleButton
//	          +- ColorButton
//	          +- FontButton
//	          +- LinkButton
//	          +- OptionMenu
//	          +- ScaleButton
//
// The Button Widget is a Bin Container that represents a focusable Drawable
// Widget that is Sensitive to event interactions.
type Button interface {
	Bin
	Activatable
	Alignable
	Buildable

	Init() (already bool)
	Build(builder Builder, element *CBuilderElement) error
	Activate() (value bool)
	Clicked() enums.EventFlag
	GetRelief() (value ReliefStyle)
	SetRelief(newStyle ReliefStyle)
	GetLabel() (value string)
	SetLabel(label string)
	GetUseStock() (value bool)
	SetUseStock(useStock bool)
	GetUseUnderline() (enabled bool)
	SetUseUnderline(enabled bool)
	GetUseMarkup() (enabled bool)
	SetUseMarkup(enabled bool)
	GetFocusOnClick() (value bool)
	SetFocusOnClick(focusOnClick bool)
	GetAlignment() (xAlign float64, yAlign float64)
	SetAlignment(xAlign float64, yAlign float64)
	GetImage() (value Widget, ok bool)
	SetImage(image Widget)
	GetImagePosition() (value PositionType)
	SetImagePosition(position PositionType)
	GetPressed() bool
	SetPressed(pressed bool)
	GetFocusChain() (focusableWidgets []interface{}, explicitlySet bool)
	GrabFocus()
	GrabEventFocus()
	CancelEvent()
	GetWidgetAt(p *ptypes.Point2I) Widget
	GetThemeRequest() (theme paint.Theme)
	GetSizeRequest() (width, height int)
}

// The CButton structure implements the Button interface and is exported to
// facilitate type embedding with custom implementations. No member variables
// are exported as the interface methods are the only intended means of
// interacting with Button objects.
type CButton struct {
	CBin

	pressed bool
}

// MakeButton is used by the Buildable system to construct a new Button with
// a default label that is empty.
func MakeButton() *CButton {
	b := NewButtonWithLabel("")
	return b
}

// NewButton is a constructor for new Box instances without any pre-set label.
func NewButton() *CButton {
	b := new(CButton)
	b.Init()
	return b
}

// NewButtonWithLabel will construct a NewButton with a NewLabel, pre-configured
// for a centered placement within the new Button instance, taking care to avoid
// focus complications and default event handling.
//
// Parameters:
// 	label	the text of the button
func NewButtonWithLabel(text string) (b *CButton) {
	b = NewButton()
	label := NewLabel(text)
	b.Add(label)
	label.SetTheme(DefaultButtonTheme)
	label.UnsetFlags(CAN_FOCUS)
	label.UnsetFlags(CAN_DEFAULT)
	label.UnsetFlags(RECEIVES_DEFAULT)
	label.SetLineWrap(false)
	label.SetLineWrapMode(enums.WRAP_NONE)
	label.SetJustify(enums.JUSTIFY_CENTER)
	label.SetAlignment(0.5, 0.5)
	label.SetSingleLineMode(true)
	label.Show()
	return b
}

// NewButtonWithMnemonic creates a NewButtonWithLabel. If the characters in the
// label are preceded by an underscore, they are underlined. If you need a
// literal underscore character in a label, use '__' (two underscores). The
// first underlined character represents a keyboard accelerator called a
// mnemonic. Pressing Alt and that key activates the button.
//
// Parameters:
// 	label	the text of the button
func NewButtonWithMnemonic(text string) (b *CButton) {
	b = NewButtonWithLabel(text)
	b.SetUseUnderline(true)
	return b
}

// NewButtonFromStock creates a NewButtonWithLabel containing the text from a
// stock item. If stock_id is unknown, it will be treated as a mnemonic label
// (as for NewWithMnemonic).
//
// Parameters:
// 	stockId	the name of the stock item
func NewButtonFromStock(stockId StockID) (value *CButton) {
	b := NewButtonWithLabel("")
	b.Init()
	if item := LookupStockItem(stockId); item != nil {
		b.SetUseStock(true)
		b.SetUseUnderline(true)
		b.SetLabel(item.Label)
	} else {
		b.SetLabel(string(stockId))
	}
	return b
}

// NewbuttonWithWidget creates a NewButton with the given Widget as the Button's
// child.
func NewButtonWithWidget(w Widget) *CButton {
	b := NewButton()
	b.Add(w)
	return b
}

// Init initializes a Button object. This must be called at least once to
// set up the necessary defaults and allocate any memory structures. Calling
// this more than once is safe though unnecessary. Only the first call will
// result in any effect upon the Button instance. Init is used in the
// NewButton constructor and only necessary when implementing a derivative
// Button type.
func (b *CButton) Init() (already bool) {
	if b.InitTypeItem(TypeButton, b) {
		return true
	}
	b.CBin.Init()
	b.flags = NULL_WIDGET_FLAG
	b.SetFlags(SENSITIVE | PARENT_SENSITIVE | CAN_DEFAULT | RECEIVES_DEFAULT | CAN_FOCUS | APP_PAINTABLE)
	b.SetTheme(DefaultButtonTheme)
	b.pressed = false
	_ = b.InstallBuildableProperty(PropertyFocusOnClick, cdk.BoolProperty, true, true)
	_ = b.InstallBuildableProperty(PropertyButtonLabel, cdk.StringProperty, true, nil)
	_ = b.InstallBuildableProperty(PropertyRelief, cdk.StructProperty, true, nil)
	_ = b.InstallBuildableProperty(PropertyUseStock, cdk.BoolProperty, true, false)
	_ = b.InstallBuildableProperty(PropertyUseUnderline, cdk.BoolProperty, true, false)
	_ = b.InstallBuildableProperty(PropertyXAlign, cdk.FloatProperty, true, 0.5)
	_ = b.InstallBuildableProperty(PropertyYAlign, cdk.FloatProperty, true, 0.5)
	b.Connect(SignalSetProperty, ButtonSetPropertyHandle, b.setProperty)
	b.Connect(SignalLostFocus, ButtonLostFocusHandle, b.lostFocus)
	b.Connect(SignalGainedFocus, ButtonGainedFocusHandle, b.gainedFocus)
	b.Connect(SignalCdkEvent, ButtonCdkEventHandle, b.event)
	b.Connect(SignalInvalidate, ButtonInvalidateHandle, b.invalidate)
	b.Connect(SignalResize, ButtonResizeHandle, b.resize)
	b.Connect(SignalDraw, ButtonDrawHandle, b.draw)
	b.Invalidate()
	return false
}

// Build provides customizations to the Buildable system for Button Widgets.
func (b *CButton) Build(builder Builder, element *CBuilderElement) error {
	b.Freeze()
	defer b.Thaw()
	if name, ok := element.Attributes["id"]; ok {
		b.SetName(name)
	}
	if v, ok := element.Properties[PropertyUseStock.String()]; ok {
		b.SetUseStock(cstrings.IsTrue(v))
		b.SetUseUnderline(true)
	}
	if v, ok := element.Properties[PropertyLabel.String()]; ok {
		b.SetLabel(v)
	}
	for k, v := range element.Properties {
		switch cdk.Property(k) {
		case PropertyLabel:
		case PropertyUseStock:
		default:
			element.ApplyProperty(k, v)
		}
	}
	element.ApplySignals()
	return nil
}

// Activate emits a SignalActivate, returning TRUE if the event was handled
func (b *CButton) Activate() (value bool) {
	return b.Emit(SignalActivate, b) == enums.EVENT_STOP
}

// Clicked emits a SignalClicked
func (b *CButton) Clicked() enums.EventFlag {
	// TODO: button Clicked() is not defined well
	return b.Emit(SignalClicked, b)
}

// GetRelief is a convenience method for returning the relief property value
// See: SetRelief()
func (b *CButton) GetRelief() (value ReliefStyle) {
	if v, err := b.GetStructProperty(PropertyRelief); err != nil {
		b.LogErr(err)
	} else {
		var ok bool
		if value, ok = v.(ReliefStyle); !ok {
			b.LogError("value stored in relief property is not of type ReliefStyle")
		}
	}
	return
}

// SetRelief is a convenience method for updating the relief property value
//
// Note that usage of this within CTK is unimplemented at this time
func (b *CButton) SetRelief(newStyle ReliefStyle) {
	if err := b.SetStructProperty(PropertyRelief, newStyle); err != nil {
		b.LogErr(err)
	}
}

// GetLabel returns the text from the label of the button, as set by SetLabel.
// If the child Widget is not a Label, the value of the button label property
// will be returned instead.
// See: SetLabel()
func (b *CButton) GetLabel() (value string) {
	if v, ok := b.GetChild().(Label); ok {
		return v.GetText()
	}
	var err error
	if value, err = b.GetStringProperty(PropertyButtonLabel); err != nil {
		b.LogErr(err)
	}
	return
}

// SetLabel will update the text of the child Label of the button to the given
// text. This text is also used to select the stock item if SetUseStock is used.
// This will also clear any previously set labels.
//
// Parameters:
// 	label	the Label text to apply
func (b *CButton) SetLabel(label string) {
	if b.GetUseStock() && label != "" {
		label = strings.ReplaceAll(label, "gtk", "ctk")
		if item := LookupStockItem(StockID(label)); item != nil {
			label = item.Label
		}
	}
	if v, ok := b.GetChild().(Label); ok {
		if strings.HasPrefix(label, "<markup") {
			if err := v.SetMarkup(label); err != nil {
				b.LogErr(err)
			}
		} else {
			v.SetText(label)
		}
	}
}

// GetUseStock is a convenience method to return the use-stock property value.
// See: SetUseStock()
func (b *CButton) GetUseStock() (value bool) {
	var err error
	if value, err = b.GetBoolProperty(PropertyUseStock); err != nil {
		b.LogErr(err)
	}
	return
}

// SetUseStock is a convenience method to update the use-stock property value.
// If TRUE, the label set on the button is used as a stock id to select the
// stock item for the button.
//
// Parameters:
// 	useStock	TRUE if the button should use a stock item
func (b *CButton) SetUseStock(useStock bool) {
	if err := b.SetBoolProperty(PropertyUseStock, useStock); err != nil {
		b.LogErr(err)
	} else {
		label := b.GetLabel()
		b.SetLabel(label)
	}
}

// GetUseUnderline is a convenience method to return the use-underline property
// value. This is whether an embedded underline in the button label indicates a
// mnemonic.
// See: SetUseUnderline()
func (b *CButton) GetUseUnderline() (enabled bool) {
	var err error
	if enabled, err = b.GetBoolProperty(PropertyUseUnderline); err != nil {
		b.LogErr(err)
	}
	return
}

// SetUseUnderline is a convenience method to update the use-underline property
// value and update the child Label settings. If true, an underline in the text
// of the button label indicates the next character should be used for the
// mnemonic accelerator key.
//
// Parameters:
// 	useUnderline	TRUE if underlines in the text indicate mnemonics
func (b *CButton) SetUseUnderline(enabled bool) {
	if err := b.SetBoolProperty(PropertyUseUnderline, enabled); err != nil {
		b.LogErr(err)
	}
	if child := b.GetChild(); child != nil {
		if label, ok := child.(Label); ok {
			label.SetUseUnderline(enabled)
		}
	}
}

// GetUseMarkup is a convenience method to return the use-markup property value.
// This is whether markup in the label text is rendered.
// See: SetUseMarkup()
func (b *CButton) GetUseMarkup() (enabled bool) {
	var err error
	if enabled, err = b.GetBoolProperty(PropertyUseMarkup); err != nil {
		b.LogErr(err)
	}
	return
}

// SetUseMarkup is a convenience method to update the use-markup property value.
// If true, any Tango markup in the text of the button label will be rendered.
//
// Parameters:
// 	enabled	TRUE if markup is rendered
func (b *CButton) SetUseMarkup(enabled bool) {
	if err := b.SetBoolProperty(PropertyUseUnderline, enabled); err != nil {
		b.LogErr(err)
	}
	if child := b.GetChild(); child != nil {
		if label, ok := child.(Label); ok {
			label.SetUseMarkup(enabled)
		}
	}
}

// GetFocusOnClick is a convenience method to return the focus-on-click property
// value. This is whether the button grabs focus when it is clicked with the
// mouse.
// See: SetFocusOnClick()
func (b *CButton) GetFocusOnClick() (value bool) {
	var err error
	if value, err = b.GetBoolProperty(PropertyFocusOnClick); err != nil {
		b.LogErr(err)
	}
	return
}

// SetFocusOnClick is a convenience method for updating the focus-on-click
// property value. This is whether the button will grab focus when it is clicked
// with the mouse. Making mouse clicks not grab focus is useful in places like
// toolbars where you don't want the keyboard focus removed from the main area
// of the application.
//
// Parameters:
// 	focusOnClick	whether the button grabs focus when clicked with the mouse
func (b *CButton) SetFocusOnClick(focusOnClick bool) {
	if err := b.SetBoolProperty(PropertyFocusOnClick, focusOnClick); err != nil {
		b.LogErr(err)
	}
}

// GetAlignment is a convenience method for returning both the x and y alignment
// property values.
//
// Parameters:
// 	xAlign	horizontal alignment
// 	yAlign	vertical alignment
func (b *CButton) GetAlignment() (xAlign float64, yAlign float64) {
	var err error
	if xAlign, err = b.GetFloatProperty(PropertyXAlign); err != nil {
		b.LogErr(err)
	}
	err = nil
	if yAlign, err = b.GetFloatProperty(PropertyYAlign); err != nil {
		b.LogErr(err)
	}
	return
}

// SetAlignment is a convenience method for updating both the x and y alignment
// values. This property has no effect unless the child Widget implements the
// Alignable interface (ie: Misc based or Alignment Widget types).
//
// Parameters:
// 	xAlign	the horizontal position of the child, 0.0 is left aligned, 1.0 is right aligned
// 	yAlign	the vertical position of the child, 0.0 is top aligned, 1.0 is bottom aligned
func (b *CButton) SetAlignment(xAlign float64, yAlign float64) {
	xAlign = cmath.ClampF(xAlign, 0.0, 1.0)
	yAlign = cmath.ClampF(yAlign, 0.0, 1.0)
	if err := b.SetProperty(PropertyXAlign, xAlign); err != nil {
		b.LogErr(err)
	}
	if err := b.SetProperty(PropertyYAlign, yAlign); err != nil {
		b.LogErr(err)
	}
	if child := b.GetChild(); child != nil {
		if ca, ok := child.(Alignable); ok {
			ca.SetAlignment(xAlign, yAlign)
		}
	}
}

// GetImage is a convenience method to return the image property value.
// See: SetImage()
func (b *CButton) GetImage() (value Widget, ok bool) {
	if w, err := b.GetStructProperty(PropertyImage); err != nil {
		b.LogErr(err)
	} else {
		if value, ok = w.(Widget); !ok {
			value = nil
			return
		}
	}
	return
}

// SetImage is a convenience method to update the image property value.
//
// Parameters:
// 	image	a widget to set as the image for the button
//
// Note that usage of this within CTK is unimplemented at this time
func (b *CButton) SetImage(image Widget) {
	if err := b.SetStructProperty(PropertyImage, image); err != nil {
		b.LogErr(err)
	}
}

// GetImagePosition is a convenience method to return the image-position
// property value.
// See: SetImagePosition()
//
// Note that usage of this within CTK is unimplemented at this time
func (b *CButton) GetImagePosition() (value PositionType) {
	if v, err := b.GetStructProperty(PropertyImagePosition); err != nil {
		b.LogErr(err)
	} else {
		var ok bool
		if value, ok = v.(PositionType); !ok {
			b.LogError("value stored in PropertyImagePosition is not a PositionType: %T", v)
		}
	}
	return
}

// SetImagePosition is a convenience method to update the image-position
// property value. This sets the position of the image relative to the text
// inside the button.
//
// Parameters:
// 	position	the position
//
// Note that usage of this within CTK is unimplemented at this time
func (b *CButton) SetImagePosition(position PositionType) {
	if err := b.SetStructProperty(PropertyImagePosition, position); err != nil {
		b.LogErr(err)
	}
}

// GetPressed returns TRUE if the Button is currently pressed, FALSE otherwise.
func (b *CButton) GetPressed() bool {
	return b.pressed
}

// SetPressed is used to change the pressed state of the Button. If TRUE, the
// Button is flagged as pressed and a SignalPressed is emitted. If FALSE, the
// Button is flagged as not being pressed and a SignalReleased is emitted.
func (b *CButton) SetPressed(pressed bool) {
	b.pressed = pressed
	b.Invalidate()
	if pressed {
		b.Emit(SignalPressed)
	} else {
		b.Emit(SignalReleased)
	}
}

// GetFocusChain overloads the Container.GetFocusChain to always return the
// Button instance as the only item in the focus chain.
func (b *CButton) GetFocusChain() (focusableWidgets []interface{}, explicitlySet bool) {
	focusableWidgets = []interface{}{b}
	return
}

// GrabFocus will take the focus of the associated Window if the Widget instance
// CanFocus(). Any previously focused Widget will emit a lost-focus signal and
// the newly focused Widget will emit a gained-focus signal. This method emits a
// grab-focus signal initially and if the listeners return EVENT_PASS, the
// changes are applied.
func (b *CButton) GrabFocus() {
	if b.CanFocus() {
		if r := b.Emit(SignalGrabFocus, b); r == enums.EVENT_PASS {
			tl := b.GetWindow()
			if tl != nil {
				var fw Widget
				focused := tl.GetFocus()
				tl.SetFocus(b)
				if focused != nil {
					var ok bool
					if fw, ok = focused.(Widget); ok && fw.ObjectID() != b.ObjectID() {
						if f := fw.Emit(SignalLostFocus, fw); f == enums.EVENT_STOP {
							fw = nil
						}
					}
				}
				if f := b.Emit(SignalGainedFocus, b, fw); f == enums.EVENT_STOP {
					if fw != nil {
						tl.SetFocus(fw)
					}
				}
				b.LogDebug("has taken focus")
			}
		}
	}
}

// GrabEventFocus will emit a grab-event-focus signal and if all signal handlers
// return enums.EVENT_PASS will set the Button instance as the Window event
// focus handler.
func (b *CButton) GrabEventFocus() {
	if window := b.GetWindow(); window != nil {
		if f := b.Emit(SignalGrabEventFocus, b, window); f == enums.EVENT_PASS {
			window.SetEventFocus(b)
		}
	}
}

// CancelEvent emits a cancel-event signal and if the signal handlers all return
// enums.EVENT_PASS, then set the button as not pressed and release any event
// focus.
func (b *CButton) CancelEvent() {
	if f := b.Emit(SignalCancelEvent, b); f == enums.EVENT_PASS {
		b.SetPressed(false)
		b.ReleaseEventFocus()
	}
}

// GetWidgetAt returns the Button instance if the position given is within the
// allocated size at the origin point of the Button. If the position given is
// not contained within the Button space, `nil` is returned.
func (b *CButton) GetWidgetAt(p *ptypes.Point2I) Widget {
	if b.HasPoint(p) && b.IsVisible() {
		return b
	}
	return nil
}

// GetThemeRequest returns the current theme for the Button, reflecting the
// pressed state of the Button. This method is only to be used within draw
// signal handlers to render the current state of the Button.
func (b *CButton) GetThemeRequest() (theme paint.Theme) {
	theme = b.CWidget.GetThemeRequest()
	if b.GetPressed() {
		theme.Content.Normal = theme.Content.Active
		theme.Content.Focused = theme.Content.Active
		theme.Border.Normal = theme.Border.Active
		theme.Border.Focused = theme.Border.Active
	} else if b.IsFocused() {
		theme.Content.Normal = theme.Content.Focused
		theme.Border.Normal = theme.Border.Focused
	}
	return
}

// GetSizeRequest returns the requested size of the Drawable Widget. This method
// is used by Container Widgets to resolve the surface space allocated for their
// child Widget instances.
func (b *CButton) GetSizeRequest() (width, height int) {
	size := ptypes.NewRectangle(b.CWidget.GetSizeRequest())
	if child := b.GetChild(); child != nil {
		labelSizeReq := ptypes.NewRectangle(child.GetSizeRequest())
		if size.W <= -1 && labelSizeReq.W > -1 {
			size.W = 2 + labelSizeReq.W + 2 // borders and bookends
		}
		if size.H <= -1 && labelSizeReq.H > -1 {
			size.H = labelSizeReq.H + 2 // borders
		}
	}
	return size.W, size.H
}

func (b *CButton) getBorderRequest() (border bool) {
	border = true
	alloc := b.GetAllocation()
	if alloc.W <= 2 || alloc.H <= 2 {
		border = false
	}
	return
}

func (b *CButton) setProperty(data []interface{}, argv ...interface{}) enums.EventFlag {
	if len(argv) == 3 {
		if key, ok := argv[1].(cdk.Property); ok {
			switch key {
			case PropertyButtonLabel:
				if val, ok := argv[2].(string); ok {
					b.SetLabel(val)
				} else {
					b.LogError("property label value is not string: %T", argv[2])
				}
			}
		}
	}
	// allow property to be set by other signal handlers
	return enums.EVENT_PASS
}

func (b *CButton) lostFocus(data []interface{}, argv ...interface{}) enums.EventFlag {
	_ = b.Invalidate()
	return enums.EVENT_PASS
}

func (b *CButton) gainedFocus(data []interface{}, argv ...interface{}) enums.EventFlag {
	_ = b.Invalidate()
	return enums.EVENT_PASS
}

func (b *CButton) event(data []interface{}, argv ...interface{}) enums.EventFlag {
	if evt, ok := argv[1].(cdk.Event); ok {
		switch e := evt.(type) {
		case *cdk.EventMouse:
			pos := ptypes.NewPoint2I(e.Position())
			switch e.State() {
			case cdk.BUTTON_PRESS, cdk.DRAG_START:
				if b.HasPoint(pos) {
					if focusOnClick, err := b.GetBoolProperty(PropertyFocusOnClick); err == nil && focusOnClick {
						b.GrabFocus()
					}
					b.GrabEventFocus()
					b.SetPressed(true)
					b.LogDebug("pressed")
					return enums.EVENT_STOP
				}
			case cdk.MOUSE_MOVE, cdk.DRAG_MOVE:
				if b.HasEventFocus() {
					if !b.HasPoint(pos) {
						b.LogDebug("out of bounds")
						b.CancelEvent()
						return enums.EVENT_STOP
					}
				}
				return enums.EVENT_PASS
			case cdk.BUTTON_RELEASE, cdk.DRAG_STOP:
				if b.HasEventFocus() {
					if !b.HasPoint(pos) {
						b.LogDebug("out of bounds")
						b.CancelEvent()
						return enums.EVENT_STOP
					}
					b.ReleaseEventFocus()
					if f := b.Clicked(); f == enums.EVENT_PASS {
						b.Activate()
					}
					b.SetPressed(false)
					b.LogDebug("released")
					return enums.EVENT_STOP
				}
			}
		case *cdk.EventKey:
			if b.HasEventFocus() {
				b.LogDebug("keypress cancelling mouse event handling")
				b.CancelEvent()
				return enums.EVENT_STOP
			}
			switch e.Key() {
			case cdk.KeyRune:
				if e.Rune() != ' ' {
					break
				}
				fallthrough
			case cdk.KeyEnter:
				if focusOnClick, err := b.GetBoolProperty(PropertyFocusOnClick); err == nil && focusOnClick {
					b.GrabFocus()
				}
				b.LogTrace("pressed")
				b.SetPressed(true)
				if f := b.Clicked(); f == enums.EVENT_PASS {
					b.Activate()
				}
				b.SetPressed(false)
				b.LogTrace("released")
				return enums.EVENT_STOP
			}
		}
	}
	return enums.EVENT_PASS
}

func (b *CButton) invalidate(data []interface{}, argv ...interface{}) enums.EventFlag {
	if child := b.GetChild(); child != nil {
		theme := b.GetThemeRequest()
		child.SetTheme(theme)
		child.Invalidate()
	}
	return enums.EVENT_STOP
}

func (b *CButton) resize(data []interface{}, argv ...interface{}) enums.EventFlag {
	theme := b.GetThemeRequest()
	alloc := b.GetAllocation()
	size := ptypes.NewRectangle(alloc.W, alloc.H)
	origin := b.GetOrigin()
	child := b.GetChild()
	if child != nil {
		if alloc.W <= 0 || alloc.H <= 0 {
			child.SetAllocation(ptypes.MakeRectangle(0, 0))
			return child.Resize()
		}
		local := ptypes.NewPoint2I(0, 0)
		if alloc.W >= 3 && alloc.H >= 3 {
			local.Add(1, 1)
			size.Sub(2, 2)
		}
		req := ptypes.MakeRectangle(child.GetSizeRequest())
		req.Clamp(0, 0, size.W, size.H)
		if label, ok := child.(Label); ok {
			w, h := label.GetPlainTextInfoAtWidth(req.W)
			xAlign := 0.5
			yAlign := 0.5
			if h >= size.H {
				yAlign = 0.0
			}
			if w >= size.W {
				label.SetJustify(enums.JUSTIFY_LEFT)
				xAlign = 0.0
			} else {
				label.SetJustify(enums.JUSTIFY_CENTER)
			}
			label.SetAlignment(xAlign, yAlign)
		}
		child.SetOrigin(origin.X+local.X, origin.Y+local.Y)
		child.SetAllocation(*size)
		child.Resize()
		if err := memphis.ConfigureSurface(child.ObjectID(), *local, *size, theme.Content.Normal); err != nil {
			child.LogErr(err)
		}
	}
	b.Invalidate()
	return enums.EVENT_PASS
}

func (b *CButton) draw(data []interface{}, argv ...interface{}) enums.EventFlag {
	if surface, ok := argv[1].(*memphis.CSurface); ok {
		b.Lock()
		defer b.Unlock()
		size := b.GetAllocation()
		if !b.IsVisible() || size.W <= 0 || size.H <= 0 {
			b.LogTrace("Draw(%v): not visible, zero width or zero height", surface)
			surface.Fill(b.GetTheme())
			return enums.EVENT_STOP
		}

		var child Widget
		var label Label
		if child = b.GetChild(); child == nil {
			b.LogError("button child (label) not found")
			// return enums.EVENT_PASS
		} else if v, ok := child.(Label); ok {
			label = v
		}

		theme := b.GetThemeRequest()
		border := b.getBorderRequest()

		surface.Box(
			ptypes.MakePoint2I(0, 0),
			ptypes.MakeRectangle(size.W, size.H),
			border, true,
			theme.Content.Overlay,
			theme.Content.FillRune,
			theme.Content.Normal,
			theme.Border.Normal,
			theme.Border.BorderRunes,
		)

		if label != nil {
			if f := label.Draw(); f == enums.EVENT_STOP {
				if err := surface.Composite(label.ObjectID()); err != nil {
					b.LogError("composite error: %v", err)
				}
			}
		} else if child != nil {
			if f := child.Draw(); f == enums.EVENT_STOP {
				if err := surface.Composite(child.ObjectID()); err != nil {
					b.LogError("composite error: %v", err)
				}
			}
		}

		if debug, _ := b.GetBoolProperty(cdk.PropertyDebug); debug {
			surface.DebugBox(paint.ColorRed, b.ObjectInfo())
		}
		return enums.EVENT_STOP
	}
	return enums.EVENT_PASS
}

// Whether the button grabs focus when it is clicked with the mouse.
// Flags: Read / Write
// Default value: TRUE
const PropertyFocusOnClick cdk.Property = "focus-on-click"

// Child widget to appear next to the button text.
// Flags: Read / Write
const PropertyImage cdk.Property = "image"

// The position of the image relative to the text inside the button.
// Flags: Read / Write
// Default value: GTK_POS_LEFT
const PropertyImagePosition cdk.Property = "image-position"

// Text of the label widget inside the button, if the button contains a label
// widget.
// Flags: Read / Write / Construct
// Default value: NULL
const PropertyButtonLabel cdk.Property = "label"

// The border relief style.
// Flags: Read / Write
// Default value: GTK_RELIEF_NORMAL
const PropertyRelief cdk.Property = "relief"

// If set, the label is used to pick a stock item instead of being displayed.
// Flags: Read / Write / Construct
// Default value: FALSE
const PropertyUseStock cdk.Property = "use-stock"

// If set, an underline in the text indicates the next character should be
// used for the mnemonic accelerator key.
// Flags: Read / Write / Construct
// Default value: FALSE
const PropertyUseUnderline cdk.Property = "use-underline"

// If the child of the button is a Misc or Alignment, this property can
// be used to control it's horizontal alignment. 0.0 is left aligned, 1.0 is
// right aligned.
// Flags: Read / Write
// Allowed values: [0,1]
// Default value: 0.5
// const PropertyXAlign cdk.Property = "xalign"

// If the child of the button is a Misc or Alignment, this property can
// be used to control it's vertical alignment. 0.0 is top aligned, 1.0 is
// bottom aligned.
// Flags: Read / Write
// Allowed values: [0,1]
// Default value: 0.5
// const PropertyYAlign cdk.Property = "yalign"

// The activate signal on Button is an action signal and emitting it causes
// the button to animate press then release. Applications should never
// connect to this signal, but use the clicked signal.
// const SignalActivate cdk.Signal = "activate"

// Emitted when the button has been activated (pressed and released).
const SignalClicked cdk.Signal = "clicked"

// Emitted when the pointer enters the button.
const SignalEnter cdk.Signal = "enter"

// Emitted when the pointer leaves the button.
const SignalLeave cdk.Signal = "leave"

// Emitted when the button is pressed.
const SignalPressed cdk.Signal = "pressed"

// Emitted when the button is released.
const SignalReleased cdk.Signal = "released"

const ButtonSetPropertyHandle = "button-set-property-handler"

const ButtonLostFocusHandle = "button-lost-focus-handler"

const ButtonGainedFocusHandle = "button-gained-focus-handler"

const ButtonCdkEventHandle = "button-cdk-event-handler"

const ButtonInvalidateHandle = "button-invalidate-handler"

const ButtonResizeHandle = "button-resize-handler"

const ButtonDrawHandle = "button-draw-handler"
