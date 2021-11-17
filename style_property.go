// Copyright 2021  The CTK Authors
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
	"strconv"
	"strings"

	"github.com/go-curses/cdk"
	"github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/cdk/lib/paint"
	"github.com/go-curses/cdk/lib/ptypes"
)

type StyleProperty string

func (p StyleProperty) String() string {
	return string(p)
}

type CStyleProperty struct {
	name      cdk.Property
	state     StateType
	kind      cdk.PropertyType
	write     bool
	buildable bool
	def       interface{}
	value     interface{}
}

func NewStyleProperty(name cdk.Property, state StateType, kind cdk.PropertyType, write bool, buildable bool, def interface{}) (property *CStyleProperty) {
	property = new(CStyleProperty)
	property.name = name
	property.state = state
	property.kind = kind
	property.write = write
	property.buildable = buildable
	property.def = def
	property.value = def
	return
}

func (p *CStyleProperty) Clone() *CStyleProperty {
	return &CStyleProperty{
		name:      p.name,
		state:     p.state,
		kind:      p.kind,
		write:     p.write,
		buildable: p.buildable,
		def:       p.def,
		value:     p.value,
	}
}

func (p *CStyleProperty) Name() cdk.Property {
	return p.name
}

func (p *CStyleProperty) State() StateType {
	return p.state
}

func (p *CStyleProperty) Type() cdk.PropertyType {
	return p.kind
}

func (p *CStyleProperty) ReadOnly() bool {
	return !p.write
}

func (p *CStyleProperty) Buildable() bool {
	return p.buildable
}

func (p *CStyleProperty) Set(value interface{}) error {
	if !p.write {
		return fmt.Errorf("cannot change read-only property: %v", p.name)
	}
	switch p.Type() {
	case cdk.BoolProperty:
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("%v value is not of bool type: %v (%T)", p.name, value, value)
		}
	case cdk.StringProperty:
		if _, ok := value.(string); !ok {
			return fmt.Errorf("%v value is not of string type: %v (%T)", p.name, value, value)
		}
	case cdk.IntProperty:
		if _, ok := value.(int); !ok {
			return fmt.Errorf("%v value is not of int type: %v (%T)", p.name, value, value)
		}
	case cdk.FloatProperty:
		if _, ok := value.(float64); !ok {
			return fmt.Errorf("%v value is not of float64 type: %v (%T)", p.name, value, value)
		}
	case cdk.ColorProperty:
		if _, ok := value.(paint.Color); !ok {
			return fmt.Errorf("%v value is not of cdk.Color type: %v (%T)", p.name, value, value)
		}
	case cdk.ThemeProperty:
		if _, ok := value.(paint.Theme); !ok {
			return fmt.Errorf("%v value is not of cdk.Theme type: %v (%T)", p.name, value, value)
		}
	case cdk.PointProperty:
		if _, ok := value.(ptypes.Point2I); !ok {
			return fmt.Errorf("%v value is not of cdk.Point2I type: %v (%T)", p.name, value, value)
		}
	case cdk.RectangleProperty:
		if _, ok := value.(ptypes.Rectangle); !ok {
			return fmt.Errorf("%v value is not of cdk.Rectangle type: %v (%T)", p.name, value, value)
		}
	case cdk.RegionProperty:
		if _, ok := value.(ptypes.Region); !ok {
			return fmt.Errorf("%v value is not of cdk.Region type: %v (%T)", p.name, value, value)
		}
	case cdk.StructProperty:
		// no checks, just pass
	default:
		return fmt.Errorf("invalid property type for %v: %v", p.name, p.Type())
	}
	p.value = value
	return nil
}

func (p *CStyleProperty) SetFromString(value string) error {
	switch p.Type() {
	case cdk.BoolProperty:
		switch strings.ToLower(value) {
		case "true", "t", "1":
			return p.Set(true)
		}
		return p.Set(false)
	case cdk.StringProperty:
		return p.Set(value)
	case cdk.IntProperty:
		if index := strings.Index(value, "px"); index > -1 {
			value = value[:index-1]
		}
		if index := strings.Index(value, "%"); index > -1 {
			value = value[:index-1]
		}
		if v, err := strconv.Atoi(value); err != nil {
			return err
		} else {
			return p.Set(v)
		}
	case cdk.FloatProperty:
		if v, err := strconv.ParseFloat(value, 64); err != nil {
			return err
		} else {
			return p.Set(v)
		}
	case cdk.ColorProperty:
		if c, ok := paint.ParseColor(value); ok {
			return p.Set(c)
		} else {
			return fmt.Errorf("invalid color value: %v", value)
		}
	case cdk.StyleProperty:
		if c, err := paint.ParseStyle(value); err != nil {
			return err
		} else {
			return p.Set(c)
		}
	case cdk.ThemeProperty:
		return fmt.Errorf("theme property not supported by builder features")
	case cdk.PointProperty:
		if v, ok := ptypes.ParsePoint2I(value); ok {
			return p.Set(v)
		} else {
			return fmt.Errorf("invalid point value: %v", value)
		}
	case cdk.RectangleProperty:
		if v, ok := ptypes.ParseRectangle(value); ok {
			return p.Set(v)
		} else {
			return fmt.Errorf("invalid rectangle value: %v", value)
		}
	case cdk.RegionProperty:
		if v, ok := ptypes.ParseRegion(value); ok {
			return p.Set(v)
		} else {
			return fmt.Errorf("invalid region value: %v", value)
		}
	case cdk.StructProperty:
		if efs, ok := p.Default().(enums.EnumFromString); ok {
			if nv, err := efs.FromString(value); err != nil {
				return err
			} else {
				return p.Set(nv)
			}
		}
		return fmt.Errorf("complex property %v not supported by builder features", p.Name())
	}
	return fmt.Errorf("error")
}

func (p *CStyleProperty) Default() (def interface{}) {
	def = p.def
	return
}

func (p *CStyleProperty) Value() (value interface{}) {
	if p.value == nil {
		value = p.def
	} else {
		value = p.value
	}
	return
}
