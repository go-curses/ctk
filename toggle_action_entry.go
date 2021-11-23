package ctk

type ToggleActionEntry struct {
	Name        string
	StockId     string
	Label       string
	Accelerator string
	Tooltip     string
	Callback    GCallback
	IsActive    bool
}
