package ctk

import (
	"github.com/go-curses/cdk"
)

const TypeClipboard cdk.CTypeTag = "ctk-settings"

func init() {
	_ = cdk.TypesManager.AddType(TypeClipboard, func() interface{} { return nil })
}

var ctkDefaultClipboard *CClipboard
var _ Clipboard = (*CClipboard)(nil)

// Clipboard Hierarchy:
//	Object
//	  +- Clipboard
type Clipboard interface {
	Object

	GetText() (text string)
	SetText(text string)
	Copy(text string)
	Paste(text string)
}

type CClipboard struct {
	CObject
}

func GetDefaultClipboard() (clipboard Clipboard) {
	if ctkDefaultClipboard == nil {
		ctkDefaultClipboard = new(CClipboard)
		ctkDefaultClipboard.Init()
	}
	return ctkDefaultClipboard
}

func (c *CClipboard) Init() (already bool) {
	if c.InitTypeItem(TypeClipboard, c) {
		return true
	}
	c.CObject.Init()
	_ = c.InstallProperty(PropertyText, cdk.StringProperty, true, "")
	return false
}

// SetText updates the clipboard's cache of pasted content
func (c *CClipboard) SetText(text string) {
	if err := c.SetStringProperty(PropertyText, text); err != nil {
		c.LogErr(err)
	}
}

// GetText retrieves the clipboard's cache of pasted content
func (c *CClipboard) GetText() (text string) {
	var err error
	if text, err = c.GetStringProperty(PropertyText); err != nil {
		c.LogErr(err)
	}
	return
}

// Copy updates the clipboard's cache of pasted content and passes the copy
// event to the underlying operating system (if supported) using OSC52 terminal
// sequences
func (c *CClipboard) Copy(text string) {
	c.SetText(text)
	c.Emit(SignalCopy, c, text)
	c.LogDebug("text: \"%v\"", text)
	d := cdk.GetDefaultDisplay()
	if d != nil {
		s := d.Screen()
		if s != nil {
			s.CopyToClipboard(text)
		}
	}
}

// Paste updates the clipboard's cache of pasted content and emits a "Paste"
// event itself
func (c *CClipboard) Paste(text string) {
	c.SetText(text)
	c.Emit(SignalPaste, c, text)
	c.LogDebug("text: \"%v\"", text)
}

const SignalCopy cdk.Signal = "copy"

const SignalPaste cdk.Signal = "paste"