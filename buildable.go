package ctk

import (
	"github.com/go-curses/cdk"
)

// TODO: decide if Buildable is even necessary as all of it's methods are present in Widget, Object or MetaData

type Buildable interface {
	Widget

	ListProperties() (known []cdk.Property)
	InitWithProperties(properties map[cdk.Property]string) (already bool, err error)
	Build(builder Builder, element *CBuilderElement) error
	SetProperties(properties map[cdk.Property]string) (err error)
	SetPropertyFromString(property cdk.Property, value string) (err error)
	Connect(signal cdk.Signal, handle string, c cdk.SignalListenerFn, data ...interface{})
}