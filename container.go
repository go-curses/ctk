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
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/go-curses/cdk"
	cenums "github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/cdk/lib/ptypes"
	"github.com/go-curses/ctk/lib/enums"
)

const TypeContainer cdk.CTypeTag = "ctk-container"

// TODO: remove Container.properties, use Widget children directly

func init() {
	_ = cdk.TypesManager.AddType(TypeContainer, nil)
}

// Container Hierarchy:
//	Object
//	  +- Widget
//	    +- Container
//	      +- Bin
//	      +- Box
//	      +- CList
//	      +- Fixed
//	      +- Paned
//	      +- IconView
//	      +- Layout
//	      +- List
//	      +- MenuShell
//	      +- Notebook
//	      +- Socket
//	      +- Table
//	      +- TextView
//	      +- Toolbar
//	      +- ToolItemGroup
//	      +- ToolPalette
//	      +- Tree
//	      +- TreeView
//
// In the Curses Tool Kit, the Container interface is an extension of the CTK
// Widget interface and for all intents and purposes, this is the base class for
// any CTK type that will contain other widgets. The Container also supports the
// tracking of focus and default widgets by maintaining two chain-list types:
// FocusChain and DefaultChain.
//
// Note that currently CTK only supports the FocusChain
type Container interface {
	Widget

	ShowAll()
	Add(w Widget)
	AddWithProperties(widget Widget, argv ...interface{})
	Remove(w Widget)
	ResizeChildren()
	ChildType() (value cdk.CTypeTag)
	GetChildren() (children []Widget)
	HasChild(widget Widget) (present bool)
	GetFocusChild() (value Widget)
	SetFocusChild(child Widget)
	GetFocusVAdjustment() (value Adjustment)
	SetFocusVAdjustment(adjustment Adjustment)
	GetFocusHAdjustment() (value Adjustment)
	SetFocusHAdjustment(adjustment Adjustment)
	ChildGet(child Widget, properties ...cdk.Property) (values []interface{})
	ChildSet(child Widget, argv ...interface{})
	GetChildProperty(child Widget, propertyName cdk.Property) (value interface{})
	SetChildProperty(child Widget, propertyName cdk.Property, value interface{})
	GetBorderWidth() (value int)
	SetBorderWidth(borderWidth int)
	GetFocusChain() (focusableWidgets []Widget, explicitlySet bool)
	SetFocusChain(focusableWidgets []Widget)
	UnsetFocusChain()
	FindChildProperty(property cdk.Property) (value *cdk.CProperty)
	InstallChildProperty(name cdk.Property, kind cdk.PropertyType, write bool, def interface{}) error
	ListChildProperties() (properties []*cdk.CProperty)
	FindAllWidgetsAt(p *ptypes.Point2I) (found []Widget)
	FindWidgetAt(p *ptypes.Point2I) (found Widget)
}

var _ Container = (*CContainer)(nil)

// The CContainer structure implements the Container interface and is exported
// to facilitate type embedding with custom implementations. No member variables
// are exported as the interface methods are the only intended means of
// interacting with Container objects.
type CContainer struct {
	CWidget

	children      []Widget
	resizeMode    cenums.ResizeMode
	properties    []*cdk.CProperty
	property      map[uuid.UUID][]*cdk.CProperty
	focusChain    []Widget
	focusChainSet bool
}

// NewContainer is the constructor for new Container instances.
func NewContainer() Container {
	a := new(CContainer)
	a.Init()
	return a
}

// Init initializes a Container object. This must be called at least once to
// set up the necessary defaults and allocate any memory structures. Calling
// this more than once is safe though unnecessary. Only the first call will
// result in any effect upon the Container instance. Init is used in the
// NewContainer constructor and only necessary when implementing a derivative
// Container type.
func (c *CContainer) Init() (already bool) {
	if c.InitTypeItem(TypeContainer, c) {
		return true
	}
	c.CWidget.Init()
	c.flags = enums.NULL_WIDGET_FLAG
	c.SetFlags(enums.PARENT_SENSITIVE | enums.APP_PAINTABLE)
	_ = c.InstallProperty(PropertyBorderWidth, cdk.IntProperty, true, 0)
	_ = c.InstallProperty(PropertyChild, cdk.StructProperty, true, nil)
	_ = c.InstallProperty(PropertyResizeMode, cdk.StructProperty, true, nil)
	c.children = make([]Widget, 0)
	c.properties = make([]*cdk.CProperty, 0)
	c.property = make(map[uuid.UUID][]*cdk.CProperty)
	c.focusChain = make([]Widget, 0)
	return false
}

// Build provides customizations to the Buildable system for Container Widgets.
func (c *CContainer) Build(builder Builder, element *CBuilderElement) error {
	c.Freeze()
	defer c.Thaw()
	if err := c.CObject.Build(builder, element); err != nil {
		return err
	}
	for _, child := range element.Children {
		if newChild := builder.Build(child); newChild != nil {
			child.Instance = newChild
			if newChildWidget, ok := newChild.(Widget); ok {
				newChildWidget.Show()
				c.Add(newChildWidget)
			} else {
				c.LogError("new child object is not a Widget type: %v (%T)")
			}
		}
	}
	return nil
}

// SetWindow sets the Container window field to the given Window and then does
// the same for each of the Widget children.
//
// Locking: write
func (c *CContainer) SetWindow(w Window) {
	c.CWidget.SetWindow(w)
	children := c.GetChildren()
	for _, child := range children {
		if wc, ok := child.Self().(Container); ok {
			wc.SetWindow(w)
		} else {
			child.SetWindow(w)
		}
	}
}

// ShowAll is a convenience method to call Show on the Container itself and then
// call ShowAll for all Widget children.
//
// Locking: write
func (c *CContainer) ShowAll() {
	c.Show()
	children := c.GetChildren()
	for _, child := range children {
		if cc, ok := child.Self().(Container); ok {
			cc.ShowAll()
		} else {
			child.Show()
		}
	}
}

func (c *CContainer) Map() {
	c.CWidget.Map()
	children := c.GetChildren()
	for _, child := range children {
		if cc, ok := child.Self().(Container); ok {
			cc.Map()
		} else {
			child.Map()
		}
	}
}

func (c *CContainer) Unmap() {
	c.CWidget.Unmap()
	children := c.GetChildren()
	for _, child := range children {
		if cc, ok := child.Self().(Container); ok {
			cc.Unmap()
		} else {
			child.Unmap()
		}
	}
}

// Add the given Widget to the container. Typically this is used for simple
// containers such as Window, Frame, or Button; for more complicated layout
// containers such as Box or Table, this function will pick default packing
// parameters that may not be correct. So consider functions such as
// Box.PackStart() and Table.Attach() as an alternative to Container.Add() in
// those cases. A Widget may be added to only one container at a time; you can't
// place the same widget inside two different containers. This method emits an
// add signal initially and if the listeners return cenums.EVENT_PASS then the
// change is applied.
//
// Parameters:
// 	widget	a widget to be placed inside container
//
// Locking: write
func (c *CContainer) Add(w Widget) {
	// TODO: if can default and no default yet, set
	if f := c.Emit(SignalAdd, c, w); f == cenums.EVENT_PASS {
		w.SetParent(c)
		if window := c.GetWindow(); window != nil {
			if wc, ok := w.Self().(Container); ok {
				wc.SetWindow(window)
			} else {
				w.SetWindow(window)
			}
		}
		w.Map()
		w.Connect(SignalLostFocus, ContainerLostFocusHandle, c.childLostFocus)
		w.Connect(SignalGainedFocus, ContainerGainedFocusHandle, c.childGainedFocus)
		w.Connect(SignalShow, ContainerChildShowHandle, c.childShow)
		w.Connect(SignalHide, ContainerChildHideHandle, c.childHide)
		c.Lock()
		c.children = append(c.children, w)
		childProps := make([]*cdk.CProperty, len(c.properties))
		for idx, prop := range c.properties {
			childProps[idx] = prop.Clone()
		}
		c.property[w.ObjectID()] = childProps
		c.Unlock()
		// log.DebugDF(1, "child added to container: %v", w.ObjectName())
	}
}

// AddWithProperties the given Widget to the Container, setting any given child
// properties at the same time.
// See: Add() and ChildSet()
//
// Parameters:
// 	widget	instance to be placed inside container
// 	argv    list of property names and values
//
// Locking: write
func (c *CContainer) AddWithProperties(widget Widget, argv ...interface{}) {
	c.Add(widget)
	c.ChildSet(widget, argv...)
}

// Remove the given Widget from the Container. Widget must be inside Container.
// This method emits a remove signal initially and if the listeners return
// cenums.EVENT_PASS, the change is applied.
//
// Parameters:
// 	widget	a current child of container
//
// Locking: write
func (c *CContainer) Remove(w Widget) {
	var children []Widget
	for _, child := range c.GetChildren() {
		if child.ObjectID() == w.ObjectID() {
			if f := c.Emit(SignalRemove, c, child); f == cenums.EVENT_PASS {
				_ = w.Disconnect(SignalLostFocus, ContainerLostFocusHandle)
				_ = w.Disconnect(SignalGainedFocus, ContainerGainedFocusHandle)
				_ = w.Disconnect(SignalShow, ContainerChildShowHandle)
				_ = w.Disconnect(SignalHide, ContainerChildHideHandle)
				w.Unmap()
				w.SetParent(nil)
				c.Lock()
				delete(c.property, w.ObjectID())
				c.Unlock()
				continue
			}
		}
		children = append(children, child)
	}
	c.Lock()
	c.children = children
	c.Unlock()
}

// ResizeChildren will call Resize on each child Widget.
//
// Locking: write
func (c *CContainer) ResizeChildren() {
	children := c.GetChildren()
	c.Lock()
	defer c.Unlock()
	for _, child := range children {
		child.Resize()
	}
}

// ChildType returns the type of the children supported by the container. Note
// that this may return TYPE_NONE to indicate that no more children can be
// added, e.g. for a Paned which already has two children.
//
// Returns:
//  tag	a cdk.CTypeTag
//
// Note that usage of this within CTK is unimplemented at this time
//
// Locking: read
func (c *CContainer) ChildType() (value cdk.CTypeTag) {
	c.RLock()
	defer c.RUnlock()
	return TypeWidget
}

// GetChildren returns the container's non-internal children.
//
// Returns:
//  children	list of Widget children
//
// Note that usage of this within CTK is unimplemented at this time
//
// Locking: read
func (c *CContainer) GetChildren() (children []Widget) {
	c.RLock()
	defer c.RUnlock()
	for _, child := range c.children {
		children = append(children, child)
	}
	return
}

func (c *CContainer) HasChild(widget Widget) (present bool) {
	allChildren := append([]Widget{}, c.GetCompositeChildren()...)
	c.RLock()
	defer c.RUnlock()
	allChildren = append(allChildren, c.children...)
	wid := widget.ObjectID()
	for _, child := range allChildren {
		if child.ObjectID() == wid {
			return true
		} else if container, ok := child.Self().(Container); ok && container.HasChild(widget) {
			return true
		}
	}
	return
}

// GetFocusChild returns the current focus child widget inside container. This
// is not the currently focused widget. That can be obtained by calling
// Window.GetFocus().
//
// Returns:
//  widget	child which will receive the focus inside container when the container is focussed, or NULL if none is set
//
// Note that usage of this within CTK is unimplemented at this time
func (c *CContainer) GetFocusChild() (value Widget) {
	return
}

// SetFocusChild updates the focus child for the Container.
//
// Parameters:
// 	child	a Widget, or `nil`
//
// Note that usage of this within CTK is unimplemented at this time
func (c *CContainer) SetFocusChild(child Widget) {}

// GetFocusVAdjustment retrieves the vertical focus adjustment for the
// container.
// See: SetFocusVAdjustment()
//
// Returns:
//  adjustment	the vertical focus adjustment, or NULL if none has been set
//
// Note that usage of this within CTK is unimplemented at this time
func (c *CContainer) GetFocusVAdjustment() (value Adjustment) {
	return nil
}

// SetFocusVAdjustment hooks up an adjustment to focus handling in a container,
// so when a child of the container is focused, the adjustment is scrolled to
// show that widget. This function sets the vertical alignment. See
// ScrolledWindow.GetVAdjustment for a typical way of obtaining the
// adjustment and SetFocusHAdjustment for setting the horizontal adjustment.
//
// Parameters:
// 	adjustment	an adjustment which should be adjusted when the focus is moved among the descendents of container
//
// Note that usage of this within CTK is unimplemented at this time
func (c *CContainer) SetFocusVAdjustment(adjustment Adjustment) {}

// GetFocusHAdjustment retrieves the horizontal focus adjustment for the
// container.
// See: SetFocusHAdjustment()
//
// Returns:
//  adjustment	the horizontal focus adjustment, or NULL if none has been set
//
// Note that usage of this within CTK is unimplemented at this time
func (c *CContainer) GetFocusHAdjustment() (value Adjustment) {
	return nil
}

// SetFocusHAdjustment hooks up an adjustment to focus handling in a container,
// so when a child of the container is focused, the adjustment is scrolled to
// show that widget. This function sets the horizontal alignment. See
// ScrolledWindow.GetHadjustment for a typical way of obtaining the
// adjustment and SetFocusVadjustment for setting the vertical adjustment.
//
// Parameters:
// 	adjustment	an adjustment which should be adjusted when the focus is moved among the descendents of container
//
// Note that usage of this within CTK is unimplemented at this time
func (c *CContainer) SetFocusHAdjustment(adjustment Adjustment) {}

// ChildGet returns the values of one or more child properties for the given
// child.
//
// Parameters:
// 	child          a widget which is a child of container
// 	properties...  one or more property names
//
// Returns:
//   an array of property values, in the order of the property names given, and
//   if the property named is not found, a Go error is returned for that
//   position of property names given
//
// Locking: read
func (c *CContainer) ChildGet(child Widget, properties ...cdk.Property) (values []interface{}) {
	c.RLock()
	defer c.RUnlock()
	for _, nextProperty := range properties {
		if nextProp := c.GetChildProperty(child, nextProperty); nextProp != nil {
			values = append(values, nextProp)
		} else {
			values = append(values, fmt.Errorf("property not found: %v", nextProperty))
		}
	}
	return
}

// ChildSet updates one or more child properties for the given child in the
// container.
//
// Parameters:
// 	child	a widget which is a child of container
//  argv    a list of property name and value pairs
//
// Locking: write
func (c *CContainer) ChildSet(child Widget, argv ...interface{}) {
	argc := len(argv)
	if argc%2 != 0 {
		c.LogError("argument list is not even")
		return
	}
	for i := 0; i < argc; i += 2 {
		if pn, ok := argv[i].(cdk.Property); ok {
			c.SetChildProperty(child, pn, argv[i+1])
		} else {
			c.LogError("expected cdk.Property, received: %v (%T)", argv[i], argv[i])
		}
	}
}

// GetChildProperty returns the value of a child property for the given child.
// Parameters:
// 	child	a widget which is a child of container
// 	propertyName	the name of the property to get
//
// Returns:
//   the value stored in the given property, or nil if the property
//
// Locking: read
func (c *CContainer) GetChildProperty(child Widget, propertyName cdk.Property) (value interface{}) {
	c.RLock()
	defer c.RUnlock()
	if properties, ok := c.property[child.ObjectID()]; ok {
		for _, cp := range properties {
			if cp.Name().String() == propertyName.String() {
				value = cp.Value()
				break
			}
		}
	}
	return
}

// SetChildProperty updates a child property for the given child.
//
// Parameters:
// 	child	a widget which is a child of container
// 	propertyName	the name of the property to set
// 	value	the value to set the property to
//
// Locking: write
func (c *CContainer) SetChildProperty(child Widget, propertyName cdk.Property, value interface{}) {
	c.Lock()
	defer c.Unlock()
	if properties, ok := c.property[child.ObjectID()]; ok {
		for _, cp := range properties {
			if cp.Name().String() == propertyName.String() {
				if err := cp.Set(value); err != nil {
					c.LogErr(err)
				}
				break
			}
		}
	}
}

// GetBorderWidth retrieves the border width of the Container.
// See: SetBorderWidth()
//
// Returns:
// 	the current border width
//
// Locking: read
//
// Note that usage of this within CTK is unimplemented at this time
func (c *CContainer) GetBorderWidth() (value int) {
	c.RLock()
	defer c.RUnlock()
	var err error
	if value, err = c.GetIntProperty(PropertyBorderWidth); err != nil {
		c.LogErr(err)
	}
	return
}

// SetBorderWidth updates the border width of the Container. The border width of
// a container is the amount of space to leave around the outside of the
// container. The only exception to this is Window; because toplevel windows
// can't leave space outside, they leave the space inside. The border is added
// on all sides of the container. To add space to only one side, one approach is
// to create a Alignment widget, call Widget.SetSizeRequest to give it a size
// and place it on the side of the container as a spacer.
//
// Parameters:
// 	borderWidth	amount of blank space to leave outside the container
//
// Locking: write
//
// Note that usage of this within CTK is unimplemented at this time
func (c *CContainer) SetBorderWidth(borderWidth int) {
	c.Lock()
	defer c.Unlock()
	if err := c.SetIntProperty(PropertyBorderWidth, borderWidth); err != nil {
		c.LogErr(err)
	}
}

// GetFocusChain retrieves the focus chain of the container, if one has been set
// explicitly. If no focus chain has been explicitly set, CTK computes the
// focus chain based on the positions of the children.
//
// Returns:
// 	focusableWidgets	widgets in the focus chain.
// 	explicitlySet       TRUE if the focus chain has been set explicitly.
//
// Locking: read
func (c *CContainer) GetFocusChain() (focusableWidgets []Widget, explicitlySet bool) {
	if c.focusChainSet {
		return c.focusChain, true
	}
	var allChildren []Widget
	if c.CanFocus() && c.IsVisible() && c.IsSensitive() {
		allChildren = append(allChildren, c)
	}
	allChildren = append(allChildren, c.GetCompositeChildren()...)
	allChildren = append(allChildren, c.children...)
	c.Lock()
	defer c.Unlock()
	for _, child := range allChildren {
		if cc, ok := child.Self().(Container); ok {
			// the container itself may be more than a Container, if so, add it
			if cc.CanFocus() && cc.IsVisible() && cc.IsSensitive() {
				focusableWidgets = append(focusableWidgets, child)
				continue
			}
			fc, _ := cc.GetFocusChain()
			for _, cChild := range fc {
				if cChild.CanFocus() && cChild.IsVisible() && cChild.IsSensitive() {
					focusableWidgets = append(focusableWidgets, cChild)
				}
			}
		} else {
			if child.CanFocus() && child.IsVisible() && child.IsSensitive() {
				focusableWidgets = append(focusableWidgets, child)
			}
		}
	}
	return
}

// SetFocusChain updates a focus chain, overriding the one computed
// automatically by CTK. In principle each widget in the chain should be a
// descendant of the container, but this is not enforced by this method, since
// it's allowed to set the focus chain before you pack the widgets, or have a
// widget in the chain that isn't always packed. The necessary checks are done
// when the focus chain is actually traversed.
//
// Parameters:
// 	focusableWidgets	the new focus chain.
//
// Locking: write
func (c *CContainer) SetFocusChain(focusableWidgets []Widget) {
	c.Lock()
	defer c.Unlock()
	c.focusChain = focusableWidgets
	c.focusChainSet = true
}

// UnsetFocusChain removes a focus chain explicitly set with SetFocusChain.
//
// Locking: write
func (c *CContainer) UnsetFocusChain() {
	c.Lock()
	defer c.Unlock()
	c.focusChain = []Widget{}
	c.focusChainSet = false
}

// FindChildProperty searches for a child property of a container by name.
//
// Parameters:
// 	 property		the name of the child property to find
//
// Returns:
//  value  the cdk.CProperty of the child property or nil if there is no child property with that name.
//
// Locking: read
func (c *CContainer) FindChildProperty(property cdk.Property) (value *cdk.CProperty) {
	c.RLock()
	defer c.RUnlock()
	for _, prop := range c.properties {
		if prop.Name().String() == property.String() {
			value = prop
			break
		}
	}
	return
}

// InstallChildProperty adds a child property on a container.
//
// Locking: write
func (c *CContainer) InstallChildProperty(name cdk.Property, kind cdk.PropertyType, write bool, def interface{}) error {
	existing := c.FindChildProperty(name)
	if existing != nil {
		return fmt.Errorf("property exists: %v", name)
	}
	c.Lock()
	defer c.Unlock()
	c.properties = append(
		c.properties,
		cdk.NewProperty(name, kind, write, false, def),
	)
	return nil
}

// ListChildProperties returns all child properties of a Container.
//
// Parameters:
//  properties	list of *cdk.Property instances
//
// Locking: read
func (c *CContainer) ListChildProperties() (properties []*cdk.CProperty) {
	c.RLock()
	defer c.RUnlock()
	for _, prop := range c.properties {
		properties = append(properties, prop)
	}
	return
}

// GetWidgetAt is a wrapper around the Widget.GetWidgetAt() method that if the
// Container has the given Point2I within it's bounds, will return itself, or
// if there is a child Widget at the given Point2I will return the child Widget.
// If the Container does not have the given point within it's bounds, will
// return nil
//
// Locking: read
func (c *CContainer) GetWidgetAt(p *ptypes.Point2I) Widget {
	if c.HasPoint(p) && c.IsVisible() {
		allChildren := append([]Widget{}, c.GetCompositeChildren()...)
		c.RLock()
		allChildren = append(allChildren, c.children...)
		c.RUnlock()
		for _, child := range allChildren {
			switch childType := child.(type) {
			case Button:
				if childType.HasPoint(p) && childType.IsVisible() {
					return childType
				}
			case Alignment, Frame, ButtonBox, Scrollbar, ScrolledViewport, Viewport, Window, Box, Container:
				if childType.HasPoint(p) && childType.IsVisible() {
					if w := childType.GetWidgetAt(p); w != nil {
						return w
					}
				}
			default:
				if child.HasPoint(p) && child.IsVisible() {
					return child
				}
			}
		}
		return c
	}
	return nil
}

func (c *CContainer) FindAllWidgetsAt(p *ptypes.Point2I) (found []Widget) {
	track := new(WidgetSlice)
	if first := c.FindWidgetAt(p); first != nil {
		track.Append(first)
		parent := first.GetParent()
		for parent != nil {
			if idx := track.IndexOf(parent); idx < 0 {
				track.Append(parent)
			}
			next := parent.GetParent()
			if next == nil || track.IndexOf(next) > -1 {
				break
			}
			parent = next
		}
		found = *track
	}
	return
}

func (c *CContainer) FindWidgetAt(p *ptypes.Point2I) (found Widget) {
	return c.findWidgetAt(c, p)
}

func (c *CContainer) findWidgetAt(parent Widget, p *ptypes.Point2I) (found Widget) {
	if parent.HasPoint(p) && parent.IsVisible() {
		if cParent, ok := parent.Self().(Container); ok {
			allChildren := append([]Widget{}, cParent.GetCompositeChildren()...)
			allChildren = append(allChildren, cParent.GetChildren()...)
			for _, child := range allChildren {
				switch childType := child.(type) {
				case Button:
					if childType.HasPoint(p) && childType.IsVisible() {
						return childType
					}
				case Alignment, Frame, ButtonBox, ScrolledViewport, Viewport, Window, Box, Container:
					if grandchild := c.findWidgetAt(childType, p); grandchild != nil && grandchild.IsVisible() {
						return grandchild
					}
				default:
					if child.HasPoint(p) && child.IsVisible() {
						return child
					}
				}
			}
		}
		found = parent
	}
	return
}

func (c *CContainer) RenderFreeze() {
	c.CWidget.RenderFreeze()
	allChildren := append([]Widget{}, c.GetCompositeChildren()...)
	allChildren = append(allChildren, c.GetChildren()...)
	for _, child := range allChildren {
		switch childType := child.(type) {
		case Container:
			childType.RenderFreeze()
		default:
			child.RenderFreeze()
		}
	}
}

func (c *CContainer) RenderThaw() {
	c.CWidget.ResumeSignal(SignalInvalidate)
	c.renderThaw(c)
	c.Resize()
	c.CWidget.ResumeSignal(SignalDraw)
}

func (c *CContainer) renderThaw(parent Container) {
	allChildren := append([]Widget{}, parent.GetCompositeChildren()...)
	allChildren = append(allChildren, parent.GetChildren()...)
	for _, child := range allChildren {
		switch childType := child.(type) {
		case Container:
			childType.ResumeSignal(SignalInvalidate, SignalDraw)
			c.renderThaw(childType)
		default:
			child.ResumeSignal(SignalInvalidate, SignalDraw)
		}
	}
}

func (c *CContainer) Destroy() {
	allChildren := append([]Widget{}, c.GetCompositeChildren()...)
	allChildren = append(allChildren, c.GetChildren()...)
	for _, child := range allChildren {
		switch childType := child.(type) {
		case Container:
			childType.Destroy()
		default:
			child.Destroy()
		}
	}
}

func (c *CContainer) InvalidateChildren() {
	var descend func(cc Container)
	descend = func(cc Container) {
		children := cc.GetChildren()
		for _, child := range children {
			// child.SetInvalidated(true)
			if cChild, ok := child.Self().(Container); ok {
				descend(cChild)
			}
			child.Invalidate()
		}
	}
	descend(c)
}

func (c *CContainer) InvalidateAllChildren() {
	var descend func(cc Container)
	descend = func(cc Container) {
		allChildren := append([]Widget{}, cc.GetCompositeChildren()...)
		allChildren = append(allChildren, cc.GetChildren()...)
		for _, child := range allChildren {
			// child.SetInvalidated(true)
			if cChild, ok := child.Self().(Container); ok {
				descend(cChild)
			}
			child.Invalidate()
		}
	}
	descend(c)
}

func (c *CContainer) childShow(data []interface{}, argv ...interface{}) cenums.EventFlag {
	c.Resize()
	return cenums.EVENT_PASS
}

func (c *CContainer) childHide(data []interface{}, argv ...interface{}) cenums.EventFlag {
	c.Resize()
	return cenums.EVENT_PASS
}

func (c *CContainer) childLostFocus(data []interface{}, argv ...interface{}) cenums.EventFlag {
	// c.Invalidate()
	return cenums.EVENT_PASS
}

func (c *CContainer) childGainedFocus(data []interface{}, argv ...interface{}) cenums.EventFlag {
	// c.Invalidate()
	return cenums.EVENT_PASS
}

// The width of the empty border outside the containers children.
// Flags: Read / Write
// Allowed values: <= 65535
// Default value: 0
const PropertyBorderWidth cdk.Property = "border-width"

// Can be used to add a new child to the container.
// Flags: Write
const PropertyChild cdk.Property = "child"

// Specify how resize events are handled.
// Flags: Read / Write
// Default value: GTK_RESIZE_PARENT
const PropertyResizeMode cdk.Property = "resize-mode"

// Listener function arguments:
// 	widget Widget
const SignalAdd cdk.Signal = "add"

const SignalCheckResize cdk.Signal = "check-resize"

// Listener function arguments:
// 	widget Widget
const SignalRemove cdk.Signal = "remove"

// Listener function arguments:
// 	widget Widget
const SignalSetFocusChild cdk.Signal = "set-focus-child"

const SignalPushCompositeChild cdk.Signal = "push-composite-child"

const SignalPopCompositeChild cdk.Signal = "pop-composite-child"

const ContainerChildShowHandle = "container-child-show-handler"

const ContainerChildHideHandle = "container-child-hide-handler"

const ContainerLostFocusHandle = "container-lost-focus-handler"

const ContainerGainedFocusHandle = "container-gained-focus-handler"