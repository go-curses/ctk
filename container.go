package ctk

import (
	"fmt"

	"github.com/go-curses/cdk"
	"github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/cdk/lib/ptypes"
	"github.com/go-curses/cdk/log"
	"github.com/gofrs/uuid"
)

// CDK type-tag for Container objects
const TypeContainer cdk.CTypeTag = "ctk-container"

func init() {
	_ = cdk.TypesManager.AddType(TypeContainer, nil)
}

type Callback = func()

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
// In the Curses Tool Kit, the Container interface is an extension of the CTK
// Widget interface and for all intents and purposes, this is the base class for
// any CTK type that will contain other widgets. The Container also supports the
// tracking of focus and default widgets by maintaining two WidgetChain types:
// FocusChain and DefaultChain.
type Container interface {
	Widget
	Buildable

	Init() (already bool)
	Add(w Widget)
	Remove(w Widget)
	AddWithProperties(widget Widget, argv ...interface{})
	ForEach(callback Callback, callbackData interface{})
	GetChildren() (children []Widget)
	GetFocusChild() (value Widget)
	SetFocusChild(child Widget)
	GetFocusVAdjustment() (value Adjustment)
	SetFocusVAdjustment(adjustment Adjustment)
	GetFocusHAdjustment() (value Adjustment)
	SetFocusHAdjustment(adjustment Adjustment)
	ChildType() (value cdk.CTypeTag)
	ChildGet(child Widget, properties ...cdk.Property) (values []interface{})
	ChildSet(child Widget, argv ...interface{})
	GetChildProperty(child Widget, propertyName cdk.Property) (value interface{})
	SetChildProperty(child Widget, propertyName cdk.Property, value interface{})
	ForAll(callback Callback, callbackData interface{})
	GetBorderWidth() (value int)
	SetBorderWidth(borderWidth int)
	GetFocusChain() (focusableWidgets []interface{}, explicitlySet bool)
	SetFocusChain(focusableWidgets []interface{})
	UnsetFocusChain()
	FindChildProperty(property cdk.Property) (value *cdk.CProperty)
	InstallChildProperty(name cdk.Property, kind cdk.PropertyType, write bool, def interface{}) error
	ListChildProperties() (properties []*cdk.CProperty)
	Build(builder Builder, element *CBuilderElement) error
	SetOrigin(x, y int)
	SetWindow(w Window)
	GetWidgetAt(p *ptypes.Point2I) Widget
	ShowAll()
}

// The CContainer structure implements the Container interface and is
// exported to facilitate type embedding with custom implementations. No member
// variables are exported as the interface methods are the only intended means
// of interacting with Container objects
type CContainer struct {
	CWidget

	children      []Widget
	resizeMode    enums.ResizeMode
	properties    []*cdk.CProperty
	property      map[uuid.UUID][]*cdk.CProperty
	focusChain    []interface{}
	focusChainSet bool
}

// Container object initialization. This must be called at least once to setup
// the necessary defaults and allocate any memory structures. Calling this more
// than once is safe though unnecessary. Only the first call will result in any
// effect upon the Container instance
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
	c.Connect(SignalLostFocus, ContainerLostFocusHandle, c.lostFocus)
	c.Connect(SignalGainedFocus, ContainerGainedFocusHandle, c.gainedFocus)
	return false
}

// Adds widget to container.Typically used for simple containers such as Window,
// Frame, or Button; for more complicated layout containers such as Box or
// Table, this function will pick default packing parameters that may not be
// correct. So consider functions such as Box.PackStart() and Table.Attach()
// as an alternative to Container.Add() in those cases. A Widget may be added to
// only one container at a time; you can't place the same widget inside two
// different containers. This method emits an add signal initially and if the
// listeners return EVENT_PASS then the change is applied
//
// Parameters:
// 	widget	a widget to be placed inside container
//
// Emits: SignalAdd, Argv=[Container instance, Widget instance]
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

// Removes Widget from container. Widget must be inside Container. This method
// emits a remove signal initially and if the listeners return EVENT_PASS then
// the change is applied
//
// Parameters:
// 	widget	a current child of container
//
// Emits: SignalRemove, Argv=[Container instance, Widget instance]
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

// Adds widget to container, setting child properties at the same time. See
// Add and ChildSet for more details.
// Parameters:
// 	widget	a widget to be placed inside container
// 	argv    a list of property names and values
//
func (c *CContainer) AddWithProperties(widget Widget, argv ...interface{}) {
	c.Add(widget)
	c.ChildSet(widget, argv...)
}

// Invokes callback on each non-internal child of container . See
// Forall for details on what constitutes an "internal"
// child. Most applications should use Foreach, rather than
// Forall.
// Parameters:
// 	callback	a callback.
// 	callbackData	callback user data
func (c *CContainer) ForEach(callback Callback, callbackData interface{}) {}

// Returns the container's non-internal children. See Forall
// for details on what constitutes an "internal" child.
// Returns:
// 	a newly-allocated list of the container's non-internal
// 	children.
// 	[element-type Widget][transfer container]
func (c *CContainer) GetChildren() (children []Widget) {
	for _, child := range c.children {
		children = append(children, child)
	}
	return
}

// Returns the current focus child widget inside container . This is not the
// currently focused widget. That can be obtained by calling
// WindowGetFocus.
// Returns:
// 	The child widget which will receive the focus inside container
// 	when the container is focussed, or NULL if none is set.
func (c *CContainer) GetFocusChild() (value Widget) {
	return
}

// Sets, or unsets if child is NULL, the focused child of container . This
// function emits the Container::set_focus_child signal of container .
// Implementations of Container can override the default behaviour by
// overriding the class closure of this signal. This is function is mostly
// meant to be used by widgets. Applications can use WidgetGrabFocus
// to manualy set the focus to a specific widget.
// Parameters:
// 	child	a Widget, or NULL.
func (c *CContainer) SetFocusChild(child Widget) {}

// Retrieves the vertical focus adjustment for the container. See
// SetFocusVAdjustment.
// Returns:
// 	the vertical focus adjustment, or NULL if none has been set.
// 	[transfer none]
func (c *CContainer) GetFocusVAdjustment() (value Adjustment) {
	return nil
}

// Hooks up an adjustment to focus handling in a container, so when a child
// of the container is focused, the adjustment is scrolled to show that
// widget. This function sets the vertical alignment. See
// ScrolledWindowGetVAdjustment for a typical way of obtaining the
// adjustment and SetFocusHAdjustment for setting the
// horizontal adjustment. The adjustments have to be in pixel units and in
// the same coordinate system as the allocation for immediate children of the
// container.
// Parameters:
// 	adjustment	an adjustment which should be adjusted when the focus
// is moved among the descendents of container
//
func (c *CContainer) SetFocusVAdjustment(adjustment Adjustment) {}

// Retrieves the horizontal focus adjustment for the container. See
// SetFocusHAdjustment.
// Returns:
// 	the horizontal focus adjustment, or NULL if none has been set.
// 	[transfer none]
func (c *CContainer) GetFocusHAdjustment() (value Adjustment) {
	return nil
}

// Hooks up an adjustment to focus handling in a container, so when a child
// of the container is focused, the adjustment is scrolled to show that
// widget. This function sets the horizontal alignment. See
// ScrolledWindowGetHadjustment for a typical way of obtaining the
// adjustment and SetFocusVadjustment for setting the
// vertical adjustment. The adjustments have to be in pixel units and in the
// same coordinate system as the allocation for immediate children of the
// container.
// Parameters:
// 	adjustment	an adjustment which should be adjusted when the focus is
// moved among the descendents of container
//
func (c *CContainer) SetFocusHAdjustment(adjustment Adjustment) {}

// func (c *CContainer) ResizeChildren() {}

// Returns the type of the children supported by the container. Note that
// this may return G_TYPE_NONE to indicate that no more children can be
// added, e.g. for a Paned which already has two children.
// Returns:
// 	a GType.
func (c *CContainer) ChildType() (value cdk.CTypeTag) {
	return TypeWidget
}

// Gets the values of one or more child properties for child and container.
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

// Sets one or more child properties for the given child in the container.
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

// Gets the value of a child property for child and container .
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

// Sets a child property for child and container .
// Parameters:
// 	child	a widget which is a child of container
//
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

// Gets the values of one or more child properties for child and container .
// Parameters:
// 	child	a widget which is a child of container
//
// 	firstPropertyName	the name of the first property to get
// 	varArgs	return location for the first property, followed
// optionally by more name/return location pairs, followed by NULL
// func (c *CContainer) ChildGetValist(child Widget, firstPropertyName string, varArgs va_list) {}

// Sets one or more child properties for child and container .
// Parameters:
// 	child	a widget which is a child of container
//
// 	firstPropertyName	the name of the first property to set
// 	varArgs	a NULL-terminated list of property names and values, starting
// with first_prop_name
//
// func (c *CContainer) ChildSetValist(child Widget, firstPropertyName string, varArgs va_list) {}

// Invokes callback on each child of container , including children that are
// considered "internal" (implementation details of the container).
// "Internal" children generally weren't added by the user of the container,
// but were added by the container implementation itself. Most applications
// should use Foreach, rather than Forall.
// Parameters:
// 	callback	a callback
// 	callbackData	callback user data
func (c *CContainer) ForAll(callback Callback, callbackData interface{}) {}

// Retrieves the border width of the container. See
// SetBorderWidth.
// Returns:
// 	the current border width
func (c *CContainer) GetBorderWidth() (value int) {
	var err error
	if value, err = c.GetIntProperty(PropertyBorderWidth); err != nil {
		c.LogErr(err)
	}
	return
}

// Sets the border width of the container. The border width of a container is
// the amount of space to leave around the outside of the container. The only
// exception to this is Window; because toplevel windows can't leave space
// outside, they leave the space inside. The border is added on all sides of
// the container. To add space to only one side, one approach is to create a
// Alignment widget, call WidgetSetSizeRequest to give it a size,
// and place it on the side of the container as a spacer.
// Parameters:
// 	borderWidth	amount of blank space to leave outside
// the container. Valid values are in the range 0-65535 pixels.
func (c *CContainer) SetBorderWidth(borderWidth int) {
	if err := c.SetIntProperty(PropertyBorderWidth, borderWidth); err != nil {
		c.LogErr(err)
	}
}

// Retrieves the focus chain of the container, if one has been set
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
			if cc.CanFocus() && cc.IsVisible() {
				focusableWidgets = append(focusableWidgets, child)
				continue
			}
			fc, _ := cc.GetFocusChain()
			for _, cChild := range fc {
				focusableWidgets = append(focusableWidgets, cChild)
			}
		} else if child.CanFocus() && child.IsVisible() {
			focusableWidgets = append(focusableWidgets, child)
		}
	}
	return
}

// Sets a focus chain, overriding the one computed automatically by CTK. In
// principle each widget in the chain should be a descendant of the
// container, but this is not enforced by this method, since it's allowed to
// set the focus chain before you pack the widgets, or have a widget in the
// chain that isn't always packed. The necessary checks are done when the
// focus chain is actually traversed.
// Parameters:
// 	focusableWidgets	the new focus chain.
func (c *CContainer) SetFocusChain(focusableWidgets []interface{}) {
	c.focusChain = focusableWidgets
	c.focusChainSet = true
}

// Removes a focus chain explicitly set with SetFocusChain.
func (c *CContainer) UnsetFocusChain() {
	c.focusChain = []interface{}{}
	c.focusChainSet = false
}

// Finds a child property of a container by name.
// Parameters:
// 	 property		the name of the child property to find
//
// Returns:
//   the cdk.CProperty of the child property or nil if there is no child
// 	 property with that name.
func (c *CContainer) FindChildProperty(property cdk.Property) (value *cdk.CProperty) {
	for _, prop := range c.properties {
		if prop.Name().String() == property.String() {
			value = prop
			break
		}
	}
	return
}

// Installs a child property on a container.
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

// Returns all child properties of a container class.
// Parameters:
// 	cclass	a ContainerClass.
// 	nProperties	location to return the number of child properties found
// 	returns	a newly
// allocated NULL-terminated array of GParamSpec*.
// The array must be freed with g_free.
func (c *CContainer) ListChildProperties() (properties []*cdk.CProperty) {
	for _, prop := range c.properties {
		properties = append(properties, prop)
	}
	return
}

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

func (c *CContainer) SetOrigin(x, y int) {
	if f := c.Emit(SignalOrigin, c, ptypes.MakePoint2I(x, y)); f == enums.EVENT_PASS {
		c.origin.Set(x, y)
		for _, child := range c.GetChildren() {
			child.SetOrigin(x, y)
		}
	}
}

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

// A wrapper around the Widget.GetWidgetAt() method that if the Container has
// the given Point2I within it's bounds, will return a itself, or if there is a
// child Widget at the given Point2I will return the child Widget. If the
// Container does not have the given point within it's bounds, will return nil
func (c *CContainer) GetWidgetAt(p *ptypes.Point2I) Widget {
	if c.HasPoint(p) && c.IsVisible() {
		for _, child := range c.children {
			switch c := child.(type) {
			case Container:
				if w := c.GetWidgetAt(p); w != nil && w.IsVisible() {
					return w
				}
			default:
				if child.HasPoint(p) {
					return child
				}
			}
		}
		return c
	}
	return nil
}

// The Container type implements a version of Widget.ShowAll() where all the
// children of the Container have their ShowAll() method called, in addition to
// calling Show() on itself first.
func (c *CContainer) ShowAll() {
	c.Show()
	for _, child := range c.GetChildren() {
		child.ShowAll()
	}
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
