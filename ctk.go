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