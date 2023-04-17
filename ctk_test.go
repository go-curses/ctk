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

func TestCtk(t *testing.T) {
	Convey("Checking CTK Rigging", t, func() {
		Convey("typical", WithApp(
			TestingWithCtkWindow,
			func(app Application) {
				So(app.Version(), ShouldEqual, "v0.0.0")
				d := app.Display()
				w := d.FocusedWindow()
				So(w, ShouldNotBeNil)
			},
		))
	})
}