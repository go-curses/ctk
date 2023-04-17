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

func TestRange(t *testing.T) {
	Convey("Testing Ranges", t, func() {
		Convey("basic checks", func() {
			r := &CRange{}
			So(r.Init(), ShouldEqual, false)
			So(r.Init(), ShouldEqual, true)
			adj := r.GetAdjustment()
			So(adj, ShouldNotBeNil)
			So(adj.GetLower(), ShouldEqual, 0)
			So(adj.GetUpper(), ShouldEqual, 0)
			So(adj.GetStepIncrement(), ShouldEqual, 0)
			So(adj.GetPageIncrement(), ShouldEqual, 0)
			So(adj.GetValue(), ShouldEqual, 0)
			// stuff
			r.SetRange(1, 9)
			So(adj.GetLower(), ShouldEqual, 1)
			So(adj.GetUpper(), ShouldEqual, 9)
			min, max := r.GetRange()
			So(min, ShouldEqual, 1)
			So(max, ShouldEqual, 9)
			r.SetIncrements(2, 5)
			So(adj.GetStepIncrement(), ShouldEqual, 2)
			So(adj.GetPageIncrement(), ShouldEqual, 5)
			r.SetValue(9)
			So(adj.GetValue(), ShouldEqual, 9)
			r.SetValue(10)
			So(adj.GetValue(), ShouldEqual, 9)
			So(r.GetValue(), ShouldEqual, 9)
			step, page := r.GetIncrements()
			So(step, ShouldEqual, 2)
			So(page, ShouldEqual, 5)
			So(r.GetRestrictToFillLevel(), ShouldEqual, false)
			r.SetRestrictToFillLevel(true)
			So(r.GetRestrictToFillLevel(), ShouldEqual, true)
			So(r.GetFillLevel(), ShouldEqual, 1.0)
			r.SetFillLevel(0.5)
			r.SetValue(4)
			So(adj.GetValue(), ShouldEqual, 4)
			r.SetValue(5)
			So(adj.GetValue(), ShouldEqual, 4)
			r.SetRestrictToFillLevel(false)
			So(r.GetRestrictToFillLevel(), ShouldEqual, false)
			So(r.GetShowFillLevel(), ShouldEqual, false)
			r.SetShowFillLevel(true)
			So(r.GetShowFillLevel(), ShouldEqual, true)
			So(r.GetInverted(), ShouldEqual, false)
			r.SetInverted(true)
			So(r.GetInverted(), ShouldEqual, true)
			So(r.GetLowerStepperSensitivity(), ShouldEqual, enums.SensitivityAuto)
			So(r.GetUpperStepperSensitivity(), ShouldEqual, enums.SensitivityAuto)
			r.SetLowerStepperSensitivity(enums.SensitivityOff)
			So(r.GetLowerStepperSensitivity(), ShouldEqual, enums.SensitivityOff)
			r.SetUpperStepperSensitivity(enums.SensitivityOn)
			So(r.GetUpperStepperSensitivity(), ShouldEqual, enums.SensitivityOn)
			So(r.GetFlippable(), ShouldEqual, false)
			r.SetFlippable(true)
			So(r.GetFlippable(), ShouldEqual, true)
			So(r.GetMinSliderLength(), ShouldEqual, 1)
			r.SetMinSliderLength(2)
			So(r.GetMinSliderLength(), ShouldEqual, 2)
			r.SetMinSliderLength(-1)
			So(r.GetMinSliderLength(), ShouldEqual, 1)
			So(r.GetSliderSizeFixed(), ShouldEqual, false)
			r.SetSliderSizeFixed(true)
			So(r.GetSliderSizeFixed(), ShouldEqual, true)
			r.SetSliderSizeFixed(false)
			So(r.GetSliderLength(), ShouldEqual, -1)
			r.SetSliderLength(1)
			So(r.GetSliderSizeFixed(), ShouldEqual, true)
			So(r.GetSliderLength(), ShouldEqual, 1)
			r.SetSliderLength(-1)
			So(r.GetSliderSizeFixed(), ShouldEqual, false)
			So(r.GetSliderLength(), ShouldEqual, -1)
			So(r.GetStepperSize(), ShouldEqual, -1)
			r.SetStepperSize(2)
			So(r.GetStepperSpacing(), ShouldEqual, 0)
			r.SetStepperSpacing(1)
			So(r.GetStepperSpacing(), ShouldEqual, 1)
			So(r.GetTroughUnderSteppers(), ShouldEqual, false)
			r.SetTroughUnderSteppers(true)
			So(r.GetTroughUnderSteppers(), ShouldEqual, true)
			// nil testing
			// r.SetAdjustment(nil)
			// So(r.GetAdjustment(), ShouldBeNil)
			// r.SetMinSliderLength(200)
			// So(r.GetMinSliderLength(), ShouldEqual, 200)
			// r.SetSliderLength(200)
			// So(r.GetSliderLength(), ShouldEqual, 200)
		})
	})
}