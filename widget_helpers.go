package ctk

import (
	cenums "github.com/go-curses/cdk/lib/enums"
)

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