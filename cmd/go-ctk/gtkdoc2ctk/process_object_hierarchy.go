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
	log "github.com/go-curses/cdk/log"
)

func ProcessObjectHierarchy(src *GtkSource, s *goquery.Selection) {
	pre := s.Find("pre.screen")
	text, err := pre.Html()
	if err != nil {
		log.Error(err)
		return
	}
	text = rxStripLineArt.ReplaceAllString(text, "")
	text = rxStripTags.ReplaceAllString(text, "")
	lines := strings.Split(text, "\n")
	last := ""
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) > 0 {
			if line == "GObject" {
				line = "Object"
			} else if line == "GInterface" {
				line = "CInterface"
			} else if line == "GInitiallyUnowned" {
				src.Hierarchy = []string{}
				continue
			} else if strings.HasPrefix(line, "Gdk") {
				line = line[3:]
			} else if strings.HasPrefix(line, "Gtk") {
				line = line[3:]
			}
			src.Hierarchy = append(src.Hierarchy, line)
			if src.Parent == "" && line == src.Name {
				src.Parent = last
			}
			last = line
			if !cstrings.StringSliceHasValue(src.Hierarchy, line) {
				src.Hierarchy = append(src.Hierarchy, line)
			}
		}
	}
}