package ctk

type ActionEntry struct {
	Name        string
	StockId     string
	Label       string
	Accelerator string
	Tooltip     string
	Callback    GCallback
}
