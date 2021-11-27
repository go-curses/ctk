// Package ctk is a curses-based graphical user interface library modeled
// after the GTK API.
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
// 	     |  |     `- Dialog
// 	     |  `- Box
// 	     |     |- HBox
// 	     |     |- VBox
//	     |     `- ButtonBox
//	     |        |- HButtonBox
//	     |        `- VButtonBox
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

// TODO: button focusOnClick raises issue of event processing vs focus handling
// TODO: Arrow needs to use Misc features for alignment etc
