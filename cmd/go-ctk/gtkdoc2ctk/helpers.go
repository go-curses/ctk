// Copyright (c) 2023  The Go-Curses Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gtkdoc2ctk

import (
	"regexp"
	"strings"

	"github.com/go-curses/cdk/lib/paint"
	"github.com/iancoleman/strcase"
)

var rxGdkClassName = regexp.MustCompile(`\b(Gdk|Gtk)([_a-zA-Z0-9]+?)\b`)
var rxGtkFuncName = regexp.MustCompile(`\b(gdk|gtk)_([_a-zA-Z0-9]+?)\b`)

func RewriteGtkThingsToCtkThings(class, input string) (output string) {
	output = input
	if rxGtkFuncName.MatchString(input) {
		matches := rxGtkFuncName.FindAllStringSubmatch(input, -1)
		snake := strcase.ToSnake(class)
		for _, m := range matches {
			text := m[0]
			// prefix := m[1]
			match := m[2]
			if strings.HasPrefix(match, snake+"_") {
				match = match[len(snake):]
			}
			output = strings.Replace(output, text, strcase.ToCamel(match), 1)
		}
	}
	if rxGdkClassName.MatchString(input) {
		output = rxGdkClassName.ReplaceAllString(output, `$2`)
	}
	output = strings.Replace(output, "()", "", -1)
	output = strings.Replace(output, "Gdk"+class, class, -1)
	output = strings.Replace(output, "Gtk"+class, class, -1)
	// output = strings.Replace(output, "::", "", -1)
	output = strings.Replace(output, "GTK+", "CTK", -1)
	return
}

func TranslateGtkType(pkg, cType string) (t GtkType) {
	t.C = cType
	switch cType {
	case "gint", "guint", "gshort", "gushort", "glong", "gsize", "goffset":
		t.GoName = "int"
		t.GoLabel = "int"
		t.GoType = 0
	case "gdouble", "gfloat":
		t.GoName = "float64"
		t.GoLabel = "float"
		t.GoType = 0.0
	case "gchar", "guchar", "gstrv":
		t.GoName = "string"
		t.GoLabel = "string"
		t.GoType = ""
	case "gboolean":
		t.GoName = "bool"
		t.GoLabel = "bool"
		t.GoType = false
	case "style":
		t.GoName = "cdk.Style"
		t.GoLabel = "style"
		t.GoType = paint.Style{}
	case "void":
		t.GoName = ""
		t.GoLabel = ""
		t.GoType = nil
	case "gpointer", "gconstpointer", "gintptr", "guintptr":
		t.GoName = "interface{}"
		t.GoType = nil
		t.GoLabel = "struct"
	case "...interface{}":
		t.GoName = "...interface{}"
		t.GoType = []interface{}{}
		t.GoLabel = "interface{}"
	default:
		t.GoType = nil
		t.GoLabel = "struct"
		if pkg == "ctk" {
			t.GoName = strings.Replace(cType, "Gdk", "", -1)
			t.GoName = strings.Replace(t.GoName, "Gtk", "", -1)
		} else if pkg == "cdk" {
			t.GoName = strings.Replace(cType, "Gdk", "", -1)
			t.GoName = strings.Replace(t.GoName, "Gtk", "ctk.", -1)
		} else {
			t.GoName = strings.Replace(cType, "Gdk", "ctk.", -1)
			t.GoName = strings.Replace(t.GoName, "Gtk", "ctk.", -1)
		}
	}
	return
}

func TranslateNamedVariable(pkg, pType, pName string) (pn string, pt GtkType) {
	pType = rxStripDigitSuffix.ReplaceAllString(pType, "")
	if strings.HasPrefix(pName, "*") {
		pName = pName[1:]
	}
	pName = strcase.ToLowerCamel(pName)
	pt, pn = TranslateGtkType(pkg, pType), pName
	return
}