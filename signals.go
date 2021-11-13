package ctk

import (
	"github.com/go-curses/cdk"
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
	SignalSetThemeRequest   cdk.Signal = "set-theme-request"
	SignalSetWindow         cdk.Signal = "set-window"
	SignalShowAll           cdk.Signal = "show-all"
	SignalSpacing           cdk.Signal = "spacing"
	SignalTextDirection     cdk.Signal = "text-direction"
	SignalUnparent          cdk.Signal = "unparent"
	SignalUnsetFlags        cdk.Signal = "unset-flags"
	SignalUnsetState        cdk.Signal = "unset-state"
)
