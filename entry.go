package ctk

import (
	"strings"
	"time"

	"github.com/gofrs/uuid"

	"github.com/go-curses/cdk/lib/sync"

	"github.com/go-curses/cdk"
	cenums "github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/cdk/lib/paint"
	"github.com/go-curses/cdk/lib/ptypes"
	cstrings "github.com/go-curses/cdk/lib/strings"
	"github.com/go-curses/cdk/memphis"

	"github.com/go-curses/ctk/lib/enums"
)

const TypeEntry cdk.CTypeTag = "ctk-entry"

func init() {
	_ = cdk.TypesManager.AddType(TypeEntry, func() interface{} { return MakeEntry() })
	ctkBuilderTranslators[TypeEntry] = func(builder Builder, widget Widget, name, value string) error {
		switch strings.ToLower(name) {
		case "wrap":
			isTrue := cstrings.IsTrue(value)
			if err := widget.SetBoolProperty(PropertyWrap, isTrue); err != nil {
				return err
			}
			if isTrue {
				if wmi, err := widget.GetStructProperty(PropertyWrapMode); err == nil {
					if wm, ok := wmi.(cenums.WrapMode); ok {
						if wm == cenums.WRAP_NONE {
							if err := widget.SetStructProperty(PropertyWrapMode, cenums.WRAP_WORD); err != nil {
								widget.LogErr(err)
							}
						}
					}
				}
			}
			return nil
		}
		return ErrFallthrough
	}
}

var (
	_ Editable  = (*CEntry)(nil)
	_ TextField = (*CEntry)(nil)
)

var (
	DefaultTextFieldTheme = paint.Theme{
		Content: paint.ThemeAspect{
			Normal:      paint.DefaultColorStyle.Foreground(paint.ColorWhite).Background(paint.ColorDarkSlateGray).Dim(true).Bold(false),
			Selected:    paint.DefaultColorStyle.Foreground(paint.ColorWhite).Background(paint.ColorDarkSlateGray).Dim(false).Bold(true),
			Active:      paint.DefaultColorStyle.Foreground(paint.ColorWhite).Background(paint.ColorDarkSlateGray).Dim(false).Bold(true).Reverse(true),
			Prelight:    paint.DefaultColorStyle.Foreground(paint.ColorWhite).Background(paint.ColorDarkSlateGray).Dim(false).Bold(false),
			Insensitive: paint.DefaultColorStyle.Foreground(paint.ColorWhite).Background(paint.ColorDarkSlateGray).Dim(false).Bold(false),
			FillRune:    paint.DefaultFillRune,
			BorderRunes: paint.DefaultBorderRune,
			ArrowRunes:  paint.DefaultArrowRune,
			Overlay:     false,
		},
		Border: paint.ThemeAspect{
			Normal:      paint.DefaultColorStyle.Foreground(paint.ColorWhite).Background(paint.ColorDarkSlateGray).Dim(false).Bold(false),
			Selected:    paint.DefaultColorStyle.Foreground(paint.ColorWhite).Background(paint.ColorDarkSlateGray).Dim(false).Bold(true),
			Active:      paint.DefaultColorStyle.Foreground(paint.ColorWhite).Background(paint.ColorDarkSlateGray).Dim(false).Bold(true).Reverse(true),
			Prelight:    paint.DefaultColorStyle.Foreground(paint.ColorWhite).Background(paint.ColorDarkSlateGray).Dim(false),
			Insensitive: paint.DefaultColorStyle.Foreground(paint.ColorWhite).Background(paint.ColorGray).Dim(false),
			FillRune:    paint.DefaultFillRune,
			BorderRunes: paint.DefaultBorderRune,
			ArrowRunes:  paint.DefaultArrowRune,
			Overlay:     false,
		},
	}
)

// TextField Hierarchy:
//	Object
//	  +- Widget
//	    +- Misc
//	      +- TextField
//	        +- AccelLabel
//	        +- TipsQuery
//
// The TextField Widget presents text to the end user.
type TextField interface {
	Misc
	Alignable
	Buildable
	Editable
	Sensitive

	Init() (already bool)
	Build(builder Builder, element *CBuilderElement) error
	SetText(text string)
	SetAttributes(attrs paint.Style)
	SetJustify(justify cenums.Justification)
	SetWidthChars(nChars int)
	SetMaxWidthChars(nChars int)
	SetLineWrap(wrap bool)
	SetLineWrapMode(wrapMode cenums.WrapMode)
	GetSelectable() (value bool)
	GetText() (value string)
	SelectRegion(startOffset int, endOffset int)
	SetSelectable(setting bool)
	GetAttributes() (value paint.Style)
	GetJustify() (value cenums.Justification)
	GetWidthChars() (value int)
	GetMaxWidthChars() (value int)
	GetLineWrap() (value bool)
	GetLineWrapMode() (value cenums.WrapMode)
	GetSingleLineMode() (value bool)
	SetSingleLineMode(singleLineMode bool)
	Settings() (singleLineMode bool, lineWrapMode cenums.WrapMode, justify cenums.Justification, maxWidthChars int)
	GetSizeRequest() (width, height int)
	CancelEvent()
}

type cTextFieldChange struct {
	name string
	argv []interface{}
}

// The CTextField structure implements the TextField interface and is exported
// to facilitate type embedding with custom implementations. No member variables
// are exported as the interface methods are the only intended means of
// interacting with TextField objects.
type CEntry struct {
	CMisc

	tid     uuid.UUID
	tRegion ptypes.Region

	offset    *ptypes.Region
	cursor    *ptypes.Point2I
	selection *ptypes.Point2I
	position  int
	queue     []*cTextFieldChange
	qLock     *sync.RWMutex
	qTimer    uuid.UUID

	tProfile *memphis.TextProfile
	tBuffer  memphis.TextBuffer
	tbStyle  paint.Style
}

// MakeEntry is used by the Buildable system to construct a new TextField.
func MakeEntry() TextField {
	return NewTextField("")
}

// NewTextField is the constructor for new TextField instances.
func NewTextField(plain string) TextField {
	l := new(CEntry)
	l.Init()
	l.SetText(plain)
	return l
}

// Init initializes a TextField object. This must be called at least once to
// set up the necessary defaults and allocate any memory structures. Calling
// this more than once is safe though unnecessary. Only the first call will
// result in any effect upon the TextField instance. Init is used in the
// NewTextField constructor and only necessary when implementing a derivative
// TextField type.
func (l *CEntry) Init() (already bool) {
	if l.InitTypeItem(TypeEntry, l) {
		return true
	}
	l.CMisc.Init()
	l.flags = enums.NULL_WIDGET_FLAG
	l.SetFlags(enums.SENSITIVE | enums.PARENT_SENSITIVE | enums.CAN_DEFAULT | enums.APP_PAINTABLE | enums.CAN_FOCUS)
	l.SetTheme(DefaultTextFieldTheme)
	_ = l.InstallProperty(PropertyAttributes, cdk.StructProperty, true, nil)
	_ = l.InstallProperty(PropertyJustify, cdk.StructProperty, true, cenums.JUSTIFY_NONE)
	_ = l.InstallProperty(PropertyText, cdk.StringProperty, true, "")
	_ = l.InstallProperty(PropertyMaxWidthChars, cdk.IntProperty, true, -1)
	_ = l.InstallProperty(PropertySelectable, cdk.BoolProperty, true, false)
	_ = l.InstallProperty(PropertySingleLineMode, cdk.BoolProperty, true, false)
	_ = l.InstallProperty(PropertyWidthChars, cdk.IntProperty, true, -1)
	_ = l.InstallProperty(PropertyWrap, cdk.BoolProperty, true, false)
	_ = l.InstallProperty(PropertyWrapMode, cdk.StructProperty, true, cenums.WRAP_NONE)
	_ = l.InstallProperty(PropertyEditable, cdk.BoolProperty, true, true)
	l.selection = nil
	l.position = 0
	l.queue = make([]*cTextFieldChange, 0)
	l.qLock = &sync.RWMutex{}
	l.qTimer = uuid.Nil
	l.offset = ptypes.NewRegion(0, 0, 0, 0)
	l.cursor = ptypes.NewPoint2I(0, 0)
	l.tProfile = memphis.NewTextProfile("")
	l.tBuffer = nil
	l.tid, _ = uuid.NewV4()
	l.tRegion = ptypes.MakeRegion(0, 0, 0, 0)
	if err := memphis.MakeSurface(l.tid, l.tRegion.Origin(), l.tRegion.Size(), paint.DefaultColorStyle); err != nil {
		l.LogErr(err)
	}
	l.Connect(SignalCdkEvent, TextFieldEventHandle, l.event)
	l.Connect(SignalLostFocus, TextFieldLostFocusHandle, l.lostFocus)
	l.Connect(SignalGainedFocus, TextFieldGainedFocusHandle, l.gainedFocus)
	l.Connect(SignalInvalidate, TextFieldInvalidateHandle, l.invalidate)
	l.Connect(SignalResize, TextFieldResizeHandle, l.resize)
	l.Connect(SignalDraw, TextFieldDrawHandle, l.draw)
	// _ = l.SetBoolProperty(PropertyDebug, true)
	l.Invalidate()
	return false
}

// Build provides customizations to the Buildable system for TextField Widgets.
func (l *CEntry) Build(builder Builder, element *CBuilderElement) error {
	l.Freeze()
	defer l.Thaw()
	if name, ok := element.Attributes["id"]; ok {
		l.SetName(name)
	}
	for k, v := range element.Properties {
		switch cdk.Property(k) {
		case PropertyText:
			l.SetText(v)
		default:
			element.ApplyProperty(k, v)
		}
	}
	element.ApplySignals()
	return nil
}

// SetText updates the text within the TextField widget. It overwrites any text that
// was there before. This will also clear any previously set mnemonic
// accelerators.
//
// Parameters:
// 	text	the text you want to set
//
// Locking: write
func (l *CEntry) SetText(text string) {
	l.setText(text)
	l.Invalidate()
	l.updateCursor()
}

func (l *CEntry) setText(text string) {
	l.Lock()
	l.tProfile.Set(text)
	l.Unlock()
	if err := l.SetStringProperty(PropertyText, l.tProfile.Get()); err != nil {
		l.LogErr(err)
	}
}

// SetAttributes updates the attributes property to be the given paint.Style.
//
// Parameters:
// 	attrs	a paint.Style
//
// Locking: write
func (l *CEntry) SetAttributes(attrs paint.Style) {
	if err := l.SetStructProperty(PropertyAttributes, attrs); err != nil {
		l.LogErr(err)
	}
}

// SetJustify updates the alignment of the lines in the text of the label
// relative to each other. JUSTIFY_LEFT is the default value when the widget is
// first created with New. If you instead want to set the alignment of the label
// as a whole, use SetAlignment instead.
//
// SetJustify has no effect on labels containing only a single line.
//
// Parameters:
// 	jtype	a Justification
//
// Locking: write
func (l *CEntry) SetJustify(justify cenums.Justification) {
	if err := l.SetStructProperty(PropertyJustify, justify); err != nil {
		l.LogErr(err)
	}
}

// SetWidthChars updates the desired width in characters of label to nChars.
//
// Parameters:
// 	nChars	the new desired width, in characters.
//
// Locking: write
func (l *CEntry) SetWidthChars(nChars int) {
	if err := l.SetIntProperty(PropertyWidthChars, nChars); err != nil {
		l.LogErr(err)
	}
}

// SetMaxWidthChars updates the desired maximum width in characters of label to
// nChars.
//
// Parameters:
// 	nChars	the new desired maximum width, in characters.
//
// Locking: write
func (l *CEntry) SetMaxWidthChars(nChars int) {
	if err := l.SetIntProperty(PropertyMaxWidthChars, nChars); err != nil {
		l.LogErr(err)
	}
}

// SetLineWrap updates the line wrapping within the TextField widget. TRUE makes it
// break lines if text exceeds the widget's size. FALSE lets the text get cut
// off by the edge of the widget if it exceeds the widget size. Note that
// setting line wrapping to TRUE does not make the label wrap at its parent
// container's width, because CTK widgets conceptually can't make their
// requisition depend on the parent container's size. For a label that wraps
// at a specific position, set the label's width using SetSizeRequest.
//
// Parameters:
// 	wrap	the setting
//
// Locking: write
func (l *CEntry) SetLineWrap(wrap bool) {
	if err := l.SetBoolProperty(PropertyWrap, wrap); err != nil {
		l.LogErr(err)
	}
}

// SetLineWrapMode updates the line wrapping if line-wrap is on (see
// SetLineWrap) this controls how the line wrapping is done. The default is
// WRAP_WORD which means wrap on word boundaries.
//
// Parameters:
// 	wrapMode	the line wrapping mode
//
// Locking: write
func (l *CEntry) SetLineWrapMode(wrapMode cenums.WrapMode) {
	if err := l.SetStructProperty(PropertyWrapMode, wrapMode); err != nil {
		l.LogErr(err)
	}
}

// GetSelectable returns the value set by SetSelectable.
//
// Locking: read
func (l *CEntry) GetSelectable() (value bool) {
	var err error
	if value, err = l.GetBoolProperty(PropertySelectable); err != nil {
		l.LogErr(err)
	}
	return
}

// GetText returns the text from a label widget, as displayed on the screen.
// This does not include any embedded underlines indicating mnemonics or Tango
// markup.
// See: GetLabel
//
// Locking: read
func (l *CEntry) GetText() (value string) {
	return l.tProfile.Get()
}

// SelectRegion selects a range of characters in the label, if the label is
// selectable. If the label is not selectable, this function has no effect. If
// start_offset or end_offset are -1, then the end of the label will be
// substituted.
// See: SetSelectable()
//
// Parameters:
// 	startOffset	start offset (in characters not bytes)
// 	endOffset	end offset (in characters not bytes)
//
func (l *CEntry) SelectRegion(startOffset int, endOffset int) {
	if l.GetSelectable() {
		l.selection = ptypes.NewPoint2I(startOffset, endOffset)
	}
}

// SetSelectable updates the selectable property for the TextField. TextFields allow the
// user to select text from the label, for copy-and-paste.
//
// Parameters:
// 	setting	TRUE to allow selecting text in the label
//
// Note that usage of this within CTK is unimplemented at this time
//
// Locking: write
func (l *CEntry) SetSelectable(setting bool) {
	if err := l.SetBoolProperty(PropertySelectable, setting); err != nil {
		l.LogErr(err)
	}
}

// GetAttributes returns the attribute list that was set on the label using
// SetAttributes, if any. This function does not reflect attributes that come
// from the TextField markup (see SetMarkup).
//
// Locking: read
func (l *CEntry) GetAttributes() (value paint.Style) {
	var ok bool
	if v, err := l.GetStructProperty(PropertyAttributes); err != nil {
		l.LogErr(err)
	} else if value, ok = v.(paint.Style); !ok {
		l.LogError("value stored in PropertyAttributes is not of paint.Style type: %v (%T)", v, v)
	}
	return
}

// GetJustify returns the justification of the label.
// See: SetJustify()
//
// Locking: read
func (l *CEntry) GetJustify() (value cenums.Justification) {
	var ok bool
	if v, err := l.GetStructProperty(PropertyJustify); err != nil {
		l.LogErr(err)
	} else if value, ok = v.(cenums.Justification); !ok {
		l.LogError("value stored in PropertyJustify is not of cenums.Justification type: %v (%T)", v, v)
	}
	return
}

// GetWidthChars retrieves the desired width of label, in characters.
// See: SetWidthChars()
//
// Locking: read
func (l *CEntry) GetWidthChars() (value int) {
	var err error
	if value, err = l.GetIntProperty(PropertyWidthChars); err != nil {
		l.LogErr(err)
	}
	return
}

// GetMaxWidthChars retrieves the desired maximum width of label, in characters.
// See: SetWidthChars()
//
// Locking: read
func (l *CEntry) GetMaxWidthChars() (value int) {
	var err error
	if value, err = l.GetIntProperty(PropertyMaxWidthChars); err != nil {
		l.LogErr(err)
	}
	return
}

// GetLineWrap returns whether lines in the label are automatically wrapped.
// See: SetLineWrap()
//
// Locking: read
func (l *CEntry) GetLineWrap() (value bool) {
	var err error
	if value, err = l.GetBoolProperty(PropertyWrap); err != nil {
		l.LogErr(err)
	}
	return
}

// GetLineWrapMode returns line wrap mode used by the label.
// See: SetLineWrapMode()
//
// Locking: read
func (l *CEntry) GetLineWrapMode() (value cenums.WrapMode) {
	var ok bool
	if v, err := l.GetStructProperty(PropertyWrapMode); err != nil {
		l.LogErr(err)
	} else if value, ok = v.(cenums.WrapMode); !ok {
		l.LogError("value stored in PropertyWrap is not of cenums.WrapMode type: %v (%T)", v, v)
	}
	return
}

// GetSingleLineMode returns whether the label is in single line mode.
//
// Locking: read
func (l *CEntry) GetSingleLineMode() (value bool) {
	var err error
	if value, err = l.GetBoolProperty(PropertySingleLineMode); err != nil {
		l.LogErr(err)
	}
	return
}

// SetSingleLineMode updates whether the label is in single line mode.
//
// Parameters:
// 	singleLineMode	TRUE if the label should be in single line mode
//
// Locking: write
func (l *CEntry) SetSingleLineMode(singleLineMode bool) {
	if err := l.SetBoolProperty(PropertySingleLineMode, singleLineMode); err != nil {
		l.LogErr(err)
	} else {
		l.Invalidate()
	}
}

// Settings is a convenience method to return the interesting settings currently
// configured on the TextField instance.
//
// Locking: read
func (l *CEntry) Settings() (singleLineMode bool, lineWrapMode cenums.WrapMode, justify cenums.Justification, maxWidthChars int) {
	singleLineMode = l.GetSingleLineMode()
	lineWrapMode = l.GetLineWrapMode()
	justify = l.GetJustify()
	maxWidthChars = l.GetMaxWidthChars()
	return
}

func (l *CEntry) GetSelectionBounds() (startPos, endPos int, ok bool) {
	if l.selection != nil {
		startPos = l.selection.X
		endPos = l.selection.Y
		ok = true
	}
	return
}

func (l *CEntry) InsertText(newText string, position int) {
	l.insertText(newText, position)
	l.Invalidate()
	l.updateCursor()
}

func (l *CEntry) insertText(newText string, position int) {
	if modified, ok := l.tProfile.Insert(newText, position); ok {
		if err := l.SetStringProperty(PropertyText, modified); err != nil {
			l.LogErr(err)
		}
	}
}

func (l *CEntry) DeleteText(startPos int, endPos int) {
	l.deleteText(startPos, endPos)
	l.Invalidate()
	l.updateCursor()
}

func (l *CEntry) deleteText(startPos int, endPos int) {
	if modified, ok := l.tProfile.Delete(startPos, endPos); ok {
		if err := l.SetStringProperty(PropertyText, modified); err != nil {
			l.LogErr(err)
		}
	}
}

func (l *CEntry) GetChars(startPos int, endPos int) (value string) {
	content := l.GetText()
	contentLength := len(content)
	if startPos >= contentLength {
		return
	}
	if contentLength <= endPos {
		endPos = contentLength - 1
	}
	value = content[startPos:endPos]
	return
}

func (l *CEntry) CutClipboard() {
	// TODO implement me
	panic("implement me")
}

func (l *CEntry) CopyClipboard() {
	// TODO implement me
	panic("implement me")
}

func (l *CEntry) PasteClipboard() {
	// TODO implement me
	panic("implement me")
}

func (l *CEntry) DeleteSelection() {
	// TODO implement me
	panic("implement me")
}

func (l *CEntry) SetPosition(position int) {
	l.setPosition(position)
	l.Invalidate()
	l.updateCursor()
}

func (l *CEntry) setPosition(position int) {
	l.Lock()
	max := l.tProfile.Len()
	if position > max {
		position = max
	}
	l.position = position
	l.Unlock()
}

func (l *CEntry) GetPosition() (value int) {
	l.RLock()
	defer l.RUnlock()
	return l.position
}

func (l *CEntry) SetEditable(isEditable bool) {
	if err := l.SetBoolProperty(PropertyEditable, isEditable); err != nil {
		l.LogErr(err)
	}
}

func (l *CEntry) GetEditable() (value bool) {
	var err error
	if value, err = l.GetBoolProperty(PropertyEditable); err != nil {
		l.LogErr(err)
	}
	return
}

// GetSizeRequest returns the requested size of the TextField taking into account
// the label's content and any padding set.
//
// Locking: read
func (l *CEntry) GetSizeRequest() (width, height int) {
	alloc := l.GetAllocation()
	size := l.CWidget.SizeRequest()
	if alloc.W > 0 && size.W > alloc.W {
		size.W = alloc.W
	}
	if alloc.H > 0 && size.H > alloc.H {
		size.H = alloc.H
	}
	return size.W, size.H
}

// CancelEvent emits a cancel-event signal and if the signal handlers all return
// cenums.EVENT_PASS, then set the button as not pressed and release any event
// focus.
func (l *CEntry) CancelEvent() {
	l.LogDebug("hit cancel event")
}

// Activate emits a SignalActivate, returning TRUE if the event was handled
func (l *CEntry) Activate() (value bool) {
	return l.Emit(SignalActivate, l) == cenums.EVENT_STOP
}

func (l *CEntry) getMaxCharsRequest() (maxWidth int) {
	alloc := l.GetAllocation()
	maxWidth = l.GetMaxWidthChars()
	if maxWidth <= -1 {
		w, _ := l.GetSizeRequest()
		if w > -1 {
			maxWidth = w
		} else {
			maxWidth = alloc.W
		}
	}
	return
}

func (l *CEntry) refreshTextBuffer() (err error) {
	style := l.GetThemeRequest().Content.Normal
	alloc := l.GetAllocation()
	pos := l.GetPosition()

	l.Lock()

	posPoint := l.tProfile.GetPointFromPosition(pos)

	// keep pos within alloc
	if posPoint.X > alloc.W {
		l.offset.X = posPoint.X - alloc.W
		l.cursor.X = posPoint.X - l.offset.X
	} else {
		l.offset.X = 0
		l.cursor.X = posPoint.X
	}
	if posPoint.Y > alloc.H {
		l.offset.Y = posPoint.Y - alloc.H
		l.cursor.Y = posPoint.Y - l.offset.Y
	} else {
		l.offset.Y = 0
		l.cursor.Y = posPoint.Y
	}
	l.offset.W = alloc.W
	l.offset.H = alloc.H

	if l.cursor.X >= alloc.W {
		l.offset.X += 1
		l.cursor.X = alloc.W - 1
	}
	if l.cursor.Y >= alloc.H {
		l.offset.Y += 1
		l.cursor.Y = alloc.H - 1
	}
	// crop text to alloc using offset
	text := l.tProfile.Crop(*l.offset)
	// l.LogDebug("pos:%v, posPoint:%v, offset:%v, cursor:%v", pos, posPoint, l.offset, l.cursor)

	l.tBuffer = memphis.NewTextBuffer(text, style, false)
	l.Unlock()
	return
}

func (l *CEntry) resize(data []interface{}, argv ...interface{}) cenums.EventFlag {
	alloc := l.GetAllocation()
	if !l.IsVisible() || alloc.W <= 0 || alloc.H <= 0 {
		l.LogTrace("not visible, zero width or zero height")
		return cenums.EVENT_PASS
	}

	theme := l.GetThemeRequest()
	origin := l.GetOrigin()
	id := l.ObjectID()
	xPad, _ := l.GetPadding()
	_, yAlign := l.GetAlignment()

	size := ptypes.NewRectangle(alloc.W, alloc.H)
	local := ptypes.MakePoint2I(xPad, 0)
	size.W = alloc.W - (xPad * 2)
	size.H = alloc.H - (xPad * 2)

	if size.H < alloc.H {
		delta := alloc.H - size.H
		local.Y += int(float64(delta) * yAlign)
	}

	l.tRegion = ptypes.MakeRegion(local.X, local.Y, size.W, size.H)

	l.LockDraw()
	if err := memphis.ConfigureSurface(id, origin, alloc, theme.Content.Normal); err != nil {
		l.LogErr(err)
	}
	if err := memphis.ConfigureSurface(l.tid, local, *size, theme.Content.Normal); err != nil {
		l.LogErr(err)
	}
	l.UnlockDraw()

	l.Invalidate()
	return cenums.EVENT_STOP
}

func (l *CEntry) invalidate(data []interface{}, argv ...interface{}) cenums.EventFlag {
	theme := l.GetThemeRequest()
	if err := l.refreshTextBuffer(); err != nil {
		l.LogErr(err)
	}
	id := l.ObjectID()
	// region := l.GetRegion()
	origin := l.GetOrigin()
	alloc := l.GetAllocation()
	l.Lock()
	if !memphis.HasSurface(id) {
		if err := memphis.MakeSurface(id, origin, alloc, theme.Content.Normal); err != nil {
			l.LogErr(err)
		}
	}
	theme.Content.FillRune = rune(0)
	if err := memphis.FillSurface(id, theme); err != nil {
		l.LogErr(err)
	}
	if !memphis.HasSurface(l.tid) {
		if err := memphis.MakeSurface(l.tid, l.tRegion.Origin(), l.tRegion.Size(), theme.Content.Normal); err != nil {
			l.LogErr(err)
		}
	}
	if err := memphis.FillSurface(l.tid, theme); err != nil {
		l.LogErr(err)
	}
	l.Unlock()
	return cenums.EVENT_PASS
}

func (l *CEntry) draw(data []interface{}, argv ...interface{}) cenums.EventFlag {
	if surface, ok := argv[1].(*memphis.CSurface); ok {
		alloc := l.GetAllocation()
		if !l.IsVisible() || alloc.W <= 0 || alloc.H <= 0 {
			l.LogTrace("not visible, zero width or zero height")
			return cenums.EVENT_PASS
		}

		l.LockDraw()
		defer l.UnlockDraw()

		theme := l.GetThemeRequest()

		singleLineMode, lineWrapMode, justify, _ := l.Settings()

		if l.tBuffer != nil {
			if tSurface, err := memphis.GetSurface(l.tid); err != nil {
				l.LogErr(err)
			} else {
				tSurface.Box(
					ptypes.MakePoint2I(0, 0),
					ptypes.MakeRectangle(alloc.W, alloc.H),
					false, true,
					theme.Content.Overlay,
					theme.Content.FillRune,
					theme.Content.Normal,
					theme.Border.Normal,
					theme.Border.BorderRunes,
				)
				if f := l.tBuffer.Draw(tSurface, singleLineMode, lineWrapMode, false, justify, cenums.ALIGN_TOP); f == cenums.EVENT_STOP {
					if err := surface.CompositeSurface(tSurface); err != nil {
						l.LogErr(err)
					}
				}
			}
		}

		if debug, _ := l.GetBoolProperty(cdk.PropertyDebug); debug {
			surface.DebugBox(paint.ColorSilver, l.ObjectInfo())
		}
		return cenums.EVENT_STOP
	}
	return cenums.EVENT_PASS
}

func (l *CEntry) appendChange(name string, argv ...interface{}) {
	l.qLock.Lock()
	l.queue = append(l.queue, &cTextFieldChange{
		name: name,
		argv: argv,
	})
	l.qLock.Unlock()
}

func (l *CEntry) moveDown(lines int) {
	pos := l.GetPosition()
	l.Lock()
	posPoint := l.tProfile.GetPointFromPosition(pos)
	posPoint.Y += lines
	newPos := l.tProfile.GetPositionFromPoint(posPoint)
	l.appendChange("SetPosition", newPos)
	l.LogDebug("moved down (%v line) position: %v (%v)", lines, newPos, posPoint)
	l.Unlock()
}

func (l *CEntry) moveUp(lines int) {
	pos := l.GetPosition()
	l.Lock()
	posPoint := l.tProfile.GetPointFromPosition(pos)
	posPoint.Y -= lines
	if posPoint.Y < 0 {
		posPoint.Y = 0
	}
	newPos := l.tProfile.GetPositionFromPoint(posPoint)
	l.appendChange("SetPosition", newPos)
	l.LogDebug("moved up (%v line) position: %v (%v)", lines, newPos, posPoint)
	l.Unlock()
}

func (l *CEntry) moveHome() {
	pos := l.GetPosition()
	l.Lock()
	posPoint := l.tProfile.GetPointFromPosition(pos)
	posPoint.X = 0
	newPos := l.tProfile.GetPositionFromPoint(posPoint)
	l.appendChange("SetPosition", newPos)
	l.LogDebug("moved to home position: %v (%v)", newPos, posPoint)
	l.Unlock()
}

func (l *CEntry) moveEnd() {
	pos := l.GetPosition()
	l.Lock()
	posPoint := l.tProfile.GetPointFromPosition(pos)
	posPoint.X = -1
	newPos := l.tProfile.GetPositionFromPoint(posPoint)
	l.appendChange("SetPosition", newPos)
	l.LogDebug("moved to end position: %v (%v)", newPos, posPoint)
	l.Unlock()
}

func (l *CEntry) moveLeft(characters int) {
	if pos := l.GetPosition(); pos > 0 {
		l.LogDebug("move left %d character(s): %v", characters, pos-characters)
		l.appendChange("SetPosition", pos-characters)
	} else {
		l.LogDebug("at the start")
	}
}

func (l *CEntry) moveRight(characters int) {
	if pos := l.GetPosition(); pos < l.tProfile.Len() {
		l.LogDebug("move right %d character(s): %v", characters, pos+characters)
		l.appendChange("SetPosition", pos+characters)
	} else {
		l.LogDebug("all the way right: %v", pos)
	}
}

func (l *CEntry) deleteForwards() {
	pos := l.GetPosition()
	if tLen := l.tProfile.Len(); tLen > 0 {
		if pos < tLen {
			l.LogDebug("deleting forwards")
			l.appendChange("DeleteText", pos, pos)
			l.appendChange("SetPosition", pos)
		} else {
			l.LogDebug("deleting forwards (EOL)")
			l.appendChange("DeleteText", tLen-1, tLen-1)
			l.appendChange("SetPosition", tLen-1)
		}
	} else {
		l.LogDebug("nothing to delete forwards")
	}
}

func (l *CEntry) deleteBackwards() {
	pos := l.GetPosition()
	if pos > 0 {
		l.LogDebug("deleting backwards")
		l.appendChange("DeleteText", pos-1, pos-1)
		l.appendChange("SetPosition", pos-1)
	} else {
		l.LogDebug("nothing to delete backwards")
	}
}

func (l *CEntry) processQueue() (changesApplied bool) {
	l.qLock.Lock()
	if l.queue == nil || len(l.queue) == 0 {
		l.qLock.Unlock()
		return false
	}
	for _, change := range l.queue {
		switch change.name {
		case "SetPosition":
			if len(change.argv) == 1 {
				if v, ok := change.argv[0].(int); ok {
					l.setPosition(v)
					changesApplied = true
				} else {
					l.LogError("argument is not an 'int' for SetPosition change: %T (%v)", change.argv[0], change.argv)
				}
			} else {
				l.LogError("too many arguments for SetPosition change: %v", change.argv)
			}
		case "InsertText":
			if len(change.argv) == 2 {
				if newText, ok := change.argv[0].(string); ok {
					if pos, ok := change.argv[1].(int); ok {
						l.insertText(newText, pos)
						changesApplied = true
					} else {
						l.LogError("second argument is not an 'int' for InsertText change: %T (%v)", change.argv[1], change.argv)
					}
				} else {
					l.LogError("first argument is not a 'string' for InsertText change: %T (%v)", change.argv[0], change.argv)
				}
			} else {
				l.LogError("too many arguments for InsertText change: %v", change.argv)
			}
		case "DeleteText":
			if len(change.argv) == 2 {
				if start, ok := change.argv[0].(int); ok {
					if end, ok := change.argv[1].(int); ok {
						l.deleteText(start, end)
						changesApplied = true
					} else {
						l.LogError("second argument is not an 'int' for DeleteText change: %T (%v)", change.argv[1], change.argv)
					}
				} else {
					l.LogError("first argument is not an 'int' for DeleteText change: %T (%v)", change.argv[0], change.argv)
				}
			} else {
				l.LogError("too many arguments for DeleteText change: %v", change.argv)
			}
		}
	}
	l.queue = nil
	l.qLock.Unlock()
	return
}

func (l *CEntry) event(data []interface{}, argv ...interface{}) cenums.EventFlag {
	if evt, ok := argv[1].(cdk.Event); ok {
		switch e := evt.(type) {
		case *cdk.EventMouse:
			pos := ptypes.NewPoint2I(e.Position())
			switch e.State() {
			// case cdk.BUTTON_PRESS, cdk.DRAG_START:
			// 	if l.HasPoint(pos) && !l.HasEventFocus() {
			// 		l.GrabEventFocus()
			// 		return cenums.EVENT_STOP
			// 	}
			// case cdk.MOUSE_MOVE, cdk.DRAG_MOVE:
			// 	if l.HasEventFocus() {
			// 		if !l.HasPoint(pos) {
			// 			l.LogDebug("moved out of bounds")
			// 			l.CancelEvent()
			// 			return cenums.EVENT_STOP
			// 		}
			// 	}
			// 	return cenums.EVENT_PASS
			// case cdk.BUTTON_RELEASE, cdk.DRAG_STOP:
			// 	if l.HasEventFocus() {
			// 		if !l.HasPoint(pos) {
			// 			l.LogDebug("released out of bounds")
			// 			l.CancelEvent()
			// 			return cenums.EVENT_STOP
			// 		}
			// 		l.ReleaseEventFocus()
			// 		l.GrabFocus()
			// 		l.LogDebug("released")
			// 		return cenums.EVENT_STOP
			// 	}
			// }
			case cdk.BUTTON_RELEASE:
				if l.HasPoint(pos) {
					local := pos.NewClone()
					local.SubPoint(l.GetOrigin())
					local.AddPoint(l.offset.Origin())
					l.SetPosition(l.tProfile.GetPositionFromPoint(*local))
					if !l.HasFocus() {
						l.GrabFocus()
					}
					l.GetDisplay().RequestDraw()
					l.GetDisplay().RequestShow()
				}
			}
		case *cdk.EventKey:
			if !l.HasFocus() {
				return cenums.EVENT_PASS
			}
			r := e.Rune()
			v := e.Name()
			m := e.Modifiers()

			pos := l.GetPosition()

			switch r {

			case 10, 13:
				if l.GetSingleLineMode() {
					l.LogDebug("activate default")
				} else {
					l.appendChange("InsertText", "\n", pos)
					l.appendChange("SetPosition", pos+1)
					l.LogDebug(`printable key: \n, at pos: %v`, pos)
				}
				return cenums.EVENT_STOP

			case 127:
				l.deleteBackwards()
				return cenums.EVENT_STOP

			case 1: // 'a':
				if m.Has(cdk.ModCtrl) {
					// ctrl + a
					l.LogDebug("move home (ctrl+a)")
					l.moveHome()
					return cenums.EVENT_STOP
				}

			case 2: // 'b':
				if m.Has(cdk.ModCtrl) {
					// ctrl + b
					l.moveLeft(1)
					return cenums.EVENT_STOP
				}

			case 4: // 'd':
				if m.Has(cdk.ModCtrl) {
					// ctrl + d
					l.deleteForwards()
					return cenums.EVENT_STOP
				}

			case 5: // 'e':
				if m.Has(cdk.ModCtrl) {
					// ctrl + e
					l.LogDebug("move end (ctrl+e)")
					l.moveEnd()
					return cenums.EVENT_STOP
				}

			case 6: // 'f':
				if m.Has(cdk.ModCtrl) {
					// ctrl + f
					l.moveRight(1)
					return cenums.EVENT_STOP
				}

			case 8: // 'h':
				if m.Has(cdk.ModCtrl) {
					// ctrl + h
					l.deleteBackwards()
					return cenums.EVENT_STOP
				}

			case 14: // 'n':
				if m.Has(cdk.ModCtrl) {
					// ctrl + n
					l.moveDown(1)
					return cenums.EVENT_STOP
				}

			case 16: // 'p':
				if m.Has(cdk.ModCtrl) {
					// ctrl + p
					l.moveUp(1)
					return cenums.EVENT_STOP
				}
			}

			if k := e.Key(); k == cdk.KeyRune {
				pk := string(r)
				l.appendChange("InsertText", pk, pos)
				l.appendChange("SetPosition", pos+1)
				l.LogDebug("printable key: %v, at pos: %v", pk, pos)
				return cenums.EVENT_STOP
			}

			alloc := l.GetAllocation()

			switch v {
			case "Home":
				l.moveHome()
				l.LogDebug("move home (Home)")
				return cenums.EVENT_STOP

			case "End":
				l.moveEnd()
				l.LogDebug("move end (End)")
				return cenums.EVENT_STOP

			case "PgUp":
				l.LogDebug("move up %d lines (PgUp)", alloc.H)
				l.moveUp(alloc.H)
				return cenums.EVENT_STOP

			case "PgDn":
				l.LogDebug("move down %d lines (PgDn)", alloc.H)
				l.moveDown(alloc.H)
				return cenums.EVENT_STOP

			case "Delete":
				l.deleteForwards()
				return cenums.EVENT_STOP

			case "Left":
				l.moveLeft(1)
				return cenums.EVENT_STOP

			case "Right":
				l.moveRight(1)
				return cenums.EVENT_STOP

			case "Up", "Down":
				if l.GetSingleLineMode() {
					l.LogDebug("cannot move %v with single line mode", v)
				} else if v == "Down" {
					l.LogDebug("move down one line")
					l.moveDown(1)
				} else {
					l.LogDebug("move up one line")
					l.moveUp(1)
				}
				return cenums.EVENT_STOP

			default:
				l.LogDebug("other key: r:%v, n:%v", r, e.Name())
				return cenums.EVENT_STOP
			}
		}
	}
	return cenums.EVENT_PASS
}

func (l *CEntry) updateCursor() {
	if l.HasFocus() {
		if w := l.GetWindow(); w != nil {
			if d := w.GetDisplay(); d != nil {
				if s := d.Screen(); s != nil {
					o := l.GetOrigin()
					l.RLock()
					x, y := o.X+l.cursor.X, o.Y+l.cursor.Y
					l.RUnlock()
					if found := w.FindWidgetAt(ptypes.NewPoint2I(x, y)); found != nil {
						if l.ObjectID() == found.ObjectID() {
							s.ShowCursor(x, y)
							// l.LogDebug("cursor x,y = %v,%v (%v) [%v] - %v", x, y, l.cursor, l.offset, found.ObjectInfo())
						} else {
							s.HideCursor()
							// l.LogDebug("hide cursor x,y = %v,%v (%v) [%v] - %v", x, y, l.cursor, l.offset, found.ObjectInfo())
						}
					}
				}
			}
		}
		return
	}
	if w := l.GetWindow(); w != nil {
		if d := w.GetDisplay(); d != nil {
			if s := d.Screen(); s != nil {
				s.HideCursor()
			}
		}
	}
}

func (l *CEntry) lostFocus([]interface{}, ...interface{}) cenums.EventFlag {
	if l.qTimer != uuid.Nil {
		cdk.StopTimeout(l.qTimer)
		l.qTimer = uuid.Nil
	}
	l.UnsetState(enums.StateSelected)
	l.Invalidate()
	l.updateCursor()
	return cenums.EVENT_STOP
}

func (l *CEntry) gainedFocus([]interface{}, ...interface{}) cenums.EventFlag {
	l.SetState(enums.StateSelected)
	l.Invalidate()
	l.updateCursor()
	if l.qTimer != uuid.Nil {
		cdk.StopTimeout(l.qTimer)
		l.qTimer = uuid.Nil
	}
	l.qTimer = cdk.AddTimeout(time.Millisecond*100, func() cenums.EventFlag {
		if l.processQueue() {
			l.Invalidate()
			l.updateCursor()
		}
		return cenums.EVENT_PASS
	})
	return cenums.EVENT_STOP
}

const PropertyText cdk.Property = "text"

const PropertyEditable cdk.Property = "editable"

const TextFieldEventHandle = "text-field-event-handler"

const TextFieldLostFocusHandle = "text-field-lost-focus-handler"

const TextFieldGainedFocusHandle = "text-field-gained-focus-handler"

const TextFieldInvalidateHandle = "text-field-invalidate-handler"

const TextFieldResizeHandle = "text-field-resize-handler"

const TextFieldDrawHandle = "text-field-draw-handler"