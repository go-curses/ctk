package ctk

import (
	cenums "github.com/go-curses/cdk/lib/enums"
)

type WidgetSlice []Widget

func (ws *WidgetSlice) IndexOf(widget Widget) (idx int) {
	wid := widget.ObjectID()
	for idx = 0; idx < len(*ws); idx++ {
		if wid == (*ws)[idx].ObjectID() {
			return
		}
	}
	idx = -1
	return
}

func (ws *WidgetSlice) Append(widget Widget) {
	if idx := ws.IndexOf(widget); idx < 0 {
		*ws = append(*ws, widget)
	}
}

func (ws *WidgetSlice) Remove(widget Widget) {
	if idx := ws.IndexOf(widget); idx > -1 {
		*ws = append((*ws)[:idx], (*ws)[idx+1:]...)
	}
}

type WidgetIteratorFn = func(target Widget) cenums.EventFlag

func WidgetDescend(widget Widget, fn WidgetIteratorFn) (rv cenums.EventFlag) {
	if widget != nil {
		for _, composite := range widget.GetCompositeChildren() {
			if rv = fn(composite); rv == cenums.EVENT_STOP {
				return
			}
		}
		if container, ok := widget.Self().(Container); ok {
			for _, grandchild := range container.GetChildren() {
				if rv = WidgetDescend(grandchild, fn); rv == cenums.EVENT_STOP {
					return
				}
			}
		}
	}
	return
}

func WidgetRecurseSetWindow(widget Widget, window Window) {
	if widget != nil {
		widget.SetWindow(window)
		WidgetDescend(widget, func(target Widget) cenums.EventFlag {
			target.SetWindow(window)
			return cenums.EVENT_PASS
		})
	}
}

func WidgetRecurseInvalidate(widget Widget) {
	if widget != nil {
		widget.SetInvalidated(true)
		WidgetDescend(widget, func(target Widget) cenums.EventFlag {
			target.SetInvalidated(true)
			return cenums.EVENT_PASS
		})
	}
}