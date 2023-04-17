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

package gtkdoc2ctk

import (
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
)

type GtkNamedType struct {
	Name  string
	Type  GtkType
	Value string
	Note  string
}

func NewNamedType(name string, t GtkType) *GtkNamedType {
	if strings.HasPrefix(name, "*") {
		name = name[1:]
	}
	return &GtkNamedType{
		Name: name,
		Type: t,
	}
}

func (n GtkNamedType) String() string {
	note := ""
	if len(n.Note) > 0 {
		note = "\t" + n.Note
	}
	return fmt.Sprintf(
		"%v %v%v",
		strcase.ToLowerCamel(n.Name),
		n.Type,
		note,
	)
}