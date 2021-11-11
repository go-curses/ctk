package ctk

type StyleSheetMedia struct {
	Conditions string
	Rules      []*StyleSheetRule
}

func (m StyleSheetMedia) String() string {
	s := "@media "
	s += m.Conditions
	s += " {"
	for _, r := range m.Rules {
		s += r.String()
	}
	s += "}"
	return s
}
