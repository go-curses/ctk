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
	cenums "github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/ctk/lib/enums"
)

const (
	TypeVBox cdk.CTypeTag = "ctk-v-box"
)

func init() {
	_ = cdk.TypesManager.AddType(TypeVBox, func() interface{} { return MakeVBox() })
}

// Basic vbox interface
type VBox interface {
	Box
}

var _ VBox = (*CVBox)(nil)

type CVBox struct {
	CBox
}

func MakeVBox() VBox {
	return NewVBox(false, 0)
}

func NewVBox(homogeneous bool, spacing int) VBox {
	b := new(CVBox)
	b.Init()
	b.SetHomogeneous(homogeneous)
	b.SetSpacing(spacing)
	return b
}

func (b *CVBox) Init() bool {
	if b.InitTypeItem(TypeVBox, b) {
		return true
	}
	b.CBox.Init()
	b.flags = enums.NULL_WIDGET_FLAG
	b.SetFlags(enums.PARENT_SENSITIVE | enums.APP_PAINTABLE)
	b.SetOrientation(cenums.ORIENTATION_VERTICAL)
	return false
}