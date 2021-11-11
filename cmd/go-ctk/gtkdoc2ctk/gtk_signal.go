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
