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
	"github.com/iancoleman/strcase"

	cstrings "github.com/go-curses/cdk/lib/strings"
)

func ProcessSignalDetails(src *GtkSource, s *goquery.Selection) {
	rs2 := s.Find("div.refsect2")
	rs2.Each(func(i int, selection *goquery.Selection) {

		// get the name
		text := selection.Find("h3 code.literal").Text()
		text = strings.TrimSpace(text)
		text = rxStripFancyQuotes.ReplaceAllString(text, "")
		// construct new signal object
		signal := NewSignal(strcase.ToCamel(text), text)
		// find the user func parameter notes
		notes := make(map[string]string, 0)
		trs := selection.Find("div.informaltable table tbody tr")
		trs.Each(func(i int, selection *goquery.Selection) {
			pName := selection.Find("td.parameter_name p").Text()
			pName = strings.TrimSpace(pName)
			if pName == "user_data" {
				return
			}
			notes[pName] = selection.Find("td.parameter_description p").Text()
			notes[pName] = strings.Replace(notes[pName], "\n", " ", -1)
			annotation := selection.Find("td.parameter_annotation p").Text()
			if annotation != "&nbsp;" && annotation != " " && annotation != "" {
				notes[pName] += " (" + annotation + ")"
			}
		})
		// decode the code block describing the user function signature
		plText := selection.Find("pre.programlisting").Text()
		plTextArgs := rxUserFnArgs.FindStringSubmatch(plText)
		if len(plTextArgs) > 1 {
			v := strings.Replace(plTextArgs[1], "\n", "", -1)
			args := strings.Split(v, ",")
			for _, arg := range args {
				m := rxUserFnArg.FindStringSubmatch(arg)
				if len(m) == 3 {
					if m[1] == "Gdk"+src.Name {
						continue
					} else if m[1] == "Gtk"+src.Name {
						continue
					}
					nt := NewNamedType(m[2], NewType(src.PackageName, m[1]))
					if v, ok := notes[nt.Name]; ok {
						nt.Note = v
					}
					if nt.Name == "user_data" {
						continue
					}
					signal.UserFnArgs = append(signal.UserFnArgs, nt)
				}
			}
		}
		warning := selection.Find("div.warning").ChildrenFiltered("p")
		if warning != nil && warning.Length() > 0 {
			if strings.Contains(warning.Text(), "is deprecated and should not be used") || strings.Contains(warning.Text(), "has been deprecated since") {
				if !src.Context.Bool("include-deprecated") {
					return
				}
			}
			signal.Docs += "WARNING:"
			warning.Each(func(i int, warnSel *goquery.Selection) {
				text := strings.Replace(warnSel.Text(), "\n", " ", -1)
				signal.Docs += "\n"
				signal.Docs += "\t" + RewriteGtkThingsToCtkThings(src.Name, text)
			})
			signal.Docs += "\n"
		}
		docStr := ""
		ps := selection.ChildrenFiltered("p")
		ps.Each(func(i int, docSel *goquery.Selection) {
			text := docSel.Text()
			if rxTagLine.MatchString(text) {
				m := rxTagLine.FindStringSubmatch(text)
				if len(m) >= 3 {
					if m[1] == "Flags" || m[1] == "Since" {
						return
					}
				}
			}
			if len(text) > 0 {
				if len(docStr) > 0 {
					docStr += "\n"
				}
				docStr += strings.Replace(text, "\n", " ", -1)
			}
		})
		if len(docStr) > 0 {
			docStr = cstrings.BasicWordWrap(docStr, 76)
		}
		signal.Docs += RewriteGtkThingsToCtkThings(src.Name, docStr)
		src.Signals = append(src.Signals, signal)
	})
}