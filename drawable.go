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