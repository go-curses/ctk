package ctk

import (
	"sync"

	"github.com/go-curses/cdk"
	"github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/cdk/lib/paint"
	"github.com/go-curses/cdk/lib/ptypes"
	"github.com/go-curses/cdk/memphis"
)

// CDK type-tag for Widget objects
const TypeWidget cdk.CTypeTag = "ctk-widget"

func init() {
	_ = cdk.TypesManager.AddType(TypeWidget, nil)
}

// Widget Hierarchy:
//	Object
//	  +- Widget
//	    +- Container
//	    +- Misc
//	    +- Calendar
//	    +- CellView
//	    +- DrawingArea
//	    +- Entry
//	    +- Ruler
//	    +- Range
//	    +- Separator
//	    +- HSV
//	    +- Invisible
//	    +- OldEditable
//	    +- Preview
//	    +- Progress
//
// Widget is the base class all widgets in CTK derive from. It manages the
// widget lifecycle, states and style. Widget introduces style properties -
// these are basically object properties that are stored not on the object, but
// in the style object associated to the widget. Style properties are set in
// resource files. This mechanism is used for configuring such things as the
// location of the scrollbar arrows through the theme, giving theme authors more
// control over the look of applications without the need to write a theme
// engine. Use InstallStyleProperty to install style properties for a Widget,
// FindStyleProperty or ListStyleProperties to get information about existing
// style properties and Style.GetProperty() or Style.Get() to obtain the value
// of a style property. The Widget implementation of the Buildable interface
// supports a custom <accelerator> element, which has attributes named key,
// modifiers and signal and allows developers to specify accelerators.
//
// Note that usage of Style mentioned above is unimplemented at this time, the
// comment is left here as a reminder of what it should be doing. Currently,
// the ctk.Style is not used at all and instead, the paint.Style and paint.Theme
// from CDK are used throughout. To change the style of a Widget, construct a
// new paint.Theme structure and apply it to the Widget with SetTheme().
type Widget interface {
	Object

	Init() (already bool)
	Destroy()
	Unparent()
	Show()
	Hide()
	ShowAll()
	Map()
	Unmap()
	Realize()
	Unrealize()
	AddAccelerator(accelSignal string, accelGroup AccelGroup, accelKey int, accelMods ModifierType, accelFlags AccelFlags)
	RemoveAccelerator(accelGroup AccelGroup, accelKey int, accelMods ModifierType) (value bool)
	SetAccelPath(accelPath string, accelGroup AccelGroup)
	CanActivateAccel(signalId int) (value bool)
	Activate() (value bool)
	Reparent(parent Container)
	Intersect(area ptypes.Rectangle, intersection ptypes.Rectangle) (value bool)
	IsFocus() (value bool)
	GrabFocus()
	GrabDefault()
	SetName(name string)
	GetName() (value string)
	SetState(state StateType)
	SetSensitive(sensitive bool)
	SetParent(parent Container)
	SetParentWindow(parentWindow Window)
	GetParentWindow() (value Window)
	SetEvents(events cdk.EventMask)
	AddEvents(events cdk.EventMask)
	GetToplevel() (value Widget)
	GetAncestor(widgetType cdk.CTypeTag) (value Widget)
	GetEvents() (value cdk.EventMask)
	GetPointer(x int, y int)
	IsAncestor(ancestor Widget) (value bool)
	TranslateCoordinates(destWidget Widget, srcX int, srcY int, destX int, destY int) (value bool)
	HideOnDelete() (value bool)
	SetStyle(style Style)
	EnsureStyle()
	GetStyle() (value Style)
	ResetRcStyles()
	GetDefaultStyle() (value Style)
	SetDirection(dir TextDirection)
	GetDirection() (value TextDirection)
	SetDefaultDirection(dir TextDirection)
	GetDefaultDirection() (value TextDirection)
	Path() (path string)
	ClassPath(pathLength int, path string, pathReversed string)
	GetCompositeName() (value string)
	ModifyStyle(style RcStyle)
	GetModifierStyle() (value RcStyle)
	ModifyFg(state StateType, color paint.Color)
	ModifyBg(state StateType, color paint.Color)
	ModifyText(state StateType, color paint.Color)
	ModifyBase(state StateType, color paint.Color)
	SetAppPaintable(appPaintable bool)
	SetDoubleBuffered(doubleBuffered bool)
	SetRedrawOnAllocate(redrawOnAllocate bool)
	SetCompositeName(name string)
	SetScrollAdjustments(hAdjustment Adjustment, vAdjustment Adjustment) (value bool)
	RegionIntersect(region ptypes.Region) (value ptypes.Region)
	SendExpose(event cdk.Event) (value int)
	SendFocusChange(event cdk.Event) (value bool)
	ChildFocus(direction DirectionType) (value bool)
	ChildNotify(childProperty string)
	FreezeChildNotify()
	GetChildVisible() (value bool)
	GetParent() (value Container)
	GetDisplay() (value cdk.Display)
	GetRootWindow() (value Window)
	GetScreen() (value cdk.Display)
	HasScreen() (value bool)
	GetSizeRequest() (width, height int)
	SizeRequest() ptypes.Rectangle
	SetChildVisible(isVisible bool)
	SetSizeRequest(width, height int)
	ThawChildNotify()
	SetNoShowAll(noShowAll bool)
	GetNoShowAll() (value bool)
	AddMnemonicLabel(label Widget)
	RemoveMnemonicLabel(label Widget)
	IsComposited() (value bool)
	ErrorBell()
	KeynavFailed(direction DirectionType) (value bool)
	GetTooltipMarkup() (value string)
	SetTooltipMarkup(markup string)
	GetTooltipText() (value string)
	SetTooltipText(text string)
	GetTooltipWindow() (value Window)
	SetTooltipWindow(customWindow Window)
	GetHasTooltip() (value bool)
	SetHasTooltip(hasTooltip bool)
	TriggerTooltipQuery()
	GetWindow() (window Window)
	GetAppPaintable() (value bool)
	GetCanDefault() (value bool)
	SetCanDefault(canDefault bool)
	GetCanFocus() (value bool)
	SetCanFocus(canFocus bool)
	GetDoubleBuffered() (value bool)
	GetHasWindow() (value bool)
	SetHasWindow(hasWindow bool)
	GetSensitive() (value bool)
	IsSensitive() bool
	GetState() (value StateType)
	GetVisible() (value bool)
	SetVisible(visible bool)
	HasDefault() (value bool)
	HasFocus() (value bool)
	HasGrab() (value bool)
	HasRcStyle() (value bool)
	IsDrawable() (value bool)
	IsToplevel() (value bool)
	SetWindow(window Window)
	SetReceivesDefault(receivesDefault bool)
	GetReceivesDefault() (value bool)
	SetRealized(realized bool)
	GetRealized() (value bool)
	SetMapped(mapped bool)
	GetMapped() (value bool)
	GetThemeRequest() (theme paint.Theme)
	SetTheme(theme paint.Theme)
	HasState(s StateType) bool
	UnsetState(v StateType)
	GetFlags() WidgetFlags
	HasFlags(f WidgetFlags) bool
	UnsetFlags(v WidgetFlags)
	SetFlags(v WidgetFlags)
	IsParentFocused() bool
	IsFocused() bool
	CanFocus() bool
	IsDefault() bool
	CanDefault() bool
	IsVisible() bool
	HasEventFocus() bool
	GrabEventFocus()
	ReleaseEventFocus()
	GetTopParent() (parent Container)
	GetWidgetAt(p *ptypes.Point2I) Widget
	Draw() enums.EventFlag
}

// The CWidget structure implements the Widget interface and is exported
// to facilitate type embedding with custom implementations. No member variables
// are exported as the interface methods are the only intended means of
// interacting with Widget objects.
type CWidget struct {
	CObject

	display   cdk.Display
	parent    interface{}
	state     StateType
	flags     WidgetFlags
	flagsLock *sync.RWMutex
}

// Init initializes a Widget object. This must be called at least once to
// set up the necessary defaults and allocate any memory structures. Calling
// this more than once is safe though unnecessary. Only the first call will
// result in any effect upon the Widget instance. Init is used in the
// NewWidget constructor and only necessary when implementing a derivative
// Widget type.
func (w *CWidget) Init() (already bool) {
	if w.InitTypeItem(TypeWidget, w) {
		return true
	}
	w.CObject.Init()
	_ = w.InstallProperty(PropertyAppPaintable, cdk.BoolProperty, true, false)
	_ = w.InstallProperty(PropertyCanDefault, cdk.BoolProperty, true, false)
	_ = w.InstallProperty(PropertyCanFocus, cdk.BoolProperty, true, false)
	_ = w.InstallProperty(PropertyCompositeChild, cdk.BoolProperty, true, false)
	_ = w.InstallProperty(PropertyDoubleBuffered, cdk.BoolProperty, true, true)
	_ = w.InstallProperty(PropertyEvents, cdk.StructProperty, true, cdk.EVENT_MASK_NONE)
	// _ = w.InstallProperty(PropertyExtensionEvents, cdk.StructProperty, true, nil)
	_ = w.InstallProperty(PropertyHasDefault, cdk.BoolProperty, true, false)
	_ = w.InstallProperty(PropertyHasFocus, cdk.BoolProperty, true, false)
	_ = w.InstallProperty(PropertyHasTooltip, cdk.BoolProperty, true, false)
	_ = w.InstallProperty(PropertyHeightRequest, cdk.IntProperty, true, -1)
	_ = w.InstallProperty(PropertyIsFocus, cdk.BoolProperty, true, false)
	_ = w.InstallProperty(PropertyNoShowAll, cdk.BoolProperty, true, false)
	_ = w.InstallProperty(PropertyParent, cdk.StructProperty, true, nil)
	_ = w.InstallProperty(PropertyReceivesDefault, cdk.BoolProperty, true, false)
	_ = w.InstallProperty(PropertySensitive, cdk.BoolProperty, true, true)
	_ = w.InstallProperty(PropertyStyle, cdk.StructProperty, true, &CStyle{})
	_ = w.InstallProperty(PropertyTooltipMarkup, cdk.StringProperty, true, "")
	_ = w.InstallProperty(PropertyTooltipText, cdk.StringProperty, true, "")
	_ = w.InstallProperty(PropertyVisible, cdk.BoolProperty, true, false)
	_ = w.InstallProperty(PropertyWidthRequest, cdk.IntProperty, true, -1)
	_ = w.InstallProperty(PropertyWindow, cdk.StructProperty, true, nil)
	w.flagsLock = &sync.RWMutex{}
	w.state = StateNormal
	w.flags = NULL_WIDGET_FLAG
	w.Connect(SignalLostFocus, WidgetLostFocusHandle, w.lostFocus)
	w.Connect(SignalGainedFocus, WidgetGainedFocusHandle, w.gainedFocus)
	// TODO: widget register surface tied to APP_PAINTABLE somehow
	if err := memphis.RegisterSurface(w.ObjectID(), ptypes.Point2I{}, ptypes.Rectangle{}, w.GetTheme().Content.Normal); err != nil {
		w.LogErr(err)
	}
	return false
}

// Destroy a widget. Equivalent to DestroyObject. When a widget is destroyed, it
// will break any references it holds to other objects. If the widget is
// inside a container, the widget will be removed from the container. If the
// widget is a toplevel (derived from Window), it will be removed from the
// list of toplevels, and the reference CTK holds to it will be removed.
// Removing a widget from its container or the list of toplevels results in
// the widget being finalized, unless you've added additional references to
// the widget with g_object_ref. In most cases, only toplevel widgets
// (windows) require explicit destruction, because when you destroy a
// toplevel its children will be destroyed as well.
func (w *CWidget) Destroy() {
	if err := w.DestroyObject(); err != nil {
		w.LogErr(err)
	}
}

// Unparent is only for use in widget implementations. Should be called by
// implementations of the remove method on Container, to dissociate a child from
// the container.
func (w *CWidget) Unparent() {
	if w.parent != nil {
		if parent, ok := w.parent.(Container); ok {
			if f := w.Emit(SignalUnparent, w, w.parent); f == enums.EVENT_PASS {
				parent.Remove(w)
			}
		}
	}
}

// Show flags a widget to be displayed. Any widget that isn't shown will not
// appear on the screen. If you want to show all the widgets in a container,
// it's easier to call ShowAll on the container, instead of
// individually showing the widgets. Remember that you have to show the
// containers containing a widget, in addition to the widget itself, before
// it will appear onscreen. When a toplevel container is shown, it is
// immediately realized and mapped; other shown widgets are realized and
// mapped when their toplevel container is realized and mapped.
func (w *CWidget) Show() {
	if !w.HasFlags(VISIBLE) {
		if r := w.Emit(SignalShow, w); r == enums.EVENT_PASS {
			w.SetFlags(VISIBLE)
			w.Invalidate()
		}
	}
}

// Hide reverses the effects of Show, causing the widget to be hidden (invisible
// to the user).
func (w *CWidget) Hide() {
	if w.HasFlags(VISIBLE) {
		if r := w.Emit(SignalHide, w); r == enums.EVENT_PASS {
			w.UnsetFlags(VISIBLE)
			w.Invalidate()
		}
	}
}

// ShowAll recursively shows a widget, and any child widgets (if the widget is a
// container).
func (w *CWidget) ShowAll() {
	w.Show()
}

// This function is only for use in widget implementations. Causes a widget
// to be mapped if it isn't already.
//
func (w *CWidget) Map() {}

// This function is only for use in widget implementations. Causes a widget
// to be unmapped if it's currently mapped.
func (w *CWidget) Unmap() {}

// Creates the GDK (windowing system) resources associated with a widget. For
// example, widget->window will be created when a widget is realized.
// Normally realization happens implicitly; if you show a widget and all its
// parent containers, then the widget will be realized and mapped
// automatically. Realizing a widget requires all the widget's parent widgets
// to be realized; calling Realize realizes the widget's parents
// in addition to widget itself. If a widget is not yet inside a toplevel
// window when you realize it, bad things will happen. This function is
// primarily used in widget implementations, and isn't very useful otherwise.
// Many times when you think you might need it, a better approach is to
// connect to a signal that will be called after the widget is realized
// automatically, such as Widget::expose-event. Or simply
// g_signal_connect to the Widget::realize signal.
func (w *CWidget) Realize() {}

// This function is only useful in widget implementations. Causes a widget to
// be unrealized (frees all GDK resources associated with the widget, such as
// widget->window ).
func (w *CWidget) Unrealize() {}

// Installs an accelerator for this widget in accel_group that causes
// accel_signal to be emitted if the accelerator is activated. The
// accel_group needs to be added to the widget's toplevel via
// WindowAddAccelGroup, and the signal must be of type G_RUN_ACTION.
// Accelerators added through this function are not user changeable during
// runtime. If you want to support accelerators that can be changed by the
// user, use AccelMapAddEntry and SetAccelPath or
// MenuItemSetAccelPath instead.
// Parameters:
// 	widget	widget to install an accelerator on
// 	accelSignal	widget signal to emit on accelerator activation
// 	accelGroup	accel group for this widget, added to its toplevel
// 	accelKey	GDK keyval of the accelerator
// 	accelMods	modifier key combination of the accelerator
// 	accelFlags	flag accelerators, e.g. GTK_ACCEL_VISIBLE
func (w *CWidget) AddAccelerator(accelSignal string, accelGroup AccelGroup, accelKey int, accelMods ModifierType, accelFlags AccelFlags) {
}

// Removes an accelerator from widget , previously installed with
// AddAccelerator.
// Parameters:
// 	widget	widget to install an accelerator on
// 	accelGroup	accel group for this widget
// 	accelKey	GDK keyval of the accelerator
// 	accelMods	modifier key combination of the accelerator
// 	returns	whether an accelerator was installed and could be removed
func (w *CWidget) RemoveAccelerator(accelGroup AccelGroup, accelKey int, accelMods ModifierType) (value bool) {
	return false
}

// Given an accelerator group, accel_group , and an accelerator path,
// accel_path , sets up an accelerator in accel_group so whenever the key
// binding that is defined for accel_path is pressed, widget will be
// activated. This removes any accelerators (for any accelerator group)
// installed by previous calls to SetAccelPath. Associating
// accelerators with paths allows them to be modified by the user and the
// modifications to be saved for future use. (See AccelMapSave.) This
// function is a low level function that would most likely be used by a menu
// creation system like UIManager. If you use UIManager, setting up
// accelerator paths will be done automatically. Even when you you aren't
// using UIManager, if you only want to set up accelerators on menu items
// MenuItemSetAccelPath provides a somewhat more convenient
// interface. Note that accel_path string will be stored in a GQuark.
// Therefore, if you pass a static string, you can save some memory by
// interning it first with g_intern_static_string.
// Parameters:
// 	accelPath	path used to look up the accelerator.
// 	accelGroup	a AccelGroup.
func (w *CWidget) SetAccelPath(accelPath string, accelGroup AccelGroup) {}

// Determines whether an accelerator that activates the signal identified by
// signal_id can currently be activated. This is done by emitting the
// can-activate-accel signal on widget ; if the signal isn't overridden
// by a handler or in a derived widget, then the default check is that the
// widget must be sensitive, and the widget and all its ancestors mapped.
// Parameters:
// 	signalId	the ID of a signal installed on widget
//
// Returns:
// 	TRUE if the accelerator can be activated.
func (w *CWidget) CanActivateAccel(signalId int) (value bool) {
	return false
}

// For widgets that can be "activated" (buttons, menu items, etc.) this
// function activates them. Activation is what happens when you press Enter
// on a widget during key navigation. If widget isn't activatable, the
// function returns FALSE.
// Returns:
// 	TRUE if the widget was activatable
func (w *CWidget) Activate() (value bool) {
	if w.IsSensitive() {
		w.Emit(SignalActivate, w)
		value = true
	}
	return
}

// Move the Widget to the given container, removing itself first from any other
// container that was currently holding it. This method emits a reparent signal
// initially and if the listeners return EVENT_PAS, the change is applied
// Parameters:
// 	newParent	a Container to move the widget into
//
// Emits: SignalReparent, Argv=[Widget instance, new parent]
func (w *CWidget) Reparent(parent Container) {
	if r := w.Emit(SignalReparent, w, parent); r == enums.EVENT_PASS {
		if w.parent != nil {
			if pc, ok := w.parent.(Container); ok {
				pc.Remove(w)
			}
		}
		parent.Add(w)
	}
}

// Computes the intersection of a widget 's area and area , storing the
// intersection in intersection , and returns TRUE if there was an
// intersection. intersection may be NULL if you're only interested in
// whether there was an intersection.
// Parameters:
// 	area	a rectangle
// 	intersection	rectangle to store intersection of widget
// and area
//
// Returns:
// 	TRUE if there was an intersection
func (w *CWidget) Intersect(area ptypes.Rectangle, intersection ptypes.Rectangle) (value bool) {
	return false
}

// Determines if the widget is the focus widget within its toplevel. (This
// does not mean that the HAS_FOCUS flag is necessarily set; HAS_FOCUS will
// only be set if the toplevel widget additionally has the global input
// focus.)
// Returns:
// 	TRUE if the widget is the focus widget.
// Returns TRUE if the Widget instance is currently the focus of it's parent
// Window, FALSE otherwise
func (w *CWidget) IsFocus() (value bool) {
	if window := w.GetWindow(); window != nil {
		if w.CanFocus() {
			if focused := window.GetFocus(); focused != nil {
				if fw, ok := focused.(Widget); ok {
					if fw.ObjectID() == w.ObjectID() {
						return true
					}
				}
			}
		}
	}
	return false
}

// Causes widget to have the keyboard focus for the Window it's inside.
// widget must be a focusable widget, such as a Entry; something like
// Frame won't work. More precisely, it must have the GTK_CAN_FOCUS flag
// set. Use SetCanFocus to modify that flag. The widget also
// needs to be realized and mapped. This is indicated by the related signals.
// Grabbing the focus immediately after creating the widget will likely fail
// and cause critical warnings.
// If the Widget instance CanFocus() then take the focus of the associated
// Window. Any previously focused Widget will emit a lost-focus signal and the
// newly focused Widget will emit a gained-focus signal. This method emits a
// grab-focus signal initially and if the listeners return EVENT_PASS, the
// changes are applied
//
// Emits: SignalGrabFocus, Argv=[Widget instance]
// Emits: SignalLostFocus, Argv=[Previous focus Widget instance], From=Previous focus Widget instance
// Emits: SignalGainedFocus, Argv=[Widget instance, previous focus Widget instance]
func (w *CWidget) GrabFocus() {
	if w.CanFocus() {
		if r := w.Emit(SignalGrabFocus, w); r == enums.EVENT_PASS {
			tl := w.GetWindow()
			if tl != nil {
				var fw Widget
				focused := tl.GetFocus()
				tl.SetFocus(w)
				if focused != nil {
					var ok bool
					if fw, ok = focused.(Widget); ok && fw.ObjectID() != w.ObjectID() {
						if f := fw.Emit(SignalLostFocus, focused); f == enums.EVENT_STOP {
							fw = nil
						}
					}
				}
				if f := w.Emit(SignalGainedFocus, w, focused); f == enums.EVENT_STOP {
					if fw != nil {
						tl.SetFocus(focused)
					}
				}
				w.LogDebug("has taken focus")
			}
		}
	}
}

// Causes widget to become the default widget. widget must have the
// GTK_CAN_DEFAULT flag set; typically you have to set this flag yourself by
// calling SetCanDefault (widget , TRUE). The default widget is
// activated when the user presses Enter in a window. Default widgets must be
// activatable, that is, Activate should affect them.
func (w *CWidget) GrabDefault() {}

// Widgets can be named, which allows you to refer to them from a gtkrc file.
// You can apply a style to widgets with a particular name in the gtkrc file.
// See the documentation for gtkrc files (on the same page as the docs for
// RcStyle). Note that widget names are separated by periods in paths (see
// Path), so names with embedded periods may cause confusion.
// Parameters:
// 	name	name for the widget
func (w *CWidget) SetName(name string) {
	if err := w.SetStringProperty(PropertyName, name); err != nil {
		w.LogErr(err)
	}
}

// Retrieves the name of a widget. See SetName for the
// significance of widget names.
// Returns:
// 	name of the widget. This string is owned by CTK and should not
// 	be modified or freed
func (w *CWidget) GetName() (value string) {
	var err error
	if value, err = w.GetStringProperty(PropertyName); err != nil {
		w.LogErr(err)
	}
	return
}

// This function is for use in widget implementations. Sets the state of a
// widget (insensitive, prelighted, etc.) Usually you should set the state
// using wrapper functions such as SetSensitive.
// Parameters:
// 	state	new state for widget
//
// Adds the given state bitmask to the Widget instance. This method emits a
// set-state signal initially and if the listeners return EVENT_PASS, the change
// is applied
//
// Emit: SignalSetState, Argv=[Widget instance, given state to set]
func (w *CWidget) SetState(state StateType) {
	if f := w.Emit(SignalSetState, w, state); f == enums.EVENT_PASS {
		w.state = w.state | state
	}
}

// Sets the sensitivity of a widget. A widget is sensitive if the user can
// interact with it. Insensitive widgets are "grayed out" and the user can't
// interact with them. Insensitive widgets are known as "inactive",
// "disabled", or "ghosted" in some other toolkits.
// Parameters:
// 	sensitive	TRUE to make the widget sensitive
// Emits: SignalSetSensitive, Argv=[Widget instance, given sensitive bool]
func (w *CWidget) SetSensitive(sensitive bool) {
	if f := w.Emit(SignalSetSensitive, w, sensitive); f == enums.EVENT_PASS {
		if !sensitive {
			w.SetState(StateInsensitive)
		} else {
			w.UnsetState(StateInsensitive)
		}
		if err := w.SetBoolProperty(PropertySensitive, sensitive); err != nil {
			w.LogErr(err)
		}
	}
}

// This function is useful only when implementing subclasses of Container.
// Sets the container as the parent of widget , and takes care of some
// details such as updating the state and style of the child to reflect its
// new location. The opposite function is Unparent.
// Parameters:
// 	parent	parent container
func (w *CWidget) SetParent(parent Container) {
	if f := w.Emit(SignalSetParent, w, w.parent, parent); f == enums.EVENT_PASS {
		if w.HasFlags(PARENT_SENSITIVE) && w.parent != nil {
			if parent, ok := w.parent.(Widget); ok {
				_ = parent.Disconnect(SignalLostFocus, WidgetLostFocusHandle)
				_ = parent.Disconnect(SignalGainedFocus, WidgetGainedFocusHandle)
			}
		}
		if err := w.SetStructProperty(PropertyParent, parent); err != nil {
			w.LogErr(err)
		} else {
			if cw, ok := parent.(Widget); parent != nil && ok && w.HasFlags(PARENT_SENSITIVE) {
				cw.Connect(SignalLostFocus, WidgetLostFocusHandle, w.lostFocus)
				cw.Connect(SignalGainedFocus, WidgetGainedFocusHandle, w.gainedFocus)
			}
		}
	}
}

// Sets a non default parent window for widget .
// Parameters:
// 	parentWindow	the new parent window.
func (w *CWidget) SetParentWindow(parentWindow Window) {}

// Gets widget 's parent window.
// Returns:
// 	the parent window of widget .
// 	[transfer none]
func (w *CWidget) GetParentWindow() (value Window) {
	return nil
}

// Sets the event mask (see EventMask) for a widget. The event mask
// determines which events a widget will receive. Keep in mind that different
// widgets have different default event masks, and by changing the event mask
// you may disrupt a widget's functionality, so be careful. This function
// must be called while a widget is unrealized. Consider
// AddEvents for widgets that are already realized, or if you
// want to preserve the existing event mask. This function can't be used with
// GTK_NO_WINDOW widgets; to get events on those widgets, place them inside a
// EventBox and receive events on the event box.
// Parameters:
// 	events	event mask
func (w *CWidget) SetEvents(events cdk.EventMask) {
	if err := w.SetStructProperty(PropertyEvents, events); err != nil {
		w.LogErr(err)
	}
}

// Adds the events in the bitfield events to the event mask for widget . See
// SetEvents for details.
// Parameters:
// 	events	an event mask, see EventMask
func (w *CWidget) AddEvents(events cdk.EventMask) {}

// This function returns the topmost widget in the container hierarchy widget
// is a part of. If widget has no parent widgets, it will be returned as the
// topmost widget. No reference will be added to the returned widget; it
// should not be unreferenced. Note the difference in behavior vs.
// GetAncestor; GetAncestor (widget,
// GTK_TYPE_WINDOW) would return NULL if widget wasn't inside a toplevel
// window, and if the window was inside a Window-derived widget which was
// in turn inside the toplevel Window. While the second case may seem
// unlikely, it actually happens when a Plug is embedded inside a
// Socket within the same application. To reliably find the toplevel
// Window, use GetToplevel and check if the TOPLEVEL flags
// is set on the result.
// Returns:
// 	the topmost ancestor of widget , or widget itself if there's no
// 	ancestor.
// 	[transfer none]
func (w *CWidget) GetToplevel() (value Widget) {
	return nil
}

// Gets the first ancestor of widget with type widget_type . For example,
// GetAncestor (widget, GTK_TYPE_BOX) gets the first Box
// that's an ancestor of widget . No reference will be added to the returned
// widget; it should not be unreferenced. See note about checking for a
// toplevel Window in the docs for GetToplevel. Note that
// unlike IsAncestor, GetAncestor considers
// widget to be an ancestor of itself.
// Parameters:
// 	widgetType	ancestor type
// Returns:
// 	the ancestor widget, or NULL if not found.
// 	[transfer none]
func (w *CWidget) GetAncestor(widgetType cdk.CTypeTag) (value Widget) {
	return nil
}

// Returns the event mask for the widget (a bitfield containing flags from
// the EventMask enumeration). These are the events that the widget will
// receive.
// Returns:
// 	event mask for widget
func (w *CWidget) GetEvents() (value cdk.EventMask) {
	if v, err := w.GetStructProperty(PropertyEvents); err != nil {
		w.LogErr(err)
	} else {
		var ok bool
		if value, ok = v.(cdk.EventMask); !ok {
			w.LogError("value stored in %v property is not of cdk.EventMask type: %v (%T)", PropertyEvents, v, v)
		}
	}
	return
}

// Obtains the location of the mouse pointer in widget coordinates. Widget
// coordinates are a bit odd; for historical reasons, they are defined as
// widget->window coordinates for widgets that are not GTK_NO_WINDOW widgets,
// and are relative to widget->allocation.x , widget->allocation.y for
// widgets that are GTK_NO_WINDOW widgets.
// Parameters:
// 	x	return location for the X coordinate, or NULL.
// 	y	return location for the Y coordinate, or NULL.
func (w *CWidget) GetPointer(x int, y int) {}

// Determines whether widget is somewhere inside ancestor , possibly with
// intermediate containers.
// Returns:
// 	TRUE if ancestor contains widget as a child, grandchild, great
// 	grandchild, etc.
func (w *CWidget) IsAncestor(ancestor Widget) (value bool) {
	return false
}

// Translate coordinates relative to src_widget 's allocation to coordinates
// relative to dest_widget 's allocations. In order to perform this
// operation, both widgets must be realized, and must share a common
// toplevel.
// Parameters:
// 	srcX	X position relative to src_widget
//
// 	srcY	Y position relative to src_widget
//
// 	destX	location to store X position relative to dest_widget
// .
// 	destY	location to store Y position relative to dest_widget
// .
// Returns:
// 	FALSE if either widget was not realized, or there was no common
// 	ancestor. In this case, nothing is stored in *dest_x and
// 	*dest_y . Otherwise TRUE.
func (w *CWidget) TranslateCoordinates(destWidget Widget, srcX int, srcY int, destX int, destY int) (value bool) {
	return false
}

// Utility function; intended to be connected to the delete-event
// signal on a Window. The function calls Hide on its
// argument, then returns TRUE. If connected to ::delete-event, the result is
// that clicking the close button for a window (on the window frame, top
// right corner usually) will hide but not destroy the window. By default,
// CTK destroys windows when ::delete-event is received.
// Returns:
// 	TRUE
func (w *CWidget) HideOnDelete() (value bool) {
	return false
}

// Sets the Style for a widget (widget->style ). You probably don't want
// to use this function; it interacts badly with themes, because themes work
// by replacing the Style. Instead, use ModifyStyle.
// Parameters:
// 	style	a Style, or NULL to remove the effect of a previous
// SetStyle and go back to the default style.
func (w *CWidget) SetStyle(style Style) {
	if err := w.SetStructProperty(PropertyStyle, style); err != nil {
		w.LogErr(err)
	}
}

// Ensures that widget has a style (widget->style ). Not a very useful
// function; most of the time, if you want the style, the widget is realized,
// and realized widgets are guaranteed to have a style already.
func (w *CWidget) EnsureStyle() {}

// Simply an accessor function that returns widget->style .
// Returns:
// 	the widget's Style.
// 	[transfer none]
func (w *CWidget) GetStyle() (value Style) {
	var ok bool
	if v, err := w.GetStructProperty(PropertyStyle); err != nil {
		w.LogErr(err)
	} else if value, ok = v.(Style); !ok {
		w.LogError("value stored in %v property is not of Style type: %v (%T)", PropertyStyle, v, v)
	}
	return
}

// Reset the styles of widget and all descendents, so when they are looked up
// again, they get the correct values for the currently loaded RC file
// settings. This function is not useful for applications.
func (w *CWidget) ResetRcStyles() {}

// Returns the default style used by all widgets initially.
// Returns:
// 	the default style. This Style object is owned by CTK and
// 	should not be modified or freed.
// 	[transfer none]
func (w *CWidget) GetDefaultStyle() (value Style) {
	return nil
}

// Sets the reading direction on a particular widget. This direction controls
// the primary direction for widgets containing text, and also the direction
// in which the children of a container are packed. The ability to set the
// direction is present in order so that correct localization into languages
// with right-to-left reading directions can be done. Generally, applications
// will let the default reading direction present, except for containers
// where the containers are arranged in an order that is explicitely visual
// rather than logical (such as buttons for text justification). If the
// direction is set to GTK_TEXT_DIR_NONE, then the value set by
// SetDefaultDirection will be used.
// Parameters:
// 	dir	the new direction
func (w *CWidget) SetDirection(dir TextDirection) {}

// Gets the reading direction for a particular widget. See
// SetDirection.
// Returns:
// 	the reading direction for the widget.
func (w *CWidget) GetDirection() (value TextDirection) {
	return TextDirLtr
}

// Sets the default reading direction for widgets where the direction has not
// been explicitly set by SetDirection.
// Parameters:
// 	dir	the new default direction. This cannot be
// GTK_TEXT_DIR_NONE.
func (w *CWidget) SetDefaultDirection(dir TextDirection) {}

// Obtains the current default reading direction. See
// SetDefaultDirection.
// Returns:
// 	the current default direction.
func (w *CWidget) GetDefaultDirection() (value TextDirection) {
	return TextDirLtr
}

// Sets a shape for this widget's GDK window. This allows for transparent
// windows etc., see WindowShapeCombineMask for more information.
// Parameters:
// 	shapeMask	shape to be added, or NULL to remove an existing shape.
// 	offsetX	X position of shape mask with respect to window
//
// 	offsetY	Y position of shape mask with respect to window
//
// func (w *CWidget) ShapeCombineMask(shapeMask Bitmap, offsetX int, offsetY int) {}

// Sets an input shape for this widget's GDK window. This allows for windows
// which react to mouse click in a nonrectangular region, see
// WindowInputShapeCombineMask for more information.
// Parameters:
// 	shapeMask	shape to be added, or NULL to remove an existing shape.
// 	offsetX	X position of shape mask with respect to window
//
// 	offsetY	Y position of shape mask with respect to window
//
// func (w *CWidget) InputShapeCombineMask(shapeMask Bitmap, offsetX int, offsetY int) {}

// Obtains the full path to widget . The path is simply the name of a widget
// and all its parents in the container hierarchy, separated by periods. The
// name of a widget comes from GetName. Paths are used to apply
// styles to a widget in gtkrc configuration files. Widget names are the type
// of the widget by default (e.g. "Button") or can be set to an
// application-specific value with SetName. By setting the name
// of a widget, you allow users or theme authors to apply styles to that
// specific widget in their gtkrc file. path_reversed_p fills in the path in
// reverse order, i.e. starting with widget 's name instead of starting with
// the name of widget 's outermost ancestor.
// Parameters:
// 	pathLength	location to store length of the path, or NULL.
// 	path	location to store allocated path string, or NULL.
// 	pathReversed	location to store allocated reverse path string, or NULL.
func (w *CWidget) Path() (path string) {
	var parents []Widget
	parent := w.GetParent()
	for {
		if parent == nil {
			break
		}
		parents = append(parents, parent)
		grandparent := parent.GetParent()
		if grandparent != nil && grandparent.ObjectID() == parent.ObjectID() {
			break
		}
	}
	for i := 0; i < len(parents); i++ {
		parent := parents[i]
		parentName := parent.GetTypeTag().ClassName()
		if name := parent.GetName(); name != "" {
			parentName += "#" + name
		}
		if len(path) > 0 {
			path += ">"
		}
		path += parentName
	}
	return
}

// Same as Path, but always uses the name of a widget's type,
// never uses a custom name set with SetName.
// Parameters:
// 	pathLength	location to store the length of the class path, or NULL.
// 	path	location to store the class path as an allocated string, or NULL.
// 	pathReversed	location to store the reverse class path as an allocated
// string, or NULL.
func (w *CWidget) ClassPath(pathLength int, path string, pathReversed string) {}

// Obtains the composite name of a widget.
// Returns:
// 	the composite name of widget , or NULL if widget is not a
// 	composite child. The string should be freed when it is no
// 	longer needed.
func (w *CWidget) GetCompositeName() (value string) {
	return ""
}

// Modifies style values on the widget. Modifications made using this
// technique take precedence over style values set via an RC file, however,
// they will be overriden if a style is explicitely set on the widget using
// SetStyle. The RcStyle structure is designed so each field
// can either be set or unset, so it is possible, using this function, to
// modify some style values and leave the others unchanged. Note that
// modifications made with this function are not cumulative with previous
// calls to ModifyStyle or with such functions as
// ModifyFg. If you wish to retain previous values, you must
// first call GetModifierStyle, make your modifications to the
// returned style, then call ModifyStyle with that style. On
// the other hand, if you first call ModifyStyle, subsequent
// calls to such functions ModifyFg will have a cumulative
// effect with the initial modifications.
// Parameters:
// 	style	the RcStyle holding the style modifications
func (w *CWidget) ModifyStyle(style RcStyle) {}

// Returns the current modifier style for the widget. (As set by
// ModifyStyle.) If no style has previously set, a new
// RcStyle will be created with all values unset, and set as the modifier
// style for the widget. If you make changes to this rc style, you must call
// ModifyStyle, passing in the returned rc style, to make sure
// that your changes take effect. Caution: passing the style back to
// ModifyStyle will normally end up destroying it, because
// ModifyStyle copies the passed-in style and sets the copy as
// the new modifier style, thus dropping any reference to the old modifier
// style. Add a reference to the modifier style if you want to keep it alive.
// Returns:
// 	the modifier style for the widget. This rc style is owned by
// 	the widget. If you want to keep a pointer to value this around,
// 	you must add a refcount using g_object_ref.
// 	[transfer none]
func (w *CWidget) GetModifierStyle() (value RcStyle) {
	return nil
}

// Sets the foreground color for a widget in a particular state. All other
// style values are left untouched. See also ModifyStyle.
// Parameters:
// 	state	the state for which to set the foreground color
// 	color	the color to assign (does not need to be allocated),
// or NULL to undo the effect of previous calls to
// of ModifyFg.
func (w *CWidget) ModifyFg(state StateType, color paint.Color) {}

// Sets the background color for a widget in a particular state. All other
// style values are left untouched. See also ModifyStyle. Note
// that "no window" widgets (which have the GTK_NO_WINDOW flag set) draw on
// their parent container's window and thus may not draw any background
// themselves. This is the case for e.g. Label. To modify the background
// of such widgets, you have to set the background color on their parent; if
// you want to set the background of a rectangular area around a label, try
// placing the label in a EventBox widget and setting the background color
// on that.
// Parameters:
// 	state	the state for which to set the background color
// 	color	the color to assign (does not need to be allocated),
// or NULL to undo the effect of previous calls to
// of ModifyBg.
func (w *CWidget) ModifyBg(state StateType, color paint.Color) {}

// Sets the text color for a widget in a particular state. All other style
// values are left untouched. The text color is the foreground color used
// along with the base color (see ModifyBase) for widgets such
// as Entry and TextView. See also ModifyStyle.
// Parameters:
// 	state	the state for which to set the text color
// 	color	the color to assign (does not need to be allocated),
// or NULL to undo the effect of previous calls to
// of ModifyText.
func (w *CWidget) ModifyText(state StateType, color paint.Color) {}

// Sets the base color for a widget in a particular state. All other style
// values are left untouched. The base color is the background color used
// along with the text color (see ModifyText) for widgets such
// as Entry and TextView. See also ModifyStyle. Note that
// "no window" widgets (which have the GTK_NO_WINDOW flag set) draw on their
// parent container's window and thus may not draw any background themselves.
// This is the case for e.g. Label. To modify the background of such
// widgets, you have to set the base color on their parent; if you want to
// set the background of a rectangular area around a label, try placing the
// label in a EventBox widget and setting the base color on that.
// Parameters:
// 	state	the state for which to set the base color
// 	color	the color to assign (does not need to be allocated),
// or NULL to undo the effect of previous calls to
// of ModifyBase.
func (w *CWidget) ModifyBase(state StateType, color paint.Color) {}

// Sets whether the application intends to draw on the widget in an
// expose-event handler. This is a hint to the widget and does not
// affect the behavior of the CTK core; many widgets ignore this flag
// entirely. For widgets that do pay attention to the flag, such as
// EventBox and Window, the effect is to suppress default themed
// drawing of the widget's background. (Children of the widget will still be
// drawn.) The application is then entirely responsible for drawing the
// widget background. Note that the background is still drawn when the widget
// is mapped. If this is not suitable (e.g. because you want to make a
// transparent window using an RGBA visual), you can work around this by
// doing:
// Parameters:
// 	appPaintable	TRUE if the application will paint on the widget
func (w *CWidget) SetAppPaintable(appPaintable bool) {
	w.SetFlags(APP_PAINTABLE)
	if err := w.SetBoolProperty(PropertyAppPaintable, appPaintable); err != nil {
		w.LogErr(err)
	}
}

// Widgets are double buffered by default; you can use this function to turn
// off the buffering. "Double buffered" simply means that
// WindowBeginPaintRegion and WindowEndPaint are called
// automatically around expose events sent to the widget.
// WindowBeginPaint diverts all drawing to a widget's window to an
// offscreen buffer, and WindowEndPaint draws the buffer to the
// screen. The result is that users see the window update in one smooth step,
// and don't see individual graphics primitives being rendered. In very
// simple terms, double buffered widgets don't flicker, so you would only use
// this function to turn off double buffering if you had special needs and
// really knew what you were doing. Note: if you turn off double-buffering,
// you have to handle expose events, since even the clearing to the
// background color or pixmap will not happen automatically (as it is done in
// WindowBeginPaint).
// Parameters:
// 	doubleBuffered	TRUE to double-buffer a widget
func (w *CWidget) SetDoubleBuffered(doubleBuffered bool) {
	if err := w.SetBoolProperty(PropertyDoubleBuffered, doubleBuffered); err != nil {
		w.LogErr(err)
	}
}

// Sets whether the entire widget is queued for drawing when its size
// allocation changes. By default, this setting is TRUE and the entire widget
// is redrawn on every size change. If your widget leaves the upper left
// unchanged when made bigger, turning this setting off will improve
// performance. Note that for NO_WINDOW widgets setting this flag to FALSE
// turns off all allocation on resizing: the widget will not even redraw if
// its position changes; this is to allow containers that don't draw anything
// to avoid excess invalidations. If you set this flag on a NO_WINDOW widget
// that does draw on widget->window , you are responsible for invalidating
// both the old and new allocation of the widget when the widget is moved and
// responsible for invalidating regions newly when the widget increases size.
// Parameters:
// 	redrawOnAllocate	if TRUE, the entire widget will be redrawn
// when it is allocated to a new size. Otherwise, only the
// new portion of the widget will be redrawn.
func (w *CWidget) SetRedrawOnAllocate(redrawOnAllocate bool) {}

// Sets a widgets composite name. The widget must be a composite child of its
// parent; see PushCompositeChild.
// Parameters:
// 	name	the name to set
func (w *CWidget) SetCompositeName(name string) {}

// For widgets that support scrolling, sets the scroll adjustments and
// returns TRUE. For widgets that don't support scrolling, does nothing and
// returns FALSE. Widgets that don't support scrolling can be scrolled by
// placing them in a Viewport, which does support scrolling.
// Parameters:
// 	hadjustment	an adjustment for horizontal scrolling, or NULL.
// 	vadjustment	an adjustment for vertical scrolling, or NULL.
// Returns:
// 	TRUE if the widget supports scrolling
func (w *CWidget) SetScrollAdjustments(hadjustment Adjustment, vadjustment Adjustment) (value bool) {
	return false
}

// Emits a draw signal, primarily used to render canvases and cause end-user
// facing display updates. Signal listeners can draw to the Canvas and return
// EVENT_STOP to cause those changes to be composited upon the larger display
// canvas
//
// Emits: SignalDraw, Argv=[Object instance, canvas]
func (o *CWidget) Draw() enums.EventFlag {
	if surface, err := memphis.GetSurface(o.ObjectID()); err != nil {
		o.LogErr(err)
	} else if o.IsDrawable() {
		return o.Emit(SignalDraw, o, surface)
	}
	return o.Emit(SignalDraw, o, nil)
}

// Emits the mnemonic-activate signal. The default handler for this
// signal activates the widget if group_cycling is FALSE, and just grabs the
// focus if group_cycling is TRUE.
// Parameters:
// 	groupCycling	TRUE if there are other widgets with the same mnemonic
// Returns:
// 	TRUE if the signal has been handled
// func (w *CWidget) MnemonicActivate(groupCycling bool) (value bool) {
// 	return false
// }

// Installs a style property on a widget class. The parser for the style
// property is determined by the value type of pspec .
// Parameters:
// 	klass	a WidgetClass
// 	pspec	the GParamSpec for the property
// func (w *CWidget) ClassInstallStyleProperty(pspec GParamSpec) {}

// Installs a style property on a widget class.
// Parameters:
// 	klass	a WidgetClass
// 	pspec	the GParamSpec for the style property
// 	parser	the parser for the style property
// func (w *CWidget) ClassInstallStylePropertyParser(pspec GParamSpec, parser RcPropertyParser) {}

// Finds a style property of a widget class by name.
// Parameters:
// 	klass	a WidgetClass
// 	propertyName	the name of the style property to find
// 	returns	the GParamSpec of the style property or
// NULL if class
// has no style property with that name.
// func (w *CWidget) ClassFindStyleProperty(propertyName string) (value GParamSpec) {
// 	return nil
// }

// Returns all style properties of a widget class.
// Parameters:
// 	klass	a WidgetClass
// 	nProperties	location to return the number of style properties found
// 	returns	an newly
// allocated array of GParamSpec*. The array must be freed with
// g_free.
// func (w *CWidget) ClassListStyleProperties(nProperties int) (value GParamSpec) {
// 	return nil
// }

// Computes the intersection of a widget 's area and region , returning the
// intersection. The result may be empty, use RegionEmpty to check.
// Parameters:
// 	region	a Region, in the same coordinate system as
// widget->allocation
// . That is, relative to widget->window
// for NO_WINDOW widgets; relative to the parent window
// of widget->window
// for widgets with their own window.
// 	returns	A newly allocated region holding the intersection of widget
// and region
// . The coordinates of the return value are
// relative to widget->window
// for NO_WINDOW widgets, and
// relative to the parent window of widget->window
// for
// widgets with their own window.
func (w *CWidget) RegionIntersect(region ptypes.Region) (value ptypes.Region) {
	return ptypes.Region{}
}

// Very rarely-used function. This function is used to emit an expose event
// signals on a widget. This function is not normally used directly. The only
// time it is used is when propagating an expose event to a child NO_WINDOW
// widget, and that is normally done using ContainerPropagateExpose.
// If you want to force an area of a window to be redrawn, use
// WindowInvalidateRect or WindowInvalidateRegion. To cause
// the redraw to be done immediately, follow that call with a call to
// WindowProcessUpdates.
// Parameters:
// 	event	a expose Event
// Returns:
// 	return from the event signal emission (TRUE if the event was
// 	handled)
func (w *CWidget) SendExpose(event cdk.Event) (value int) {
	return 0
}

// Sends the focus change event to widget This function is not meant to be
// used by applications. The only time it should be used is when it is
// necessary for a Widget to assign focus to a widget that is semantically
// owned by the first widget even though it's not a direct child - for
// instance, a search entry in a floating window similar to the quick search
// in TreeView. An example of its usage is:
// Parameters:
// 	event	a Event of type GDK_FOCUS_CHANGE
// Returns:
// 	the return value from the event signal emission: TRUE if the
// 	event was handled, and FALSE otherwise
func (w *CWidget) SendFocusChange(event cdk.Event) (value bool) {
	return false
}

// Gets the values of a multiple style properties of widget .
// Parameters:
// 	firstPropertyName	the name of the first property to get
// 	varargs	pairs of property names and locations to
// return the property values, starting with the location for
// first_property_name
// , terminated by NULL.
// func (w *CWidget) StyleGet(firstPropertyName string, argv ...interface{}) {}

// Gets the value of a style property of widget .
// Parameters:
// 	propertyName	the name of a style property
// 	value	location to return the property value
// func (w *CWidget) StyleGetProperty(propertyName string, value GValue) {}

// Non-vararg variant of StyleGet. Used primarily by language
// bindings.
// Parameters:
// 	firstPropertyName	the name of the first property to get
// 	varArgs	a va_list of pairs of property names and
// locations to return the property values, starting with the location
// for first_property_name
// .
// func (w *CWidget) StyleGetValist(firstPropertyName string, varArgs va_list) {}

// This function attaches the widget's Style to the widget's Window. It
// is a replacement for and should only ever be called in a derived widget's
// "realize" implementation which does not chain up to its parent class'
// "realize" implementation, because one of the parent classes (finally
// Widget) would attach the style itself.
// func (w *CWidget) StyleAttach() {}

// This function is used by custom widget implementations; if you're writing
// an app, you'd use GrabFocus to move the focus to a
// particular widget, and ContainerSetFocusChain to change the focus
// tab order. So you may want to investigate those functions instead.
// ChildFocus is called by containers as the user moves around
// the window using keyboard shortcuts. direction indicates what kind of
// motion is taking place (up, down, left, right, tab forward, tab backward).
// ChildFocus emits the focus signal; widgets override
// the default handler for this signal in order to implement appropriate
// focus behavior. The default ::focus handler for a widget should return
// TRUE if moving in direction left the focus on a focusable location inside
// that widget, and FALSE if moving in direction moved the focus outside the
// widget. If returning TRUE, widgets normally call GrabFocus
// to place the focus accordingly; if returning FALSE, they don't modify the
// current focus location. This function replaces ContainerFocus from
// CTK 1.2. It was necessary to check that the child was visible, sensitive,
// and focusable before calling ContainerFocus.
// ChildFocus returns FALSE if the widget is not currently in a
// focusable state, so there's no need for those checks.
// Parameters:
// 	direction	direction of focus movement
// Returns:
// 	TRUE if focus ended up inside widget
func (w *CWidget) ChildFocus(direction DirectionType) (value bool) {
	return false
}

// Emits a child-notify signal for the on widget . This is the analogue
// of g_object_notify for child properties.
// Parameters:
// 	childProperty	the name of a child property installed on the
// class of widget
// 's parent
func (w *CWidget) ChildNotify(childProperty string) {}

// Stops emission of child-notify signals on widget . The signals are
// queued until ThawChildNotify is called on widget . This is
// the analogue of g_object_freeze_notify for child properties.
func (w *CWidget) FreezeChildNotify() {}

// Gets the value set with SetChildVisible. If you feel a need
// to use this function, your code probably needs reorganization. This
// function is only useful for container implementations and never should be
// called by an application.
// Returns:
// 	TRUE if the widget is mapped with the parent.
func (w *CWidget) GetChildVisible() (value bool) {
	return false
}

// Returns the parent container of widget .
// Returns:
// 	the parent container of widget , or NULL.
// 	[transfer none]
func (w *CWidget) GetParent() (value Container) {
	if v, err := w.GetStructProperty(PropertyParent); err != nil {
		w.LogErr(err)
	} else {
		var ok bool
		if value, ok = v.(Container); !ok && v != nil {
			w.LogError("value stored in %v property is not a Container type: %v (%T)", PropertyParent, v, v)
		}
	}
	return
}

// Gets the settings object holding the settings (global property settings,
// RC file information, etc) used for this widget. Note that this function
// can only be called when the Widget is attached to a toplevel, since the
// settings object is specific to a particular Screen.
// Returns:
// 	the relevant Settings object.
// 	[transfer none]
// func (w *CWidget) GetSettings() (value Settings) {
// 	return nil
// }

// Returns the clipboard object for the given selection to be used with
// widget . widget must have a Display associated with it, so must be
// attached to a toplevel window.
// Parameters:
// 	selection	a Atom which identifies the clipboard
// to use. GDK_SELECTION_CLIPBOARD gives the
// default clipboard. Another common value
// is GDK_SELECTION_PRIMARY, which gives
// the primary X selection.
// Returns:
// 	the appropriate clipboard object. If no clipboard already
// 	exists, a new one will be created. Once a clipboard object has
// 	been created, it is persistent for all time.
// 	[transfer none]
// func (w *CWidget) GetClipboard(selection Atom) (value Clipboard) {
// 	return nil
// }

// Get the Display for the toplevel window associated with this widget.
// This function can only be called after the widget has been added to a
// widget hierarchy with a Window at the top. In general, you should only
// create display specific resources when a widget has been realized, and you
// should free those resources when the widget is unrealized.
// Returns:
// 	the Display for the toplevel for this widget.
// 	[transfer none]
func (w *CWidget) GetDisplay() (value cdk.Display) {
	return w.display
}

// Get the root window where this widget is located. This function can only
// be called after the widget has been added to a widget hierarchy with
// Window at the top. The root window is useful for such purposes as
// creating a popup Window associated with the window. In general, you
// should only create display specific resources when a widget has been
// realized, and you should free those resources when the widget is
// unrealized.
// Returns:
// 	the Window root window for the toplevel for this widget.
// 	[transfer none]
func (w *CWidget) GetRootWindow() (value Window) {
	return nil
}

// Get the Screen from the toplevel window associated with this widget.
// This function can only be called after the widget has been added to a
// widget hierarchy with a Window at the top. In general, you should only
// create screen specific resources when a widget has been realized, and you
// should free those resources when the widget is unrealized.
// Returns:
// 	the Screen for the toplevel for this widget.
// 	[transfer none]
func (w *CWidget) GetScreen() (value cdk.Display) {
	return nil
}

// Checks whether there is a Screen is associated with this widget. All
// toplevel widgets have an associated screen, and all widgets added into a
// hierarchy with a toplevel window at the top.
// Returns:
// 	TRUE if there is a Screen associcated with the widget.
func (w *CWidget) HasScreen() (value bool) {
	return false
}

// Gets the size request that was explicitly set for the widget using
// SetSizeRequest. A value of -1 stored in width or height
// indicates that that dimension has not been set explicitly and the natural
// requisition of the widget will be used intead. See
// SetSizeRequest. To get the size a widget will actually use,
// call SizeRequest instead of this function.
// Parameters:
// 	width	return location for width, or NULL.
// 	height	return location for height, or NULL.
func (w *CWidget) GetSizeRequest() (width, height int) {
	var err error
	if width, err = w.GetIntProperty(PropertyWidthRequest); err != nil {
		w.LogErr(err)
	}
	if height, err = w.GetIntProperty(PropertyHeightRequest); err != nil {
		w.LogErr(err)
	}
	return
}

// Returns the currently requested size
func (w *CWidget) SizeRequest() ptypes.Rectangle {
	return ptypes.MakeRectangle(w.GetSizeRequest())
}

// Sets whether widget should be mapped along with its when its parent is
// mapped and widget has been shown with Show. The child
// visibility can be set for widget before it is added to a container with
// SetParent, to avoid mapping children unnecessary before
// immediately unmapping them. However it will be reset to its default state
// of TRUE when the widget is removed from a container. Note that changing
// the child visibility of a widget does not queue a resize on the widget.
// Most of the time, the size of a widget is computed from all visible
// children, whether or not they are mapped. If this is not the case, the
// container can queue a resize itself. This function is only useful for
// container implementations and never should be called by an application.
// Parameters:
// 	isVisible	if TRUE, widget
// should be mapped along with its parent.
func (w *CWidget) SetChildVisible(isVisible bool) {}

// Sets the minimum size of a widget; that is, the widget's size request will
// be width by height . You can use this function to force a widget to be
// either larger or smaller than it normally would be. In most cases,
// WindowSetDefaultSize is a better choice for toplevel windows than
// this function; setting the default size will still allow users to shrink
// the window. Setting the size request will force them to leave the window
// at least as large as the size request. When dealing with window sizes,
// WindowSetGeometryHints can be a useful function as well. Note the
// inherent danger of setting any fixed size - themes, translations into
// other languages, different fonts, and user action can all change the
// appropriate size for a given widget. So, it's basically impossible to
// hardcode a size that will always be correct. The size request of a widget
// is the smallest size a widget can accept while still functioning well and
// drawing itself correctly. However in some strange cases a widget may be
// allocated less than its requested size, and in many cases a widget may be
// allocated more space than it requested. If the size request in a given
// direction is -1 (unset), then the "natural" size request of the widget
// will be used instead. Widgets can't actually be allocated a size less than
// 1 by 1, but you can pass 0,0 to this function to mean "as small as
// possible."
// Parameters:
// 	width	width widget should request, or -1 to unset
// 	height	height widget should request, or -1 to unset
//
// Emits: SignalSetSizeRequest, Argv=[Widget instance, given size]
func (w *CWidget) SetSizeRequest(width, height int) {
	if f := w.Emit(SignalSetSizeRequest, w, ptypes.MakeRectangle(width, height)); f == enums.EVENT_PASS {
		if err := w.SetIntProperty(PropertyWidthRequest, width); err != nil {
			w.LogErr(err)
		}
		if err := w.SetIntProperty(PropertyHeightRequest, height); err != nil {
			w.LogErr(err)
		}
	}
}

// Reverts the effect of a previous call to FreezeChildNotify.
// This causes all queued child-notify signals on widget to be emitted.
func (w *CWidget) ThawChildNotify() {}

// Sets the no-show-all property, which determines whether calls to
// ShowAll and HideAll will affect this widget.
// This is mostly for use in constructing widget hierarchies with externally
// controlled visibility, see UIManager.
// Parameters:
// 	noShowAll	the new value for the "no-show-all" property
func (w *CWidget) SetNoShowAll(noShowAll bool) {
	if err := w.SetBoolProperty(PropertyNoShowAll, noShowAll); err != nil {
		w.LogErr(err)
	}
}

// Returns the current value of the Widget:no-show-all property, which
// determines whether calls to ShowAll and
// HideAll will affect this widget.
// Returns:
// 	the current value of the "no-show-all" property.
func (w *CWidget) GetNoShowAll() (value bool) {
	var err error
	if value, err = w.GetBoolProperty(PropertyNoShowAll); err != nil {
		w.LogErr(err)
	}
	return
}

// Adds a widget to the list of mnemonic labels for this widget. (See
// ListMnemonicLabels). Note the list of mnemonic labels for
// the widget is cleared when the widget is destroyed, so the caller must
// make sure to update its internal state at this point as well, by using a
// connection to the destroy signal or a weak notifier.
func (w *CWidget) AddMnemonicLabel(label Widget) {}

// Removes a widget from the list of mnemonic labels for this widget. (See
// ListMnemonicLabels). The widget must have previously been
// added to the list with AddMnemonicLabel.
func (w *CWidget) RemoveMnemonicLabel(label Widget) {}

// Whether widget can rely on having its alpha channel drawn correctly. On
// X11 this function returns whether a compositing manager is running for
// widget 's screen. Please note that the semantics of this call will change
// in the future if used on a widget that has a composited window in its
// hierarchy (as set by WindowSetComposited).
// Returns:
// 	TRUE if the widget can rely on its alpha channel being drawn
// 	correctly.
func (w *CWidget) IsComposited() (value bool) {
	return false
}

// Notifies the user about an input-related error on this widget. If the
// gtk-error-bell setting is TRUE, it calls WindowBeep,
// otherwise it does nothing. Note that the effect of WindowBeep can
// be configured in many ways, depending on the windowing backend and the
// desktop environment or window manager that is used.
func (w *CWidget) ErrorBell() {}

// This function should be called whenever keyboard navigation within a
// single widget hits a boundary. The function emits the keynav-failed
// signal on the widget and its return value should be interpreted in a way
// similar to the return value of ChildFocus: When TRUE is
// returned, stay in the widget, the failed keyboard navigation is Ok and/or
// there is nowhere we can/should move the focus to. When FALSE is returned,
// the caller should continue with keyboard navigation outside the widget,
// e.g. by calling ChildFocus on the widget's toplevel. The
// default ::keynav-failed handler returns TRUE for GTK_DIR_TAB_FORWARD and
// GTK_DIR_TAB_BACKWARD. For the other values of DirectionType, it looks
// at the gtk-keynav-cursor-only setting and returns FALSE if the
// setting is TRUE. This way the entire user interface becomes
// cursor-navigatable on input devices such as mobile phones which only have
// cursor keys but no tab key. Whenever the default handler returns TRUE, it
// also calls ErrorBell to notify the user of the failed
// keyboard navigation. A use case for providing an own implementation of
// ::keynav-failed (either by connecting to it or by overriding it) would be
// a row of Entry widgets where the user should be able to navigate the
// entire row with the cursor keys, as e.g. known from user interfaces that
// require entering license keys.
// Parameters:
// 	direction	direction of focus movement
// Returns:
// 	TRUE if stopping keyboard navigation is fine, FALSE if the
// 	emitting widget should try to handle the keyboard navigation
// 	attempt in its parent container(s).
func (w *CWidget) KeynavFailed(direction DirectionType) (value bool) {
	return false
}

// Gets the contents of the tooltip for widget .
// Returns:
// 	the tooltip text, or NULL. You should free the returned string
// 	with g_free when done.
func (w *CWidget) GetTooltipMarkup() (value string) {
	var err error
	if value, err = w.GetStringProperty(PropertyTooltipMarkup); err != nil {
		w.LogErr(err)
	}
	return
}

// Sets markup as the contents of the tooltip, which is marked up with the
// Tango text markup language. This function will take care of setting
// Widget:has-tooltip to TRUE and of the default handler for the
// Widget::query-tooltip signal. See also the Widget:tooltip-markup
// property and TooltipSetMarkup.
// Parameters:
// 	markup	the contents of the tooltip for widget
// , or NULL.
func (w *CWidget) SetTooltipMarkup(markup string) {
	if err := w.SetStringProperty(PropertyTooltipMarkup, markup); err != nil {
		w.LogErr(err)
	}
}

// Gets the contents of the tooltip for widget .
// Returns:
// 	the tooltip text, or NULL. You should free the returned string
// 	with g_free when done.
func (w *CWidget) GetTooltipText() (value string) {
	var err error
	if value, err = w.GetStringProperty(PropertyTooltipText); err != nil {
		w.LogErr(err)
	}
	return
}

// Sets text as the contents of the tooltip. This function will take care of
// setting Widget:has-tooltip to TRUE and of the default handler for the
// Widget::query-tooltip signal. See also the Widget:tooltip-text
// property and TooltipSetText.
// Parameters:
// 	text	the contents of the tooltip for widget
//
func (w *CWidget) SetTooltipText(text string) {
	if err := w.SetStringProperty(PropertyTooltipText, text); err != nil {
		w.LogErr(err)
	}
}

// Returns the Window of the current tooltip. This can be the Window
// created by default, or the custom tooltip window set using
// SetTooltipWindow.
// Returns:
// 	The Window of the current tooltip.
// 	[transfer none]
func (w *CWidget) GetTooltipWindow() (value Window) {
	return nil
}

// Replaces the default, usually yellow, window used for displaying tooltips
// with custom_window . CTK will take care of showing and hiding
// custom_window at the right moment, to behave likewise as the default
// tooltip window. If custom_window is NULL, the default tooltip window will
// be used. If the custom window should have the default theming it needs to
// have the name "gtk-tooltip", see SetName.
// Parameters:
// 	customWindow	a Window, or NULL.
func (w *CWidget) SetTooltipWindow(customWindow Window) {}

// Returns the current value of the has-tooltip property. See
// Widget:has-tooltip for more information.
// Returns:
// 	current value of has-tooltip on widget .
func (w *CWidget) GetHasTooltip() (value bool) {
	var err error
	if value, err = w.GetBoolProperty(PropertyHasTooltip); err != nil {
		w.LogErr(err)
	}
	return
}

// Sets the has-tooltip property on widget to has_tooltip . See
// Widget:has-tooltip for more information.
// Parameters:
// 	hasTooltip	whether or not widget
// has a tooltip.
func (w *CWidget) SetHasTooltip(hasTooltip bool) {
	if err := w.SetBoolProperty(PropertyHasTooltip, hasTooltip); err != nil {
		w.LogErr(err)
	}
}

// Triggers a tooltip query on the display where the toplevel of widget is
// located. See TooltipTriggerTooltipQuery for more information.
func (w *CWidget) TriggerTooltipQuery() {}

// Returns the widget's window if it is realized, NULL otherwise
// Returns:
// 	widget 's window.
// 	[transfer none]
// Returns the Window instance associated with this Widget instance, nil
// otherwise
func (w *CWidget) GetWindow() (window Window) {
	if v, err := w.GetStructProperty(PropertyWindow); err != nil {
		w.LogErr(err)
		return nil
	} else {
		var ok bool
		if window, ok = v.(Window); !ok && v != nil {
			w.LogError("value stored in %v property is not of Window type: %v (%T)", PropertyWindow, v, v)
			return nil
		}
	}
	if window == nil && w.parent != nil {
		p := w.parent
		for {
			if wc, ok := p.(Window); ok {
				window = wc
				break
			} else if pc, ok := p.(Container); ok {
				p = pc.GetParent()
			} else {
				break
			}
		}
		if pw, ok := p.(Window); ok && window == nil {
			window = pw
		}
	}
	if dm := w.GetDisplay(); dm != nil {
		if window != nil {
			if overlay := dm.GetWindowTopOverlay(window.ObjectID()); overlay != nil {
				if wo, ok := overlay.(Window); ok {
					window = wo
				}
			}
		} else { // window == nil
			if tw := dm.ActiveWindow(); tw != nil {
				window, _ = tw.(Window)
			}
		}
	}
	return
}

// Retrieves the widget's allocation.
// Parameters:
// 	allocation	a pointer to a Allocation to copy to.
// func (w *CWidget) GetAllocation(allocation ptypes.Rectangle) {}

// Sets the widget's allocation. This should not be used directly, but from
// within a widget's size_allocate method.
// Parameters:
// 	allocation	a pointer to a Allocation to copy from
// func (w *CWidget) SetAllocation(allocation ptypes.Rectangle) {}

// Determines whether the application intends to draw on the widget in an
// expose-event handler. See SetAppPaintable
// Returns:
// 	TRUE if the widget is app paintable
func (w *CWidget) GetAppPaintable() (value bool) {
	var err error
	if value, err = w.GetBoolProperty(PropertyAppPaintable); err != nil {
		w.LogErr(err)
	}
	return
}

// Determines whether widget can be a default widget. See
// SetCanDefault.
// Returns:
// 	TRUE if widget can be a default widget, FALSE otherwise
func (w *CWidget) GetCanDefault() (value bool) {
	var err error
	if value, err = w.GetBoolProperty(PropertyCanDefault); err != nil {
		w.LogErr(err)
	}
	return
}

// Specifies whether widget can be a default widget. See
// GrabDefault for details about the meaning of "default".
// Parameters:
// 	canDefault	whether or not widget
// can be a default widget.
func (w *CWidget) SetCanDefault(canDefault bool) {
	if err := w.SetBoolProperty(PropertyCanDefault, canDefault); err != nil {
		w.LogErr(err)
	}
}

// Determines whether widget can own the input focus. See
// SetCanFocus.
// Returns:
// 	TRUE if widget can own the input focus, FALSE otherwise
func (w *CWidget) GetCanFocus() (value bool) {
	var err error
	if value, err = w.GetBoolProperty(PropertyCanFocus); err != nil {
		w.LogErr(err)
	}
	return
}

// Specifies whether widget can own the input focus. See
// GrabFocus for actually setting the input focus on a widget.
// Parameters:
// 	canFocus	whether or not widget
// can own the input focus.
func (w *CWidget) SetCanFocus(canFocus bool) {
	if err := w.SetBoolProperty(PropertyCanFocus, canFocus); err != nil {
		w.LogErr(err)
	}
}

// Determines whether the widget is double buffered. See
// SetDoubleBuffered
// Returns:
// 	TRUE if the widget is double buffered
func (w *CWidget) GetDoubleBuffered() (value bool) {
	var err error
	if value, err = w.GetBoolProperty(PropertyDoubleBuffered); err != nil {
		w.LogErr(err)
	}
	return
}

// Determines whether widget has a Window of its own. See
// SetHasWindow.
// Returns:
// 	TRUE if widget has a window, FALSE otherwise
func (w *CWidget) GetHasWindow() (value bool) {
	return false
}

// Specifies whether widget has a Window of its own. Note that all
// realized widgets have a non-NULL "window" pointer (GetWindow
// never returns a NULL window when a widget is realized), but for many of
// them it's actually the Window of one of its parent widgets. Widgets
// that do not create a window for themselves in Widget::realize must
// announce this by calling this function with has_window = FALSE. This
// function should only be called by widget implementations, and they should
// call it in their init function.
// Parameters:
// 	hasWindow	whether or not widget
// has a window.
func (w *CWidget) SetHasWindow(hasWindow bool) {}

// Returns the widget's sensitivity (in the sense of returning the value that
// has been set using SetSensitive). The effective sensitivity
// of a widget is however determined by both its own and its parent widget's
// sensitivity. See IsSensitive.
// Returns:
// 	TRUE if the widget is sensitive
func (w *CWidget) GetSensitive() (value bool) {
	var err error
	if value, err = w.GetBoolProperty(PropertySensitive); err != nil {
		w.LogErr(err)
	}
	return
}

// Returns the widget's effective sensitivity, which means it is sensitive
// itself and also its parent widget is sensitive
// Returns:
// 	TRUE if the widget is effectively sensitive
func (w *CWidget) IsSensitive() bool {
	if w.HasState(StateInsensitive) {
		return false
	}
	if parent := w.GetParent(); parent != nil {
		if parent.HasState(StateInsensitive) {
			return false
		}
	}
	return true
}

// Returns the widget's state. See SetState.
// Returns:
// 	the state of the widget.
func (w *CWidget) GetState() (value StateType) {
	return w.state
}

// Determines whether the widget is visible. Note that this doesn't take into
// account whether the widget's parent is also visible or the widget is
// obscured in any way. See SetVisible.
// Returns:
// 	TRUE if the widget is visible
func (w *CWidget) GetVisible() (value bool) {
	var err error
	if value, err = w.GetBoolProperty(PropertyVisible); err != nil {
		w.LogErr(err)
	}
	return
}

// Sets the visibility state of widget . Note that setting this to TRUE
// doesn't mean the widget is actually viewable, see
// GetVisible. This function simply calls Show or
// Hide but is nicer to use when the visibility of the widget
// depends on some condition.
// Parameters:
// 	visible	whether the widget should be shown or not
func (w *CWidget) SetVisible(visible bool) {
	if err := w.SetBoolProperty(PropertyVisible, visible); err != nil {
		w.LogErr(err)
	}
}

// Determines whether widget is the current default widget within its
// toplevel. See SetCanDefault.
// Returns:
// 	TRUE if widget is the current default widget within its
// 	toplevel, FALSE otherwise
func (w *CWidget) HasDefault() (value bool) {
	return false
}

// Determines if the widget has the global input focus. See
// IsFocus for the difference between having the global input
// focus, and only having the focus within a toplevel.
// Returns:
// 	TRUE if the widget has the global input focus.
func (w *CWidget) HasFocus() (value bool) {
	return false
}

// Determines whether the widget is currently grabbing events, so it is the
// only widget receiving input events (keyboard and mouse). See also
// GrabAdd.
// Returns:
// 	TRUE if the widget is in the grab_widgets stack
func (w *CWidget) HasGrab() (value bool) {
	return false
}

// Determines if the widget style has been looked up through the rc
// mechanism.
// Returns:
// 	TRUE if the widget has been looked up through the rc mechanism,
// 	FALSE otherwise.
func (w *CWidget) HasRcStyle() (value bool) {
	return false
}

// Determines whether widget can be drawn to. A widget can be drawn to if it
// is mapped and visible.
// Returns:
// 	TRUE if widget is drawable, FALSE otherwise
// Returns TRUE if the APP_PAINTABLE flag is set, FALSE otherwise
func (w *CWidget) IsDrawable() (value bool) {
	return w.HasFlags(APP_PAINTABLE)
}

// Determines whether widget is a toplevel widget. Currently only Window
// and Invisible are toplevel widgets. Toplevel widgets have no parent
// widget.
// Returns:
// 	TRUE if widget is a toplevel, FALSE otherwise
func (w *CWidget) IsToplevel() (value bool) {
	return w.HasFlags(TOPLEVEL)
}

// Sets a widget's window. This function should only be used in a widget's
// Widget::realize implementation. The window passed is usually either
// new window created with WindowNew, or the window of its parent
// widget as returned by GetParentWindow. Widgets must
// indicate whether they will create their own Window by calling
// SetHasWindow. This is usually done in the widget's init
// function.
// Parameters:
// 	window	a Window
//
// Emits: SignalSetWindow, Argv=[Widget instance, given window]
func (w *CWidget) SetWindow(window Window) {
	if f := w.Emit(SignalSetWindow, w, window); f == enums.EVENT_PASS {
		if err := w.SetStructProperty(PropertyWindow, window); err != nil {
			w.LogErr(err)
		}
	}
}

// Specifies whether widget will be treated as the default widget within its
// toplevel when it has the focus, even if another widget is the default. See
// GrabDefault for details about the meaning of "default".
// Parameters:
// 	receivesDefault	whether or not widget
// can be a default widget.
func (w *CWidget) SetReceivesDefault(receivesDefault bool) {
	if err := w.SetBoolProperty(PropertyReceivesDefault, receivesDefault); err != nil {
		w.LogErr(err)
	}
}

// Determines whether widget is alyways treated as default widget withing its
// toplevel when it has the focus, even if another widget is the default. See
// SetReceivesDefault.
// Returns:
// 	TRUE if widget acts as default widget when focussed, FALSE
// 	otherwise
func (w *CWidget) GetReceivesDefault() (value bool) {
	var err error
	if value, err = w.GetBoolProperty(PropertyReceivesDefault); err != nil {
		w.LogErr(err)
	}
	return
}

// Marks the widget as being realized. This function should only ever be
// called in a derived widget's "realize" or "unrealize" implementation.
// Parameters:
// 	realized	TRUE to mark the widget as realized
func (w *CWidget) SetRealized(realized bool) {}

// Determines whether widget is realized.
// Returns:
// 	TRUE if widget is realized, FALSE otherwise
func (w *CWidget) GetRealized() (value bool) {
	return false
}

// Marks the widget as being realized. This function should only ever be
// called in a derived widget's "map" or "unmap" implementation.
// Parameters:
// 	mapped	TRUE to mark the widget as mapped
func (w *CWidget) SetMapped(mapped bool) {}

// Whether the widget is mapped.
// Returns:
// 	TRUE if the widget is mapped, FALSE otherwise.
func (w *CWidget) GetMapped() (value bool) {
	return false
}

// Retrieves the widget's requisition. This function should only be used by
// widget implementations in order to figure whether the widget's requisition
// has actually changed after some internal state change (so that they can
// call QueueResize instead of QueueDraw).
// Normally, SizeRequest should be used.
// Parameters:
// 	requisition	a pointer to a Requisition to copy to.
// func (w *CWidget) GetRequisition(requisition Requisition) {}

// Copies a Requisition.
// Parameters:
// 	requisition	a Requisition
// Returns:
// 	a copy of requisition
// func (w *CWidget) RequisitionCopy(requisition Requisition) (value Requisition) {
// 	return nil
// }

// Frees a Requisition.
// Parameters:
// 	requisition	a Requisition
// func (w *CWidget) RequisitionFree(requisition Requisition) {}

// Returns the current theme, adjusted for Widget focus and accounting for
// any PARENT_SENSITIVE conditions. This method is primarily useful in drawable
// Widget types during the Invalidate() and Draw() stages of the Widget
// lifecycle
func (w *CWidget) GetThemeRequest() (theme paint.Theme) {
	theme = w.GetTheme()
	if (w.CanFocus() && w.IsFocused()) || w.IsParentFocused() {
		theme.Content.Normal = theme.Content.Focused
		theme.Border.Normal = theme.Border.Focused
	}
	return
}

// Set the Theme for the Widget instance. This will also refresh the requested
// theme. A request theme is a transient theme, based on the actually set theme
// and adjusted for focus. If the given theme is equivalent to the current theme
// then no action is taken. After verifying that the given theme is different,
// this method emits a set-theme signal and if the listeners return EVENT_PASS,
// the changes are applied and the Widget.Invalidate() method is called
func (w *CWidget) SetTheme(theme paint.Theme) {
	if theme.String() != w.GetTheme().String() {
		if f := w.Emit(SignalSetTheme, w, theme); f == enums.EVENT_PASS {
			w.CObject.SetTheme(theme)
			w.Invalidate()
		}
	}
}

// Returns TRUE if the Widget has the given StateType, FALSE otherwise
func (w *CWidget) HasState(s StateType) bool {
	return w.state&s != 0
}

// Removes the given state bitmask from the Widget instance. This method emits
// an unset-state signal initially and if the listeners return EVENT_PASS, the
// change is applied
//
// Emit: SignalUnsetState, Argv=[Widget instance, given state to unset]
func (w *CWidget) UnsetState(v StateType) {
	if f := w.Emit(SignalUnsetState, w, v); f == enums.EVENT_PASS {
		w.state = w.state &^ v
	}
}

// Returns the current flags for the Widget instance
func (w *CWidget) GetFlags() WidgetFlags {
	return w.flags
}

// Returns TRUE if the Widget instance has the given flag, FALSE otherwise
func (w *CWidget) HasFlags(f WidgetFlags) bool {
	w.flagsLock.RLock()
	defer w.flagsLock.RUnlock()
	return w.flags.HasBit(f)
}

// Removes the given flags from the Widget instance. This method emits an
// unset-flags signal initially and if the listeners return EVENT_PASS, the
// change is applied
//
// Emits: SignalUnsetFlags, Argv=[Widget instance, given flags to unset]
func (w *CWidget) UnsetFlags(v WidgetFlags) {
	if f := w.Emit(SignalUnsetFlags, w, v); f == enums.EVENT_PASS {
		w.flagsLock.Lock()
		defer w.flagsLock.Unlock()
		w.flags = w.flags &^ v
	}
}

// Sets the given flags on the Widget instance. This method emits a set-flags
// signal initially and if the listeners return EVENT_PASS, the change is
// applied
//
// Emits: SignalSetFlags, Argv=[Widget instance, given flags to set]
func (w *CWidget) SetFlags(v WidgetFlags) {
	if f := w.Emit(SignalSetFlags, w, w.flags, v); f == enums.EVENT_PASS {
		w.flagsLock.Lock()
		defer w.flagsLock.Unlock()
		w.flags = w.flags | v
	}
}

// If the Widget instance is PARENT_SENSITIVE and one of it's parents are the
// focus for the associated Window, return TRUE and FALSE otherwise
func (w *CWidget) IsParentFocused() bool {
	if w.HasFlags(PARENT_SENSITIVE) {
		var lastParent Widget
		parent, _ := w.GetParent().(Widget)
		for parent != nil {
			if parent.IsFocus() {
				return true
			}
			if !parent.HasFlags(PARENT_SENSITIVE) || parent.IsToplevel() {
				// don't recurse
				break
			}
			lastParent = parent
			if parent, _ = parent.GetParent().(Widget); parent != nil {
				if lastParent.ObjectID() == parent.ObjectID() {
					break // stop
				}
			}
		}
	}
	return false
}

// Returns TRUE if the Widget instance or it's parent are the current focus of
// the associated Window
func (w *CWidget) IsFocused() bool {
	return w.IsFocus() || w.IsParentFocused()
}

// Returns TRUE if the Widget instance has the CAN_FOCUS flag, FALSE otherwise
func (w *CWidget) CanFocus() bool {
	return w.HasFlags(CAN_FOCUS)
}

// Returns TRUE if the Widget instance CanDefault() and the HAS_DEFAULT flag is
// set, returns FALSE otherwise
func (w *CWidget) IsDefault() bool {
	if w.CanDefault() {
		return w.HasFlags(HAS_DEFAULT)
	}
	return false
}

// Returns TRUE if the Widget instance IsSensitive() and the CAN_DEFAULT flag is
// set, returns FALSE otherwise
func (w *CWidget) CanDefault() bool {
	if w.IsSensitive() {
		return w.HasFlags(CAN_DEFAULT)
	}
	return false
}

// Returns TRUE if the VISIBLE flag is set, FALSE otherwise
func (w *CWidget) IsVisible() bool {
	return w.HasFlags(VISIBLE)
}

func (w *CWidget) HasEventFocus() bool {
	if window := w.GetWindow(); window != nil {
		if ef := window.GetEventFocus(); ef != nil {
			if wef, ok := ef.(Widget); ok {
				if wef.ObjectID() == w.ObjectID() {
					return true
				}
			}
		}
	}
	return false
}

func (w *CWidget) GrabEventFocus() {
	if window := w.GetWindow(); window != nil {
		if f := w.Emit(SignalGrabEventFocus, w, window); f == enums.EVENT_PASS {
			window.SetEventFocus(w)
		}
	}
}

func (w *CWidget) ReleaseEventFocus() {
	if window := w.GetWindow(); window != nil {
		if ef := window.GetEventFocus(); ef != nil {
			if wef, ok := ef.(Widget); ok && wef.ObjectID() == w.ObjectID() {
				if f := w.Emit(SignalReleaseEventFocus, w, window); f == enums.EVENT_PASS {
					window.SetEventFocus(nil)
				}
			}
		}
	}
}

// Returns the top-most parent in the Widget instance's parent hierarchy.
// Returns nil if the Widget has no parent container
func (w *CWidget) GetTopParent() (parent Container) {
	var ok bool
	if parent, ok = w.GetParent().(Container); !ok {
		return
	}
	for {
		if parent != nil {
			parent, _ = w.GetParent().(Container)
		} else {
			break
		}
	}
	return
}

// A wrapper around the Object.GetObjectAt() method, only returning Widget
// instance types or nil otherwise
func (w *CWidget) GetWidgetAt(p *ptypes.Point2I) Widget {
	if o := w.CObject.GetObjectAt(p); o != nil {
		if ow, ok := o.(Widget); ok && ow.IsVisible() {
			return ow
		}
	}
	return nil
}

func (w *CWidget) lostFocus(_ []interface{}, _ ...interface{}) enums.EventFlag {
	theme := w.GetTheme()
	w.SetThemeRequest(theme)
	w.Invalidate()
	return enums.EVENT_PASS
}

func (w *CWidget) gainedFocus(_ []interface{}, _ ...interface{}) enums.EventFlag {
	theme := w.GetTheme()
	theme.Content.Normal = theme.Content.Focused
	theme.Border.Normal = theme.Border.Focused
	w.SetThemeRequest(theme)
	w.Invalidate()
	return enums.EVENT_PASS
}

// Whether the application will paint directly on the widget.
// Flags: Read / Write
// Default value: FALSE
const PropertyAppPaintable cdk.Property = "app-paintable"

// Whether the widget can be the default widget.
// Flags: Read / Write
// Default value: FALSE
const PropertyCanDefault cdk.Property = "can-default"

// Whether the widget can accept the input focus.
// Flags: Read / Write
// Default value: FALSE
const PropertyCanFocus cdk.Property = "can-focus"

// Whether the widget is part of a composite widget.
// Flags: Read
// Default value: FALSE
const PropertyCompositeChild cdk.Property = "composite-child"

// Whether or not the widget is double buffered.
// Flags: Read / Write
// Default value: TRUE
const PropertyDoubleBuffered cdk.Property = "double-buffered"

// The event mask that decides what kind of Events this widget gets.
// Flags: Read / Write
// Default value: GDK_STRUCTURE_MASK
const PropertyEvents cdk.Property = "events"

// The mask that decides what kind of extension events this widget gets.
// Flags: Read / Write
// Default value: GDK_EXTENSION_EVENTS_NONE
const PropertyExtensionEvents cdk.Property = "extension-events"

// Whether the widget is the default widget.
// Flags: Read / Write
// Default value: FALSE
const PropertyHasDefault cdk.Property = "has-default"

// Whether the widget has the input focus.
// Flags: Read / Write
// Default value: FALSE
const PropertyHasFocus cdk.Property = "has-focus"

// Enables or disables the emission of query-tooltip on widget . A
// value of TRUE indicates that widget can have a tooltip, in this case the
// widget will be queried using query-tooltip to determine whether it
// will provide a tooltip or not. Note that setting this property to TRUE for
// the first time will change the event masks of the Windows of this
// widget to include leave-notify and motion-notify events. This cannot and
// will not be undone when the property is set to FALSE again.
// Flags: Read / Write
// Default value: FALSE
const PropertyHasTooltip cdk.Property = "has-tooltip"

// Override for height request of the widget, or -1 if natural request should
// be used.
// Flags: Read / Write
// Allowed values: >= -1
// Default value: -1
const PropertyHeightRequest cdk.Property = "height-request"

// Whether the widget is the focus widget within the toplevel.
// Flags: Read / Write
// Default value: FALSE
const PropertyIsFocus cdk.Property = "is-focus"

// The name of the widget.
// Flags: Read / Write
// Default value: NULL
const PropertyName cdk.Property = "name"

// Whether ShowAll should not affect this widget.
// Flags: Read / Write
// Default value: FALSE
const PropertyNoShowAll cdk.Property = "no-show-all"

// The parent widget of this widget. Must be a Container widget.
// Flags: Read / Write
const PropertyParent cdk.Property = "parent"

// If TRUE, the widget will receive the default action when it is focused.
// Flags: Read / Write
// Default value: FALSE
const PropertyReceivesDefault cdk.Property = "receives-default"

// Whether the widget responds to input.
// Flags: Read / Write
// Default value: TRUE
const PropertySensitive cdk.Property = "sensitive"

// The style of the widget, which contains information about how it will look
// (colors etc).
// Flags: Read / Write
const PropertyStyle cdk.Property = "style"

// Sets the text of tooltip to be the given string, which is marked up with
// the Tango text markup language. Also see TooltipSetMarkup. This is
// a convenience property which will take care of getting the tooltip shown
// if the given string is not NULL: has-tooltip will automatically be
// set to TRUE and there will be taken care of query-tooltip in the
// default signal handler.
// Flags: Read / Write
// Default value: NULL
const PropertyTooltipMarkup cdk.Property = "tooltip-markup"

// Sets the text of tooltip to be the given string. Also see
// TooltipSetText. This is a convenience property which will take
// care of getting the tooltip shown if the given string is not NULL:
// has-tooltip will automatically be set to TRUE and there will be
// taken care of query-tooltip in the default signal handler.
// Flags: Read / Write
// Default value: NULL
const PropertyTooltipText cdk.Property = "tooltip-text"

// Whether the widget is visible.
// Flags: Read / Write
// Default value: FALSE
const PropertyVisible cdk.Property = "visible"

// Override for width request of the widget, or -1 if natural request should
// be used.
// Flags: Read / Write
// Allowed values: >= -1
// Default value: -1
const PropertyWidthRequest cdk.Property = "width-request"

// The widget's window if it is realized, NULL otherwise.
// Flags: Read
const PropertyWindow cdk.Property = "window"

const SignalAccelClosuresChanged cdk.Signal = "accel-closures-changed"

// The ::button-press-event signal will be emitted when a button (typically
// from a mouse) is pressed. To receive this signal, the Window associated
// to the widget needs to enable the GDK_BUTTON_PRESS_MASK mask. This signal
// will be sent to the grab widget if there is one.
const SignalButtonPressEvent cdk.Signal = "button-press-event"

// The ::button-release-event signal will be emitted when a button (typically
// from a mouse) is released. To receive this signal, the Window
// associated to the widget needs to enable the GDK_BUTTON_RELEASE_MASK mask.
// This signal will be sent to the grab widget if there is one.
const SignalButtonReleaseEvent cdk.Signal = "button-release-event"

// Determines whether an accelerator that activates the signal identified by
// signal_id can currently be activated. This signal is present to allow
// applications and derived widgets to override the default Widget
// handling for determining whether an accelerator can be activated.
const SignalCanActivateAccel cdk.Signal = "can-activate-accel"

// The ::child-notify signal is emitted for each changed on an object. The
// signal's detail holds the property name.
// Listener function arguments:
//      pspec GParamSpec        the GParamSpec of the changed child property
const SignalChildNotify cdk.Signal = "child-notify"

// The ::client-event will be emitted when the widget 's window receives a
// message (via a ClientMessage event) from another application.
const SignalClientEvent cdk.Signal = "client-event"

// The ::composited-changed signal is emitted when the composited status of
// widget s screen changes. See ScreenIsComposited.
const SignalCompositedChanged cdk.Signal = "composited-changed"

// The ::configure-event signal will be emitted when the size, position or
// stacking of the widget 's window has changed. To receive this signal, the
// Window associated to the widget needs to enable the GDK_STRUCTURE_MASK
// mask. GDK will enable this mask automatically for all new windows.
const SignalConfigureEvent cdk.Signal = "configure-event"

// Emitted when a redirected window belonging to widget gets drawn into. The
// region/area members of the event shows what area of the redirected
// drawable was drawn into.
const SignalDamageEvent cdk.Signal = "damage-event"

// The ::delete-event signal is emitted if a user requests that a toplevel
// window is closed. The default handler for this signal destroys the window.
// Connecting HideOnDelete to this signal will cause the
// window to be hidden instead, so that it can later be shown again without
// reconstructing it.
const SignalDeleteEvent cdk.Signal = "delete-event"

// The ::destroy-event signal is emitted when a Window is destroyed. You
// rarely get this signal, because most widgets disconnect themselves from
// their window before they destroy it, so no widget owns the window at
// destroy time. To receive this signal, the Window associated to the
// widget needs to enable the GDK_STRUCTURE_MASK mask. GDK will enable this
// mask automatically for all new windows.
const SignalDestroyEvent cdk.Signal = "destroy-event"

// The ::direction-changed signal is emitted when the text direction of a
// widget changes.
// Listener function arguments:
//      previousDirection TextDirection the previous text direction of widget
const SignalDirectionChanged cdk.Signal = "direction-changed"

// The ::drag-begin signal is emitted on the drag source when a drag is
// started. A typical reason to connect to this signal is to set up a custom
// drag icon with DragSourceSetIcon. Note that some widgets set up a
// drag icon in the default handler of this signal, so you may have to use
// g_signal_connect_after to override what the default handler did.
// Listener function arguments:
//      dragContext DragContext the drag context
const SignalDragBegin cdk.Signal = "drag-begin"

// The ::drag-data-delete signal is emitted on the drag source when a drag
// with the action GDK_ACTION_MOVE is successfully completed. The signal
// handler is responsible for deleting the data that has been dropped. What
// "delete" means depends on the context of the drag operation.
// Listener function arguments:
//      dragContext DragContext the drag context
const SignalDragDataDelete cdk.Signal = "drag-data-delete"

// The ::drag-data-get signal is emitted on the drag source when the drop
// site requests the data which is dragged. It is the responsibility of the
// signal handler to fill data with the data in the format which is indicated
// by info . See SelectionDataSet and SelectionDataSetText.
// Listener function arguments:
//      dragContext DragContext the drag context
//      data SelectionData      the GtkSelectionData to be filled with the dragged data
//      info int        the info that has been registered with the target in the GtkTargetList
//      time int        the timestamp at which the data was requested
const SignalDragDataGet cdk.Signal = "drag-data-get"

// The ::drag-data-received signal is emitted on the drop site when the
// dragged data has been received. If the data was received in order to
// determine whether the drop will be accepted, the handler is expected to
// call DragStatus and not finish the drag. If the data was received
// in response to a drag-drop signal (and this is the last target to be
// received), the handler for this signal is expected to process the received
// data and then call DragFinish, setting the success parameter
// depending on whether the data was processed successfully. The handler may
// inspect and modify drag_context->action before calling DragFinish,
// e.g. to implement GDK_ACTION_ASK as shown in the following example:
// Listener function arguments:
//      dragContext DragContext the drag context
//      x int   where the drop happened
//      y int   where the drop happened
//      data SelectionData      the received data
//      info int        the info that has been registered with the target in the GtkTargetList
//      time int        the timestamp at which the data was received
const SignalDragDataReceived cdk.Signal = "drag-data-received"

// The ::drag-drop signal is emitted on the drop site when the user drops the
// data onto the widget. The signal handler must determine whether the cursor
// position is in a drop zone or not. If it is not in a drop zone, it returns
// FALSE and no further processing is necessary. Otherwise, the handler
// returns TRUE. In this case, the handler must ensure that DragFinish
// is called to let the source know that the drop is done. The call to
// DragFinish can be done either directly or in a
// drag-data-received handler which gets triggered by calling
// DragGetData to receive the data for one or more of the supported
// targets.
const SignalDragDrop cdk.Signal = "drag-drop"

// The ::drag-end signal is emitted on the drag source when a drag is
// finished. A typical reason to connect to this signal is to undo things
// done in drag-begin.
// Listener function arguments:
//      dragContext DragContext the drag context
const SignalDragEnd cdk.Signal = "drag-end"

// The ::drag-failed signal is emitted on the drag source when a drag has
// failed. The signal handler may hook custom code to handle a failed DND
// operation based on the type of error, it returns TRUE is the failure has
// been already handled (not showing the default "drag operation failed"
// animation), otherwise it returns FALSE.
const SignalDragFailed cdk.Signal = "drag-failed"

// The ::drag-leave signal is emitted on the drop site when the cursor leaves
// the widget. A typical reason to connect to this signal is to undo things
// done in drag-motion, e.g. undo highlighting with
// DragUnhighlight
// Listener function arguments:
//      dragContext DragContext the drag context
//      time int        the timestamp of the motion event
const SignalDragLeave cdk.Signal = "drag-leave"

// The drag-motion signal is emitted on the drop site when the user moves the
// cursor over the widget during a drag. The signal handler must determine
// whether the cursor position is in a drop zone or not. If it is not in a
// drop zone, it returns FALSE and no further processing is necessary.
// Otherwise, the handler returns TRUE. In this case, the handler is
// responsible for providing the necessary information for displaying
// feedback to the user, by calling DragStatus. If the decision
// whether the drop will be accepted or rejected can't be made based solely
// on the cursor position and the type of the data, the handler may inspect
// the dragged data by calling DragGetData and defer the
// DragStatus call to the drag-data-received handler. Note that
// you cannot not pass GTK_DEST_DEFAULT_DROP, GTK_DEST_DEFAULT_MOTION or
// GTK_DEST_DEFAULT_ALL to DragDestSet when using the drag-motion
// signal that way. Also note that there is no drag-enter signal. The drag
// receiver has to keep track of whether he has received any drag-motion
// signals since the last drag-leave and if not, treat the drag-motion
// signal as an "enter" signal. Upon an "enter", the handler will typically
// highlight the drop site with DragHighlight.
const SignalDragMotion cdk.Signal = "drag-motion"

// The ::enter-notify-event will be emitted when the pointer enters the
// widget 's window. To receive this signal, the Window associated to the
// widget needs to enable the GDK_ENTER_NOTIFY_MASK mask. This signal will be
// sent to the grab widget if there is one.
const SignalEnterNotifyEvent cdk.Signal = "enter-notify-event"

// The CTK main loop will emit three signals for each GDK event delivered to
// a widget: one generic ::event signal, another, more specific, signal that
// matches the type of event delivered (e.g. key-press-event) and
// finally a generic event-after signal.
const SignalEvent cdk.Signal = "event"

// After the emission of the event signal and (optionally) the second
// more specific signal, ::event-after will be emitted regardless of the
// previous two signals handlers return values.
// Listener function arguments:
//      event Event     the GdkEvent which triggered this signal
const SignalEventAfter cdk.Signal = "event-after"

// The ::expose-event signal is emitted when an area of a previously obscured
// Window is made visible and needs to be redrawn. GTK_NO_WINDOW widgets
// will get a synthesized event from their parent widget. To receive this
// signal, the Window associated to the widget needs to enable the
// GDK_EXPOSURE_MASK mask. Note that the ::expose-event signal has been
// replaced by a ::draw signal in CTK 3. The CTK 3 migration guide for
// hints on how to port from ::expose-event to ::draw.
const SignalExposeEvent cdk.Signal = "expose-event"

const SignalFocus cdk.Signal = "focus"

// The ::focus-in-event signal will be emitted when the keyboard focus enters
// the widget 's window. To receive this signal, the Window associated to
// the widget needs to enable the GDK_FOCUS_CHANGE_MASK mask.
const SignalFocusInEvent cdk.Signal = "focus-in-event"

// The ::focus-out-event signal will be emitted when the keyboard focus
// leaves the widget 's window. To receive this signal, the Window
// associated to the widget needs to enable the GDK_FOCUS_CHANGE_MASK mask.
const SignalFocusOutEvent cdk.Signal = "focus-out-event"

// Emitted when a pointer or keyboard grab on a window belonging to widget
// gets broken. On X11, this happens when the grab window becomes unviewable
// (i.e. it or one of its ancestors is unmapped), or if the same application
// grabs the pointer or keyboard again.
const SignalGrabBrokenEvent cdk.Signal = "grab-broken-event"

const SignalGrabFocus cdk.Signal = "grab-focus"

// The ::grab-notify signal is emitted when a widget becomes shadowed by a
// CTK grab (not a pointer or keyboard grab) on another widget, or when it
// becomes unshadowed due to a grab being removed. A widget is shadowed by a
// GrabAdd when the topmost grab widget in the grab stack of its
// window group is not its ancestor.
// Listener function arguments:
//      wasGrabbed bool FALSE if the widget becomes shadowed, TRUE if it becomes unshadowed
const SignalGrabNotify cdk.Signal = "grab-notify"

const SignalHide cdk.Signal = "hide"

// The ::hierarchy-changed signal is emitted when the anchored state of a
// widget changes. A widget is anchored when its toplevel ancestor is a
// Window. This signal is emitted when a widget changes from un-anchored
// to anchored or vice-versa.
const SignalHierarchyChanged cdk.Signal = "hierarchy-changed"

// The ::key-press-event signal is emitted when a key is pressed. To receive
// this signal, the Window associated to the widget needs to enable the
// GDK_KEY_PRESS_MASK mask. This signal will be sent to the grab widget if
// there is one.
const SignalKeyPressEvent cdk.Signal = "key-press-event"

// The ::key-release-event signal is emitted when a key is pressed. To
// receive this signal, the Window associated to the widget needs to
// enable the GDK_KEY_RELEASE_MASK mask. This signal will be sent to the grab
// widget if there is one.
const SignalKeyReleaseEvent cdk.Signal = "key-release-event"

// Gets emitted if keyboard navigation fails. See KeynavFailed
// for details.
const SignalKeynavFailed cdk.Signal = "keynav-failed"

// The ::leave-notify-event will be emitted when the pointer leaves the
// widget 's window. To receive this signal, the Window associated to the
// widget needs to enable the GDK_LEAVE_NOTIFY_MASK mask. This signal will be
// sent to the grab widget if there is one.
const SignalLeaveNotifyEvent cdk.Signal = "leave-notify-event"

const SignalMap cdk.Signal = "map"

// The ::map-event signal will be emitted when the widget 's window is
// mapped. A window is mapped when it becomes visible on the screen. To
// receive this signal, the Window associated to the widget needs to
// enable the GDK_STRUCTURE_MASK mask. GDK will enable this mask
// automatically for all new windows.
const SignalMapEvent cdk.Signal = "map-event"

const SignalMnemonicActivate cdk.Signal = "mnemonic-activate"

// The ::motion-notify-event signal is emitted when the pointer moves over
// the widget's Window. To receive this signal, the Window associated
// to the widget needs to enable the GDK_POINTER_MOTION_MASK mask. This
// signal will be sent to the grab widget if there is one.
const SignalMotionNotifyEvent cdk.Signal = "motion-notify-event"

// Listener function arguments:
//      direction DirectionType
const SignalMoveFocus cdk.Signal = "move-focus"

// The ::no-expose-event will be emitted when the widget 's window is drawn
// as a copy of another Drawable (with DrawDrawable or
// WindowCopyArea) which was completely unobscured. If the source
// window was partially obscured EventExpose events will be generated for
// those areas.
const SignalNoExposeEvent cdk.Signal = "no-expose-event"

// The ::parent-set signal is emitted when a new parent has been set on a
// widget.
const SignalParentSet cdk.Signal = "parent-set"

// This signal gets emitted whenever a widget should pop up a context menu.
// This usually happens through the standard key binding mechanism; by
// pressing a certain key while a widget is focused, the user can cause the
// widget to pop up a menu. For example, the Entry widget creates a menu
// with clipboard commands. See the section called Implement
// Widget::popup_menu for an example of how to use this signal.
const SignalPopupMenu cdk.Signal = "popup-menu"

// The ::property-notify-event signal will be emitted when a property on the
// widget 's window has been changed or deleted. To receive this signal, the
// Window associated to the widget needs to enable the
// GDK_PROPERTY_CHANGE_MASK mask.
const SignalPropertyNotifyEvent cdk.Signal = "property-notify-event"

// To receive this signal the Window associated to the widget needs to
// enable the GDK_PROXIMITY_IN_MASK mask. This signal will be sent to the
// grab widget if there is one.
const SignalProximityInEvent cdk.Signal = "proximity-in-event"

// To receive this signal the Window associated to the widget needs to
// enable the GDK_PROXIMITY_OUT_MASK mask. This signal will be sent to the
// grab widget if there is one.
const SignalProximityOutEvent cdk.Signal = "proximity-out-event"

// Emitted when has-tooltip is TRUE and the gtk-tooltip-timeout
// has expired with the cursor hovering "above" widget ; or emitted when
// widget got focus in keyboard mode. Using the given coordinates, the signal
// handler should determine whether a tooltip should be shown for widget . If
// this is the case TRUE should be returned, FALSE otherwise. Note that if
// keyboard_mode is TRUE, the values of x and y are undefined and should not
// be used. The signal handler is free to manipulate tooltip with the
// therefore destined function calls.
const SignalQueryTooltip cdk.Signal = "query-tooltip"

const SignalRealize cdk.Signal = "realize"

// The ::screen-changed signal gets emitted when the screen of a widget has
// changed.
// Listener function arguments:
//      previousScreen Screen   the previous screen, or NULL if the widget was not associated with a screen before.
const SignalScreenChanged cdk.Signal = "screen-changed"

// The ::scroll-event signal is emitted when a button in the 4 to 7 range is
// pressed. Wheel mice are usually configured to generate button press events
// for buttons 4 and 5 when the wheel is turned. To receive this signal, the
// Window associated to the widget needs to enable the
// GDK_BUTTON_PRESS_MASK mask. This signal will be sent to the grab widget if
// there is one.
const SignalScrollEvent cdk.Signal = "scroll-event"

// The ::selection-clear-event signal will be emitted when the the widget 's
// window has lost ownership of a selection.
const SignalSelectionClearEvent cdk.Signal = "selection-clear-event"

// Listener function arguments:
//      data SelectionData
//      info int
//      time int
const SignalSelectionGet cdk.Signal = "selection-get"

const SignalSelectionNotifyEvent cdk.Signal = "selection-notify-event"

// Listener function arguments:
//      data SelectionData
//      time int
const SignalSelectionReceived cdk.Signal = "selection-received"

// The ::selection-request-event signal will be emitted when another client
// requests ownership of the selection owned by the widget 's window.
const SignalSelectionRequestEvent cdk.Signal = "selection-request-event"

const SignalShow cdk.Signal = "show"

const SignalShowHelp cdk.Signal = "show-help"

// Listener function arguments:
//      allocation Rectangle
const SignalSizeAllocate cdk.Signal = "size-allocate"

// Listener function arguments:
//      requisition Requisition
const SignalSizeRequest cdk.Signal = "size-request"

// The ::state-changed signal is emitted when the widget state changes. See
// GetState.
// Listener function arguments:
//      state StateType the previous state
const SignalStateChanged cdk.Signal = "state-changed"

// The ::style-set signal is emitted when a new style has been set on a
// widget. Note that style-modifying functions like ModifyBase
// also cause this signal to be emitted.
// Listener function arguments:
//      previousStyle Style     the previous style, or NULL if the widget just got its initial style.
const SignalStyleSet cdk.Signal = "style-set"

const SignalUnmap cdk.Signal = "unmap"

// The ::unmap-event signal will be emitted when the widget 's window is
// unmapped. A window is unmapped when it becomes invisible on the screen. To
// receive this signal, the Window associated to the widget needs to
// enable the GDK_STRUCTURE_MASK mask. GDK will enable this mask
// automatically for all new windows.
const SignalUnmapEvent cdk.Signal = "unmap-event"

const SignalUnrealize cdk.Signal = "unrealize"

// The ::visibility-notify-event will be emitted when the widget 's window is
// obscured or unobscured. To receive this signal the Window associated to
// the widget needs to enable the GDK_VISIBILITY_NOTIFY_MASK mask.
const SignalVisibilityNotifyEvent cdk.Signal = "visibility-notify-event"

// The ::window-state-event will be emitted when the state of the toplevel
// window associated to the widget changes. To receive this signal the
// Window associated to the widget needs to enable the GDK_STRUCTURE_MASK
// mask. GDK will enable this mask automatically for all new windows.
const SignalWindowStateEvent cdk.Signal = "window-state-event"

const WidgetLostFocusHandle = "widget-lost-focus-handler"

const WidgetGainedFocusHandle = "widget-gained-focus-handler"
