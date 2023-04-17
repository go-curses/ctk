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
	"github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/cdk/lib/paint"
	"github.com/go-curses/cdk/lib/ptypes"
)

const TypeFakeWindow cdk.CTypeTag = "ctk-fake-window"

type WithFakeWindowFn = func(w Window)

type CFakeWindow struct {
	CWindow
}

func (f *CFakeWindow) Init() (already bool) {
	if f.InitTypeItem(TypeFakeWindow, f) {
		return true
	}
	f.SetAllocation(ptypes.MakeRectangle(80, 24))
	theme, _ := paint.GetTheme(paint.NilTheme)
	f.SetTheme(theme)
	f.Connect(SignalDraw, "fake-window-draw-handler", f.draw)
	return false
}

func (f *CFakeWindow) draw(data []interface{}, argv ...interface{}) enums.EventFlag {
	return enums.EVENT_STOP
}

func WithFakeWindow(fn WithFakeWindowFn) func() {
	fakeWindow := new(CFakeWindow)
	fakeWindow.Init()
	fakeWindow.SetTheme(paint.GetDefaultMonoTheme())
	return func() {
		fn(fakeWindow)
	}
}

func WithFakeWindowOptions(w, h int, theme paint.Theme, fn WithFakeWindowFn) func() {
	fakeWindow := new(CFakeWindow)
	fakeWindow.Init()
	fakeWindow.SetAllocation(ptypes.MakeRectangle(w, h))
	fakeWindow.SetTheme(theme)
	return func() {
		fn(fakeWindow)
	}
}