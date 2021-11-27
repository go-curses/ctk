package ctk

import (
	"github.com/go-curses/cdk"
	cenums "github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/ctk/lib/enums"
)

const TypeHButtonBox cdk.CTypeTag = "ctk-h-button-box"

func init() {
	_ = cdk.TypesManager.AddType(TypeHButtonBox, func() interface{} { return MakeHButtonBox() })
}

type HButtonBox interface {
	ButtonBox

	Init() bool
}

type CHButtonBox struct {
	CButtonBox
}

func MakeHButtonBox() HButtonBox {
	return NewHButtonBox(false, 0)
}

func NewHButtonBox(homogeneous bool, spacing int) HButtonBox {
	b := new(CHButtonBox)
	b.Init()
	b.SetHomogeneous(homogeneous)
	b.SetSpacing(spacing)
	return b
}

func (b *CHButtonBox) Init() bool {
	if b.InitTypeItem(TypeHButtonBox, b) {
		return true
	}
	b.CButtonBox.Init()
	b.flags = enums.NULL_WIDGET_FLAG
	b.SetFlags(enums.PARENT_SENSITIVE | enums.APP_PAINTABLE)
	b.SetOrientation(cenums.ORIENTATION_HORIZONTAL)
	return false
}
