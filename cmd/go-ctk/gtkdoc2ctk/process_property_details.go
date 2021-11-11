package gtkdoc2ctk

import (
	"fmt"
	"html"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/iancoleman/strcase"

	cstrings "github.com/go-curses/cdk/lib/strings"
)

func ProcessPropertyDetails(src *GtkSource, s *goquery.Selection) {
	rs2 := s.Find("div.refsect2")
	rs2.Each(func(i int, selection *goquery.Selection) {
		text := selection.Find("h3 code.literal").Text()
		text = strings.TrimSpace(text)
		text = rxStripFancyQuotes.ReplaceAllString(text, "")
		n, tag := strcase.ToCamel(text), text
		text = selection.Find("pre.programlisting span.type").Text()
		text = strings.TrimSpace(text)
		t := NewType(src.PackageName, text)
		// Flags: Read / Write
		write := false
		def := "nil"
		docStr := ""
		var tagLines []string
		ps := selection.ChildrenFiltered("p")
		ps.Each(func(i int, selection *goquery.Selection) {
			st := selection.Text()
			if rxTagLine.MatchString(st) {
				m := rxTagLine.FindStringSubmatch(st)
				if len(m) >= 3 {
					switch m[1] {
					case "Flags":
						if strings.Contains(m[2], "Write") {
							write = true
						}
					case "Default value":
						def = st[15:]
					case "Since":
						return
					default:
					}
					tagLines = append(tagLines, strings.TrimSpace(st))
				}
			} else {
				if len(docStr) > 0 {
					docStr += "\n"
				}
				docStr += strings.TrimSpace(st)
			}
		})
		docStr = cstrings.BasicWordWrap(docStr, 76)
		docStr += "\n"
		docStr += strings.Join(tagLines, "\n")
		def = html.UnescapeString(def)
		if strings.HasPrefix(def, "\"") {
			def = strings.Replace(def, "\"", "", -1)
		}
		switch def {
		case "FALSE":
			def = "false"
		case "TRUE":
			def = "true"
		case "":
			def = "\"\""
		default:
			if !rxIsNumbers.MatchString(def) {
				def = "nil"
			}
		}
		warning := selection.Find("div.warning").ChildrenFiltered("p")
		if warning != nil && warning.Length() > 0 {
			if strings.Contains(warning.Text(), "is deprecated and should not be used") || strings.Contains(warning.Text(), "has been deprecated since") {
				if !src.Context.Bool("include-deprecated") {
					return
				}
			}
			if len(docStr) > 0 {
				if docStr[len(docStr)-1] != '\n' {
					docStr += "\n\n"
				} else {
					docStr += "\n"
				}
			}
			docStr += "WARNING:"
			warning.Each(func(i int, warnSel *goquery.Selection) {
				text := strings.Replace(warnSel.Text(), "\n", " ", -1)
				docStr += "\n"
				docStr += "\t" + text
			})
		}
		docStr = RewriteGtkThingsToCtkThings(src.Name, docStr)
		decl := ""
		for _, line := range strings.Split(docStr, "\n") {
			if len(decl) > 0 {
				decl += "\n"
			}
			decl += "// " + line
		}
		if len(decl) > 0 {
			decl += "\n"
		}
		decl += fmt.Sprintf("const Property%s cdk.Property = \"%s\"", n, tag)
		prop := &GtkProperty{
			Name:    n,
			Tag:     tag,
			Decl:    decl,
			Type:    t,
			Write:   write,
			Default: def,
			Docs:    docStr,
		}
		src.Properties = append(src.Properties, prop)
	})
}
