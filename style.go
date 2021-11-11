package ctk

import (
	"github.com/go-curses/cdk"
	"github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/cdk/lib/paint"
	"github.com/go-curses/cdk/lib/ptypes"
)

// CDK type-tag for Style objects
const TypeStyle cdk.CTypeTag = "ctk-style"

func init() {
	_ = cdk.TypesManager.AddType(TypeStyle, func() interface{} { return MakeStyle() })
}

// Style Hierarchy:
//	Object
//	  +- Style
type Style interface {
	Object

	Init() (already bool)
	Copy() (value Style)
	Attach(window Window) (value Style)
	Detach()
	ApplyDefaultBackground(window Window, setBg bool, stateType StateType, area ptypes.Rectangle, x int, y int, width int, height int)
	LookupColor(colorName string, color paint.Color) (value bool)
	Get(widgetType cdk.CTypeTag, firstPropertyName string, argv ...interface{})
	PaintArrow(window Window, stateType StateType, shadowType ShadowType, area ptypes.Rectangle, widget Widget, detail string, arrowType ArrowType, fill bool, x int, y int, width int, height int)
	PaintBox(window Window, stateType StateType, shadowType ShadowType, area ptypes.Rectangle, widget Widget, detail string, x int, y int, width int, height int)
	PaintBoxGap(window Window, stateType StateType, shadowType ShadowType, area ptypes.Rectangle, widget Widget, detail string, x int, y int, width int, height int, gapSide PositionType, gapX int, gapWidth int)
	PaintCheck(window Window, stateType StateType, shadowType ShadowType, area ptypes.Rectangle, widget Widget, detail string, x int, y int, width int, height int)
	PaintDiamond(window Window, stateType StateType, shadowType ShadowType, area ptypes.Rectangle, widget Widget, detail string, x int, y int, width int, height int)
	PaintExtension(window Window, stateType StateType, shadowType ShadowType, area ptypes.Rectangle, widget Widget, detail string, x int, y int, width int, height int, gapSide PositionType)
	PaintFlatBox(window Window, stateType StateType, shadowType ShadowType, area ptypes.Rectangle, widget Widget, detail string, x int, y int, width int, height int)
	PaintFocus(window Window, stateType StateType, area ptypes.Rectangle, widget Widget, detail string, x int, y int, width int, height int)
	PaintHandle(window Window, stateType StateType, shadowType ShadowType, area ptypes.Rectangle, widget Widget, detail string, x int, y int, width int, height int, orientation enums.Orientation)
	PaintHLine(window Window, stateType StateType, area ptypes.Rectangle, widget Widget, detail string, x1 int, x2 int, y int)
	PaintOption(window Window, stateType StateType, shadowType ShadowType, area ptypes.Rectangle, widget Widget, detail string, x int, y int, width int, height int)
	PaintPolygon(window Window, stateType StateType, shadowType ShadowType, area ptypes.Rectangle, widget Widget, detail string, points ptypes.Point2I, nPoints int, fill bool)
	PaintShadow(window Window, stateType StateType, shadowType ShadowType, area ptypes.Rectangle, widget Widget, detail string, x int, y int, width int, height int)
	PaintShadowGap(window Window, stateType StateType, shadowType ShadowType, area ptypes.Rectangle, widget Widget, detail string, x int, y int, width int, height int, gapSide PositionType, gapX int, gapWidth int)
	PaintSlider(window Window, stateType StateType, shadowType ShadowType, area ptypes.Rectangle, widget Widget, detail string, x int, y int, width int, height int, orientation enums.Orientation)
	PaintSpinner(window Window, stateType StateType, area ptypes.Rectangle, widget Widget, detail string, step int, x int, y int, width int, height int)
	PaintTab(window Window, stateType StateType, shadowType ShadowType, area ptypes.Rectangle, widget Widget, detail string, x int, y int, width int, height int)
	PaintVLine(window Window, stateType StateType, area ptypes.Rectangle, widget Widget, detail string, y1 int, y2 int, x int)
	PaintExpander(window Window, stateType StateType, area ptypes.Rectangle, widget Widget, detail string, x int, y int, expanderStyle ExpanderStyle)
	// PaintLayout(window Window, stateType StateType, useText bool, area ptypes.Rectangle, widget Widget, detail string, x int, y int, layout PangoLayout)
	PaintResizeGrip(window Window, stateType StateType, area ptypes.Rectangle, widget Widget, detail string, edge WindowEdge, x int, y int, width int, height int)
	BorderNew() (value paint.Border)
	BorderCopy(border paint.Border) (value paint.Border)
	BorderFree(border paint.Border)
}

// The CStyle structure implements the Style interface and is
// exported to facilitate type embedding with custom implementations. No member
// variables are exported as the interface methods are the only intended means
// of interacting with Style objects
type CStyle struct {
	CObject
}

// Default constructor for Style objects
func MakeStyle() *CStyle {
	return NewStyle()
}

// Constructor for Style objects
func NewStyle() (value *CStyle) {
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

// Creates a copy of the passed in Style object.
// Returns:
// 	a copy of style .
// 	[transfer full]
func (s *CStyle) Copy() (value Style) {
	return nil
}

// Attaches a style to a window; this process allocates the colors and
// creates the GC's for the style - it specializes it to a particular visual
// and colormap. The process may involve the creation of a new style if the
// style has already been attached to a window with a different style and
// colormap. Since this function may return a new object, you have to use it
// in the following way: style = Attach (style, window)
// Parameters:
// 	window	a Window.
// Returns:
// 	Either style , or a newly-created Style. If the style is
// 	newly created, the style parameter will be unref'ed, and the
// 	new style will have a reference count belonging to the caller.
func (s *CStyle) Attach(window Window) (value Style) {
	return nil
}

// Detaches a style from a window. If the style is not attached to any
// windows anymore, it is unrealized. See Attach.
func (s *CStyle) Detach() {}

//
// Parameters:
// 	area	.
func (s *CStyle) ApplyDefaultBackground(window Window, setBg bool, stateType StateType, area ptypes.Rectangle, x int, y int, width int, height int) {
}

// Looks up color_name in the style's logical color mappings, filling in
// color and returning TRUE if found, otherwise returning FALSE. Do not cache
// the found mapping, because it depends on the Style and might change
// when a theme switch occurs.
// Parameters:
// 	colorName	the name of the logical color to look up
// 	color	the Color to fill in.
// Returns:
// 	TRUE if the mapping was found.
func (s *CStyle) LookupColor(colorName string, color paint.Color) (value bool) {
	return false
}

// Looks up stock_id in the icon factories associated with style and the
// default icon factory, returning an icon set if found, otherwise NULL.
// Parameters:
// 	stockId	an icon name
// Returns:
// 	icon set of stock_id .
// 	[transfer none]
// func (s *CStyle) LookupIconSet(stockId string) (value IconSet) {
// 	return nil
// }

// Renders the icon specified by source at the given size according to the
// given parameters and returns the result in a pixbuf.
// Parameters:
// 	source	the IconSource specifying the icon to render
// 	direction	a text direction
// 	state	a state
// 	size	(type int) the size to render the icon at. A size of
// (IconSize)-1 means render at the size of the source and
// don't scale.
// 	widget	the widget.
// 	detail	a style detail.
// Returns:
// 	a newly-created Pixbuf containing the rendered icon.
// 	[transfer full]
// func (s *CStyle) RenderIcon(source IconSource, direction TextDirection, state StateType, size IconSize, widget Widget, detail string) (value Pixbuf) {
// 	return nil
// }

// Queries the value of a style property corresponding to a widget class is
// in the given style.
// Parameters:
// 	widgetType	the GType of a descendant of Widget
// 	propertyName	the name of the style property to get
// 	value	a GValue where the value of the property being
// queried will be stored
func (s *CStyle) GetStyleProperty(widgetType cdk.CTypeTag, propertyName string) (value interface{}) {
	return
}

// Non-vararg variant of Get. Used primarily by language
// bindings.
// Parameters:
// 	widgetType	the GType of a descendant of Widget
// 	firstPropertyName	the name of the first style property to get
// 	varArgs	a va_list of pairs of property names and
// locations to return the property values, starting with the
// location for first_property_name
// .
// func (s *CStyle) GetValist(widgetType GType, firstPropertyName string, varArgs va_list) {}

// Gets the values of a multiple style properties for widget_type from style
// .
// Parameters:
// 	widgetType	the GType of a descendant of Widget
// 	firstPropertyName	the name of the first style property to get
// 	varargs	pairs of property names and locations to
// return the property values, starting with the location for
// first_property_name
// , terminated by NULL.
// func (s *CStyle) Get(widgetType GType, firstPropertyName string, argv ...interface{}) {}

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
func (s *CStyle) PaintArrow(window Window, stateType StateType, shadowType ShadowType, area ptypes.Rectangle, widget Widget, detail string, arrowType ArrowType, fill bool, x int, y int, width int, height int) {
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
func (s *CStyle) PaintBox(window Window, stateType StateType, shadowType ShadowType, area ptypes.Rectangle, widget Widget, detail string, x int, y int, width int, height int) {
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
func (s *CStyle) PaintBoxGap(window Window, stateType StateType, shadowType ShadowType, area ptypes.Rectangle, widget Widget, detail string, x int, y int, width int, height int, gapSide PositionType, gapX int, gapWidth int) {
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
func (s *CStyle) PaintCheck(window Window, stateType StateType, shadowType ShadowType, area ptypes.Rectangle, widget Widget, detail string, x int, y int, width int, height int) {
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
func (s *CStyle) PaintDiamond(window Window, stateType StateType, shadowType ShadowType, area ptypes.Rectangle, widget Widget, detail string, x int, y int, width int, height int) {
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
func (s *CStyle) PaintExtension(window Window, stateType StateType, shadowType ShadowType, area ptypes.Rectangle, widget Widget, detail string, x int, y int, width int, height int, gapSide PositionType) {
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
func (s *CStyle) PaintFlatBox(window Window, stateType StateType, shadowType ShadowType, area ptypes.Rectangle, widget Widget, detail string, x int, y int, width int, height int) {
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
func (s *CStyle) PaintFocus(window Window, stateType StateType, area ptypes.Rectangle, widget Widget, detail string, x int, y int, width int, height int) {
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
func (s *CStyle) PaintHandle(window Window, stateType StateType, shadowType ShadowType, area ptypes.Rectangle, widget Widget, detail string, x int, y int, width int, height int, orientation enums.Orientation) {
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
func (s *CStyle) PaintOption(window Window, stateType StateType, shadowType ShadowType, area ptypes.Rectangle, widget Widget, detail string, x int, y int, width int, height int) {
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
func (s *CStyle) PaintPolygon(window Window, stateType StateType, shadowType ShadowType, area ptypes.Rectangle, widget Widget, detail string, points []ptypes.Point2I, nPoints int, fill bool) {
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
func (s *CStyle) PaintShadow(window Window, stateType StateType, shadowType ShadowType, area ptypes.Rectangle, widget Widget, detail string, x int, y int, width int, height int) {
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
func (s *CStyle) PaintShadowGap(window Window, stateType StateType, shadowType ShadowType, area ptypes.Rectangle, widget Widget, detail string, x int, y int, width int, height int, gapSide PositionType, gapX int, gapWidth int) {
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
func (s *CStyle) PaintSlider(window Window, stateType StateType, shadowType ShadowType, area ptypes.Rectangle, widget Widget, detail string, x int, y int, width int, height int, orientation enums.Orientation) {
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
func (s *CStyle) PaintSpinner(window Window, stateType StateType, area ptypes.Rectangle, widget Widget, detail string, step int, x int, y int, width int, height int) {
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
func (s *CStyle) PaintTab(window Window, stateType StateType, shadowType ShadowType, area ptypes.Rectangle, widget Widget, detail string, x int, y int, width int, height int) {
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
func (s *CStyle) PaintVLine(window Window, stateType StateType, area ptypes.Rectangle, widget Widget, detail string, y1 int, y2 int, x int) {
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
func (s *CStyle) PaintExpander(window Window, stateType StateType, area ptypes.Rectangle, widget Widget, detail string, x int, y int, expanderStyle ExpanderStyle) {
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
func (s *CStyle) PaintResizeGrip(window Window, stateType StateType, area ptypes.Rectangle, widget Widget, detail string, edge WindowEdge, x int, y int, width int, height int) {
}

// // Allocates a new Border structure and initializes its elements to zero.
// // Returns:
// // 	a new empty Border. The newly allocated Border should be
// // 	freed with BorderFree
// func (s *CStyle) BorderNew() (value Border) {
// 	return nil
// }
//
// // Copies a Border structure.
// // Parameters:
// // 	border	a Border.
// // 	returns	a copy of border_
// // .
// func (s *CStyle) BorderCopy(border paint.Border) (value paint.Border) {
// 	return nil
// }
//
// // Frees a Border structure.
// // Parameters:
// // 	border	a Border.
// func (s *CStyle) BorderFree(border paint.Border) {}
//
// // func (s *CStyle) GtkRcPropertyParser(pspec GParamSpec, rcString GString, propertyValue GValue) (value bool) {
// // 	return false
// // }

// Emitted when the style has been initialized for a particular colormap and
// depth. Connecting to this signal is probably seldom useful since most of
// the time applications and widgets only deal with styles that have been
// already realized.
const SignalStyleRealize cdk.Signal = "realize"

// Emitted when the aspects of the style specific to a particular colormap
// and depth are being cleaned up. A connection to this signal can be useful
// if a widget wants to cache objects like a GC as object data on
// Style. This signal provides a convenient place to free such cached
// objects.
const SignalStyleUnrealize cdk.Signal = "unrealize"
