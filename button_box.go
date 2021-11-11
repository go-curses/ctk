package ctk

import (
	"strconv"
	"strings"

	"github.com/go-curses/cdk"
	"github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/cdk/lib/paint"
	cstrings "github.com/go-curses/cdk/lib/strings"
	"github.com/go-curses/cdk/memphis"
)

// CDK type-tag for ButtonBox objects
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
type ButtonBox interface {
	Box
	Buildable
	Orientable

	Init() (already bool)
	Build(builder Builder, element *CBuilderElement) error
	GetLayout() (value ButtonBoxStyle)
	GetChildSecondary(w Widget) (isSecondary bool)
	GetChildPrimary(w Widget) (isPrimary bool)
	SetLayout(layoutStyle ButtonBoxStyle)
	SetChildSecondary(child Widget, isSecondary bool)
	Show()
	ShowAll()
	Hide()
	PackStart(w Widget, expand, fill bool, padding int)
	PackEnd(w Widget, expand, fill bool, padding int)
	SetChildPacking(child Widget, expand bool, fill bool, padding int, packType PackType)
	Add(w Widget)
	Remove(w Widget)
	Resize() enums.EventFlag
	Invalidate() enums.EventFlag
}

// The CButtonBox structure implements the ButtonBox interface and is
// exported to facilitate type embedding with custom implementations. No member
// variables are exported as the interface methods are the only intended means
// of interacting with ButtonBox objects
type CButtonBox struct {
	CBox
}

func MakeButtonBox() *CButtonBox {
	return NewButtonBox(enums.ORIENTATION_HORIZONTAL, false, 0)
}

func NewButtonBox(orientation enums.Orientation, homogeneous bool, spacing int) *CButtonBox {
	b := new(CButtonBox)
	b.Init()
	b.Freeze()
	b.SetOrientation(orientation)
	b.SetHomogeneous(homogeneous)
	b.SetSpacing(spacing)
	b.Connect(SignalDraw, ButtonBoxDrawHandle, b.draw)
	b.Thaw()
	return b
}

// ButtonBox object initialization. This must be called at least once to setup
// the necessary defaults and allocate any memory structures. Calling this more
// than once is safe though unnecessary. Only the first call will result in any
// effect upon the ButtonBox instance
func (b *CButtonBox) Init() (already bool) {
	if b.InitTypeItem(TypeButtonBox, b) {
		return true
	}
	b.CBox.Init()
	b.flags = NULL_WIDGET_FLAG
	b.SetFlags(PARENT_SENSITIVE | APP_PAINTABLE)
	orientation := b.GetOrientation()
	spacing := b.GetSpacing()
	b.CBox.PackStart(NewBox(orientation, false, spacing), true, true, 0)
	b.CBox.PackEnd(NewBox(orientation, false, spacing), true, true, 0)
	_ = b.InstallProperty(PropertyLayoutStyle, cdk.StructProperty, true, LayoutStart)
	// b.Connect(SignalDraw, ButtonBoxDrawHandle, b.draw)
	return false
}

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
					var packType = PackStart
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
								packType = PackEnd
							default:
								b.LogError("invalid pack-type given: %v, must be either start or end")
							}
						}
					}
					if packType == PackStart {
						b.PackStart(newChildWidget, expand, fill, padding)
					} else {
						b.PackEnd(newChildWidget, expand, fill, padding)
					}
				} else {
					b.Add(newChildWidget)
				}
				if newChildWidget.HasFlags(HAS_FOCUS) {
					newChildWidget.GrabFocus()
				}
			} else {
				b.LogError("new child object is not a Widget type: %v (%T)")
			}
		}
	}
	return nil
}
func (b *CButtonBox) GetLayout() (value ButtonBoxStyle) {
	if v, err := b.GetStructProperty(PropertyLayoutStyle); err == nil {
		var ok bool
		if value, ok = v.(ButtonBoxStyle); !ok {
			b.LogError("value stored in %v is not of ButtonBoxStyle type: %v (%T)", v, v)
		}
	} else {
		b.LogErr(err)
	}
	return
}

// Returns whether child should appear in a secondary group of children.
// Parameters:
// 	child	a child of widget
// Returns:
// 	whether child should appear in a secondary group of children.
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

// Returns whether child should appear in the primary group of children.
// Parameters:
// 	child	a child of widget
// Returns:
// 	whether child should appear in the primary group of children.
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

// Sets the layout for the ButtonBox. The normal expand and fill packing options
// are ignored in the context of a ButtonBox. The layout style instead defines
// the visual placement of the child widgets. The size request for each child
// determines it's allocation in the ButtonBox with the exception of the
// "expand" layout style in which the children are evenly allocated to consume
// all available space in the ButtonBox.
//
// The different styles are as follows:
//   "start"      group Widgets from the starting edge
//   "end"        group Widgets from the ending edge
//   "center"     group Widgets together in the center, away from edges
//   "spread"     spread Widgets evenly and centered, away form edges
//   "edge"       spread Widgets evenly and centered, flush with edges
//   "expand"     expand all Widgets to evenly consume all available space
func (b *CButtonBox) SetLayout(layoutStyle ButtonBoxStyle) {
	if err := b.SetStructProperty(PropertyLayoutStyle, layoutStyle); err != nil {
		b.LogErr(err)
	}
}

// Sets whether child should appear in a secondary group of children.
//   A typical use of a secondary child is the help button in a dialog.
// This group appears after the other children if the style
//   is CTK_BUTTONBOX_START, CTK_BUTTONBOX_SPREAD or
//   CTK_BUTTONBOX_EDGE, and before the other children if the style
//   is CTK_BUTTONBOX_END. For horizontal button boxes, the definition
//   of before/after depends on direction of the widget (see
//   gtk_widget_set_direction()). If the style is CTK_BUTTONBOX_START
//   or CTK_BUTTONBOX_END, then the secondary children are aligned at
//   the other end of the button box from the main children. For the
//   other styles, they appear immediately next to the main children.
// Parameters:
// 	child	a child of widget
// 	secondary	if TRUE, the child appears in a secondary group of the button box.
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

func (b *CButtonBox) Show() {
	b.CBox.Show()
	if primary := b.getPrimary(); primary != nil {
		primary.Show()
	}
	if secondary := b.getSecondary(); secondary != nil {
		secondary.Show()
	}
}

func (b *CButtonBox) ShowAll() {
	b.CBox.ShowAll()
	if primary := b.getPrimary(); primary != nil {
		primary.ShowAll()
	}
	if secondary := b.getSecondary(); secondary != nil {
		secondary.ShowAll()
	}
}

func (b *CButtonBox) Hide() {
	b.CBox.Hide()
	if primary := b.getPrimary(); primary != nil {
		primary.Hide()
	}
	if secondary := b.getSecondary(); secondary != nil {
		secondary.Hide()
	}
}

func (b *CButtonBox) PackStart(w Widget, expand, fill bool, padding int) {
	if primary := b.getPrimary(); primary != nil {
		primary.PackStart(w, expand, fill, padding)
	}
}

func (b *CButtonBox) PackEnd(w Widget, expand, fill bool, padding int) {
	if secondary := b.getSecondary(); secondary != nil {
		secondary.PackEnd(w, expand, fill, padding)
	}
}

// Sets the way child is packed into box .
// Parameters:
// 	child	the Widget of the child to set
// 	expand	the new value of the expand child property
// 	fill	the new value of the fill child property
// 	padding	the new value of the padding child property
// 	packType	the new value of the pack-type child property
func (b *CButtonBox) SetChildPacking(child Widget, expand bool, fill bool, padding int, packType PackType) {
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

func (b *CButtonBox) Add(w Widget) {
	if primary := b.getPrimary(); primary != nil {
		primary.PackStart(w, true, true, 0)
	}
}

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

func (b *CButtonBox) getPrimary() (box Box) {
	children := b.GetChildren()
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
	children := b.GetChildren()
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

func (b *CButtonBox) Resize() enums.EventFlag {
	return b.CBox.Resize()
}

func (b *CButtonBox) Invalidate() enums.EventFlag {
	return b.CBox.Invalidate()
}

func (b *CButtonBox) draw(data []interface{}, argv ...interface{}) enums.EventFlag {
	if surface, ok := argv[1].(*memphis.CSurface); ok {
		b.Lock()
		defer b.Unlock()
		alloc := b.GetAllocation()
		if !b.IsVisible() || alloc.W <= 0 || alloc.H <= 0 {
			b.LogTrace("not visible, zero width or zero height")
			return enums.EVENT_PASS
		}
		debug, _ := b.GetBoolProperty(cdk.PropertyDebug)
		debugChildren, _ := b.GetBoolProperty(PropertyDebugChildren)
		orientation := b.GetOrientation()
		children := b.getBoxChildren()
		surface.Fill(b.GetTheme())
		for _, child := range children {
			if child.widget.IsVisible() {
				child.widget.Draw()
				if childSurface, err := memphis.GetSurface(child.widget.ObjectID()); err != nil {
					child.widget.LogErr(err)
				} else {
					if debugChildren && orientation == enums.ORIENTATION_VERTICAL {
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
		if debug && orientation == enums.ORIENTATION_VERTICAL {
			surface.DebugBox(paint.ColorPink, b.ObjectInfo())
		} else if debug {
			surface.DebugBox(paint.ColorPurple, b.ObjectInfo())
		}
	}
	return enums.EVENT_PASS
}

// How to lay out the buttons in the box. Possible values are: default, spread, edge, start and end.
// Flags: Read / Write
// Default value: GTK_BUTTONBOX_DEFAULT_STYLE
const PropertyLayoutStyle cdk.Property = "layout-style"

const ButtonBoxDrawHandle = "button-box-draw-handler"
