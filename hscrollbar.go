package ctk

import (
	"github.com/go-curses/cdk"
	"github.com/go-curses/cdk/lib/enums"
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
	s.orientation = enums.ORIENTATION_HORIZONTAL
	s.Init()
	return s
}

func (s *CHScrollbar) Init() (already bool) {
	if s.InitTypeItem(TypeHScrollbar, s) {
		return true
	}
	s.CScrollbar.Init()
	s.SetFlags(SENSITIVE | PARENT_SENSITIVE | CAN_FOCUS | APP_PAINTABLE)
	s.SetTheme(DefaultColorScrollbarTheme)
	return false
}
