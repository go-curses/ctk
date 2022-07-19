package ctk

import (
	"github.com/go-curses/cdk"
	"github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/cdk/lib/paint"
	"github.com/go-curses/cdk/lib/ptypes"
)

const TypeFakeWindow cdk.CTypeTag = "ctk-fake-window"

type WithFakeWindowFn = func(w Window)

type CFakeWindow struct {
	CWindow
}

func (f *CFakeWindow) Init() (already bool) {
	if f.InitTypeItem(TypeFakeWindow, f) {
		return true
	}
	f.SetAllocation(ptypes.MakeRectangle(80, 24))
	theme, _ := paint.GetDefaultTheme(paint.NilTheme)
	f.SetTheme(theme)
	f.Connect(SignalDraw, "fake-window-draw-handler", f.draw)
	return false
}

func (f *CFakeWindow) draw(data []interface{}, argv ...interface{}) enums.EventFlag {
	return enums.EVENT_STOP
}

func WithFakeWindow(fn WithFakeWindowFn) func() {
	fakeWindow := new(CFakeWindow)
	fakeWindow.Init()
	fakeWindow.SetTheme(paint.GetDefaultMonoTheme())
	return func() {
		fn(fakeWindow)
	}
}

func WithFakeWindowOptions(w, h int, theme paint.Theme, fn WithFakeWindowFn) func() {
	fakeWindow := new(CFakeWindow)
	fakeWindow.Init()
	fakeWindow.SetAllocation(ptypes.MakeRectangle(w, h))
	fakeWindow.SetTheme(theme)
	return func() {
		fn(fakeWindow)
	}
}