package ctk

import (
	"github.com/go-curses/cdk"
	cenums "github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/ctk/lib/enums"
)

const (
	TypeHBox cdk.CTypeTag = "ctk-h-box"
)

func init() {
	_ = cdk.TypesManager.AddType(TypeHBox, func() interface{} { return MakeHBox() })
}

type HBox interface {
	Box
}

var _ HBox = (*CHBox)(nil)

type CHBox struct {
	CBox
}

func MakeHBox() HBox {
	return NewHBox(false, 0)
}

func NewHBox(homogeneous bool, spacing int) HBox {
	b := new(CHBox)
	b.Init()
	b.SetHomogeneous(homogeneous)
	b.SetSpacing(spacing)
	return b
}

func (b *CHBox) Init() bool {
	if b.InitTypeItem(TypeHBox, b) {
		return true
	}
	b.CBox.Init()
	b.flags = enums.NULL_WIDGET_FLAG
	b.SetFlags(enums.PARENT_SENSITIVE)
	b.SetFlags(enums.APP_PAINTABLE)
	b.SetOrientation(cenums.ORIENTATION_HORIZONTAL)
	return false
}