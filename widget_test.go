package ctk

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestWidget(t *testing.T) {
	cw := &CWidget{}
	Convey("widget basics", t, func() {
		So(cw.Init(), ShouldEqual, false)
		So(cw.Init(), ShouldEqual, true)
	})

	Convey("widget states", t, func() {
		So(cw.GetState(), ShouldEqual, StateNormal)
		cw.SetState(StateActive)
		So(cw.GetState(), ShouldEqual, StateNormal|StateActive)
		cw.UnsetState(StateActive)
		cw.SetState(StateActive)
		So(cw.GetState(), ShouldEqual, StateNormal|StateActive)
		cw.SetState(StateNone)
		So(cw.GetState(), ShouldEqual, StateNormal)
		cw.SetState(StateActive)
		So(cw.GetState(), ShouldEqual, StateNormal|StateActive)
		cw.UnsetState(StateNone)
		So(cw.GetState(), ShouldEqual, StateNormal|StateActive)
	})

	Convey("widget flags", t, func() {
		So(cw.GetFlags(), ShouldEqual, NULL_WIDGET_FLAG)
		cw.SetFlags(TOPLEVEL)
		So(cw.GetFlags(), ShouldEqual, TOPLEVEL)
		cw.SetFlags(VISIBLE)
		So(cw.GetFlags(), ShouldEqual, TOPLEVEL|VISIBLE)
		cw.UnsetFlags(VISIBLE)
		So(cw.GetFlags(), ShouldEqual, TOPLEVEL)
	})

	Convey("widget safe methods", t, func() {
		w := NewWindowWithTitle("test")
		cw.SetWindow(w)
		// default
		So(cw.CanDefault(), ShouldEqual, false)
		So(cw.IsDefault(), ShouldEqual, false)
		cw.SetFlags(CAN_DEFAULT)
		// w.rebuildFocusChain()
		So(cw.CanDefault(), ShouldEqual, true)
		So(cw.IsDefault(), ShouldEqual, false)
		// focus
		So(cw.CanFocus(), ShouldEqual, false)
		So(cw.IsFocus(), ShouldEqual, false)
		cw.SetFlags(CAN_FOCUS)
		// w.rebuildFocusChain()
		So(cw.CanFocus(), ShouldEqual, true)
		So(cw.IsFocus(), ShouldEqual, false)
		// cw.GrabFocus()
		// So(cw.IsFocus(), ShouldEqual, true)

	})
}
