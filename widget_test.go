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
	"testing"

	"github.com/go-curses/ctk/lib/enums"
	. "github.com/smartystreets/goconvey/convey"
)

func TestWidget(t *testing.T) {
	cw := &CWidget{}
	Convey("widget basics", t, func() {
		So(cw.Init(), ShouldEqual, false)
		So(cw.Init(), ShouldEqual, true)
	})

	Convey("widget states", t, func() {
		So(cw.GetState(), ShouldEqual, enums.StateNormal)
		cw.SetState(enums.StateActive)
		So(cw.GetState(), ShouldEqual, enums.StateNormal|enums.StateActive)
		cw.UnsetState(enums.StateActive)
		cw.SetState(enums.StateActive)
		So(cw.GetState(), ShouldEqual, enums.StateNormal|enums.StateActive)
		cw.SetState(enums.StateNone)
		So(cw.GetState(), ShouldEqual, enums.StateNormal)
		cw.SetState(enums.StateActive)
		So(cw.GetState(), ShouldEqual, enums.StateNormal|enums.StateActive)
		cw.UnsetState(enums.StateNone)
		So(cw.GetState(), ShouldEqual, enums.StateNormal|enums.StateActive)
	})

	Convey("widget flags", t, func() {
		So(cw.GetFlags(), ShouldEqual, enums.NULL_WIDGET_FLAG)
		cw.SetFlags(enums.TOPLEVEL)
		So(cw.GetFlags(), ShouldEqual, enums.TOPLEVEL)
		cw.SetFlags(enums.VISIBLE)
		So(cw.GetFlags(), ShouldEqual, enums.TOPLEVEL|enums.VISIBLE)
		cw.UnsetFlags(enums.VISIBLE)
		So(cw.GetFlags(), ShouldEqual, enums.TOPLEVEL)
	})

	Convey("widget safe methods", t, func() {
		w := NewWindowWithTitle("test")
		cw.SetWindow(w)
		// default
		So(cw.CanDefault(), ShouldEqual, false)
		So(cw.IsDefault(), ShouldEqual, false)
		cw.SetFlags(enums.CAN_DEFAULT)
		// w.rebuildFocusChain()
		So(cw.CanDefault(), ShouldEqual, true)
		So(cw.IsDefault(), ShouldEqual, false)
		// focus
		So(cw.CanFocus(), ShouldEqual, false)
		So(cw.IsFocus(), ShouldEqual, false)
		cw.SetFlags(enums.CAN_FOCUS)
		// w.rebuildFocusChain()
		So(cw.CanFocus(), ShouldEqual, true)
		So(cw.IsFocus(), ShouldEqual, false)
		// cw.GrabFocus()
		// So(cw.IsFocus(), ShouldEqual, true)

	})
}