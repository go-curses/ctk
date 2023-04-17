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

type GtkSignal struct {
	Name       string
	Tag        string
	Docs       string
	UserFnArgs []*GtkNamedType
}

func NewSignal(name, tag string) *GtkSignal {
	return &GtkSignal{
		Name:       name,
		Tag:        tag,
		UserFnArgs: make([]*GtkNamedType, 0),
	}
}

func (s GtkSignal) String() string {
	docs := ""
	if len(s.Docs) > 0 {
		for _, line := range strings.Split(s.Docs, "\n") {
			docs += "// " + line + "\n"
		}
	}
	if len(s.UserFnArgs) > 0 {
		docs += "// Listener function arguments:\n"
		for _, arg := range s.UserFnArgs {
			docs += "// \t" + arg.String() + "\n"
		}
	}
	return docs + fmt.Sprintf(
		"const Signal%v cdk.Signal = \"%v\"\n",
		s.Name,
		s.Tag,
	)
}