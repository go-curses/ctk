package ctk

import (
	"github.com/go-curses/cdk"
	cenums "github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/ctk/lib/enums"
)

const TypeSeparator cdk.CTypeTag = "ctk-separator"

func init() {
	_ = cdk.TypesManager.AddType(TypeSeparator, func() interface{} { return MakeSeparator() })
}

var _ Separator = (*CSeparator)(nil)

// Separator Hierarchy:
//	Object
//	  +- Widget
//	    +- Container
//	      +- Bin
//	        +- Separator
//
// The Separator Widget is used to capture Widget events (mouse, keyboard)
// without needing having any defined user-interface.
type Separator interface {
	Bin
	Buildable
	Orientable

	Init() (already bool)
}

// The CSeparator structure implements the Separator interface and is exported
// to facilitate type embedding with custom implementations. No member variables
// are exported as the interface methods are the only intended means of
// interacting with Separator objects.
type CSeparator struct {
	CBin
}

// MakeSeparator is used by the Buildable system to construct a new Separator.
func MakeSeparator() Separator {
	return NewSeparator()
}

// NewSeparator is the constructor for new Separator instances.
func NewSeparator() (value Separator) {
	s := new(CSeparator)
	s.Init()
	return s
}

// Init initializes a Separator object. This must be called at least once to
// set up the necessary defaults and allocate any memory structures. Calling
// this more than once is safe though unnecessary. Only the first call will
// result in any effect upon the Separator instance. Init is used in the
// NewSeparator constructor and only necessary when implementing a derivative
// Separator type.
func (s *CSeparator) Init() (already bool) {
	if s.InitTypeItem(TypeSeparator, s) {
		return true
	}
	s.CBin.Init()
	s.SetFlags(enums.APP_PAINTABLE)
	_ = s.InstallBuildableProperty(PropertyOrientation, cdk.StructProperty, true, cenums.ORIENTATION_HORIZONTAL)
	s.Connect(SignalDraw, TextFieldDrawHandle, s.draw)
	return false
}

// GetOrientation is a convenience method for returning the orientation property
// value.
// See: SetOrientation()
//
// Locking: read
func (s *CSeparator) GetOrientation() (orientation cenums.Orientation) {
	s.RLock()
	defer s.RUnlock()
	var ok bool
	if v, err := s.GetStructProperty(PropertyOrientation); err != nil {
		s.LogErr(err)
	} else if orientation, ok = v.(cenums.Orientation); !ok && v != nil {
		s.LogError("invalid value stored in %v: %v (%T)", PropertyOrientation, v, v)
	}
	return
}

// SetOrientation is a convenience method for updating the orientation property
// value.
//
// Parameters:
//  orientation  the desired cenums.Orientation to use
//
// Locking: write
func (s *CSeparator) SetOrientation(orientation cenums.Orientation) {
	s.Lock()
	if err := s.SetStructProperty(PropertyOrientation, orientation); err != nil {
		s.LogErr(err)
	}
	s.Unlock()
	s.Resize()
}

func (s *CSeparator) draw(data []interface{}, argv ...interface{}) cenums.EventFlag {
	return cenums.EVENT_STOP
}