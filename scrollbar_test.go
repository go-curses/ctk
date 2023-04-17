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

	. "github.com/smartystreets/goconvey/convey"

	"github.com/go-curses/cdk/lib/enums"
)

func TestVScrollbar(t *testing.T) {
	Convey("Testing raw scrollbar", t, func() {
		s := &CScrollbar{}
		Convey("default to vertical", func() {
			So(s.orientation, ShouldEqual, 0)
			So(s.Init(), ShouldEqual, false)
			So(s.Init(), ShouldEqual, true)
			So(s.orientation, ShouldEqual, enums.ORIENTATION_VERTICAL)
		})
	})
	Convey("Testing vertical scrollbars", t, func() {
		Convey("basic checks", func() {
			vs := NewVScrollbar()
			So(vs, ShouldNotBeNil)
			So(vs.GetHasBackwardStepper(), ShouldEqual, true)
			vs.SetHasBackwardStepper(false)
			So(vs.GetHasBackwardStepper(), ShouldEqual, false)
		})
	})
}