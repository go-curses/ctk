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
	"github.com/go-curses/cdk"
	"github.com/go-curses/cdk/lib/enums"
)

// Activatable Hierarchy:
//	CInterface
//	  +- Activatable
type Activatable interface {
	Activate() (value bool)
	Clicked() enums.EventFlag
	GrabFocus()
}

const SignalActivate cdk.Signal = "activate"

// The action that this activatable will activate and receive updates from
// for various states and possibly appearance.
// Flags: Read / Write
const PropertyRelatedAction cdk.Property = "related-action"

// Whether this activatable should reset its layout and appearance when
// setting the related action or when the action changes appearance. See the
// Action documentation directly to find which properties should be
// ignored by the Activatable when this property is FALSE.
// Flags: Read / Write
// Default value: TRUE
const PropertyUseActionAppearance cdk.Property = "use-action-appearance"