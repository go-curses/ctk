package ctk

import (
	"fmt"
	"strings"

	"github.com/gofrs/uuid"

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
	DefaultEntryTheme = paint.Theme{
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

// Entry Hierarchy:
//	Object
//	  +- Widget
//	    +- Misc
//	      +- Entry
//	        +- AccelLabel
//	        +- TipsQuery
//
// The Entry Widget presents text to the end user.
type Entry interface {
	Misc
	Alignable
	Buildable
	Editable
	Sensitive

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
}

var _ Entry = (*CEntry)(nil)

type cEntryChange struct {
	name string
	argv []interface{}
}

// The CTextField structure implements the Entry interface and is exported
// to facilitate type embedding with custom implementations. No member variables
// are exported as the interface methods are the only intended means of
// interacting with Entry objects.
type CEntry struct {
	CMisc

	tid     uuid.UUID
	tRegion ptypes.Region

	offset    *ptypes.Region
	cursor    *ptypes.Point2I
	selection *ptypes.Range
	position  int

	tProfile *memphis.TextProfile
	tBuffer  memphis.TextBuffer
	tbStyle  paint.Style
}

// MakeEntry is used by the Buildable system to construct a new Entry.
func MakeEntry() Entry {
	return NewEntry("")
}

// NewEntry is the constructor for new Entry instances.
func NewEntry(plain string) Entry {
	l := new(CEntry)
	l.Init()
	l.SetText(plain)
	return l
}

// Init initializes a Entry object. This must be called at least once to
// set up the necessary defaults and allocate any memory structures. Calling
// this more than once is safe though unnecessary. Only the first call will
// result in any effect upon the Entry instance. Init is used in the
// NewEntry constructor and only necessary when implementing a derivative
// Entry type.
func (l *CEntry) Init() (already bool) {
	if l.InitTypeItem(TypeEntry, l) {
		return true
	}
	l.CMisc.Init()
	l.flags = enums.NULL_WIDGET_FLAG
	l.SetFlags(enums.SENSITIVE | enums.PARENT_SENSITIVE | enums.CAN_DEFAULT | enums.APP_PAINTABLE | enums.CAN_FOCUS)
	l.SetTheme(DefaultEntryTheme)

	l.selection = nil
	l.position = 0
	l.offset = ptypes.NewRegion(0, 0, 0, 0)
	l.cursor = ptypes.NewPoint2I(0, 0)
	l.tProfile = memphis.NewTextProfile("")
	l.tBuffer = nil
	l.tid, _ = uuid.NewV4()
	l.tRegion = ptypes.MakeRegion(0, 0, 0, 0)
	if err := memphis.MakeSurface(l.tid, l.tRegion.Origin(), l.tRegion.Size(), paint.DefaultColorStyle); err != nil {
		l.LogErr(err)
	}

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

	l.Connect(SignalCdkEvent, TextFieldEventHandle, l.event)
	l.Connect(SignalLostFocus, TextFieldLostFocusHandle, l.lostFocus)
	l.Connect(SignalGainedFocus, TextFieldGainedFocusHandle, l.gainedFocus)
	l.Connect(SignalResize, TextFieldResizeHandle, l.resize)
	l.Connect(SignalDraw, TextFieldDrawHandle, l.draw)
	// _ = l.SetBoolProperty(PropertyDebug, true)
	return false
}

// Build provides customizations to the Buildable system for Entry Widgets.
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

// SetText updates the text within the Entry widget. It overwrites any text that
// was there before. This will also clear any previously set mnemonic
// accelerators.
//
// Parameters:
// 	text	the text you want to set
//
// Locking: write
func (l *CEntry) SetText(text string) {
	l.setText(text)
}

func (l *CEntry) setText(text string) {
	l.Lock()
	l.tProfile.Set(text)
	l.Unlock()
	if err := l.SetStringProperty(PropertyText, l.tProfile.Get()); err != nil {
		l.LogErr(err)
	} else {
		l.refresh()
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
	} else {
		l.refresh()
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
	} else {
		l.refresh()
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
	} else {
		l.refresh()
	}
}

// SetLineWrap updates the line wrapping within the Entry widget. TRUE makes it
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
	} else {
		l.refresh()
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
	} else {
		l.refresh()
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
		l.Lock()
		l.selection = ptypes.NewRange(startOffset, endOffset)
		l.Unlock()
	}
}

// SetSelectable updates the selectable property for the Entry. TextFields allow the
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
	} else {
		l.refresh()
	}
}

// GetAttributes returns the attribute list that was set on the label using
// SetAttributes, if any. This function does not reflect attributes that come
// from the Entry markup (see SetMarkup).
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
		l.refresh()
	}
}

// Settings is a convenience method to return the interesting settings currently
// configured on the Entry instance.
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
	l.RLock()
	defer l.RUnlock()
	if l.selection != nil {
		startPos = l.selection.Start
		endPos = l.selection.End
		ok = true
	}
	return
}

func (l *CEntry) InsertTextAndSetPosition(newText string, index, position int) {
	l.insertTextAndSetPosition(newText, index, position)
}

func (l *CEntry) insertTextAndSetPosition(newText string, index, position int) {
	l.insertText(newText, index)
	l.setPosition(position)
}

func (l *CEntry) InsertText(newText string, position int) {
	l.insertText(newText, position)
}

func (l *CEntry) insertText(newText string, position int) {
	if modified, ok := l.tProfile.Insert(newText, position); ok {
		if err := l.SetStringProperty(PropertyText, modified); err != nil {
			l.LogErr(err)
		} else {
			l.refresh()
			l.Emit(SignalChangedText, l, modified)
		}
	}
}

func (l *CEntry) DeleteTextAndSetPosition(start, end, position int) {
	l.deleteTextAndSetPosition(start, end, position)
}

func (l *CEntry) deleteTextAndSetPosition(start, end, position int) {
	l.deleteText(start, end)
	l.setPosition(position)
}

func (l *CEntry) DeleteText(startPos int, endPos int) {
	l.deleteText(startPos, endPos)
}

func (l *CEntry) deleteText(startPos int, endPos int) {
	if modified, ok := l.tProfile.Delete(startPos, endPos); ok {
		if err := l.SetStringProperty(PropertyText, modified); err != nil {
			l.LogErr(err)
		} else {
			l.refresh()
			l.Emit(SignalChangedText, l, modified)
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
		value = content[startPos:]
	} else {
		value = content[startPos:endPos]
	}
	return
}

func (l *CEntry) CutClipboard() {
	value := ""
	l.RLock()
	if l.selection != nil && l.tProfile != nil {
		if l.tProfile.Len() > 0 {
			value = l.tProfile.Select(l.selection.Start, l.selection.End)
		}
		l.RUnlock()
		l.deleteTextAndSetPosition(l.selection.Start, l.selection.End, l.selection.Start)
	} else {
		l.RUnlock()
	}
	if d := l.GetDisplay(); d != nil {
		clipboard := d.GetClipboard()
		clipboard.Copy(value)
	}
	l.LogDebug("cut to clipboard: \"%v\"", value)
	l.clearSelection()
}

func (l *CEntry) CopyClipboard() {
	value := ""
	l.RLock()
	if l.selection != nil && l.tProfile != nil {
		if l.tProfile.Len() > 0 {
			value = l.tProfile.Select(l.selection.Start, l.selection.End)
		}
	}
	l.RUnlock()
	if d := l.GetDisplay(); d != nil {
		clipboard := d.GetClipboard()
		clipboard.Copy(value)
	}
	l.LogDebug("copied to clipboard: \"%v\"", value)
	l.clearSelection()
}

func (l *CEntry) PasteClipboard() {
	var value string
	if d := l.GetDisplay(); d != nil {
		clipboard := d.GetClipboard()
		value = clipboard.GetText()
	}
	pos := l.GetPosition()
	l.RLock()
	var selection *ptypes.Range
	if l.selection != nil {
		selection = l.selection.NewClone()
	}
	l.RUnlock()
	if selection != nil {
		l.deleteTextAndSetPosition(selection.Start, selection.End, selection.Start)
		pos = selection.Start
	}
	l.insertTextAndSetPosition(value, pos, pos+len(value))
	l.LogDebug("pasted from clipboard: \"%v\"", value)
	l.clearSelection()
}

func (l *CEntry) DeleteSelection() {
	l.RLock()
	selection := l.selection
	l.RUnlock()
	if selection != nil {
		l.deleteTextAndSetPosition(l.selection.Start, l.selection.End, l.selection.Start)
		l.LogDebug("selection deleted")
	}
	l.clearSelection()
}

func (l *CEntry) SetPosition(position int) {
	l.setPosition(position)
}

func (l *CEntry) setPosition(position int) {
	l.Lock()
	max := l.tProfile.Len()
	if position > max {
		position = max
	}
	l.position = position
	l.Unlock()
	l.refresh()
}

func (l *CEntry) GetPosition() (value int) {
	l.RLock()
	defer l.RUnlock()
	return l.position
}

func (l *CEntry) SetEditable(isEditable bool) {
	if err := l.SetBoolProperty(PropertyEditable, isEditable); err != nil {
		l.LogErr(err)
	} else {
		l.refresh()
	}
}

func (l *CEntry) GetEditable() (value bool) {
	var err error
	if value, err = l.GetBoolProperty(PropertyEditable); err != nil {
		l.LogErr(err)
	}
	return
}

// GetSizeRequest returns the requested size of the Entry taking into account
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

	if l.tBuffer != nil {
		l.tBuffer.Set(text, style)
	} else {
		l.tBuffer = memphis.NewTextBuffer(text, style, false)
	}

	l.Unlock()
	return
}

func (l *CEntry) refresh() {
	if err := l.refreshTextBuffer(); err != nil {
		l.LogErr(err)
	}
	l.updateCursor()
	l.Invalidate()
}

func (l *CEntry) resize(data []interface{}, argv ...interface{}) cenums.EventFlag {

	alloc := l.GetAllocation()
	if !l.IsVisible() || alloc.W <= 0 || alloc.H <= 0 {
		l.LogTrace("not visible, zero width or zero height")
		return cenums.EVENT_PASS
	}

	origin := l.GetOrigin()
	xPad, _ := l.GetPadding()
	_, yAlign := l.GetAlignment()

	size := ptypes.NewRectangle(alloc.W, alloc.H)
	local := ptypes.MakePoint2I(origin.X+xPad, origin.Y)
	size.W = alloc.W - (xPad * 2)
	size.H = alloc.H - (xPad * 2)

	if size.H < alloc.H {
		delta := alloc.H - size.H
		local.Y += int(float64(delta) * yAlign)
	}

	l.Lock()
	l.tRegion = ptypes.MakeRegion(local.X, local.Y, size.W, size.H)
	l.Unlock()

	theme := l.GetThemeRequest()
	if err := memphis.FillSurface(l.ObjectID(), theme); err != nil {
		l.LogErr(err)
	}
	if err := memphis.MakeConfigureSurface(l.tid, l.tRegion.Origin(), l.tRegion.Size(), theme.Content.Normal); err != nil {
		l.LogErr(err)
	} else if err := memphis.FillSurface(l.tid, theme); err != nil {
		l.LogErr(err)
	}

	l.refresh()
	return cenums.EVENT_STOP
}

func (l *CEntry) draw(data []interface{}, argv ...interface{}) cenums.EventFlag {

	if surface, ok := argv[1].(*memphis.CSurface); ok {
		alloc := l.GetAllocation()
		if !l.IsVisible() || alloc.W <= 0 || alloc.H <= 0 {
			l.LogTrace("not visible, zero width or zero height")
			return cenums.EVENT_PASS
		}

		theme := l.GetThemeRequest()
		singleLineMode, lineWrapMode, justify, _ := l.Settings()

		surface.Fill(theme)

		if tBuffer := l.tBuffer.Clone(); tBuffer != nil {
			tBuffer.SetStyle(theme.Content.Normal)
			if l.selection != nil {
				crop := l.tProfile.GetCropSelect(*l.selection, *l.offset)
				tBuffer.Select(crop.Start, crop.End)
			}

			if tSurface, err := memphis.GetSurface(l.tid); err != nil {
				l.LogErr(err)
			} else {
				tSurface.Fill(theme)
				tBuffer.Draw(tSurface, singleLineMode, lineWrapMode, false, justify, cenums.ALIGN_TOP)
				if err := surface.CompositeSurface(tSurface); err != nil {
					l.LogErr(err)
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

func (l *CEntry) updateSelection(oldPos, newPos int) (note string) {
	l.Lock()
	if l.tProfile != nil && l.tProfile.Len() > 0 {

		isMovingBackwards := oldPos > newPos

		if l.selection != nil {

			// moving selection start backwards
			// moving selection start forwards
			// moving selection end backwards
			// moving selection end forwards

			isMovingStart := newPos <= l.selection.Start
			isFlipping := l.selection.Start >= l.selection.End

			if isFlipping {
				if isMovingBackwards {
					l.selection.End = l.selection.Start
					l.selection.Start = newPos
					note = fmt.Sprintf("moving selection flip backwards: %v [%v,%v]", l.selection, oldPos, newPos)
				} else {
					l.selection.Start = l.selection.End
					l.selection.End = newPos
					note = fmt.Sprintf("moving selection flip forwards: %v [%v,%v]", l.selection, oldPos, newPos)
				}
			} else if isMovingStart {
				if isMovingBackwards {
					l.selection.Start = newPos + 1
					note = fmt.Sprintf("moving selection start backwards: %v [%v,%v]", l.selection, oldPos, newPos)
				} else {
					l.selection.Start = newPos + 1
					note = fmt.Sprintf("moving selection start forwards: %v [%v,%v]", l.selection, oldPos, newPos)
				}
			} else {
				if isMovingBackwards {
					l.selection.End = newPos - 1
					note = fmt.Sprintf("moving selection end backwards: %v [%v,%v]", l.selection, oldPos, newPos)
				} else {
					l.selection.End = newPos - 1
					note = fmt.Sprintf("moving selection end forwards: %v [%v,%v]", l.selection, oldPos, newPos)
				}
			}

		} else {

			if isMovingBackwards {
				l.selection = ptypes.NewRange(newPos+1, oldPos)
				note = fmt.Sprintf("started new selection backwards: %v [%v,%v]", l.selection, oldPos, newPos)
			} else {
				l.selection = ptypes.NewRange(oldPos, newPos-1)
				note = fmt.Sprintf("started new selection forwards: %v [%v,%v]", l.selection, oldPos, newPos)
			}

		}
	} else {
		note = fmt.Sprintf("cannot select range of zero-length string")
	}
	l.Unlock()
	l.Invalidate()
	return
}

func (l *CEntry) unselectAll() {
	if l.selectedAll() {
		l.clearSelection()
		l.setPosition(0)
	}
}

func (l *CEntry) selectAll() {
	l.Lock()
	end := l.tProfile.Len() - 1
	if l.selection == nil {
		l.selection = ptypes.NewRange(0, end)
		l.LogDebug("new select all (ctrl+a): %v", l.selection)
	} else {
		l.selection.Start = 0
		l.selection.End = end
		l.LogDebug("rpl select all (ctrl+a): %v", l.selection)
	}
	l.Unlock()
	l.setPosition(end + 1)
}

func (l *CEntry) selectedAll() bool {
	l.RLock()
	defer l.RUnlock()
	if l.selection == nil {
		return false
	}
	end := l.tProfile.Len() - 1
	return l.selection.Start == 0 && l.selection.End == end
}

func (l *CEntry) clearSelection() {
	l.Lock()
	if l.selection != nil {
		l.selection = nil
		l.Unlock()
		l.LogDebug("selection cleared")
		l.Invalidate()
	} else {
		l.Unlock()
	}
}

func (l *CEntry) moveSelection(oldPos, newPos int, shift bool) (note string) {
	note = "not selecting"
	if shift {
		note = l.updateSelection(oldPos, newPos)
	} else {
		l.clearSelection()
	}
	return
}

func (l *CEntry) moveDown(lines int, shift bool) {
	pos := l.GetPosition()
	l.RLock()
	posPoint := l.tProfile.GetPointFromPosition(pos)
	posPoint.Y += lines
	newPos := l.tProfile.GetPositionFromPoint(posPoint)
	l.RUnlock()
	note := l.moveSelection(pos, newPos, shift)
	l.setPosition(newPos)
	l.LogDebug("moved down (%v line) position: %v (%v) [%v]", lines, newPos, posPoint, note)
}

func (l *CEntry) moveUp(lines int, shift bool) {
	pos := l.GetPosition()
	l.RLock()
	posPoint := l.tProfile.GetPointFromPosition(pos)
	posPoint.Y -= lines
	if posPoint.Y < 0 {
		posPoint.Y = 0
	}
	newPos := l.tProfile.GetPositionFromPoint(posPoint)
	l.RUnlock()
	note := l.moveSelection(pos, newPos, shift)
	l.setPosition(newPos)
	l.LogDebug("moved up (%v line) position: %v (%v) [%v]", lines, newPos, posPoint, note)
}

func (l *CEntry) moveHome(shift bool) {
	pos := l.GetPosition()
	l.RLock()
	posPoint := l.tProfile.GetPointFromPosition(pos)
	posPoint.X = 0
	newPos := l.tProfile.GetPositionFromPoint(posPoint)
	l.RUnlock()
	note := l.moveSelection(pos, newPos-1, shift)
	l.setPosition(newPos)
	l.LogDebug("moved to home position: %v (%v) [%v]", newPos, posPoint, note)
}

func (l *CEntry) moveEnd(shift bool) {
	pos := l.GetPosition()
	l.RLock()
	posPoint := l.tProfile.GetPointFromPosition(pos)
	posPoint.X = -1
	newPos := l.tProfile.GetPositionFromPoint(posPoint)
	l.RUnlock()
	note := l.moveSelection(pos, newPos, shift)
	l.setPosition(newPos)
	l.LogDebug("moved to end position: %v (%v) [%v]", newPos, posPoint, note)
}

func (l *CEntry) moveLeft(characters int, shift bool) {
	if pos := l.GetPosition(); pos > 0 {
		newPos := pos - characters
		note := l.moveSelection(pos, newPos, shift)
		l.setPosition(newPos)
		l.LogDebug("move left %d character(s): %v [%v]", characters, newPos, note)
	} else {
		note := l.moveSelection(pos, -1, shift)
		l.LogDebug("at the start [%v]", note)
	}
}

func (l *CEntry) moveRight(characters int, shift bool) {
	if pos := l.GetPosition(); pos < l.tProfile.Len() {
		newPos := pos + characters
		note := l.moveSelection(pos, newPos, shift)
		l.setPosition(newPos)
		l.LogDebug("move right %d character(s): %v [%v]", characters, newPos, note)
	} else {
		l.LogDebug("all the way right: %v", pos)
	}
}

func (l *CEntry) deleteForwards() {
	pos := l.GetPosition()
	l.RLock()
	var selection *ptypes.Range
	if l.selection != nil {
		selection = l.selection.NewClone()
	}
	l.RUnlock()
	if tLen := l.tProfile.Len(); tLen > 0 {
		if selection != nil {
			l.LogDebug("deleting selection")
			l.deleteTextAndSetPosition(selection.Start, selection.End, selection.Start-1)
			l.clearSelection()
		} else {
			if pos < tLen {
				l.LogDebug("deleting forwards")
				l.deleteTextAndSetPosition(pos, pos, pos)
			} else {
				l.LogDebug("deleting forwards (EOL)")
				l.deleteTextAndSetPosition(tLen-1, tLen-1, tLen-1)
			}
		}
	} else {
		l.LogDebug("nothing to delete forwards")
	}
}

func (l *CEntry) deleteBackwards() {
	pos := l.GetPosition()
	l.RLock()
	var selection *ptypes.Range
	if l.selection != nil {
		selection = l.selection.NewClone()
	}
	l.RUnlock()
	if tLen := l.tProfile.Len(); tLen > 0 {
		if selection != nil {
			l.LogDebug("deleting selection")
			l.deleteTextAndSetPosition(selection.Start, selection.End, selection.Start)
			l.clearSelection()
		} else if pos > 0 {
			l.LogDebug("deleting backwards")
			l.deleteTextAndSetPosition(pos-1, pos-1, pos-1)
		} else {
			l.LogDebug("nothing to delete backwards")
		}
	} else {
		l.LogDebug("nothing to delete backwards")
	}
}

func (l *CEntry) event(data []interface{}, argv ...interface{}) cenums.EventFlag {
	l.LockEvent()
	defer l.UnlockEvent()

	if evt, ok := argv[1].(cdk.Event); ok {
		switch e := evt.(type) {
		case *cdk.EventPaste:
			if !l.HasFocus() {
				return cenums.EVENT_PASS
			}
			l.PasteClipboard()
			return cenums.EVENT_STOP

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
					l.insertTextAndSetPosition("\n", pos, pos+1)
					l.LogDebug(`printable key: \n, at pos: %v`, pos)
				}
				return cenums.EVENT_STOP

			case 127:
				l.deleteBackwards()
				return cenums.EVENT_STOP

			case 1: // 'a':
				if m.Has(cdk.ModCtrl) {
					// ctrl + a
					if l.selectedAll() {
						l.unselectAll()
					} else {
						l.selectAll()
					}
					return cenums.EVENT_STOP
				}

			case 2: // 'b':
				if m.Has(cdk.ModCtrl) {
					// ctrl + b
					l.moveLeft(1, m.Has(cdk.ModShift))
					return cenums.EVENT_STOP
				}

			case 3: // 'c':
				if m.Has(cdk.ModCtrl) {
					// ctrl + c
					l.CopyClipboard()
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
					l.moveEnd(m.Has(cdk.ModShift))
					return cenums.EVENT_STOP
				}

			case 6: // 'f':
				if m.Has(cdk.ModCtrl) {
					// ctrl + f
					l.moveRight(1, m.Has(cdk.ModShift))
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
					l.moveDown(1, m.Has(cdk.ModShift))
					return cenums.EVENT_STOP
				}

			case 16: // 'p':
				if m.Has(cdk.ModCtrl) {
					// ctrl + p
					l.moveUp(1, m.Has(cdk.ModShift))
					return cenums.EVENT_STOP
				}

			case 22: // 'v':
				if m.Has(cdk.ModCtrl) {
					// ctrl + v
					l.PasteClipboard()
					return cenums.EVENT_STOP
				}

			case 24: // 'x':
				if m.Has(cdk.ModCtrl) {
					// ctrl + x
					l.CutClipboard()
					return cenums.EVENT_STOP
				}
			}

			if k := e.Key(); k == cdk.KeyRune {
				pk := string(r)
				l.RLock()
				var selection *ptypes.Range
				if l.selection != nil {
					selection = l.selection.NewClone()
				}
				l.RUnlock()
				if selection != nil {
					pos = l.selection.Start
					l.deleteText(l.selection.Start, l.selection.End)
					l.clearSelection()
					l.LogDebug("replacing selection with printable key...")
				}
				l.insertTextAndSetPosition(pk, pos, pos+1)
				l.LogDebug("printable key: %v, at pos: %v", pk, pos)
				return cenums.EVENT_STOP
			}

			alloc := l.GetAllocation()

			switch v {
			case "Home", "Shift+Home":
				l.moveHome(m.Has(cdk.ModShift))
				l.LogDebug("move home (Home)")
				return cenums.EVENT_STOP

			case "End", "Shift+End":
				l.moveEnd(m.Has(cdk.ModShift))
				l.LogDebug("move end (End)")
				return cenums.EVENT_STOP

			case "PgUp", "Shift+PgUp":
				l.LogDebug("move up %d lines (PgUp)", alloc.H)
				l.moveUp(alloc.H, m.Has(cdk.ModShift))
				return cenums.EVENT_STOP

			case "PgDn", "Shift+PgDn":
				l.LogDebug("move down %d lines (PgDn)", alloc.H)
				l.moveDown(alloc.H, m.Has(cdk.ModShift))
				return cenums.EVENT_STOP

			case "Delete", "Shift+Delete":
				l.deleteForwards()
				return cenums.EVENT_STOP

			case "Left", "Shift+Left":
				l.moveLeft(1, m.Has(cdk.ModShift))
				return cenums.EVENT_STOP

			case "Right", "Shift+Right":
				l.moveRight(1, m.Has(cdk.ModShift))
				return cenums.EVENT_STOP

			case "Up", "Down", "Shift+Up", "Shift+Down":
				vv := strings.Replace(v, "Shift+", "", 1)
				if l.GetSingleLineMode() {
					l.LogDebug("cannot move %v with single line mode", v)
				} else if vv == "Down" {
					l.LogDebug("move down one line")
					l.moveDown(1, m.Has(cdk.ModShift))
				} else {
					l.LogDebug("move up one line")
					l.moveUp(1, m.Has(cdk.ModShift))
				}
				return cenums.EVENT_STOP

			default:
				l.LogDebug("other key: r:%v, n:%v", r, e.Name())
				return cenums.EVENT_STOP
			}

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
					if !l.HasFocus() {
						l.GrabFocus()
					}
					local := pos.NewClone()
					local.SubPoint(l.GetOrigin())
					local.AddPoint(l.offset.Origin())
					l.clearSelection()
					l.SetPosition(l.tProfile.GetPositionFromPoint(*local))
					return cenums.EVENT_STOP
				}
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
	l.UnsetState(enums.StateSelected)
	l.refresh()
	return cenums.EVENT_PASS
}

func (l *CEntry) gainedFocus([]interface{}, ...interface{}) cenums.EventFlag {
	l.SetState(enums.StateSelected)
	l.refresh()
	return cenums.EVENT_PASS
}

const PropertyText cdk.Property = "text"

const PropertyEditable cdk.Property = "editable"

const TextFieldEventHandle = "text-field-event-handler"

const TextFieldLostFocusHandle = "text-field-lost-focus-handler"

const TextFieldGainedFocusHandle = "text-field-gained-focus-handler"

const TextFieldInvalidateHandle = "text-field-invalidate-handler"

const TextFieldResizeHandle = "text-field-resize-handler"

const TextFieldDrawHandle = "text-field-draw-handler"