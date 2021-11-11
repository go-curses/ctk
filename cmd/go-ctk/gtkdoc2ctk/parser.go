package gtkdoc2ctk

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/iancoleman/strcase"
	"github.com/urfave/cli/v2"
)

var (
	rxStripLineArt     = regexp.MustCompile(`(?msi)<span[^>]+?>.+?</span>`)
	rxStripTags        = regexp.MustCompile(`(?msi)<[^>]+?>`)
	rxStripDigitSuffix = regexp.MustCompile(`(?msi)\d+\s*$`)
	rxStripFancyQuotes = regexp.MustCompile(`(?msi)[“”]`)
	rxUserFnArgs       = regexp.MustCompile(`(?msi)void\s*user_function\s*\(\s*(.+?)\s*\)\s*`)
	rxUserFnArg        = regexp.MustCompile(`^\s*(\S+?)\s*(\S+?)\s*$`)
	rxIsNumbers        = regexp.MustCompile(`^\s*-??\d+\.??\d*\s*$`)
	rxTagLine          = regexp.MustCompile(`^\s*([^:]+?): (.+?)\s*$`)
)

func ParseGtk2Doc(snakedName, contents string, c *cli.Context) (source *GtkSource, err error) {
	camelName := strcase.ToCamel(snakedName)
	src := &GtkSource{
		Name:        camelName,
		Flat:        snakedName,
		Tag:         strcase.ToDelimited(snakedName, '-'),
		This:        strings.ToLower(snakedName[:1]),
		PackageName: c.String("package-name"),
		Description: "",
		Implements:  []string{},
		Hierarchy:   []string{},
		Properties:  make([]*GtkProperty, 0),
		Constructor: nil,
		Functions:   make([]*GtkFunc, 0),
		Factories:   make([]*GtkFunc, 0),
		Context:     c,
	}

	body := strings.NewReader(contents)
	var doc *goquery.Document
	if doc, err = goquery.NewDocumentFromReader(body); err != nil {
		return nil, err
	}

	doc.Find("div.refsect1").Each(func(_ int, s *goquery.Selection) {
		a := s.Find("a")
		if name, ok := a.Attr("name"); ok {
			switch name {
			case fmt.Sprintf("Gdk%s.object-hierarchy", src.Name):
				fallthrough
			case fmt.Sprintf("Gtk%s.object-hierarchy", src.Name):
				ProcessObjectHierarchy(src, s)
			case fmt.Sprintf("Gdk%s.implemented-interfaces", src.Name):
				fallthrough
			case fmt.Sprintf("Gtk%s.implemented-interfaces", src.Name):
				ProcessImplementedInterfaces(src, s)
			case fmt.Sprintf("Gdk%s.description", src.Name):
				fallthrough
			case fmt.Sprintf("Gtk%s.description", src.Name):
				ProcessDescription(src, s)
			case fmt.Sprintf("Gdk%s.signal-details", src.Name):
				fallthrough
			case fmt.Sprintf("Gtk%s.signal-details", src.Name):
				ProcessSignalDetails(src, s)
			case fmt.Sprintf("Gdk%s.property-details", src.Name):
				fallthrough
			case fmt.Sprintf("Gtk%s.property-details", src.Name):
				ProcessPropertyDetails(src, s)
			case fmt.Sprintf("Gdk%s.functions_details", src.Name):
				fallthrough
			case fmt.Sprintf("Gtk%s.functions_details", src.Name):
				ProcessFunctionDetails(src, s)
			}
		}
	})

	// sort out the property arguments for the constructor and any factories
	if src.Constructor != nil {
		for _, arg := range src.Constructor.Argv {
			for _, prop := range src.Properties {
				if strcase.ToLowerCamel(prop.Name) == arg.Name {
					arg.Value = prop.Default
					break
				}
			}
			if arg.Value == "nil" {
				switch arg.Type.GoName {
				case "string":
					arg.Value = "\"\""
				case "int":
					arg.Value = "0"
				case "float64":
					arg.Value = "0.0"
				case "bool":
					arg.Value = "false"
				default:
					// leave it nil
				}
			}
		}
	}
	if len(src.Factories) > 0 {
		for _, factory := range src.Factories {
			for _, arg := range factory.Argv {
				for _, prop := range src.Properties {
					if prop.Name == arg.Name {
						arg.Value = prop.Default
						if arg.Value == "" && arg.Type.GoName == "string" {
							arg.Value = "\"\""
						}
						break
					}
				}
			}
		}
	}

	// with all parameters, signals and functions parsed, sort out the
	// function bodies

	for _, fn := range src.Functions {
		fn.Body = ""
		wrote := false
		if len(src.Properties) > 0 {
			for _, prop := range src.Properties {
				switch fn.Name {
				case fmt.Sprintf("Get%s", strcase.ToCamel(prop.Name)):
					getType := strcase.ToCamel(prop.Type.GoName)
					switch getType {
					case "String", "Bool", "Int", "Float64":
					default:
						getType = "Struct"
					}
					fn.Body += "\tvar err error\n"
					fn.Body += "\tif value, err = " + src.This + ".Get" + getType + "Property(Property" + prop.Name + "); err != nil {\n"
					fn.Body += "\t\t" + src.This + ".LogErr(err)\n"
					fn.Body += "\t}\n"
					fn.Body += "\treturn"
					wrote = true
				case fmt.Sprintf("Set%s", strcase.ToCamel(prop.Name)):
					// setter, take argv and assign to property?
					for _, arg := range fn.Argv {
						check0 := strings.Contains(strcase.ToSnake(arg.Name), strcase.ToSnake(prop.Name))
						check1 := strings.Contains(strcase.ToSnake(prop.Name), strcase.ToSnake(arg.Name))
						check2 := arg.Type.GoName == prop.Type.GoName
						if check0 || check1 || check2 {
							fn.Body += "\tif err := " + src.This + ".Set" + strcase.ToCamel(prop.Type.GoLabel) + "Property(Property" + prop.Name + ", " + fn.Argv[0].Name + "); err != nil {\n"
							fn.Body += "\t\t" + src.This + ".LogErr(err)\n"
							fn.Body += "\t}\n"
							wrote = true
							break
						}
					}
				}
				if wrote {
					break
				}
			}
			for _, sig := range src.Signals {
				switch fn.Name {
				case fmt.Sprintf("Emit%v", strcase.ToCamel(sig.Name)):
					fn.Body += "\tif f := " + src.This + ".Emit(Signal" + sig.Name + "); f == cdk.EVENT_STOP {\n"
					fn.Body += "\t\t" + src.This + ".LogTrace(\"Signal" + sig.Name + " was stopped\")\n"
					fn.Body += "\t}\n"
				}
			}
		}
		if !wrote && fn.Retv.GoName != "" {
			switch ft := fn.Retv.GoType.(type) {
			case string:
				fn.Body = fmt.Sprintf("\treturn \"%s\"", ft)
			case nil:
				fn.Body = "\treturn nil"
			default:
				fn.Body = fmt.Sprintf("\treturn %v", ft)
			}
		}
		// handle other retv
		if strings.HasSuffix(fn.Body, "\n") {
			fn.Body = fn.Body[:len(fn.Body)-1]
		}
	}

	return src, nil
}
