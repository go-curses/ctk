package ctk

import (
	"github.com/go-curses/cdk"
	cenums "github.com/go-curses/cdk/lib/enums"
	cmath "github.com/go-curses/cdk/lib/math"
	"github.com/go-curses/cdk/lib/paint"
	"github.com/go-curses/cdk/lib/ptypes"
	"github.com/go-curses/cdk/memphis"

	"github.com/go-curses/ctk/lib/enums"
)

const TypeViewport cdk.CTypeTag = "ctk-viewport"

func init() {
	_ = cdk.TypesManager.AddType(TypeViewport, func() interface{} { return MakeViewport() })
}

// Viewport Hierarchy:
//	Object
//	  +- Widget
//	    +- Container
//	      +- Bin
//	        +- Viewport
//
// The Viewport widget acts as an adaptor class, implementing scrollability
// for child widgets that lack their own scrolling capabilities. Use Viewport
// to scroll child widgets such as Table, Box, and so on. If a widget has
// native scrolling abilities, such as TextView, TreeView IconView, it can
// be added to a ScrolledWindow with Container.Add. If a widget does not, you
// must first add the widget to a Viewport, then add the viewport to the
// scrolled window. The convenience function ScrolledWindow.AddWithViewport does
// exactly this, so you can ignore the presence of the viewport.
type Viewport interface {
	Bin

	Init() (already bool)
	GetHAdjustment() (adjustment Adjustment)
	SetHAdjustment(adjustment Adjustment)
	GetVAdjustment() (adjustment Adjustment)
	SetVAdjustment(adjustment Adjustment)
	GetShadowType() (shadowType enums.ShadowType)
	SetShadowType(shadowType enums.ShadowType)
	GetBinWindow() (value Window)
	GetViewWindow() (value Window)
}

// The CViewport structure implements the Viewport interface and is exported
// to facilitate type embedding with custom implementations. No member variables
// are exported as the interface methods are the only intended means of
// interacting with Viewport objects.
type CViewport struct {
	CBin

	hAdjustment Adjustment
	vAdjustment Adjustment
}

// MakeViewport is used by the Buildable system to construct a new Viewport.
func MakeViewport() Viewport {
	return NewViewport(nil, nil)
}

// NewViewport is the constructor for new Viewport instances.
func NewViewport(hAdjustment, vAdjustment Adjustment) Viewport {
	v := new(CViewport)
	v.hAdjustment = hAdjustment
	v.vAdjustment = vAdjustment
	v.Init()
	return v
}

// Init initializes a Viewport object. This must be called at least once to
// set up the necessary defaults and allocate any memory structures. Calling
// this more than once is safe though unnecessary. Only the first call will
// result in any effect upon the Viewport instance. Init is used in the
// NewViewport constructor and only necessary when implementing a derivative
// Viewport type.
func (v *CViewport) Init() (already bool) {
	if v.InitTypeItem(TypeViewport, v) {
		return true
	}
	v.CBin.Init()
	v.flags = enums.NULL_WIDGET_FLAG
	v.SetFlags(enums.SENSITIVE | enums.PARENT_SENSITIVE | enums.APP_PAINTABLE)
	if v.vAdjustment == nil {
		v.vAdjustment = NewAdjustment(0, 0, 0, 0, 0, 0)
	}
	if v.hAdjustment == nil {
		v.hAdjustment = NewAdjustment(0, 0, 0, 0, 0, 0)
	}
	_ = v.InstallProperty(PropertyHAdjustment, cdk.StructProperty, true, v.hAdjustment)
	_ = v.InstallProperty(PropertyViewportShadowType, cdk.StructProperty, true, nil)
	_ = v.InstallProperty(PropertyVAdjustment, cdk.StructProperty, true, v.vAdjustment)
	v.Connect(SignalResize, ViewportResizeHandle, v.resize)
	v.Connect(SignalDraw, ViewportDrawHandle, v.draw)
	return false
}

// GetHAdjustment returns the horizontal adjustment of the viewport.
// See: SetHAdjustment()
func (v *CViewport) GetHAdjustment() (adjustment Adjustment) {
	var ok bool
	if value, err := v.GetStructProperty(PropertyHAdjustment); err != nil {
		v.LogErr(err)
	} else if adjustment, ok = value.(*CAdjustment); !ok {
		v.LogError("value stored in %v property is not of *CAdjustment type: %v (%T)", PropertyHAdjustment, value, value)
	}
	return
}

// SetHAdjustment replaces the horizontal adjustment of the viewport with the
// given adjustment.
func (v *CViewport) SetHAdjustment(adjustment Adjustment) {
	if err := v.SetStructProperty(PropertyHAdjustment, adjustment); err != nil {
		v.LogErr(err)
	} else {
		v.hAdjustment = adjustment
	}
}

// GetVAdjustment returns the vertical adjustment of the viewport.
// See: SetVAdjustment()
func (v *CViewport) GetVAdjustment() (adjustment Adjustment) {
	var ok bool
	if value, err := v.GetStructProperty(PropertyVAdjustment); err != nil {
		v.LogErr(err)
	} else if adjustment, ok = value.(*CAdjustment); !ok {
		v.LogError("value stored in %v property is not of *CAdjustment type: %v (%T)", PropertyVAdjustment, value, value)
	}
	return
}

// SetHAdjustment replaces the horizontal adjustment of the viewport with the
// given adjustment.
func (v *CViewport) SetVAdjustment(adjustment Adjustment) {
	if err := v.SetStructProperty(PropertyVAdjustment, adjustment); err != nil {
		v.LogErr(err)
	} else {
		v.vAdjustment = adjustment
	}
}

// GetShadowType returns the shadow type of the Viewport.
// See: SetShadowType()
//
// Note that usage of this within CTK is unimplemented at this time
func (v *CViewport) GetShadowType() (shadowType enums.ShadowType) {
	var ok bool
	if value, err := v.GetStructProperty(PropertyViewportShadowType); err != nil {
		v.LogErr(err)
	} else if shadowType, ok = value.(enums.ShadowType); !ok {
		v.LogError("value stored in %v property is not of ShadowType type: %v (%T)", PropertyViewportShadowType, value, value)
	}
	return
}

// SetShadowType updates the shadow type of the viewport.
//
// Note that usage of this within CTK is unimplemented at this time
func (v *CViewport) SetShadowType(shadowType enums.ShadowType) {
	if err := v.SetStructProperty(PropertyViewportShadowType, shadowType); err != nil {
		v.LogErr(err)
	}
}

// GetBinWindow returns the bin window of the Viewport.
//
// Note that usage of this within CTK is unimplemented at this time
func (v *CViewport) GetBinWindow() (value Window) {
	v.LogWarn("method unimplemented")
	return nil
}

// GetViewWindow returns the view window of the Viewport.
//
// Note that usage of this within CTK is unimplemented at this time
func (v *CViewport) GetViewWindow() (value Window) {
	v.LogError("method unimplemented")
	return nil
}

func (v *CViewport) resize(data []interface{}, argv ...interface{}) cenums.EventFlag {
	v.LockDraw()
	defer v.UnlockDraw()

	alloc := v.GetAllocation()
	child := v.GetChild()
	horizontal, vertical := v.GetHAdjustment(), v.GetVAdjustment()
	if alloc.W == 0 || alloc.H == 0 {
		if child != nil {
			child.SetAllocation(ptypes.MakeRectangle(0, 0))
			child.Resize()
		}
		if horizontal != nil {
			horizontal.Configure(0, 0, 0, 0, 0, 0)
		}
		if vertical != nil {
			vertical.Configure(0, 0, 0, 0, 0, 0)
		}
		return cenums.EVENT_STOP
	}

	hValue, hLower, hUpper, hStepIncrement, hPageIncrement, hPageSize := 0, 0, 0, 0, 0, 0
	vValue, vLower, vUpper, vStepIncrement, vPageIncrement, vPageSize := 0, 0, 0, 0, 0, 0

	if child != nil {

		childSize := ptypes.NewRectangle(child.GetSizeRequest())
		if childSize.W <= -1 {
			childSize.W = alloc.W
		}
		if childSize.H <= -1 {
			childSize.H = alloc.H
		}

		hChanged, vChanged := false, false
		if childSize.W >= alloc.W {
			delta := childSize.W - alloc.W
			hLower, hUpper, hStepIncrement, hPageIncrement, hPageSize = 0, delta, 1, alloc.W/2, alloc.W
			if horizontal != nil {
				hValue = cmath.ClampI(horizontal.GetValue(), 0, hUpper)
				horizontal.Configure(hValue, hLower, hUpper, hStepIncrement, hPageIncrement, hPageSize)
			}
			hChanged = true
		}
		if childSize.H >= alloc.H {
			delta := childSize.H - alloc.H
			vLower, vUpper, vStepIncrement, vPageIncrement, vPageSize = 0, delta, 1, alloc.H/2, alloc.H
			if vertical != nil {
				vValue = cmath.ClampI(vertical.GetValue(), 0, vUpper)
				vertical.Configure(vValue, vLower, vUpper, vStepIncrement, vPageIncrement, vPageSize)
			}
			vChanged = true
		}

		origin := v.GetOrigin()
		childOrigin := child.GetOrigin()
		if hChanged {
			childOrigin.X = origin.X - hValue
		}
		if vChanged {
			childOrigin.Y = origin.Y - vValue
		}

		child.SetOrigin(childOrigin.X, childOrigin.Y)
		child.SetAllocation(*childSize)
		child.Resize()
		// v.LogDebug("child resized: origin=%v, alloc=%v", child.GetOrigin(), child.GetAllocation())

		if hChanged {
			if horizontal != nil {
				horizontal.Changed()
			}
		}
		if vChanged {
			if vertical != nil {
				vertical.Changed()
			}
		}
	}

	v.Invalidate()
	return cenums.EVENT_STOP
}

func (v *CViewport) draw(data []interface{}, argv ...interface{}) cenums.EventFlag {
	v.LockDraw()
	defer v.UnlockDraw()

	if surface, ok := argv[1].(*memphis.CSurface); ok {
		size := v.GetAllocation()
		if !v.IsVisible() || size.W <= 0 || size.H <= 0 {
			v.LogTrace("not visible, zero width or zero height")
			return cenums.EVENT_PASS
		}

		surface.Fill(v.GetThemeRequest())

		if child := v.GetChild(); child != nil {
			child.Draw()
			if err := surface.Composite(child.ObjectID()); err != nil {
				v.LogError("composite error: %v", err)
			}
		}

		if debug, _ := v.GetBoolProperty(cdk.PropertyDebug); debug {
			surface.DebugBox(paint.ColorSilver, v.ObjectInfo())
		}

		return cenums.EVENT_STOP
	}
	return cenums.EVENT_PASS
}

// The Adjustment that determines the values of the horizontal position
// for this viewport.
// Flags: Read / Write / Construct
const PropertyViewportHAdjustment cdk.Property = "hadjustment"

// Determines how the shadowed box around the viewport is drawn.
// Flags: Read / Write
// Default value: GTK_SHADOW_IN
const PropertyViewportShadowType cdk.Property = "shadow-type"

// The Adjustment that determines the values of the vertical position for
// this viewport.
// Flags: Read / Write / Construct
const PropertyViewportVAdjustment cdk.Property = "vadjustment"

// Set the scroll adjustments for the viewport. Usually scrolled containers
// like ScrolledWindow will emit this signal to connect two instances of
// Scrollbar to the scroll directions of the Viewport.
// Listener function arguments:
// 	vertical Adjustment	the vertical GtkAdjustment
// 	arg2 Adjustment
const SignalSetScrollAdjustments cdk.Signal = "set-scroll-adjustments"

const ViewportInvalidateHandle = "viewport-invalidate-handler"

const ViewportResizeHandle = "viewport-resize-handler"

const ViewportDrawHandle = "viewport-draw-handler"