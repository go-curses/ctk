package ctk

import (
	"github.com/go-curses/cdk"
	cmath "github.com/go-curses/cdk/lib/math"
)

// CDK type-tag for Misc objects
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
type Misc interface {
	Widget
	Buildable

	Init() (already bool)
	SetAlignment(xAlign float64, yAlign float64)
	SetPadding(xPad int, yPad int)
	GetAlignment() (xAlign float64, yAlign float64)
	GetPadding() (xPad int, yPad int)
}

// The CMisc structure implements the Misc interface and is
// exported to facilitate type embedding with custom implementations. No member
// variables are exported as the interface methods are the only intended means
// of interacting with Misc objects
type CMisc struct {
	CWidget
}

// Misc object initialization. This must be called at least once to setup
// the necessary defaults and allocate any memory structures. Calling this more
// than once is safe though unnecessary. Only the first call will result in any
// effect upon the Misc instance
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

func (m *CMisc) SetPadding(xPad int, yPad int) {
	if err := m.SetIntProperty(PropertyXPad, xPad); err != nil {
		m.LogErr(err)
	}
	if err := m.SetIntProperty(PropertyYPad, yPad); err != nil {
		m.LogErr(err)
	}
}

// Gets the X and Y alignment of the widget within its allocation. See
// SetAlignment.
// Parameters:
//      xAlign  location to store X alignment of misc
// , or NULL.
//      yAlign  location to store Y alignment of misc
// , or NULL.
func (m *CMisc) GetAlignment() (xAlign float64, yAlign float64) {
	xAlign, _ = m.GetFloatProperty(PropertyXAlign)
	yAlign, _ = m.GetFloatProperty(PropertyYAlign)
	return
}

// Gets the padding in the X and Y directions of the widget. See
// SetPadding.
// Parameters:
//      xPad    location to store padding in the X
// direction, or NULL.
//      yPad    location to store padding in the Y
// direction, or NULL.
func (m *CMisc) GetPadding() (xPad int, yPad int) {
	xPad, _ = m.GetIntProperty(PropertyXPad)
	yPad, _ = m.GetIntProperty(PropertyYPad)
	return
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
