package ctk

import (
	"strconv"
	"strings"

	"github.com/go-curses/cdk"
	cenums "github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/cdk/lib/paint"
	cstrings "github.com/go-curses/cdk/lib/strings"
	"github.com/go-curses/cdk/memphis"
	"github.com/go-curses/ctk/lib/enums"
)

const TypeButtonBox cdk.CTypeTag = "ctk-button-box"

func init() {
	_ = cdk.TypesManager.AddType(TypeButtonBox, func() interface{} { return MakeButtonBox() })
}

// ButtonBox Hierarchy:
//	Object
//	  +- Widget
//	    +- Container
//	      +- Box
//	        +- ButtonBox
//	          +- HButtonBox
//	          +- VButtonBox
//
// The ButtonBox Widget is a Box Container that has a primary and a secondary
// grouping of its Widget children. These are typically used by the Dialog
// Widget to implement the action buttons.
type ButtonBox interface {
	Box
	Buildable
	Orientable

	Init() (already bool)
	Build(builder Builder, element *CBuilderElement) error
	GetLayout() (value enums.ButtonBoxStyle)
	SetLayout(layoutStyle enums.ButtonBoxStyle)
	GetChildren() (children []Widget)
	Add(w Widget)
	Remove(w Widget)
	PackStart(w Widget, expand, fill bool, padding int)
	PackEnd(w Widget, expand, fill bool, padding int)
	GetChildPrimary(w Widget) (isPrimary bool)
	GetChildSecondary(w Widget) (isSecondary bool)
	SetChildSecondary(child Widget, isSecondary bool)
	SetChildPacking(child Widget, expand bool, fill bool, padding int, packType enums.PackType)
	SetSpacing(spacing int)
}

// The CButtonBox structure implements the ButtonBox interface and is exported
// to facilitate type embedding with custom implementations. No member variables
// are exported as the interface methods are the only intended means of
// interacting with ButtonBox objects.
type CButtonBox struct {
	CBox
}

// MakeButtonBox is used by the Buildable system to construct a new horizontal
// homogeneous ButtonBox with no spacing between the Widget children.
func MakeButtonBox() ButtonBox {
	return NewButtonBox(cenums.ORIENTATION_HORIZONTAL, false, 0)
}

// NewButtonBox is a constructor for new Box instances.
func NewButtonBox(orientation cenums.Orientation, homogeneous bool, spacing int) ButtonBox {
	b := new(CButtonBox)
	b.Init()
	b.Freeze()
	b.SetOrientation(orientation)
	b.SetHomogeneous(homogeneous)
	b.SetSpacing(spacing)
	b.Thaw()
	return b
}

// Init initializes a ButtonBox object. This must be called at least once to
// set up the necessary defaults and allocate any memory structures. Calling
// this more than once is safe though unnecessary. Only the first call will
// result in any effect upon the ButtonBox instance. Init is used in the
// NewButtonBox constructor and only necessary when implementing a derivative
// ButtonBox type.
func (b *CButtonBox) Init() (already bool) {
	if b.InitTypeItem(TypeButtonBox, b) {
		return true
	}
	b.CBox.Init()
	b.flags = enums.NULL_WIDGET_FLAG
	b.SetFlags(enums.PARENT_SENSITIVE | enums.APP_PAINTABLE)
	orientation := b.GetOrientation()
	spacing := b.GetSpacing()
	primary := NewBox(orientation, false, spacing)
	primary.Show()
	b.CBox.PackStart(primary, true, true, 0)
	secondary := NewBox(orientation, false, spacing)
	secondary.Show()
	b.CBox.PackEnd(secondary, true, true, 0)
	_ = b.InstallProperty(PropertyLayoutStyle, cdk.StructProperty, true, enums.LayoutStart)
	b.Connect(SignalDraw, ButtonBoxDrawHandle, b.draw)
	return false
}

// Build provides customizations to the Buildable system for ButtonBox Widgets.
func (b *CButtonBox) Build(builder Builder, element *CBuilderElement) error {
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
					var expand, fill = false, true
					var padding = 0
					var packType = enums.PackStart
					for k, v := range child.Packing {
						switch k {
						case "expand":
							expand = cstrings.IsTrue(v)
						case "fill":
							fill = cstrings.IsTrue(v)
						case "padding":
							var err error
							if padding, err = strconv.Atoi(v); err != nil {
								b.LogErr(err)
								padding = 0
							}
						case "pack-type":
							switch strings.ToLower(v) {
							case "start":
							case "end":
								packType = enums.PackEnd
							default:
								b.LogError("invalid pack-type given: %v, must be either start or end")
							}
						}
					}
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

// GetLayout is a convenience method for returning the layout-style property
// value as the ButtonBoxStyle type.
// See: SetLayout()
func (b *CButtonBox) GetLayout() (value enums.ButtonBoxStyle) {
	if v, err := b.GetStructProperty(PropertyLayoutStyle); err == nil {
		var ok bool
		if value, ok = v.(enums.ButtonBoxStyle); !ok {
			b.LogError("value stored in %v is not of ButtonBoxStyle type: %v (%T)", v, v)
		}
	} else {
		b.LogErr(err)
	}
	return
}

// SetLayout is a convenience method for updating the layout-style property for
// the ButtonBox. The normal expand and fill packing options are ignored in the
// context of a ButtonBox. The layout style instead defines the visual placement
// of the child widgets. The size request for each child determines it's
// allocation in the ButtonBox with the exception of the "expand" layout style
// in which the children are evenly allocated to consume all available space in
// the ButtonBox.
//
// The different styles are as follows:
//   "start"      group Widgets from the starting edge
//   "end"        group Widgets from the ending edge
//   "center"     group Widgets together in the center, away from edges
//   "spread"     spread Widgets evenly and centered, away form edges
//   "edge"       spread Widgets evenly and centered, flush with edges
//   "expand"     expand all Widgets to evenly consume all available space
//
// Note that usage of this within CTK is unimplemented at this time
func (b *CButtonBox) SetLayout(layoutStyle enums.ButtonBoxStyle) {
	if err := b.SetStructProperty(PropertyLayoutStyle, layoutStyle); err != nil {
		b.LogErr(err)
	}
}

// GetChildren returns the children of the primary and secondary groupings.
func (b *CButtonBox) GetChildren() (children []Widget) {
	for _, child := range b.getPrimary().GetChildren() {
		children = append(children, child)
	}
	for _, child := range b.getSecondary().GetChildren() {
		children = append(children, child)
	}
	return
}

// Add is a convenience method for adding the given Widget to the primary group
// with default PackStart configuration of: expand=true, fill=true and padding=0
func (b *CButtonBox) Add(w Widget) {
	if primary := b.getPrimary(); primary != nil {
		primary.PackStart(w, true, true, 0)
	}
}

// Remove the given Widget from the ButtonBox
func (b *CButtonBox) Remove(w Widget) {
	if b.GetChildSecondary(w) {
		if secondary := b.getSecondary(); secondary != nil {
			secondary.Remove(w)
		}
	} else {
		if primary := b.getPrimary(); primary != nil {
			primary.Remove(w)
		}
	}
}

// PackStart will add the given Widget to the primary group with the given Box
// packing configuration.
func (b *CButtonBox) PackStart(w Widget, expand, fill bool, padding int) {
	if primary := b.getPrimary(); primary != nil {
		primary.PackStart(w, expand, fill, padding)
	}
}

// PackEnd will add the given Widget to the secondary group with the given Box
// packing configuration.
func (b *CButtonBox) PackEnd(w Widget, expand, fill bool, padding int) {
	if secondary := b.getSecondary(); secondary != nil {
		secondary.PackEnd(w, expand, fill, padding)
	}
}

// GetChildPrimary is a convenience method that returns TRUE if the given Widget
// is in the primary grouping and returns FALSE otherwise.
//
// Parameters:
// 	child	a child of widget
func (b *CButtonBox) GetChildPrimary(w Widget) (isPrimary bool) {
	if box := b.getPrimary(); box != nil {
		for _, child := range box.GetChildren() {
			if child.ObjectID() == w.ObjectID() {
				return true
			}
		}
	}
	return false
}

// GetChildSecondary is a convenience method that returns TRUE if the given
// Widget is in the primary grouping and returns FALSE otherwise.
//
// Parameters:
// 	child	a child of widget
func (b *CButtonBox) GetChildSecondary(w Widget) (isSecondary bool) {
	if box := b.getSecondary(); box != nil {
		for _, child := range box.GetChildren() {
			if child.ObjectID() == w.ObjectID() {
				return true
			}
		}
	}
	return false
}

// SetChildSecondary will ensure the given Widget is in the secondary grouping
// if the isSecondary argument is TRUE. If isSecondary is FALSE, this will
// ensure that the Widget is in the primary grouping.
//
// Parameters:
// 	child	a child of widget
// 	secondary	TRUE, if the child appears in a secondary group
func (b *CButtonBox) SetChildSecondary(child Widget, isSecondary bool) {
	defer b.Invalidate()
	alreadySecondary := b.GetChildSecondary(child)
	if isSecondary && alreadySecondary {
		return
	}
	if isSecondary {
		if pBox := b.getPrimary(); pBox != nil {
			pBox.Remove(child)
		}
		if sBox := b.getSecondary(); sBox != nil {
			sBox.PackEnd(child, false, false, 0)
		}
		return
	}
	if b.GetChildPrimary(child) {
		return
	}
	if sBox := b.getSecondary(); sBox != nil {
		sBox.Remove(child)
	}
	if pBox := b.getPrimary(); pBox != nil {
		pBox.PackStart(child, false, false, 0)
	}
}

// SetChildPacking is a convenience method to set the packing configuration for
// the given child Widget of the ButtonBox, regardless of which grouping the
// Widget is in.
//
// Parameters:
// 	child	the Widget of the child to set
// 	expand	the new value of the expand child property
// 	fill	the new value of the fill child property
// 	padding	the new value of the padding child property
// 	packType	the new value of the pack-type child property
func (b *CButtonBox) SetChildPacking(child Widget, expand bool, fill bool, padding int, packType enums.PackType) {
	if b.GetChildPrimary(child) {
		if primary := b.getPrimary(); primary != nil {
			primary.SetChildPacking(child, expand, fill, padding, packType)
		}
	} else if b.GetChildSecondary(child) {
		if secondary := b.getSecondary(); secondary != nil {
			secondary.SetChildPacking(child, expand, fill, padding, packType)
		}
	} else {
		b.LogError("%v is not a child of %v", child, b)
	}
}

func (b *CButtonBox) SetSpacing(spacing int) {
	b.CBox.SetSpacing(spacing)
	b.getPrimary().SetSpacing(spacing)
	b.getSecondary().SetSpacing(spacing)
}

func (b *CButtonBox) getPrimary() (box Box) {
	children := b.CBox.GetChildren()
	if len(children) > 0 {
		var ok bool
		if box, ok = children[0].(Box); !ok {
			box = nil
		} else {
			return
		}
	}
	b.LogError("button box missing primary container")
	return nil
}

func (b *CButtonBox) getSecondary() (box Box) {
	children := b.CBox.GetChildren()
	if len(children) > 1 {
		var ok bool
		if box, ok = children[1].(Box); !ok {
			box = nil
		} else {
			return
		}
	}
	b.LogError("button box missing secondary container")
	return nil
}

func (b *CButtonBox) draw(data []interface{}, argv ...interface{}) cenums.EventFlag {
	if surface, ok := argv[1].(*memphis.CSurface); ok {
		alloc := b.GetAllocation()
		if !b.IsVisible() || alloc.W <= 0 || alloc.H <= 0 {
			b.LogTrace("not visible, zero width or zero height")
			return cenums.EVENT_PASS
		}

		b.LockDraw()
		defer b.UnlockDraw()

		debug, _ := b.GetBoolProperty(cdk.PropertyDebug)
		debugChildren, _ := b.GetBoolProperty(PropertyDebugChildren)
		orientation := b.GetOrientation()
		children := b.getBoxChildren()
		theme := b.GetThemeRequest()
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
		return cenums.EVENT_STOP
	}
	return cenums.EVENT_PASS
}

// How to lay out the buttons in the box. Possible values are: default, spread, edge, start and end.
// Flags: Read / Write
// Default value: GTK_BUTTONBOX_DEFAULT_STYLE
const PropertyLayoutStyle cdk.Property = "layout-style"

const ButtonBoxDrawHandle = "button-box-draw-handler"
