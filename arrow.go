package ctk

import (
	"unicode/utf8"

	"github.com/go-curses/cdk"
	"github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/cdk/lib/paint"
	"github.com/go-curses/cdk/lib/ptypes"
	"github.com/go-curses/cdk/memphis"
)

// CDK type-tag for Arrow objects
const TypeArrow cdk.CTypeTag = "ctk-arrow"

func init() {
	_ = cdk.TypesManager.AddType(TypeArrow, func() interface{} { return MakeArrow() })
}

// Arrow Hierarchy:
//	Object
//	  +- Widget
//	    +- Misc
//	      +- Arrow
// Arrow should be used to draw simple arrows that need to point in one of
// the four cardinal directions (up, down, left, or right). The style of the
// arrow can be one of shadow in, shadow out, etched in, or etched out. Note
// that these directions and style types may be amended in versions of CTK
// to come. Arrow will fill any space allotted to it, but since it is
// inherited from Misc, it can be padded and/or aligned, to fill exactly the
// space the programmer desires. Arrows are created with a call to
// New. The direction or style of an arrow can be changed after
// creation by using Set.
type Arrow interface {
	Misc
	Buildable

	Init() bool
	SetArrowType(arrow ArrowType)
	GetArrowType() (arrow ArrowType)
	GetArrowRune() (r rune, width int)
	GetSizeRequest() (width, height int)
}

// The CArrow structure implements the Arrow interface and is
// exported to facilitate type embedding with custom implementations. No member
// variables are exported as the interface methods are the only intended means
// of interacting with Arrow objects
type CArrow struct {
	CMisc
}

// Default constructor for Arrow objects
func MakeArrow() *CArrow {
	return NewArrow(ArrowRight)
}

// Constructor for Arrow objects
func NewArrow(arrow ArrowType) *CArrow {
	a := new(CArrow)
	a.Init()
	a.SetArrowType(arrow)
	return a
}

// Arrow object initialization. This must be called at least once to setup
// the necessary defaults and allocate any memory structures. Calling this more
// than once is safe though unnecessary. Only the first call will result in any
// effect upon the Arrow instance
func (a *CArrow) Init() bool {
	if a.InitTypeItem(TypeArrow, a) {
		return true
	}
	a.CMisc.Init()
	a.flags = NULL_WIDGET_FLAG
	a.SetFlags(PARENT_SENSITIVE)
	a.SetFlags(APP_PAINTABLE)
	_ = a.InstallBuildableProperty(PropertyArrowType, cdk.StructProperty, true, nil)
	_ = a.InstallBuildableProperty(PropertyArrowShadowType, cdk.StructProperty, true, nil)
	a.Connect(cdk.SignalDraw, ArrowDrawHandle, a.draw)
	return false
}

// Sets the direction of the Arrow.
// Parameters:
// 	arrowType	a valid ArrowType.
func (a *CArrow) SetArrowType(arrow ArrowType) {
	if err := a.SetStructProperty(PropertyArrowType, arrow); err != nil {
		a.LogErr(err)
	} else {
		a.Invalidate()
	}
}

func (a *CArrow) GetArrowType() (arrow ArrowType) {
	arrow = ArrowRight // default
	var ok bool
	if sa, err := a.GetStructProperty(PropertyArrowType); err != nil {
		a.LogErr(err)
	} else if arrow, ok = sa.(ArrowType); !ok {
		a.LogErr(err)
	}
	return
}

func (a *CArrow) GetArrowRune() (r rune, width int) {
	theme := a.GetTheme()
	arrowRunes := theme.Border.ArrowRunes
	switch a.GetArrowType() {
	case ArrowUp:
		r = arrowRunes.Up
	case ArrowLeft:
		r = arrowRunes.Left
	case ArrowDown:
		r = arrowRunes.Down
	case ArrowRight:
		r = arrowRunes.Right
	}
	width = utf8.RuneLen(r)
	return
}

func (a *CArrow) GetSizeRequest() (width, height int) {
	size := ptypes.NewRectangle(a.CWidget.GetSizeRequest())
	_, runeWidth := a.GetArrowRune()
	xPad, yPad := a.GetPadding()
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

func (a *CArrow) draw(data []interface{}, argv ...interface{}) enums.EventFlag {
	if surface, ok := argv[1].(*memphis.CSurface); ok {
		a.Lock()
		defer a.Unlock()
		alloc := a.GetAllocation()
		if !a.IsVisible() || alloc.W <= 0 || alloc.H <= 0 {
			a.LogTrace("not visible, zero width or zero height")
			return enums.EVENT_PASS
		}
		r, _ := a.GetArrowRune()
		xAlign, yAlign := a.GetAlignment()
		xPad, yPad := a.GetPadding()
		size := ptypes.MakeRectangle(alloc.W-(xPad*2), alloc.H-(yPad*2))
		point := ptypes.MakePoint2I(xPad, yPad)
		if size.W < alloc.W {
			delta := alloc.W - size.W
			point.X = int(float64(delta) * xAlign)
		}
		if size.H < alloc.H {
			delta := alloc.H - size.H
			point.Y = int(float64(delta) * yAlign)
		}
		style := a.GetThemeRequest().Content.Normal
		if err := surface.SetRune(point.X, point.Y, r, style); err != nil {
			a.LogError("set rune error: %v", err)
		}
		if debug, _ := a.GetBoolProperty(cdk.PropertyDebug); debug {
			surface.DebugBox(paint.ColorSilver, a.ObjectInfo())
		}
		return enums.EVENT_STOP
	}
	return enums.EVENT_PASS
}

// The direction the arrow should point.
// Flags: Read / Write
// Default value: GTK_ARROW_RIGHT
const PropertyArrowType cdk.Property = "arrow-type"

// Appearance of the shadow surrounding the arrow.
// Flags: Read / Write
// Default value: GTK_SHADOW_OUT
const PropertyArrowShadowType cdk.Property = "shadow-type"

const ArrowDrawHandle = "arrow-draw-handler"
