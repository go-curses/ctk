package ctk

import (
	"github.com/go-curses/cdk"
	"github.com/go-curses/cdk/lib/enums"
)

const TypeVButtonBox cdk.CTypeTag = "ctk-v-button-box"

func init() {
	_ = cdk.TypesManager.AddType(TypeVButtonBox, func() interface{} { return MakeVButtonBox() })
}

type VButtonBox interface {
	ButtonBox

	Init() bool
}

type CVButtonBox struct {
	CButtonBox
}

func MakeVButtonBox() *CVButtonBox {
	return NewVButtonBox(false, 0)
}

func NewVButtonBox(homogeneous bool, spacing int) *CVButtonBox {
	b := new(CVButtonBox)
	b.Init()
	b.SetHomogeneous(homogeneous)
	b.SetSpacing(spacing)
	return b
}

func (b *CVButtonBox) Init() bool {
	if b.InitTypeItem(TypeVButtonBox, b) {
		return true
	}
	b.CButtonBox.Init()
	b.flags = NULL_WIDGET_FLAG
	b.SetFlags(PARENT_SENSITIVE | APP_PAINTABLE)
	b.SetOrientation(enums.ORIENTATION_VERTICAL)
	return false
}
