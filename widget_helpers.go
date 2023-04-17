// Copyright (c) 2023  The Go-Curses Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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