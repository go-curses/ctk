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
	"strings"

	"github.com/PuerkitoBio/goquery"

	cstrings "github.com/go-curses/cdk/lib/strings"
)

func ProcessDescription(src *GtkSource, s *goquery.Selection) {
	ps := s.Find("p")
	ps.Each(func(i int, pSel *goquery.Selection) {
		desc := pSel.Text()
		if cstrings.IsEmpty(desc) {
			return
		}
		if len(src.Description) > 0 {
			src.Description += "\n"
		}
		desc = strings.Replace(desc, "Gdk"+src.Name, src.Name, -1)
		desc = strings.Replace(desc, "Gtk"+src.Name, src.Name, -1)
		if src.PackageName == "ctk" {
			desc = strings.Replace(desc, "Gdk", "", -1)
			desc = strings.Replace(desc, "Gtk", "", -1)
		} else {
			desc = strings.Replace(desc, "Gdk", "ctk.", -1)
			desc = strings.Replace(desc, "Gtk", "ctk.", -1)
		}
		src.Description += desc
	})
	if !cstrings.IsEmpty(src.Description) {
		src.Description = cstrings.BasicWordWrap(src.Description, 76)
		desc := ""
		for _, line := range strings.Split(src.Description, "\n") {
			if len(desc) > 0 {
				desc += "\n"
			}
			desc += "// " + line
		}
		if strings.TrimSpace(desc) != "//" {
			desc = RewriteGtkThingsToCtkThings(src.Name, desc)
			src.Description = desc
		}
	}
}