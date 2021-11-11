package ctk

import (
	"github.com/go-curses/cdk"
	"github.com/go-curses/cdk/lib/enums"
)

type Orientable interface {
	GetOrientation() (orientation enums.Orientation)
	SetOrientation(orientation enums.Orientation)
}

// The orientation of the orientable.
// Flags: Read / Write
// Default value: ORIENTATION_HORIZONTAL
const PropertyOrientation cdk.Property = "orientation"
