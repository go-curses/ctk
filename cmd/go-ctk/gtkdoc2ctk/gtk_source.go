package gtkdoc2ctk

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
)

type GtkSource struct {
	Name        string // CamelCase
	Flat        string // snake_case
	Tag         string // tag-case
	This        string // method this name
	Parent      string
	PackageName string
	Description string
	Hierarchy   []string
	Implements  []string
	Signals     []*GtkSignal
	Properties  []*GtkProperty
	Functions   []*GtkFunc
	Constructor *GtkFunc
	Factories   []*GtkFunc
	Context     *cli.Context
}

func (s GtkSource) ObjectHierarchy() (output string) {
	output = "// " + s.Name + " Hierarchy:\n"
	depth := 0
	found := false
	for _, thing := range s.Hierarchy {
		output += "//\t"
		for i := 0; i < depth; i++ {
			output += "  "
		}
		if depth > 0 {
			output += "+- "
		}
		output += thing + "\n"
		if thing == s.Name {
			found = true
			depth += 1
		} else if !found {
			depth += 1
		}
	}
	if strings.HasSuffix(output, "\n") {
		output = output[:len(output)-1]
	}
	return
}

func (s GtkSource) String() string {
	return fmt.Sprintf(
		"GtkSource={Name:%v;Flat:%v;DocLines:%d;Hierarchy:%v;Signals:%d;Properties:%d;Functions:%d;};",
		s.Name,
		s.Flat,
		len(strings.Split(s.Description, "\n")),
		strings.Join(s.Hierarchy, ","),
		len(s.Signals),
		len(s.Properties),
		len(s.Functions),
	)
}
