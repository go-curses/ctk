package ctk

import (
	"github.com/go-curses/cdk"
	"github.com/go-curses/cdk/lib/enums"
)

const (
	TypeVScrollbar cdk.CTypeTag = "ctk-v-scrollbar"
)

func init() {
	_ = cdk.TypesManager.AddType(TypeVScrollbar, func() interface{} { return MakeVScrollbar() })
}

type VScrollbar interface {
	Scrollbar

	Init() (already bool)
}

type CVScrollbar struct {
	CScrollbar
}

func MakeVScrollbar() *CVScrollbar {
	return NewVScrollbar()
}

func NewVScrollbar() *CVScrollbar {
	v := &CVScrollbar{}
	v.orientation = enums.ORIENTATION_VERTICAL
	v.Init()
	return v
}

func (v *CVScrollbar) Init() (already bool) {
	if v.InitTypeItem(TypeVScrollbar, v) {
		return true
	}
	v.CScrollbar.Init()
	v.SetFlags(SENSITIVE | PARENT_SENSITIVE | CAN_FOCUS | APP_PAINTABLE)
	v.SetTheme(DefaultScrollbarTheme)
	return false
}
