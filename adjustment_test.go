package ctk

import (
	"testing"

	"github.com/go-curses/cdk/lib/enums"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAdjustment(t *testing.T) {
	Convey("Testing Adjustments", t, func() {
		Convey("Initialization", func() {
			a := &CAdjustment{}
			So(a.Init(), ShouldEqual, false)
			So(a.Init(), ShouldEqual, true)
		})
		Convey("Basics", func() {
			a := NewAdjustment(5, 0, 10, 1, 2, 4)
			So(a, ShouldNotBeNil)
			So(a.GetValue(), ShouldEqual, 5)
			So(a.GetLower(), ShouldEqual, 0)
			So(a.GetUpper(), ShouldEqual, 10)
			So(a.GetStepIncrement(), ShouldEqual, 1)
			So(a.GetPageIncrement(), ShouldEqual, 2)
			So(a.GetPageSize(), ShouldEqual, 4)
		})
		Convey("Signals...", func() {
			Convey("Value Changed", func() {
				a := NewAdjustment(5, 0, 10, 1, 2, 4)
				signalReceived := false
				a.Connect(
					SignalValueChanged,
					"test-value-changed",
					func(data []interface{}, argv ...interface{}) enums.EventFlag {
						signalReceived = true
						return enums.EVENT_PASS
					},
				)
				a.SetValue(6)
				So(signalReceived, ShouldEqual, true)
				So(a.GetValue(), ShouldEqual, 6)
				signalReceived = false
				a.SetUpper(1)
				So(signalReceived, ShouldEqual, false)
				a.SetLower(1)
				So(signalReceived, ShouldEqual, false)
				a.SetStepIncrement(1)
				So(signalReceived, ShouldEqual, false)
				a.SetPageIncrement(1)
				So(signalReceived, ShouldEqual, false)
				a.SetPageSize(1)
				So(signalReceived, ShouldEqual, false)
				a.Configure(6, 1, 1, 1, 1, 1)
				So(signalReceived, ShouldEqual, false)
				a.Configure(5, 1, 1, 1, 1, 1)
				So(signalReceived, ShouldEqual, true)
				So(a.GetValue(), ShouldEqual, 5)
				signalReceived = false
				a.ClampPage(1, 1)
				So(signalReceived, ShouldEqual, false)
			})
			Convey("Changed", func() {
				a := NewAdjustment(5, 0, 10, 1, 2, 4)
				signalReceived := false
				a.Connect(
					SignalChanged,
					"test-changed",
					func(data []interface{}, argv ...interface{}) enums.EventFlag {
						signalReceived = true
						return enums.EVENT_PASS
					},
				)
				a.SetValue(6)
				So(signalReceived, ShouldEqual, false)
				So(a.GetValue(), ShouldEqual, 6)
				signalReceived = false
				a.SetLower(1)
				So(signalReceived, ShouldEqual, true)
				signalReceived = false
				a.SetUpper(1)
				So(signalReceived, ShouldEqual, true)
				signalReceived = false
				a.SetStepIncrement(1)
				So(signalReceived, ShouldEqual, true)
				signalReceived = false
				a.SetPageIncrement(1)
				So(signalReceived, ShouldEqual, true)
				signalReceived = false
				a.SetPageSize(1)
				So(signalReceived, ShouldEqual, true)
				signalReceived = false
				// a.ClampPage(1, 1)
				// So(signalReceived, ShouldEqual, true)

				// testing configure signals
				signalReceived = false
				a.Configure(6, 1, 1, 1, 1, 1)
				So(signalReceived, ShouldEqual, false)
				a.Configure(5, 1, 1, 1, 1, 1)
				So(signalReceived, ShouldEqual, false)
				a.Configure(5, 0, 10, 1, 2, 4)
				So(signalReceived, ShouldEqual, true)
				So(a.GetValue(), ShouldEqual, 5)
				So(a.GetLower(), ShouldEqual, 0)
				So(a.GetUpper(), ShouldEqual, 10)
				So(a.GetStepIncrement(), ShouldEqual, 1)
				So(a.GetPageIncrement(), ShouldEqual, 2)
				So(a.GetPageSize(), ShouldEqual, 4)
			})
		})
	})
}
