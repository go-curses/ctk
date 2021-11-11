package ctk

import (
	"github.com/go-curses/cdk"
)

// CDK type-tag for Bin objects
const TypeBin cdk.CTypeTag = "ctk-bin"

func init() {
	_ = cdk.TypesManager.AddType(TypeBin, nil)
}

// Bin Hierarchy:
//	Object
//	  +- Widget
//	    +- Container
//	      +- Bin
//	        +- Window
//	        +- Alignment
//	        +- Frame
//	        +- Button
//	        +- Item
//	        +- ComboBox
//	        +- EventBox
//	        +- Expander
//	        +- HandleBox
//	        +- ToolItem
//	        +- ScrolledWindow
//	        +- Viewport
// The Bin widget is a container with just one child. It is not very useful
// itself, but it is useful for deriving subclasses, since it provides common
// code needed for handling a single child widget. Many CTK widgets are
// subclasses of Bin, including Window, Button, Frame, HandleBox or
// ScrolledWindow.
type Bin interface {
	Container
	Buildable

	Init() (already bool)
	GetChild() (value Widget)
	Add(w Widget)
	ShowAll()
}

// The CBin structure implements the Bin interface and is
// exported to facilitate type embedding with custom implementations. No member
// variables are exported as the interface methods are the only intended means
// of interacting with Bin objects
type CBin struct {
	CContainer
}

// Bin object initialization. This must be called at least once to setup
// the necessary defaults and allocate any memory structures. Calling this more
// than once is safe though unnecessary. Only the first call will result in any
// effect upon the Bin instance
func (b *CBin) Init() (already bool) {
	if b.InitTypeItem(TypeBin, b) {
		return true
	}
	b.CContainer.Init()
	b.flags = NULL_WIDGET_FLAG
	b.SetFlags(PARENT_SENSITIVE)
	return false
}

// Gets the child of the Bin, or NULL if the bin contains no child widget.
// The returned widget does not have a reference added, so you do not need to
// unref it.
// Returns:
// 	pointer to child of the Bin.
// 	[transfer none]
func (b *CBin) GetChild() (value Widget) {
	if len(b.children) > 0 {
		value = b.children[0]
	}
	return
}

// Add the given widget to the Bin, if the Bin is full (has one child already)
// the given Widget replaces the existing Widget.
func (b *CBin) Add(w Widget) {
	if len(b.children) > 0 {
		children := b.GetChildren()
		for _, child := range children {
			b.Remove(child)
		}
	}
	b.CContainer.Add(w)
}

func (b *CBin) ShowAll() {
	b.Show()
	for _, child := range b.children {
		child.ShowAll()
	}
}
