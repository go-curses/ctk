package gtkdoc2ctk

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/iancoleman/strcase"

	cstrings "github.com/go-curses/cdk/lib/strings"
	log "github.com/go-curses/cdk/log"
)

func ProcessFunctionDetails(src *GtkSource, s *goquery.Selection) {
	blobs := s.Find("div.refsect2")
	blobs.Each(func(i int, selection *goquery.Selection) {
		fn := new(GtkFunc)
		// method name
		h3 := selection.Find("h3").Text()
		h3 = strings.Replace(h3, "()", "", 1)
		h3 = strings.Replace(h3, src.Flat, "", 1)
		h3 = strings.TrimSpace(h3)
		fnName := h3
		if strings.HasPrefix(fnName, "gdk_") {
			fnName = strcase.ToCamel(fnName[4:])
		} else if strings.HasPrefix(fnName, "gtk_") {
			fnName = strcase.ToCamel(fnName[4:])
		}
		fn.Name = strcase.ToCamel(fnName)
		// documentation
		fn.Docs = ""
		docStr := ""
		warning := selection.Find("div.warning").ChildrenFiltered("p")
		if warning != nil && warning.Length() > 0 {
			if strings.Contains(warning.Text(), "is deprecated and should not be used") || strings.Contains(warning.Text(), "has been deprecated since") {
				if !src.Context.Bool("include-deprecated") {
					return
				}
			}
			docStr += "WARNING:"
			warning.Each(func(i int, warnSel *goquery.Selection) {
				text := strings.Replace(warnSel.Text(), "\n", " ", -1)
				docStr += "\n"
				docStr += "\t" + text
			})
		}
		var tagLines []string
		ps := selection.ChildrenFiltered("p")
		ps.Each(func(i int, docSel *goquery.Selection) {
			text := docSel.Text()
			if rxTagLine.MatchString(text) {
				m := rxTagLine.FindStringSubmatch(text)
				if len(m) >= 3 {
					switch m[1] {
					case "Since":
						return
					default:
					}
					tagLines = append(tagLines, strings.TrimSpace(text))
				}
			} else {
				if len(docStr) > 0 {
					docStr += "\n"
				}
				docStr += text
			}
		})
		// return value type
		rvText := selection.Find("pre.programlisting span.returnvalue").Text()
		fn.Retv = TranslateGtkType(src.PackageName, rvText)
		// parameter and return value docs?
		var paramLines []string
		var rvLines []string
		rs3 := selection.Find("div.refsect3")
		rs3.Each(func(i int, rsSel *goquery.Selection) {
			h4 := strings.TrimSpace(rsSel.Find("h4").Text())
			log.DebugF("func rs3 h4: %v", h4)
			switch h4 {
			case "Parameters":
				trs := rsSel.Find("div.informaltable > table > tbody > tr")
				if trs.Length() == 0 {
					return
				}
				tmpStr := ""
				trs.Each(func(i int, trSel *goquery.Selection) {
					n := trSel.Find("td.parameter_name > p").Text()
					n = strcase.ToLowerCamel(n)
					t := trSel.Find("td.parameter_description span.type").Text()
					t = strings.TrimSpace(t)
					if t == "Gdk"+src.Name || t == "Gtk"+src.Name {
						return
					}
					d := trSel.Find("td.parameter_description > p").Text()
					tmpStr += "\n"
					tmpStr += "\t" + n + "\t" + d
				})
				if len(tmpStr) > 0 {
					paramLines = append(paramLines, tmpStr)
				}
			case "Returns":
				ps := rsSel.ChildrenFiltered("p")
				ps.Each(func(i int, psSel *goquery.Selection) {
					lines := strings.TrimSpace(strings.Replace(psSel.Text(), "\n", " ", -1))
					wrapped := cstrings.BasicWordWrap(lines, 65)
					for _, line := range strings.Split(wrapped, "\n") {
						rvLines = append(rvLines, "\t"+line)
					}
				})
			}
		})
		// function signature, arguments
		prePl := selection.ChildrenFiltered("pre.programlisting")
		emp := prePl.ChildrenFiltered("em.parameter")
		ems := emp.ChildrenFiltered("code")
		fn.Argv = make([]*GtkNamedType, 0)
		ems.Each(func(i int, emSel *goquery.Selection) {
			text := emSel.Text()
			if i == 0 && (strings.HasPrefix(text, "Gdk"+src.Name) || strings.HasPrefix(text, "Gtk"+src.Name)) || text == "void" {
				return
			}
			if text == "..." {
				text = "...interface{} argv"
			}
			if strings.HasPrefix(text, "const ") {
				text = text[6:]
			}
			parts := strings.Split(text, " ")
			if len(parts) != 2 {
				log.ErrorF("too many parts in argument signature: \"%v\", parts=%v\n", emSel.Text(), parts)
				return
			}
			nv := NewNamedType(TranslateNamedVariable(src.PackageName, parts[0], parts[1]))
			fn.Argv = append(fn.Argv, nv)
		})
		if fn.Name == "New" {
			fn.Name = "New" + src.Name
			fn.Retv.GoName = "*C" + src.Name
			src.Constructor = fn
		} else if strings.HasPrefix(fn.Name, "New") {
			keep := fn.Name[3:]
			fn.Name = "New" + src.Name + keep
			fn.Retv.GoName = "*C" + src.Name
			src.Factories = append(src.Factories, fn)
		} else {
			src.Functions = append(src.Functions, fn)
		}
		// update the documentation string
		docStr = cstrings.BasicWordWrap(docStr, 76)
		if len(tagLines) > 0 {
			docStr += "\n"
			docStr += strings.Join(tagLines, "\n")
		}
		if len(paramLines) > 0 {
			docStr += "\nParameters:"
			docStr += strings.Join(paramLines, "\n")
		}
		if len(rvLines) > 0 {
			docStr += "\nReturns:\n"
			docStr += strings.Join(rvLines, "\n")
		}
		if !cstrings.IsEmpty(docStr) {
			for _, line := range strings.Split(docStr, "\n") {
				if len(fn.Docs) > 0 {
					fn.Docs += "\n"
				}
				fn.Docs += "// " + line
			}
			fn.Docs = RewriteGtkThingsToCtkThings(src.Name, fn.Docs)
		}
	})
}
