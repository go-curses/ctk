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

func newStyleSheetSelectorFromPath(path string) (selector *StyleSheetSelector) {
	selector = &StyleSheetSelector{
		Name:  "",
		Type:  "",
		Class: "",
		State: "normal",
	}
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

func (s StyleSheetSelector) Match(selector *StyleSheetSelector) (match bool) {
	var wClass, wType, wName bool
	var mClass, mType, mName bool
	if len(selector.Class) > 0 {
		wClass = true
		if len(s.Class) > 0 {
			// has classes to compare
			mClass = s.Class == selector.Class
		} else if selector.Class == "normal" {
			// normal class can be omitted
			mClass = true
		} else {
			return false
		}
	}
	if len(selector.Type) > 0 {
		wType = true
		if len(s.Type) > 0 {
			// has types to compare
			mType = s.Type == selector.Type
		} else {
			return false
		}
	}
	if len(selector.Name) > 0 {
		wName = true
		if len(s.Name) > 0 {
			// has names to compare
			mName = s.Name == selector.Name
		} else {
			return false
		}
	}
	return (!wClass || (wClass && mClass)) && (!wType || (wType && mType)) && (!wName || (wName && mName))
}

var (
	rxSelectorName  = regexp.MustCompile(`#([a-zA-Z][-_a-zA-Z0-9]+)`)
	rxSelectorClass = regexp.MustCompile(`\.([a-zA-Z][-_a-zA-Z0-9]+)`)
	rxSelectorState = regexp.MustCompile(`:([a-zA-Z][-_a-zA-Z0-9]+)`)
)
