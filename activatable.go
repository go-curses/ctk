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
