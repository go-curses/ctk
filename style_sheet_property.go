package ctk

import (
	"fmt"

	tcss "github.com/tdewolff/parse/v2/css"
)

type StyleSheetProperty struct {
	Key   string
	Value string
	Type  tcss.TokenType
}

func (e StyleSheetProperty) String() string {
	return fmt.Sprintf("%v:%v (%v);", e.Key, e.Value, e.Type)
}
