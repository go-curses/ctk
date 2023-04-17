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