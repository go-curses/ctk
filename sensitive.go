package ctk

import (
	"github.com/go-curses/cdk"
	"github.com/go-curses/cdk/lib/enums"
)

type Sensitive interface {
	Object

	GetWindow() Window
	CanFocus() bool
	IsFocus() bool
	IsFocused() bool
	IsVisible() bool
	GrabFocus()
	CancelEvent()
	IsSensitive() bool
	SetSensitive(sensitive bool)
	ProcessEvent(evt cdk.Event) enums.EventFlag
}
