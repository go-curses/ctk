package ctk

import (
	"github.com/go-curses/cdk"
	"github.com/go-curses/cdk/lib/enums"
)

const (
	TypeVBox cdk.CTypeTag = "ctk-v-box"
)

func init() {
	_ = cdk.TypesManager.AddType(TypeVBox, func() interface{} { return MakeVBox() })
}

// Basic vbox interface
type VBox interface {
	Box

	Init() bool
}

type CVBox struct {
	CBox
}

func MakeVBox() *CVBox {
	return NewVBox(false, 0)
}

func NewVBox(homogeneous bool, spacing int) *CVBox {
	b := new(CVBox)
	b.Init()
	b.SetHomogeneous(homogeneous)
	b.SetSpacing(spacing)
	return b
}

func (b *CVBox) Init() bool {
	if b.InitTypeItem(TypeVBox, b) {
		return true
	}
	b.CBox.Init()
	b.flags = NULL_WIDGET_FLAG
	b.SetFlags(PARENT_SENSITIVE | APP_PAINTABLE)
	b.SetOrientation(enums.ORIENTATION_VERTICAL)
	return false
}
