package gtkdoc2ctk

type GtkType struct {
	C       string
	GoName  string
	GoLabel string
	GoType  interface{}
}

func NewType(pkg, c string) GtkType {
	return TranslateGtkType(pkg, c)
}

func (t GtkType) String() string {
	return t.GoName
}
