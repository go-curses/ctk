package ctk

type StyleSheetRule struct {
	Selector   string
	Properties []*StyleSheetProperty
}

func (r StyleSheetRule) String() string {
	s := r.Selector
	s += " {"
	for _, e := range r.Properties {
		s += e.String()
	}
	s += "}"
	return s
}
