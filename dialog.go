package ctk

// TODO: SetWindow calls should only happen in expose/realize events
// TODO: implement expose/realize concept

import (
	"fmt"

	"github.com/go-curses/cdk"
	cenums "github.com/go-curses/cdk/lib/enums"
	cmath "github.com/go-curses/cdk/lib/math"
	"github.com/go-curses/cdk/lib/ptypes"
	"github.com/go-curses/cdk/log"
	"github.com/go-curses/ctk/lib/enums"
)

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
//
// The Dialog Widget is a Window with actionable Buttons, typically intended to
// be used as a transient for another Window rather than a Window on its own.
type Dialog interface {
	Window
	Buildable

	Run() (response chan enums.ResponseType)
	RunFunc(fn func(response enums.ResponseType, argv ...interface{}), argv ...interface{})
	Response(responseId enums.ResponseType)
	GetDialogFlags() (flags enums.DialogFlags)
	SetDialogFlags(flags enums.DialogFlags)
	AddButton(buttonText string, responseId enums.ResponseType) (button Button)
	AddButtons(argv ...interface{})
	AddActionWidget(child Widget, responseId enums.ResponseType)
	AddSecondaryActionWidget(child Widget, responseId enums.ResponseType)
	SetDefaultResponse(responseId enums.ResponseType)
	SetResponseSensitive(responseId enums.ResponseType, sensitive bool)
	GetResponseForWidget(widget Widget) (value enums.ResponseType)
	GetWidgetForResponse(responseId enums.ResponseType) (value Widget)
	GetActionArea() (value ButtonBox)
	GetContentArea() (value VBox)
}

var _ Dialog = (*CDialog)(nil)

// The CDialog structure implements the Dialog interface and is exported
// to facilitate type embedding with custom implementations. No member variables
// are exported as the interface methods are the only intended means of
// interacting with Dialog objects.
type CDialog struct {
	CWindow

	dialogFlags enums.DialogFlags
	defResponse enums.ResponseType
	content     VBox
	action      ButtonBox
	widgets     map[enums.ResponseType][]Widget

	done     chan bool
	response enums.ResponseType
}

// MakeDialog is used by the Buildable system to construct a new Dialog.
func MakeDialog() Dialog {
	return NewDialog()
}

// NewDialog is the constructor for new Dialog instances.
func NewDialog() (value Dialog) {
	d := new(CDialog)
	d.Init()
	return d
}

// NewDialogWithButtons creates a new Dialog with title, transient parent, a
// bitmask of DialogFlags and a variadic list of paired items. The items are
// the button ResponseType paired with a Button label string (which can be a
// ctk.StockID for access to the stock Buttons in CTK).
//
// The `flags` argument can be used to make the dialog modal (ctk.DialogModal)
// and/or to have it destroyed along with its transient parent
// (ctk.DialogDestroyWithParent).
//
// If the user clicks one of these Dialog Buttons, the Dialog will emit the
// response signal with the corresponding response ID. Buttons are from left to
// right, so the first button in the list will be the leftmost button in the
// Dialog.
//
// Parameters:
// 	title	label for the dialog
// 	parent	Transient parent of the dialog, or `nil`
// 	flags	from DialogFlags
// 	argv	response ID with label pairs
func NewDialogWithButtons(title string, parent Window, flags enums.DialogFlags, argv ...interface{}) (dialog Dialog) {
	d := new(CDialog)
	d.dialogFlags = flags
	d.Init()
	d.SetTitle(title)
	if parent != nil {
		d.SetTransientFor(parent)
		d.SetParent(parent)
		if err := d.ImportStylesFromString(parent.ExportStylesToString()); err != nil {
			d.LogErr(err)
		}
	}
	d.SetWindow(d)
	if len(argv) > 0 {
		d.AddButtons(argv...)
	}
	return d
}

func NewMessageDialog(title, message string) (dialog Dialog) {
	d := new(CDialog)
	d.Init()
	d.SetTitle(title)
	d.SetWindow(d)
	d.AddButton(string(StockClose), enums.ResponseClose)
	d.SetDefaultResponse(enums.ResponseClose)
	var label Label
	var err error
	if label, err = NewLabelWithMarkup(message); err != nil {
		log.Error(err)
		label = NewLabel(message)
	}
	label.Show()
	label.SetSingleLineMode(false)
	label.SetLineWrap(true)
	label.SetLineWrapMode(cenums.WRAP_WORD)
	d.GetContentArea().PackStart(label, true, true, 0)
	w, h := label.GetSizeRequest()
	d.SetSizeRequest(w+4, h+4)
	return d
}

func NewYesNoDialog(title, message string, defaultNo bool) (dialog Dialog) {
	d := new(CDialog)
	d.Init()
	d.SetTitle(title)
	d.SetWindow(d)
	d.AddButton(string(StockYes), enums.ResponseYes)
	d.AddButton(string(StockNo), enums.ResponseNo)
	if defaultNo {
		d.SetDefaultResponse(enums.ResponseNo)
	} else {
		d.SetDefaultResponse(enums.ResponseYes)
	}
	var label Label
	var err error
	if label, err = NewLabelWithMarkup(message); err != nil {
		log.Error(err)
		label = NewLabel(message)
	}
	label.Show()
	label.SetSingleLineMode(false)
	label.SetLineWrap(true)
	label.SetLineWrapMode(cenums.WRAP_WORD)
	d.GetContentArea().PackStart(label, true, true, 0)
	w, h := label.GetSizeRequest()
	d.SetSizeRequest(w+4, h+4)
	return d
}

type buttonMenuOption struct {
	label    string
	response enums.ResponseType
}

func NewButtonMenuDialog(title, message string, argv ...interface{}) (dialog Dialog) {
	d := new(CDialog)
	d.Init()
	d.SetTitle(title)
	d.SetWindow(d)
	d.AddButton(string(StockCancel), enums.ResponseCancel)
	d.SetDefaultResponse(enums.ResponseCancel)

	contentArea := d.GetContentArea()

	var options []buttonMenuOption
	optionsWidth := 0

	argc := len(argv)

	if argc%2 != 0 {
		d.LogError("invalid button menu arguments: %v", argv)
		return nil
	}

	for i := 0; i < argc; i += 2 {
		bmo := buttonMenuOption{}
		if stock, ok := argv[i].(StockID); ok {
			bmo.label = string(stock)
		} else if label, ok := argv[i].(string); ok {
			bmo.label = label
		} else {
			d.LogError("invalid button menu label argument: %v (%T)", argv[i], argv[i])
			continue
		}

		if response, ok := argv[i+1].(enums.ResponseType); ok {
			bmo.response = response
		} else if response, ok := argv[i+1].(int); ok {
			bmo.response = enums.ResponseType(response)
		} else {
			d.LogError("invalid button menu response argument: %v (%T)", argv[i], argv[i])
			continue
		}

		labelLen := len(bmo.label)
		if labelLen > optionsWidth {
			optionsWidth = labelLen
		}
		options = append(options, bmo)
	}
	numOptions := len(options)

	var label Label
	var err error
	if label, err = NewLabelWithMarkup(message); err != nil {
		log.Error(err)
		label = NewLabel(message)
	}
	label.Show()
	label.SetSingleLineMode(false)
	label.SetLineWrap(true)
	label.SetLineWrapMode(cenums.WRAP_WORD)
	contentArea.PackStart(label, false, false, 0)

	w, h := label.GetSizeRequest()
	if optionsWidth+2 > w {
		w = optionsWidth + 2
	} else {
		optionsWidth = w - 2
	}
	h += cmath.CeilI(numOptions, 5)

	buttonScroll := NewScrolledViewport()
	buttonScroll.Show()
	if numOptions > 5 {
		buttonScroll.SetPolicy(enums.PolicyAutomatic, enums.PolicyAutomatic)
		optionsWidth -= 2
	} else {
		buttonScroll.SetPolicy(enums.PolicyNever, enums.PolicyNever)
	}
	contentArea.PackStart(buttonScroll, true, true, 0)

	buttonBox := NewVBox(false, 0)
	buttonBox.Show()
	buttonBox.SetSizeRequest(w, numOptions)
	buttonScroll.Add(buttonBox)

	for _, option := range options {
		button := NewButtonWithLabel(option.label)
		button.Show()
		button.SetSizeRequest(optionsWidth, 1)
		button.Connect(SignalActivate, "dialog-button-menu-activate-handler", func(data []interface{}, argv ...interface{}) cenums.EventFlag {
			if len(data) == 1 {
				if opt, ok := data[0].(buttonMenuOption); ok {
					d.Response(opt.response)
					return cenums.EVENT_STOP
				}
			}
			d.LogError("button menu activate handler, invalid data: %v", data)
			return cenums.EVENT_STOP
		}, option)
		buttonBox.PackStart(button, false, false, 0)
	}

	d.SetSizeRequest(w+2, h+4)
	return d
}

// Init initializes a Dialog object. This must be called at least once to
// set up the necessary defaults and allocate any memory structures. Calling
// this more than once is safe though unnecessary. Only the first call will
// result in any effect upon the Dialog instance. Init is used in the
// NewDialog constructor and only necessary when implementing a derivative
// Dialog type.
func (d *CDialog) Init() (already bool) {
	if d.InitTypeItem(TypeDialog, d) {
		return true
	}
	d.CWindow.Init()
	if d.dialogFlags == 0 {
		d.dialogFlags = enums.DialogModal | enums.DialogDestroyWithParent
	}
	d.flags = enums.NULL_WIDGET_FLAG
	d.SetFlags(enums.PARENT_SENSITIVE | enums.APP_PAINTABLE | enums.TOPLEVEL)
	d.SetParent(d)
	d.SetWindow(d)
	d.SetDecorated(false)
	d.defResponse = enums.ResponseNone
	d.done = make(chan bool, 1)
	d.response = enums.ResponseNone
	d.widgets = make(map[enums.ResponseType][]Widget)

	vbox := d.GetVBox()
	vbox.SetSpacing(1)

	d.content = NewVBox(false, 0)
	d.content.Show()
	vbox.PackStart(d.content, true, true, 0)

	d.action = NewHButtonBox(false, 1)
	d.action.Show()
	vbox.PackEnd(d.action, false, true, 1)

	d.Connect(SignalResponse, DialogResponseHandle, func(data []interface{}, argv ...interface{}) cenums.EventFlag {
		if len(argv) == 1 {
			if value, ok := argv[0].(enums.ResponseType); ok {
				d.response = value
			} else {
				d.LogError("response signal received invalid ResponseType: %v (%T)", argv[0], argv[0])
				d.response = enums.ResponseNone
			}
		} else {
			d.response = d.defResponse
		}
		d.done <- true
		return cenums.EVENT_STOP
	})
	d.Connect(SignalCdkEvent, DialogEventHandle, d.event)
	return false
}

// Build provides customizations to the Buildable system for Dialog Widgets.
func (d *CDialog) Build(builder Builder, element *CBuilderElement) error {
	d.Freeze()
	defer d.Thaw()
	if err := d.CWindow.Build(builder, element); err != nil {
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

// Run in CTK, unlike the GTK equivalent, does not block the main thread.
// Run maintains its own internal main-loop process and returns a ResponseType
// channel so that once the user presses one of the action-buttons or closes the
// Dialog with ESC for example, the response channel can deliver the user-input
// to the Dialog calling code.
//
// Before entering the recursive main loop, Run calls Show on the Dialog for
// you. Note that you still need to Show any children of the Dialog yourself.
// You can force Run to return at any time by calling Response to emit the
// ::response signal directly. Destroying the dialog during Run is a very
// bad idea, because your post-run code won't know whether the dialog was
// destroyed or not and there would likely be a closed-chan issue with the
// ResponseType channel.
//
// After Run returns, you are responsible for hiding or destroying the dialog if
// you wish to do so.
func (d *CDialog) Run() (response chan enums.ResponseType) {
	response = make(chan enums.ResponseType)
	display := d.GetDisplay()
	if display == nil {
		d.LogError("display not found")
		response <- enums.ResponseNone
		return
	}
	if !display.IsRunning() {
		d.LogError("display not running")
		response <- enums.ResponseNone
		return
	}
	if d.defResponse != enums.ResponseNone {
		if ab, ok := d.widgets[d.defResponse]; ok {
			last := len(ab) - 1
			if ab[last] != nil {
				ab[last].GrabFocus()
			}
		}
	}

	if transient := d.GetTransientFor(); transient != nil {
		// make sure the dialog will be "on top" of the correct window
		display.FocusWindow(transient)
	} else {
		display.FocusWindow(d)
	}

	d.SetRegion(d.getDialogRegion())
	d.Resize()
	d.Show()
	display.RequestDraw()
	display.RequestSync()
	cdk.Go(func() {
		// wait for the response event
		select {
		case <-d.done:
		}
		response <- d.response
		display.RequestDraw()
		display.RequestShow()
	})
	return
}

func (d *CDialog) RunFunc(fn func(response enums.ResponseType, argv ...interface{}), argv ...interface{}) {
	r := d.Run()
	cdk.Go(func() {
		response := <-r
		fn(response, argv...)
		d.Destroy()
		d.RequestDrawAndSync()
	})
}

// Response emits the response signal with the given response ID. Used to
// indicate that the user has responded to the dialog in some way; typically
// either you or Run will be monitoring the ::response signal and take
// appropriate action.
//
// Parameters:
// 	responseId	ResponseType identifier
func (d *CDialog) Response(responseId enums.ResponseType) {
	d.Emit(SignalResponse, responseId)
}

func (d *CDialog) Add(w Widget) {
	d.CBin.Add(w)
	w.SetWindow(d)
}

func (d *CDialog) GetWindow() Window {
	return d
}

func (d *CDialog) GetDialogFlags() (flags enums.DialogFlags) {
	d.RLock()
	flags = d.dialogFlags
	d.RUnlock()
	return
}

func (d *CDialog) SetDialogFlags(flags enums.DialogFlags) {
	d.Lock()
	d.dialogFlags = d.dialogFlags.Set(flags)
	d.Unlock()
}

// AddButton is a convenience method for AddActionWidget to create a Button with
// the given text (or a stock button, if button_text is a StockID) and set
// things up so that clicking the button will emit the response signal with the
// given ResponseType. The Button is appended to the end of the dialog's action
// area. The button widget is returned, but usually you don't need it.
//
// Parameters:
// 	buttonText	text of button, or stock ID
// 	responseId	response ID for the button
func (d *CDialog) AddButton(buttonText string, responseId enums.ResponseType) (button Button) {
	if item := LookupStockItem(StockID(buttonText)); item != nil {
		button = NewButtonFromStock(StockID(buttonText))
	} else {
		button = NewButtonWithLabel(buttonText)
	}
	button.Show()
	button.SetUseUnderline(true)
	button.SetSizeRequest(-1, 1)
	d.AddActionWidget(button, responseId)
	return
}

// AddButtons is a convenience method for AddButton to create many buttons, in
// the same way as calling AddButton repeatedly. Each Button must have both
// ResponseType and label text provided.
//
// Parameters:
// 	argv	response ID with label pairs
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
		var responseId enums.ResponseType
		if responseId, ok = argv[i+1].(enums.ResponseType); !ok {
			d.LogError("invalid ResponseType argument: %v (%T)", argv[i])
			continue
		}
		d.AddButton(text, responseId)
	}
}

// AddActionWidget adds the given activatable widget to the action area of a
// Dialog, connecting a signal handler that will emit the response signal on the
// Dialog when the widget is activated. The widget is appended to the end of the
// Dialog's action area. If you want to add a non-activatable widget, simply
// pack it into the action_area field of the Dialog struct.
//
// Parameters:
// 	child	an activatable widget
// 	responseId	response ID for child
func (d *CDialog) AddActionWidget(child Widget, responseId enums.ResponseType) {
	child.Connect(SignalActivate, DialogActivateHandle, func(data []interface{}, argv ...interface{}) cenums.EventFlag {
		d.LogDebug("responding with: %v", responseId)
		d.Response(responseId)
		return cenums.EVENT_STOP
	})
	d.action.PackStart(child, false, false, 0)
	d.widgets[responseId] = append(d.widgets[responseId], child)
}

// AddSecondaryActionWidget is the same as AddActionWidget with the exception of
// adding the given Widget to the secondary action Button grouping instead of
// the primary grouping as with AddActionWidget.
func (d *CDialog) AddSecondaryActionWidget(child Widget, responseId enums.ResponseType) {
	child.Connect(SignalActivate, DialogActivateHandle, func(data []interface{}, argv ...interface{}) cenums.EventFlag {
		d.LogDebug("responding with: %v", responseId)
		d.Response(responseId)
		return cenums.EVENT_STOP
	})
	d.action.PackEnd(child, false, false, 0)
	d.widgets[responseId] = append(d.widgets[responseId], child)
}

// SetDefaultResponse updates which action Widget is activated when the user
// presses the ENTER key without changing the focused Widget first. The last
// Widget in the Dialog's action area with the given ResponseType as the default
// widget for the dialog.
//
// Parameters:
// 	responseId	a response ID
func (d *CDialog) SetDefaultResponse(responseId enums.ResponseType) {
	d.Lock()
	d.defResponse = responseId
	d.Unlock()
}

// SetResponseSensitive calls Widget.SetSensitive for each widget in the
// Dialog's action area with the given Responsetype. A convenient way to
// sensitize/desensitize Dialog Buttons.
//
// Parameters:
// 	responseId	a response ID
// 	setting	TRUE for sensitive
func (d *CDialog) SetResponseSensitive(responseId enums.ResponseType, sensitive bool) {
	if list, ok := d.widgets[responseId]; ok {
		for _, w := range list {
			w.SetSensitive(sensitive)
		}
	}
}

// GetResponseForWidget is a convenience method for looking up the ResponseType
// associated with the given Widget in the Dialog action area. Returns
// ResponseNone if the Widget is not found in the action area of the Dialog.
//
// Parameters:
// 	widget	a widget in the action area of dialog
func (d *CDialog) GetResponseForWidget(widget Widget) (value enums.ResponseType) {
	d.RLock()
	dwidgets := d.widgets
	d.RUnlock()
	for response, widgets := range dwidgets {
		for _, w := range widgets {
			if w.ObjectID() == widget.ObjectID() {
				return response
			}
		}
	}
	return enums.ResponseNone
}

// GetWidgetForResponse returns the last Widget Button that uses the given
// ResponseType in the action area of a Dialog.
//
// Parameters:
// 	responseId	the response ID used by the dialog widget
func (d *CDialog) GetWidgetForResponse(responseId enums.ResponseType) (value Widget) {
	d.RLock()
	defer d.RUnlock()
	if widgets, ok := d.widgets[responseId]; ok {
		if last := len(widgets) - 1; last > -1 {
			value = widgets[last]
		}
	}
	return
}

// GetActionArea returns the action area ButtonBox of a Dialog instance.
func (d *CDialog) GetActionArea() (value ButtonBox) {
	d.RLock()
	defer d.RUnlock()
	return d.action
}

// GetContentArea returns the content area VBox of a Dialog instance.
func (d *CDialog) GetContentArea() (value VBox) {
	d.RLock()
	defer d.RUnlock()
	return d.content
}

// Show ensures that the Dialog, content and action areas are all set to VISIBLE
func (d *CDialog) Show() {
	d.CWindow.Show()
	if widgets, ok := d.widgets[d.defResponse]; ok {
		if len(widgets) > 0 {
			widgets[0].GrabFocus()
		}
	}
}

// ShowAll calls ShowAll upon the Dialog, content area, action area and all the
// action Widget children.
func (d *CDialog) ShowAll() {
	d.Show()
	d.CWindow.ShowAll()
}

// Destroy hides the Dialog, removes it from any transient Window associations,
// removes the Dialog from the Display and finally emits the destroy-event
// signal.
func (d *CDialog) Destroy() {
	d.Hide()
	display := d.GetDisplay()
	if tf := d.GetTransientFor(); tf != nil {
		tf.SetTransientFor(nil)
		d.SetTransientFor(nil)
		display.UnmapWindow(d)
		display.FocusWindow(tf)
	} else {
		display.UnmapWindow(d)
	}
	d.Emit(SignalDestroyEvent, d)
}

func (d *CDialog) getDialogRegion() (region ptypes.Region) {
	if display := cdk.GetDefaultDisplay(); display != nil && display.IsRunning() {
		alloc := ptypes.MakeRectangle(display.Screen().Size())
		region = d.getDialogRegionForAllocation(alloc)
	}
	return
}

func (d *CDialog) getDialogRegionForAllocation(alloc ptypes.Rectangle) (region ptypes.Region) {
	if display := cdk.GetDefaultDisplay(); display != nil && display.IsRunning() {
		origin := ptypes.MakePoint2I(0, 0)
		req := d.SizeRequest()
		if req.W <= -1 || req.W > alloc.W {
			req.W = alloc.W
		}
		if req.H <= -1 || req.H > alloc.H {
			req.H = alloc.H
		}
		// center placement gravity
		if alloc.W > req.W {
			delta := alloc.W - req.W
			origin.X = int(float64(delta) * 0.5)
		}
		if alloc.H > req.H {
			delta := alloc.H - req.H
			origin.Y = int(float64(delta) * 0.5)
		}
		region = ptypes.MakeRegion(origin.X, origin.Y, req.W, req.H)
	}
	return
}

func (d *CDialog) event(data []interface{}, argv ...interface{}) cenums.EventFlag {
	if evt, ok := argv[1].(cdk.Event); ok {
		switch e := evt.(type) {
		case *cdk.EventKey:
			switch e.Key() {
			case cdk.KeyEscape:
				if f := d.Emit(SignalClose); f == cenums.EVENT_PASS {
					d.Response(enums.ResponseClose)
				}
				return cenums.EVENT_STOP
			}

		case *cdk.EventResize:
			alloc := ptypes.MakeRectangle(e.Size())
			region := d.getDialogRegionForAllocation(alloc)
			d.SetOrigin(region.X, region.Y)
			d.SetAllocation(region.Size())
			return d.Resize()
		}
		// do not block parent event handlers with cenums.EVENT_STOP, always
		// allow event-pass-through so that normal Window events can be handled
		// without duplicating the features here
	}
	return cenums.EVENT_PASS
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

const DialogEventHandle = "dialog-event-handler"

const DialogInvalidateHandle = "dialog-invalidate-handler"

const DialogResizeHandle = "dialog-resize-handler"

const DialogDrawHandle = "dialog-draw-handler"

const DialogActivateHandle = "dialog-activate-handler"