package ctk

import (
	"github.com/go-curses/cdk"
)

// Stock items represent commonly-used menu or toolbar items such as "Open" or
// "Exit". Each stock item is identified by a stock ID; stock IDs are just
// strings, and constants such as StockOpen are provided to avoid typing
// mistakes in the strings. Applications can register their own stock items in
// addition to those built-in to CTK.
//
// Each stock ID can be associated with a StockItem, which contains the
// user-visible label, keyboard accelerator, and translation domain of the menu
// or toolbar item; and/or with an icon stored in a GtkIconFactory. See
// IconFactory for more information on stock icons. The connection between a
// StockItem and stock icons is purely conventional (by virtue of using the same
// stock ID); it's possible to register a stock item but no icon, and vice
// versa. Stock icons may have a RTL variant which gets used for right-to-left
// locales.
type StockID string

const (
	StockAbout          StockID = "ctk-about"
	StockAdd            StockID = "ctk-add"
	StockApply          StockID = "ctk-apply"
	StockBold           StockID = "ctk-bold"
	StockCancel         StockID = "ctk-cancel"
	StockClear          StockID = "ctk-clear"
	StockClose          StockID = "ctk-close"
	StockConvert        StockID = "ctk-convert"
	StockConnect        StockID = "ctk-connect"
	StockCopy           StockID = "ctk-copy"
	StockCut            StockID = "ctk-cut"
	StockDelete         StockID = "ctk-delete"
	StockDialogError    StockID = "ctk-dialog-error"
	StockDialogInfo     StockID = "ctk-dialog-info"
	StockDialogQuestion StockID = "ctk-dialog-question"
	StockDialogWarning  StockID = "ctk-dialog-warning"
	StockDirectory      StockID = "ctk-directory"
	StockDiscard        StockID = "ctk-discard"
	StockDisconnect     StockID = "ctk-disconnect"
	StockEdit           StockID = "ctk-edit"
	StockExecute        StockID = "ctk-execute"
	StockFile           StockID = "ctk-file"
	StockFind           StockID = "ctk-find"
	StockFindAndReplace StockID = "ctk-find-and-replace"
	StockGotoBottom     StockID = "ctk-goto-bottom"
	StockGotoFirst      StockID = "ctk-goto-first"
	StockGotoLast       StockID = "ctk-goto-last"
	StockGotoTop        StockID = "ctk-goto-top"
	StockGoBack         StockID = "ctk-go-back"
	StockGoDown         StockID = "ctk-go-down"
	StockGoForward      StockID = "ctk-go-forward"
	StockGoUp           StockID = "ctk-go-up"
	StockHelp           StockID = "ctk-help"
	StockHome           StockID = "ctk-home"
	StockIndent         StockID = "ctk-indent"
	StockIndex          StockID = "ctk-index"
	StockInfo           StockID = "ctk-info"
	StockItalic         StockID = "ctk-italic"
	StockJumpTo         StockID = "ctk-jump-to"
	StockJustifyCenter  StockID = "ctk-justify-center"
	StockJustifyFill    StockID = "ctk-justify-fill"
	StockJustifyLeft    StockID = "ctk-justify-left"
	StockJustifyRight   StockID = "ctk-justify-right"
	StockMediaForward   StockID = "ctk-media-forward"
	StockMediaNext      StockID = "ctk-media-next"
	StockMediaPause     StockID = "ctk-media-pause"
	StockMediaPlay      StockID = "ctk-media-play"
	StockMediaPrevious  StockID = "ctk-media-previous"
	StockMediaRecord    StockID = "ctk-media-record"
	StockMediaRewind    StockID = "ctk-media-rewind"
	StockMediaStop      StockID = "ctk-media-stop"
	StockNew            StockID = "ctk-new"
	StockNo             StockID = "ctk-no"
	StockOk             StockID = "ctk-ok"
	StockOpen           StockID = "ctk-open"
	StockPaste          StockID = "ctk-paste"
	StockPreferences    StockID = "ctk-preferences"
	StockProperties     StockID = "ctk-properties"
	StockQuit           StockID = "ctk-quit"
	StockRedo           StockID = "ctk-redo"
	StockRefresh        StockID = "ctk-refresh"
	StockRemove         StockID = "ctk-remove"
	StockRevertToSaved  StockID = "ctk-revert-to-saved"
	StockSave           StockID = "ctk-save"
	StockSaveAs         StockID = "ctk-save-as"
	StockSelectAll      StockID = "ctk-select-all"
	StockSelectColor    StockID = "ctk-select-color"
	StockSelectFont     StockID = "ctk-select-font"
	StockSortAscending  StockID = "ctk-sort-ascending"
	StockSortDescending StockID = "ctk-sort-descending"
	StockStop           StockID = "ctk-stop"
	StockStrikethrough  StockID = "ctk-strikethrough"
	StockUndelete       StockID = "ctk-undelete"
	StockUnderline      StockID = "ctk-underline"
	StockUndo           StockID = "ctk-undo"
	StockUnindent       StockID = "ctk-unindent"
	StockYes            StockID = "ctk-yes"
	StockZoom100        StockID = "ctk-zoom-100"
	StockZoomFit        StockID = "ctk-zoom-fit"
	StockZoomIn         StockID = "ctk-zoom-in"
	StockZoomOut        StockID = "ctk-zoom-out"
)

var ctkStockItemRegistry = map[StockID]*StockItem{
	StockDialogInfo:     {ID: StockDialogInfo, Label: "Information"},
	StockDialogWarning:  {ID: StockDialogWarning, Label: "Warning"},
	StockDialogError:    {ID: StockDialogError, Label: "Error"},
	StockDialogQuestion: {ID: StockDialogQuestion, Label: "Question"},
	StockAbout:          {ID: StockAbout, Label: "_About"},
	StockAdd:            {ID: StockAdd, Label: "_Add"},
	StockApply:          {ID: StockApply, Label: "_Apply"},
	StockBold:           {ID: StockBold, Label: "_Bold"},
	StockCancel:         {ID: StockCancel, Label: "_Cancel"},
	StockClear:          {ID: StockClear, Label: "_Clear"},
	StockClose:          {ID: StockClose, Label: "_Close"},
	StockConnect:        {ID: StockConnect, Label: "C_onnect"},
	StockConvert:        {ID: StockConvert, Label: "_Convert"},
	StockCopy:           {ID: StockCopy, Label: "_Copy"},
	StockCut:            {ID: StockCut, Label: "Cu_t"},
	StockDelete:         {ID: StockDelete, Label: "_Delete"},
	StockDirectory:      {ID: StockDirectory, Label: "_Directory"},
	StockDiscard:        {ID: StockDiscard, Label: "_Discard"},
	StockDisconnect:     {ID: StockDisconnect, Label: "_Disconnect"},
	StockExecute:        {ID: StockExecute, Label: "_Execute"},
	StockEdit:           {ID: StockEdit, Label: "_Edit"},
	StockFile:           {ID: StockFile, Label: "_File"},
	StockFind:           {ID: StockFind, Label: "_Find"},
	StockFindAndReplace: {ID: StockFindAndReplace, Label: "Find and _Replace"},
	StockGotoBottom:     {ID: StockGotoBottom, Label: "_Bottom"},
	StockGotoFirst:      {ID: StockGotoFirst, Label: "_First"},
	StockGotoLast:       {ID: StockGotoLast, Label: "_Last"},
	StockGotoTop:        {ID: StockGotoTop, Label: "_Top"},
	StockGoBack:         {ID: StockGoBack, Label: "_Back"},
	StockGoDown:         {ID: StockGoDown, Label: "_Down"},
	StockGoForward:      {ID: StockGoForward, Label: "_Forward"},
	StockGoUp:           {ID: StockGoUp, Label: "_Up"},
	StockHelp:           {ID: StockHelp, Label: "_Help"},
	StockHome:           {ID: StockHome, Label: "_Home"},
	StockIndent:         {ID: StockIndent, Label: "Increase Indent"},
	StockUnindent:       {ID: StockUnindent, Label: "Decrease Indent"},
	StockIndex:          {ID: StockIndex, Label: "_Index"},
	StockInfo:           {ID: StockInfo, Label: "_Information"},
	StockItalic:         {ID: StockItalic, Label: "_Italic"},
	StockJumpTo:         {ID: StockJumpTo, Label: "_Jump to"},
	StockJustifyCenter:  {ID: StockJustifyCenter, Label: "_Center"},
	StockJustifyFill:    {ID: StockJustifyFill, Label: "_Fill"},
	StockJustifyLeft:    {ID: StockJustifyLeft, Label: "_Left"},
	StockJustifyRight:   {ID: StockJustifyRight, Label: "_Right"},
	StockMediaForward:   {ID: StockMediaForward, Label: "_Forward"},
	StockMediaNext:      {ID: StockMediaNext, Label: "_Next"},
	StockMediaPause:     {ID: StockMediaPause, Label: "P_ause"},
	StockMediaPlay:      {ID: StockMediaPlay, Label: "_Play"},
	StockMediaPrevious:  {ID: StockMediaPrevious, Label: "Pre_vious"},
	StockMediaRecord:    {ID: StockMediaRecord, Label: "_Record"},
	StockMediaRewind:    {ID: StockMediaRewind, Label: "R_ewind"},
	StockMediaStop:      {ID: StockMediaStop, Label: "_Stop"},
	StockNew:            {ID: StockNew, Label: "_New"},
	StockNo:             {ID: StockNo, Label: "_No"},
	StockOk:             {ID: StockOk, Label: "_OK"},
	StockOpen:           {ID: StockOpen, Label: "_Open"},
	StockPaste:          {ID: StockPaste, Label: "_Paste"},
	StockPreferences:    {ID: StockPreferences, Label: "_Preferences"},
	StockProperties:     {ID: StockProperties, Label: "_Properties"},
	StockQuit:           {ID: StockQuit, Label: "_Quit"},
	StockRedo:           {ID: StockRedo, Label: "_Redo"},
	StockRefresh:        {ID: StockRefresh, Label: "_Refresh"},
	StockRemove:         {ID: StockRemove, Label: "_Remove"},
	StockRevertToSaved:  {ID: StockRevertToSaved, Label: "_Revert"},
	StockSave:           {ID: StockSave, Label: "_Save"},
	StockSaveAs:         {ID: StockSaveAs, Label: "Save _As"},
	StockSelectAll:      {ID: StockSelectAll, Label: "Select _All"},
	StockSelectColor:    {ID: StockSelectColor, Label: "_Color"},
	StockSelectFont:     {ID: StockSelectFont, Label: "_Font"},
	StockSortAscending:  {ID: StockSortAscending, Label: "_Ascending"},
	StockSortDescending: {ID: StockSortDescending, Label: "_Descending"},
	StockStop:           {ID: StockStop, Label: "_Stop"},
	StockStrikethrough:  {ID: StockStrikethrough, Label: "_Strikethrough"},
	StockUndelete:       {ID: StockUndelete, Label: "_Undelete"},
	StockUnderline:      {ID: StockUnderline, Label: "_Underline"},
	StockUndo:           {ID: StockUndo, Label: "_Undo"},
	StockYes:            {ID: StockYes, Label: "_Yes"},
	StockZoom100:        {ID: StockZoom100, Label: "_Normal Size"},
	StockZoomFit:        {ID: StockZoomFit, Label: "Best _Fit"},
	StockZoomIn:         {ID: StockZoomIn, Label: "Zoom _In"},
	StockZoomOut:        {ID: StockZoomOut, Label: "Zoom _Out"},
}

type StockItem struct {
	ID         StockID
	Label      string
	Modifier   ModifierType
	KeyVal     cdk.Key
	I18nDomain string
}

// Registers each of the stock items in items. If an item already exists with
// the same stock ID as one of the items, the old item gets replaced.
func AddStockItems(items ...*StockItem) {
	for _, item := range items {
		ctkStockItemRegistry[item.ID] = item
	}
}

// Retrieves a list of all known stock IDs added to an IconFactory or registered
// with StockAdd().
func ListStockIDs() (list []StockID) {
	for id := range ctkStockItemRegistry {
		list = append(list, id)
	}
	return
}

// Retrieve a stock item by ID. Returns nil if item not found.
func LookupStockItem(id StockID) (item *StockItem) {
	var ok bool
	if item, ok = ctkStockItemRegistry[id]; !ok {
		item = nil
		// log.ErrorDF(1, "stock id not found: \"%v\"", id)
	}
	return
}

func LookupStockLabel(label string) (item *StockItem) {
	for _, m := range ctkStockItemRegistry {
		if m.Label == label {
			item = m
			break
		}
	}
	return
}
