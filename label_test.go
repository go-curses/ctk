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

	"github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/cdk/lib/ptypes"
	"github.com/go-curses/cdk/memphis"
	. "github.com/smartystreets/goconvey/convey"
)

func TestLabel(t *testing.T) {
	Convey("Testing Labels", t, func() {
		Convey("basics: justification", func() {
			l := NewLabel("test")
			So(l, ShouldNotBeNil)
			So(l.GetTheme().String(), ShouldEqual, "{Content={Normal={white[#ffffff],navy[#000080],0},Selected={white[#ffffff],navy[#000080],0},Active={white[#ffffff],navy[#000080],0},Prelight={white[#ffffff],navy[#000080],0},Insensitive={white[#ffffff],navy[#000080],0},FillRune=32,BorderRunes={BorderRunes=9488,9472,9484,9474,9492,9472,9496,9474},ArrowRunes={ArrowRunes=8593,8592,8595,8594},Overlay=false},Border={Normal={white[#ffffff],navy[#000080],0},Selected={white[#ffffff],navy[#000080],0},Active={white[#ffffff],navy[#000080],0},Prelight={white[#ffffff],navy[#000080],0},Insensitive={white[#ffffff],navy[#000080],0},FillRune=32,BorderRunes={BorderRunes=9488,9472,9484,9474,9492,9472,9496,9474},ArrowRunes={ArrowRunes=8593,8592,8595,8594},Overlay=false}}")
			l.Show()
			l.SetOrigin(0, 0)
			l.SetSizeRequest(10, 1)
			l.SetAllocation(ptypes.MakeRectangle(10, 1))
			l.Resize()
			Convey("left justification (default)", func() {
				if surface, err := memphis.GetSurface(l.ObjectID()); err != nil {
					So(err, ShouldBeNil)
				} else {
					size := surface.GetSize()
					So(size.W, ShouldEqual, 10)
					So(size.H, ShouldEqual, 1)
					So(l.Draw(), ShouldEqual, enums.EVENT_STOP)
					for i := 0; i < 10; i++ {
						cell := surface.GetContent(i, 0)
						So(cell, ShouldNotBeNil)
						switch i {
						case 0:
							So(cell.Value(), ShouldEqual, 't')
						case 1:
							So(cell.Value(), ShouldEqual, 'e')
						case 2:
							So(cell.Value(), ShouldEqual, 's')
						case 3:
							So(cell.Value(), ShouldEqual, 't')
						default:
							So(cell.IsSpace(), ShouldEqual, true)
						}
					}
				}
			})
			Convey("right justification", func() {
				l.SetJustify(enums.JUSTIFY_RIGHT)
				if surface, err := memphis.GetSurface(l.ObjectID()); err != nil {
					So(err, ShouldBeNil)
				} else {
					size := surface.GetSize()
					So(size.W, ShouldEqual, 10)
					So(size.H, ShouldEqual, 1)
					So(l.Draw(), ShouldEqual, enums.EVENT_STOP)
					for i := 0; i < 10; i++ {
						cell := surface.GetContent(i, 0)
						So(cell, ShouldNotBeNil)
						switch i {
						case 6:
							So(cell.Value(), ShouldEqual, 't')
						case 7:
							So(cell.Value(), ShouldEqual, 'e')
						case 8:
							So(cell.Value(), ShouldEqual, 's')
						case 9:
							So(cell.Value(), ShouldEqual, 't')
						default:
							So(cell.IsSpace(), ShouldEqual, true)
						}
					}
				}
			})
			Convey("center justification", func() {
				l.SetJustify(enums.JUSTIFY_CENTER)
				if surface, err := memphis.GetSurface(l.ObjectID()); err != nil {
					So(err, ShouldBeNil)
				} else {
					size := surface.GetSize()
					So(size.W, ShouldEqual, 10)
					So(size.H, ShouldEqual, 1)
					So(l.Draw(), ShouldEqual, enums.EVENT_STOP)
					for i := 0; i < 10; i++ {
						cell := surface.GetContent(i, 0)
						So(cell, ShouldNotBeNil)
						switch i {
						case 3:
							So(cell.Value(), ShouldEqual, 't')
						case 4:
							So(cell.Value(), ShouldEqual, 'e')
						case 5:
							So(cell.Value(), ShouldEqual, 's')
						case 6:
							So(cell.Value(), ShouldEqual, 't')
						default:
							So(cell.IsSpace(), ShouldEqual, true)
						}
					}
				}
			})
		})
	})
}