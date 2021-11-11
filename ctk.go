// CTK - Curses Tool Kit
//
// The Curses Tool Kit a curses-based graphical user interface library modeled
// after the GTK project.
//
// The purpose of this project is to provide a similar API to the actual GTK
// project and instead of interacting with an X11 server, the Curses Tool Kit
// interacts with a terminal display, managed by the Curses Development Kit.
//
// CTK Type Hierarchy
//
// 	Object
// 	  |- Adjustment
// 	  `- Widget
// 	     |- Container
// 	     |  |- Bin
// 	     |  |  |- Button
// 	     |  |  |- EventBox
// 	     |  |  |- Frame
// 	     |  |  |- Viewport
// 	     |  |  |  `- ScrolledViewport
// 	     |  |  `- Window
// 	     |  `- Box
// 	     |     |- HBox
// 	     |     `- VBox
// 	     |- Misc
// 	     |  |- Arrow
// 	     |  `- Label
// 	     |- Range
// 	     |  |- Scale
// 	     |  |  |- HScale
// 	     |  |  `- VScale
// 	     |  `- Scrollbar
// 	     |     |- HScrollbar
// 	     |     `- VScrollbar
// 	     `- Sensitive
package ctk

// TODO: refactor for more parity with Gtk version
// TODO: style properties?
// TODO: Invalidate, Resize, Draw as signal handlers
// TODO: remove Theme completely and implement CSS things
// TODO: ObjectID as CSS id, name as CSS name attribute, etc
// TODO: implement a "get selector" method which finds the parent-path to object
// TODO: sensitive things can process events
// TODO: windows can LoadStyleSheetFromString()
// TODO: button focusOnClick raises issue of event processing vs focus handling
// TODO: Arrow needs to use Misc features for alignment etc
// TODO: style properties
// TODO: current theme getters
// TODO: widgets do not manipulate theme/style (no setters outside CSS)
// TODO: states like :hover :focus and so on effect current theme
// TODO: refactor enum types for BitFlags and so on, standardize convention
