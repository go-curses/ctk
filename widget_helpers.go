package ctk

func WidgetRecurseSetWindow(child Widget, window Window) {
	if child != nil {
		child.SetWindow(window)
		for _, composite := range child.GetCompositeChildren() {
			composite.SetWindow(window)
		}
		if container, ok := child.Self().(Container); ok {
			for _, grandchild := range container.GetChildren() {
				WidgetRecurseSetWindow(grandchild, window)
			}
		}
	}
}

func WidgetRecurseInvalidate(child Widget) {
	if child != nil {
		child.SetInvalidated(true)
		for _, composite := range child.GetCompositeChildren() {
			composite.SetInvalidated(true)
		}
		if container, ok := child.Self().(Container); ok {
			for _, grandchild := range container.GetChildren() {
				WidgetRecurseInvalidate(grandchild)
			}
		}
	}
}