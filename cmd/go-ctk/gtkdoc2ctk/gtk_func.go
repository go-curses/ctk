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

type GtkFunc struct {
	Name string
	Docs string
	Argv []*GtkNamedType
	Retv GtkType
	Body string
}

func (f *GtkFunc) String() string {
	if f == nil {
		return ""
	}
	argv := ""
	for _, arg := range f.Argv {
		if len(argv) > 0 {
			argv += ", "
		}
		argv += arg.String()
	}
	retv := ""
	if len(f.Retv.GoName) > 0 {
		retv = " (value " + f.Retv.GoName + ")"
	}
	return fmt.Sprintf(
		"%v(%v)%v",
		f.Name,
		argv,
		retv,
	)
}

func (f *GtkFunc) InterfaceString() string {
	if f == nil {
		return ""
	}
	docs := ""
	for _, line := range strings.Split(f.Docs, "\n") {
		docs += "\t" + line + "\n"
	}
	return fmt.Sprintf(
		"%v\t%v",
		docs,
		f.String(),
	)
}