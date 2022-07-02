package ctk

import (
	"github.com/go-curses/cdk"
	cenums "github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/ctk/lib/enums"
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
}

var _ VBox = (*CVBox)(nil)

type CVBox struct {
	CBox
}

func MakeVBox() VBox {
	return NewVBox(false, 0)
}

func NewVBox(homogeneous bool, spacing int) VBox {
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
	b.flags = enums.NULL_WIDGET_FLAG
	b.SetFlags(enums.PARENT_SENSITIVE | enums.APP_PAINTABLE)
	b.SetOrientation(cenums.ORIENTATION_VERTICAL)
	return false
}