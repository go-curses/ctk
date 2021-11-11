package ctk

import (
	"regexp"
)

type StyleSheetSelector struct {
	Name  string
	Type  string
	Class string
	State string
}

func (s StyleSheetSelector) String() string {
	str := ""
	if len(s.Type) > 0 {
		str += s.Type
	}
	if len(s.Name) > 0 {
		str += "#" + s.Name
	}
	if len(s.Class) > 0 {
		str += "." + s.Class
	}
	if len(s.State) > 0 {
		str += ":" + s.State
	}
	return str
}

func (s StyleSheetSelector) Match(selector *StyleSheetSelector) (score int) {
	if len(s.Type) > 0 && len(selector.Type) > 0 {
		if s.Type == selector.Type {
			score += 1
		}
	}
	if len(s.Name) > 0 && len(selector.Name) > 0 {
		if s.Name == selector.Name {
			score += 1
		} else {
			score -= 1
		}
	}
	if len(s.Class) > 0 && len(selector.Class) > 0 {
		if s.Class == selector.Class {
			score += 1
		}
	}
	if len(s.State) > 0 && len(selector.State) > 0 {
		if s.State == selector.State {
			score += 1
		} else {
			score -= 1
		}
	} else {
		if len(s.State) > 0 || len(selector.State) > 0 {
			score -= 1
		}
	}
	return
}

var (
	rxSelectorName  = regexp.MustCompile(`#([a-zA-Z][-_a-zA-Z0-9]+)`)
	rxSelectorClass = regexp.MustCompile(`\.([a-zA-Z][-_a-zA-Z0-9]+)`)
	rxSelectorState = regexp.MustCompile(`:([a-zA-Z][-_a-zA-Z0-9]+)`)
)

func ParseSelector(path string) (selector *StyleSheetSelector) {
	selector = &StyleSheetSelector{}
	altered := path
	if rxSelectorName.MatchString(path) {
		m := rxSelectorName.FindStringSubmatch(path)
		if len(m) > 1 {
			selector.Name = m[1]
			altered = rxSelectorName.ReplaceAllString(altered, "")
		}
	}
	if rxSelectorClass.MatchString(path) {
		m := rxSelectorClass.FindStringSubmatch(path)
		if len(m) > 1 {
			selector.Class = m[1]
			altered = rxSelectorClass.ReplaceAllString(altered, "")
		}
	}
	if rxSelectorState.MatchString(path) {
		m := rxSelectorState.FindStringSubmatch(path)
		if len(m) > 1 {
			selector.State = m[1]
			altered = rxSelectorState.ReplaceAllString(altered, "")
		}
	}
	if len(altered) > 0 {
		selector.Type = altered
	}
	return
}
