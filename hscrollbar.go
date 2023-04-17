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
	TypeHScrollbar cdk.CTypeTag = "ctk-h-scrollbar"
)

func init() {
	_ = cdk.TypesManager.AddType(TypeHScrollbar, func() interface{} { return MakeHScrollbar() })
}

type HScrollbar interface {
	Scrollbar
}

var _ HScrollbar = (*CHScrollbar)(nil)

type CHScrollbar struct {
	CScrollbar
}

func MakeHScrollbar() HScrollbar {
	return NewHScrollbar()
}

func NewHScrollbar() HScrollbar {
	s := &CHScrollbar{}
	s.orientation = cenums.ORIENTATION_HORIZONTAL
	s.Init()
	return s
}

func (s *CHScrollbar) Init() (already bool) {
	if s.InitTypeItem(TypeHScrollbar, s) {
		return true
	}
	s.CScrollbar.Init()
	s.SetFlags(enums.SENSITIVE | enums.PARENT_SENSITIVE | enums.CAN_FOCUS | enums.APP_PAINTABLE)
	return false
}