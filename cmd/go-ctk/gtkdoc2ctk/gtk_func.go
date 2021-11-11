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
