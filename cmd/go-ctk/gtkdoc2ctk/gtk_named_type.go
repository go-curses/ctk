package gtkdoc2ctk

import (
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
)

type GtkNamedType struct {
	Name  string
	Type  GtkType
	Value string
	Note  string
}

func NewNamedType(name string, t GtkType) *GtkNamedType {
	if strings.HasPrefix(name, "*") {
		name = name[1:]
	}
	return &GtkNamedType{
		Name: name,
		Type: t,
	}
}

func (n GtkNamedType) String() string {
	note := ""
	if len(n.Note) > 0 {
		note = "\t" + n.Note
	}
	return fmt.Sprintf(
		"%v %v%v",
		strcase.ToLowerCamel(n.Name),
		n.Type,
		note,
	)
}
