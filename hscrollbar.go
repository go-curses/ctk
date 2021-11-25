package ctk

import (
	"github.com/go-curses/cdk"
	cenums "github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/ctk/lib/enums"
)

const (
	TypeHScrollbar cdk.CTypeTag = "ctk-h-scrollbar"
)

func init() {
	_ = cdk.TypesManager.AddType(TypeHScrollbar, func() interface{} { return MakeHScrollbar() })
}

type HScrollbar interface {
	Scrollbar

	Init() (already bool)
}

type CHScrollbar struct {
	CScrollbar
}

func MakeHScrollbar() *CHScrollbar {
	return NewHScrollbar()
}

func NewHScrollbar() *CHScrollbar {
	s := &CHScrollbar{}
	s.orientation = cenums.ORIENTATION_HORIZONTAL
	s.Init()
	return s
}

func (s *CHScrollbar) Init() (already bool) {
	if s.InitTypeItem(TypeHScrollbar, s) {
		return true
	}
	s.CScrollbar.Init()
	s.SetFlags(enums.SENSITIVE | enums.PARENT_SENSITIVE | enums.CAN_FOCUS | enums.APP_PAINTABLE)
	s.SetTheme(DefaultScrollbarTheme)
	return false
}
