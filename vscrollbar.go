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
	TypeVScrollbar cdk.CTypeTag = "ctk-v-scrollbar"
)

func init() {
	_ = cdk.TypesManager.AddType(TypeVScrollbar, func() interface{} { return MakeVScrollbar() })
}

type VScrollbar interface {
	Scrollbar
}

var _ VScrollbar = (*CVScrollbar)(nil)

type CVScrollbar struct {
	CScrollbar
}

func MakeVScrollbar() VScrollbar {
	return NewVScrollbar()
}

func NewVScrollbar() VScrollbar {
	v := &CVScrollbar{}
	v.orientation = cenums.ORIENTATION_VERTICAL
	v.Init()
	return v
}

func (v *CVScrollbar) Init() (already bool) {
	if v.InitTypeItem(TypeVScrollbar, v) {
		return true
	}
	v.CScrollbar.Init()
	v.SetFlags(enums.SENSITIVE | enums.PARENT_SENSITIVE | enums.CAN_FOCUS | enums.APP_PAINTABLE)
	return false
}