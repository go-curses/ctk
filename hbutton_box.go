package ctk

import (
	"github.com/go-curses/cdk"
	"github.com/go-curses/cdk/lib/enums"
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

func MakeHButtonBox() *CHButtonBox {
	return NewHButtonBox(false, 0)
}

func NewHButtonBox(homogeneous bool, spacing int) *CHButtonBox {
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
	b.flags = NULL_WIDGET_FLAG
	b.SetFlags(PARENT_SENSITIVE | APP_PAINTABLE)
	b.SetOrientation(enums.ORIENTATION_HORIZONTAL)
	return false
}
