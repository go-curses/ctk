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
