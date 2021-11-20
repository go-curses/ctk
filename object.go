// Copyright 2020 The CDK Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use file except in compliance with the License.
// You may obtain a copy of the license at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ctk

import (
	"fmt"
	"strings"

	"github.com/go-curses/cdk"
	"github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/cdk/lib/paint"
	"github.com/go-curses/cdk/lib/ptypes"
)

const TypeObject cdk.CTypeTag = "ctk-object"

func init() {
	_ = cdk.TypesManager.AddType(TypeObject, nil)
}

// Object in the Curses Tool Kit, is an extension of the CDK Object type and for
// all intents and purposes, this is the base class for any CTK type with no
// other CTK type embedding a CDK type directly.
type Object interface {
	cdk.Object

	Init() (already bool)
	Build(builder Builder, element *CBuilderElement) error
	ObjectInfo() string
	SetOrigin(x, y int)
	GetOrigin() ptypes.Point2I
	SetAllocation(size ptypes.Rectangle)
	GetAllocation() ptypes.Rectangle
	GetObjectAt(p *ptypes.Point2I) Object
	HasPoint(p *ptypes.Point2I) (contains bool)
	Invalidate() enums.EventFlag
	ProcessEvent(evt cdk.Event) enums.EventFlag
	Resize() enums.EventFlag
	GetTextDirection() (direction TextDirection)
	SetTextDirection(direction TextDirection)
	CssSelector() (selector string)
	InstallCssProperty(name cdk.Property, state StateType, kind cdk.PropertyType, write bool, def interface{}) (err error)
	SetCssPropertyFromStyle(key, value string) (err error)
	GetCssProperty(name cdk.Property, state StateType) (property *CStyleProperty)
	GetCssProperties() (properties map[StateType][]*CStyleProperty)
	GetCssBool(name cdk.Property, state StateType) (value bool, err error)
	GetCssString(name cdk.Property, state StateType) (value string, err error)
	GetCssInt(name cdk.Property, state StateType) (value int, err error)
	GetCssFloat(name cdk.Property, state StateType) (value float64, err error)
	GetCssColor(name cdk.Property, state StateType) (value paint.Color, err error)
}

// The CObject structure implements the Object interface and is exported
// to facilitate type embedding with custom implementations. No member variables
// are exported as the interface methods are the only intended means of
// interacting with Object objects.
type CObject struct {
	cdk.CObject

	origin        *ptypes.Point2I
	allocation    *ptypes.Rectangle
	textDirection TextDirection
	css           map[StateType]map[cdk.Property]*CStyleProperty
}

// Init initializes an Object instance. This must be called at least once to
// set up the necessary defaults and allocate any memory structures. Calling
// this more than once is safe though unnecessary. Only the first call will
// result in any effect upon the Object instance. Init is used in the
// NewObject constructor and only necessary when implementing a derivative
// Object type.
func (o *CObject) Init() (already bool) {
	if o.InitTypeItem(TypeObject, o) {
		return true
	}
	o.CObject.Init()
	o.origin = &ptypes.Point2I{}
	o.allocation = &ptypes.Rectangle{}
	o.css = make(map[StateType]map[cdk.Property]*CStyleProperty)
	_ = o.InstallProperty(PropertyParent, cdk.StructProperty, true, nil)
	return false
}

// Build provides customizations to the Buildable system for Object Widgets.
func (o *CObject) Build(builder Builder, element *CBuilderElement) error {
	o.Freeze()
	defer o.Thaw()
	if name, ok := element.Attributes["id"]; ok {
		o.SetName(name)
	} else if name, ok := element.Attributes["name"]; ok {
		o.SetName(name)
	}
	tt := o.GetTypeTag().Tag()
	if fn, ok := ctkBuilderTranslators[tt]; ok {
		for k, v := range element.Properties {
			if widget, ok := element.Instance.(Widget); ok {
				if err := fn(builder, widget, k, v); err != nil {
					if err != ErrFallthrough {
						return fmt.Errorf("%v property translator error: %v", tt, err)
					} else {
						element.ApplyProperty(k, v)
					}
				}
			}
		}
	} else {
		element.ApplyProperties()
	}
	element.ApplySignals()
	return nil
}

// ObjectInfo is a convenience method to return a string identifying the Object
// instance with its type, unique identifier, name if set (see SetName()), the
// origin point and current size allocation.
func (o *CObject) ObjectInfo() string {
	o.RLock()
	origin := o.origin
	o.RUnlock()
	info := fmt.Sprintf("%v %v,%v %v", o.ObjectID(), origin, o.GetAllocation(), o.ObjectName())
	return info
}

// SetOrigin updates the origin of this instance in display space. This method
// emits an origin signal initially and if the listeners return EVENT_PASS then
// the change is applied.
//
// Emits: SignalOrigin, Argv=[Object instance, new origin]
func (o *CObject) SetOrigin(x, y int) {
	if f := o.Emit(SignalOrigin, o, ptypes.MakePoint2I(x, y)); f == enums.EVENT_PASS {
		o.Lock()
		o.origin.Set(x, y)
		o.Unlock()
	}
}

// GetOrigin returns the current origin point of the Object instance
func (o *CObject) GetOrigin() (origin ptypes.Point2I) {
	o.RLock()
	origin = o.origin.Clone()
	o.RUnlock()
	return
}

// SetAllocation updates the allocated size of the Object instance. This method
// is only useful for custom CTK types that need to render child Widgets. This
// method emits an allocation signal initially and if the listeners return
// EVENT_PASS the change is applied and constrained to a minimum width and
// height of zero.
func (o *CObject) SetAllocation(size ptypes.Rectangle) {
	if f := o.Emit(SignalAllocation, o.allocation, size); f == enums.EVENT_PASS {
		o.Lock()
		o.allocation.Set(size.W, size.H)
		o.allocation.Floor(0, 0)
		o.Unlock()
	}
}

// GetAllocation returns the current allocation size of the Object instance.
func (o *CObject) GetAllocation() (alloc ptypes.Rectangle) {
	o.RLock()
	alloc = o.allocation.Clone()
	o.RUnlock()
	return
}

// GetObjectAt returns the Object's instance if the given point is within the
// Object's display space bounds. This method is mainly used by Window objects
// and other event processing Widgets that need to find a Widget by mouse-cursor
// coordinates for example. If this Object does not encompass the point given,
// it returns `nil`.
func (o *CObject) GetObjectAt(p *ptypes.Point2I) Object {
	if o.HasPoint(p) {
		return o
	}
	return nil
}

// HasPoint determines whether the given point is within the Object's display
// space bounds.
func (o *CObject) HasPoint(p *ptypes.Point2I) (contains bool) {
	origin := o.GetOrigin()
	size := o.GetAllocation()
	o.RLock()
	if p.X >= origin.X && p.X < (origin.X+size.W) {
		if p.Y >= origin.Y && p.Y < (origin.Y+size.H) {
			o.RUnlock()
			return true
		}
	}
	o.RUnlock()
	return false
}

// Invalidate emits an invalidate signal, primarily used in other CTK types
// which are drawable and need an opportunity to invalidate the memphis surfaces
// so that the next CTK draw cycle can reflect the latest changes to the Object
// instance.
//
// Locking: expected read/write
func (o *CObject) Invalidate() enums.EventFlag {
	return o.Emit(SignalInvalidate, o)
}

// ProcessEvent emits a cdk-event signal, primarily used to consume CDK events
// received such as mouse or key events in other CTK and custom types that embed
// CObject.
//
// Locking: expected read/write
func (o *CObject) ProcessEvent(evt cdk.Event) enums.EventFlag {
	return o.Emit(SignalCdkEvent, o, evt)
}

// Resize emits a resize signal, primarily used to make adjustments or otherwise
// reallocate resources necessary for subsequent draw events.
//
// Locking: read
func (o *CObject) Resize() enums.EventFlag {
	size := o.GetAllocation()
	return o.Emit(SignalResize, o, size)
}

// GetTextDirection returns the current text direction for this Object instance.
//
// Note that usage of this within CTK is unimplemented at this time
func (o *CObject) GetTextDirection() (direction TextDirection) {
	o.RLock()
	direction = o.textDirection
	o.RUnlock()
	return
}

// SetTextDirection updates text direction for this Object instance. This method
// emits a text-direction signal initially and if the listeners return
// EVENT_PASS, the change is applied.
//
// Note that usage of this within CTK is unimplemented at this time
func (o *CObject) SetTextDirection(direction TextDirection) {
	if f := o.Emit(SignalTextDirection, o, direction); f == enums.EVENT_PASS {
		o.Lock()
		o.textDirection = direction
		o.Unlock()
	}
}

// CssSelector returns a selector string identifying this exact Object instance.
func (o *CObject) CssSelector() (selector string) {
	selector += o.GetTypeTag().String()
	name := o.GetName()
	if name != "" {
		selector += "#" + name
	}
	return
}

// InstallCssProperty installs a new cdk.Property in a secondary CSS-focused
// property list.
//
// Note that usage of this within CTK is unimplemented at this time
func (o *CObject) InstallCssProperty(name cdk.Property, state StateType, kind cdk.PropertyType, write bool, def interface{}) (err error) {
	switch kind {
	case cdk.BoolProperty, cdk.StringProperty, cdk.IntProperty, cdk.FloatProperty, cdk.ColorProperty:
	default:
		return fmt.Errorf("unsupported css property type: %v", kind)
	}
	if existing := o.GetCssProperty(name, state); existing != nil {
		return fmt.Errorf("property exists: %v", name)
	}
	o.Lock()
	if _, ok := o.css[state]; !ok {
		o.css[state] = make(map[cdk.Property]*CStyleProperty)
	}
	o.css[state][name] = NewStyleProperty(name, state, kind, write, false, def)
	o.Unlock()
	return nil
}

func (o *CObject) SetCssPropertyFromStyle(key, value string) (err error) {
	o.Lock()
	state := StateNormal
	if strings.Contains(key, ":") {
		parts := strings.Split(key, ":")
		state = StateTypeFromString(parts[1])
		key = parts[0]
	}
	if property, ok := o.css[state][cdk.Property(key)]; ok {
		err = property.SetFromString(value)
	} else {
		err = fmt.Errorf("css property not found: %v", key)
	}
	o.Unlock()
	return
}

// GetCssProperty returns the cdk.Property instance of the property found with
// the name given, returning `nil` if no property by the name given is found.
//
// Note that usage of this within CTK is unimplemented at this time
func (o *CObject) GetCssProperty(name cdk.Property, state StateType) (property *CStyleProperty) {
	o.RLock()
	var ok bool
	if property, ok = o.css[state][name]; !ok {
		property = nil
	}
	o.RUnlock()
	return
}

// GetCssProperties returns all the installed CSS properties for the Object.
//
// Note that usage of this within CTK is unimplemented at this time
func (o *CObject) GetCssProperties() (properties map[StateType][]*CStyleProperty) {
	o.RLock()
	properties = make(map[StateType][]*CStyleProperty)
	for s, _ := range o.css {
		properties[s] = make([]*CStyleProperty, len(o.css[s]))
		for _, v := range o.css[s] {
			properties[s] = append(properties[s], v)
		}
	}
	o.RUnlock()
	return
}

// GetCssValue returns the value of the property found with the same name as the
// given name.
//
// Note that usage of this within CTK is unimplemented at this time
func (o *CObject) GetCssValue(name cdk.Property, state StateType) (value interface{}) {
	o.RLock()
	if v, ok := o.css[state][name]; ok {
		value = v.Value()
	}
	o.RUnlock()
	return
}

// GetCssBool is a convenience method to return a boolean value for the CSS
// property of the given name.
//
// Note that usage of this within CTK is unimplemented at this time
func (o *CObject) GetCssBool(name cdk.Property, state StateType) (value bool, err error) {
	if prop := o.GetCssProperty(name, state); prop != nil {
		o.RLock()
		if prop.Type() == cdk.BoolProperty {
			if v, ok := prop.Value().(bool); ok {
				o.RUnlock()
				return v, nil
			}
			if v, ok := prop.Default().(bool); ok {
				o.RUnlock()
				return v, nil
			}
		}
		err = fmt.Errorf("%v.(%v) css property is not a bool", name, prop.Type())
		o.RUnlock()
		return
	}
	err = fmt.Errorf("css property not found: %v", name)
	return
}

// GetCssString is a convenience method to return a string value for the CSS
// property of the given name.
//
// Note that usage of this within CTK is unimplemented at this time
func (o *CObject) GetCssString(name cdk.Property, state StateType) (value string, err error) {
	if prop := o.GetCssProperty(name, state); prop != nil {
		o.RLock()
		if prop.Type() == cdk.StringProperty {
			if v, ok := prop.Value().(string); ok {
				o.RUnlock()
				return v, nil
			}
			if v, ok := prop.Default().(string); ok {
				o.RUnlock()
				return v, nil
			}
		}
		err = fmt.Errorf("%v.(%v) css property is not a string", name, prop.Type())
		o.RUnlock()
		return
	}
	err = fmt.Errorf("css property not found: %v", name)
	return
}

// GetCssInt is a convenience method to return a int value for the CSS
// property of the given name.
//
// Note that usage of this within CTK is unimplemented at this time
func (o *CObject) GetCssInt(name cdk.Property, state StateType) (value int, err error) {
	if prop := o.GetCssProperty(name, state); prop != nil {
		o.RLock()
		if prop.Type() == cdk.IntProperty {
			if v, ok := prop.Value().(int); ok {
				o.RUnlock()
				return v, nil
			}
			if v, ok := prop.Default().(int); ok {
				o.RUnlock()
				return v, nil
			}
		}
		err = fmt.Errorf("%v.(%v) css property is not a int", name, prop.Type())
		o.RUnlock()
		return
	}
	err = fmt.Errorf("css property not found: %v", name)
	return
}

// GetCssFloat is a convenience method to return a float value for the CSS
// property of the given name.
//
// Note that usage of this within CTK is unimplemented at this time
func (o *CObject) GetCssFloat(name cdk.Property, state StateType) (value float64, err error) {
	if prop := o.GetCssProperty(name, state); prop != nil {
		o.RLock()
		if prop.Type() == cdk.FloatProperty {
			if v, ok := prop.Value().(float64); ok {
				o.RUnlock()
				return v, nil
			}
			if v, ok := prop.Default().(float64); ok {
				o.RUnlock()
				return v, nil
			}
		}
		err = fmt.Errorf("%v.(%v) css property is not a float64", name, prop.Type())
		o.RUnlock()
		return
	}
	err = fmt.Errorf("css property not found: %v", name)
	return
}

// GetCssColor is a convenience method to return a paint.Color value for the CSS
// property of the given name.
//
// Note that usage of this within CTK is unimplemented at this time
func (o *CObject) GetCssColor(name cdk.Property, state StateType) (value paint.Color, err error) {
	if prop := o.GetCssProperty(name, state); prop != nil {
		o.RLock()
		if prop.Type() == cdk.ColorProperty {
			if v, ok := prop.Value().(paint.Color); ok {
				o.RUnlock()
				return v, nil
			}
			if v, ok := prop.Default().(paint.Color); ok {
				o.RUnlock()
				return v, nil
			}
		}
		err = fmt.Errorf("%v.(%v) css property is not a Color", name, prop.Type())
		o.RUnlock()
		return
	}
	err = fmt.Errorf("css property not found: %v", name)
	return
}

const PropertyDebug cdk.Property = cdk.PropertyDebug
