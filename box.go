package ctk

import (
	"github.com/go-curses/cdk"
	cenums "github.com/go-curses/cdk/lib/enums"
	cmath "github.com/go-curses/cdk/lib/math"
	"github.com/go-curses/cdk/lib/paint"
	"github.com/go-curses/cdk/lib/ptypes"
	"github.com/go-curses/cdk/memphis"
	"github.com/go-curses/ctk/lib/enums"
)

const TypeBox cdk.CTypeTag = "ctk-box"

func init() {
	_ = cdk.TypesManager.AddType(TypeBox, func() interface{} { return MakeBox() })
}

// Box Hierarchy:
//	Object
//	  +- Widget
//	    +- Container
//	      +- Box
//	        +- ButtonBox
//	        +- VBox
//	        +- HBox
//
// The Box Widget is a Container for organizing one or more child Widgets. A Box
// displays either a horizontal row or vertical column of the visible children
// contained within.
type Box interface {
	Container
	Buildable
	Orientable

	Init() (already bool)
	Build(builder Builder, element *CBuilderElement) error
	GetOrientation() (orientation cenums.Orientation)
	SetOrientation(orientation cenums.Orientation)
	GetHomogeneous() (value bool)
	SetHomogeneous(homogeneous bool)
	GetSpacing() (value int)
	SetSpacing(spacing int)
	Add(child Widget)
	Remove(w Widget)
	PackStart(child Widget, expand, fill bool, padding int)
	PackEnd(child Widget, expand, fill bool, padding int)
	ReorderChild(child Widget, position int)
	QueryChildPacking(child Widget) (expand bool, fill bool, padding int, packType enums.PackType)
	SetChildPacking(child Widget, expand bool, fill bool, padding int, packType enums.PackType)
	GetFocusChain() (focusableWidgets []interface{}, explicitlySet bool)
	GetSizeRequest() (width, height int)
}

// The CBox structure implements the Box interface and is exported to
// facilitate type embedding with custom implementations. No member variables
// are exported as the interface methods are the only intended means of
// interacting with Box objects.
type CBox struct {
	CContainer
}

// The cBoxChild is an internal structure used for tracking the per-child
// packing configuration. This should never need to be accessed by developers
// directly.
type cBoxChild struct {
	widget   Widget
	expand   bool
	fill     bool
	padding  int
	packType enums.PackType
}

// MakeBox is used by the Buildable system to construct a new Box with default
// settings of: horizontal orientation, dynamically sized (not homogeneous) and
// no extra spacing.
func MakeBox() (box Box) {
	box = NewBox(cenums.ORIENTATION_HORIZONTAL, false, 0)
	return
}

// NewBox is the constructor for new Box instances.
//
// Parameters:
//  orientation  the orientation of the Box vertically or horizontally
//  homogeneous  whether each child receives an equal size allocation or not
//  spacing      extra spacing to include between children
func NewBox(orientation cenums.Orientation, homogeneous bool, spacing int) Box {
	b := new(CBox)
	b.Init()
	b.SetOrientation(orientation)
	b.SetHomogeneous(homogeneous)
	b.SetSpacing(spacing)
	return b
}

// Init initializes a Box object. This must be called at least once to
// set up the necessary defaults and allocate any memory structures. Calling
// this more than once is safe though unnecessary. Only the first call will
// result in any effect upon the Box instance. Init is used in the
// NewBox constructor and only necessary when implementing a derivative
// Box type.
func (b *CBox) Init() (already bool) {
	if b.InitTypeItem(TypeBox, b) {
		return true
	}
	b.CContainer.Init()
	b.flags = enums.NULL_WIDGET_FLAG
	b.SetFlags(enums.PARENT_SENSITIVE | enums.APP_PAINTABLE)
	_ = b.InstallBuildableProperty(PropertyDebugChildren, cdk.BoolProperty, true, false)
	_ = b.InstallBuildableProperty(PropertyOrientation, cdk.StructProperty, true, cenums.ORIENTATION_HORIZONTAL)
	_ = b.InstallBuildableProperty(PropertyHomogeneous, cdk.BoolProperty, true, false)
	_ = b.InstallBuildableProperty(PropertySpacing, cdk.IntProperty, true, 0)
	_ = b.InstallChildProperty(PropertyBoxChildPackType, cdk.StructProperty, true, enums.PackStart)
	_ = b.InstallChildProperty(PropertyBoxChildExpand, cdk.BoolProperty, true, false)
	_ = b.InstallChildProperty(PropertyBoxChildFill, cdk.BoolProperty, true, true)
	_ = b.InstallChildProperty(PropertyBoxChildPadding, cdk.IntProperty, true, 0)
	b.Connect(SignalInvalidate, BoxInvalidateHandle, b.invalidate)
	b.Connect(SignalResize, BoxResizeHandle, b.resize)
	b.Connect(SignalDraw, BoxDrawHandle, b.draw)
	return false
}

// Build provides customizations to the Buildable system for Box Widgets.
func (b *CBox) Build(builder Builder, element *CBuilderElement) error {
	b.Freeze()
	defer b.Thaw()
	if err := b.CObject.Build(builder, element); err != nil {
		return err
	}
	for _, child := range element.Children {
		if newChild := builder.Build(child); newChild != nil {
			child.Instance = newChild
			if newChildWidget, ok := newChild.(Widget); ok {
				newChildWidget.Show()
				if len(child.Packing) > 0 {
					expand, fill, padding, packType := builder.ParsePacking(child)
					if packType == enums.PackStart {
						b.PackStart(newChildWidget, expand, fill, padding)
					} else {
						b.PackEnd(newChildWidget, expand, fill, padding)
					}
				} else {
					b.Add(newChildWidget)
				}
				if newChildWidget.HasFlags(enums.HAS_FOCUS) {
					newChildWidget.GrabFocus()
				}
			} else {
				b.LogError("new child object is not a Widget type: %v (%T)")
			}
		}
	}
	return nil
}

// GetOrientation is a convenience method for returning the orientation property
// value.
// See: SetOrientation()
//
// Locking: read
func (b *CBox) GetOrientation() (orientation cenums.Orientation) {
	b.RLock()
	defer b.RUnlock()
	var ok bool
	if v, err := b.GetStructProperty(PropertyOrientation); err != nil {
		b.LogErr(err)
	} else if orientation, ok = v.(cenums.Orientation); !ok && v != nil {
		b.LogError("invalid value stored in %v: %v (%T)", PropertyOrientation, v, v)
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
func (b *CBox) SetOrientation(orientation cenums.Orientation) {
	b.Lock()
	if err := b.SetStructProperty(PropertyOrientation, orientation); err != nil {
		b.LogErr(err)
	}
	b.Unlock()
	b.Resize()
}

// GetHomogeneous is a convenience method for returning the homogeneous property
// value.
// See: SetHomogeneous()
//
// Locking: read
func (b *CBox) GetHomogeneous() (value bool) {
	b.RLock()
	defer b.RUnlock()
	var err error
	if value, err = b.GetBoolProperty(PropertyHomogeneous); err != nil {
		b.LogErr(err)
	}
	return
}

// SetHomogeneous is a convenience method for updating the homogeneous property
// of the Box instance, controlling whether or not all children of box are given
// equal space in the box.
//
// Parameters:
// 	homogeneous	 TRUE to create equal allotments, FALSE for variable allotments
//
// Locking: write
func (b *CBox) SetHomogeneous(homogeneous bool) {
	b.Lock()
	defer b.Unlock()
	if err := b.SetBoolProperty(PropertyHomogeneous, homogeneous); err != nil {
		b.LogErr(err)
	}
}

// GetSpacing is a convenience method for returning the spacing property value.
// See: SetSpacing()
//
// Locking: read
func (b *CBox) GetSpacing() (value int) {
	b.RLock()
	defer b.RUnlock()
	var err error
	if value, err = b.GetIntProperty(PropertySpacing); err != nil {
		b.LogErr(err)
	}
	return
}

// SetSpacing is a convenience method to update the spacing property value.
//
// Parameters:
// 	spacing	 the number of characters to put between children
//
// Locking: write
func (b *CBox) SetSpacing(spacing int) {
	b.Lock()
	defer b.Unlock()
	if err := b.SetIntProperty(PropertySpacing, spacing); err != nil {
		b.LogErr(err)
	}
}

// Add the given Widget to the Box using PackStart() with default settings of:
// expand=false, fill=true and padding=0
//
// Locking: write
func (b *CBox) Add(child Widget) {
	b.PackStart(child, false, true, 0)
}

// Remove the given Widget from the Box Container, disconnecting any signal
// handlers in the process.
//
// Locking: write
func (b *CBox) Remove(w Widget) {
	_ = b.Disconnect(SignalShow, BoxChildShowHandle)
	_ = b.Disconnect(SignalHide, BoxChildHideHandle)
	b.CContainer.Remove(w)
}

// PackStart adds child to box, packed with reference to the start of box. The
// child is packed after any other child packed with reference to the start of
// box.
//
// Parameters
//
// child     the Widget to be added to box
// expand    TRUE if the new child is to be given extra space allocated to box.
//           The extra space will be divided evenly between all children of box
//           that use this option
// fill      TRUE if space given to child by the expand option is actually
//           allocated to child, rather than just padding it. This parameter has
//           no effect if expand is set to FALSE. A child is always allocated
//           the full height of an HBox and the full width of a VBox. This
//           option affects the other dimension
// padding   extra space in pixels to put between this child and its neighbors,
//           over and above the global amount specified by spacing property.
//           If child is a widget at one of the reference ends of box , then
//           padding pixels are also put between child and the reference edge of
//           box
//
// Locking: write
func (b *CBox) PackStart(child Widget, expand, fill bool, padding int) {
	b.LogDebug("expand=%v, fill=%v, padding=%v, child=%v", expand, fill, padding, child.ObjectName())
	if f := b.Emit(SignalPackStart, b, child, expand, fill, padding); f == cenums.EVENT_PASS {
		child.SetParent(b)
		child.SetWindow(b.GetWindow())
		b.CContainer.AddWithProperties(child,
			PropertyBoxChildPackType, enums.PackStart,
			PropertyBoxChildExpand, expand,
			PropertyBoxChildFill, fill,
			PropertyBoxChildPadding, padding,
		)
		b.Resize()
	}
}

// PackEnd adds child to box, packed with reference to the end of box. The child
// is packed after (away from end of) any other child packed with reference to
// the end of box.
//
// Parameters
//
// child     the Widget to be added to box
// expand    TRUE if the new child is to be given extra space allocated to box.
//           The extra space will be divided evenly between all children of box
//           that use this option
// fill      TRUE if space given to child by the expand option is actually
//           allocated to child, rather than just padding it. This parameter
//           has no effect if expand is set to FALSE. A child is always
//           allocated the full height of an HBox and the full width of a VBox.
//           This option affects the other dimension
// padding   extra space in pixels to put between this child and its neighbors,
//           over and above the global amount specified by spacing property.
//           If child is a widget at one of the reference ends of box, then
//           padding pixels are also put between child and the reference edge of
//           box
//
// Locking: write
func (b *CBox) PackEnd(child Widget, expand, fill bool, padding int) {
	b.LogDebug("expand=%v, fill=%v, padding=%v, child=%v", expand, fill, padding, child.ObjectName())
	if f := b.Emit(SignalPackEnd, b, child, expand, fill, padding); f == cenums.EVENT_PASS {
		child.SetParent(b)
		child.SetWindow(b.GetWindow())
		b.CContainer.AddWithProperties(child,
			PropertyBoxChildPackType, enums.PackEnd,
			PropertyBoxChildExpand, expand,
			PropertyBoxChildFill, fill,
			PropertyBoxChildPadding, padding,
		)
		b.Resize()
	}
}

// ReorderChild moves the given child to a new position in the list of Box
// children. The list is the children field of Box, and contains both widgets
// packed PACK_START as well as widgets packed PACK_END, in the order that these
// widgets were added to the box. A widget's position in the Box children list
// determines where the widget is packed into box. A child widget at some
// position in the list will be packed just after all other widgets of the
// same packing type that appear earlier in the list. The children field is not
// exported and only the interface methods are able to manipulate the field.
//
// Parameters:
// 	child	    the Widget to move
// 	position	the new position for child in the list of children of box starting from 0. If negative, indicates the end of the list
//
// Locking: write
func (b *CBox) ReorderChild(child Widget, position int) {
	childId := child.ObjectID()
	b.Lock()
	var children []Widget
	if position < 0 {
		position = len(b.children) - 1 + position
	}
	for idx, c := range b.children {
		b.Unlock()
		cId := c.ObjectID()
		b.Lock()
		if idx == position {
			children = append(children, child)
		} else if cId != childId {
			children = append(children, c)
		}
	}
	b.children = children
	b.Unlock()
}

// QueryChildPacking obtains information about how the child is packed into the
// Box. If the given child Widget is not contained within the Box an error is
// logged and the return values will all be their `nil` equivalents.
//
// Parameters:
// 	child	the Widget of the child to query
//
// Locking: read
func (b *CBox) QueryChildPacking(child Widget) (expand bool, fill bool, padding int, packType enums.PackType) {
	b.RLock()
	if cps, ok := b.property[child.ObjectID()]; ok {
		for _, cp := range cps {
			switch cp.Name() {
			case PropertyBoxChildExpand:
				if v, ok := cp.Value().(bool); ok {
					expand = v
				}
			case PropertyBoxChildFill:
				if v, ok := cp.Value().(bool); ok {
					fill = v
				}
			case PropertyBoxChildPadding:
				if v, ok := cp.Value().(int); ok {
					padding = v
				}
			case PropertyBoxChildPackType:
				if v, ok := cp.Value().(enums.PackType); ok {
					packType = v
				}
			}
		}
	} else {
		b.LogError("%v is not a child of %v", child, b)
	}
	b.RUnlock()
	return
}

// SetChildPacking updates the information about how the child is packed into
// the Box. If the given child Widget is not contained within the Box an error
// is logged and no action is taken.
//
// Parameters:
// 	child	the Widget of the child to set
// 	expand	the new value of the “expand” child property
// 	fill	the new value of the “fill” child property
// 	padding	the new value of the “padding” child property
// 	packType	the new value of the “pack-type” child property
//
// Locking: write
func (b *CBox) SetChildPacking(child Widget, expand bool, fill bool, padding int, packType enums.PackType) {
	if cps, ok := b.property[child.ObjectID()]; ok {
		b.Lock()
		defer b.Unlock()
		for _, cp := range cps {
			switch cp.Name() {
			case PropertyBoxChildExpand:
				if err := cp.Set(expand); err != nil {
					b.LogErr(err)
				}
			case PropertyBoxChildFill:
				if err := cp.Set(fill); err != nil {
					b.LogErr(err)
				}
			case PropertyBoxChildPadding:
				if err := cp.Set(padding); err != nil {
					b.LogErr(err)
				}
			case PropertyBoxChildPackType:
				if err := cp.Set(packType); err != nil {
					b.LogErr(err)
				}
			}
		}
	} else {
		b.LogError("%v is not a child of %v", child, b)
	}
}

// GetFocusChain retrieves the focus chain of the Box, if one has been set
// explicitly. If no focus chain has been explicitly set, CTK computes the
// focus chain based on the positions of the children, taking into account the
// child packing configuration.
//
// Returns:
// 	focusableWidgets	widgets in the focus chain.
// 	explicitlySet       TRUE if the focus chain has been set explicitly.
//
// Locking: read
func (b *CBox) GetFocusChain() (focusableWidgets []interface{}, explicitlySet bool) {
	b.RLock()
	if b.focusChainSet {
		b.RUnlock()
		return b.focusChain, true
	}
	b.RUnlock()
	boxChildren := b.getBoxChildren()
	b.RLock()
	var children []interface{}
	for _, child := range boxChildren {
		if child.packType == enums.PackStart {
			children = append(children, child.widget)
		}
	}
	for _, child := range boxChildren {
		if child.packType == enums.PackEnd {
			children = append(children, child.widget)
		}
	}
	for _, child := range children {
		if cc, ok := child.(Container); ok {
			fc, _ := cc.GetFocusChain()
			for _, cChild := range fc {
				focusableWidgets = append(focusableWidgets, cChild)
			}
		} else if cw, ok := child.(Widget); ok {
			if cw.CanFocus() && cw.IsVisible() && cw.IsSensitive() {
				focusableWidgets = append(focusableWidgets, child)
			}
		}
	}
	b.RUnlock()
	return
}

// GetSizeRequest returns the requested size of the Drawable Widget. This method
// is used by Container Widgets to resolve the surface space allocated for their
// child Widget instances.
//
// Locking: read
func (b *CBox) GetSizeRequest() (width, height int) {
	children := b.getBoxChildren()

	b.RLock()
	nChildren := len(children)
	if nChildren <= 0 {
		b.RUnlock()
		return
	}
	b.RUnlock()

	orientation := b.GetOrientation()
	spacing := b.GetSpacing()
	isHomogeneous := b.GetHomogeneous()

	b.Lock()
	defer b.Unlock()

	isVertical := orientation == cenums.ORIENTATION_VERTICAL
	var w, h int

	if isHomogeneous {
		// get the size of the largest child and request that for all children
		for _, child := range children {
			req := ptypes.MakeRectangle(child.widget.GetSizeRequest())
			if w < req.W {
				w = req.W
				if !isVertical && child.padding > 0 {
					w += child.padding * 2
				}
			}
			if h < req.H {
				h = req.H
				if isVertical && child.padding > 0 {
					h += child.padding * 2
				}
			}
		}
		if isVertical {
			width = w
			height = (nChildren * h) + cmath.FloorI((nChildren-1)*spacing, 0)
		} else {
			width = (nChildren * w) + cmath.FloorI((nChildren-1)*spacing, 0)
			height = h
		}
		return
	}
	// add up the sizes of all children, including spacing and child padding
	sizes := make([]*ptypes.Rectangle, nChildren)
	tally := ptypes.NewRectangle(0, 0)
	for idx, child := range children {
		sizes[idx] = ptypes.NewRectangle(child.widget.GetSizeRequest())
		sizes[idx].Floor(0, 0)
		if w < sizes[idx].W {
			w = sizes[idx].W
			if !isVertical && child.padding > 0 {
				w += child.padding * 2
			}
		}
		tally.W += sizes[idx].W
		if !isVertical && child.padding > 0 {
			tally.W += child.padding * 2
		}
		if h < sizes[idx].H {
			h = sizes[idx].H
			if isVertical && child.padding > 0 {
				h += child.padding * 2
			}
		}
		tally.H += sizes[idx].H
		if !isVertical && child.padding > 0 {
			tally.H += child.padding * 2
		}
	}
	if isVertical {
		width = w
		height = tally.H
	} else {
		width = tally.W
		height = h
	}
	return
}

func (b *CBox) invalidate(data []interface{}, argv ...interface{}) cenums.EventFlag {
	origin := b.GetOrigin()
	style := b.GetThemeRequest().Content.Normal
	for _, child := range b.getBoxChildren() {
		b.Lock()
		local := child.widget.GetOrigin()
		local.SubPoint(origin)
		alloc := child.widget.GetAllocation()
		if err := memphis.ConfigureSurface(child.widget.ObjectID(), local, alloc, style); err != nil {
			child.widget.LogErr(err)
		}
		b.Unlock()
	}
	return cenums.EVENT_STOP
}

func (b *CBox) resize(data []interface{}, argv ...interface{}) cenums.EventFlag {
	children := b.getBoxChildren()
	numChildren := len(children)
	if numChildren == 0 {
		return cenums.EVENT_STOP
	}
	spacing := b.GetSpacing()
	origin := b.GetOrigin().NewClone()
	alloc := b.GetAllocation().NewClone()
	orientation := b.GetOrientation()
	isVertical := orientation == cenums.ORIENTATION_VERTICAL
	homogeneous := b.GetHomogeneous()
	b.Lock()
	// intermediaries
	var increment int
	var gaps []int
	if isVertical {
		increment, gaps = cmath.SolveSpaceAlloc(numChildren, alloc.H, spacing)
	} else {
		increment, gaps = cmath.SolveSpaceAlloc(numChildren, alloc.W, spacing)
	}
	nextPoint := origin.NewClone()
	b.Unlock()
	if homogeneous {
		return b.resizeHomogeneous(isVertical, gaps, increment, numChildren, origin, nextPoint, alloc, children)
	}
	return b.resizeDynamicAlloc(isVertical, gaps, increment, spacing, numChildren, origin, nextPoint, alloc, children)
}

func (b *CBox) resizeHomogeneous(isVertical bool, gaps []int, increment, numChildren int, origin, nextPoint *ptypes.Point2I, alloc *ptypes.Rectangle, children []*cBoxChild) cenums.EventFlag {
	style := b.GetThemeRequest().Content.Normal
	// b.Lock()
	// defer b.Unlock()

	// assume child.expand == true
	var consumed int
	tracking := make([]struct {
		x, y, w, h int
		rw, rh     int
		aw, ah     int
		extra      int
		overflow   int
	}, numChildren)

	// first: build up tracking dataset

	for idx, child := range children {
		req := ptypes.NewRectangle(child.widget.GetSizeRequest())
		if child.fill {
			if isVertical {
				tracking[idx].w = alloc.W
				tracking[idx].h = increment
			} else { // horizontal
				tracking[idx].w = increment
				tracking[idx].h = alloc.H
			}
			tracking[idx].aw = tracking[idx].w
			tracking[idx].ah = tracking[idx].h
		} else { // !child.fill
			if isVertical {
				tracking[idx].w = alloc.W
				tracking[idx].aw = alloc.W
				tracking[idx].ah = increment
				if req.H <= -1 && req.H > increment {
					tracking[idx].h = increment
					req.H = increment
				} else {
					tracking[idx].h = req.H
				}
			} else { // horizontal
				tracking[idx].h = alloc.H
				tracking[idx].aw = increment
				tracking[idx].ah = alloc.H
				if req.W <= -1 && req.W > increment {
					tracking[idx].w = increment
					req.W = increment
				} else {
					tracking[idx].w = req.W
				}
			} // if isVertical
		} // if child.fill
		req.Floor(0, 0)
		tracking[idx].rw = req.W
		tracking[idx].rh = req.H
	}
	for idx, _ := range children {
		if isVertical {
			consumed += tracking[idx].ah
		} else {
			consumed += tracking[idx].aw
		}
	} // for each child

	// solve positions

	for idx, child := range children {
		if isVertical {
			if tracking[idx].h < tracking[idx].ah {
				delta := tracking[idx].ah - tracking[idx].h
				if ca, ok := child.widget.(Alignable); ok {
					_, yAlign := ca.GetAlignment()
					pad := int(float64(delta) * yAlign)
					tracking[idx].y += pad
					tracking[idx].overflow += delta - pad
				} else {
					tracking[idx].overflow += delta
				}
			}
		} else { // horizontal
			if tracking[idx].w < tracking[idx].aw {
				delta := tracking[idx].aw - tracking[idx].w
				if ca, ok := child.widget.(Alignable); ok {
					xAlign, _ := ca.GetAlignment()
					pad := int(float64(delta) * xAlign)
					tracking[idx].x += pad
					tracking[idx].overflow += delta - pad
				} else {
					tracking[idx].overflow += delta
				}
			}
		} // if isVertical
	} // for each child

	// allocate children and update canvas

	for idx, child := range children {
		local := ptypes.NewPoint2I(tracking[idx].x, tracking[idx].y)
		childSize := ptypes.NewRectangle(tracking[idx].w, tracking[idx].h)
		childAlloc := ptypes.NewRectangle(tracking[idx].aw, tracking[idx].ah)
		nextPoint.Add(local.X, local.Y)
		x := nextPoint.X - origin.X
		y := nextPoint.Y - origin.Y
		if err := memphis.ConfigureSurface(child.widget.ObjectID(), ptypes.MakePoint2I(x, y), ptypes.MakeRectangle(childAlloc.W, childAlloc.H), style); err != nil {
			child.widget.LogErr(err)
		}
		child.widget.SetOrigin(nextPoint.X, nextPoint.Y)
		child.widget.SetAllocation(*childSize)
		child.widget.Resize()
		if isVertical {
			nextPoint.Y += tracking[idx].h + tracking[idx].overflow
		} else {
			nextPoint.X += tracking[idx].w + tracking[idx].overflow
		}
		if len(gaps) > idx {
			if isVertical {
				nextPoint.Y += gaps[idx]
			} else {
				nextPoint.X += gaps[idx]
			}
		}
	}

	return cenums.EVENT_STOP
}

func (b *CBox) resizeDynamicAlloc(isVertical bool, gaps []int, increment, spacing, numChildren int, origin, nextPoint *ptypes.Point2I, alloc *ptypes.Rectangle, children []*cBoxChild) cenums.EventFlag {
	style := b.GetThemeRequest().Content.Normal
	// b.Lock()
	// defer b.Unlock()
	var (
		extraSpace   int
		numExpanding int
		consumed     int
	)
	tracking := make([]struct {
		x, y, w, h int
		rw, rh     int
		aw, ah     int
		extra      int
		overflow   int
	}, numChildren)

	// TODO: PackEnd children need resizeDynamicAlloc to go RtL, right aligned

	for idx, child := range children {
		req := ptypes.NewRectangle(child.widget.GetSizeRequest())
		if child.expand { // expand
			numExpanding += 1
			if child.fill { // expand && fill
				if isVertical {
					tracking[idx].w = alloc.W
					tracking[idx].h = increment
				} else {
					tracking[idx].w = increment
					tracking[idx].h = alloc.H
				}
				tracking[idx].aw = tracking[idx].w
				tracking[idx].ah = tracking[idx].h
			} else { // expand && !fill
				if isVertical { // expand && !fill && vertical
					if req.H <= -1 || req.H > increment {
						req.H = increment
					}
					tracking[idx].w = alloc.W
					tracking[idx].h = req.H
					tracking[idx].aw = alloc.W
					tracking[idx].ah = increment
				} else { // expand && !fill && horizontal
					if req.W <= -1 || req.W > increment {
						req.W = increment
					}
					tracking[idx].w = req.W
					tracking[idx].h = alloc.H
					tracking[idx].aw = increment
					tracking[idx].ah = alloc.H
				}
			} // else expand, !fill
		} else { // if !expand (assume !fill)
			if isVertical { // !expand, !fill, vertical
				req.W = alloc.W // force width
				if req.H <= -1 || req.H > increment {
					req.H = increment
				}
				if req.H < increment {
					delta := increment - req.H
					extraSpace += delta
				}
			} else { // !expand, !fill, horizontal
				req.H = alloc.H // force height
				if req.W <= -1 || req.W > increment {
					req.W = increment
				}
				if req.W < increment {
					delta := increment - req.W
					extraSpace += delta
				}
			}
			tracking[idx].w = req.W
			tracking[idx].h = req.H
			tracking[idx].aw = req.W
			tracking[idx].ah = req.H
		} // else expand
		tracking[idx].rw = req.W
		tracking[idx].rh = req.H
	}

	for idx, _ := range children {
		if isVertical {
			consumed += tracking[idx].ah
		} else {
			consumed += tracking[idx].aw
		}
	}
	var total = 0
	if isVertical {
		total = alloc.H - consumed
	} else {
		total = alloc.W - consumed
	}
	if extraSpace > 0 {
		var (
			dist []int
			err  error
		)
		if dist, gaps, err = cmath.Distribute(total, extraSpace, numExpanding, numChildren, spacing); err != nil {
			b.LogErr(err)
		} else {
			di, nDist := 0, len(dist)
			for idx, child := range children {
				if child.expand {
					if di < nDist {
						if child.fill {
							if isVertical {
								tracking[idx].ah += dist[di]
								tracking[idx].h = tracking[idx].ah
							} else {
								tracking[idx].aw += dist[di]
								tracking[idx].w = tracking[idx].aw
							}
						} else { // !fill
							if isVertical {
								tracking[idx].ah += dist[di]
							} else {
								tracking[idx].aw += dist[di]
							}
						}
						di += 1
					}
				}
			}
		}
	}

	// solve positions

	for idx, child := range children {
		if ca, ok := child.widget.(Alignable); ok {
			xAlign, yAlign := ca.GetAlignment()
			if isVertical {
				if tracking[idx].h < tracking[idx].ah {
					delta := tracking[idx].ah - tracking[idx].h
					inc := int(float64(delta) * yAlign)
					tracking[idx].y += inc
					tracking[idx].overflow += delta - inc
				}
			} else {
				if tracking[idx].w < tracking[idx].aw {
					delta := tracking[idx].aw - tracking[idx].w
					inc := int(float64(delta) * xAlign)
					tracking[idx].x += inc
					tracking[idx].overflow += delta - inc
				}
			}
		}
	}

	// allocate space and update canvas

	for idx, child := range children {
		track := tracking[idx]
		nextPoint.X += track.x
		nextPoint.Y += track.y
		child.widget.SetOrigin(nextPoint.X, nextPoint.Y)
		child.widget.SetAllocation(ptypes.MakeRectangle(track.w, track.h))
		child.widget.Resize()
		if err := memphis.ConfigureSurface(child.widget.ObjectID(), ptypes.MakePoint2I(nextPoint.X-origin.X, nextPoint.Y-origin.Y), ptypes.MakeRectangle(track.w, track.h), style); err != nil {
			child.widget.LogErr(err)
		}
		if isVertical {
			nextPoint.Y += track.h + track.overflow
		} else {
			nextPoint.X += track.w + track.overflow
		}
		if len(gaps) > idx {
			if isVertical {
				nextPoint.Y += gaps[idx]
			} else {
				nextPoint.X += gaps[idx]
			}
		}
	}

	return cenums.EVENT_STOP
}

func (b *CBox) draw(data []interface{}, argv ...interface{}) cenums.EventFlag {
	if surface, ok := argv[1].(*memphis.CSurface); ok {
		alloc := b.GetAllocation()
		if !b.IsVisible() || alloc.W <= 0 || alloc.H <= 0 {
			b.LogTrace("not visible, zero width or zero height")
			return cenums.EVENT_PASS
		}
		b.RLock()
		debug, _ := b.GetBoolProperty(cdk.PropertyDebug)
		debugChildren, _ := b.GetBoolProperty(PropertyDebugChildren)
		b.RUnlock()
		orientation := b.GetOrientation()
		children := b.getBoxChildren()
		theme := b.GetThemeRequest()
		b.Lock()
		defer b.Unlock()
		surface.Fill(theme)
		for _, child := range children {
			if child.widget.IsVisible() {
				if f := child.widget.Draw(); f == cenums.EVENT_STOP {
					if childSurface, err := memphis.GetSurface(child.widget.ObjectID()); err != nil {
						child.widget.LogErr(err)
					} else {
						if debugChildren && orientation == cenums.ORIENTATION_VERTICAL {
							childSurface.DebugBox(paint.ColorPink, child.widget.ObjectInfo()+" ["+b.ObjectInfo()+"]")
						} else if debugChildren {
							childSurface.DebugBox(paint.ColorPurple, child.widget.ObjectInfo()+" ["+b.ObjectInfo()+"]")
						}
						if err := surface.CompositeSurface(childSurface); err != nil {
							b.LogError("composite error: %v", err)
						}
					}
				}
			}
		}
		if debug && orientation == cenums.ORIENTATION_VERTICAL {
			surface.DebugBox(paint.ColorPink, b.ObjectInfo())
		} else if debug {
			surface.DebugBox(paint.ColorPurple, b.ObjectInfo())
		}
	}
	return cenums.EVENT_STOP
}

func (b *CBox) getBoxChildren() (children []*cBoxChild) {
	bChildren := b.GetChildren()
	b.RLock()
	nChildren := len(bChildren)
	expand := make([]bool, nChildren)
	fill := make([]bool, nChildren)
	padding := make([]int, nChildren)
	packType := make([]enums.PackType, nChildren)
	for idx, child := range bChildren {
		b.RUnlock()
		expand[idx], fill[idx], padding[idx], packType[idx] = b.QueryChildPacking(child)
		b.RLock()
		if child.IsVisible() && packType[idx] == enums.PackStart {
			children = append(children, &cBoxChild{
				widget:   child,
				expand:   expand[idx],
				fill:     fill[idx],
				padding:  padding[idx],
				packType: packType[idx],
			})
		}
	}
	for idx, child := range bChildren {
		if child.IsVisible() && packType[idx] == enums.PackEnd {
			children = append(children, &cBoxChild{
				widget:   child,
				expand:   expand[idx],
				fill:     fill[idx],
				padding:  padding[idx],
				packType: packType[idx],
			})
		}
	}
	b.RUnlock()
	return
}

// Whether the children should all be the same size.
// Flags: Read / Write
// Default value: FALSE
const PropertyHomogeneous cdk.Property = "homogeneous"

// The amount of space between children.
// Flags: Read / Write
// Allowed values: >= 0
// Default value: 0
const PropertySpacing cdk.Property = "spacing"

const PropertyDebugChildren cdk.Property = "debug-children"

const PropertyBoxChildPackType cdk.Property = "box-child--pack-type"

const PropertyBoxChildExpand cdk.Property = "box-child--expand"

const PropertyBoxChildFill cdk.Property = "box-child--fill"

const PropertyBoxChildPadding cdk.Property = "box-child--padding"

const BoxChildShowHandle = "box-child-show-handler"

const BoxChildHideHandle = "box-child-hide-handler"

const BoxInvalidateHandle = "box-invalidate-handler"

const BoxResizeHandle = "box-resize-handler"

const BoxDrawHandle = "box-draw-handler"
