package ctk

// TODO: mnemonics support, Accel

import (
	"strings"

	"github.com/go-curses/cdk"
	cenums "github.com/go-curses/cdk/lib/enums"
	cmath "github.com/go-curses/cdk/lib/math"
	"github.com/go-curses/cdk/lib/paint"
	"github.com/go-curses/cdk/lib/ptypes"
	cstrings "github.com/go-curses/cdk/lib/strings"
	"github.com/go-curses/cdk/memphis"

	"github.com/go-curses/ctk/lib/enums"
)

const (
	TypeButton       cdk.CTypeTag    = "ctk-button"
	ButtonMonoTheme  paint.ThemeName = "button-mono"
	ButtonColorTheme paint.ThemeName = "button-color"
)

func init() {
	_ = cdk.TypesManager.AddType(TypeButton, func() interface{} { return MakeButton() })

	borders, _ := paint.GetDefaultBorder(paint.StockBorder)
	arrows, _ := paint.GetDefaultArrow(paint.StockArrow)

	style := paint.GetDefaultColorStyle()
	styleNormal := style.Foreground(paint.ColorWhite).Background(paint.ColorFireBrick)
	styleActive := style.Foreground(paint.ColorWhite).Background(paint.ColorDarkRed)
	styleInsensitive := style.Foreground(paint.ColorDarkSlateGray).Background(paint.ColorRosyBrown)

	paint.SetDefaultTheme(ButtonColorTheme, paint.Theme{
		Content: paint.ThemeAspect{
			Normal:      styleNormal.Dim(true).Bold(false),
			Selected:    styleActive.Dim(false).Bold(true),
			Active:      styleActive.Dim(false).Bold(true).Reverse(true),
			Prelight:    styleActive.Dim(false),
			Insensitive: styleInsensitive.Dim(true),
			FillRune:    paint.DefaultFillRune,
			BorderRunes: borders,
			ArrowRunes:  arrows,
			Overlay:     false,
		},
		Border: paint.ThemeAspect{
			Normal:      styleNormal.Dim(true).Bold(false),
			Selected:    styleActive.Dim(false).Bold(true),
			Active:      styleActive.Dim(false).Bold(true).Reverse(true),
			Prelight:    styleActive.Dim(false),
			Insensitive: styleInsensitive.Dim(true),
			FillRune:    paint.DefaultFillRune,
			BorderRunes: borders,
			ArrowRunes:  arrows,
			Overlay:     false,
		},
	})

	style = paint.GetDefaultMonoStyle()
	styleNormal = style.Foreground(paint.ColorWhite).Background(paint.ColorBlack)
	styleActive = style.Foreground(paint.ColorBlack).Background(paint.ColorWhite)
	styleInsensitive = style.Foreground(paint.ColorLightGray).Background(paint.ColorDarkGray)

	paint.SetDefaultTheme(ButtonMonoTheme, paint.Theme{
		Content: paint.ThemeAspect{
			Normal:      styleNormal.Dim(true).Bold(false),
			Selected:    styleActive.Dim(false).Bold(true),
			Active:      styleActive.Dim(false).Bold(true).Reverse(true),
			Prelight:    styleActive.Dim(false),
			Insensitive: styleInsensitive.Dim(true),
			FillRune:    paint.DefaultFillRune,
			BorderRunes: borders,
			ArrowRunes:  arrows,
			Overlay:     false,
		},
		Border: paint.ThemeAspect{
			Normal:      styleNormal.Dim(true).Bold(false),
			Selected:    styleActive.Dim(false).Bold(true),
			Active:      styleActive.Dim(false).Bold(true).Reverse(true),
			Prelight:    styleActive.Dim(false),
			Insensitive: styleInsensitive.Dim(true),
			FillRune:    paint.DefaultFillRune,
			BorderRunes: borders,
			ArrowRunes:  arrows,
			Overlay:     false,
		},
	})
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
	Sensitive

	Activate() (value bool)
	Clicked() cenums.EventFlag
	GetRelief() (value enums.ReliefStyle)
	SetRelief(newStyle enums.ReliefStyle)
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
	GetImagePosition() (value enums.PositionType)
	SetImagePosition(position enums.PositionType)
	GetPressed() bool
	SetPressed(pressed bool)
	CancelEvent()
}

var _ Button = (*CButton)(nil)

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
func MakeButton() Button {
	b := NewButtonWithLabel("")
	return b
}

// NewButton is a constructor for new Button instances without a label Widget.
func NewButton() Button {
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
func NewButtonWithLabel(text string) (b Button) {
	b = NewButton()
	label := NewLabel(text)
	label.Show()
	label.SetTheme(b.GetTheme())
	label.UnsetFlags(enums.CAN_FOCUS)
	label.UnsetFlags(enums.CAN_DEFAULT)
	label.UnsetFlags(enums.RECEIVES_DEFAULT)
	label.SetLineWrap(false)
	label.SetLineWrapMode(cenums.WRAP_NONE)
	label.SetJustify(cenums.JUSTIFY_CENTER)
	label.SetAlignment(0.5, 0.5)
	label.SetSingleLineMode(true)
	b.Add(label)
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
func NewButtonWithMnemonic(text string) (b Button) {
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
func NewButtonFromStock(stockId StockID) (value Button) {
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

// NewButtonWithWidget creates a NewButton with the given Widget as the Button's
// child.
func NewButtonWithWidget(w Widget) Button {
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
	b.flags = enums.NULL_WIDGET_FLAG
	b.SetFlags(enums.SENSITIVE | enums.PARENT_SENSITIVE | enums.CAN_DEFAULT | enums.RECEIVES_DEFAULT | enums.CAN_FOCUS | enums.APP_PAINTABLE | enums.COMPOSITE_PARENT)

	_ = b.InstallBuildableProperty(PropertyFocusOnClick, cdk.BoolProperty, true, true)
	_ = b.InstallBuildableProperty(PropertyButtonLabel, cdk.StringProperty, true, nil)
	_ = b.InstallBuildableProperty(PropertyRelief, cdk.StructProperty, true, nil)
	_ = b.InstallBuildableProperty(PropertyUseStock, cdk.BoolProperty, true, false)
	_ = b.InstallBuildableProperty(PropertyUseUnderline, cdk.BoolProperty, true, false)
	_ = b.InstallBuildableProperty(PropertyXAlign, cdk.FloatProperty, true, 0.5)
	_ = b.InstallBuildableProperty(PropertyYAlign, cdk.FloatProperty, true, 0.5)

	b.pressed = false

	theme, _ := paint.GetDefaultTheme(ButtonColorTheme)
	b.SetTheme(theme)

	b.Connect(SignalSetProperty, ButtonSetPropertyHandle, b.setProperty)
	b.Connect(SignalCdkEvent, ButtonCdkEventHandle, b.event)
	b.Connect(SignalLostFocus, ButtonLostFocusHandle, b.lostFocus)
	b.Connect(SignalGainedFocus, ButtonGainedFocusHandle, b.gainedFocus)
	b.Connect(SignalInvalidate, ButtonInvalidateHandle, b.invalidate)
	b.Connect(SignalResize, ButtonResizeHandle, b.resize)
	b.Connect(SignalDraw, ButtonDrawHandle, b.draw)
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

func (b *CButton) Add(w Widget) {
	b.CBin.Add(w)
	w.SetTheme(b.GetTheme())
}

func (b *CButton) SetTheme(theme paint.Theme) {
	b.CBin.SetTheme(theme)
	if child := b.GetChild(); child != nil {
		child.SetTheme(theme)
	}
	WidgetRecurseInvalidate(b)
}

func (b *CButton) SetState(state enums.StateType) {
	b.CBin.SetState(state)
	if child := b.GetChild(); child != nil {
		child.SetState(state)
	}
	WidgetRecurseInvalidate(b)
}

func (b *CButton) UnsetState(state enums.StateType) {
	b.CBin.UnsetState(state)
	if child := b.GetChild(); child != nil {
		child.UnsetState(state)
	}
	WidgetRecurseInvalidate(b)
}

// Activate emits a SignalActivate, returning TRUE if the event was handled
func (b *CButton) Activate() (value bool) {
	if b.IsSensitive() {
		return b.Emit(SignalActivate, b) == cenums.EVENT_STOP
	}
	return false
}

// Clicked emits a SignalClicked
func (b *CButton) Clicked() cenums.EventFlag {
	// TODO: button Clicked() is not defined well
	if b.IsSensitive() {
		return b.Emit(SignalClicked, b)
	}
	return cenums.EVENT_PASS
}

// GetRelief is a convenience method for returning the relief property value
// See: SetRelief()
//
// Locking: read
func (b *CButton) GetRelief() (value enums.ReliefStyle) {
	if v, err := b.GetStructProperty(PropertyRelief); err != nil {
		b.LogErr(err)
	} else {
		var ok bool
		if value, ok = v.(enums.ReliefStyle); !ok {
			b.LogError("value stored in relief property is not of type ReliefStyle")
		}
	}
	return
}

// SetRelief is a convenience method for updating the relief property value
//
// Note that usage of this within CTK is unimplemented at this time
//
// Locking: write
func (b *CButton) SetRelief(newStyle enums.ReliefStyle) {
	if err := b.SetStructProperty(PropertyRelief, newStyle); err != nil {
		b.LogErr(err)
	}
}

// GetLabel returns the text from the label of the button, as set by SetLabel.
// If the child Widget is not a Label, the value of the button label property
// will be returned instead.
// See: SetLabel()
//
// Locking: read
func (b *CButton) GetLabel() (value string) {
	if v, ok := b.GetChild().Self().(Label); ok {
		return v.GetText()
	}
	b.RLock()
	defer b.RUnlock()
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
//
// Locking: write
func (b *CButton) SetLabel(label string) {
	if b.GetUseStock() && label != "" {
		b.Lock()
		label = strings.ReplaceAll(label, "gtk", "ctk")
		if item := LookupStockItem(StockID(label)); item != nil {
			label = item.Label
		}
		b.Unlock()
	}
	if v, ok := b.GetChild().Self().(Label); ok {
		if strings.HasPrefix(label, "<markup") {
			if err := v.SetMarkup(label); err != nil {
				b.LogErr(err)
			}
		} else {
			v.SetText(label)
		}
	} else if err := b.SetStringProperty(PropertyButtonLabel, label); err != nil {
		b.LogErr(err)
	}
	b.Invalidate()
}

// GetUseStock is a convenience method to return the use-stock property value.
// See: SetUseStock()
//
// Locking: read
func (b *CButton) GetUseStock() (value bool) {
	b.RLock()
	defer b.RUnlock()
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
		if label, ok := child.Self().(Label); ok {
			if enabled {
				label.SetMnemonicWidget(b)
			} else {
				label.SetMnemonicWidget(nil)
			}
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
	} else {
		if child := b.GetChild(); child != nil {
			if label, ok := child.(Label); ok {
				label.SetUseMarkup(enabled)
			}
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
		if ca, ok := child.Self().(Alignable); ok {
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
func (b *CButton) GetImagePosition() (value enums.PositionType) {
	if v, err := b.GetStructProperty(PropertyImagePosition); err != nil {
		b.LogErr(err)
	} else {
		var ok bool
		if value, ok = v.(enums.PositionType); !ok {
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
func (b *CButton) SetImagePosition(position enums.PositionType) {
	if err := b.SetStructProperty(PropertyImagePosition, position); err != nil {
		b.LogErr(err)
	}
}

// GetPressed returns TRUE if the Button is currently pressed, FALSE otherwise.
func (b *CButton) GetPressed() bool {
	b.RLock()
	defer b.RUnlock()
	return b.pressed
}

// SetPressed is used to change the pressed state of the Button. If TRUE, the
// Button is flagged as pressed and a SignalPressed is emitted. If FALSE, the
// Button is flagged as not being pressed and a SignalReleased is emitted.
func (b *CButton) SetPressed(pressed bool) {
	b.Lock()
	b.pressed = pressed
	b.Unlock()
	if pressed {
		b.SetState(enums.StateActive)
		b.Emit(SignalPressed)
	} else {
		b.UnsetState(enums.StateActive)
		b.Emit(SignalReleased)
	}
}

// GetFocusChain overloads the Container.GetFocusChain to always return the
// Button instance as the only item in the focus chain.
func (b *CButton) GetFocusChain() (focusableWidgets []Widget, explicitlySet bool) {
	focusableWidgets = []Widget{b}
	return
}

// CancelEvent emits a cancel-event signal and if the signal handlers all return
// cenums.EVENT_PASS, then set the button as not pressed and release any event
// focus.
func (b *CButton) CancelEvent() {
	if f := b.Emit(SignalCancelEvent, b); f == cenums.EVENT_PASS {
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

// GetSizeRequest returns the requested size of the Drawable Widget. This method
// is used by Container Widgets to resolve the surface space allocated for their
// child Widget instances.
func (b *CButton) GetSizeRequest() (width, height int) {
	req := ptypes.NewRectangle(b.CWidget.GetSizeRequest())
	if child := b.GetChild(); child != nil {
		childReq := ptypes.NewRectangle(child.GetSizeRequest())
		lw, lh := -1, -1
		if label, ok := child.Self().(Label); ok {
			lw, lh = label.GetPlainTextInfo()
		}
		if req.W <= -1 {
			if childReq.W > 0 {
				req.W = 1 + childReq.W + 1 // borders, bookends
			} else if lw > 0 {
				req.W = 1 + lw + 1 // borders, bookends
			}
		}
		if req.H <= -1 {
			if childReq.H > 0 {
				req.H = childReq.H
			} else if lh > 0 {
				req.H = 1 + lw + 1 // borders
			} else {
				req.H = 3
			}
		}
	}
	return req.W, req.H
}

func (b *CButton) getBorderRequest() (border bool) {
	border = true
	alloc := b.GetAllocation()
	if alloc.W <= 2 || alloc.H <= 2 {
		border = false
	}
	return
}

func (b *CButton) setProperty(data []interface{}, argv ...interface{}) cenums.EventFlag {
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
	return cenums.EVENT_PASS
}

func (b *CButton) lostFocus(data []interface{}, argv ...interface{}) cenums.EventFlag {
	WidgetRecurseInvalidate(b)
	return cenums.EVENT_PASS
}

func (b *CButton) gainedFocus(data []interface{}, argv ...interface{}) cenums.EventFlag {
	WidgetRecurseInvalidate(b)
	return cenums.EVENT_PASS
}

func (b *CButton) event(data []interface{}, argv ...interface{}) cenums.EventFlag {
	if !b.IsSensitive() {
		return cenums.EVENT_PASS
	}
	if evt, ok := argv[1].(cdk.Event); ok {
		switch e := evt.(type) {
		case *cdk.EventMouse:
			pos := ptypes.NewPoint2I(e.Position())
			switch e.State() {
			case cdk.BUTTON_PRESS, cdk.DRAG_START:
				if b.HasPoint(pos) {
					b.GrabEventFocus()
					b.SetPressed(true)
					b.LogDebug("pressed")
					return cenums.EVENT_STOP
				}
			case cdk.MOUSE_MOVE, cdk.DRAG_MOVE:
				if b.HasEventFocus() {
					if !b.HasPoint(pos) {
						b.LogDebug("out of bounds")
						b.CancelEvent()
						return cenums.EVENT_STOP
					}
				}
				return cenums.EVENT_PASS
			case cdk.BUTTON_RELEASE, cdk.DRAG_STOP:
				if b.HasEventFocus() {
					if !b.HasPoint(pos) {
						b.LogDebug("out of bounds")
						b.CancelEvent()
						return cenums.EVENT_STOP
					}
					b.ReleaseEventFocus()
					if focusOnClick, err := b.GetBoolProperty(PropertyFocusOnClick); err == nil && focusOnClick {
						b.GrabFocus()
					}
					if f := b.Clicked(); f == cenums.EVENT_PASS {
						b.Activate()
					}
					b.SetPressed(false)
					b.LogDebug("released")
					return cenums.EVENT_STOP
				}
			}
		case *cdk.EventKey:
			if b.HasEventFocus() {
				b.LogDebug("keypress cancelling mouse event handling")
				b.CancelEvent()
				return cenums.EVENT_STOP
			}
			switch cdk.Key(e.Rune()) {
			case cdk.KeyEnter, cdk.KeySpace:
				if focusOnClick, err := b.GetBoolProperty(PropertyFocusOnClick); err == nil && focusOnClick {
					b.GrabFocus()
				}
				b.LogTrace("pressed")
				b.SetPressed(true)
				if f := b.Clicked(); f == cenums.EVENT_PASS {
					b.Activate()
				}
				b.SetPressed(false)
				b.LogTrace("released")
				return cenums.EVENT_STOP
			}
		}
	}
	return cenums.EVENT_PASS
}

func (b *CButton) invalidate(data []interface{}, argv ...interface{}) cenums.EventFlag {
	if child := b.GetChild(); child != nil {
		WidgetRecurseInvalidate(child)
	}
	return cenums.EVENT_PASS
}

func (b *CButton) resize(data []interface{}, argv ...interface{}) cenums.EventFlag {

	alloc := b.GetAllocation()
	size := ptypes.NewRectangle(alloc.W, alloc.H)
	origin := b.GetOrigin()
	child := b.GetChild()

	var label Label = nil
	if childLabel, ok := child.(Label); child != nil && ok {
		label = childLabel
	}

	if child != nil {

		if alloc.W <= 0 || alloc.H <= 0 {
			child.SetAllocation(ptypes.MakeRectangle(0, 0))
			rv := child.Resize()
			return rv
		}

		local := ptypes.NewPoint2I(0, 0)
		if alloc.W >= 3 && alloc.H >= 3 {
			local.Add(1, 1)
			size.Sub(2, 2)
		} else if alloc.W >= 3 {
			local.Add(1, 0)
			size.Sub(2, 0)
		}

		req := ptypes.MakeRectangle(child.GetSizeRequest())
		req.Clamp(0, 0, size.W, size.H)

		if label != nil {
			req = ptypes.MakeRectangle(label.GetSizeRequest())
			req.Clamp(0, 0, size.W, size.H)
			w, h := label.GetPlainTextInfoAtWidth(req.W)
			xAlign := 0.5
			yAlign := 0.5
			if h >= size.H {
				yAlign = 0.0
			}
			if w >= size.W {
				label.SetJustify(cenums.JUSTIFY_LEFT)
				xAlign = 0.0
			} else {
				label.SetJustify(cenums.JUSTIFY_CENTER)
			}
			label.SetAlignment(xAlign, yAlign)
		}

		child.SetOrigin(origin.X+local.X, origin.Y+local.Y)
		child.SetAllocation(*size)
		child.Resize()
	}

	b.Invalidate()
	return cenums.EVENT_STOP
}

func (b *CButton) draw(data []interface{}, argv ...interface{}) cenums.EventFlag {

	if surface, ok := argv[1].(*memphis.CSurface); ok {
		size := b.GetAllocation()
		if !b.IsVisible() || size.W <= 0 || size.H <= 0 {
			b.LogTrace("not visible, zero width or zero height", surface)
			surface.Fill(b.GetTheme())
			return cenums.EVENT_STOP
		}

		var child Widget
		var label Label
		if child = b.GetChild(); child == nil {
			b.LogError("button child (label) not found")
		} else if v, ok := child.Self().(Label); ok {
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
			label.SetTheme(theme)
			label.Draw()
			label.LockDraw()
			if err := surface.Composite(label.ObjectID()); err != nil {
				b.LogError("composite error: %v", err)
			}
			label.UnlockDraw()
		} else if child != nil {
			child.SetTheme(theme)
			child.Draw()
			child.LockDraw()
			if err := surface.Composite(child.ObjectID()); err != nil {
				b.LogError("composite error: %v", err)
			}
			child.UnlockDraw()
		}

		if debug, _ := b.GetBoolProperty(cdk.PropertyDebug); debug {
			surface.DebugBox(paint.ColorRed, b.ObjectInfo())
		}

		return cenums.EVENT_STOP
	}
	return cenums.EVENT_PASS
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

// The ::button-press-event signal will be emitted when a button (typically
// from a mouse) is pressed. To receive this signal, the Window associated
// to the widget needs to enable the GDK_BUTTON_PRESS_MASK mask. This signal
// will be sent to the grab widget if there is one.
const SignalButtonPressEvent cdk.Signal = "button-press-event"

// The ::button-release-event signal will be emitted when a button (typically
// from a mouse) is released. To receive this signal, the Window
// associated to the widget needs to enable the GDK_BUTTON_RELEASE_MASK mask.
// This signal will be sent to the grab widget if there is one.
const SignalButtonReleaseEvent cdk.Signal = "button-release-event"

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

const ButtonCdkEventHandle = "button-cdk-event-handler"

const ButtonLostFocusHandle = "button-lost-focus-handler"

const ButtonGainedFocusHandle = "button-gained-focus-handler"

const ButtonInvalidateHandle = "button-invalidate-handler"

const ButtonResizeHandle = "button-resize-handler"

const ButtonDrawHandle = "button-draw-handler"