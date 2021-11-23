package ctk

import (
	"regexp"
	"strings"

	"github.com/go-curses/cdk"
	"github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/cdk/lib/paint"
	"github.com/go-curses/cdk/lib/ptypes"
	cstrings "github.com/go-curses/cdk/lib/strings"
	"github.com/go-curses/cdk/memphis"
	"github.com/gofrs/uuid"
)

const TypeLabel cdk.CTypeTag = "ctk-label"

var (
	rxLabelPlainText = regexp.MustCompile(`(?msi)(_)([A-Za-z0-9])`)
	rxLabelMnemonic  = regexp.MustCompile(`(?msi)_([A-Za-z0-9])`)
)

func init() {
	_ = cdk.TypesManager.AddType(TypeLabel, func() interface{} { return MakeLabel() })
	ctkBuilderTranslators[TypeLabel] = func(builder Builder, widget Widget, name, value string) error {
		switch strings.ToLower(name) {
		case "wrap":
			isTrue := cstrings.IsTrue(value)
			if err := widget.SetBoolProperty(PropertyWrap, isTrue); err != nil {
				return err
			}
			if isTrue {
				if wmi, err := widget.GetStructProperty(PropertyWrapMode); err == nil {
					if wm, ok := wmi.(enums.WrapMode); ok {
						if wm == enums.WRAP_NONE {
							if err := widget.SetStructProperty(PropertyWrapMode, enums.WRAP_WORD); err != nil {
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

// Label Hierarchy:
//	Object
//	  +- Widget
//	    +- Misc
//	      +- Label
//	        +- AccelLabel
//	        +- TipsQuery
//
// The Label Widget presents text to the end user.
type Label interface {
	Misc
	Alignable
	Buildable

	Init() (already bool)
	Build(builder Builder, element *CBuilderElement) error
	SetText(text string)
	SetAttributes(attrs paint.Style)
	SetMarkup(text string) (parseError error)
	SetMarkupWithMnemonic(str string) (err error)
	SetJustify(justify enums.Justification)
	SetEllipsize(mode bool)
	SetWidthChars(nChars int)
	SetMaxWidthChars(nChars int)
	SetLineWrap(wrap bool)
	SetLineWrapMode(wrapMode enums.WrapMode)
	GetMnemonicKeyVal() (value rune)
	GetSelectable() (value bool)
	GetText() (value string)
	SelectRegion(startOffset int, endOffset int)
	SetMnemonicWidget(widget Widget)
	SetSelectable(setting bool)
	SetTextWithMnemonic(str string)
	GetAttributes() (value paint.Style)
	GetJustify() (value enums.Justification)
	GetEllipsize() (value bool)
	GetWidthChars() (value int)
	GetMaxWidthChars() (value int)
	GetLabel() (value string)
	GetLineWrap() (value bool)
	GetLineWrapMode() (value enums.WrapMode)
	GetMnemonicWidget() (value Widget)
	GetSelectionBounds() (start int, end int, nonEmpty bool)
	GetUseMarkup() (value bool)
	GetUseUnderline() (value bool)
	GetSingleLineMode() (value bool)
	SetLabel(str string)
	SetUseMarkup(setting bool)
	SetUseUnderline(setting bool)
	SetSingleLineMode(singleLineMode bool)
	GetCurrentUri() (value string)
	SetTrackVisitedLinks(trackLinks bool)
	GetTrackVisitedLinks() (value bool)
	GetClearText() (text string)
	GetPlainText() (text string)
	GetCleanText() (text string)
	GetPlainTextInfo() (maxWidth, lineCount int)
	GetPlainTextInfoAtWidth(width int) (maxWidth, lineCount int)
	GetSizeRequest() (width, height int)
}

// The CLabel structure implements the Label interface and is exported
// to facilitate type embedding with custom implementations. No member variables
// are exported as the interface methods are the only intended means of
// interacting with Label objects.
type CLabel struct {
	CMisc

	text string

	tid     uuid.UUID
	tbuffer memphis.TextBuffer
	tbStyle paint.Style
}

// MakeLabel is used by the Buildable system to construct a new Label.
func MakeLabel() *CLabel {
	return NewLabel("")
}

// NewLabel is the constructor for new Label instances.
func NewLabel(plain string) *CLabel {
	l := new(CLabel)
	l.Init()
	l.SetText(plain)
	return l
}

// NewLabelWithMnemonic creates a new Label, containing the text given. If
// characters in the string are preceded by an underscore, they are underlined.
// If you need a literal underscore character in a label, use '__' (two
// underscores). The first underlined character represents a keyboard
// accelerator called a mnemonic. The mnemonic key can be used to activate
// another widget, chosen automatically, or explicitly using SetMnemonicWidget.
// If SetMnemonicWidget is not called, then the first activatable ancestor of
// the Label will be chosen as the mnemonic widget. For instance, if the label
// is inside a button or menu item, the button or menu item will automatically
// become the mnemonic widget and be activated by the mnemonic.
//
// Parameters:
// 	label	text, with an underscore in front of the mnemonic character
func NewLabelWithMnemonic(label string) (value *CLabel) {
	l := new(CLabel)
	l.Init()
	l.SetTextWithMnemonic(label)
	return l
}

// NewLabelWithMarkup creates a new Label, containing the text given and if the
// text contains Tango markup, the rendered text will display accordingly.
func NewLabelWithMarkup(markup string) (label *CLabel, err error) {
	label = new(CLabel)
	label.Init()
	err = label.SetMarkup(markup)
	return
}

// Init initializes a Label object. This must be called at least once to
// set up the necessary defaults and allocate any memory structures. Calling
// this more than once is safe though unnecessary. Only the first call will
// result in any effect upon the Label instance. Init is used in the
// NewLabel constructor and only necessary when implementing a derivative
// Label type.
func (l *CLabel) Init() (already bool) {
	if l.InitTypeItem(TypeLabel, l) {
		return true
	}
	l.CMisc.Init()
	l.flags = NULL_WIDGET_FLAG
	l.SetFlags(PARENT_SENSITIVE | APP_PAINTABLE)
	_ = l.InstallProperty(PropertyAttributes, cdk.StructProperty, true, nil)
	_ = l.InstallProperty(PropertyCursorPosition, cdk.IntProperty, false, 0)
	_ = l.InstallProperty(PropertyEllipsize, cdk.BoolProperty, true, false)
	_ = l.InstallProperty(PropertyJustify, cdk.StructProperty, true, enums.JUSTIFY_LEFT)
	_ = l.InstallProperty(PropertyLabel, cdk.StringProperty, true, "")
	_ = l.InstallProperty(PropertyMaxWidthChars, cdk.IntProperty, true, -1)
	_ = l.InstallProperty(PropertyMnemonicKeyVal, cdk.IntProperty, false, rune(0))
	_ = l.InstallProperty(PropertyMnemonicWidget, cdk.StructProperty, true, nil)
	_ = l.InstallProperty(PropertySelectable, cdk.BoolProperty, true, false)
	_ = l.InstallProperty(PropertySelectionBound, cdk.IntProperty, false, 0)
	_ = l.InstallProperty(PropertySingleLineMode, cdk.BoolProperty, true, false)
	_ = l.InstallProperty(PropertyTrackVisitedLinks, cdk.BoolProperty, true, true)
	_ = l.InstallProperty(PropertyUseMarkup, cdk.BoolProperty, true, false)
	_ = l.InstallProperty(PropertyUseUnderline, cdk.BoolProperty, true, false)
	_ = l.InstallProperty(PropertyWidthChars, cdk.IntProperty, true, -1)
	_ = l.InstallProperty(PropertyWrap, cdk.BoolProperty, true, false)
	_ = l.InstallProperty(PropertyWrapMode, cdk.StructProperty, true, enums.WRAP_WORD)
	l.Connect(SignalInvalidate, LabelInvalidateHandle, l.invalidate)
	l.Connect(SignalResize, LabelResizeHandle, l.resize)
	l.Connect(SignalDraw, LabelDrawHandle, l.draw)
	l.text = ""
	l.tbuffer = nil
	l.tid, _ = uuid.NewV4()
	if err := memphis.RegisterSurface(l.tid, ptypes.Point2I{}, ptypes.Rectangle{}, paint.DefaultColorStyle); err != nil {
		l.LogErr(err)
	}
	// _ = l.SetBoolProperty(PropertyDebug, true)
	l.Invalidate()
	return false
}

// Build provides customizations to the Buildable system for Label Widgets.
func (l *CLabel) Build(builder Builder, element *CBuilderElement) error {
	l.Freeze()
	defer l.Thaw()
	if name, ok := element.Attributes["id"]; ok {
		l.SetName(name)
	}
	for k, v := range element.Properties {
		switch cdk.Property(k) {
		case PropertyLabel:
			l.SetLabel(v)
		default:
			element.ApplyProperty(k, v)
		}
	}
	element.ApplySignals()
	return nil
}

// SetText updates the text within the Label widget. It overwrites any text that
// was there before. This will also clear any previously set mnemonic
// accelerators.
//
// Parameters:
// 	text	the text you want to set
//
// Locking: write
func (l *CLabel) SetText(text string) {
	if err := l.SetBoolProperty(PropertyUseMarkup, false); err != nil {
		l.LogErr(err)
	}
	l.Lock()
	l.text = text
	l.Unlock()
	l.Invalidate()
}

// SetAttributes updates the attributes property to be the given paint.Style.
//
// Parameters:
// 	attrs	a paint.Style
//
// Locking: write
func (l *CLabel) SetAttributes(attrs paint.Style) {
	if err := l.SetStructProperty(PropertyAttributes, attrs); err != nil {
		l.LogErr(err)
	}
}

// SetMarkup parses text which is marked up with the Tango text markup language,
// setting the Label's text and attribute list based on the parse results.
//
// Parameters:
// 	text	a markup string (see Tango markup format)
//
// Locking: write
func (l *CLabel) SetMarkup(text string) (parseError error) {
	l.Lock()
	l.text = text
	l.Unlock()
	if err := l.SetBoolProperty(PropertyUseMarkup, true); err != nil {
		l.LogErr(err)
	}
	// Invalidate will call refreshTextBuffer again, we do this once before to
	// see if there's any errors generated because we can't return any errors
	// encountered in signal handlers (beyond the logging).
	if parseError = l.refreshTextBuffer(); parseError == nil {
		l.Invalidate()
	}
	return
}

// SetMarkupWithMnemonic parses str which is marked up with the Tango text
// markup language, setting the label's text and attribute list based on the
// parse results. If characters in str are preceded by an underscore, they are
// underlined indicating that they represent a keyboard accelerator called a
// mnemonic. The mnemonic key can be used to activate another widget, chosen
// automatically, or explicitly using SetMnemonicWidget.
//
// Parameters:
// 	str	a markup string (see Pango markup format)
//
// Locking: write
func (l *CLabel) SetMarkupWithMnemonic(str string) (err error) {
	l.SetUseUnderline(true)
	err = l.SetMarkup(str)
	return
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
func (l *CLabel) SetJustify(justify enums.Justification) {
	if err := l.SetStructProperty(PropertyJustify, justify); err != nil {
		l.LogErr(err)
	}
}

// SetEllipsize updates the mode used to ellipsize (add an ellipsis: "...") to
// the text if there is not enough space to render the entire string.
//
// Parameters:
// 	mode	bool
//
// Locking: write
func (l *CLabel) SetEllipsize(mode bool) {
	if err := l.SetBoolProperty(PropertyEllipsize, mode); err != nil {
		l.LogErr(err)
	}
}

// SetWidthChars updates the desired width in characters of label to nChars.
//
// Parameters:
// 	nChars	the new desired width, in characters.
//
// Locking: write
func (l *CLabel) SetWidthChars(nChars int) {
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
func (l *CLabel) SetMaxWidthChars(nChars int) {
	if err := l.SetIntProperty(PropertyMaxWidthChars, nChars); err != nil {
		l.LogErr(err)
	}
}

// SetLineWrap updates the line wrapping within the Label widget. TRUE makes it
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
func (l *CLabel) SetLineWrap(wrap bool) {
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
func (l *CLabel) SetLineWrapMode(wrapMode enums.WrapMode) {
	if err := l.SetStructProperty(PropertyWrapMode, wrapMode); err != nil {
		l.LogErr(err)
	}
}

// GetMnemonicKeyVal returns the mnemonic character in the Label text if the
// Label has been set so that it has an mnemonic key. Rhis function
// returns the keyval used for the mnemonic accelerator. If there is no
// mnemonic set up it returns `rune(0)`.
//
// Returns:
//  value	keyval usable for accelerators
//
// Locking: read
func (l *CLabel) GetMnemonicKeyVal() (value rune) {
	if l.GetUseUnderline() {
		label := l.GetClearText()
		l.RLock()
		if rxLabelMnemonic.MatchString(label) {
			m := rxLabelMnemonic.FindStringSubmatch(label)
			if len(m) > 1 {
				l.RUnlock()
				return rune(strings.ToLower(m[1])[0])
			}
		}
		l.RUnlock()
	}
	return
}

// GetSelectable returns the value set by SetSelectable.
//
// Locking: read
func (l *CLabel) GetSelectable() (value bool) {
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
func (l *CLabel) GetText() (value string) {
	value = l.GetCleanText()
	return
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
// Note that usage of this within CTK is unimplemented at this time
//
// Method stub, unimplemented
func (l *CLabel) SelectRegion(startOffset int, endOffset int) {}

// SetMnemonicWidget updates the mnemonic-widget property with the given Widget.
// If the label has been set so that it has a mnemonic key (using i.e.
// SetMarkupWithMnemonic, SetTextWithMnemonic,
// NewWithMnemonic or the use-underline property) the label can be associated
// with a widget that is the target of the mnemonic. When the label is inside a
// widget (like a Button or a Notebook tab) it is automatically associated with
// the correct widget, but sometimes (i.e. when the target is a Entry next to
// the label) you need to set it explicitly using this function. The target
// widget will be accelerated by emitting the mnemonic-activate signal on it.
// The default handler for this signal will activate the widget if there are no
// mnemonic collisions and toggle focus between the colliding widgets otherwise.
//
// Parameters:
// 	widget	the target Widget.
//
// Note that usage of this within CTK is unimplemented at this time
//
// Locking: write
func (l *CLabel) SetMnemonicWidget(widget Widget) {
	if err := l.SetStructProperty(PropertyMnemonicWidget, widget); err != nil {
		l.LogErr(err)
	} else {
		l.Invalidate()
	}
}

// SetSelectable updates the selectable property for the Label. Labels allow the
// user to select text from the label, for copy-and-paste.
//
// Parameters:
// 	setting	TRUE to allow selecting text in the label
//
// Note that usage of this within CTK is unimplemented at this time
//
// Locking: write
func (l *CLabel) SetSelectable(setting bool) {
	if err := l.SetBoolProperty(PropertySelectable, setting); err != nil {
		l.LogErr(err)
	}
}

// SetTextWithMnemonic updates the Label's text from the string str. If
// characters in str are preceded by an underscore, they are underlined
// indicating that they represent a keyboard accelerator called a mnemonic. The
// mnemonic key can be used to activate another widget, chosen automatically, or
// explicitly using SetMnemonicWidget.
//
// Parameters:
// 	str	a string
//
// Locking: write
func (l *CLabel) SetTextWithMnemonic(str string) {
	l.SetUseUnderline(true)
	l.SetText(str)
}

// GetAttributes returns the attribute list that was set on the label using
// SetAttributes, if any. This function does not reflect attributes that come
// from the Label markup (see SetMarkup).
//
// Locking: read
func (l *CLabel) GetAttributes() (value paint.Style) {
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
func (l *CLabel) GetJustify() (value enums.Justification) {
	var ok bool
	if v, err := l.GetStructProperty(PropertyJustify); err != nil {
		l.LogErr(err)
	} else if value, ok = v.(enums.Justification); !ok {
		l.LogError("value stored in PropertyJustify is not of enums.Justification type: %v (%T)", v, v)
	}
	return
}

// GetEllipsize returns the ellipsizing state of the label.
// See: SetEllipsize()
//
// Locking: read
func (l *CLabel) GetEllipsize() (value bool) {
	var err error
	if value, err = l.GetBoolProperty(PropertyEllipsize); err != nil {
		l.LogErr(err)
	}
	return
}

// GetWidthChars retrieves the desired width of label, in characters.
// See: SetWidthChars()
//
// Locking: read
func (l *CLabel) GetWidthChars() (value int) {
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
func (l *CLabel) GetMaxWidthChars() (value int) {
	var err error
	if value, err = l.GetIntProperty(PropertyMaxWidthChars); err != nil {
		l.LogErr(err)
	}
	return
}

// GetLabel returns the text from a label widget including any embedded
// underlines indicating mnemonics and Tango markup.
// See: GetText()
//
// Locking: read
func (l *CLabel) GetLabel() (value string) {
	var err error
	if value, err = l.GetStringProperty(PropertyLabel); err != nil {
		l.LogErr(err)
	}
	return
}

// GetLineWrap returns whether lines in the label are automatically wrapped.
// See: SetLineWrap()
//
// Locking: read
func (l *CLabel) GetLineWrap() (value bool) {
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
func (l *CLabel) GetLineWrapMode() (value enums.WrapMode) {
	var ok bool
	if v, err := l.GetStructProperty(PropertyWrapMode); err != nil {
		l.LogErr(err)
	} else if value, ok = v.(enums.WrapMode); !ok {
		l.LogError("value stored in PropertyWrap is not of enums.WrapMode type: %v (%T)", v, v)
	}
	return
}

// GetMnemonicWidget retrieves the target of the mnemonic (keyboard shortcut) of
// this Label.
// See: SetMnemonicWidget()
//
// Locking: read
func (l *CLabel) GetMnemonicWidget() (value Widget) {
	if v, err := l.GetStructProperty(PropertyMnemonicWidget); err == nil {
		value, _ = v.(Widget)
	} else {
		l.LogErr(err)
	}
	return
}

// GetSelectionBounds returns the selected range of characters in the label,
// returning TRUE for nonEmpty if there's a selection.
//
// Note that usage of this within CTK is unimplemented at this time
//
// Method stub, unimplemented
func (l *CLabel) GetSelectionBounds() (start int, end int, nonEmpty bool) {
	return 0, 0, false
}

// GetUseMarkup returns whether the label's text is interpreted as marked up
// with the Tango text markup language.
// See: SetUseMarkup()
//
// Locking: read
func (l *CLabel) GetUseMarkup() (value bool) {
	var err error
	if value, err = l.GetBoolProperty(PropertyUseMarkup); err != nil {
		l.LogErr(err)
	}
	return
}

// GetUseUnderline returns whether an embedded underline in the label indicates
// a mnemonic.
// See: SetUseUnderline()
//
// Locking: read
func (l *CLabel) GetUseUnderline() (value bool) {
	var err error
	if value, err = l.GetBoolProperty(PropertyUseUnderline); err != nil {
		l.LogErr(err)
	}
	return
}

// GetSingleLineMode returns whether the label is in single line mode.
//
// Locking: read
func (l *CLabel) GetSingleLineMode() (value bool) {
	var err error
	if value, err = l.GetBoolProperty(PropertySingleLineMode); err != nil {
		l.LogErr(err)
	}
	return
}

// SetLabel updates the text of the label. The label is interpreted as including
// embedded underlines and/or Pango markup depending on the values of the
// use-underline and use-markup properties.
//
// Parameters:
// 	str	the new text to set for the label
//
// Locking: write
func (l *CLabel) SetLabel(str string) {
	if err := l.SetStringProperty(PropertyLabel, str); err != nil {
		l.LogErr(err)
	} else {
		if l.GetUseMarkup() {
			if err := l.SetMarkup(str); err != nil {
				l.LogErr(err)
			}
		} else {
			l.SetText(str)
		}
	}
}

// SetUseMarkup updates whether the text of the label contains markup in Tango's
// text markup language.
// See: SetMarkup()
//
// Parameters:
// 	setting	TRUE if the label's text should be parsed for markup.
//
// Locking: write
func (l *CLabel) SetUseMarkup(setting bool) {
	if err := l.SetBoolProperty(PropertyUseMarkup, setting); err != nil {
		l.LogErr(err)
	} else {
		l.Invalidate()
	}
}

// SetUseUnderline updates the use-underline property for the Lable. If TRUE, an
// underline in the text indicates the next character should be used for the
// mnemonic accelerator key.
//
// Parameters:
// 	setting	TRUE if underlines in the text indicate mnemonics
//
// Locking: write
func (l *CLabel) SetUseUnderline(setting bool) {
	if err := l.SetBoolProperty(PropertyUseUnderline, setting); err != nil {
		l.LogErr(err)
	} else {
		l.Invalidate()
	}
}

// SetSingleLineMode updates whether the label is in single line mode.
//
// Parameters:
// 	singleLineMode	TRUE if the label should be in single line mode
//
// Locking: write
func (l *CLabel) SetSingleLineMode(singleLineMode bool) {
	if err := l.SetBoolProperty(PropertySingleLineMode, singleLineMode); err != nil {
		l.LogErr(err)
	} else {
		l.Invalidate()
	}
}

// GetCurrentUri returns the URI for the currently active link in the label. The
// active link is the one under the mouse pointer or, in a selectable label, the
// link in which the text cursor is currently positioned. This function is
// intended for use in a activate-link handler or for use in a
// query-tooltip handler.
//
// Note that usage of this within CTK is unimplemented at this time
//
// Method stub, unimplemented
func (l *CLabel) GetCurrentUri() (value string) {
	return ""
}

// SetTrackVisitedLinks updates whether the label should keep track of clicked
// links (and use a different color for them).
//
// Parameters:
// 	trackLinks	TRUE to track visited links
//
// Note that usage of this within CTK is unimplemented at this time
//
// Locking: write
func (l *CLabel) SetTrackVisitedLinks(trackLinks bool) {
	if err := l.SetBoolProperty(PropertyTrackVisitedLinks, trackLinks); err != nil {
		l.LogErr(err)
	}
}

// GetTrackVisitedLinks returns whether the label is currently keeping track of
// clicked links.
//
// Returns:
// 	TRUE if clicked links are remembered
//
// Note that usage of this within CTK is unimplemented at this time
//
// Locking: read
func (l *CLabel) GetTrackVisitedLinks() (value bool) {
	var err error
	if value, err = l.GetBoolProperty(PropertyTrackVisitedLinks); err != nil {
		l.LogErr(err)
	}
	return
}

// Settings is a convenience method to return the interesting settings currently
// configured on the Label instance.
//
// Locking: read
func (l *CLabel) Settings() (singleLineMode bool, lineWrapMode enums.WrapMode, ellipsize bool, justify enums.Justification, maxWidthChars int) {
	singleLineMode = l.GetSingleLineMode()
	lineWrapMode = l.GetLineWrapMode()
	ellipsize = l.GetEllipsize()
	justify = l.GetJustify()
	maxWidthChars = l.GetMaxWidthChars()
	return
}

// GetClearText returns the Label's text, stripped of markup.
//
// Locking: read
func (l *CLabel) GetClearText() (text string) {
	if l.tbuffer == nil {
		return ""
	}
	singleLineMode, lineWrapMode, ellipsize, justify, maxWidthChars := l.Settings()
	l.RLock()
	text = l.tbuffer.ClearText(lineWrapMode, ellipsize, justify, maxWidthChars)
	if singleLineMode {
		if strings.Contains(text, "\n") {
			if idx := strings.Index(text, "\n"); idx >= 0 {
				text = text[:idx]
			}
		}
	}
	l.RUnlock()
	return
}

// GetPlainText returns the Label's text, stripped of markup and mnemonics.
//
// Locking: read
func (l *CLabel) GetPlainText() (text string) {
	if l.tbuffer == nil {
		return ""
	}
	singleLineMode, lineWrapMode, ellipsize, justify, maxWidthChars := l.Settings()
	l.RLock()
	text = l.tbuffer.PlainText(lineWrapMode, ellipsize, justify, maxWidthChars)
	if singleLineMode {
		if strings.Contains(text, "\n") {
			if idx := strings.Index(text, "\n"); idx >= 0 {
				text = text[:idx]
			}
		}
	}
	l.RUnlock()
	return
}

// GetCleanText filters the result of GetClearText to strip leading underscores.
//
// Locking: read
func (l *CLabel) GetCleanText() (text string) {
	plain := l.GetPlainText()
	l.RLock()
	text = rxLabelPlainText.ReplaceAllString(plain, "$2")
	l.RUnlock()
	return
}

// GetPlainTextInfo returns the maximum line width and line count.
//
// Locking: read
func (l *CLabel) GetPlainTextInfo() (maxWidth, lineCount int) {
	maxWidthChars := l.GetMaxWidthChars()
	maxWidth, lineCount = l.GetPlainTextInfoAtWidth(maxWidthChars)
	return
}

// GetPlainTextInfoAtWidth returns the maximum line width and line count, with the given width as an override to the
// value set with SetMaxWidthChars. This is used primarily for pre-rendering stages like Resize to determine the size
// allocations without having to render the actual text with Tango first.
//
// Locking: read
func (l *CLabel) GetPlainTextInfoAtWidth(width int) (maxWidth, lineCount int) {
	if l.tbuffer == nil {
		return -1, -1
	}
	_, lineWrapMode, ellipsize, justify, _ := l.Settings()
	l.RLock()
	content := l.tbuffer.PlainText(lineWrapMode, ellipsize, justify, width)
	lines := strings.Split(content, "\n")
	lineCount = len(lines)
	for _, line := range lines {
		size := len(line)
		if size > maxWidth {
			maxWidth = size
		}
	}
	l.RUnlock()
	return
}

// GetSizeRequest returns the requested size of the Label taking into account
// the label's content and any padding set.
//
// Locking: read
func (l *CLabel) GetSizeRequest() (width, height int) {
	size := ptypes.NewRectangle(l.CWidget.GetSizeRequest())
	if size.W <= -1 {
		if wc := l.GetWidthChars(); wc <= -1 {
			if mwc := l.GetMaxWidthChars(); mwc <= -1 {
				alloc := l.GetAllocation()
				if alloc.W > 0 {
					size.W, _ = l.GetPlainTextInfoAtWidth(alloc.W)
				} else {
					size.W, _ = l.GetPlainTextInfo()
				}
			} else {
				size.W = mwc
			}
		} else {
			size.W = wc
		}
	}
	if size.H <= -1 {
		_, size.H = l.GetPlainTextInfoAtWidth(size.W)
	}
	// add padding
	xPadding, yPadding := l.GetPadding()
	l.RLock()
	size.W += xPadding * 2
	size.H += yPadding * 2
	l.RUnlock()
	return size.W, size.H
}

func (l *CLabel) getMaxCharsRequest() (maxWidth int) {
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

func (l *CLabel) refreshMnemonics() {
	if w := l.GetWindow(); w != nil {
		if widget := l.GetMnemonicWidget(); widget != nil {
			w.RemoveWidgetMnemonics(widget)
		} else {
			if parent := l.GetParent(); parent != nil {
				w.RemoveWidgetMnemonics(parent)
			}
		}
	}
	if !GetDefaultSettings().GetEnableMnemonics() {
		return
	}
	if l.GetUseUnderline() {
		if w := l.GetWindow(); w != nil {
			if keyval := l.GetMnemonicKeyVal(); keyval > 0 {
				if widget := l.GetMnemonicWidget(); widget != nil {
					w.AddMnemonic(keyval, widget)
				} else {
					if parent := l.GetParent(); parent != nil {
						if pw, ok := parent.(Sensitive); ok && pw.IsSensitive() && pw.IsVisible() {
							w.AddMnemonic(keyval, pw)
						}
					}
				}
			}
		}
	}
}

func (l *CLabel) refreshTextBuffer() (err error) {
	style := l.GetThemeRequest().Content.Normal
	useUnderline := l.GetUseUnderline()
	l.Lock()
	var markup bool
	if markup, err = l.GetBoolProperty(PropertyUseMarkup); err == nil && markup {
		var m memphis.Tango
		if m, err = memphis.NewMarkup(l.text, style); err != nil {
			// tbuffer must always be valid, default to plain text on error
			l.tbuffer = memphis.NewTextBuffer(l.text, style, useUnderline)
		} else {
			// use the markup tbuffer
			l.tbuffer = m.TextBuffer(useUnderline)
		}
	} else {
		// plain text tbuffer
		l.tbuffer = memphis.NewTextBuffer(l.text, style, useUnderline)
	}
	l.Unlock()
	return
}

func (l *CLabel) resize(data []interface{}, argv ...interface{}) enums.EventFlag {
	alloc := l.GetAllocation()
	if !l.IsVisible() || alloc.W <= 0 || alloc.H <= 0 {
		l.LogTrace("not visible, zero width or zero height")
		return enums.EVENT_PASS
	}

	theme := l.GetThemeRequest()
	origin := l.GetOrigin()
	id := l.ObjectID()
	xPad, _ := l.GetPadding()
	_, yAlign := l.GetAlignment()

	size := ptypes.NewRectangle(alloc.W, alloc.H)
	local := ptypes.MakePoint2I(xPad, 0)
	_, size.H = l.GetPlainTextInfoAtWidth(alloc.W - (xPad * 2))

	l.Lock()

	if err := memphis.ConfigureSurface(id, origin, alloc, theme.Content.Normal); err != nil {
		l.LogErr(err)
	}

	if size.H < alloc.H {
		delta := alloc.H - size.H
		local.Y += int(float64(delta) * yAlign)
	}

	if err := memphis.ConfigureSurface(l.tid, local, *size, theme.Content.Normal); err != nil {
		l.LogErr(err)
	}

	l.Unlock()
	l.Invalidate()
	return enums.EVENT_STOP
}

func (l *CLabel) invalidate(data []interface{}, argv ...interface{}) enums.EventFlag {
	theme := l.GetThemeRequest()
	if err := l.refreshTextBuffer(); err != nil {
		l.LogErr(err)
	}
	l.refreshMnemonics()
	id := l.ObjectID()
	l.Lock()
	theme.Content.FillRune = rune(0)
	if err := memphis.FillSurface(id, theme); err != nil {
		l.LogErr(err)
	}
	if err := memphis.FillSurface(l.tid, theme); err != nil {
		l.LogErr(err)
	}
	l.Unlock()
	return enums.EVENT_PASS
}

func (l *CLabel) draw(data []interface{}, argv ...interface{}) enums.EventFlag {
	if surface, ok := argv[1].(*memphis.CSurface); ok {
		alloc := l.GetAllocation()
		if !l.IsVisible() || alloc.W <= 0 || alloc.H <= 0 {
			l.LogTrace("not visible, zero width or zero height")
			return enums.EVENT_PASS
		}

		singleLineMode, lineWrapMode, ellipsize, justify, _ := l.Settings()

		l.Lock()
		defer l.Unlock()

		if l.tbuffer != nil {
			if tSurface, err := memphis.GetSurface(l.tid); err != nil {
				l.LogErr(err)
			} else {
				if f := l.tbuffer.Draw(tSurface, singleLineMode, lineWrapMode, ellipsize, justify, enums.ALIGN_TOP); f == enums.EVENT_STOP {
					if err := surface.CompositeSurface(tSurface); err != nil {
						l.LogErr(err)
					}
				}
			}
		}

		if debug, _ := l.GetBoolProperty(cdk.PropertyDebug); debug {
			surface.DebugBox(paint.ColorSilver, l.ObjectInfo())
		}
		return enums.EVENT_STOP
	}
	return enums.EVENT_PASS
}

// A list of style attributes to apply to the text of the label.
// Flags: Read / Write
const PropertyAttributes cdk.Property = "attributes"

// The current position of the insertion cursor in chars.
// Flags: Read
// Allowed values: >= 0
// Default value: 0
const PropertyCursorPosition cdk.Property = "cursor-position"

// The preferred place to ellipsize the string, if the label does not have
// enough room to display the entire string, specified as a bool.
// Flags: Read / Write
// Default value: false
const PropertyEllipsize cdk.Property = "ellipsize"

// The alignment of the lines in the text of the label relative to each
// other. This does NOT affect the alignment of the label within its
// allocation. See Misc::xAlign for that.
// Flags: Read / Write
// Default value: enums.JUSTIFY_LEFT
const PropertyJustify cdk.Property = "justify"

// The text of the label.
// Flags: Read / Write
// Default value: ""
const PropertyLabel cdk.Property = "label"

// The desired maximum width of the label, in characters. If this property is
// set to -1, the width will be calculated automatically, otherwise the label
// will request space for no more than the requested number of characters. If
// the “width-chars” property is set to a positive value, then the
// "max-width-chars" property is ignored.
// Flags: Read / Write
// Allowed values: >= -1
// Default value: -1
const PropertyMaxWidthChars cdk.Property = "max-width-chars"

// The mnemonic accelerator key for this label.
// Flags: Read
// Default value: 16777215
const PropertyMnemonicKeyVal cdk.Property = "mnemonic-key-val"

// The widget to be activated when the label's mnemonic key is pressed.
// Flags: Read / Write
const PropertyMnemonicWidget cdk.Property = "mnemonic-widget"

// Whether the label text can be selected with the mouse.
// Flags: Read / Write
// Default value: FALSE
const PropertySelectable cdk.Property = "selectable"

// The position of the opposite end of the selection from the cursor in
// chars.
// Flags: Read
// Allowed values: >= 0
// Default value: 0
const PropertySelectionBound cdk.Property = "selection-bound"

// Whether the label is in single line mode. In single line mode, the height
// of the label does not depend on the actual text, it is always set to
// ascent + descent of the font. This can be an advantage in situations where
// resizing the label because of text changes would be distracting, e.g. in a
// statusbar.
// Flags: Read / Write
// Default value: FALSE
const PropertySingleLineMode cdk.Property = "single-line-mode"

// Set this property to TRUE to make the label track which links have been
// clicked. It will then apply the ::visited-link-color color, instead of
// ::link-color.
// Flags: Read / Write
// Default value: TRUE
const PropertyTrackVisitedLinks cdk.Property = "track-visited-links"

// The text of the label includes XML markup. See pango_parse_markup.
// Flags: Read / Write
// Default value: FALSE
const PropertyUseMarkup cdk.Property = "use-markup"

// If set, an underline in the text indicates the next character should be
// used for the mnemonic accelerator key.
// Flags: Read / Write
// Default value: FALSE
const PropertyLabelUseUnderline cdk.Property = "use-underline"

// The desired width of the label, in characters. If this property is set to
// -1, the width will be calculated automatically, otherwise the label will
// request either 3 characters or the property value, whichever is greater.
// If the "width-chars" property is set to a positive value, then the
// “max-width-chars” property is ignored.
// Flags: Read / Write
// Allowed values: >= -1
// Default value: -1
const PropertyWidthChars cdk.Property = "width-chars"

// If set, wrap lines if the text becomes too wide.
// Flags: Read / Write
// Default value: FALSE
const PropertyWrap cdk.Property = "wrap"

// If line wrapping is on (see the “wrap” property) this controls how the
// line wrapping is done. The default is PANGO_WRAP_WORD, which means wrap on
// word boundaries.
// Flags: Read / Write
// Default value: PANGO_WRAP_WORD
const PropertyWrapMode cdk.Property = "wrap-mode"

// A keybinding signal which gets emitted when the user activates a link in
// the label. Applications may also emit the signal with
// g_signal_emit_by_name if they need to control activation of URIs
// programmatically. The default bindings for this signal are all forms of
// the Enter key.
const SignalActivateCurrentLink cdk.Signal = "activate-current-link"

// The signal which gets emitted to activate a URI. Applications may connect
// to it to override the default behaviour, which is to call ShowUri.
const SignalActivateLink cdk.Signal = "activate-link"

// The ::copy-clipboard signal is a which gets emitted to copy the selection
// to the clipboard. The default binding for this signal is Ctrl-c.
const SignalCopyClipboard cdk.Signal = "copy-clipboard"

// The ::move-cursor signal is a which gets emitted when the user initiates a
// cursor movement. If the cursor is not visible in entry , this signal
// causes the viewport to be moved instead. Applications should not connect
// to it, but may emit it with g_signal_emit_by_name if they need to
// control the cursor programmatically. The default bindings for this signal
// come in two variants, the variant with the Shift modifier extends the
// selection, the variant without the Shift modifer does not. There are too
// many key combinations to list them all here.
// Listener function arguments:
// 	step MovementStep	the granularity of the move, as a GtkMovementStep
// 	count int	the number of step units to move
// 	extendSelection bool	TRUE if the move should extend the selection
const SignalMoveCursor cdk.Signal = "move-cursor"

// The ::populate-popup signal gets emitted before showing the context menu
// of the label. Note that only selectable labels have context menus. If you
// need to add items to the context menu, connect to this signal and append
// your menuitems to the menu .
// Listener function arguments:
// 	menu Menu	the menu that is being populated
const SignalPopulatePopup cdk.Signal = "populate-popup"

const LabelInvalidateHandle = "label-invalidate-handler"

const LabelResizeHandle = "label-resize-handler"

const LabelDrawHandle = "label-draw-handler"
