package ctk

import (
	"github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/cdk/lib/paint"
	"github.com/go-curses/cdk/lib/ptypes"
)

type Drawable interface {
	Hide()
	Show()
	ShowAll()
	IsVisible() bool
	HasPoint(p *ptypes.Point2I) bool
	GetWidgetAt(p *ptypes.Point2I) (instance interface{})
	GetSizeRequest() (size ptypes.Rectangle)
	SetSizeRequest(x, y int)
	GetTheme() (theme paint.Theme)
	SetTheme(theme paint.Theme)
	GetThemeRequest() (theme paint.Theme)
	GetOrigin() (origin ptypes.Point2I)
	SetOrigin(x, y int)
	GetAllocation() (alloc ptypes.Rectangle)
	SetAllocation(alloc ptypes.Rectangle)
	Invalidate() enums.EventFlag
	Resize() enums.EventFlag
	Draw() enums.EventFlag
}
