package ctk

import (
	"fmt"
	"strconv"

	"github.com/go-curses/cdk"
	cstrings "github.com/go-curses/cdk/lib/strings"
	"github.com/go-curses/ctk/lib/enums"
)

type BuilderElement interface {
	String() string
	ApplyProperties()
	ApplyProperty(k, v string) (set bool)
}

type CBuilderElement struct {
	TagName    string
	Content    string
	Builder    Builder
	Instance   interface{}
	Attributes map[string]string
	Properties map[string]string
	Signals    map[string]string
	Packing    map[string]string
	Children   []*CBuilderElement
}

func newBuilderElement(tagName string, builder Builder) *CBuilderElement {
	b := new(CBuilderElement)
	b.init()
	b.Builder = builder
	b.TagName = tagName
	return b
}

func (b *CBuilderElement) init() {
	b.Content = ""
	b.Instance = nil
	b.Attributes = make(map[string]string)
	b.Properties = make(map[string]string)
	b.Signals = make(map[string]string)
	b.Packing = make(map[string]string)
	b.Children = make([]*CBuilderElement, 0)
}

func (b *CBuilderElement) String() string {
	if len(b.Children) == 0 {
		return fmt.Sprintf("<%v %v>%v</%v>", b.TagName, b.Attributes, b.Content, b.TagName)
	}
	children := ""
	for _, child := range b.Children {
		children += child.String() + "\n"
	}
	return fmt.Sprintf("<%v %v>\n%v</%v>", b.TagName, b.Attributes, children, b.TagName)
}

func (b *CBuilderElement) ApplySignals() {
	for k, v := range b.Signals {
		b.ApplySignal(k, v)
	}
}

func (b *CBuilderElement) ApplySignal(k, v string) {
	if buildable, ok := b.Instance.(Buildable); ok {
		ks := cdk.Signal(k)
		if fn := b.Builder.LookupNamedSignalHandler(v); fn != nil {
			buildable.Connect(ks, v, fn)
		} else {
			b.Builder.LogError("missing named signal handler: %v", v)
		}
	}
}

func (b *CBuilderElement) ApplyProperties() {
	for k, v := range b.Properties {
		b.ApplyProperty(k, v)
	}
}

func (b *CBuilderElement) ApplyProperty(k, v string) (set bool) {
	if buildableWidget, ok := b.Instance.(Widget); ok {
		tt := buildableWidget.GetTypeTag()
		if fn, ok := ctkBuilderTranslators[tt]; ok {
			if err := fn(b.Builder, buildableWidget, k, v); err != nil {
				if err != ErrFallthrough {
					buildableWidget.LogError("%v property translator error: %v", tt, err)
					return false
				}
			} else {
				return true
			}
		}
	}
	if buildableWidget, ok := b.Instance.(Buildable); ok {
		switch k {
		case "sensitive":
			buildableWidget.SetSensitive(cstrings.IsTrue(v))
		case "has-focus", "has_focus", "is-focus", "is_focus":
			if cstrings.IsTrue(v) {
				buildableWidget.SetFlags(enums.HAS_FOCUS)
			} else {
				buildableWidget.UnsetFlags(enums.HAS_FOCUS)
			}
		case "can-focus", "can_focus":
			if cstrings.IsTrue(v) {
				buildableWidget.SetFlags(enums.CAN_FOCUS)
			} else {
				buildableWidget.UnsetFlags(enums.CAN_FOCUS)
			}
		case "can-default", "can_default":
			if cstrings.IsTrue(v) {
				buildableWidget.SetFlags(enums.CAN_DEFAULT)
			} else {
				buildableWidget.UnsetFlags(enums.CAN_DEFAULT)
			}
		case "has-default", "has_default":
			if cstrings.IsTrue(v) {
				buildableWidget.SetFlags(enums.HAS_DEFAULT)
			} else {
				buildableWidget.UnsetFlags(enums.HAS_DEFAULT)
			}
		case "app-paintable":
			if cstrings.IsTrue(v) {
				buildableWidget.SetFlags(enums.APP_PAINTABLE)
			} else {
				buildableWidget.UnsetFlags(enums.APP_PAINTABLE)
			}
		case "width_request", "width-request":
			if bw, ok := b.Instance.(Widget); ok {
				w, _ := strconv.Atoi(v)
				h, _ := bw.GetIntProperty(PropertyHeightRequest)
				bw.SetSizeRequest(w, h)
			}
		case "height_request", "height-request":
			if bw, ok := b.Instance.(Widget); ok {
				w, _ := bw.GetIntProperty(PropertyWidthRequest)
				h, _ := strconv.Atoi(v)
				bw.SetSizeRequest(w, h)
			}
		default:
			kp := cdk.Property(k)
			if err := buildableWidget.SetPropertyFromString(kp, v); err != nil {
				buildableWidget.LogErr(err)
				return false
			}
		}
	}
	return true
}
