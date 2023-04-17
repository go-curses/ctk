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

const TypeHButtonBox cdk.CTypeTag = "ctk-h-button-box"

func init() {
	_ = cdk.TypesManager.AddType(TypeHButtonBox, func() interface{} { return MakeHButtonBox() })
}

type HButtonBox interface {
	ButtonBox
}

var _ ButtonBox = (*CButtonBox)(nil)

type CHButtonBox struct {
	CButtonBox
}

func MakeHButtonBox() HButtonBox {
	return NewHButtonBox(false, 0)
}

func NewHButtonBox(homogeneous bool, spacing int) HButtonBox {
	b := new(CHButtonBox)
	b.Init()
	b.SetHomogeneous(homogeneous)
	b.SetSpacing(spacing)
	return b
}

func (b *CHButtonBox) Init() bool {
	if b.InitTypeItem(TypeHButtonBox, b) {
		return true
	}
	b.CButtonBox.Init()
	b.flags = enums.NULL_WIDGET_FLAG
	b.SetFlags(enums.PARENT_SENSITIVE | enums.APP_PAINTABLE)
	b.SetOrientation(cenums.ORIENTATION_HORIZONTAL)
	return false
}