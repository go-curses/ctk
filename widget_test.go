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
