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
	"time"

	"github.com/go-curses/cdk"
	cenums "github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/cdk/lib/paint"
	"github.com/go-curses/cdk/lib/ptypes"
	"github.com/go-curses/cdk/memphis"

	"github.com/go-curses/ctk/lib/enums"
)

const TypeSpinner cdk.CTypeTag = "ctk-spinner"

func init() {
	_ = cdk.TypesManager.AddType(TypeSpinner, func() interface{} { return MakeSpinner() })
}

// Spinner Hierarchy:
//
//	Object
//	  +- Widget
//	    +- Misc
//	      +- Spinner
//
// The Spinner Widget needs documentation.
type Spinner interface {
	Misc
	Buildable

	StartSpinning()
	StopSpinning()
	IsSpinning() (running bool)
	IncrementSpinner()
	GetSpinnerRune() (r rune)
	SetSpinnerRunes(runes ...rune)
}

var _ Spinner = (*CSpinner)(nil)

// The CSpinner structure implements the Spinner interface and is exported to
// facilitate type embedding with custom implementations. No member variables
// are exported as the interface methods are the only intended means of
// interacting with Spinner objects.
type CSpinner struct {
	CMisc

	index int
	runes []rune

	ticker *time.Ticker
}

// MakeSpinner is used by the Buildable system to construct a new Spinner with a
// default SpinnerType setting of SpinnerRight.
func MakeSpinner() Spinner {
	return NewSpinner()
}

// NewSpinner is the constructor for new Spinner instances.
func NewSpinner() Spinner {
	s := new(CSpinner)
	s.Init()
	return s
}

// Init initializes a Spinner object. This must be called at least once to
// set up the necessary defaults and allocate any memory structures. Calling
// this more than once is safe though unnecessary. Only the first call will
// result in any effect upon the Spinner instance. Init is used in the
// NewSpinner constructor and only necessary when implementing a derivative
// Spinner type.
func (s *CSpinner) Init() bool {
	if s.InitTypeItem(TypeSpinner, s) {
		return true
	}
	s.CMisc.Init()
	s.flags = enums.NULL_WIDGET_FLAG
	s.SetFlags(enums.PARENT_SENSITIVE)
	s.SetFlags(enums.APP_PAINTABLE)
	s.Connect(SignalResize, SpinnerResizeHandle, s.resize)
	s.Connect(SignalDraw, SpinnerDrawHandle, s.draw)
	s.runes, _ = paint.GetSpinners(paint.SevenDotSpinner)
	return false
}

func (s *CSpinner) StartSpinning() {
	if s.ticker != nil {
		return
	}
	s.Lock()
	s.ticker = time.NewTicker(time.Millisecond * 250)
	s.Unlock()
	s.Emit(SignalSpinnerStart, s, string(s.GetSpinnerRune()))
	cdk.Go(func() {
		for range s.ticker.C {
			s.IncrementSpinner()
			s.Invalidate()
			if d := s.GetDisplay(); d != nil {
				d.RequestDraw()
				d.RequestShow()
			}
			if !s.IsSpinning() {
				break
			}
			s.Emit(SignalSpinnerTick, s, string(s.GetSpinnerRune()))
		}
	})
}

func (s *CSpinner) StopSpinning() {
	s.Lock()
	if s.ticker == nil {
		s.Unlock()
		return
	}
	s.ticker.Stop()
	s.ticker = nil
	s.Unlock()
	s.Emit(SignalSpinnerStop, s, string(s.GetSpinnerRune()))
}

func (s *CSpinner) IsSpinning() (running bool) {
	s.RLock()
	defer s.RUnlock()
	running = s.ticker != nil
	return
}

func (s *CSpinner) IncrementSpinner() {
	s.Lock()
	defer s.Unlock()
	if s.index += 1; s.index >= len(s.runes) {
		s.index = 0
	}
}

func (s *CSpinner) GetSpinnerRune() (r rune) {
	s.RLock()
	defer s.RUnlock()
	r = s.runes[s.index]
	return
}

func (s *CSpinner) SetSpinnerRunes(runes ...rune) {
	s.Lock()
	defer s.Unlock()
	s.runes = runes
}

// GetSizeRequest returns the requested size of the Drawable Widget. This method
// is used by Container Widgets to resolve the surface space allocated for their
// child Widget instances.
//
// Locking: read
func (s *CSpinner) GetSizeRequest() (width, height int) {
	size := ptypes.NewRectangle(s.CWidget.GetSizeRequest())
	runeWidth := 1
	xPad, yPad := s.GetPadding()
	s.RLock()
	defer s.RUnlock()
	if size.W <= -1 {
		size.W = runeWidth // variable width rune supported
		size.W += xPad * 2
	}
	if size.H <= -1 {
		size.H = 1 // always one high
		size.H += yPad * 2
	}
	size.Floor(1, 1)
	return size.W, size.H
}

func (s *CSpinner) resize(data []interface{}, argv ...interface{}) cenums.EventFlag {
	s.Invalidate()
	return cenums.EVENT_STOP
}

func (s *CSpinner) draw(data []interface{}, argv ...interface{}) cenums.EventFlag {

	if surface, ok := argv[1].(*memphis.CSurface); ok {
		alloc := s.GetAllocation()

		if !s.IsVisible() || alloc.W <= 0 || alloc.H <= 0 {
			s.LogTrace("not visible, zero width or zero height")
			return cenums.EVENT_PASS
		}

		theme := s.GetThemeRequest()
		style := theme.Content.Normal
		r := s.GetSpinnerRune()
		xAlign, yAlign := s.GetAlignment()
		xPad, yPad := s.GetPadding()
		size := ptypes.MakeRectangle(alloc.W-(xPad*2), alloc.H-(yPad*2))
		point := ptypes.MakePoint2I(xPad, yPad)

		surface.Fill(theme)

		if size.W < alloc.W {
			delta := alloc.W - size.W
			point.X = int(float64(delta) * xAlign)
		}
		if size.H < alloc.H {
			delta := alloc.H - size.H
			point.Y = int(float64(delta) * yAlign)
		}

		if err := surface.SetRune(point.X, point.Y, r, style); err != nil {
			s.LogError("set rune error: %v", err)
		}

		if debug, _ := s.GetBoolProperty(cdk.PropertyDebug); debug {
			surface.DebugBox(paint.ColorSilver, s.ObjectInfo())
		}

		return cenums.EVENT_STOP
	}
	return cenums.EVENT_PASS
}

// ArgvSpinnerEvents is a convenience function of recasting the arguments for SignalSpinnerStart, SignalSpinnerTick and
// SignalSpinnerStop events emitted by Spinner widgets
func ArgvSpinnerEvents(argv []interface{}) (spinner Spinner, symbol string) {
	if len(argv) >= 2 {
		var ok bool
		if spinner, ok = argv[0].(Spinner); ok {
			symbol, _ = argv[1].(string)
		}
	}
	return
}

const (
	SignalSpinnerStart cdk.Signal = "spinner-start"
	SignalSpinnerTick  cdk.Signal = "spinner-tick"
	SignalSpinnerStop  cdk.Signal = "spinner-stop"
)

const SpinnerResizeHandle = "spinner-resize-handler"

const SpinnerDrawHandle = "spinner-draw-handler"