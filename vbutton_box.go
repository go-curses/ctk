package ctk

import (
	"github.com/go-curses/cdk"
	cenums "github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/ctk/lib/enums"
)

const TypeVButtonBox cdk.CTypeTag = "ctk-v-button-box"

func init() {
	_ = cdk.TypesManager.AddType(TypeVButtonBox, func() interface{} { return MakeVButtonBox() })
}

type VButtonBox interface {
	ButtonBox
}

var _ VButtonBox = (*CVButtonBox)(nil)

type CVButtonBox struct {
	CButtonBox
}

func MakeVButtonBox() VButtonBox {
	return NewVButtonBox(false, 0)
}

func NewVButtonBox(homogeneous bool, spacing int) VButtonBox {
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
	b.flags = enums.NULL_WIDGET_FLAG
	b.SetFlags(enums.PARENT_SENSITIVE | enums.APP_PAINTABLE)
	b.SetOrientation(cenums.ORIENTATION_VERTICAL)
	return false
}