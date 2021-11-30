package ctk

import (
	"github.com/go-curses/cdk"
	"github.com/go-curses/ctk/lib/enums"
)

type Buildable interface {
	Object

	ListProperties() (known []cdk.Property)
	InitWithProperties(properties map[cdk.Property]string) (already bool, err error)
	Build(builder Builder, element *CBuilderElement) error
	SetProperties(properties map[cdk.Property]string) (err error)
	SetPropertyFromString(property cdk.Property, value string) (err error)
	SetSensitive(sensitive bool)
	SetFlags(widgetFlags enums.WidgetFlags)
	UnsetFlags(widgetFlags enums.WidgetFlags)
	Connect(signal cdk.Signal, handle string, c cdk.SignalListenerFn, data ...interface{})
	LogErr(err error)
	Show()
	GrabFocus()
}
