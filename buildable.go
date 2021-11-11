package ctk

import (
	"github.com/go-curses/cdk"
)

type Buildable interface {
	ListProperties() (known []cdk.Property)
	InitWithProperties(properties map[cdk.Property]string) (already bool, err error)
	Build(builder Builder, element *CBuilderElement) error
	SetProperties(properties map[cdk.Property]string) (err error)
	SetPropertyFromString(property cdk.Property, value string) (err error)
	SetSensitive(sensitive bool)
	SetFlags(widgetFlags WidgetFlags)
	UnsetFlags(widgetFlags WidgetFlags)
	Connect(signal cdk.Signal, handle string, c cdk.SignalListenerFn, data ...interface{})
	LogErr(err error)
	Show()
	GrabFocus()
	// GetSizeRequest() (size cdk.Rectangle)
	// SetSizeRequest(w, h int)
}
