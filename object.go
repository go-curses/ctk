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

	"github.com/go-curses/cdk"
	"github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/cdk/lib/paint"
	"github.com/go-curses/cdk/lib/ptypes"
)

// CDK type tag for CTK Objects
const TypeObject cdk.CTypeTag = "ctk-object"

func init() {
	_ = cdk.TypesManager.AddType(TypeObject, nil)
}

// In the Curses Tool Kit, the Object type is an extension of the CDK Object
// type and for all intents and purposes, this is the base class for any CTK
// type with no other CTK type embedding a CDK type directly
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
	Draw() enums.EventFlag
	Resize() enums.EventFlag
	GetTextDirection() (direction TextDirection)
	SetTextDirection(direction TextDirection)
	CssSelector() (selector string)
	InstallCssProperty(name cdk.Property, kind cdk.PropertyType, write bool, def interface{}) (err error)
	GetCssProperty(name cdk.Property) (property *cdk.CProperty)
	GetCssProperties() (properties []*cdk.CProperty)
	GetCssBool(name cdk.Property) (value bool, err error)
	GetCssString(name cdk.Property) (value string, err error)
	GetCssInt(name cdk.Property) (value int, err error)
	GetCssFloat(name cdk.Property) (value float64, err error)
	GetCssColor(name cdk.Property) (value paint.Color, err error)
}

// The CObject structure implements the Object interface and is exported to
// facilitate type embedding with custom implementations. No member variables
// are exported as the interface methods are the only intended means of
// interacting with CTK objects
type CObject struct {
	cdk.CObject

	origin        *ptypes.Point2I
	allocation    *ptypes.Rectangle
	textDirection TextDirection
	css           map[cdk.Property]*cdk.CProperty
}

// CTK object initialization. This must be called at least once to setup the
// necessary defaults and allocate any memory structures. Calling this more
// than once is safe though unnecessary. Only the first call will result in any
// effect upon the Object instance
func (o *CObject) Init() (already bool) {
	if o.InitTypeItem(TypeObject, o) {
		return true
	}
	o.CObject.Init()
	o.origin = &ptypes.Point2I{}
	o.allocation = &ptypes.Rectangle{}
	o.css = make(map[cdk.Property]*cdk.CProperty)
	_ = o.InstallProperty(PropertyParent, cdk.StructProperty, true, nil)
	_ = o.InstallCssProperty(PropertyClass, cdk.StringProperty, true, "")
	_ = o.InstallCssProperty(PropertyWidth, cdk.IntProperty, true, -1)
	_ = o.InstallCssProperty(PropertyHeight, cdk.IntProperty, true, -1)
	_ = o.InstallCssProperty(PropertyColor, cdk.ColorProperty, true, "#ffffff")
	_ = o.InstallCssProperty(PropertyBackgroundColor, cdk.ColorProperty, true, "#000000")
	_ = o.InstallCssProperty(PropertyBackgroundFillContent, cdk.StringProperty, true, " ")
	_ = o.InstallCssProperty(PropertyBorder, cdk.BoolProperty, true, false)
	_ = o.InstallCssProperty(PropertyBorderColor, cdk.ColorProperty, true, "#ffffff")
	_ = o.InstallCssProperty(PropertyBorderBackgroundColor, cdk.ColorProperty, true, "#000000")
	_ = o.InstallCssProperty(PropertyBold, cdk.BoolProperty, true, false)
	_ = o.InstallCssProperty(PropertyBlink, cdk.BoolProperty, true, false)
	_ = o.InstallCssProperty(PropertyReverse, cdk.BoolProperty, true, false)
	_ = o.InstallCssProperty(PropertyUnderline, cdk.BoolProperty, true, false)
	_ = o.InstallCssProperty(PropertyDim, cdk.BoolProperty, true, false)
	_ = o.InstallCssProperty(PropertyItalic, cdk.BoolProperty, true, false)
	_ = o.InstallCssProperty(PropertyStrike, cdk.BoolProperty, true, false)
	_ = o.InstallCssProperty(PropertyBorderTopLeftContent, cdk.StringProperty, true, "+")
	_ = o.InstallCssProperty(PropertyBorderTopContent, cdk.StringProperty, true, "-")
	_ = o.InstallCssProperty(PropertyBorderTopRightContent, cdk.StringProperty, true, "+")
	_ = o.InstallCssProperty(PropertyBorderLeftContent, cdk.StringProperty, true, "|")
	_ = o.InstallCssProperty(PropertyBorderRightContent, cdk.StringProperty, true, "|")
	_ = o.InstallCssProperty(PropertyBorderBottomLeftContent, cdk.StringProperty, true, "+")
	_ = o.InstallCssProperty(PropertyBorderBottomContent, cdk.StringProperty, true, "-")
	_ = o.InstallCssProperty(PropertyBorderBottomRightContent, cdk.StringProperty, true, "+")
	return false
}

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

// Convenience method to return a string identifying the Object instance with
// it's type, unique identifier, name if set (see SetName()), the origin point
// and current size allocation
func (o *CObject) ObjectInfo() string {
	info := fmt.Sprintf("%v %v,%v %v", o.ObjectID(), o.origin, o.GetAllocation(), o.ObjectName())
	return info
}

// Set the origin of this instance in display space. This method emits an origin
// signal initially and if the listeners return EVENT_PASS then the change is
// applied
//
// Emits: SignalOrigin, Argv=[Object instance, new origin]
func (o *CObject) SetOrigin(x, y int) {
	if f := o.Emit(SignalOrigin, o, ptypes.MakePoint2I(x, y)); f == enums.EVENT_PASS {
		o.origin.Set(x, y)
	}
}

// Returns the current origin point of the Object instance
func (o *CObject) GetOrigin() ptypes.Point2I {
	return *o.origin
}

// Set the allocated size of the Object instance. This method is only useful for
// custom CTK types that need to render child Widgets. This method emits an
// allocation signal initially and if the listeners return EVENT_PASS the change
// is applied and constrained to a minimum width and height of zero
func (o *CObject) SetAllocation(size ptypes.Rectangle) {
	if f := o.Emit(SignalAllocation, o.allocation, size); f == enums.EVENT_PASS {
		o.allocation.Set(size.W, size.H)
		o.allocation.Floor(0, 0)
	}
}

// Returns the current allocation size of the Object instance
func (o *CObject) GetAllocation() ptypes.Rectangle {
	return *o.allocation
}

// Returns the Object's instance if the given point is within the Object's
// display space bounds. This method is mainly used by Window objects and other
// event processing Widgets that need to find a Widget by mouse-cursor
// coordinates for example. If this Object does not encompasses the point given,
// it returns nil
func (o *CObject) GetObjectAt(p *ptypes.Point2I) Object {
	if o.HasPoint(p) {
		return o
	}
	return nil
}

// Determines whether or not the given point is within the Object's display
// space bounds
func (o *CObject) HasPoint(p *ptypes.Point2I) (contains bool) {
	origin := o.GetOrigin()
	size := o.GetAllocation()
	if p.X >= origin.X && p.X < (origin.X+size.W) {
		if p.Y >= origin.Y && p.Y < (origin.Y+size.H) {
			return true
		}
	}
	return false
}

// Emits an invalidate signal, primarily used in other CTK types which are
// drawable and need an opportunity to invalidate the render canvases so that
// the next CTK draw cycle can reflect the latest changes to the Object instance
//
// Emits: SignalInvalidate, Argv=[Object instance]
func (o *CObject) Invalidate() enums.EventFlag {
	return o.Emit(SignalInvalidate, o)
}

// Emits a cdk-event signal, primarily used to consume CDK events received such
// as mouse or key events in other CTK and custom types that embed CObject
//
// Emits: SignalCdkEvent, Argv=[Object instance, event]
func (o *CObject) ProcessEvent(evt cdk.Event) enums.EventFlag {
	return o.Emit(SignalCdkEvent, o, evt)
}

// Emits a draw signal, primarily used to render canvases and cause end-user
// facing display updates. Signal listeners can draw to the Canvas and return
// EVENT_STOP to cause those changes to be composited upon the larger display
// canvas
//
// Emits: SignalDraw, Argv=[Object instance, canvas]
func (o *CObject) Draw() enums.EventFlag {
	return o.Emit(SignalDraw, o, nil)
}

// Emits a resize signal, primarily used to make adjustments or otherwise
// reallocate resources necessary for subsequent draw events
//
// Emits: SignalResize, Argv=[Object instance, allocated size]
func (o *CObject) Resize() enums.EventFlag {
	size := o.GetAllocation()
	return o.Emit(SignalResize, o, size)
}

// Returns the current text direction for this Object instance
func (o *CObject) GetTextDirection() (direction TextDirection) {
	return o.textDirection
}

// Set the text direction for this Object instance. This method emits a
// text-direction signal initially and if the listeners return EVENT_PASS, the
// change is applied
//
// Emit: SignalTextDirection, Argv=[Object instance, new text direction]
func (o *CObject) SetTextDirection(direction TextDirection) {
	if f := o.Emit(SignalTextDirection, o, direction); f == enums.EVENT_PASS {
		o.textDirection = direction
	}
}

// return a selector string identifying this exact Object instance.
func (o *CObject) CssSelector() (selector string) {
	selector += o.GetTypeTag().String()
	name := o.GetName()
	if name != "" {
		selector += "#" + name
	}
	return
}

func (o *CObject) InstallCssProperty(name cdk.Property, kind cdk.PropertyType, write bool, def interface{}) (err error) {
	switch kind {
	case cdk.BoolProperty, cdk.StringProperty, cdk.IntProperty, cdk.FloatProperty, cdk.ColorProperty:
	default:
		return fmt.Errorf("unsupported css property type: %v", kind)
	}
	if existing := o.GetCssProperty(name); existing != nil {
		return fmt.Errorf("property exists: %v", name)
	}
	o.css[name] = cdk.NewProperty(name, kind, write, false, def)
	return nil
}

func (o *CObject) GetCssProperty(name cdk.Property) (property *cdk.CProperty) {
	var ok bool
	if property, ok = o.css[name]; !ok {
		property = nil
	}
	return
}

func (o *CObject) GetCssProperties() (properties []*cdk.CProperty) {
	for _, v := range o.css {
		properties = append(properties, v)
	}
	return
}

// func (o *CObject) GetCssValue(name Property) (value interface{}) {
// 	if v, ok := o.css[name]; ok {
// 		value = v.Value()
// 	}
// 	return
// }

func (o *CObject) GetCssBool(name cdk.Property) (value bool, err error) {
	if prop := o.GetCssProperty(name); prop != nil {
		if prop.Type() == cdk.BoolProperty {
			if v, ok := prop.Value().(bool); ok {
				return v, nil
			}
			if v, ok := prop.Default().(bool); ok {
				return v, nil
			}
		}
		err = fmt.Errorf("%v.(%v) css property is not a bool", name, prop.Type())
		return
	}
	err = fmt.Errorf("css property not found: %v", name)
	return
}

func (o *CObject) GetCssString(name cdk.Property) (value string, err error) {
	if prop := o.GetCssProperty(name); prop != nil {
		if prop.Type() == cdk.StringProperty {
			if v, ok := prop.Value().(string); ok {
				return v, nil
			}
			if v, ok := prop.Default().(string); ok {
				return v, nil
			}
		}
		err = fmt.Errorf("%v.(%v) css property is not a string", name, prop.Type())
		return
	}
	err = fmt.Errorf("css property not found: %v", name)
	return
}

func (o *CObject) GetCssInt(name cdk.Property) (value int, err error) {
	if prop := o.GetCssProperty(name); prop != nil {
		if prop.Type() == cdk.IntProperty {
			if v, ok := prop.Value().(int); ok {
				return v, nil
			}
			if v, ok := prop.Default().(int); ok {
				return v, nil
			}
		}
		err = fmt.Errorf("%v.(%v) css property is not a int", name, prop.Type())
		return
	}
	err = fmt.Errorf("css property not found: %v", name)
	return
}

func (o *CObject) GetCssFloat(name cdk.Property) (value float64, err error) {
	if prop := o.GetCssProperty(name); prop != nil {
		if prop.Type() == cdk.FloatProperty {
			if v, ok := prop.Value().(float64); ok {
				return v, nil
			}
			if v, ok := prop.Default().(float64); ok {
				return v, nil
			}
		}
		err = fmt.Errorf("%v.(%v) css property is not a float64", name, prop.Type())
		return
	}
	err = fmt.Errorf("css property not found: %v", name)
	return
}

func (o *CObject) GetCssColor(name cdk.Property) (value paint.Color, err error) {
	if prop := o.GetCssProperty(name); prop != nil {
		if prop.Type() == cdk.ColorProperty {
			if v, ok := prop.Value().(paint.Color); ok {
				return v, nil
			}
			if v, ok := prop.Default().(paint.Color); ok {
				return v, nil
			}
		}
		err = fmt.Errorf("%v.(%v) css property is not a Color", name, prop.Type())
		return
	}
	err = fmt.Errorf("css property not found: %v", name)
	return
}

// grouping label
const PropertyClass cdk.Property = "class"

// convenience wrapper for cdk.PropertyDebug
const PropertyDebug cdk.Property = cdk.PropertyDebug

// css properties
const PropertyWidth cdk.Property = "width"
const PropertyHeight cdk.Property = "height"
const PropertyColor cdk.Property = "color"
const PropertyBackgroundColor cdk.Property = "background-color"
const PropertyBackgroundFillContent cdk.Property = "background-fill-content"
const PropertyBorder cdk.Property = "border"
const PropertyBorderColor cdk.Property = "border-color"
const PropertyBorderBackgroundColor cdk.Property = "border-background-color"
const PropertyBold cdk.Property = "bold"
const PropertyBlink cdk.Property = "blink"
const PropertyReverse cdk.Property = "reverse"
const PropertyUnderline cdk.Property = "underline"
const PropertyDim cdk.Property = "dim"
const PropertyItalic cdk.Property = "italic"
const PropertyStrike cdk.Property = "strike"
const PropertyBorderTopLeftContent cdk.Property = "border-top-left-content"
const PropertyBorderTopContent cdk.Property = "border-top-content"
const PropertyBorderTopRightContent cdk.Property = "border-top-right-content"
const PropertyBorderBottomLeftContent cdk.Property = "border-bottom-left-content"
const PropertyBorderBottomRightContent cdk.Property = "border-bottom-right-content"
const PropertyBorderLeftContent cdk.Property = "border-left-content"
const PropertyBorderRightContent cdk.Property = "border-right-content"
const PropertyBorderBottomContent cdk.Property = "border-bottom-content"
