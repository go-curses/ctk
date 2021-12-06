package ctk

import (
	"github.com/go-curses/cdk"
	cenums "github.com/go-curses/cdk/lib/enums"
)

const (
	SignalAllocation        cdk.Signal = "allocation"
	SignalCancelEvent       cdk.Signal = "cancel-event"
	SignalCdkEvent          cdk.Signal = "cdk-event"
	SignalDraw              cdk.Signal = cdk.SignalDraw
	SignalError             cdk.Signal = "error"
	SignalEventKey          cdk.Signal = "key-event"
	SignalEventMouse        cdk.Signal = "mouse-event"
	SignalGainedEventFocus  cdk.Signal = "gained-event-focus"
	SignalGainedFocus       cdk.Signal = "gained-focus"
	SignalGrabEventFocus    cdk.Signal = "grab-event-focus"
	SignalHomogeneous       cdk.Signal = "homogeneous"
	SignalInvalidate        cdk.Signal = "invalidate"
	SignalLostEventFocus    cdk.Signal = "lost-event-focus"
	SignalLostFocus         cdk.Signal = "lost-focus"
	SignalName              cdk.Signal = "name"
	SignalOrigin            cdk.Signal = "origin"
	SignalPackEnd           cdk.Signal = "pack-end"
	SignalPackStart         cdk.Signal = "pack-start"
	SignalReleaseEventFocus cdk.Signal = "set-event-focus"
	SignalReorderChild      cdk.Signal = "reorder-child"
	SignalReparent          cdk.Signal = "reparent"
	SignalResize            cdk.Signal = "resize"
	SignalSetEventFocus     cdk.Signal = "set-event-focus"
	SignalSetFlags          cdk.Signal = "set-flags"
	SignalSetParent         cdk.Signal = "set-parent"
	SignalSetProperty       cdk.Signal = cdk.SignalSetProperty
	SignalSetSensitive      cdk.Signal = "set-sensitive"
	SignalSetSizeRequest    cdk.Signal = "set-size-request"
	SignalSetState          cdk.Signal = "set-state"
	SignalSetTheme          cdk.Signal = "set-theme"
	SignalSetWindow         cdk.Signal = "set-window"
	SignalShowAll           cdk.Signal = "show-all"
	SignalSpacing           cdk.Signal = "spacing"
	SignalTextDirection     cdk.Signal = "text-direction"
	SignalUnparent          cdk.Signal = "unparent"
	SignalUnsetFlags        cdk.Signal = "unset-flags"
	SignalUnsetState        cdk.Signal = "unset-state"
)

func WithArgvNoneSignal(fn func(), eventFlag cenums.EventFlag) cdk.SignalListenerFn {
	return func(_ []interface{}, _ ...interface{}) cenums.EventFlag {
		fn()
		return eventFlag
	}
}

func WithArgvNoneWithFlagsSignal(fn func() cenums.EventFlag) cdk.SignalListenerFn {
	return func(_ []interface{}, _ ...interface{}) cenums.EventFlag {
		return fn()
	}
}

func ArgvSignalEvent(argv ...interface{}) (object Object, event cdk.Event, ok bool) {
	if len(argv) == 2 {
		if object, ok = argv[0].(Object); ok {
			if event, ok = argv[1].(cdk.Event); ok {
				return
			}
			event = nil
		}
		object = nil
	}
	return
}

func WithArgvSignalEvent(fn SignalEventFn) cdk.SignalListenerFn {
	return func(_ []interface{}, argv ...interface{}) cenums.EventFlag {
		if widget, event, ok := ArgvSignalEvent(argv...); ok {
			return fn(widget, event)
		}
		return cenums.EVENT_STOP
	}
}

type SignalEventFn = func(object Object, event cdk.Event) cenums.EventFlag
