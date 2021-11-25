package ctk

import (
	"github.com/go-curses/ctk/lib/enums"
)

type ActionEntry struct {
	Name        string
	StockId     string
	Label       string
	Accelerator string
	Tooltip     string
	Callback    enums.GCallback
}
