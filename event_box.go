package ctk

import (
	"github.com/go-curses/cdk"
	cenums "github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/ctk/lib/enums"
)

const TypeEventBox cdk.CTypeTag = "ctk-event-box"

func init() {
	_ = cdk.TypesManager.AddType(TypeEventBox, func() interface{} { return MakeEventBox() })
}

// EventBox Hierarchy:
//	Object
//	  +- Widget
//	    +- Container
//	      +- Bin
//	        +- EventBox
//
// The EventBox Widget is used to capture Widget events (mouse, keyboard)
// without needing having any defined user-interface.
type EventBox interface {
	Bin
	Buildable

	Init() (already bool)
	SetAboveChild(aboveChild bool)
	GetAboveChild() (value bool)
	SetVisibleWindow(visibleWindow bool)
	GetVisibleWindow() (value bool)
	Activate() (value bool)
	CancelEvent()
	ProcessEvent(evt cdk.Event) cenums.EventFlag
}

// The CEventBox structure implements the EventBox interface and is exported
// to facilitate type embedding with custom implementations. No member variables
// are exported as the interface methods are the only intended means of
// interacting with EventBox objects.
type CEventBox struct {
	CBin
}

// MakeEventBox is used by the Buildable system to construct a new EventBox.
func MakeEventBox() EventBox {
	return NewEventBox()
}

// NewEventBox is the constructor for new EventBox instances.
func NewEventBox() (value EventBox) {
	e := new(CEventBox)
	e.Init()
	return e
}

// Init initializes a EventBox object. This must be called at least once to
// set up the necessary defaults and allocate any memory structures. Calling
// this more than once is safe though unnecessary. Only the first call will
// result in any effect upon the EventBox instance. Init is used in the
// NewEventBox constructor and only necessary when implementing a derivative
// EventBox type.
func (b *CEventBox) Init() (already bool) {
	if b.InitTypeItem(TypeEventBox, b) {
		return true
	}
	b.CBin.Init()
	b.SetFlags(enums.SENSITIVE | enums.PARENT_SENSITIVE | enums.CAN_DEFAULT | enums.RECEIVES_DEFAULT | enums.CAN_FOCUS | enums.APP_PAINTABLE)
	_ = b.InstallProperty(PropertyAboveChild, cdk.BoolProperty, true, false)
	_ = b.InstallProperty(PropertyVisibleWindow, cdk.BoolProperty, true, true)
	return false
}

// SetAboveChild updates whether the event box window is positioned above the
// windows of its child, as opposed to below it. If the window is above, all
// events inside the event box will go to the event box. If the window is below,
// events in windows of child widgets will first got to that widget, and then to
// its parents. The default is to keep the window below the child.
//
// Parameters:
// 	aboveChild	TRUE if the event box window is above the windows of its child
func (b *CEventBox) SetAboveChild(aboveChild bool) {
	if err := b.SetBoolProperty(PropertyAboveChild, aboveChild); err != nil {
		b.LogErr(err)
	}
}

// GetAboveChild returns whether the event box window is above or below the
// windows of its child.
// See: SetAboveChild()
func (b *CEventBox) GetAboveChild() (value bool) {
	var err error
	if value, err = b.GetBoolProperty(PropertyAboveChild); err != nil {
		b.LogErr(err)
	}
	return
}

// SetVisibleWindow updates whether the event box uses a visible or invisible
// child window. The default is to use visible windows. In an invisible window
// event box, the window that the event box creates is a GDK_INPUT_ONLY window,
// which means that it is invisible and only serves to receive events. A visible
// window event box creates a visible (GDK_INPUT_OUTPUT) window that acts as the
// parent window for all the widgets contained in the event box. You should
// generally make your event box invisible if you just want to trap events.
// Creating a visible window may cause artifacts that are visible to the
// user. The main reason to create a non input-only event box is if you want to
// set the background to a different color or draw on it.
//
// Parameters:
// 	visibleWindow	boolean value
func (b *CEventBox) SetVisibleWindow(visibleWindow bool) {
	if err := b.SetBoolProperty(PropertyVisibleWindow, visibleWindow); err != nil {
		b.LogErr(err)
	}
}

// GetVisibleWindow returns whether the event box has a visible window.
// See: SetVisibleWindow()
func (b *CEventBox) GetVisibleWindow() (value bool) {
	var err error
	if value, err = b.GetBoolProperty(PropertyVisibleWindow); err != nil {
		b.LogErr(err)
	}
	return
}

// Activate will emit an activate signal and return TRUE if the signal handlers
// return EVENT_STOP indicating that the event was in fact handled.
func (b *CEventBox) Activate() (value bool) {
	return b.Emit(SignalActivate, b) == cenums.EVENT_STOP
}

// CancelEvent will emit a cancel-event signal.
func (b *CEventBox) CancelEvent() {
	b.Emit(SignalCancelEvent, b)
}

// ProcessEvent manages the processing of events, current this is just emitting
// a cdk-event signal and returning the result.
func (b *CEventBox) ProcessEvent(evt cdk.Event) cenums.EventFlag {
	return b.Emit(SignalCdkEvent, b, evt)
}

// Whether the event-trapping window of the eventbox is above the window of
// the child widget as opposed to below it.
// Flags: Read / Write
// Default value: FALSE
const PropertyAboveChild cdk.Property = "above-child"

// Whether the event box is visible, as opposed to invisible and only used to
// trap events.
// Flags: Read / Write
// Default value: TRUE
const PropertyVisibleWindow cdk.Property = "visible-window"
