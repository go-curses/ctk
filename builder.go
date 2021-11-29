package ctk

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-curses/ctk/lib/enums"
	"github.com/iancoleman/strcase"

	"github.com/go-curses/cdk"
	cstrings "github.com/go-curses/cdk/lib/strings"
	"github.com/go-curses/cdk/log"
)

type BuilderTranslationFn = func(builder Builder, widget Widget, name, value string) error

var ctkBuilderTranslators map[cdk.TypeTag]BuilderTranslationFn

func init() {
	ctkBuilderTranslators = make(map[cdk.TypeTag]BuilderTranslationFn)
}

func BuilderRegisterConstructor(tag cdk.TypeTag, fn BuilderTranslationFn) {
	if _, ok := ctkBuilderTranslators[tag]; ok {
		log.WarnF("overwriting existing NewObject constructor for %v objects", tag)
	}
	ctkBuilderTranslators[tag] = fn
}

const TypeBuilder cdk.CTypeTag = "ctk-builder"

func init() {
	_ = cdk.TypesManager.AddType(TypeBuilder, nil)
}

type Builder interface {
	cdk.Object

	Init() (already bool)
	BuildableTypeTags() map[string]cdk.TypeTag
	LookupNamedSignalHandler(name string) (fn cdk.SignalListenerFn)
	AddNamedSignalHandler(name string, fn cdk.SignalListenerFn)
	GetWidget(name string) (w interface{})
	GetWidgetsBuiltByType(tag cdk.CTypeTag) (widgets []interface{})
	ParsePacking(packing *CBuilderElement) (expand, fill bool, padding int, packType enums.PackType)
	LoadFromString(raw string) (topElement *CBuilderElement, err error)
	Build(element *CBuilderElement) (newObject interface{})
}

type CBuilder struct {
	cdk.CObject

	handlers  map[string]cdk.SignalListenerFn
	objects   []*CBuilderElement
	buildable map[string]cdk.TypeTag
	built     []interface{}
}

func NewBuilder() (builder Builder) {
	builder = &CBuilder{}
	_ = builder.Init()
	return
}

func (b *CBuilder) Init() (already bool) {
	if b.InitTypeItem(TypeBuilder, b) {
		return true
	}
	b.CObject.Init()
	b.buildable = cdk.TypesManager.GetBuildableInfo()
	b.handlers = make(map[string]cdk.SignalListenerFn)
	b.built = make([]interface{}, 0)
	return
}

func (b *CBuilder) BuildableTypeTags() map[string]cdk.TypeTag {
	return b.buildable
}

func (b *CBuilder) LookupNamedSignalHandler(name string) (fn cdk.SignalListenerFn) {
	fn, _ = b.handlers[name]
	return
}

func (b *CBuilder) AddNamedSignalHandler(name string, fn cdk.SignalListenerFn) {
	if _, ok := b.handlers[name]; ok {
		b.LogWarn("overwriting existing handler: %v", name)
	}
	log.DebugDF(1, "adding handler: %v", name)
	b.handlers[name] = fn
}

func (b *CBuilder) GetWidget(name string) (w interface{}) {
	w = cdk.TypesManager.GetTypeItemByName(name)
	log.DebugDF(1, "GetWidget(%v) = %v", name, w)
	return
}

func (b *CBuilder) GetWidgetsBuiltByType(tag cdk.CTypeTag) (widgets []interface{}) {
	widgets = make([]interface{}, 0)
	for _, object := range b.built {
		if wo, ok := object.(Widget); ok {
			if wo.GetTypeTag().Equals(tag) {
				widgets = append(widgets, object)
			}
		}
	}
	return
}

func (b *CBuilder) ParsePacking(packing *CBuilderElement) (expand, fill bool, padding int, packType enums.PackType) {
	expand, fill = false, true
	padding = 0
	packType = enums.PackStart
	for k, v := range packing.Packing {
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
	return
}

func (b *CBuilder) LoadFromString(raw string) (topElement *CBuilderElement, err error) {
	b.LogDebug("known buildable types: %v", b.buildable)
	r := strings.NewReader(raw)
	parser := xml.NewDecoder(r)
	var n BuilderNode
	if err = parser.Decode(&n); err != nil {
		return nil, err
	}
	topElement = b.walkElements(n)
	b.LogDebug("see report:\n[report]\n%v[/report]", b.report(0, topElement))
	_ = b.Build(topElement)
	return topElement, nil
}

func (b *CBuilder) Build(element *CBuilderElement) (newObject interface{}) {
	switch element.TagName {
	case "requires":
	case "ctk-builder", "interface":
		if len(element.Children) > 0 {
			for _, child := range element.Children {
				newObject = b.Build(child)
				child.Instance = newObject
			}
		}
	case "object":
		if class, ok := element.Attributes["class"]; ok {
			if tt, ok := b.buildable[class]; ok {
				ct, _ := cdk.TypesManager.GetType(tt)
				newObject = ct.New()
				b.built = append(b.built, newObject)
				var newBuildable Buildable
				var ok bool
				if newBuildable, ok = newObject.(Buildable); !ok {
					b.LogError("new object is not a Buildable type: %v (%T)", newObject, newObject)
					newObject = nil
					break
				}
				element.Instance = newObject
				newBuildable.Show()
				if err := newBuildable.Build(b, element); err != nil {
					b.LogErr(err)
				}
			} else {
				b.LogError("ctk class not implemented: %v", class)
			}
		}
	default:
		b.LogError("unexpected element type: %v", element.TagName)
	}
	return
}

func (b *CBuilder) report(depth int, element *CBuilderElement) (output string) {
	pad := ""
	for i := 0; i < depth; i++ {
		pad += "\t"
	}
	nChildren := len(element.Children)
	properties := ""
	for k, v := range element.Properties {
		if len(properties) > 0 {
			properties += " "
		}
		properties += k + "=\"" + v + "\""
	}
	for k, v := range element.Attributes {
		if len(properties) > 0 {
			properties += " "
		}
		properties += k + "=\"" + v + "\""
	}
	switch element.TagName {
	case "requires":
		output += fmt.Sprintf("%v<requires %v/>\n", pad, properties)
	case "ctk-builder", "interface":
		output += fmt.Sprintf("%v<ctk-builder>", pad)
		if nChildren > 0 {
			output += "\n"
			for _, child := range element.Children {
				output += b.report(depth+1, child)
			}
		}
		output += fmt.Sprintf("%v</ctk-builder>\n", pad)
	case "object":
		if class, ok := element.Attributes["class"]; ok {
			if name, ok := element.Attributes["id"]; ok {
				if len(properties) > 0 {
					properties = fmt.Sprintf(`name="%v" %v`, name, properties)
				} else {
					properties = fmt.Sprintf(`name="%v"`, name)
				}
			}
			if nChildren > 0 {
				output += fmt.Sprintf("%v<%v %v>\n", pad, class, properties)
				for _, child := range element.Children {
					output += b.report(depth+1, child)
				}
				output += fmt.Sprintf("%v</%v>\n", pad, class)
			} else {
				output += fmt.Sprintf("%v<%v %v/>\n", pad, class, properties)
			}
		}
	default:
		output += fmt.Sprintf("%v<!-- unexpected tag: %v -->\n", pad, element.TagName)
	}
	return
}

func (b *CBuilder) walkElements(n BuilderNode, argv ...interface{}) (be *CBuilderElement) {
	if n.XMLName.Local == "object" {
		return b.walkObject(n)
	}
	be = newBuilderElement(n.XMLName.Local, b)
	be.Attributes = b.parseTagAttributes(n.Attrs)
	be.Content = string(n.Content)
	for _, cn := range n.Nodes {
		switch cn.XMLName.Local {
		case "object":
			if cbe := b.walkObject(cn); cbe != nil {
				be.Children = append(be.Children, cbe)
			}
		default:
			if cbe := b.walkElements(cn, argv...); cbe != nil {
				be.Children = append(be.Children, cbe)
			}
		}
	}
	return
}

func (b *CBuilder) walkObject(n BuilderNode) (be *CBuilderElement) {
	tagStr := fmt.Sprintf("<%v %v>%v</%v>", n.XMLName.Local, n.Attrs, string(n.Content), n.XMLName.Local)
	be = newBuilderElement(n.XMLName.Local, b)
	be.Attributes = b.parseTagAttributes(n.Attrs)
	be.Content = string(n.Content)
	for _, cn := range n.Nodes {
		switch cn.XMLName.Local {
		case "property":
			attrs := b.parseTagAttributes(cn.Attrs)
			if v, ok := attrs["name"]; ok {
				vs := strcase.ToDelimited(v, '-')
				be.Properties[vs] = string(cn.Content)
			} else {
				b.LogError("missing name attribute on property tag: %v", tagStr)
			}
		case "signal":
			attrs := b.parseTagAttributes(cn.Attrs)
			if name, ok := attrs["name"]; ok {
				if handler, ok := attrs["handler"]; ok {
					be.Signals[name] = handler
				}
			} else {
				b.LogError("missing name attribute on property tag: %v", tagStr)
			}
		case "child":
			b.walkObjectChild(cn, be)
		default:
			b.LogError("ignoring unexpected tag: %v", cn.XMLName.Local)
		}
	}
	return
}

func (b *CBuilder) walkObjectChild(n BuilderNode, parentElement *CBuilderElement) {
	var object *CBuilderElement
	packing := make(map[string]string)
	for _, cn := range n.Nodes {
		switch cn.XMLName.Local {
		case "object":
			if cbe := b.walkObject(cn); cbe != nil {
				object = cbe
				parentElement.Children = append(
					parentElement.Children,
					object,
				)
			}
		case "packing":
			for _, ccn := range cn.Nodes {
				attrs := b.parseTagAttributes(ccn.Attrs)
				if v, ok := attrs["name"]; ok {
					packing[v] = string(ccn.Content)
				}
			}
		case "placeholder":
			return
		default:
			b.LogError("ignoring unexpected tag: %v", cn.XMLName.Local)
		}
	}
	if object != nil {
		for k, v := range packing {
			object.Packing[k] = v
		}
	} else {
		b.LogError("object not found in child tag children")
	}
	return
}

func (b *CBuilder) parseTagAttributes(attrs []xml.Attr) (properties map[string]string) {
	properties = make(map[string]string)
	for _, attr := range attrs {
		properties[attr.Name.Local] = attr.Value
	}
	return
}

const PropertyID cdk.Property = "id"
const PropertyHandler cdk.Property = "handler"
const PropertySwapped cdk.Property = "swapped"
