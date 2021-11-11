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
