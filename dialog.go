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

	Init() (already bool)
	Build(builder Builder, element *CBuilderElement) error
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
	Show()
	ShowAll()
	Destroy()
}

// The CDialog structure implements the Dialog interface and is exported
// to facilitate type embedding with custom implementations. No member variables
// are exported as the interface methods are the only intended means of
// interacting with Dialog objects.
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

// MakeDialog is used by the Buildable system to construct a new Dialog.
func MakeDialog() *CDialog {
	return NewDialog()
}

// NewDialog is the constructor for new Dialog instances.
func NewDialog() (value *CDialog) {
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
func NewDialogWithButtons(title string, parent Window, flags DialogFlags, argv ...interface{}) (value *CDialog) {
	d := new(CDialog)
	d.dialogFlags = flags
	d.Init()
	d.SetTitle(title)
	d.SetTransientFor(parent)
	d.SetParent(parent)
	d.SetWindow(parent)
	if len(argv) > 0 {
		d.AddButtons(argv...)
	}
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
	d.action.Show()
	vbox.PackEnd(d.action, false, true, 0)
	d.done = make(chan bool, 1)
	d.response = ResponseNone
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
	d.Connect(SignalCdkEvent, DialogEventHandle, d.event)
	d.Connect(SignalInvalidate, DialogInvalidateHandle, d.invalidate)
	d.Connect(SignalResize, DialogResizeHandle, d.resize)
	d.Connect(SignalDraw, DialogDrawHandle, d.draw)
	d.widgets = make(map[ResponseType][]Widget)
	return false
}

// Build provides customizations to the Buildable system for Dialog Widgets.
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
func (d *CDialog) Run() (response chan ResponseType) {
	d.Show()
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
	if d.defResponse != ResponseNone {
		if ab, ok := d.widgets[d.defResponse]; ok {
			last := len(ab) - 1
			if ab[last] != nil {
				ab[last].GrabFocus()
			}
		}
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

// Response emits the response signal with the given response ID. Used to
// indicate that the user has responded to the dialog in some way; typically
// either you or Run will be monitoring the ::response signal and take
// appropriate action.
//
// Parameters:
// 	responseId	ResponseType identifier
func (d *CDialog) Response(responseId ResponseType) {
	d.Emit(SignalResponse, responseId)
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
		var responseId ResponseType
		if responseId, ok = argv[i+1].(ResponseType); !ok {
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
func (d *CDialog) AddActionWidget(child Widget, responseId ResponseType) {
	child.Connect(SignalActivate, DialogActivateHandle, func(data []interface{}, argv ...interface{}) enums.EventFlag {
		d.LogDebug("responding with: %v", responseId)
		d.Response(responseId)
		return enums.EVENT_STOP
	})
	d.action.PackStart(child, false, false, 0)
	d.widgets[responseId] = append(d.widgets[responseId], child)
}

// AddSecondaryActionWidget is the same as AddActionWidget with the exception of
// adding the given Widget to the secondary action Button grouping instead of
// the primary grouping as with AddActionWidget.
func (d *CDialog) AddSecondaryActionWidget(child Widget, responseId ResponseType) {
	child.Connect(SignalActivate, DialogActivateHandle, func(data []interface{}, argv ...interface{}) enums.EventFlag {
		d.LogDebug("responding with: %v", responseId)
		d.Response(responseId)
		return enums.EVENT_STOP
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
func (d *CDialog) SetDefaultResponse(responseId ResponseType) {
	d.defResponse = responseId
}

// SetResponseSensitive calls Widget.SetSensitive for each widget in the
// Dialog's action area with the given Responsetype. A convenient way to
// sensitize/desensitize Dialog Buttons.
//
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

// GetResponseForWidget is a convenience method for looking up the ResponseType
// associated with the given Widget in the Dialog action area. Returns
// ResponseNone if the Widget is not found in the action area of the Dialog.
//
// Parameters:
// 	widget	a widget in the action area of dialog
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

// GetWidgetForResponse returns the last Widget Button that uses the given
// ResponseType in the action area of a Dialog.
//
// Parameters:
// 	responseId	the response ID used by the dialog widget
func (d *CDialog) GetWidgetForResponse(responseId ResponseType) (value Widget) {
	if widgets, ok := d.widgets[responseId]; ok {
		if last := len(widgets) - 1; last > -1 {
			value = widgets[last]
		}
	}
	return
}

// GetActionArea returns the action area ButtonBox of a Dialog instance.
func (d *CDialog) GetActionArea() (value ButtonBox) {
	return d.action
}

// GetContentArea returns the content area VBox of a Dialog instance.
func (d *CDialog) GetContentArea() (value VBox) {
	return d.content
}

// Show ensures that the Dialog, content and action areas are all set to VISIBLE
func (d *CDialog) Show() {
	d.SetFlags(VISIBLE)
	d.CWindow.Show()
	d.content.Show()
	d.action.Show()
}

// ShowAll calls ShowAll upon the Dialog, content area, action area and all the
// action Widget children.
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

// Destroy hides the Dialog, removes it from any transient Window associations,
// removes the Dialog from the Display and finally emits the destroy-event
// signal.
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

// TODO: set-size-request makes dialog window size, truncated by actual size
// TODO: local canvas / child alloc and origin issues

func (d *CDialog) event(data []interface{}, argv ...interface{}) enums.EventFlag {
	if evt, ok := argv[1].(cdk.Event); ok {
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
					point := ptypes.NewPoint2I(e.Position())
					point.AddPoint(d.GetOrigin())
					if mw := child.GetWidgetAt(point); mw != nil {
						if ms, ok := mw.(Sensitive); ok && ms.IsSensitive() && ms.IsVisible() {
							if f := ms.ProcessEvent(evt); f == enums.EVENT_STOP {
								return enums.EVENT_STOP
							}
						}
					}
				}
			}
		}
		// do not block parent event handlers with enums.EVENT_STOP, always
		// allow event-pass-through so that normal Window events can be handled
		// without duplicating the features here
	}
	return enums.EVENT_PASS
}

func (d *CDialog) invalidate(data []interface{}, argv ...interface{}) enums.EventFlag {
	// d.rebuildFocusChain()
	origin := d.GetOrigin()
	alloc := d.GetAllocation()
	if err := memphis.ConfigureSurface(d.ObjectID(), origin, alloc, d.GetThemeRequest().Content.Normal); err != nil {
		d.LogErr(err)
	}
	return enums.EVENT_STOP
}

func (d *CDialog) resize(data []interface{}, argv ...interface{}) enums.EventFlag {
	region := d.getDialogRegion().NewClone()
	if tf := d.GetTransientFor(); tf != nil {
		if dm := d.GetDisplay(); dm != nil {
			dm.SetWindowOverlayRegion(tf.ObjectID(), d.ObjectID(), *region)
		}
	}
	d.SetOrigin(region.X, region.Y)
	d.SetAllocation(region.Size())

	if child := d.GetChild(); child != nil {
		local := ptypes.MakePoint2I(1, 1)
		child.SetOrigin(region.X+local.X, region.Y+local.Y)
		alloc := region.Size().NewClone()
		alloc.Sub(2, 2)
		child.SetAllocation(*alloc)
		child.Resize()
		if err := memphis.ConfigureSurface(child.ObjectID(), local, *alloc, child.GetThemeRequest().Content.Normal); err != nil {
			child.LogErr(err)
		}
	}

	d.Invalidate()
	return enums.EVENT_STOP
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
		// _ = d.GetVBox().SetBoolProperty(cdk.PropertyDebug, true)
		// _ = d.content.SetBoolProperty(cdk.PropertyDebug, true)
		// _ = d.action.SetBoolProperty(cdk.PropertyDebug, true)
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

const DialogEventHandle = "dialog-event-handler"

const DialogInvalidateHandle = "dialog-invalidate-handler"

const DialogResizeHandle = "dialog-resize-handler"

const DialogDrawHandle = "dialog-draw-handler"

const DialogActivateHandle = "dialog-activate-handler"
