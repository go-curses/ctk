package ctk

// TODO: SetWindow calls should only happen in expose/realize events
// TODO: implement expose/realize concept

import (
	"fmt"

	"github.com/go-curses/cdk"
	"github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/cdk/lib/paint"
	"github.com/go-curses/cdk/lib/ptypes"
	"github.com/go-curses/cdk/memphis"
	"github.com/gofrs/uuid"
	"github.com/jtolio/gls"
)

// CDK type-tag for Dialog objects
const TypeDialog cdk.CTypeTag = "ctk-dialog"

func init() {
	_ = cdk.TypesManager.AddType(TypeDialog, func() interface{} { return MakeDialog() })
	ctkBuilderTranslators[TypeDialog] = func(builder Builder, widget Widget, name, value string) error {
		if fn, ok := ctkBuilderTranslators[TypeWindow]; ok {
			return fn(builder, widget, name, value)
		}
		return fmt.Errorf("dialog property translator not implemented")
	}
}

// Dialog Hierarchy:
//	Object
//	  +- Widget
//	    +- Container
//	      +- Bin
//	        +- Window
//	          +- Dialog
//	            +- AboutDialog
//	            +- ColorSelectionDialog
//	            +- FileChooserDialog
//	            +- FileSelection
//	            +- FontSelectionDialog
//	            +- InputDialog
//	            +- MessageDialog
//	            +- PageSetupUnixDialog
//	            +- PrintUnixDialog
//	            +- RecentChooserDialog
type Dialog interface {
	Window
	Buildable

	Init() (already bool)
	Run() (response chan ResponseType)
	Response(responseId ResponseType)
	AddButton(buttonText string, responseId ResponseType) (value Button)
	AddButtons(argv ...interface{})
	AddActionWidget(child Widget, responseId ResponseType)
	AddSecondaryActionWidget(child Widget, responseId ResponseType)
	SetDefaultResponse(responseId ResponseType)
	SetResponseSensitive(responseId ResponseType, sensitive bool)
	GetResponseForWidget(widget Widget) (value ResponseType)
	GetWidgetForResponse(responseId ResponseType) (value Widget)
	GetActionArea() (value ButtonBox)
	GetContentArea() (value VBox)
	Build(builder Builder, element *CBuilderElement) error
	Show()
	Destroy()
	ProcessEvent(evt cdk.Event) enums.EventFlag
	Resize() enums.EventFlag
	Invalidate() enums.EventFlag
}

// The CDialog structure implements the Dialog interface and is
// exported to facilitate type embedding with custom implementations. No member
// variables are exported as the interface methods are the only intended means
// of interacting with Dialog objects
type CDialog struct {
	CWindow

	dialogFlags DialogFlags
	defResponse ResponseType
	content     VBox
	action      ButtonBox
	widgets     map[ResponseType][]Widget

	done     chan bool
	response ResponseType
}

// Default constructor for Dialog objects
func MakeDialog() *CDialog {
	return NewDialog()
}

// Constructor for Dialog objects
func NewDialog() (value *CDialog) {
	d := new(CDialog)
	d.Init()
	return d
}

// Creates a new Dialog with title title (or NULL for the default title;
// see WindowSetTitle) and transient parent parent (or NULL for none;
// see WindowSetTransientFor). The flags argument can be used to
// make the dialog modal (GTK_DIALOG_MODAL) and/or to have it destroyed along
// with its transient parent (GTK_DIALOG_DESTROY_WITH_PARENT). After flags ,
// button text/response ID pairs should be listed, with a NULL pointer ending
// the list. Button text can be either a stock ID such as GTK_STOCK_OK, or
// some arbitrary text. A response ID can be any positive number, or one of
// the values in the ResponseType enumeration. If the user clicks one of
// these dialog buttons, Dialog will emit the response signal with
// the corresponding response ID. If a Dialog receives the
// delete-event signal, it will emit ::response with a response ID of
// GTK_RESPONSE_DELETE_EVENT. However, destroying a dialog does not emit the
// ::response signal; so be careful relying on ::response when using the
// GTK_DIALOG_DESTROY_WITH_PARENT flag. Buttons are from left to right, so
// the first button in the list will be the leftmost button in the dialog.
// Here's a simple example:
// Parameters:
// 	title	Title of the dialog, or NULL.
// 	parent	Transient parent of the dialog, or NULL.
// 	flags	from DialogFlags
// 	firstButtonText	stock ID or text to go in first button, or NULL.
// 	varargs	response ID for first button, then additional buttons, ending with NULL
// Returns:
// 	a new Dialog
func NewDialogWithButtons(title string, parent Window, flags DialogFlags, argv ...interface{}) (value *CDialog) {
	d := new(CDialog)
	d.dialogFlags = flags
	d.Init()
	d.SetTitle(title)
	// d.SetParent(d)
	// d.SetWindow(d)
	d.SetTransientFor(parent)
	if len(argv) > 0 {
		d.AddButtons(argv...)
	}
	return d
}

// Dialog object initialization. This must be called at least once to setup
// the necessary defaults and allocate any memory structures. Calling this more
// than once is safe though unnecessary. Only the first call will result in any
// effect upon the Dialog instance
func (d *CDialog) Init() (already bool) {
	if d.InitTypeItem(TypeDialog, d) {
		return true
	}
	d.CWindow.Init()
	if d.dialogFlags == 0 {
		d.dialogFlags = DialogModal | DialogDestroyWithParent
	}
	d.flags = NULL_WIDGET_FLAG
	d.SetFlags(PARENT_SENSITIVE | APP_PAINTABLE)
	d.parent = d
	d.defResponse = ResponseNone
	vbox := d.GetVBox()
	vbox.Show()
	d.content = NewVBox(false, 0)
	d.content.Show()
	vbox.PackStart(d.content, true, true, 0)
	d.action = NewHButtonBox(false, 0)
	// d.action.SetSizeRequest(-1, 3)
	d.action.Show()
	vbox.PackEnd(d.action, false, true, 0)
	d.done = make(chan bool, 1)
	d.response = ResponseNone
	d.Connect(SignalDraw, DialogDrawHandle, d.draw)
	d.Connect(SignalResponse, DialogResponseHandle, func(data []interface{}, argv ...interface{}) enums.EventFlag {
		if len(argv) == 1 {
			if value, ok := argv[0].(ResponseType); ok {
				d.response = value
			} else {
				d.LogError("response signal received invalid ResponseType: %v (%T)", argv[0], argv[0])
				d.response = ResponseNone
			}
		} else {
			d.response = d.defResponse
		}
		d.done <- true
		return enums.EVENT_PASS
	})
	d.widgets = make(map[ResponseType][]Widget)
	return false
}

// Blocks in a recursive main loop until the dialog either emits the
// response signal, or is destroyed. If the dialog is destroyed during
// the call to Run, Run returns GTK_RESPONSE_NONE.
// Otherwise, it returns the response ID from the ::response signal emission.
// Before entering the recursive main loop, Run calls
// WidgetShow on the dialog for you. Note that you still need to show
// any children of the dialog yourself. During Run, the default
// behavior of delete-event is disabled; if the dialog receives
// ::delete_event, it will not be destroyed as windows usually are, and
// Run will return GTK_RESPONSE_DELETE_EVENT. Also, during
// Run the dialog will be modal. You can force Run
// to return at any time by calling Response to emit the
// ::response signal. Destroying the dialog during Run is a very
// bad idea, because your post-run code won't know whether the dialog was
// destroyed or not. After Run returns, you are responsible for
// hiding or destroying the dialog if you wish to do so. Typical usage of
// this function might be: Note that even though the recursive main loop
// gives the effect of a modal dialog (it prevents the user from interacting
// with other windows in the same window group while the dialog is run),
// callbacks such as timeouts, IO channel watches, DND drops, etc, will be
// triggered during a Run call.
// Returns:
// 	response ID
func (d *CDialog) Run() (response chan ResponseType) {
	response = make(chan ResponseType, 1)
	display := d.GetDisplay()
	screen := display.Screen()
	if screen == nil {
		d.LogError("screen not found")
		response <- ResponseNone
		return
	}
	parentId := uuid.Nil
	dw, dh := screen.Size()
	previousWindow := display.ActiveWindow()
	if transient := d.GetTransientFor(); transient != nil {
		parentId = transient.ObjectID()
		display.AddWindowOverlay(parentId, d, d.getDialogRegion())
		d.Resize()
		display.SetActiveWindow(transient)
	} else {
		d.SetAllocation(ptypes.MakeRectangle(dw, dh))
		d.Resize()
		display.SetActiveWindow(d)
	}
	display.RequestDraw()
	display.RequestShow()
	gls.Go(func() {
		// wait for the response event
		select {
		case <-d.done:
		}
		response <- d.response
		if parentId != uuid.Nil {
			display.RemoveWindowOverlay(parentId, d.ObjectID())
		}
		display.SetActiveWindow(previousWindow)
		display.RequestDraw()
		display.RequestShow()
	})
	return
}

// Emits the response signal with the given response ID. Used to
// indicate that the user has responded to the dialog in some way; typically
// either you or Run will be monitoring the ::response signal
// and take appropriate action.
// Parameters:
// 	responseId	response ID
func (d *CDialog) Response(responseId ResponseType) {
	d.Emit(SignalResponse, responseId)
}

// Adds a button with the given text (or a stock button, if button_text is a
// stock ID) and sets things up so that clicking the button will emit the
// response signal with the given response_id . The button is appended
// to the end of the dialog's action area. The button widget is returned, but
// usually you don't need it.
// Parameters:
// 	buttonText	text of button, or stock ID
// 	responseId	response ID for the button
// Returns:
// 	the button widget that was added.
// 	[transfer none]
func (d *CDialog) AddButton(buttonText string, responseId ResponseType) (button Button) {
	if item := LookupStockItem(StockID(buttonText)); item != nil {
		button = NewButtonFromStock(StockID(buttonText))
	} else {
		button = NewButtonWithLabel(buttonText)
	}
	button.Show()
	d.AddActionWidget(button, responseId)
	return
}

// Adds more buttons, same as calling AddButton repeatedly. The
// variable argument list should be NULL-terminated as with
// NewWithButtons. Each button must have both text and
// response ID.
// Parameters:
// 	firstButtonText	button text or stock ID
// 	varargs	response ID for first button, then more text-response_id pairs
func (d *CDialog) AddButtons(argv ...interface{}) {
	if len(argv)%2 != 0 {
		d.LogError("not an even number of arguments given")
		return
	}
	for i := 0; i < len(argv); i += 2 {
		var ok bool
		var text string
		if text, ok = argv[i].(string); !ok {
			if stockId, ok := argv[i].(StockID); ok {
				text = string(stockId)
			} else {
				d.LogError("invalid text argument: %v (%T)", argv[i])
				continue
			}
		}
		var responseId ResponseType
		if responseId, ok = argv[i+1].(ResponseType); !ok {
			d.LogError("invalid ResponseType argument: %v (%T)", argv[i])
			continue
		}
		d.AddButton(text, responseId)
	}
}

// Adds an activatable widget to the action area of a Dialog, connecting a
// signal handler that will emit the response signal on the dialog when
// the widget is activated. The widget is appended to the end of the dialog's
// action area. If you want to add a non-activatable widget, simply pack it
// into the action_area field of the Dialog struct.
// Parameters:
// 	child	an activatable widget
// 	responseId	response ID for child
//
func (d *CDialog) AddActionWidget(child Widget, responseId ResponseType) {
	child.Connect(SignalActivate, DialogActivateHandle, func(data []interface{}, argv ...interface{}) enums.EventFlag {
		d.LogDebug("responding with: %v", responseId)
		d.Response(responseId)
		return enums.EVENT_STOP
	})
	d.action.PackStart(child, false, false, 0)
	d.widgets[responseId] = append(d.widgets[responseId], child)
}

func (d *CDialog) AddSecondaryActionWidget(child Widget, responseId ResponseType) {
	child.Connect(SignalActivate, DialogActivateHandle, func(data []interface{}, argv ...interface{}) enums.EventFlag {
		d.LogDebug("responding with: %v", responseId)
		d.Response(responseId)
		return enums.EVENT_STOP
	})
	d.action.PackEnd(child, false, false, 0)
	d.widgets[responseId] = append(d.widgets[responseId], child)
}

// Sets the last widget in the dialog's action area with the given
// response_id as the default widget for the dialog. Pressing "Enter"
// normally activates the default widget.
// Parameters:
// 	responseId	a response ID
func (d *CDialog) SetDefaultResponse(responseId ResponseType) {
	d.defResponse = responseId
}

// Calls WidgetSetSensitive (widget, setting ) for each widget in the
// dialog's action area with the given response_id . A convenient way to
// sensitize/desensitize dialog buttons.
// Parameters:
// 	responseId	a response ID
// 	setting	TRUE for sensitive
func (d *CDialog) SetResponseSensitive(responseId ResponseType, sensitive bool) {
	if list, ok := d.widgets[responseId]; ok {
		for _, w := range list {
			w.SetSensitive(sensitive)
		}
	}
}

// Gets the response id of a widget in the action area of a dialog.
// Parameters:
// 	widget	a widget in the action area of dialog
//
// Returns:
// 	the response id of widget , or GTK_RESPONSE_NONE if widget
// 	doesn't have a response id set.
func (d *CDialog) GetResponseForWidget(widget Widget) (value ResponseType) {
	for response, widgets := range d.widgets {
		for _, w := range widgets {
			if w.ObjectID() == widget.ObjectID() {
				return response
			}
		}
	}
	return ResponseNone
}

// Gets the widget button that uses the given response ID in the action area
// of a dialog.
// Parameters:
// 	responseId	the response ID used by the dialog
// widget
// Returns:
// 	the widget button that uses the given response_id , or NULL.
// 	[transfer none]
func (d *CDialog) GetWidgetForResponse(responseId ResponseType) (value Widget) {
	if widgets, ok := d.widgets[responseId]; ok {
		if last := len(widgets) - 1; last > -1 {
			value = widgets[last]
		}
	}
	return
}

// Returns the action area of dialog .
// Returns:
// 	the action area.
// 	[transfer none]
func (d *CDialog) GetActionArea() (value ButtonBox) {
	return d.action
}

// Returns the content area of dialog .
// Returns:
// 	the content area VBox.
// 	[transfer none]
func (d *CDialog) GetContentArea() (value VBox) {
	return d.content
}

// // Returns TRUE if dialogs are expected to use an alternative button order on
// // the screen screen . See SetAlternativeButtonOrder for more
// // details about alternative button order. If you need to use this function,
// // you should probably connect to the ::notify:gtk-alternative-button-order
// // signal on the Settings object associated to screen , in order to be
// // notified if the button order setting changes.
// // Parameters:
// // 	screen	a Screen, or NULL to use the default screen.
// // Returns:
// // 	Whether the alternative button order should be used
// func (d *CDialog) AlternativeButtonOrder(screen Screen) (value bool) {
// 	return false
// }
//
// // Sets an alternative button order. If the
// // gtk-alternative-button-order setting is set to TRUE, the dialog
// // buttons are reordered according to the order of the response ids passed to
// // this function. By default, CTK dialogs use the button order advocated by
// // the Gnome right, and the cancel button left of it. But the builtin CTK
// // dialogs and MessageDialogs do provide an alternative button order,
// // which is more suitable on some platforms, e.g. Windows. Use this function
// // after adding all the buttons to your dialog, as the following example
// // shows:
// // Parameters:
// // 	firstResponseId	a response id used by one dialog
// // 's buttons
// // 	varargs	a list of more response ids of dialog
// // 's buttons, terminated by -1
// func (d *CDialog) SetAlternativeButtonOrder(firstResponseId int, argv ...interface{}) {}
//
// // Sets an alternative button order. If the
// // gtk-alternative-button-order setting is set to TRUE, the dialog
// // buttons are reordered according to the order of the response ids in
// // new_order . See SetAlternativeButtonOrder for more
// // information. This function is for use by language bindings.
// // Parameters:
// // 	nParams	the number of response ids in new_order
// //
// // 	newOrder	an array of response ids of
// // dialog
// // 's buttons.
// func (d *CDialog) SetAlternativeButtonOrderFromArray(nParams int, newOrder int) {}

func (d *CDialog) Build(builder Builder, element *CBuilderElement) error {
	d.Freeze()
	defer d.Thaw()
	if err := d.CObject.Build(builder, element); err != nil {
		return err
	}
	if len(element.Children) > 0 {
		contentBox := d.GetContentArea()
		actionBox := d.GetActionArea()
		internalBox := element.Children[0]
		for _, child := range internalBox.Children {
			var class string
			var ok bool
			if class, ok = child.Attributes["class"]; !ok {
				builder.LogError("missing class attribute: %v", child.String())
				continue
			}
			switch class {
			case "GtkBox":
				if err := contentBox.Build(builder, child); err != nil {
					builder.LogErr(err)
				}
			case "GtkButtonBox":
				if err := actionBox.Build(builder, child); err != nil {
					builder.LogErr(err)
				}
			default:
				builder.LogError("unexpected child element: %v", child.String())
				continue
			}
		}
	}
	return nil
}

func (d *CDialog) Hide() {
	d.UnsetFlags(VISIBLE)
	d.CWindow.Hide()
}

func (d *CDialog) Show() {
	d.SetFlags(VISIBLE)
	d.CWindow.Show()
	d.content.Show()
	d.action.Show()
}

func (d *CDialog) ShowAll() {
	d.SetFlags(VISIBLE)
	d.CWindow.ShowAll()
	d.content.ShowAll()
	d.action.ShowAll()
	for _, child := range d.GetChildren() {
		child.ShowAll()
	}
	for _, children := range d.widgets {
		for _, child := range children {
			child.ShowAll()
		}
	}
}

func (d *CDialog) Destroy() {
	d.Hide()
	dm := d.GetDisplay()
	if tf := d.GetTransientFor(); tf != nil {
		dm.RemoveWindowOverlay(tf.ObjectID(), d.ObjectID())
	} else {
		dm.RemoveWindow(d.ObjectID())
	}
	d.Emit(SignalDestroyEvent, d)
}

func (d *CDialog) ProcessEvent(evt cdk.Event) enums.EventFlag {
	// d.Lock()
	// defer d.Unlock()
	switch e := evt.(type) {
	case *cdk.EventKey:
		switch e.Key() {
		case cdk.KeyEscape:
			if f := d.Emit(SignalClose); f == enums.EVENT_PASS {
				d.Response(ResponseClose)
			}
			return enums.EVENT_STOP
		}
	case *cdk.EventResize:
		if tw := d.GetTransientFor(); tw != nil {
			tw.ProcessEvent(evt)
		}
		return d.Resize()
	case *cdk.EventMouse:
		if f := d.Emit(SignalEventMouse, d, evt); f == enums.EVENT_PASS {
			if child := d.GetChild(); child != nil {
				if mw := child.GetWidgetAt(ptypes.NewPoint2I(e.Position())); mw != nil {
					if ms, ok := mw.(Sensitive); ok && ms.IsSensitive() && ms.IsVisible() {
						if f := ms.ProcessEvent(evt); f == enums.EVENT_STOP {
							return enums.EVENT_STOP
						}
					}
				}
			}
		}
	}
	return d.CWindow.ProcessEvent(evt)
}

// TODO: set-size-request makes dialog window size, truncated by actual size
// TODO: local canvas / child alloc and origin issues

func (d *CDialog) getDialogRegion() (region ptypes.Region) {
	if dm := cdk.GetDefaultDisplay(); dm != nil {
		var origin ptypes.Point2I
		alloc := ptypes.MakeRectangle(dm.Screen().Size())
		if tf := d.GetTransientFor(); tf != nil {
			req := d.SizeRequest()
			if req.W <= -1 || req.W > alloc.W {
				req.W = alloc.W
			}
			if req.H <= -1 || req.H > alloc.H {
				req.H = alloc.H
			}
			if alloc.W > req.W {
				delta := alloc.W - req.W
				origin.X = int(float64(delta) * 0.5)
			}
			if alloc.H > req.H {
				delta := alloc.H - req.H
				origin.Y = int(float64(delta) * 0.5)
			}
			region = ptypes.MakeRegion(origin.X, origin.Y, req.W, req.H)
		} else {
			region = ptypes.MakeRegion(origin.X, origin.Y, alloc.W, alloc.H)
		}
	}
	return
}

func (d *CDialog) Resize() enums.EventFlag {
	region := d.getDialogRegion().NewClone()
	if tf := d.GetTransientFor(); tf != nil {
		if dm := d.GetDisplay(); dm != nil {
			dm.SetWindowOverlayRegion(tf.ObjectID(), d.ObjectID(), *region)
		}
	}
	d.SetOrigin(region.X, region.Y)
	d.SetAllocation(region.Size())

	if child := d.GetChild(); child != nil {
		origin := ptypes.MakePoint2I(1, 1)
		child.SetOrigin(origin.X, origin.Y)
		alloc := region.Size().NewClone()
		alloc.Sub(2, 2)
		child.SetAllocation(*alloc)
		child.Resize()
		if err := memphis.ConfigureSurface(child.ObjectID(), origin, *alloc, child.GetThemeRequest().Content.Normal); err != nil {
			child.LogErr(err)
		}
	}

	d.Invalidate()
	return enums.EVENT_STOP
}

func (d *CDialog) Invalidate() enums.EventFlag {
	// d.rebuildFocusChain()
	origin := d.GetOrigin()
	alloc := d.GetAllocation()
	if err := memphis.ConfigureSurface(d.ObjectID(), origin, alloc, d.GetThemeRequest().Content.Normal); err != nil {
		d.LogErr(err)
	}
	return d.Emit(SignalInvalidate, d)
}

func (d *CDialog) draw(data []interface{}, argv ...interface{}) enums.EventFlag {
	if surface, ok := argv[1].(*memphis.CSurface); ok {
		d.Lock()
		defer d.Unlock()
		size := d.GetAllocation()
		if !d.IsVisible() || size.W == 0 || size.H == 0 {
			d.LogDebug("not visible, zero width or zero height")
			return enums.EVENT_PASS
		}
		d.LogTrace("%v", size)
		if d.GetTitle() != "" {
			surface.FillBorderTitle(false, d.GetTitle(), enums.JUSTIFY_CENTER, d.GetThemeRequest())
		} else {
			surface.FillBorder(false, true, d.GetThemeRequest())
		}
		// d.content.SetBoolProperty(cdk.PropertyDebug, true)
		vbox := d.GetVBox()
		if r := vbox.Draw(); r == enums.EVENT_STOP {
			if err := surface.Composite(vbox.ObjectID()); err != nil {
				vbox.LogErr(err)
			}
		}
		if debug, _ := d.GetBoolProperty(cdk.PropertyDebug); debug {
			surface.DebugBox(paint.ColorNavy, d.ObjectInfo())
		}
		return enums.EVENT_STOP
	}
	return enums.EVENT_PASS
}

// The ::close signal is a which gets emitted when the user uses a keybinding
// to close the dialog. The default binding for this signal is the Escape
// key.
const SignalClose cdk.Signal = "close"

// Emitted when an action widget is clicked, the dialog receives a delete
// event, or the application programmer calls Response. On a
// delete event, the response ID is GTK_RESPONSE_DELETE_EVENT. Otherwise, it
// depends on which action widget was clicked.
// Listener function arguments:
// 	responseId int	the response ID
const SignalResponse cdk.Signal = "response"

const DialogResponseHandle = "dialog-response-handler"
const DialogDrawHandle = "dialog-draw-handler"
const DialogActivateHandle = "dialog-activate-handler"
