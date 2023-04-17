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
)

func TestViewport(t *testing.T) {
	Convey("Testing Viewports", t, func() {
		Convey("Initialization", func() {
			v := &CViewport{}
			So(v.Init(), ShouldEqual, false)
			So(v.Init(), ShouldEqual, true)
		})

		Convey("Basics", func() {
			ha := NewAdjustment(0, 0, 10, 1, 5, 5)
			va := NewAdjustment(0, 0, 10, 1, 5, 5)
			v := NewViewport(ha, va)
			So(v, ShouldNotBeNil)
			So(v.GetHAdjustment(), ShouldEqual, ha)
			So(v.GetVAdjustment(), ShouldEqual, va)
			nha := NewAdjustment(1, 0, 10, 1, 5, 5)
			nva := NewAdjustment(1, 0, 10, 1, 5, 5)
			v.SetHAdjustment(nha)
			So(v.GetHAdjustment(), ShouldEqual, nha)
			v.SetVAdjustment(nva)
			So(v.GetVAdjustment(), ShouldEqual, nva)
		})
	})
}