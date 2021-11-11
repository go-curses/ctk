package ctk

// TODO: new from stock id
// TODO: mnemonics support, GtkAccel?

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

// CDK type-tag for Button objects
const TypeButton cdk.CTypeTag = "ctk-button"

var (
	DefaultMonoButtonTheme = paint.Theme{
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
	DefaultColorButtonTheme = paint.Theme{
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
type Button interface {
	Bin
	Activatable
	Alignable
	Buildable

	Init() (already bool)
	Build(builder Builder, element *CBuilderElement) error
	Activate() (value bool)
	Clicked() enums.EventFlag
	SetRelief(newStyle ReliefStyle)
	GetRelief() (value ReliefStyle)
	GetLabel() (value string)
	SetLabel(label string)
	GetUseStock() (value bool)
	SetUseStock(useStock bool)
	GetUseUnderline() (enabled bool)
	SetUseUnderline(enabled bool)
	GetUseMarkup() (enabled bool)
	SetUseMarkup(enabled bool)
	SetFocusOnClick(focusOnClick bool)
	GetFocusOnClick() (value bool)
	SetAlignment(xAlign float64, yAlign float64)
	GetAlignment() (xAlign float64, yAlign float64)
	SetImage(image Widget)
	GetImage() (value Widget, ok bool)
	SetImagePosition(position PositionType)
	GetImagePosition() (value PositionType)
	Add(w Widget)
	Remove(w Widget)
	SetPressed(pressed bool)
	GetPressed() bool
	GrabFocus()
	GetFocusChain() (focusableWidgets []interface{}, explicitlySet bool)
	GetDefaultChildren() []Widget
	GetWidgetAt(p *ptypes.Point2I) Widget
	CancelEvent()
	GrabEventFocus()
	ProcessEvent(evt cdk.Event) enums.EventFlag
	Invalidate() enums.EventFlag
	GetThemeRequest() (theme paint.Theme)
	GetSizeRequest() (width, height int)
	Resize() enums.EventFlag
}

// The CButton structure implements the Button interface and is
// exported to facilitate type embedding with custom implementations. No member
// variables are exported as the interface methods are the only intended means
// of interacting with Button objects
type CButton struct {
	CBin

	pressed bool
}

// Default constructor for Button objects
func MakeButton() *CButton {
	b := NewButtonWithLabel("")
	return b
}

// Constructor for Button objects
func NewButton() *CButton {
	b := new(CButton)
	b.Init()
	return b
}

func NewButtonWithLabel(text string) (b *CButton) {
	b = NewButton()
	label := NewLabel(text)
	b.Add(label)
	label.SetTheme(DefaultColorButtonTheme)
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

// Creates a new Button containing a label. If characters in label are
// preceded by an underscore, they are underlined. If you need a literal
// underscore character in a label, use '__' (two underscores). The first
// underlined character represents a keyboard accelerator called a mnemonic.
// Pressing Alt and that key activates the button.
// Parameters:
// 	label	The text of the button, with an underscore in front of the
// mnemonic character
func NewButtonWithMnemonic(text string) (b *CButton) {
	b = NewButtonWithLabel(text)
	b.SetUseUnderline(true)
	return b
}

// Creates a new Button containing the image and text from a stock item.
// Some stock ids have preprocessor macros like GTK_STOCK_OK and
// GTK_STOCK_APPLY. If stock_id is unknown, then it will be treated as a
// mnemonic label (as for NewWithMnemonic).
// Parameters:
// 	stockId	the name of the stock item
// Returns:
// 	a new Button
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

// Constructor with Widget for Button objects, uncertain this works as expected
// due to struct type information loss on interface filter
func NewButtonWithWidget(w Widget) *CButton {
	b := new(CButton)
	b.Init()
	b.Add(w)
	return b
}

// Button object initialization. This must be called at least once to setup
// the necessary defaults and allocate any memory structures. Calling this more
// than once is safe though unnecessary. Only the first call will result in any
// effect upon the Button instance
func (b *CButton) Init() (already bool) {
	if b.InitTypeItem(TypeButton, b) {
		return true
	}
	b.CBin.Init()
	b.flags = NULL_WIDGET_FLAG
	b.SetFlags(SENSITIVE | PARENT_SENSITIVE | CAN_DEFAULT | RECEIVES_DEFAULT | CAN_FOCUS | APP_PAINTABLE)
	b.SetTheme(DefaultColorButtonTheme)
	b.pressed = false
	_ = b.InstallBuildableProperty(PropertyFocusOnClick, cdk.BoolProperty, true, true)
	_ = b.InstallBuildableProperty(PropertyButtonLabel, cdk.StringProperty, true, nil)
	_ = b.InstallBuildableProperty(PropertyRelief, cdk.StructProperty, true, nil)
	_ = b.InstallBuildableProperty(PropertyUseStock, cdk.BoolProperty, true, false)
	_ = b.InstallBuildableProperty(PropertyUseUnderline, cdk.BoolProperty, true, false)
	_ = b.InstallBuildableProperty(PropertyXAlign, cdk.FloatProperty, true, 0.5)
	_ = b.InstallBuildableProperty(PropertyYAlign, cdk.FloatProperty, true, 0.5)
	b.Connect(SignalLostFocus, ButtonLostFocusHandle, b.lostFocus)
	b.Connect(SignalGainedFocus, ButtonGainedFocusHandle, b.gainedFocus)
	b.Connect(SignalDraw, ButtonDrawHandle, b.draw)
	b.Connect(cdk.SignalSetProperty, ButtonSetPropertyHandle, func(data []interface{}, argv ...interface{}) enums.EventFlag {
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
		// allow property to be set
		return enums.EVENT_PASS
	})
	b.Invalidate()
	return false
}

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

func (b *CButton) Activate() (value bool) {
	return b.Emit(SignalActivate, b) == enums.EVENT_STOP
}

// TODO: button Clicked() is not defined well

func (b *CButton) Clicked() enums.EventFlag {
	return b.Emit(SignalClicked, b)
}

func (b *CButton) SetRelief(newStyle ReliefStyle) {
	if err := b.SetStructProperty(PropertyRelief, newStyle); err != nil {
		b.LogErr(err)
	}
}

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

// Fetches the text from the label of the button, as set by
// SetLabel. If the label text has not been set the return
// value will be NULL. This will be the case if you create an empty button
// with New to use as a container.
// Returns:
// 	The text of the label widget. This string is owned by the
// 	widget and must not be modified or freed.
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

// Sets the text of the label of the button to str . This text is also used
// to select the stock item if SetUseStock is used. This will
// also clear any previously set labels.
// Parameters:
// 	label	a string
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

// Returns whether the button label is a stock item.
// Returns:
// 	TRUE if the button label is used to select a stock item instead
// 	of being used directly as the label text.
func (b *CButton) GetUseStock() (value bool) {
	var err error
	if value, err = b.GetBoolProperty(PropertyUseStock); err != nil {
		b.LogErr(err)
	}
	return
}

// If TRUE, the label set on the button is used as a stock id to select the
// stock item for the button.
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

// Returns whether an embedded underline in the button label indicates a
// mnemonic. See SetUseUnderline.
// Returns:
// 	TRUE if an embedded underline in the button label indicates the
// 	mnemonic accelerator keys.
func (b *CButton) GetUseUnderline() (enabled bool) {
	var err error
	if enabled, err = b.GetBoolProperty(PropertyUseUnderline); err != nil {
		b.LogErr(err)
	}
	return
}

// Returns whether markup in the label text is rendered. See SetUseMarkup.
// Returns:
// 	TRUE if markup is rendered
func (b *CButton) GetUseMarkup() (enabled bool) {
	var err error
	if enabled, err = b.GetBoolProperty(PropertyUseMarkup); err != nil {
		b.LogErr(err)
	}
	return
}

// If true, an underline in the text of the button label indicates the next
// character should be used for the mnemonic accelerator key.
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

// If true, any tango markup in the text of the button label will be rendered.
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

// Sets whether the button will grab focus when it is clicked with the mouse.
// Making mouse clicks not grab focus is useful in places like toolbars where
// you don't want the keyboard focus removed from the main area of the
// application.
// Parameters:
// 	focusOnClick	whether the button grabs focus when clicked with the mouse
func (b *CButton) SetFocusOnClick(focusOnClick bool) {
	if err := b.SetBoolProperty(PropertyFocusOnClick, focusOnClick); err != nil {
		b.LogErr(err)
	}
}

// Returns whether the button grabs focus when it is clicked with the mouse.
// See SetFocusOnClick.
// Returns:
// 	TRUE if the button grabs focus when it is clicked with the
// 	mouse.
func (b *CButton) GetFocusOnClick() (value bool) {
	var err error
	if value, err = b.GetBoolProperty(PropertyFocusOnClick); err != nil {
		b.LogErr(err)
	}
	return
}

// Sets the alignment of the child. This property has no effect unless the
// child is a Misc or a Alignment.
// Parameters:
// 	xAlign	the horizontal position of the child, 0.0 is left aligned,
// 1.0 is right aligned
// 	yAlign	the vertical position of the child, 0.0 is top aligned,
// 1.0 is bottom aligned
func (b *CButton) SetAlignment(xAlign float64, yAlign float64) {
	xAlign = cmath.ClampF(xAlign, 0.0, 1.0)
	yAlign = cmath.ClampF(yAlign, 0.0, 1.0)
	if err := b.SetProperty(PropertyXAlign, xAlign); err != nil {
		b.LogErr(err)
	}
	if err := b.SetProperty(PropertyYAlign, yAlign); err != nil {
		b.LogErr(err)
	}
}

// Gets the alignment of the child in the button.
// Parameters:
// 	xAlign	return location for horizontal alignment.
// 	yAlign	return location for vertical alignment.
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

// Set the image of button to the given widget. Note that it depends on the
// gtk-button-images setting whether the image will be displayed or
// not, you don't have to call WidgetShow on image yourself.
// Parameters:
// 	image	a widget to set as the image for the button
func (b *CButton) SetImage(image Widget) {
	if err := b.SetStructProperty(PropertyImage, image); err != nil {
		b.LogErr(err)
	}
}

// Gets the widget that is currently set as the image of button . This may
// have been explicitly set by SetImage or constructed by
// NewFromStock.
// Returns:
// 	a Widget or NULL in case there is no image.
// 	[transfer none]
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

// Sets the position of the image relative to the text inside the button.
// Parameters:
// 	position	the position
func (b *CButton) SetImagePosition(position PositionType) {
	if err := b.SetStructProperty(PropertyImagePosition, position); err != nil {
		b.LogErr(err)
	}
}

// Gets the position of the image relative to the text inside the button.
// Returns:
// 	the position
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

func (b *CButton) Add(w Widget) {
	if len(b.children) == 0 {
		b.CBin.Add(w)
		b.Invalidate()
	} else {
		b.LogError("button bin is full, failed to add: %v", w.ObjectName())
	}
}

func (b *CButton) Remove(w Widget) {
	if len(b.children) > 0 {
		b.CBin.Remove(w)
		b.Invalidate()
	} else {
		b.LogError("button bin is empty, failed to remove: %v", w.ObjectName())
	}
}

func (b *CButton) SetPressed(pressed bool) {
	b.pressed = pressed
	b.Invalidate()
	if pressed {
		b.Emit(SignalPressed)
	} else {
		b.Emit(SignalReleased)
	}
}

func (b *CButton) GetPressed() bool {
	return b.pressed
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

func (b *CButton) GetFocusChain() (focusableWidgets []interface{}, explicitlySet bool) {
	focusableWidgets = []interface{}{b}
	return
}

func (b *CButton) GetDefaultChildren() []Widget {
	return []Widget{b}
}

func (b *CButton) GetWidgetAt(p *ptypes.Point2I) Widget {
	if b.HasPoint(p) && b.IsVisible() {
		return b
	}
	return nil
}

func (b *CButton) CancelEvent() {
	if f := b.Emit(SignalCancelEvent, b); f == enums.EVENT_PASS {
		b.SetPressed(false)
		b.ReleaseEventFocus()
	}
}

func (b *CButton) GrabEventFocus() {
	if window := b.GetWindow(); window != nil {
		if f := b.Emit(SignalGrabEventFocus, b, window); f == enums.EVENT_PASS {
			window.SetEventFocus(b)
		}
	}
}

func (b *CButton) ProcessEvent(evt cdk.Event) enums.EventFlag {
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
	return enums.EVENT_PASS
}

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

func (b *CButton) getBorderRequest() (border bool) {
	border = true
	alloc := b.GetAllocation()
	if alloc.W <= 2 || alloc.H <= 2 {
		border = false
	}
	return
}

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

func (b *CButton) Resize() enums.EventFlag {
	// our allocation has been set prior to Resize() being called
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

func (b *CButton) Invalidate() enums.EventFlag {
	if child := b.GetChild(); child != nil {
		theme := b.GetThemeRequest()
		child.SetTheme(theme)
		child.Invalidate()
	}
	return enums.EVENT_STOP
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

func (b *CButton) lostFocus(data []interface{}, argv ...interface{}) enums.EventFlag {
	_ = b.Invalidate()
	return enums.EVENT_PASS
}

func (b *CButton) gainedFocus(data []interface{}, argv ...interface{}) enums.EventFlag {
	_ = b.Invalidate()
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

const ButtonLostFocusHandle = "button-lost-focus-handler"
const ButtonGainedFocusHandle = "button-gained-focus-handler"
const ButtonDrawHandle = "button-draw-handler"
const ButtonSetPropertyHandle = "button-set-property-handler"
