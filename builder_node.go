package ctk

import (
	"encoding/xml"
)

type BuilderNode struct {
	XMLName xml.Name
	Attrs   []xml.Attr    `xml:"-"`
	Content []byte        `xml:",innerxml"`
	Nodes   []BuilderNode `xml:",any"`
}

func (n *BuilderNode) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	n.Attrs = start.Attr
	type node BuilderNode

	return d.DecodeElement((*node)(n), &start)
}
