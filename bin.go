package ctk

import (
	"github.com/go-curses/cdk"
	"github.com/go-curses/ctk/lib/enums"
)

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
//
// The Bin Widget is a Container with just one child. It is not very useful
// itself, but it is useful for deriving subclasses, since it provides common
// code needed for handling a single child widget. Many CTK widgets are
// subclasses of Bin, including Window, Button, Frame or ScrolledWindow.
type Bin interface {
	Container
	Buildable

	GetChild() (value Widget)
}

var _ Bin = (*CBin)(nil)

// The CBin structure implements the Bin interface and is exported to
// facilitate type embedding with custom implementations. No member variables
// are exported as the interface methods are the only intended means of
// interacting with Bin objects.
type CBin struct {
	CContainer
}

// MakeBin is used by the Buildable system to construct a new Bin.
func MakeBin() Bin {
	return NewBin()
}

// NewBin is the constructor for new Bin instances.
func NewBin() Bin {
	a := new(CBin)
	a.Init()
	return a
}

// Init initializes a Bin object. This must be called at least once to
// set up the necessary defaults and allocate any memory structures. Calling
// this more than once is safe though unnecessary. Only the first call will
// result in any effect upon the Bin instance. Init is used in the
// NewBin constructor and only necessary when implementing a derivative
// Bin type.
func (b *CBin) Init() (already bool) {
	if b.InitTypeItem(TypeBin, b) {
		return true
	}
	b.CContainer.Init()
	b.flags = enums.NULL_WIDGET_FLAG
	b.SetFlags(enums.PARENT_SENSITIVE)
	return false
}

// GetChild is a convenience method to return the first child in the Bin
// Container. Returns the Widget or `nil` if the Bin contains no child widget.
//
// Locking: read
func (b *CBin) GetChild() (value Widget) {
	b.RLock()
	defer b.RUnlock()
	if len(b.children) > 0 {
		value = b.children[0]
	}
	return
}

// Add the given widget to the Bin, if the Bin is full (has one child already)
// the given Widget replaces the existing Widget.
//
// Locking: write
func (b *CBin) Add(w Widget) {
	children := b.GetChildren()
	for _, child := range children {
		b.Remove(child)
	}
	b.CContainer.Add(w)
}