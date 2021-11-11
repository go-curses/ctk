package ctk

import (
	"testing"

	"github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/cdk/lib/paint"
	"github.com/go-curses/cdk/lib/ptypes"
	"github.com/go-curses/cdk/memphis"
	. "github.com/smartystreets/goconvey/convey"
)

func TestLabel(t *testing.T) {
	Convey("Testing Labels", t, func() {
		Convey("basics: justification", func() {
			l := NewLabel("test")
			So(l, ShouldNotBeNil)
			So(l.GetTheme().String(), ShouldEqual, paint.DefaultColorTheme.String())
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
