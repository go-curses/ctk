package ctk

import (
	"github.com/go-curses/cdk"
	cenums "github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/cdk/lib/ptypes"
	"github.com/go-curses/ctk/lib/enums"
)

const TypeStyle cdk.CTypeTag = "ctk-style"

func init() {
	_ = cdk.TypesManager.AddType(TypeStyle, func() interface{} { return MakeStyle() })
}

// Style Hierarchy:
//	Object
//	  +- Style
type Style interface {
	Object

	PaintArrow(window Window, stateType enums.StateType, shadowType enums.ShadowType, area ptypes.Rectangle, widget Widget, detail string, arrowType enums.ArrowType, fill bool, x int, y int, width int, height int)
	PaintBox(window Window, stateType enums.StateType, shadowType enums.ShadowType, area ptypes.Rectangle, widget Widget, detail string, x int, y int, width int, height int)
	PaintBoxGap(window Window, stateType enums.StateType, shadowType enums.ShadowType, area ptypes.Rectangle, widget Widget, detail string, x int, y int, width int, height int, gapSide enums.PositionType, gapX int, gapWidth int)
	PaintCheck(window Window, stateType enums.StateType, shadowType enums.ShadowType, area ptypes.Rectangle, widget Widget, detail string, x int, y int, width int, height int)
	PaintDiamond(window Window, stateType enums.StateType, shadowType enums.ShadowType, area ptypes.Rectangle, widget Widget, detail string, x int, y int, width int, height int)
	PaintExtension(window Window, stateType enums.StateType, shadowType enums.ShadowType, area ptypes.Rectangle, widget Widget, detail string, x int, y int, width int, height int, gapSide enums.PositionType)
	PaintFlatBox(window Window, stateType enums.StateType, shadowType enums.ShadowType, area ptypes.Rectangle, widget Widget, detail string, x int, y int, width int, height int)
	PaintFocus(window Window, stateType enums.StateType, area ptypes.Rectangle, widget Widget, detail string, x int, y int, width int, height int)
	PaintHandle(window Window, stateType enums.StateType, shadowType enums.ShadowType, area ptypes.Rectangle, widget Widget, detail string, x int, y int, width int, height int, orientation cenums.Orientation)
	PaintOption(window Window, stateType enums.StateType, shadowType enums.ShadowType, area ptypes.Rectangle, widget Widget, detail string, x int, y int, width int, height int)
	PaintPolygon(window Window, stateType enums.StateType, shadowType enums.ShadowType, area ptypes.Rectangle, widget Widget, detail string, points []ptypes.Point2I, nPoints int, fill bool)
	PaintShadow(window Window, stateType enums.StateType, shadowType enums.ShadowType, area ptypes.Rectangle, widget Widget, detail string, x int, y int, width int, height int)
	PaintShadowGap(window Window, stateType enums.StateType, shadowType enums.ShadowType, area ptypes.Rectangle, widget Widget, detail string, x int, y int, width int, height int, gapSide enums.PositionType, gapX int, gapWidth int)
	PaintSlider(window Window, stateType enums.StateType, shadowType enums.ShadowType, area ptypes.Rectangle, widget Widget, detail string, x int, y int, width int, height int, orientation cenums.Orientation)
	PaintSpinner(window Window, stateType enums.StateType, area ptypes.Rectangle, widget Widget, detail string, step int, x int, y int, width int, height int)
	PaintTab(window Window, stateType enums.StateType, shadowType enums.ShadowType, area ptypes.Rectangle, widget Widget, detail string, x int, y int, width int, height int)
	PaintVLine(window Window, stateType enums.StateType, area ptypes.Rectangle, widget Widget, detail string, y1 int, y2 int, x int)
	PaintExpander(window Window, stateType enums.StateType, area ptypes.Rectangle, widget Widget, detail string, x int, y int, expanderStyle enums.ExpanderStyle)
	PaintResizeGrip(window Window, stateType enums.StateType, area ptypes.Rectangle, widget Widget, detail string, edge enums.WindowEdge, x int, y int, width int, height int)
}

var _ Style = (*CStyle)(nil)

// The CStyle structure implements the Style interface and is
// exported to facilitate type embedding with custom implementations. No member
// variables are exported as the interface methods are the only intended means
// of interacting with Style objects
type CStyle struct {
	CObject
}

// Default constructor for Style objects
func MakeStyle() Style {
	return NewStyle()
}

// Constructor for Style objects
func NewStyle() (value Style) {
	s := new(CStyle)
	s.Init()
	return s
}

// Style object initialization. This must be called at least once to setup
// the necessary defaults and allocate any memory structures. Calling this more
// than once is safe though unnecessary. Only the first call will result in any
// effect upon the Style instance
func (s *CStyle) Init() (already bool) {
	if s.InitTypeItem(TypeStyle, s) {
		return true
	}
	s.CObject.Init()
	return false
}

// Draws an arrow in the given rectangle on window using the given
// parameters. arrow_type determines the direction of the arrow.
// Parameters:
// 	window	a Window
// 	stateType	a state
// 	shadowType	the type of shadow to draw
// 	area	clip rectangle, or NULL if the
// output should not be clipped.
// 	widget	the widget.
// 	detail	a style detail.
// 	arrowType	the type of arrow to draw
// 	fill	TRUE if the arrow tip should be filled
// 	x	x origin of the rectangle to draw the arrow in
// 	y	y origin of the rectangle to draw the arrow in
// 	width	width of the rectangle to draw the arrow in
// 	height	height of the rectangle to draw the arrow in
func (s *CStyle) PaintArrow(window Window, stateType enums.StateType, shadowType enums.ShadowType, area ptypes.Rectangle, widget Widget, detail string, arrowType enums.ArrowType, fill bool, x int, y int, width int, height int) {
}

// Draws a box on window with the given parameters.
// Parameters:
// 	window	a Window
// 	stateType	a state
// 	shadowType	the type of shadow to draw
// 	area	clip rectangle, or NULL if the
// output should not be clipped.
// 	widget	the widget.
// 	detail	a style detail.
// 	x	x origin of the box
// 	y	y origin of the box
// 	width	the width of the box
// 	height	the height of the box
func (s *CStyle) PaintBox(window Window, stateType enums.StateType, shadowType enums.ShadowType, area ptypes.Rectangle, widget Widget, detail string, x int, y int, width int, height int) {
}

// Draws a box in window using the given style and state and shadow type,
// leaving a gap in one side.
// Parameters:
// 	window	a Window
// 	stateType	a state
// 	shadowType	type of shadow to draw
// 	area	clip rectangle, or NULL if the
// output should not be clipped.
// 	widget	the widget.
// 	detail	a style detail.
// 	x	x origin of the rectangle
// 	y	y origin of the rectangle
// 	width	width of the rectangle
// 	height	width of the rectangle
// 	gapSide	side in which to leave the gap
// 	gapX	starting position of the gap
// 	gapWidth	width of the gap
func (s *CStyle) PaintBoxGap(window Window, stateType enums.StateType, shadowType enums.ShadowType, area ptypes.Rectangle, widget Widget, detail string, x int, y int, width int, height int, gapSide enums.PositionType, gapX int, gapWidth int) {
}

// Draws a check button indicator in the given rectangle on window with the
// given parameters.
// Parameters:
// 	window	a Window
// 	stateType	a state
// 	shadowType	the type of shadow to draw
// 	area	clip rectangle, or NULL if the
// output should not be clipped.
// 	widget	the widget.
// 	detail	a style detail.
// 	x	x origin of the rectangle to draw the check in
// 	y	y origin of the rectangle to draw the check in
// 	width	the width of the rectangle to draw the check in
// 	height	the height of the rectangle to draw the check in
func (s *CStyle) PaintCheck(window Window, stateType enums.StateType, shadowType enums.ShadowType, area ptypes.Rectangle, widget Widget, detail string, x int, y int, width int, height int) {
}

// Draws a diamond in the given rectangle on window using the given
// parameters.
// Parameters:
// 	window	a Window
// 	stateType	a state
// 	shadowType	the type of shadow to draw
// 	area	clip rectangle, or NULL if the
// output should not be clipped.
// 	widget	the widget.
// 	detail	a style detail.
// 	x	x origin of the rectangle to draw the diamond in
// 	y	y origin of the rectangle to draw the diamond in
// 	width	width of the rectangle to draw the diamond in
// 	height	height of the rectangle to draw the diamond in
func (s *CStyle) PaintDiamond(window Window, stateType enums.StateType, shadowType enums.ShadowType, area ptypes.Rectangle, widget Widget, detail string, x int, y int, width int, height int) {
}

// Draws an extension, i.e. a notebook tab.
// Parameters:
// 	window	a Window
// 	stateType	a state
// 	shadowType	type of shadow to draw
// 	area	clip rectangle, or NULL if the
// output should not be clipped.
// 	widget	the widget.
// 	detail	a style detail.
// 	x	x origin of the extension
// 	y	y origin of the extension
// 	width	width of the extension
// 	height	width of the extension
// 	gapSide	the side on to which the extension is attached
func (s *CStyle) PaintExtension(window Window, stateType enums.StateType, shadowType enums.ShadowType, area ptypes.Rectangle, widget Widget, detail string, x int, y int, width int, height int, gapSide enums.PositionType) {
}

// Draws a flat box on window with the given parameters.
// Parameters:
// 	window	a Window
// 	stateType	a state
// 	shadowType	the type of shadow to draw
// 	area	clip rectangle, or NULL if the
// output should not be clipped.
// 	widget	the widget.
// 	detail	a style detail.
// 	x	x origin of the box
// 	y	y origin of the box
// 	width	the width of the box
// 	height	the height of the box
func (s *CStyle) PaintFlatBox(window Window, stateType enums.StateType, shadowType enums.ShadowType, area ptypes.Rectangle, widget Widget, detail string, x int, y int, width int, height int) {
}

// Draws a focus indicator around the given rectangle on window using the
// given style.
// Parameters:
// 	window	a Window
// 	stateType	a state
// 	area	clip rectangle, or NULL if the
// output should not be clipped.
// 	widget	the widget.
// 	detail	a style detail.
// 	x	the x origin of the rectangle around which to draw a focus indicator
// 	y	the y origin of the rectangle around which to draw a focus indicator
// 	width	the width of the rectangle around which to draw a focus indicator
// 	height	the height of the rectangle around which to draw a focus indicator
func (s *CStyle) PaintFocus(window Window, stateType enums.StateType, area ptypes.Rectangle, widget Widget, detail string, x int, y int, width int, height int) {
}

// Draws a handle as used in HandleBox and Paned.
// Parameters:
// 	window	a Window
// 	stateType	a state
// 	shadowType	type of shadow to draw
// 	area	clip rectangle, or NULL if the
// output should not be clipped.
// 	widget	the widget.
// 	detail	a style detail.
// 	x	x origin of the handle
// 	y	y origin of the handle
// 	width	with of the handle
// 	height	height of the handle
// 	orientation	the orientation of the handle
func (s *CStyle) PaintHandle(window Window, stateType enums.StateType, shadowType enums.ShadowType, area ptypes.Rectangle, widget Widget, detail string, x int, y int, width int, height int, orientation cenums.Orientation) {
}

// Draws a radio button indicator in the given rectangle on window with the
// given parameters.
// Parameters:
// 	window	a Window
// 	stateType	a state
// 	shadowType	the type of shadow to draw
// 	area	clip rectangle, or NULL if the
// output should not be clipped.
// 	widget	the widget.
// 	detail	a style detail.
// 	x	x origin of the rectangle to draw the option in
// 	y	y origin of the rectangle to draw the option in
// 	width	the width of the rectangle to draw the option in
// 	height	the height of the rectangle to draw the option in
func (s *CStyle) PaintOption(window Window, stateType enums.StateType, shadowType enums.ShadowType, area ptypes.Rectangle, widget Widget, detail string, x int, y int, width int, height int) {
}

// Draws a polygon on window with the given parameters.
// Parameters:
// 	window	a Window
// 	stateType	a state
// 	shadowType	type of shadow to draw
// 	area	clip rectangle, or NULL if the
// output should not be clipped.
// 	widget	the widget.
// 	detail	a style detail.
// 	points	an array of Points
// 	nPoints	length of points
//
// 	fill	TRUE if the polygon should be filled
func (s *CStyle) PaintPolygon(window Window, stateType enums.StateType, shadowType enums.ShadowType, area ptypes.Rectangle, widget Widget, detail string, points []ptypes.Point2I, nPoints int, fill bool) {
}

// Draws a shadow around the given rectangle in window using the given style
// and state and shadow type.
// Parameters:
// 	window	a Window
// 	stateType	a state
// 	shadowType	type of shadow to draw
// 	area	clip rectangle or NULL if the
// output should not be clipped.
// 	widget	the widget.
// 	detail	a style detail.
// 	x	x origin of the rectangle
// 	y	y origin of the rectangle
// 	width	width of the rectangle
// 	height	width of the rectangle
func (s *CStyle) PaintShadow(window Window, stateType enums.StateType, shadowType enums.ShadowType, area ptypes.Rectangle, widget Widget, detail string, x int, y int, width int, height int) {
}

// Draws a shadow around the given rectangle in window using the given style
// and state and shadow type, leaving a gap in one side.
// Parameters:
// 	window	a Window
// 	stateType	a state
// 	shadowType	type of shadow to draw
// 	area	clip rectangle, or NULL if the
// output should not be clipped.
// 	widget	the widget.
// 	detail	a style detail.
// 	x	x origin of the rectangle
// 	y	y origin of the rectangle
// 	width	width of the rectangle
// 	height	width of the rectangle
// 	gapSide	side in which to leave the gap
// 	gapX	starting position of the gap
// 	gapWidth	width of the gap
func (s *CStyle) PaintShadowGap(window Window, stateType enums.StateType, shadowType enums.ShadowType, area ptypes.Rectangle, widget Widget, detail string, x int, y int, width int, height int, gapSide enums.PositionType, gapX int, gapWidth int) {
}

// Draws a slider in the given rectangle on window using the given style and
// orientation.
// Parameters:
// 	window	a Window
// 	stateType	a state
// 	shadowType	a shadow
// 	area	clip rectangle, or NULL if the
// output should not be clipped.
// 	widget	the widget.
// 	detail	a style detail.
// 	x	the x origin of the rectangle in which to draw a slider
// 	y	the y origin of the rectangle in which to draw a slider
// 	width	the width of the rectangle in which to draw a slider
// 	height	the height of the rectangle in which to draw a slider
// 	orientation	the orientation to be used
func (s *CStyle) PaintSlider(window Window, stateType enums.StateType, shadowType enums.ShadowType, area ptypes.Rectangle, widget Widget, detail string, x int, y int, width int, height int, orientation cenums.Orientation) {
}

// Draws a spinner on window using the given parameters.
// Parameters:
// 	window	a Window
// 	stateType	a state
// 	area	clip rectangle, or NULL if the
// output should not be clipped.
// 	widget	the widget (may be NULL).
// 	detail	a style detail (may be NULL).
// 	step	the nth step, a value between 0 and “num-steps”
// 	x	the x origin of the rectangle in which to draw the spinner
// 	y	the y origin of the rectangle in which to draw the spinner
// 	width	the width of the rectangle in which to draw the spinner
// 	height	the height of the rectangle in which to draw the spinner
func (s *CStyle) PaintSpinner(window Window, stateType enums.StateType, area ptypes.Rectangle, widget Widget, detail string, step int, x int, y int, width int, height int) {
}

// Draws an option menu tab (i.e. the up and down pointing arrows) in the
// given rectangle on window using the given parameters.
// Parameters:
// 	window	a Window
// 	stateType	a state
// 	shadowType	the type of shadow to draw
// 	area	clip rectangle, or NULL if the
// output should not be clipped.
// 	widget	the widget.
// 	detail	a style detail.
// 	x	x origin of the rectangle to draw the tab in
// 	y	y origin of the rectangle to draw the tab in
// 	width	the width of the rectangle to draw the tab in
// 	height	the height of the rectangle to draw the tab in
func (s *CStyle) PaintTab(window Window, stateType enums.StateType, shadowType enums.ShadowType, area ptypes.Rectangle, widget Widget, detail string, x int, y int, width int, height int) {
}

// Draws a vertical line from (x , y1_ ) to (x , y2_ ) in window using the
// given style and state.
// Parameters:
// 	window	a Window
// 	stateType	a state
// 	area	rectangle to which the output is clipped, or NULL if the
// output should not be clipped.
// 	widget	the widget.
// 	detail	a style detail.
// 	y1	the starting y coordinate
// 	y2	the ending y coordinate
// 	x	the x coordinate
func (s *CStyle) PaintVLine(window Window, stateType enums.StateType, area ptypes.Rectangle, widget Widget, detail string, y1 int, y2 int, x int) {
}

// Draws an expander as used in TreeView. x and y specify the center the
// expander. The size of the expander is determined by the "expander-size"
// style property of widget . (If widget is not specified or doesn't have an
// "expander-size" property, an unspecified default size will be used, since
// the caller doesn't have sufficient information to position the expander,
// this is likely not useful.) The expander is expander_size pixels tall in
// the collapsed position and expander_size pixels wide in the expanded
// position.
// Parameters:
// 	window	a Window
// 	stateType	a state
// 	area	clip rectangle, or NULL if the
// output should not be clipped.
// 	widget	the widget.
// 	detail	a style detail.
// 	x	the x position to draw the expander at
// 	y	the y position to draw the expander at
// 	expanderStyle	the style to draw the expander in; determines
// whether the expander is collapsed, expanded, or in an
// intermediate state.
func (s *CStyle) PaintExpander(window Window, stateType enums.StateType, area ptypes.Rectangle, widget Widget, detail string, x int, y int, expanderStyle enums.ExpanderStyle) {
}

// Draws a layout on window using the given parameters.
// Parameters:
// 	window	a Window
// 	stateType	a state
// 	useText	whether to use the text or foreground
// graphics context of style
//
// 	area	clip rectangle, or NULL if the
// output should not be clipped.
// 	widget	the widget.
// 	detail	a style detail.
// 	x	x origin
// 	y	y origin
// 	layout	the layout to draw
// func (s *CStyle) PaintLayout(window Window, stateType StateType, useText bool, area ptypes.Rectangle, widget Widget, detail string, x int, y int, layout PangoLayout) {
// }

// Draws a resize grip in the given rectangle on window using the given
// parameters.
// Parameters:
// 	window	a Window
// 	stateType	a state
// 	area	clip rectangle, or NULL if the
// output should not be clipped.
// 	widget	the widget.
// 	detail	a style detail.
// 	edge	the edge in which to draw the resize grip
// 	x	the x origin of the rectangle in which to draw the resize grip
// 	y	the y origin of the rectangle in which to draw the resize grip
// 	width	the width of the rectangle in which to draw the resize grip
// 	height	the height of the rectangle in which to draw the resize grip
func (s *CStyle) PaintResizeGrip(window Window, stateType enums.StateType, area ptypes.Rectangle, widget Widget, detail string, edge enums.WindowEdge, x int, y int, width int, height int) {
}