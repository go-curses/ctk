package ctk

import (
	"github.com/go-curses/cdk"
	cmath "github.com/go-curses/cdk/lib/math"
)

const TypeMisc cdk.CTypeTag = "ctk-misc"

func init() {
	_ = cdk.TypesManager.AddType(TypeMisc, nil)
}

// Misc Hierarchy:
//      Object
//        +- Widget
//          +- Misc
//            +- Label
//            +- Arrow
//            +- Image
//            +- Pixmap
//
// The Misc Widget is intended primarily as a base type for other Widget
// implementations where there is a necessity for alignment and padding
// properties.
type Misc interface {
	Widget
	Buildable

	Init() (already bool)
	GetAlignment() (xAlign float64, yAlign float64)
	SetAlignment(xAlign float64, yAlign float64)
	GetPadding() (xPad int, yPad int)
	SetPadding(xPad int, yPad int)
}

// The CMisc structure implements the Misc interface and is exported
// to facilitate type embedding with custom implementations. No member variables
// are exported as the interface methods are the only intended means of
// interacting with Misc objects.
type CMisc struct {
	CWidget
}

// Init initializes a Misc object. This must be called at least once to
// set up the necessary defaults and allocate any memory structures. Calling
// this more than once is safe though unnecessary. Only the first call will
// result in any effect upon the Misc instance. Init is used in the
// NewMisc constructor and only necessary when implementing a derivative
// Misc type.
func (m *CMisc) Init() (already bool) {
	if m.InitTypeItem(TypeMisc, m) {
		return true
	}
	m.CWidget.Init()
	_ = m.InstallProperty(PropertyXAlign, cdk.FloatProperty, true, 0.0)
	_ = m.InstallProperty(PropertyXPad, cdk.IntProperty, true, 0)
	_ = m.InstallProperty(PropertyYAlign, cdk.FloatProperty, true, 0.0)
	_ = m.InstallProperty(PropertyYPad, cdk.IntProperty, true, 0)
	return false
}

// GetAlignment returns the X and Y alignment of the widget within its
// allocation.
// See: SetAlignment()
func (m *CMisc) GetAlignment() (xAlign float64, yAlign float64) {
	xAlign, _ = m.GetFloatProperty(PropertyXAlign)
	yAlign, _ = m.GetFloatProperty(PropertyYAlign)
	return
}

// SetAlignment is a convenience method to set the x-align and y-align
// properties of the Misc Widget.
func (m *CMisc) SetAlignment(xAlign float64, yAlign float64) {
	xAlign = cmath.ClampF(xAlign, 0.0, 1.0)
	yAlign = cmath.ClampF(yAlign, 0.0, 1.0)
	if err := m.SetFloatProperty(PropertyXAlign, xAlign); err != nil {
		m.LogErr(err)
	}
	if err := m.SetFloatProperty(PropertyYAlign, yAlign); err != nil {
		m.LogErr(err)
	}
}

// GetPadding returns the padding in the X and Y directions of the widget.
// See: SetPadding()
func (m *CMisc) GetPadding() (xPad int, yPad int) {
	xPad, _ = m.GetIntProperty(PropertyXPad)
	yPad, _ = m.GetIntProperty(PropertyYPad)
	return
}

// SetPadding is a convenience method to set x-pad and y-pad properties of the
// Misc Widget.
func (m *CMisc) SetPadding(xPad int, yPad int) {
	if err := m.SetIntProperty(PropertyXPad, xPad); err != nil {
		m.LogErr(err)
	}
	if err := m.SetIntProperty(PropertyYPad, yPad); err != nil {
		m.LogErr(err)
	}
}

// The horizontal alignment, from 0 (left) to 1 (right). Reversed for RTL
// layouts.
// Flags: Read / Write
// Allowed values: [0,1]
// Default value: 0.5
// const PropertyXAlign cdk.Property = "xalign"

// The amount of space to add on the left and right of the widget, in pixels.
// Flags: Read / Write
// Allowed values: >= 0
// Default value: 0
const PropertyXPad cdk.Property = "x-pad"

// The vertical alignment, from 0 (top) to 1 (bottom).
// Flags: Read / Write
// Allowed values: [0,1]
// Default value: 0.5
// const PropertyYAlign cdk.Property = "y-align"

// The amount of space to add on the top and bottom of the widget, in pixels.
// Flags: Read / Write
// Allowed values: >= 0
// Default value: 0
const PropertyYPad cdk.Property = "y-pad"
