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
	_ "embed"
	"fmt"

	"github.com/go-curses/cdk"
	cenums "github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/cdk/lib/paint"
	"github.com/go-curses/cdk/lib/ptypes"
	"github.com/go-curses/cdk/lib/sync"
	"github.com/go-curses/cdk/memphis"

	"github.com/go-curses/ctk/lib/enums"
)

const TypeWindow cdk.CTypeTag = "ctk-window"

func init() {
	_ = cdk.TypesManager.AddType(TypeWindow, func() interface{} { return MakeWindow() })
	ctkBuilderTranslators[TypeWindow] = func(builder Builder, widget Widget, name, value string) error {
		switch name {
		case "transient_for", "transient-for":
			if tfw := builder.GetWidget(value); tfw != nil {
				if err := widget.SetStructProperty(PropertyTransientFor, tfw); err != nil {
					return err
				}
				return nil
			} else {
				builder.LogError("failed to set transient-for, unknown widget id/name: %v", value)
			}
		}
		return ErrFallthrough
	}
}

//go:embed ctk.default.styles
var DefaultStyles string

// Window Hierarchy:
//
//	Object
//	  +- Widget
//	    +- Container
//	      +- Bin
//	        +- Window
//	          +- Dialog
//	          +- Assistant
//	          +- OffscreenWindow
//	          +- Plug
//
// In the Curses Tool Kit, the Window type is an extension of the CTK Bin
// type and also implements the cdk.Window interface so that it can be utilized
// within the Curses Development Kit framework. A Window is a TOPLEVEL Widget
// that can contain other widgets.
type Window interface {
	cdk.Window
	Bin

	ImportStylesFromString(css string) (err error)
	ReplaceStylesFromString(css string) (err error)
	ExportStylesToString() (css string)
	ApplyStylesTo(widget Widget)
	ReApplyStyles()
	SetTitle(title string)
	SetResizable(resizable bool)
	GetResizable() (value bool)
	AddAccelGroup(accelGroup AccelGroup)
	RemoveAccelGroup(accelGroup AccelGroup)
	ActivateFocus() (value bool)
	ActivateDefault() (value bool)
	SetModal(modal bool)
	SetPosition(position enums.WindowPosition)
	SetTransientFor(parent Window)
	SetDestroyWithParent(setting bool)
	IsActive() (active bool)
	HasToplevelFocus() (focused bool)
	ListTopLevels() (value []Window)
	AddMnemonic(keyval rune, target Widget)
	RemoveMnemonic(keyval rune, target Widget)
	RemoveWidgetMnemonics(target Widget)
	ActivateMnemonic(keyval rune, modifier cdk.ModMask) (activated bool)
	ActivateKey(event cdk.EventKey) (value bool)
	PropagateKeyEvent(event cdk.EventKey) (value bool)
	GetFocus() (focus Widget)
	SetFocus(focus Widget)
	GetDefaultWidget() (value Widget)
	SetDefault(defaultWidget Widget)
	Present()
	PresentWithTime(timestamp int)
	Iconify()
	Deiconify()
	Stick()
	Unstick()
	Maximize()
	Unmaximize()
	Fullscreen()
	Unfullscreen()
	SetKeepAbove(setting bool)
	SetKeepBelow(setting bool)
	SetDecorated(setting bool)
	SetDeletable(setting bool)
	SetMnemonicModifier(modifier cdk.ModMask)
	SetWindowType(hint cenums.WindowType)
	SetTypeHint(hint enums.WindowTypeHint)
	SetSkipTaskbarHint(setting bool)
	SetSkipPagerHint(setting bool)
	SetUrgencyHint(setting bool)
	SetAcceptFocus(setting bool)
	SetFocusOnMap(setting bool)
	SetStartupId(startupId string)
	SetRole(role string)
	GetDecorated() (value bool)
	GetDeletable() (value bool)
	GetDefaultSize(width int, height int)
	GetDestroyWithParent() (value bool)
	GetMnemonicModifier() (value cdk.ModMask)
	GetModal() (value bool)
	GetPosition(rootX int, rootY int)
	GetRole() (value string)
	GetSize() (width, height int)
	GetTitle() (value string)
	GetTransientFor() (value Window)
	GetWindowType() (value cenums.WindowType)
	GetTypeHint() (value enums.WindowTypeHint)
	GetSkipTaskbarHint() (value bool)
	GetSkipPagerHint() (value bool)
	GetUrgencyHint() (value bool)
	GetAcceptFocus() (value bool)
	GetFocusOnMap() (value bool)
	HasGroup() (value bool)
	Move(x int, y int)
	ParseGeometry(geometry string) (value bool)
	ReshowWithInitialSize()
	SetAutoStartupNotification(setting bool)
	GetOpacity() (value float64)
	SetOpacity(opacity float64)
	GetMnemonicsVisible() (value bool)
	SetMnemonicsVisible(setting bool)
	GetDisplay() (dm cdk.Display)
	SetDisplay(dm cdk.Display)
	RequestDrawAndShow()
	RequestDrawAndSync()
	GetVBox() (vbox VBox)
	GetNextFocus() (next Widget)
	GetPreviousFocus() (previous Widget)
	FocusNext() cenums.EventFlag
	FocusPrevious() cenums.EventFlag
	GetEventFocus() (o cdk.Object)
	SetEventFocus(o cdk.Object)
}

var _ Window = (*CWindow)(nil)

// The CWindow structure implements the Window interface and is exported to
// facilitate type embedding with custom implementations. No member variables
// are exported as the interface methods are the only intended means of
// interacting with Window objects.
type CWindow struct {
	CBin

	display cdk.Display

	prevMouseEvent *cdk.EventMouse
	eventFocus     interface{}
	hoverFocus     Widget
	hoverFocused   *WidgetSlice
	accelGroups    AccelGroups
	mnemonics      []*mnemonicEntry
	mnemonicMod    cdk.ModMask
	mnemonicLock   *sync.RWMutex
	receivingPaste bool
	pasteBuffer    *string

	styleSheet *cStyleSheet
}

type mnemonicEntry struct {
	key    rune
	target Widget
}

// MakeWindow is used by the Buildable system to construct a new Window.
func MakeWindow() Window {
	return NewWindow()
}

// NewWindow is a constructor for new Window instances.
func NewWindow() (w Window) {
	w = new(CWindow)
	w.Init()
	return
}

// NewWindowWithTitle is a constructor for new Window instances that also sets
// the Window title to the string given.
func NewWindowWithTitle(title string) (w Window) {
	w = NewWindow()
	w.SetTitle(title)
	return
}

// Init initializes a Window object. This must be called at least once to
// set up the necessary defaults and allocate any memory structures. Calling
// this more than once is safe though unnecessary. Only the first call will
// result in any effect upon the Window instance. Init is used in the
// NewWindow constructor and only necessary when implementing a derivative
// Window type.
func (w *CWindow) Init() (already bool) {
	if w.InitTypeItem(TypeWindow, w) {
		return true
	}
	w.CBin.Init()
	w.SetState(0)
	w.flags = enums.NULL_WIDGET_FLAG
	w.SetFlags(enums.TOPLEVEL | enums.SENSITIVE | enums.APP_PAINTABLE)
	w.prevMouseEvent = cdk.NewEventMouse(0, 0, 0, 0)
	w.parent = w
	if d := cdk.GetDefaultDisplay(); d != nil {
		w.display = d
	} else {
		w.display = nil
	}
	w.origin.X = 0
	w.origin.Y = 0
	w.accelGroups = NewAccelGroups()
	w.mnemonics = make([]*mnemonicEntry, 0)
	w.mnemonicMod = cdk.ModAlt
	w.mnemonicLock = &sync.RWMutex{}
	w.hoverFocus = nil
	w.hoverFocused = new(WidgetSlice)
	w.receivingPaste = false
	w.pasteBuffer = nil

	_ = w.InstallProperty(PropertyWindowType, cdk.StructProperty, true, cenums.WINDOW_TOPLEVEL)
	_ = w.InstallProperty(PropertyAcceptFocus, cdk.BoolProperty, true, true)
	_ = w.InstallProperty(PropertyDecorated, cdk.BoolProperty, true, true)
	_ = w.InstallProperty(PropertyDefaultHeight, cdk.IntProperty, true, -1)
	_ = w.InstallProperty(PropertyDefaultWidth, cdk.IntProperty, true, -1)
	_ = w.InstallProperty(PropertyDeletable, cdk.BoolProperty, true, true)
	_ = w.InstallProperty(PropertyDestroyWithParent, cdk.BoolProperty, true, false)
	_ = w.InstallProperty(PropertyFocusOnMap, cdk.BoolProperty, true, true)
	_ = w.InstallProperty(PropertyFocusedWidget, cdk.StructProperty, true, nil)
	_ = w.InstallProperty(PropertyGravity, cdk.StructProperty, true, enums.GravityNorthWest)
	_ = w.InstallProperty(PropertyHasToplevelFocus, cdk.BoolProperty, false, false)
	_ = w.InstallProperty(PropertyIcon, cdk.StructProperty, true, nil)
	_ = w.InstallProperty(PropertyIconName, cdk.StringProperty, true, "")
	_ = w.InstallProperty(PropertyIsActive, cdk.BoolProperty, false, false)
	_ = w.InstallProperty(PropertyMnemonicsVisible, cdk.BoolProperty, true, true)
	_ = w.InstallProperty(PropertyModal, cdk.BoolProperty, true, false)
	_ = w.InstallProperty(PropertyOpacity, cdk.FloatProperty, true, 1)
	_ = w.InstallProperty(PropertyResizable, cdk.BoolProperty, true, true)
	_ = w.InstallProperty(PropertyRole, cdk.StringProperty, true, "")
	_ = w.InstallProperty(PropertyScreen, cdk.StructProperty, true, nil)
	_ = w.InstallProperty(PropertySkipPagerHint, cdk.BoolProperty, true, false)
	_ = w.InstallProperty(PropertySkipTaskbarHint, cdk.BoolProperty, true, false)
	_ = w.InstallProperty(PropertyStartupId, cdk.StringProperty, true, "")
	_ = w.InstallProperty(PropertyTitle, cdk.StringProperty, true, "")
	_ = w.InstallProperty(PropertyTransientFor, cdk.StructProperty, true, nil)
	_ = w.InstallProperty(PropertyTypeHint, cdk.StructProperty, true, enums.WindowTypeHintNormal)
	_ = w.InstallProperty(PropertyUrgencyHint, cdk.BoolProperty, true, false)
	_ = w.InstallProperty(PropertyWindowPosition, cdk.StructProperty, true, enums.WinPosNone)

	var err error
	if w.styleSheet, err = newStyleSheetFromString(DefaultStyles); err != nil {
		w.LogErr(err)
	} else {
		w.styleSheet = newStyleSheet()
	}
	w.SetWindow(w) // after stylesheet setup

	w.Connect(SignalCdkEvent, WindowEventHandle, w.event)
	w.Connect(SignalResize, WindowResizeHandle, w.resize)
	w.Connect(SignalDraw, WindowDrawHandle, w.draw)

	if err := w.SetProperty(PropertyWindow, w); err != nil {
		w.LogErr(err)
	}
	_ = w.GetVBox()

	if w.display != nil {
		w.display.Connect(cdk.SignalFocusedWindow, WindowDisplayFocusHandle, w.displayFocusedWindow)
	}

	return false
}

// Build provides customizations to the Buildable system for Window Widgets.
func (w *CWindow) Build(builder Builder, element *CBuilderElement) error {
	w.Freeze()
	defer w.Thaw()
	if err := w.CBin.Build(builder, element); err != nil {
		return err
	}
	if len(element.Children) > 0 {
		contentBox := w.GetVBox()
		for _, child := range element.Children {
			if newChild := builder.Build(child); newChild != nil {
				child.Instance = newChild
				if newChildWidget, ok := newChild.(Widget); ok {
					newChildWidget.Show()
					// if len(child.Packing) > 0 {
					expand, fill, padding, packType := builder.ParsePacking(child)
					if packType == enums.PackStart {
						contentBox.PackStart(newChildWidget, expand, fill, padding)
					} else {
						contentBox.PackEnd(newChildWidget, expand, fill, padding)
					}
					if newChildWidget.HasFlags(enums.HAS_FOCUS) {
						newChildWidget.GrabFocus()
					}
				} else {
					contentBox.LogError("new child object is not a Widget type: %v (%T)")
				}
			}
		}
	}
	return nil
}

func (w *CWindow) ImportStylesFromString(css string) (err error) {
	w.Lock()
	if err = w.styleSheet.ParseString(css); err != nil {
		w.LogErr(err)
	}
	w.Unlock()
	return
}

func (w *CWindow) ReplaceStylesFromString(css string) (err error) {
	w.Lock()
	var ss *cStyleSheet
	if ss, err = newStyleSheetFromString(DefaultStyles); err != nil {
		w.LogErr(err)
	} else {
		w.styleSheet = ss
	}
	w.Unlock()
	return
}

func (w *CWindow) ExportStylesToString() (css string) {
	css = w.styleSheet.String()
	return
}

func (w *CWindow) ApplyStylesTo(widget Widget) {
	w.styleSheet.ApplyStylesTo(widget)
}

func (w *CWindow) ReApplyStyles() {
	var recurse func(top Widget)
	recurse = func(top Widget) {
		w.styleSheet.ApplyStylesTo(top)
		if container, ok := top.Self().(Container); ok {
			for _, child := range container.GetChildren() {
				recurse(child)
			}
		}
	}
	recurse(w)
}

func (w *CWindow) Show() {
	w.CBin.Show()
	w.GetVBox().Show()
	if display := w.GetDisplay(); display != nil {
		w.Map()
		display.MapWindowWithRegion(w, w.GetRegion())
	}
	w.SetState(enums.StatePrelight)
	w.InvalidateAllChildren()
}

func (w *CWindow) ShowAll() {
	w.GetVBox().ShowAll()
	w.Show()
}

func (w *CWindow) Hide() {
	w.CBin.Hide()
	if display := w.GetDisplay(); display != nil {
		// w.Unmap()
		display.UnmapWindow(w)
	}
}

// SetTitle updates the title of the Window. The title of a window will be
// displayed in its title bar; on the X Window System, the title bar is rendered
// by the window manager, so exactly how the title appears to users may vary
// according to a user's exact configuration. The title should help a user
// distinguish this window from other windows they may have open. A good
// title might include the application name and current document filename,
// for example.
//
// Parameters:
//
//	title	text for the title of the window
func (w *CWindow) SetTitle(title string) {
	if err := w.SetStringProperty(PropertyTitle, title); err != nil {
		w.LogErr(err)
	}
}

// SetResizable updates whether the user can resize a window. Windows are user
// resizable by default.
//
// Parameters:
//
//	resizable	TRUE if the user can resize this window
func (w *CWindow) SetResizable(resizable bool) {
	if err := w.SetBoolProperty(PropertyResizable, resizable); err != nil {
		w.LogErr(err)
	}
}

// GetResizable returns the value set by SetResizable.
func (w *CWindow) GetResizable() (value bool) {
	var err error
	if value, err = w.GetBoolProperty(PropertyResizable); err != nil {
		w.LogErr(err)
	}
	return
}

// AddAccelGroup associates accel_group with window, such that calling
// AccelGroupsActivate on the window will activate accelerators in
// accel_group.
//
// Parameters:
//
//	accelGroup	a AccelGroup
func (w *CWindow) AddAccelGroup(accelGroup AccelGroup) {
	w.Lock()
	w.accelGroups.AddAccelGroup(w, accelGroup)
	w.Unlock()
}

// RemoveAccelGroup reverses the effects of AddAccelGroup.
//
// Parameters:
//
//	accelGroup	a AccelGroup
func (w *CWindow) RemoveAccelGroup(accelGroup AccelGroup) {
	w.Lock()
	w.accelGroups.RemoveAccelGroup(w, accelGroup)
	w.Unlock()
}

// Activates the current focused widget within the window.
// Returns:
//
//	TRUE if a widget got activated.
func (w *CWindow) ActivateFocus() (value bool) {
	if focused := w.GetFocus(); focused != nil {
		if activatable, ok := focused.Self().(Activatable); ok {
			activatable.Activate()
			return true
		}
	}
	return
}

// Activates the default widget for the window, unless the current focused
// widget has been configured to receive the default action (see
// WidgetSetReceivesDefault), in which case the focused widget is
// activated.
// Returns:
//
//	TRUE if a widget got activated.
func (w *CWindow) ActivateDefault() (value bool) {
	return false
}

// Sets a window modal or non-modal. Modal windows prevent interaction with
// other windows in the same application. To keep modal dialogs on top of
// main application windows, use SetTransientFor to make the
// dialog transient for the parent; most window managers will then disallow
// lowering the dialog below the parent.
// Parameters:
//
//	modal	whether the window is modal
func (w *CWindow) SetModal(modal bool) {
	if err := w.SetBoolProperty(PropertyModal, modal); err != nil {
		w.LogErr(err)
	}
}

// This function sets up hints about how a window can be resized by the user.
// You can set a minimum and maximum size; allowed resize increments (e.g.
// for xterm, you can only resize by the size of a character); aspect ratios;
// and more. See the Geometry struct.
// Parameters:
// 	geometryWidget	widget the geometry hints will be applied to
// 	geometry	struct containing geometry information
// 	geomMask	mask indicating which struct fields should be paid attention to
// func (w *CWindow) SetGeometryHints(geometryWidget Widget, geometry Geometry, geomMask WindowHints) {}

// Window gravity defines the meaning of coordinates passed to
// Move. See Move and Gravity for more details.
// The default window gravity is GDK_GRAVITY_NORTH_WEST which will typically
// "do what you mean."
// Parameters:
// 	gravity	window gravity
// func (w *CWindow) SetGravity(gravity Gravity) {
// 	if err := w.SetStructProperty(PropertyGravity, gravity); err != nil {
// 		w.LogErr(err)
// 	}
// }

// Gets the value set by SetGravity.
// Returns:
// 	window gravity.
// 	[transfer none]
// func (w *CWindow) GetGravity() (value Gravity) {
// 	var err error
// 	if value, err = w.GetStructProperty(PropertyGravity); err != nil {
// 		w.LogErr(err)
// 	}
// 	return
// }

// Sets a position constraint for this window. If the old or new constraint
// is GTK_WIN_POS_CENTER_ALWAYS, this will also cause the window to be
// repositioned to satisfy the new constraint.
// Parameters:
//
//	position	a position constraint.
func (w *CWindow) SetPosition(position enums.WindowPosition) {}

// Dialog windows should be set transient for the main application window
// they were spawned from. This allows window managers to e.g. keep the
// dialog on top of the main window, or center the dialog over the main
// window. DialogNewWithButtons and other convenience functions in
// CTK will sometimes call SetTransientFor on your behalf.
// Passing NULL for parent unsets the current transient window. On Windows,
// this function puts the child window on top of the parent, much as the
// window manager would have done on X.
// Parameters:
//
//	parent	parent window, or NULL.
func (w *CWindow) SetTransientFor(parent Window) {
	if err := w.SetStructProperty(PropertyTransientFor, parent); err != nil {
		w.LogErr(err)
	}
}

// If setting is TRUE, then destroying the transient parent of window will
// also destroy window itself. This is useful for dialogs that shouldn't
// persist beyond the lifetime of the main window they're associated with,
// for example.
// Parameters:
//
//	setting	whether to destroy window
//
// with its transient parent
func (w *CWindow) SetDestroyWithParent(setting bool) {
	if err := w.SetBoolProperty(PropertyDestroyWithParent, setting); err != nil {
		w.LogErr(err)
	}
}

// // Sets the Screen where the window is displayed; if the window is already
// // mapped, it will be unmapped, and then remapped on the new screen.
// // Parameters:
// // 	screen	a Screen.
// func (w *CWindow) SetScreen(screen Screen) {
// 	if err := w.SetStructProperty(PropertyScreen, screen); err != nil {
// 		w.LogErr(err)
// 	}
// }
//
// // Returns the Screen associated with window .
// // Returns:
// // 	a Screen.
// // 	[transfer none]
// func (w *CWindow) GetScreen() (value Screen) {
// 	var err error
// 	if value, err = w.GetStructProperty(PropertyScreen); err != nil {
// 		w.LogErr(err)
// 	}
// 	return
// }

// Returns whether the window is part of the current active toplevel. (That
// is, the toplevel window receiving keystrokes.) The return value is TRUE if
// the window is active toplevel itself, but also if it is, say, a Plug
// embedded in the active toplevel. You might use this function if you wanted
// to draw a widget differently in an active window from a widget in an
// inactive window. See HasToplevelFocus
// Returns:
//
//	TRUE if the window part of the current active window.
func (w *CWindow) IsActive() (active bool) {
	if focused := w.display.FocusedWindow(); focused != nil {
		if focused.ObjectID() == w.ObjectID() {
			return true
		}
	}
	return false
}

// Returns whether the input focus is within this Window. For real
// toplevel windows, this is identical to IsActive, but for
// embedded windows, like Plug, the results will differ.
// Returns:
//
//	TRUE if the input focus is within this Window
func (w *CWindow) HasToplevelFocus() (focused bool) {
	return
}

// Returns a list of all existing toplevel windows. The widgets in the list
// are not individually referenced. If you want to iterate through the list
// and perform actions involving callbacks that might destroy the widgets,
// you must call g_list_foreach (result, (GFunc)g_object_ref, NULL) first,
// and then unref all the widgets afterwards.
// Returns:
//
//	list of toplevel widgets.
//	[element-type Widget][transfer container]
func (w *CWindow) ListTopLevels() (value []Window) {
	w.LogError("method unimplemented")
	return []Window{w}
}

// Adds a mnemonic to this window.
// Parameters:
//
//	keyval	the mnemonic
//	target	the widget that gets activated by the mnemonic
func (w *CWindow) AddMnemonic(keyval rune, target Widget) {
	w.mnemonicLock.Lock()
	for _, entry := range w.mnemonics {
		if entry.key == keyval {
			if widget, ok := entry.target.Self().(Sensitive); ok {
				if tw, ok := target.Self().(Sensitive); ok {
					if widget.ObjectID() == tw.ObjectID() {
						w.mnemonicLock.Unlock()
						return
					}
				}
			}
		}
	}
	if _, ok := target.Self().(Sensitive); ok {
		w.mnemonics = append(w.mnemonics, &mnemonicEntry{
			key:    keyval,
			target: target,
		})
	} else {
		w.LogError("target is not a Widget: %v (%T)", target, target)
	}
	w.mnemonicLock.Unlock()
}

// Removes a mnemonic from this window.
// Parameters:
//
//	keyval	the mnemonic
//	target	the widget that gets activated by the mnemonic
func (w *CWindow) RemoveMnemonic(keyval rune, target Widget) {
	w.mnemonicLock.Lock()
	if tw, ok := target.Self().(Sensitive); ok {
		var mnemonics []*mnemonicEntry
		for _, entry := range w.mnemonics {
			if entry.key == keyval {
				if widget, ok := entry.target.Self().(Sensitive); ok {
					if widget.ObjectID() != tw.ObjectID() {
						mnemonics = append(mnemonics, entry)
					}
				}
			}
		}
		w.mnemonics = mnemonics
	} else {
		w.LogError("target is not a Widget: %v (%T)", target, target)
	}
	w.mnemonicLock.Unlock()
}

// Removes all mnemonics from this window for the target Widget.
// Parameters:
//
//	target	the widget that gets activated by the mnemonic
func (w *CWindow) RemoveWidgetMnemonics(target Widget) {
	w.mnemonicLock.Lock()
	if tw, ok := target.Self().(Widget); ok {
		var mnemonics []*mnemonicEntry
		for _, entry := range w.mnemonics {
			if widget, ok := entry.target.Self().(Widget); ok {
				if widget.ObjectID() != tw.ObjectID() {
					mnemonics = append(mnemonics, entry)
				}
			}
		}
		w.mnemonics = mnemonics
	} else {
		w.LogError("target is not a Widget: %v (%T)", target, target)
	}
	w.mnemonicLock.Unlock()
}

// Activates the targets associated with the mnemonic.
// Parameters:
//
//	keyval	the mnemonic
//	modifier	the modifiers
//	returns	TRUE if the activation is done.
func (w *CWindow) ActivateMnemonic(keyval rune, modifier cdk.ModMask) (activated bool) {
	if !GetDefaultSettings().GetEnableMnemonics() {
		return false
	}
	w.mnemonicLock.Lock()
	if modifier == w.mnemonicMod {
		for _, entry := range w.mnemonics {
			if entry.key == keyval {
				if sa, ok := entry.target.Self().(Sensitive); ok && sa.IsSensitive() && sa.IsVisible() {
					w.mnemonicLock.Unlock()
					sa.GrabFocus()
					sa.Activate()
					return true
				}
			}
		}
	}
	w.mnemonicLock.Unlock()
	return
}

// Activates mnemonics and accelerators for this Window. This is normally
// called by the default ::key_press_event handler for toplevel windows,
// however in some cases it may be useful to call this directly when
// overriding the standard key handling for a toplevel window.
// Parameters:
//
//	event	a EventKey
//
// Returns:
//
//	TRUE if a mnemonic or accelerator was found and activated.
func (w *CWindow) ActivateKey(event cdk.EventKey) (value bool) {
	return false
}

// Propagate a key press or release event to the focus widget and up the
// focus container chain until a widget handles event . This is normally
// called by the default ::key_press_event and ::key_release_event handlers
// for toplevel windows, however in some cases it may be useful to call this
// directly when overriding the standard key handling for a toplevel window.
// Parameters:
//
//	event	a EventKey
//
// Returns:
//
//	TRUE if a widget in the focus chain handled the event.
func (w *CWindow) PropagateKeyEvent(event cdk.EventKey) (value bool) {
	return false
}

// Retrieves the current focused widget within the window. Note that this is
// the widget that would have the focus if the toplevel window focused; if
// the toplevel window is not focused then WidgetHasFocus (widget) will
// not be TRUE for the widget.
// Returns:
//
//	the currently focused widget, or NULL if there is none.
//	[transfer none]
func (w *CWindow) GetFocus() (focus Widget) {
	var err error
	var ok bool
	var focusStruct interface{}
	if focusStruct, err = w.GetStructProperty(PropertyFocusedWidget); err != nil {
		w.LogErr(err)
	} else if focus, ok = focusStruct.(Widget); !ok && focusStruct != nil {
		w.LogError("value stored in %v property is not of Widget type: %v (%T)", PropertyFocusedWidget, focusStruct, focusStruct)
	} else if ok {
		return
	}
	fc, _ := w.GetFocusChain()
	if len(fc) > 0 {
		w.SetFocus(fc[0])
		focus = fc[0]
	}
	return
}

// If focus is not the current focus widget, and is focusable, sets it as the
// focus widget for the window. If focus is NULL, unsets the focus widget for
// this window. To set the focus to a particular widget in the toplevel, it
// is usually more convenient to use WidgetGrabFocus instead of this
// function.
//
// Parameters:
//
//	focus	widget to be the new focus widget, or NULL to unset
func (w *CWindow) SetFocus(focus Widget) {
	// if transient := w.GetTransientFor(); transient != nil && w.ObjectID() != transient.ObjectID() {
	// 	transient.SetFocus(focus)
	// 	return
	// }

	if focus == nil {
		if err := w.SetStructProperty(PropertyFocusedWidget, nil); err != nil {
			w.LogErr(err)
		} else {
			w.Emit(SignalFocusChanged, nil)
		}
		return
	}

	if sensitive, ok := focus.Self().(Sensitive); ok {
		if sensitive.CanFocus() && sensitive.IsVisible() && sensitive.IsSensitive() {
			if err := w.SetStructProperty(PropertyFocusedWidget, focus); err != nil {
				w.LogErr(err)
			} else {
				w.Emit(SignalFocusChanged, w, focus)
			}
		} else {
			w.LogError("cannot focus, not visible or not sensitive")
		}
	} else {
		w.LogError("does not implement Sensitive interface: %T %v", focus, focus)
	}
}

// Returns the default widget for window . See SetDefault for
// more details.
// Returns:
//
//	the default widget, or NULL if there is none.
//	[transfer none]
func (w *CWindow) GetDefaultWidget() (value Widget) {
	return nil
}

// The default widget is the widget that's activated when the user presses
// Enter in a dialog (for example). This function sets or unsets the default
// widget for a Window about. When setting (rather than unsetting) the
// default widget it's generally easier to call WidgetGrabFocus on
// the widget. Before making a widget the default widget, you must set the
// GTK_CAN_DEFAULT flag on the widget you'd like to make the default using
// GTK_WIDGET_SET_FLAGS.
// Parameters:
//
//	defaultWidget	widget to be the default, or NULL to unset the
//
// default widget for the toplevel.
func (w *CWindow) SetDefault(defaultWidget Widget) {}

// Presents a window to the user. This may mean raising the window in the
// stacking order, deiconifying it, moving it to the current desktop, and/or
// giving it the keyboard focus, possibly dependent on the user's platform,
// window manager, and preferences. If window is hidden, this function calls
// WidgetShow as well. This function should be used when the user
// tries to open a window that's already open. Say for example the
// preferences dialog is currently open, and the user chooses Preferences
// from the menu a second time; use Present to move the
// already-open dialog where the user can see it. If you are calling this
// function in response to a user interaction, it is preferable to use
// PresentWithTime.
func (w *CWindow) Present() {}

// Presents a window to the user in response to a user interaction. If you
// need to present a window without a timestamp, use Present.
// See Present for details.
// Parameters:
//
//	timestamp	the timestamp of the user interaction (typically a
//
// button or key press event) which triggered this call
func (w *CWindow) PresentWithTime(timestamp int) {}

// Asks to iconify (i.e. minimize) the specified window . Note that you
// shouldn't assume the window is definitely iconified afterward, because
// other entities (e.g. the user or window manager) could deiconify it again,
// or there may not be a window manager in which case iconification isn't
// possible, etc. But normally the window will end up iconified. Just don't
// write code that crashes if not. It's permitted to call this function
// before showing a window, in which case the window will be iconified before
// it ever appears onscreen. You can track iconification via the
// "window-state-event" signal on Widget.
func (w *CWindow) Iconify() {}

// Asks to deiconify (i.e. unminimize) the specified window . Note that you
// shouldn't assume the window is definitely deiconified afterward, because
// other entities (e.g. the user or window manager) could iconify it again
// before your code which assumes deiconification gets to run. You can track
// iconification via the "window-state-event" signal on Widget.
func (w *CWindow) Deiconify() {}

// Asks to stick window , which means that it will appear on all user
// desktops. Note that you shouldn't assume the window is definitely stuck
// afterward, because other entities (e.g. the user or window manager) could
// unstick it again, and some window managers do not support sticking
// windows. But normally the window will end up stuck. Just don't write code
// that crashes if not. It's permitted to call this function before showing a
// window. You can track stickiness via the "window-state-event" signal on
// Widget.
func (w *CWindow) Stick() {}

// Asks to unstick window , which means that it will appear on only one of
// the user's desktops. Note that you shouldn't assume the window is
// definitely unstuck afterward, because other entities (e.g. the user or
// window manager) could stick it again. But normally the window will end up
// stuck. Just don't write code that crashes if not. You can track stickiness
// via the "window-state-event" signal on Widget.
func (w *CWindow) Unstick() {}

// Asks to maximize window , so that it becomes full-screen. Note that you
// shouldn't assume the window is definitely maximized afterward, because
// other entities (e.g. the user or window manager) could unmaximize it
// again, and not all window managers support maximization. But normally the
// window will end up maximized. Just don't write code that crashes if not.
// It's permitted to call this function before showing a window, in which
// case the window will be maximized when it appears onscreen initially. You
// can track maximization via the "window-state-event" signal on Widget.
func (w *CWindow) Maximize() {}

// Asks to unmaximize window . Note that you shouldn't assume the window is
// definitely unmaximized afterward, because other entities (e.g. the user or
// window manager) could maximize it again, and not all window managers honor
// requests to unmaximize. But normally the window will end up unmaximized.
// Just don't write code that crashes if not. You can track maximization via
// the "window-state-event" signal on Widget.
func (w *CWindow) Unmaximize() {}

// Asks to place window in the fullscreen state. Note that you shouldn't
// assume the window is definitely full screen afterward, because other
// entities (e.g. the user or window manager) could unfullscreen it again,
// and not all window managers honor requests to fullscreen windows. But
// normally the window will end up fullscreen. Just don't write code that
// crashes if not. You can track the fullscreen state via the
// "window-state-event" signal on Widget.
func (w *CWindow) Fullscreen() {}

// Asks to toggle off the fullscreen state for window . Note that you
// shouldn't assume the window is definitely not full screen afterward,
// because other entities (e.g. the user or window manager) could fullscreen
// it again, and not all window managers honor requests to unfullscreen
// windows. But normally the window will end up restored to its normal state.
// Just don't write code that crashes if not. You can track the fullscreen
// state via the "window-state-event" signal on Widget.
func (w *CWindow) Unfullscreen() {}

// Asks to keep window above, so that it stays on top. Note that you
// shouldn't assume the window is definitely above afterward, because other
// entities (e.g. the user or window manager) could not keep it above, and
// not all window managers support keeping windows above. But normally the
// window will end kept above. Just don't write code that crashes if not.
// It's permitted to call this function before showing a window, in which
// case the window will be kept above when it appears onscreen initially. You
// can track the above state via the "window-state-event" signal on
// Widget. Note that, according to the Extended Window Manager Hints
// specification, the above state is mainly meant for user preferences and
// should not be used by applications e.g. for drawing attention to their
// dialogs.
// Parameters:
//
//	setting	whether to keep window
//
// above other windows
func (w *CWindow) SetKeepAbove(setting bool) {}

// Asks to keep window below, so that it stays in bottom. Note that you
// shouldn't assume the window is definitely below afterward, because other
// entities (e.g. the user or window manager) could not keep it below, and
// not all window managers support putting windows below. But normally the
// window will be kept below. Just don't write code that crashes if not. It's
// permitted to call this function before showing a window, in which case the
// window will be kept below when it appears onscreen initially. You can
// track the below state via the "window-state-event" signal on Widget.
// Note that, according to the Extended Window Manager Hints specification,
// the above state is mainly meant for user preferences and should not be
// used by applications e.g. for drawing attention to their dialogs.
// Parameters:
//
//	setting	whether to keep window
//
// below other windows
func (w *CWindow) SetKeepBelow(setting bool) {}

// Starts resizing a window. This function is used if an application has
// window resizing controls. When GDK can support it, the resize will be done
// using the standard mechanism for the window manager or windowing system.
// Otherwise, GDK will try to emulate window resizing, potentially not all
// that well, depending on the windowing system.
// Parameters:
// 	button	mouse button that initiated the drag
// 	edge	position of the resize control
// 	rootX	X position where the user clicked to initiate the drag, in root window coordinates
// 	rootY	Y position where the user clicked to initiate the drag
// 	timestamp	timestamp from the click event that initiated the drag
// func (w *CWindow) BeginResizeDrag(edge WindowEdge, button int, rootX int, rootY int, timestamp int) {}

// Starts moving a window. This function is used if an application has window
// movement grips. When GDK can support it, the window movement will be done
// using the standard mechanism for the window manager or windowing system.
// Otherwise, GDK will try to emulate window movement, potentially not all
// that well, depending on the windowing system.
// Parameters:
// 	button	mouse button that initiated the drag
// 	rootX	X position where the user clicked to initiate the drag, in root window coordinates
// 	rootY	Y position where the user clicked to initiate the drag
// 	timestamp	timestamp from the click event that initiated the drag
// func (w *CWindow) BeginMoveDrag(button int, rootX int, rootY int, timestamp int) {}

// By default, windows are decorated with a title bar, resize controls, etc.
// Some window managers allow CTK to disable these decorations, creating a
// borderless window. If you set the decorated property to FALSE using this
// function, CTK will do its best to convince the window manager not to
// decorate the window. Depending on the system, this function may not have
// any effect when called on a window that is already visible, so you should
// call it before calling Show. On Windows, this function always
// works, since there's no window manager policy involved.
// Parameters:
//
//	setting	TRUE to decorate the window
func (w *CWindow) SetDecorated(setting bool) {
	if err := w.SetBoolProperty(PropertyDecorated, setting); err != nil {
		w.LogErr(err)
	}
}

// By default, windows have a close button in the window frame. Some disable
// this button. If you set the deletable property to FALSE using this
// function, CTK will do its best to convince the window manager not to show
// a close button. Depending on the system, this function may not have any
// effect when called on a window that is already visible, so you should call
// it before calling Show. On Windows, this function always
// works, since there's no window manager policy involved.
// Parameters:
//
//	setting	TRUE to decorate the window as deletable
func (w *CWindow) SetDeletable(setting bool) {
	if err := w.SetBoolProperty(PropertyDeletable, setting); err != nil {
		w.LogErr(err)
	}
}

// Sets the mnemonic modifier for this window.
// Parameters:
//
//	modifier	the modifier mask used to activate
//
// mnemonics on this window.
func (w *CWindow) SetMnemonicModifier(modifier cdk.ModMask) {
	w.mnemonicLock.Lock()
	w.mnemonicMod = modifier
	w.mnemonicLock.Unlock()
}

// SetWindowType updates the type of the window.
// See: enums.WindowType
func (w *CWindow) SetWindowType(hint cenums.WindowType) {
	if err := w.SetStructProperty(PropertyWindowType, hint); err != nil {
		w.LogErr(err)
	}
}

// By setting the type hint for the window, you allow the window manager to
// decorate and handle the window in a way which is suitable to the function
// of the window in your application. This function should be called before
// the window becomes visible. DialogNewWithButtons and other
// convenience functions in CTK will sometimes call
// SetTypeHint on your behalf.
// Parameters:
//
//	hint	the window type
func (w *CWindow) SetTypeHint(hint enums.WindowTypeHint) {
	if err := w.SetStructProperty(PropertyTypeHint, hint); err != nil {
		w.LogErr(err)
	}
}

// Windows may set a hint asking the desktop environment not to display the
// window in the task bar. This function sets this hint.
// Parameters:
//
//	setting	TRUE to keep this window from appearing in the task bar
func (w *CWindow) SetSkipTaskbarHint(setting bool) {
	if err := w.SetBoolProperty(PropertySkipTaskbarHint, setting); err != nil {
		w.LogErr(err)
	}
}

// Windows may set a hint asking the desktop environment not to display the
// window in the pager. This function sets this hint. (A "pager" is any
// desktop navigation tool such as a workspace switcher that displays a
// thumbnail representation of the windows on the screen.)
// Parameters:
//
//	setting	TRUE to keep this window from appearing in the pager
func (w *CWindow) SetSkipPagerHint(setting bool) {
	if err := w.SetBoolProperty(PropertySkipPagerHint, setting); err != nil {
		w.LogErr(err)
	}
}

// Windows may set a hint asking the desktop environment to draw the users
// attention to the window. This function sets this hint.
// Parameters:
//
//	setting	TRUE to mark this window as urgent
func (w *CWindow) SetUrgencyHint(setting bool) {
	if err := w.SetBoolProperty(PropertyUrgencyHint, setting); err != nil {
		w.LogErr(err)
	}
}

// Windows may set a hint asking the desktop environment not to receive the
// input focus. This function sets this hint.
// Parameters:
//
//	setting	TRUE to let this window receive input focus
func (w *CWindow) SetAcceptFocus(setting bool) {
	if err := w.SetBoolProperty(PropertyAcceptFocus, setting); err != nil {
		w.LogErr(err)
	}
}

// Windows may set a hint asking the desktop environment not to receive the
// input focus when the window is mapped. This function sets this hint.
// Parameters:
//
//	setting	TRUE to let this window receive input focus on map
func (w *CWindow) SetFocusOnMap(setting bool) {
	if err := w.SetBoolProperty(PropertyFocusOnMap, setting); err != nil {
		w.LogErr(err)
	}
}

// Startup notification identifiers are used by desktop environment to track
// application startup, to provide user feedback and other features. This
// function changes the corresponding property on the underlying Window.
// Normally, startup identifier is managed automatically and you should only
// use this function in special cases like transferring focus from other
// processes. You should use this function before calling
// Present or any equivalent function generating a window map
// event. This function is only useful on X11, not with other CTK targets.
// Parameters:
//
//	startupId	a string with startup-notification identifier
func (w *CWindow) SetStartupId(startupId string) {
	if err := w.SetStringProperty(PropertyStartupId, startupId); err != nil {
		w.LogErr(err)
	}
}

// This function is only useful on X11, not with other CTK targets. In
// combination with the window title, the window role allows a same" window
// when an application is restarted. So for example you might set the
// "toolbox" role on your app's toolbox window, so that when the user
// restarts their session, the window manager can put the toolbox back in the
// same place. If a window already has a unique title, you don't need to set
// the role, since the WM can use the title to identify the window when
// restoring the session.
// Parameters:
//
//	role	unique identifier for the window to be used when restoring a session
func (w *CWindow) SetRole(role string) {
	if err := w.SetStringProperty(PropertyRole, role); err != nil {
		w.LogErr(err)
	}
}

// Returns whether the window has been set to have decorations such as a
// title bar via SetDecorated.
// Returns:
//
//	TRUE if the window has been set to have decorations
func (w *CWindow) GetDecorated() (value bool) {
	var err error
	if value, err = w.GetBoolProperty(PropertyDecorated); err != nil {
		w.LogErr(err)
	}
	return
}

// Returns whether the window has been set to have a close button via
// SetDeletable.
// Returns:
//
//	TRUE if the window has been set to have a close button
func (w *CWindow) GetDeletable() (value bool) {
	var err error
	if value, err = w.GetBoolProperty(PropertyDeletable); err != nil {
		w.LogErr(err)
	}
	return
}

// Gets the default size of the window. A value of -1 for the width or height
// indicates that a default size has not been explicitly set for that
// dimension, so the "natural" size of the window will be used.
// Parameters:
//
//	width	location to store the default width, or NULL.
//	height	location to store the default height, or NULL.
func (w *CWindow) GetDefaultSize(width int, height int) {}

// Returns whether the window will be destroyed with its transient parent.
// See SetDestroyWithParent.
// Returns:
//
//	TRUE if the window will be destroyed with its transient parent.
func (w *CWindow) GetDestroyWithParent() (value bool) {
	var err error
	if value, err = w.GetBoolProperty(PropertyDestroyWithParent); err != nil {
		w.LogErr(err)
	}
	return
}

// Returns the mnemonic modifier for this window. See
// SetMnemonicModifier.
// Returns:
//
//	the modifier mask used to activate mnemonics on this window.
func (w *CWindow) GetMnemonicModifier() (value cdk.ModMask) {
	w.mnemonicLock.RLock()
	value = w.mnemonicMod
	w.mnemonicLock.RUnlock()
	return
}

// Returns whether the window is modal. See SetModal.
// Returns:
//
//	TRUE if the window is set to be modal and establishes a grab
//	when shown
func (w *CWindow) GetModal() (value bool) {
	var err error
	if value, err = w.GetBoolProperty(PropertyModal); err != nil {
		w.LogErr(err)
	}
	return
}

// This function returns the position you need to pass to Move
// to keep window in its current position. This means that the meaning of the
// returned value varies with window gravity. See Move for more
// details. If you haven't changed the window gravity, its gravity will be
// GDK_GRAVITY_NORTH_WEST. This means that GetPosition gets the
// position of the top-left corner of the window manager frame for the
// window. Move sets the position of this same top-left corner.
// GetPosition is not 100% reliable because the X Window System
// does not specify a way to obtain the geometry of the decorations placed on
// a window by the window manager. Thus CTK is using a "best guess" that
// works with most window managers. Moreover, nearly all window managers are
// historically broken with respect to their handling of window gravity. So
// moving a window to its current position as returned by
// GetPosition tends to result in moving the window slightly.
// Window managers are slowly getting better over time. If a window has
// gravity GDK_GRAVITY_STATIC the window manager frame is not relevant, and
// thus GetPosition will always produce accurate results.
// However you can't use static gravity to do things like place a window in a
// corner of the screen, because static gravity ignores the window manager
// decorations. If you are saving and restoring your application's window
// positions, you should know that it's impossible for applications to do
// this without getting it somewhat wrong because applications do not have
// sufficient knowledge of window manager state. The Correct Mechanism is to
// support the session management protocol (see the "GnomeClient" object in
// the GNOME libraries for example) and allow the window manager to save your
// window sizes and positions.
// Parameters:
//
//	rootX	return location for X coordinate of gravity-determined reference point.
//	rootY	return location for Y coordinate of gravity-determined reference point.
func (w *CWindow) GetPosition(rootX int, rootY int) {}

// Returns the role of the window. See SetRole for further
// explanation.
// Returns:
//
//	the role of the window if set, or NULL. The returned is owned
//	by the widget and must not be modified or freed.
func (w *CWindow) GetRole() (value string) {
	var err error
	if value, err = w.GetStringProperty(PropertyRole); err != nil {
		w.LogErr(err)
	}
	return
}

// Obtains the current size of window . If window is not onscreen, it returns
// the size CTK will suggest to the window manager for the initial window
// size (but this is not reliably the same as the size the window manager
// will actually select). The size obtained by GetSize is the
// last size received in a EventConfigure, that is, CTK uses its
// locally-stored size, rather than querying the X server for the size. As a
// result, if you call Resize then immediately call
// GetSize, the size won't have taken effect yet. After the
// window manager processes the resize request, CTK receives notification
// that the size has changed via a configure event, and the size of the
// window gets updated. Note 1: Nearly any use of this function creates a
// race condition, because the size of the window may change between the time
// that you get the size and the time that you perform some action assuming
// that size is the current size. To avoid race conditions, connect to
// "configure-event" on the window and adjust your size-dependent state to
// match the size delivered in the EventConfigure. Note 2: The returned
// size does not include the size of the window manager decorations (aka the
// window frame or border). Those are not drawn by CTK and CTK has no
// reliable method of determining their size. Note 3: If you are getting a
// window size in order to position the window onscreen, there may be a
// better way. The preferred way is to simply set the window's semantic type
// with SetTypeHint, which allows the window manager to e.g.
// center dialogs. Also, if you set the transient parent of dialogs with
// SetTransientFor window managers will often center the
// dialog over its parent window. It's much preferred to let the window
// manager handle these things rather than doing it yourself, because all
// apps will behave consistently and according to user prefs if the window
// manager handles it. Also, the window manager can take the size of the
// window decorations/border into account, while your application cannot. In
// any case, if you insist on application-specified window positioning,
// there's still a better way than doing it yourself -
// SetPosition will frequently handle the details for you.
// Parameters:
//
//	width	return location for width, or NULL.
//	height	return location for height, or NULL.
func (w *CWindow) GetSize() (width, height int) {
	alloc := w.GetAllocation()
	width, height = alloc.W, alloc.H
	return
}

// Retrieves the title of the window. See SetTitle.
// Returns:
//
//	the title of the window, or NULL if none has been set
//	explicitely. The returned string is owned by the widget and
//	must not be modified or freed.
func (w *CWindow) GetTitle() (value string) {
	var err error
	if value, err = w.GetStringProperty(PropertyTitle); err != nil {
		w.LogErr(err)
	}
	return
}

// Fetches the transient parent for this window. See
// SetTransientFor.
// Returns:
//
//	the transient parent for this window, or NULL if no transient
//	parent has been set.
//	[transfer none]
func (w *CWindow) GetTransientFor() (value Window) {
	var ok bool
	if v, err := w.GetStructProperty(PropertyTransientFor); err != nil {
		w.LogErr(err)
	} else if value, ok = v.(Window); !ok && v != nil {
		w.LogError("value stored in %v is not of Window type: %v (%T)", PropertyTransientFor, v, v)
	}
	return
}

// GetWindowType returns the type of the window.
// See: enums.WindowType.
func (w *CWindow) GetWindowType() (value cenums.WindowType) {
	var ok bool
	if v, err := w.GetStructProperty(PropertyWindowType); err != nil {
		w.LogErr(err)
	} else if value, ok = v.(cenums.WindowType); !ok {
		value = cenums.WINDOW_TOPLEVEL // default is top-level?
		w.LogError("value stored in %v is not of cenums.WindowType: %v (%T)", PropertyWindowType, v, v)
	}
	return
}

// Gets the type hint for this window. See SetTypeHint.
// Returns:
//
//	the type hint for window .
func (w *CWindow) GetTypeHint() (value enums.WindowTypeHint) {
	var ok bool
	if v, err := w.GetStructProperty(PropertyTypeHint); err != nil {
		w.LogErr(err)
	} else if value, ok = v.(enums.WindowTypeHint); !ok {
		w.LogError("value stored in %v is not of enums.WindowTypeHint: %v (%T)", PropertyTypeHint, v, v)
	}
	return
}

// Gets the value set by SetSkipTaskbarHint
// Returns:
//
//	TRUE if window shouldn't be in taskbar
func (w *CWindow) GetSkipTaskbarHint() (value bool) {
	var err error
	if value, err = w.GetBoolProperty(PropertySkipTaskbarHint); err != nil {
		w.LogErr(err)
	}
	return
}

// Gets the value set by SetSkipPagerHint.
// Returns:
//
//	TRUE if window shouldn't be in pager
func (w *CWindow) GetSkipPagerHint() (value bool) {
	var err error
	if value, err = w.GetBoolProperty(PropertySkipPagerHint); err != nil {
		w.LogErr(err)
	}
	return
}

// Gets the value set by SetUrgencyHint
// Returns:
//
//	TRUE if window is urgent
func (w *CWindow) GetUrgencyHint() (value bool) {
	var err error
	if value, err = w.GetBoolProperty(PropertyUrgencyHint); err != nil {
		w.LogErr(err)
	}
	return
}

// Gets the value set by SetAcceptFocus.
// Returns:
//
//	TRUE if window should receive the input focus
func (w *CWindow) GetAcceptFocus() (value bool) {
	var err error
	if value, err = w.GetBoolProperty(PropertyAcceptFocus); err != nil {
		w.LogErr(err)
	}
	return
}

// Gets the value set by SetFocusOnMap.
// Returns:
//
//	TRUE if window should receive the input focus when mapped.
func (w *CWindow) GetFocusOnMap() (value bool) {
	var err error
	if value, err = w.GetBoolProperty(PropertyFocusOnMap); err != nil {
		w.LogErr(err)
	}
	return
}

// Returns the group for window or the default group, if window is NULL or if
// window does not have an explicit window group.
// Returns:
// 	the WindowGroup for a window or the default group.
// 	[transfer none]
// func (w *CWindow) GetGroup() (value WindowGroup) {
// 	return nil
// }

// Returns whether window has an explicit window group.
// Returns:
//
//	TRUE if window has an explicit window group.
//	Since 2.22
func (w *CWindow) HasGroup() (value bool) {
	return false
}

// Asks the window manager to move window to the given position. Window
// managers are free to ignore this; most window managers ignore requests for
// initial window positions (instead using a user-defined placement
// algorithm) and honor requests after the window has already been shown.
// Note the position is the position of the gravity-determined reference
// point for the window. The gravity determines two things: first, the
// location of the reference point in root window coordinates; and second,
// which point on the window is positioned at the reference point. By default
// the gravity is GDK_GRAVITY_NORTH_WEST, so the reference point is simply
// the x , y supplied to Move. The top-left corner of the window
// decorations (aka window frame or border) will be placed at x , y .
// Therefore, to position a window at the top left of the screen, you want to
// use the default gravity (which is GDK_GRAVITY_NORTH_WEST) and move the
// window to 0,0. To position a window at the bottom right corner of the
// screen, you would set GDK_GRAVITY_SOUTH_EAST, which means that the
// reference point is at x + the window width and y + the window height, and
// the bottom-right corner of the window border will be placed at that
// reference point. So, to place a window in the bottom right corner you
// would first set gravity to south east, then write: Move
// (window, ScreenWidth - window_width, ScreenHeight -
// window_height) (note that this example does not take multi-head scenarios
// into account). The Extended Window Manager Hints specification at
// http://www.freedesktop.org/Standards/wm-spec has a nice table of gravities
// in the "implementation notes" section. The GetPosition
// documentation may also be relevant.
// Parameters:
//
//	x	X coordinate to move window to
//	y	Y coordinate to move window to
func (w *CWindow) Move(x int, y int) {
	// update the origin and surface configuration?
	w.SetOrigin(x, y)
	if w.IsMapped() {
		if surface, err := memphis.GetSurface(w.ObjectID()); err != nil {
			w.LogErr(err)
		} else {
			surface.SetOrigin(w.GetOrigin())
		}
	}
	w.Invalidate()
}

// Parses a standard X Window System geometry string - see the manual page
// for X (type 'man X') for details on this. ParseGeometry does
// work on all CTK ports including Win32 but is primarily intended for an X
// environment. If either a size or a position can be extracted from the
// geometry string, ParseGeometry returns TRUE and calls
// SetDefaultSize and/or Move to resize/move the
// window. If ParseGeometry returns TRUE, it will also set the
// GDK_HINT_USER_POS and/or GDK_HINT_USER_SIZE hints indicating to the window
// manager that the size/position of the window was user-specified. This
// causes most window managers to honor the geometry. Note that for
// ParseGeometry to work as expected, it has to be called when
// the window has its "final" size, i.e. after calling WidgetShowAll
// on the contents and SetGeometryHints on the window.
// Parameters:
//
//	geometry	geometry string
//
// Returns:
//
//	TRUE if string was parsed successfully
func (w *CWindow) ParseGeometry(geometry string) (value bool) {
	return false
}

// Hides window , then reshows it, resetting the default size and position of
// the window. Used by GUI builders only.
func (w *CWindow) ReshowWithInitialSize() {}

// By default, after showing the first Window, CTK calls
// NotifyStartupComplete. Call this function to disable the automatic
// startup notification. You might do this if your first window is a splash
// screen, and you want to delay notification until after your real main
// window has been shown, for example. In that example, you would disable
// startup notification temporarily, show your splash screen, then re-enable
// it so that showing the main window would automatically result in
// notification.
// Parameters:
//
//	setting	TRUE to automatically do startup notification
func (w *CWindow) SetAutoStartupNotification(setting bool) {}

// Fetches the requested opacity for this window. See
// SetOpacity.
// Returns:
//
//	the requested opacity for this window.
func (w *CWindow) GetOpacity() (value float64) {
	var err error
	if value, err = w.GetFloatProperty(PropertyOpacity); err != nil {
		w.LogErr(err)
	}
	return
}

// Request the windowing system to make window partially transparent, with
// opacity 0 being fully transparent and 1 fully opaque. (Values of the
// opacity parameter are clamped to the [0,1] range.) On X11 this has any
// effect only on X screens with a compositing manager running. See
// WidgetIsComposited. On Windows it should work always. Note that
// setting a window's opacity after the window has been shown causes it to
// flicker once on Windows.
// Parameters:
//
//	opacity	desired opacity, between 0 and 1
func (w *CWindow) SetOpacity(opacity float64) {
	if err := w.SetFloatProperty(PropertyOpacity, opacity); err != nil {
		w.LogErr(err)
	}
}

func (w *CWindow) GetMnemonicsVisible() (value bool) {
	var err error
	if value, err = w.GetBoolProperty(PropertyMnemonicsVisible); err != nil {
		w.LogErr(err)
	}
	return
}

// Sets the mnemonics-visible property.
// Parameters:
//
//	setting	the new value
func (w *CWindow) SetMnemonicsVisible(setting bool) {
	if err := w.SetBoolProperty(PropertyMnemonicsVisible, setting); err != nil {
		w.LogErr(err)
	}
}

func (w *CWindow) GetDisplay() (dm cdk.Display) {
	w.RLock()
	defer w.RUnlock()
	return w.display
}

func (w *CWindow) SetDisplay(dm cdk.Display) {
	w.Lock()
	w.display = dm
	w.Unlock()
}

func (w *CWindow) RequestDrawAndShow() {
	if d := w.GetDisplay(); d != nil {
		// w.Invalidate()
		d.RequestDraw()
		d.RequestShow()
	}
}

func (w *CWindow) RequestDrawAndSync() {
	if d := w.GetDisplay(); d != nil {
		// w.Invalidate()
		d.RequestDraw()
		d.RequestSync()
	}
}

func (w *CWindow) GetVBox() (vbox VBox) {
	// bin child must be an internal VBox
	if child := w.GetChild(); child != nil {
		var ok bool
		if vbox, ok = child.Self().(VBox); ok {
			vbox.SetFlags(enums.VISIBLE)
			return // exists, use it
		}
		w.LogError("removing internal widget: %v (%T) - not a VBox", child, child)
		w.Remove(child)
	}
	// new VBox required, either first run or erroneous child cleared
	vbox = NewVBox(false, 0)
	vbox.Show()
	w.Add(vbox)
	return
}

func (w *CWindow) GetNextFocus() (next Widget) {
	fc, _ := w.CBin.GetFocusChain()
	if focused := w.GetFocus(); focused != nil {
		w.RLock()
		found := false
		for _, fci := range fc {
			if !found {
				if fci.ObjectID() == focused.ObjectID() {
					found = true
				}
			} else {
				next = fci
				w.RUnlock()
				return
			}
		}
		if len(fc) > 0 {
			next = fc[0]
		}
		w.RUnlock()
	} else if len(fc) > 0 {
		next = fc[0]
	}
	return
}

func (w *CWindow) GetPreviousFocus() (previous Widget) {
	fc, _ := w.CBin.GetFocusChain()
	nfc := len(fc)
	if focused := w.GetFocus(); focused != nil {
		w.RLock()
		found := false
		for _, fci := range fc {
			if !found && fci.ObjectID() == focused.ObjectID() {
				found = true
				break
			}
			previous = fci
		}
		if previous == nil && nfc > 0 {
			previous = fc[nfc-1]
		}
		w.RUnlock()
	} else if nfc > 0 {
		previous = fc[nfc-1]
	}
	return
}

func (w *CWindow) FocusNext() cenums.EventFlag {
	if next := w.GetNextFocus(); next != nil {
		next.GrabFocus()
		return cenums.EVENT_STOP
	}
	w.LogError("no widgets to focus next")
	return cenums.EVENT_PASS
}

func (w *CWindow) FocusPrevious() cenums.EventFlag {
	if prev := w.GetPreviousFocus(); prev != nil {
		prev.GrabFocus()
		return cenums.EVENT_STOP
	}
	w.LogError("no widgets to focus previous")
	return cenums.EVENT_PASS
}

func (w *CWindow) GetEventFocus() (o cdk.Object) {
	if dm := w.GetDisplay(); dm != nil {
		o = dm.GetEventFocus()
	}
	return
}

func (w *CWindow) SetEventFocus(o cdk.Object) {
	if f := w.Emit(SignalSetEventFocus, w, o); f == cenums.EVENT_PASS {
		if dm := w.GetDisplay(); dm != nil {
			if err := dm.SetEventFocus(o); err != nil {
				w.LogError("error setting event focus: %v", err)
			}
		}
	}
}

func (w *CWindow) Draw() cenums.EventFlag {
	return w.CWidget.Draw()
}

func (w *CWindow) event(data []interface{}, argv ...interface{}) cenums.EventFlag {
	if evt, ok := argv[1].(cdk.Event); ok {
		switch e := evt.(type) {

		case *cdk.EventError:
			if f := w.Emit(SignalError, w, e); f == cenums.EVENT_PASS {
				w.LogError(e.Error())
			}

		case *cdk.EventPaste:
			starting := e.Start()
			if starting {
				w.Lock()
				w.pasteBuffer = new(string)
				w.receivingPaste = true
				w.Unlock()
			} else {
				w.Lock()
				if w.display != nil {
					clipboard := w.display.GetClipboard()
					clipboard.Paste(*w.pasteBuffer)
				}
				w.receivingPaste = false
				w.pasteBuffer = nil
				w.Unlock()
				if fi := w.GetFocus(); fi != nil {
					if sw, ok := fi.Self().(Sensitive); ok && sw.IsSensitive() && sw.IsVisible() {
						if f := sw.ProcessEvent(evt); f == cenums.EVENT_STOP {
							return cenums.EVENT_STOP
						}
					}
				}
			}
			return cenums.EVENT_STOP

		case *cdk.EventKey:
			w.RLock()
			if w.receivingPaste {
				w.RUnlock()
				w.Lock()
				*w.pasteBuffer += string(e.Rune())
				w.Unlock()
				return cenums.EVENT_STOP
			}
			w.RUnlock()

			if f := w.Emit(SignalEventKey, w, e); f == cenums.EVENT_PASS {
				// check for mnemonics
				if w.ActivateMnemonic(e.Rune(), e.Modifiers()) {
					return cenums.EVENT_STOP
				}
				// check for accelerators
				if w.accelGroups.Activate(w, e.Key(), e.Modifiers()) {
					return cenums.EVENT_STOP
				}

				// check focus change
				switch {
				case e.RuneAsKey() == cdk.KeyBacktab || e.Key() == cdk.KeyBacktab:
					if e.Modifiers().Has(cdk.ModShift) {
						w.LogDebug("shift+back-tab key caught")
						w.FocusNext()
					} else {
						w.LogDebug("back-tab key caught")
						w.FocusPrevious()
					}
					return cenums.EVENT_STOP
				case e.RuneAsKey() == cdk.KeyTAB || e.Key() == cdk.KeyTAB:
					w.LogDebug("tab key caught")
					if e.Modifiers().Has(cdk.ModShift) {
						w.FocusPrevious()
					} else {
						w.FocusNext()
					}
					return cenums.EVENT_STOP
				}

				// check focused
				if fi := w.GetFocus(); fi != nil {
					if sw, ok := fi.Self().(Sensitive); ok && sw.IsSensitive() && sw.IsVisible() {
						if f := sw.ProcessEvent(evt); f == cenums.EVENT_STOP {
							return cenums.EVENT_STOP
						}
					}
				}
			}

		case *cdk.EventMouse:
			// need to track enter/leave widget states
			if f := w.Emit(SignalEventMouse, w, e); f == cenums.EVENT_PASS {
				mousePosition := ptypes.NewPoint2I(e.Position())

				if fw := w.FindWidgetAt(mousePosition); fw != nil {
					var wantRefresh bool

					fwid := fw.ObjectID()

					for _, hfw := range *w.hoverFocused {
						if fwid != hfw.ObjectID() {
							if !hfw.HasPoint(mousePosition) {
								w.Lock()
								w.hoverFocused.Remove(hfw)
								w.Unlock()
								hfw.Emit(SignalLeave, e)
								wantRefresh = true
								hfw.Invalidate()
							}
						}
					}

					w.Lock()
					w.hoverFocus = fw
					w.Unlock()

					found := w.FindAllWidgetsAt(mousePosition)
					for i := len(found) - 1; i >= 0; i-- {
						aw := found[i]
						if idx := w.hoverFocused.IndexOf(aw); idx < 0 {
							w.Lock()
							w.hoverFocused.Append(aw)
							w.Unlock()
							aw.Emit(SignalEnter, e)
							wantRefresh = true
							aw.Invalidate()
						}
					}

					if wantRefresh {
						if d := w.GetDisplay(); d != nil {
							d.RequestDraw()
							d.RequestShow()
						}
					}
				}

				if mw := w.GetWidgetAt(mousePosition); mw != nil {
					if ms, ok := mw.Self().(Sensitive); ok {
						if ms.IsSensitive() && ms.IsVisible() {
							if f := ms.ProcessEvent(e); f == cenums.EVENT_STOP {
								if d := w.GetDisplay(); d != nil {
									d.RequestDraw()
									d.RequestSync()
								}
								return cenums.EVENT_STOP
							}
						}
					}
				}
				return cenums.EVENT_PASS
			}

		case *cdk.EventResize:
			alloc := ptypes.MakeRectangle(e.Size())
			origin := ptypes.MakePoint2I(0, 0)
			if rv := w.Emit(SignalResize, w, origin, alloc); rv == cenums.EVENT_STOP {
				w.LogDebug("ProcessEvent(EventResize): (%v) %v", alloc, e)
				if d := w.GetDisplay(); d != nil {
					d.RequestDraw()
					d.RequestSync()
				}
			}
		}

		w.LogTrace("ProcessEvent(cdk.Event): %v", evt)
		return cenums.EVENT_STOP
	}
	return cenums.EVENT_PASS
}

func (w *CWindow) resize(data []interface{}, argv ...interface{}) cenums.EventFlag {

	argc := len(argv)

	origin := w.GetOrigin()
	if argc >= 2 {
		if eOrigin, ok := argv[1].(ptypes.Point2I); ok {
			origin = eOrigin
		} else {
			w.LogError("expected Point2I, resize event received: %T %v - %v", argv[1], argv[1], argv)
			return cenums.EVENT_PASS
		}
	}

	alloc := w.GetAllocation()
	if argc >= 3 {
		if eAlloc, ok := argv[2].(ptypes.Rectangle); ok {
			alloc = eAlloc
		} else {
			w.LogError("expected Rectangle, resize event received: %T %v - %v", argv[1], argv[1], argv)
			return cenums.EVENT_PASS
		}
	}

	if w.GetWindowType() != cenums.WINDOW_POPUP {
		w.SetOrigin(origin.X, origin.Y)
		w.SetAllocation(alloc)
	}

	title := w.GetTitle()
	decorated := w.GetDecorated()

	childAlloc := alloc.Clone()
	childOrigin := origin.Clone()

	if title != "" {
		if decorated {
			// standard desktop layout, full border, title row, hline row
			childAlloc.W -= 2
			childAlloc.H -= 4
			childOrigin.X += 1
			childOrigin.Y += 3
		} else {
			// no border with title inline and centred on top edge
			// childAlloc.W -= 2
			childAlloc.H -= 3
			// childOrigin.X += 1
			childOrigin.Y += 3
		}
	} else if decorated {
		// border
		childAlloc.W -= 2
		childAlloc.H -= 2
		childOrigin.X += 1
		childOrigin.Y += 1
	} else {
		// fullscreen child (vbox)
	}

	if child := w.GetChild(); child != nil {
		child.SetOrigin(childOrigin.X, childOrigin.Y)
		child.SetAllocation(childAlloc)
		child.Resize()
	}

	w.LogDebug("window resized: origin=%v, alloc=%v", w.GetOrigin(), w.GetAllocation())
	w.Invalidate()
	return cenums.EVENT_STOP
}

func (w *CWindow) draw(data []interface{}, argv ...interface{}) cenums.EventFlag {

	oid := w.ObjectID()
	origin := w.GetOrigin()
	alloc := w.GetAllocation()
	theme := w.GetThemeRequest()

	if err := memphis.MakeConfigureSurface(oid, origin, alloc, theme.Content.Normal); err != nil {
		w.LogErr(err)
	} else if surface, err := memphis.GetSurface(oid); err != nil {
		w.LogErr(err)
	} else {

		if !w.IsVisible() || alloc.W == 0 || alloc.H == 0 {
			w.LogDebug("not visible, zero width or zero height")
			return cenums.EVENT_PASS
		}

		title := w.GetTitle()
		decorated := w.GetDecorated()

		_, _, attr := theme.Content.Normal.Decompose()
		dim := attr.IsDim()

		surface.FillBorder(dim, decorated, theme) // outer border

		if title != "" {
			if decorated {
				// standard desktop layout
				titleOrigin := origin.Clone()
				titleOrigin.X += 2
				titleOrigin.Y += 1
				surface.DrawSingleLineText(titleOrigin, alloc.W-4, false, cenums.JUSTIFY_LEFT, theme.Content.Normal, false, false, title)
				lineOrigin := origin.Clone()
				lineOrigin.X += 1
				lineOrigin.Y += 2
				for idx := 1; idx < alloc.W-1; idx += 1 {
					_ = surface.SetRune(idx, 2, paint.RuneHLine, theme.Content.Normal)
				}
				_ = surface.SetRune(0, 2, paint.RuneLTee, theme.Content.Normal)
				_ = surface.SetRune(alloc.W-1, 2, paint.RuneRTee, theme.Content.Normal)
				_ = surface.SetRune(alloc.W-1, alloc.H-1, theme.Border.BorderRunes.BottomRight, theme.Content.Normal)
			} else {
				// no border with title inline and centred on top edge
				surface.DrawHorizontalLine(origin, alloc.W, theme.Content.Normal, paint.RuneHLine)
				titleOrigin := origin.Clone()
				titleOrigin.X += 1
				titleOrigin.Y += 1
				surface.DrawSingleLineText(titleOrigin, alloc.W-2, false, cenums.JUSTIFY_LEFT, theme.Content.Normal, false, false, title)
				lineOrigin := origin.Clone()
				lineOrigin.X += 0
				lineOrigin.Y += 2
				surface.DrawHorizontalLine(lineOrigin, alloc.W, theme.Content.Normal, paint.RuneHLine)
			}
		}

		child := w.GetChild()
		if child != nil && child.IsVisible() {
			child.Draw()
			child.LockDraw()
			if err := surface.Composite(child.ObjectID()); err != nil {
				w.LogError("composite error: %v", err)
			}
			child.UnlockDraw()
		}

		if debug, _ := w.GetBoolProperty(cdk.PropertyDebug); debug {
			surface.DebugBox(paint.ColorNavy, w.ObjectInfo())
		}

		return cenums.EVENT_STOP
	}

	return cenums.EVENT_PASS
}

func (w *CWindow) displayFocusedWindow(data []interface{}, argv ...interface{}) cenums.EventFlag {
	if len(argv) >= 1 {
		if focused, ok := argv[0].(cdk.Window); ok {
			if focused.ObjectID() == w.ObjectID() {
				if internalFocused := w.GetFocus(); internalFocused != nil {
					w.SetFocus(nil)
					internalFocused.GrabFocus()
				}
			} else {
				w.display.Screen().HideCursor()
			}
		}
	}
	WidgetRecurseInvalidate(w)
	w.Invalidate()
	return cenums.EVENT_PASS
}

// Whether the window should receive the input focus.
// Flags: Read / Write
// Default value: TRUE
const PropertyAcceptFocus cdk.Property = "accept-focus"

// Whether the window should be decorated by the window manager.
// Flags: Read / Write
// Default value: TRUE
const PropertyDecorated cdk.Property = "decorated"

// The default height of the window, used when initially showing the window.
// Flags: Read / Write
// Allowed values: >= -1
// Default value: -1
const PropertyDefaultHeight cdk.Property = "default-height"

// The default width of the window, used when initially showing the window.
// Flags: Read / Write
// Allowed values: >= -1
// Default value: -1
const PropertyDefaultWidth cdk.Property = "default-width"

// Whether the window frame should have a close button.
// Flags: Read / Write
// Default value: TRUE
const PropertyDeletable cdk.Property = "deletable"

// If this window should be destroyed when the parent is destroyed.
// Flags: Read / Write
// Default value: FALSE
const PropertyDestroyWithParent cdk.Property = "destroy-with-parent"

// Whether the window should receive the input focus when mapped.
// Flags: Read / Write
// Default value: TRUE
const PropertyFocusOnMap cdk.Property = "focus-on-map"

// The window gravity of the window. See Move and Gravity for
// more details about window gravity.
// Flags: Read / Write
// Default value: GDK_GRAVITY_NORTH_WEST
const PropertyGravity cdk.Property = "gravity"

// Whether the input focus is within this Window.
// Flags: Read
// Default value: FALSE
const PropertyHasToplevelFocus cdk.Property = "has-toplevel-focus"

// Icon for this window.
// Flags: Read / Write
const PropertyIcon cdk.Property = "icon"

// The :icon-name property specifies the name of the themed icon to use as
// the window icon. See IconTheme for more details.
// Flags: Read / Write
// Default value: NULL
const PropertyIconName cdk.Property = "icon-name"

// Whether the toplevel is the current active window.
// Flags: Read
// Default value: FALSE
const PropertyIsActive cdk.Property = "is-active"

// Whether mnemonics are currently visible in this window.
// Flags: Read / Write
// Default value: TRUE
const PropertyMnemonicsVisible cdk.Property = "mnemonics-visible"

// If TRUE, the window is modal (other windows are not usable while this one
// is up).
// Flags: Read / Write
// Default value: FALSE
const PropertyModal cdk.Property = "modal"

// The requested opacity of the window. See SetOpacity for more
// details about window opacity.
// Flags: Read / Write
// Allowed values: [0,1]
// Default value: 1
const PropertyOpacity cdk.Property = "opacity"

// If TRUE, users can resize the window.
// Flags: Read / Write
// Default value: TRUE
const PropertyResizable cdk.Property = "resizable"

// Unique identifier for the window to be used when restoring a session.
// Flags: Read / Write
// Default value: NULL
const PropertyRole cdk.Property = "role"

// The screen where this window will be displayed.
// Flags: Read / Write
const PropertyScreen cdk.Property = "screen"

// TRUE if the window should not be in the pager.
// Flags: Read / Write
// Default value: FALSE
const PropertySkipPagerHint cdk.Property = "skip-pager-hint"

// TRUE if the window should not be in the task bar.
// Flags: Read / Write
// Default value: FALSE
const PropertySkipTaskbarHint cdk.Property = "skip-taskbar-hint"

// The :startup-id is a write-only property for setting window's startup
// notification identifier. See SetStartupId for more details.
// Flags: Write
// Default value: NULL
const PropertyStartupId cdk.Property = "startup-id"

// The title of the window.
// Flags: Read / Write
// Default value: NULL
const PropertyTitle cdk.Property = "title"

// The transient parent of the window. See SetTransientFor for
// more details about transient windows.
// Flags: Read / Write / Construct
const PropertyTransientFor cdk.Property = "transient-for"

// The type of the window.
// Flags: Read / Write / Construct Only
// Default value: GTK_WINDOW_TOPLEVEL
const PropertyWindowType = cdk.PropertyWindowType

// Hint to help the desktop environment understand what kind of window this
// is and how to treat it.
// Flags: Read / Write
// Default value: GDK_WINDOW_TYPE_HINT_NORMAL
const PropertyTypeHint cdk.Property = "type-hint"

// TRUE if the window should be brought to the user's attention.
// Flags: Read / Write
// Default value: FALSE
const PropertyUrgencyHint cdk.Property = "urgency-hint"

// The initial position of the window.
// Flags: Read / Write
// Default value: GTK_WIN_POS_NONE
const PropertyWindowPosition cdk.Property = "window-position"

const PropertyFocusedWidget = "focused-widget"

// The ::activate-default signal is a which gets emitted when the user
// activates the default widget of window .
const SignalActivateDefault cdk.Signal = "activate-default"

// The ::activate-focus signal is a which gets emitted when the user
// activates the currently focused widget of window .
const SignalActivateFocus cdk.Signal = "activate-focus"

const SignalFrameEvent cdk.Signal = "frame-event"

// The ::keys-changed signal gets emitted when the set of accelerators or
// mnemonics that are associated with window changes.
const SignalKeysChanged cdk.Signal = "keys-changed"

// Listener function arguments:
//
//	widget Widget
const SignalFocusChanged cdk.Signal = "focus-changed"

var ErrFallthrough = fmt.Errorf("fallthrough")

const WindowEventHandle = "window-event-handler"

const WindowDisplayFocusHandle = "window-display-focus-handler"

const WindowInvalidateHandle = "window-invalidate-handler"

const WindowResizeHandle = "window-resize-handler"

const WindowDrawHandle = "window-draw-handler"