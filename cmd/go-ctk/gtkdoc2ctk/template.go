package gtkdoc2ctk

import (
	"fmt"

	"github.com/gobuffalo/plush"
	"github.com/iancoleman/strcase"
	"github.com/urfave/cli/v2"
	"golang.org/x/net/html"
)

func GenerateCtkSource(c *cli.Context, src *GtkSource) (output string, err error) {
	ctx := plush.NewContext()
	ctx.Set("ctx", c)
	ctx.Set("src", src)
	ctx.Set("sprintf", fmt.Sprintf)
	ctx.Set("CamelCase", func(arg string) string {
		return strcase.ToCamel(arg)
	})
	ctx.Set("lowerCamel", func(arg string) string {
		return strcase.ToLowerCamel(arg)
	})
	ctx.Set("snake_case", func(arg string) string {
		return strcase.ToSnake(arg)
	})
	tmpl := CtkSourceTemplate
	if len(src.Hierarchy) > 0 && src.Hierarchy[0] == "CInterface" {
		tmpl = CtkInterfaceTemplate
	}
	if output, err = plush.Render(tmpl, ctx); err == nil {
		output = html.UnescapeString(output)
	} else {
		err = cli.Exit(fmt.Sprintf("error generating Go source code: %v", err), 1)
	}
	return
}
