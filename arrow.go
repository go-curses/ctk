// Copyright (c) 2023  The Go-Curses Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ctk

import (
	"github.com/mattn/go-runewidth"

	"github.com/go-curses/cdk"
	cenums "github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/cdk/lib/paint"
	"github.com/go-curses/cdk/lib/ptypes"
	"github.com/go-curses/cdk/memphis"
	"github.com/go-curses/ctk/lib/enums"
)

const TypeArrow cdk.CTypeTag = "ctk-arrow"

func init() {
	_ = cdk.TypesManager.AddType(TypeArrow, func() interface{} { return MakeArrow() })
}

// Arrow Hierarchy:
//
//	Object
//	  +- Widget
//	    +- Misc
//	      +- Arrow
//
// The Arrow Widget should be used to draw simple arrows that need to point in
// one of the four cardinal directions (up, down, left, or right). The style of
// the arrow can be one of shadow in, shadow out, etched in, or etched out. Note
// that these directions and style types may be amended in versions of CTK
// to come. Arrow will fill any space allotted to it, but since it is
// inherited from Misc, it can be padded and/or aligned, to fill exactly the
// space the programmer desires. Arrows are created with a call to
// NewArrow. The direction or style of an arrow can be changed after
// creation by using Set.
type Arrow interface {
	Misc
	Buildable

	GetArrowType() (arrow enums.ArrowType)
	SetArrowType(arrow enums.ArrowType)
	GetArrowRune() (r rune, width int)
	GetArrowRuneSet() (ars paint.ArrowRuneSet)
	SetArrowRuneSet(ars paint.ArrowRuneSet)
}

var _ Arrow = (*CArrow)(nil)

// The CArrow structure implements the Arrow interface and is exported to
// facilitate type embedding with custom implementations. No member variables
// are exported as the interface methods are the only intended means of
// interacting with Arrow objects.
type CArrow struct {
	CMisc

	arrowRuneSet *paint.ArrowRuneSet
}

// MakeArrow is used by the Buildable system to construct a new Arrow with a
// default ArrowType setting of ArrowRight.
func MakeArrow() Arrow {
	return NewArrow(enums.ArrowRight)
}

// NewArrow is the constructor for new Arrow instances.
func NewArrow(arrow enums.ArrowType) Arrow {
	a := new(CArrow)
	a.Init()
	a.SetArrowType(arrow)
	return a
}

// Init initializes an Arrow object. This must be called at least once to
// set up the necessary defaults and allocate any memory structures. Calling
// this more than once is safe though unnecessary. Only the first call will
// result in any effect upon the Arrow instance. Init is used in the
// NewArrow constructor and only necessary when implementing a derivative
// Arrow type.
func (a *CArrow) Init() bool {
	if a.InitTypeItem(TypeArrow, a) {
		return true
	}
	a.CMisc.Init()
	a.flags = enums.NULL_WIDGET_FLAG
	a.SetFlags(enums.PARENT_SENSITIVE)
	a.SetFlags(enums.APP_PAINTABLE)
	a.arrowRuneSet = nil
	_ = a.InstallBuildableProperty(PropertyArrowType, cdk.StructProperty, true, nil)
	_ = a.InstallBuildableProperty(PropertyArrowShadowType, cdk.StructProperty, true, nil)
	a.Connect(SignalResize, ArrowResizeHandle, a.resize)
	a.Connect(SignalDraw, ArrowDrawHandle, a.draw)
	return false
}

// GetArrowType is a convenience method for returning the ArrowType property
//
// Locking: read
func (a *CArrow) GetArrowType() (arrow enums.ArrowType) {
	a.RLock()
	defer a.RUnlock()
	arrow = enums.ArrowRight // default
	var ok bool
	if sa, err := a.GetStructProperty(PropertyArrowType); err != nil {
		a.LogErr(err)
	} else if arrow, ok = sa.(enums.ArrowType); !ok {
		a.LogErr(err)
	}
	return
}

// SetArrowType is a convenience method for updating the ArrowType property
//
// Parameters:
//
//	arrowType	a valid ArrowType.
//
// Locking: write
func (a *CArrow) SetArrowType(arrow enums.ArrowType) {
	a.Lock()
	if err := a.SetStructProperty(PropertyArrowType, arrow); err != nil {
		a.Unlock()
		a.LogErr(err)
	} else {
		a.Unlock()
		a.Invalidate()
	}
}

// GetArrowRune is a Curses-specific method for returning the go `rune`
// character and its byte width.
//
// Locking: read
func (a *CArrow) GetArrowRune() (r rune, width int) {
	arrowType := a.GetArrowType()

	a.RLock()
	defer a.RUnlock()

	var ars paint.ArrowRuneSet
	if a.arrowRuneSet == nil {
		ars, _ = paint.GetArrows(paint.WideArrow)
	} else {
		ars = *a.arrowRuneSet
	}

	switch arrowType {
	case enums.ArrowUp:
		r = ars.Up
	case enums.ArrowLeft:
		r = ars.Left
	case enums.ArrowDown:
		r = ars.Down
	case enums.ArrowRight:
		r = ars.Right
	}

	width = runewidth.RuneWidth(r)
	return
}

func (a *CArrow) GetArrowRuneSet() (ars paint.ArrowRuneSet) {
	theme := a.GetTheme()

	a.RLock()
	defer a.RUnlock()

	if a.arrowRuneSet != nil {
		ars = *a.arrowRuneSet
	} else {
		ars = theme.Content.ArrowRunes
	}
	return
}

func (a *CArrow) UnsetArrowRuneSet() {
	a.Lock()
	defer a.Unlock()
	a.arrowRuneSet = nil
}

func (a *CArrow) SetArrowRuneSet(ars paint.ArrowRuneSet) {
	a.Lock()
	defer a.Unlock()
	a.arrowRuneSet = &ars
}

// GetSizeRequest returns the requested size of the Drawable Widget. This method
// is used by Container Widgets to resolve the surface space allocated for their
// child Widget instances.
//
// Locking: read
func (a *CArrow) GetSizeRequest() (width, height int) {
	size := ptypes.NewRectangle(a.CWidget.GetSizeRequest())
	_, runeWidth := a.GetArrowRune()
	xPad, yPad := a.GetPadding()
	a.RLock()
	defer a.RUnlock()
	if size.W <= -1 {
		size.W = runeWidth // variable width rune supported
		size.W += xPad * 2
	}
	if size.H <= -1 {
		size.H = 1 // always one high
		size.H += yPad * 2
	}
	size.Floor(1, 1)
	return size.W, size.H
}

func (a *CArrow) resize(data []interface{}, argv ...interface{}) cenums.EventFlag {
	a.Invalidate()
	return cenums.EVENT_STOP
}

func (a *CArrow) draw(data []interface{}, argv ...interface{}) cenums.EventFlag {

	if surface, ok := argv[1].(*memphis.CSurface); ok {
		alloc := a.GetAllocation()

		if !a.IsVisible() || alloc.W <= 0 || alloc.H <= 0 {
			a.LogTrace("not visible, zero width or zero height")
			return cenums.EVENT_PASS
		}

		theme := a.GetThemeRequest()
		style := theme.Content.Normal
		r, _ := a.GetArrowRune()
		xAlign, yAlign := a.GetAlignment()
		xPad, yPad := a.GetPadding()
		size := ptypes.MakeRectangle(alloc.W-(xPad*2), alloc.H-(yPad*2))
		point := ptypes.MakePoint2I(xPad, yPad)

		surface.Fill(theme)

		if size.W < alloc.W {
			delta := alloc.W - size.W
			point.X = int(float64(delta) * xAlign)
		}
		if size.H < alloc.H {
			delta := alloc.H - size.H
			point.Y = int(float64(delta) * yAlign)
		}

		if err := surface.SetRune(point.X, point.Y, r, style); err != nil {
			a.LogError("set rune error: %v", err)
		}

		if debug, _ := a.GetBoolProperty(cdk.PropertyDebug); debug {
			surface.DebugBox(paint.ColorSilver, a.ObjectInfo())
		}

		return cenums.EVENT_STOP
	}
	return cenums.EVENT_PASS
}

// The direction the arrow should point.
// Flags: Read / Write
// Default value: GTK_ARROW_RIGHT
const PropertyArrowType cdk.Property = "arrow-type"

// Appearance of the shadow surrounding the arrow.
// Flags: Read / Write
// Default value: GTK_SHADOW_OUT
const PropertyArrowShadowType cdk.Property = "shadow-type"

const ArrowResizeHandle = "arrow-resize-handler"

const ArrowDrawHandle = "arrow-draw-handler"