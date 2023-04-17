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
)

type GtkProperty struct {
	Name    string
	Tag     string
	Decl    string
	Write   bool
	Type    GtkType
	Default string
	Docs    string
}

func (s GtkProperty) Registration() string {
	return fmt.Sprintf("\"%s\", ", s.Tag)
}

func (s GtkProperty) String() string {
	docs := ""
	for _, line := range strings.Split(s.Docs, "\n") {
		if strings.HasPrefix(line, "Since:") {
			continue
		}
		docs += "// " + line + "\n"
	}
	return fmt.Sprintf(
		"// \"%v\" %v\n",
		s.Tag,
		s.Type,
	) + docs
}