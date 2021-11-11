package ctk

import (
	"github.com/go-curses/cdk"
	"github.com/go-curses/cdk/lib/enums"
)

const (
	TypeHBox cdk.CTypeTag = "ctk-h-box"
)

func init() {
	_ = cdk.TypesManager.AddType(TypeHBox, func() interface{} { return MakeHBox() })
}

type HBox interface {
	Box

	Init() bool
}

type CHBox struct {
	CBox
}

func MakeHBox() *CHBox {
	return NewHBox(false, 0)
}

func NewHBox(homogeneous bool, spacing int) *CHBox {
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
	b.flags = NULL_WIDGET_FLAG
	b.SetFlags(PARENT_SENSITIVE)
	b.SetFlags(APP_PAINTABLE)
	b.SetOrientation(enums.ORIENTATION_HORIZONTAL)
	return false
}
