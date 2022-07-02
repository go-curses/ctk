package ctk

import (
	"github.com/go-curses/cdk"
	cenums "github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/ctk/lib/enums"
)

const (
	TypeVScrollbar cdk.CTypeTag = "ctk-v-scrollbar"
)

func init() {
	_ = cdk.TypesManager.AddType(TypeVScrollbar, func() interface{} { return MakeVScrollbar() })
}

type VScrollbar interface {
	Scrollbar
}

var _ VScrollbar = (*CVScrollbar)(nil)

type CVScrollbar struct {
	CScrollbar
}

func MakeVScrollbar() VScrollbar {
	return NewVScrollbar()
}

func NewVScrollbar() VScrollbar {
	v := &CVScrollbar{}
	v.orientation = cenums.ORIENTATION_VERTICAL
	v.Init()
	return v
}

func (v *CVScrollbar) Init() (already bool) {
	if v.InitTypeItem(TypeVScrollbar, v) {
		return true
	}
	v.CScrollbar.Init()
	v.SetFlags(enums.SENSITIVE | enums.PARENT_SENSITIVE | enums.CAN_FOCUS | enums.APP_PAINTABLE)
	v.SetTheme(DefaultScrollbarTheme)
	return false
}