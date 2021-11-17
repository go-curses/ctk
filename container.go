package ctk

import (
	"fmt"

	"github.com/go-curses/cdk"
	"github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/cdk/lib/ptypes"
	"github.com/go-curses/cdk/log"
	"github.com/gofrs/uuid"
)

const TypeContainer cdk.CTypeTag = "ctk-container"

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
	Buildable

	Init() (already bool)
	Build(builder Builder, element *CBuilderElement) error
	SetOrigin(x, y int)
	SetWindow(w Window)
	ShowAll()
	Add(w Widget)
	AddWithProperties(widget Widget, argv ...interface{})
	Remove(w Widget)
	ResizeChildren()
	ChildType() (value cdk.CTypeTag)
	GetChildren() (children []Widget)
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
	GetFocusChain() (focusableWidgets []interface{}, explicitlySet bool)
	SetFocusChain(focusableWidgets []interface{})
	UnsetFocusChain()
	FindChildProperty(property cdk.Property) (value *cdk.CProperty)
	InstallChildProperty(name cdk.Property, kind cdk.PropertyType, write bool, def interface{}) error
	ListChildProperties() (properties []*cdk.CProperty)
	GetWidgetAt(p *ptypes.Point2I) Widget
}

// The CContainer structure implements the Container interface and is exported
// to facilitate type embedding with custom implementations. No member variables
// are exported as the interface methods are the only intended means of
// interacting with Container objects.
type CContainer struct {
	CWidget

	children      []Widget
	resizeMode    enums.ResizeMode
	properties    []*cdk.CProperty
	property      map[uuid.UUID][]*cdk.CProperty
	focusChain    []interface{}
	focusChainSet bool
}

// MakeContainer is used by the Buildable system to construct a new Container.
func MakeContainer() *CContainer {
	return NewContainer()
}

// NewContainer is the constructor for new Container instances.
func NewContainer() *CContainer {
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
	c.flags = NULL_WIDGET_FLAG
	c.SetFlags(PARENT_SENSITIVE | APP_PAINTABLE)
	_ = c.InstallProperty(PropertyBorderWidth, cdk.IntProperty, true, 0)
	_ = c.InstallProperty(PropertyChild, cdk.StructProperty, true, nil)
	_ = c.InstallProperty(PropertyResizeMode, cdk.StructProperty, true, nil)
	c.children = make([]Widget, 0)
	c.properties = make([]*cdk.CProperty, 0)
	c.property = make(map[uuid.UUID][]*cdk.CProperty)
	c.focusChain = make([]interface{}, 0)
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

// SetOrigin emits an origin signal and if all signal handlers return
// enums.EVENT_PASS, updates the Container origin field while also setting all
// children to the same origin.
//
// Note that this method's behaviour may change.
func (c *CContainer) SetOrigin(x, y int) {
	if f := c.Emit(SignalOrigin, c, ptypes.MakePoint2I(x, y)); f == enums.EVENT_PASS {
		c.origin.Set(x, y)
		for _, child := range c.GetChildren() {
			child.SetOrigin(x, y)
		}
	}
	c.Invalidate()
}

// SetWindow sets the Container window field to the given Window and then does
// the same for each of the Widget children.
func (c *CContainer) SetWindow(w Window) {
	c.CWidget.SetWindow(w)
	for _, child := range c.GetChildren() {
		if wc, ok := child.(Container); ok {
			wc.SetWindow(w)
		} else {
			child.SetWindow(w)
		}
	}
}

// ShowAll is a convenience method to call Show on the Container itself and then
// call ShowAll for all Widget children.
func (c *CContainer) ShowAll() {
	c.Show()
	for _, child := range c.GetChildren() {
		child.ShowAll()
	}
}

// Add the given Widget to the container. Typically this is used for simple
// containers such as Window, Frame, or Button; for more complicated layout
// containers such as Box or Table, this function will pick default packing
// parameters that may not be correct. So consider functions such as
// Box.PackStart() and Table.Attach() as an alternative to Container.Add() in
// those cases. A Widget may be added to only one container at a time; you can't
// place the same widget inside two different containers. This method emits an
// add signal initially and if the listeners return enums.EVENT_PASS then the
// change is applied.
//
// Parameters:
// 	widget	a widget to be placed inside container
func (c *CContainer) Add(w Widget) {
	// TODO: if can default and no default yet, set
	if f := c.Emit(SignalAdd, c, w); f == enums.EVENT_PASS {
		log.DebugDF(1, "Container.Add(%v)", w)
		w.SetParent(c)
		if wc, ok := w.(Container); ok {
			wc.SetWindow(c.GetWindow())
		} else {
			w.SetWindow(c.GetWindow())
		}
		w.Connect(SignalLostFocus, ContainerLostFocusHandle, c.lostFocus)
		w.Connect(SignalGainedFocus, ContainerGainedFocusHandle, c.gainedFocus)
		c.children = append(c.children, w)
		childProps := make([]*cdk.CProperty, len(c.properties))
		for idx, prop := range c.properties {
			childProps[idx] = prop.Clone()
		}
		c.property[w.ObjectID()] = childProps
		c.Resize()
	}
}

// AddWithProperties the given Widget to the Container, setting any given child
// properties at the same time.
// See: Add() and ChildSet()
//
// Parameters:
// 	widget	widget to be placed inside container
// 	argv    list of property names and values
func (c *CContainer) AddWithProperties(widget Widget, argv ...interface{}) {
	c.Add(widget)
	c.ChildSet(widget, argv...)
}

// Remove the given Widget from the Container. Widget must be inside Container.
// This method emits a remove signal initially and if the listeners return
// enums.EVENT_PASS, the change is applied.
//
// Parameters:
// 	widget	a current child of container
func (c *CContainer) Remove(w Widget) {
	var children []Widget
	resize := false
	for _, child := range c.children {
		if child.ObjectID() == w.ObjectID() {
			if f := c.Emit(SignalRemove, c, child); f == enums.EVENT_PASS {
				_ = w.Disconnect(SignalLostFocus, ContainerLostFocusHandle)
				_ = w.Disconnect(SignalGainedFocus, ContainerGainedFocusHandle)
				w.SetParent(nil)
				delete(c.property, w.ObjectID())
				resize = true
				continue
			}
		}
		children = append(children, child)
	}
	if len(children) == 0 {
		if len(c.children) > 0 {
			c.children = make([]Widget, 0)
		}
	} else {
		c.children = children
	}
	if resize {
		c.Resize()
	}
}

// ResizeChildren will call Resize on each child Widget.
func (c *CContainer) ResizeChildren() {
	for _, child := range c.GetChildren() {
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
func (c *CContainer) ChildType() (value cdk.CTypeTag) {
	return TypeWidget
}

// GetChildren returns the container's non-internal children.
//
// Returns:
//  children	list of Widget children
//
// Note that usage of this within CTK is unimplemented at this time
func (c *CContainer) GetChildren() (children []Widget) {
	for _, child := range c.children {
		children = append(children, child)
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
func (c *CContainer) ChildGet(child Widget, properties ...cdk.Property) (values []interface{}) {
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
func (c *CContainer) GetChildProperty(child Widget, propertyName cdk.Property) (value interface{}) {
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
func (c *CContainer) SetChildProperty(child Widget, propertyName cdk.Property, value interface{}) {
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
// Note that usage of this within CTK is unimplemented at this time
func (c *CContainer) GetBorderWidth() (value int) {
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
// Note that usage of this within CTK is unimplemented at this time
func (c *CContainer) SetBorderWidth(borderWidth int) {
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
func (c *CContainer) GetFocusChain() (focusableWidgets []interface{}, explicitlySet bool) {
	if c.focusChainSet {
		return c.focusChain, true
	}
	for _, child := range c.children {
		if cc, ok := child.(Container); ok {
			if cc.CanFocus() && cc.IsVisible() && cc.IsSensitive() {
				focusableWidgets = append(focusableWidgets, child)
				continue
			}
			fc, _ := cc.GetFocusChain()
			for _, cChild := range fc {
				focusableWidgets = append(focusableWidgets, cChild)
			}
		} else if child.CanFocus() && child.IsVisible() && cc.IsSensitive() {
			focusableWidgets = append(focusableWidgets, child)
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
func (c *CContainer) SetFocusChain(focusableWidgets []interface{}) {
	c.focusChain = focusableWidgets
	c.focusChainSet = true
}

// UnsetFocusChain removes a focus chain explicitly set with SetFocusChain.
func (c *CContainer) UnsetFocusChain() {
	c.focusChain = []interface{}{}
	c.focusChainSet = false
}

// FindChildProperty searches for a child property of a container by name.
//
// Parameters:
// 	 property		the name of the child property to find
//
// Returns:
//  value  the cdk.CProperty of the child property or nil if there is no child property with that name.
func (c *CContainer) FindChildProperty(property cdk.Property) (value *cdk.CProperty) {
	for _, prop := range c.properties {
		if prop.Name().String() == property.String() {
			value = prop
			break
		}
	}
	return
}

// InstallChildProperty adds a child property on a container.
func (c *CContainer) InstallChildProperty(name cdk.Property, kind cdk.PropertyType, write bool, def interface{}) error {
	existing := c.FindChildProperty(name)
	if existing != nil {
		return fmt.Errorf("property exists: %v", name)
	}
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
func (c *CContainer) ListChildProperties() (properties []*cdk.CProperty) {
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
func (c *CContainer) GetWidgetAt(p *ptypes.Point2I) Widget {
	if c.HasPoint(p) && c.IsVisible() {
		for _, child := range c.children {
			switch c := child.(type) {
			case Container:
				if w := c.GetWidgetAt(p); w != nil && w.IsVisible() {
					return w
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

func (c *CContainer) lostFocus(data []interface{}, argv ...interface{}) enums.EventFlag {
	c.Invalidate()
	return enums.EVENT_PASS
}

func (c *CContainer) gainedFocus(data []interface{}, argv ...interface{}) enums.EventFlag {
	c.Invalidate()
	return enums.EVENT_PASS
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

const ContainerLostFocusHandle = "container-lost-focus-handler"

const ContainerGainedFocusHandle = "container-gained-focus-handler"
