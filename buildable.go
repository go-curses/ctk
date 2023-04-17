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