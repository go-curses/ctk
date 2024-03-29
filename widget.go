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

package ctk

import (
	"github.com/gofrs/uuid"

	"github.com/go-curses/cdk"
	cenums "github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/cdk/lib/paint"
	"github.com/go-curses/cdk/lib/ptypes"
	"github.com/go-curses/cdk/lib/sync"
	"github.com/go-curses/cdk/memphis"

	"github.com/go-curses/ctk/lib/enums"
)

const (
	TypeWidget        cdk.CTypeTag    = "ctk-widget"
	TooltipColorStyle paint.StyleName = "tooltip-color"
	TooltipColorTheme paint.ThemeName = "tooltip-color"
)

func init() {
	_ = cdk.TypesManager.AddType(TypeWidget, nil)

	style := paint.StyleDefault.Normal().
		Background(paint.ColorYellow).
		Foreground(paint.ColorDarkSlateBlue).
		Dim(false)
	paint.RegisterStyle(TooltipColorStyle, style)

	borders, _ := paint.GetDefaultBorderRunes(paint.StockBorder)
	arrows, _ := paint.GetArrows(paint.StockArrow)

	tooltipThemeAspect := paint.ThemeAspect{
		Normal:      style,
		Selected:    style,
		Active:      style,
		Prelight:    style,
		Insensitive: style,
		FillRune:    paint.DefaultFillRune,
		BorderRunes: borders,
		ArrowRunes:  arrows,
		Overlay:     false,
	}

	paint.RegisterTheme(TooltipColorTheme, paint.Theme{
		Content: tooltipThemeAspect,
		Border:  tooltipThemeAspect,
	})
}

// Widget Hierarchy:
//	Object
//	  +- Widget
//	    +- Container
//	    +- Misc
//	    +- Calendar
//	    +- CellView
//	    +- DrawingArea
//	    +- Entry
//	    +- Ruler
//	    +- Range
//	    +- Separator
//	    +- HSV
//	    +- Invisible
//	    +- OldEditable
//	    +- Preview
//	    +- Progress
//
// Widget is the base class all widgets in CTK derive from. It manages the
// widget lifecycle, states and style.
type Widget interface {
	Object

	Unparent()
	Map()
	Unmap()
	IsMapped() (mapped bool)
	Show()
	Hide()
	GetRegion() (region ptypes.Region)
	LockDraw()
	UnlockDraw()
	LockEvent()
	UnlockEvent()
	AddAccelerator(accelSignal string, accelGroup AccelGroup, accelKey int, accelMods enums.ModifierType, accelFlags enums.AccelFlags)
	RemoveAccelerator(accelGroup AccelGroup, accelKey int, accelMods enums.ModifierType) (value bool)
	SetAccelPath(accelPath string, accelGroup AccelGroup)
	CanActivateAccel(signalId int) (value bool)
	Activate() (value bool)
	Reparent(parent Widget)
	IsFocus() (value bool)
	GrabFocus()
	GrabDefault()
	SetSensitive(sensitive bool)
	CssFullPath() (selector string)
	CssState() (state enums.StateType)
	SetParent(parent Widget)
	GetParentWindow() (value Window)
	SetEvents(events cdk.EventMask)
	AddEvents(events cdk.EventMask)
	GetToplevel() (value Widget)
	GetAncestor(widgetType cdk.CTypeTag) (value Widget)
	GetEvents() (value cdk.EventMask)
	GetPointer(x int, y int)
	IsAncestor(ancestor Widget) (value bool)
	TranslateCoordinates(destWidget Widget, srcX int, srcY int, destX int, destY int) (value bool)
	HideOnDelete() (value bool)
	SetDirection(dir enums.TextDirection)
	GetDirection() (value enums.TextDirection)
	SetDefaultDirection(dir enums.TextDirection)
	GetDefaultDirection() (value enums.TextDirection)
	Path() (path string)
	ClassPath(pathLength int, path string, pathReversed string)
	GetCompositeName() (value string)
	SetAppPaintable(appPaintable bool)
	SetDoubleBuffered(doubleBuffered bool)
	SetRedrawOnAllocate(redrawOnAllocate bool)
	SetCompositeName(name string)
	SetScrollAdjustments(hadjustment Adjustment, vadjustment Adjustment) (value bool)
	Draw() cenums.EventFlag
	MnemonicActivate(groupCycling bool) (value bool)
	SendExpose(event cdk.Event) (value int)
	SendFocusChange(event cdk.Event) (value bool)
	ChildFocus(direction enums.DirectionType) (value bool)
	ChildNotify(childProperty string)
	FreezeChildNotify()
	GetChildVisible() (value bool)
	GetParent() (value Widget)
	GetAllParents() (parents []Widget)
	GetDisplay() (value cdk.Display)
	GetRootWindow() (value Window)
	GetScreen() (value cdk.Display)
	HasScreen() (value bool)
	GetSizeRequest() (width, height int)
	SizeRequest() ptypes.Rectangle
	SetSizeRequest(width, height int)
	SetNoShowAll(noShowAll bool)
	GetNoShowAll() (value bool)
	AddMnemonicLabel(label Widget)
	RemoveMnemonicLabel(label Widget)
	ErrorBell()
	KeynavFailed(direction enums.DirectionType) (value bool)
	GetTooltipMarkup() (value string)
	SetTooltipMarkup(markup string)
	GetTooltipText() (value string)
	SetTooltipText(text string)
	GetTooltipWindow() (value Window)
	SetTooltipWindow(customWindow Window)
	GetHasTooltip() (value bool)
	SetHasTooltip(hasTooltip bool)
	GetWindow() (window Window)
	GetAppPaintable() (value bool)
	GetCanDefault() (value bool)
	SetCanDefault(canDefault bool)
	GetCanFocus() (value bool)
	SetCanFocus(canFocus bool)
	GetHasWindow() (ok bool)
	GetSensitive() (value bool)
	IsSensitive() bool
	GetVisible() (value bool)
	SetVisible(visible bool)
	HasDefault() (value bool)
	HasFocus() (value bool)
	HasGrab() (value bool)
	IsDrawable() (value bool)
	IsToplevel() (value bool)
	SetWindow(window Window)
	SetReceivesDefault(receivesDefault bool)
	GetReceivesDefault() (value bool)
	SetRealized(realized bool)
	GetRealized() (value bool)
	SetMapped(mapped bool)
	GetMapped() (value bool)
	GetThemeRequest() (theme paint.Theme)
	GetState() (value enums.StateType)
	SetState(state enums.StateType)
	HasState(s enums.StateType) bool
	UnsetState(state enums.StateType)
	GetFlags() enums.WidgetFlags
	HasFlags(f enums.WidgetFlags) bool
	UnsetFlags(v enums.WidgetFlags)
	SetFlags(v enums.WidgetFlags)
	IsParentFocused() bool
	IsFocused() bool
	CanFocus() bool
	IsDefault() bool
	CanDefault() bool
	IsVisible() bool
	HasEventFocus() bool
	GrabEventFocus()
	ReleaseEventFocus()
	GetTopParent() (parent Widget)
	GetWidgetAt(p *ptypes.Point2I) Widget
	PushCompositeChild(child Widget)
	PopCompositeChild(child Widget)
	GetCompositeChildren() []Widget
	RenderFrozen() bool
	RenderFreeze()
	RenderThaw()
	RequestDrawAndShow()
	RequestDrawAndSync()
}

var _ Widget = (*CWidget)(nil)

// The CWidget structure implements the Widget interface and is exported
// to facilitate type embedding with custom implementations. No member variables
// are exported as the interface methods are the only intended means of
// interacting with Widget objects.
type CWidget struct {
	CObject

	renderFrozen int
	composites   []Widget
	parent       Widget
	state        enums.StateType
	flags        enums.WidgetFlags
	flagsLock    *sync.RWMutex
	drawLock     *sync.Mutex
	eventLock    *sync.Mutex

	tooltipWindow      Window
	tooltipTimer       uuid.UUID
	tooltipBrowseTimer uuid.UUID
}

// Init initializes a Widget object. This must be called at least once to
// set up the necessary defaults and allocate any memory structures. Calling
// this more than once is safe though unnecessary. Only the first call will
// result in any effect upon the Widget instance. Init is used in the
// NewWidget constructor and only necessary when implementing a derivative
// Widget type.
func (w *CWidget) Init() (already bool) {
	if w.InitTypeItem(TypeWidget, w) {
		return true
	}
	w.CObject.Init()
	_ = w.InstallProperty(PropertyAppPaintable, cdk.BoolProperty, true, false)
	_ = w.InstallProperty(PropertyCanDefault, cdk.BoolProperty, true, false)
	_ = w.InstallProperty(PropertyCanFocus, cdk.BoolProperty, true, false)
	_ = w.InstallProperty(PropertyCompositeChild, cdk.BoolProperty, true, false)
	_ = w.InstallProperty(PropertyDoubleBuffered, cdk.BoolProperty, true, true)
	_ = w.InstallProperty(PropertyEvents, cdk.StructProperty, true, cdk.EVENT_MASK_NONE)
	_ = w.InstallProperty(PropertyHasDefault, cdk.BoolProperty, true, false)
	_ = w.InstallProperty(PropertyHasFocus, cdk.BoolProperty, true, false)
	_ = w.InstallProperty(PropertyHasTooltip, cdk.BoolProperty, true, false)
	_ = w.InstallProperty(PropertyHeightRequest, cdk.IntProperty, true, -1)
	_ = w.InstallProperty(PropertyIsFocus, cdk.BoolProperty, true, false)
	_ = w.InstallProperty(PropertyNoShowAll, cdk.BoolProperty, true, false)
	_ = w.InstallProperty(PropertyParent, cdk.StructProperty, true, nil)
	_ = w.InstallProperty(PropertyReceivesDefault, cdk.BoolProperty, true, false)
	_ = w.InstallProperty(PropertySensitive, cdk.BoolProperty, true, true)
	_ = w.InstallProperty(PropertyTooltipMarkup, cdk.StringProperty, true, "")
	_ = w.InstallProperty(PropertyTooltipText, cdk.StringProperty, true, "")
	_ = w.InstallProperty(PropertyVisible, cdk.BoolProperty, true, false)
	_ = w.InstallProperty(PropertyWidthRequest, cdk.IntProperty, true, -1)
	_ = w.InstallProperty(PropertyWindow, cdk.StructProperty, true, nil)

	theme := paint.GetDefaultColorTheme()

	getDefContentColors := func(state enums.StateType) (fg, bg paint.Color) {
		switch state {
		case enums.StateNormal:
			fg, bg, _ = theme.Content.Normal.Decompose()
		case enums.StateActive:
			fg, bg, _ = theme.Content.Active.Decompose()
		case enums.StatePrelight:
			fg, bg, _ = theme.Content.Prelight.Decompose()
		case enums.StateSelected:
			fg, bg, _ = theme.Content.Selected.Decompose()
		case enums.StateInsensitive:
			fg, bg, _ = theme.Content.Insensitive.Decompose()
		}
		return
	}
	getDefBorderColors := func(state enums.StateType) (fg, bg paint.Color) {
		switch state {
		case enums.StateNormal:
			fg, bg, _ = theme.Border.Normal.Decompose()
		case enums.StateActive:
			fg, bg, _ = theme.Border.Active.Decompose()
		case enums.StatePrelight:
			fg, bg, _ = theme.Border.Prelight.Decompose()
		case enums.StateSelected:
			fg, bg, _ = theme.Border.Selected.Decompose()
		case enums.StateInsensitive:
			fg, bg, _ = theme.Border.Insensitive.Decompose()
		}
		return
	}

	for _, state := range []enums.StateType{enums.StateNormal, enums.StateActive, enums.StatePrelight, enums.StateSelected, enums.StateInsensitive} {
		_ = w.InstallCssProperty(CssPropertyClass, state, cdk.StringProperty, true, "")
		_ = w.InstallCssProperty(CssPropertyWidth, state, cdk.IntProperty, true, -1)
		_ = w.InstallCssProperty(CssPropertyHeight, state, cdk.IntProperty, true, -1)
		fg, bg := getDefContentColors(state)
		_ = w.InstallCssProperty(CssPropertyColor, state, cdk.ColorProperty, true, fg)
		_ = w.InstallCssProperty(CssPropertyBackgroundColor, state, cdk.ColorProperty, true, bg)
		fg, bg = getDefBorderColors(state)
		_ = w.InstallCssProperty(CssPropertyBorderColor, state, cdk.ColorProperty, true, fg)
		_ = w.InstallCssProperty(CssPropertyBorderBackgroundColor, state, cdk.ColorProperty, true, bg)
		_ = w.InstallCssProperty(CssPropertyBold, state, cdk.BoolProperty, true, false)
		_ = w.InstallCssProperty(CssPropertyBlink, state, cdk.BoolProperty, true, false)
		_ = w.InstallCssProperty(CssPropertyReverse, state, cdk.BoolProperty, true, false)
		_ = w.InstallCssProperty(CssPropertyUnderline, state, cdk.BoolProperty, true, false)
		_ = w.InstallCssProperty(CssPropertyDim, state, cdk.BoolProperty, true, false)
		_ = w.InstallCssProperty(CssPropertyItalic, state, cdk.BoolProperty, true, false)
		_ = w.InstallCssProperty(CssPropertyStrike, state, cdk.BoolProperty, true, false)
	}

	w.renderFrozen = 0
	w.composites = make([]Widget, 0)
	w.flagsLock = &sync.RWMutex{}
	w.drawLock = &sync.Mutex{}
	w.eventLock = &sync.Mutex{}
	w.state = enums.StateNormal
	w.flags = enums.NULL_WIDGET_FLAG
	w.tooltipWindow = nil
	w.tooltipTimer = uuid.Nil
	w.tooltipBrowseTimer = uuid.Nil
	w.SetTheme(theme)

	w.Connect(SignalLostFocus, WidgetLostFocusHandle, w.lostFocus)
	w.Connect(SignalGainedFocus, WidgetGainedFocusHandle, w.gainedFocus)
	w.Connect(SignalEnter, WidgetEnterHandle, w.enter)
	w.Connect(SignalLeave, WidgetLeaveHandle, w.leave)
	return false
}

func (w *CWidget) openTooltip() {
	settings := GetDefaultSettings()
	if settings.GetEnableTooltips() {
		if w.GetHasTooltip() {
			if w.tooltipBrowseTimer != uuid.Nil {
				cdk.StopTimeout(w.tooltipBrowseTimer)
				w.tooltipBrowseTimer = uuid.Nil
			}
			timeout := settings.GetTooltipTimeout()
			w.Lock()
			w.tooltipTimer = cdk.AddTimeout(
				timeout,
				w.openTooltipHandler,
			)
			w.LogDebug("openTooltip setting up timeout handler: %v (%v)", timeout, w.tooltipTimer)
			w.Unlock()
		}
	}
}

func (w *CWidget) openTooltipHandler() cenums.EventFlag {
	if display := w.GetDisplay(); display != nil && w.GetHasTooltip() {
		tooltipWindow := w.GetTooltipWindow()
		if screen := display.Screen(); screen != nil {
			if w.tooltipBrowseTimer != uuid.Nil {
				cdk.StopTimeout(w.tooltipBrowseTimer)
				w.tooltipBrowseTimer = uuid.Nil
			}
			position, _ := display.CursorPosition()
			tw, th := w.tooltipTextBufferInfo()
			if w.HasPoint(&position) {
				if tw > -1 && th > -1 {
					sw, sh := screen.Size()
					if position.X+2+tw < sw {
						position.X += 2
					}
					if position.Y+1+th < sh {
						position.Y += 1
					}
					if position.X+tw >= sw {
						// too far to the right
						position.X = sw - tw
					}
					if position.Y+th >= sh {
						// too far down
						position.Y = sh - th
					}
					tooltipWindow.Move(position.X, position.Y)
					tooltipWindow.SetAllocation(ptypes.MakeRectangle(tw, th))
					w.LogDebug("opening tooltip: %v", tooltipWindow.ObjectInfo())
					settings := GetDefaultSettings()
					browseModeTimeout := settings.GetTooltipBrowseModeTimeout()
					if markup := w.GetTooltipMarkup(); markup != "" {
						tooltipWindow.Resize()
						tooltipWindow.Show()
						w.tooltipBrowseTimer = cdk.AddTimeout(browseModeTimeout, func() cenums.EventFlag {
							w.closeTooltip()
							return cenums.EVENT_STOP
						})
						if window := w.GetWindow(); window != nil {
							window.Connect(SignalCdkEvent, WidgetTooltipWindowEventHandle, w.processTooltipWindowEvent)
						}
						return cenums.EVENT_STOP
					}
					if text := w.GetTooltipText(); text != "" {
						tooltipWindow.Resize()
						tooltipWindow.Show()
						w.tooltipBrowseTimer = cdk.AddTimeout(browseModeTimeout, func() cenums.EventFlag {
							w.closeTooltip()
							return cenums.EVENT_STOP
						})
						if window := w.GetWindow(); window != nil {
							window.Connect(SignalCdkEvent, WidgetTooltipWindowEventHandle, w.processTooltipWindowEvent)
						}
						return cenums.EVENT_STOP
					}
				}

			}
		}
		tooltipWindow.Hide()
	}
	return cenums.EVENT_STOP
}

func (w *CWidget) processTooltipWindowEvent(data []interface{}, argv ...interface{}) cenums.EventFlag {
	if len(argv) >= 2 {
		if evt, ok := argv[1].(cdk.Event); ok {
			switch evt.(type) {
			case *cdk.EventMouse, *cdk.EventKey:
				w.closeTooltip()
			}
		}
	}
	return cenums.EVENT_PASS
}

func (w *CWidget) closeTooltip() {
	if window := w.GetWindow(); window != nil {
		_ = window.Disconnect(SignalCdkEvent, WidgetTooltipWindowEventHandle)
	}
	if w.tooltipBrowseTimer != uuid.Nil {
		cdk.StopTimeout(w.tooltipBrowseTimer)
		w.tooltipBrowseTimer = uuid.Nil
	}
	if w.tooltipTimer != uuid.Nil {
		w.LogDebug("closing tooltip")
		cdk.StopTimeout(w.tooltipTimer)
		w.tooltipTimer = uuid.Nil
	}
	if w.tooltipWindow != nil {
		w.tooltipWindow.Hide()
	}
}

func (w *CWidget) newTooltipWindow() (tooltipWindow Window) {
	tooltipWindow = NewWindow()
	tooltipWindow.SetWindowType(cenums.WINDOW_POPUP)
	tooltipWindow.SetFlags(enums.TOPLEVEL)
	tooltipWindow.SetDecorated(false)
	theme, _ := paint.GetTheme(TooltipColorTheme)
	tooltipWindow.SetTheme(theme)
	tooltipWindow.Connect(SignalCdkEvent, "widget-tooltip-window-resize-handler", w.tooltipEvent)
	tooltipWindow.Connect(SignalResize, "widget-tooltip-window-resize-handler", w.tooltipResize)
	tooltipWindow.Connect(SignalDraw, "widget-tooltip-window-draw-handler", w.tooltipDraw)
	return
}

func (w *CWidget) tooltipTextBufferInfo() (longestLine, lineCount int) {
	longestLine, lineCount = -1, -1
	style := w.GetThemeRequest().Content.Normal
	var tb memphis.TextBuffer
	if markup := w.GetTooltipMarkup(); markup != "" {
		if m, err := memphis.NewMarkup(markup, style); err != nil {
			w.LogErr(err)
			return
		} else {
			tb = m.TextBuffer(false)
		}
	} else if text := w.GetTooltipText(); text != "" {
		tb = memphis.NewTextBuffer(text, style, false)
	} else {
		return
	}
	longestLine, lineCount = tb.PlainTextInfo(cenums.WRAP_WORD, false, cenums.JUSTIFY_LEFT, -1)
	return
}

func (w *CWidget) tooltipEvent(data []interface{}, argv ...interface{}) cenums.EventFlag {
	if _, event, ok := ArgvSignalEvent(argv...); ok {
		if tooltipWindow := w.GetTooltipWindow(); tooltipWindow != nil {
			switch e := event.(type) {
			case *cdk.EventMouse:
				point := e.Point2I().NewClone()
				if !w.HasPoint(point) {
					w.closeTooltip()
					return cenums.EVENT_STOP
				} else {
					x, y := e.Position()
					mw, mh := w.tooltipTextBufferInfo()
					if x > mw {
						x -= mw
					} else {
						x = mw - x
					}
					if y > mh {
						y -= mh
					} else {
						y = mh - y
					}
					tooltipWindow.Move(x, y)
				}
			}
		}
	}
	return cenums.EVENT_PASS
}

func (w *CWidget) tooltipResize(data []interface{}, argv ...interface{}) cenums.EventFlag {
	alloc := ptypes.MakeRectangle(w.tooltipTextBufferInfo())
	if tooltipWindow := w.GetTooltipWindow(); tooltipWindow != nil {
		tooltipWindow.SetAllocation(alloc)
	}
	return cenums.EVENT_STOP
}

func (w *CWidget) tooltipDraw(data []interface{}, argv ...interface{}) cenums.EventFlag {
	if tooltipWindow := w.GetTooltipWindow(); tooltipWindow != nil {

		if surface, ok := argv[1].(*memphis.CSurface); ok {
			size := surface.GetSize()
			if !w.IsVisible() || size.W == 0 || size.H == 0 {
				w.LogDebug("not visible, zero width or zero height")
				return cenums.EVENT_STOP
			}

			theme := tooltipWindow.GetThemeRequest()
			style := theme.Content.Normal

			var tb memphis.TextBuffer
			if markup := w.GetTooltipMarkup(); markup != "" {
				if m, err := memphis.NewMarkup(markup, style); err != nil {
					w.LogErr(err)
					return cenums.EVENT_STOP
				} else {
					tb = m.TextBuffer(false)
				}
			} else if text := w.GetTooltipText(); text != "" {
				tb = memphis.NewTextBuffer(text, style, false)
			} else {
				w.LogError("tooltip window is drawing without text?!")
				return cenums.EVENT_STOP
			}

			surface.Fill(theme)

			if tb != nil {
				tb.Draw(surface, false, cenums.WRAP_WORD, false, cenums.JUSTIFY_LEFT, cenums.ALIGN_TOP)
			}

			if debug, _ := w.GetBoolProperty(cdk.PropertyDebug); debug {
				surface.DebugBox(paint.ColorNavy, tooltipWindow.ObjectInfo())
			}
		}
	}
	return cenums.EVENT_STOP
}

// Destroy a widget. Equivalent to DestroyObject. When a widget is destroyed, it
// will break any references it holds to other objects. If the widget is
// inside a container, the widget will be removed from the container. If the
// widget is a toplevel (derived from Window), it will be removed from the
// list of toplevels, and the reference CTK holds to it will be removed.
// Removing a widget from its container or the list of toplevels results in
// the widget being finalized, unless you've added additional references to
// the widget with g_object_ref. In most cases, only toplevel widgets
// (windows) require explicit destruction, because when you destroy a
// toplevel its children will be destroyed as well.
func (w *CWidget) Destroy() {
	w.closeTooltip()
	w.Emit(SignalDestroyEvent, w)
	w.DisconnectAll()
	w.Hide()
	if w.tooltipWindow != nil {
		w.tooltipWindow.Destroy()
	}
	if err := w.DestroyObject(); err != nil {
		w.LogErr(err)
	}
}

// Unparent is only for use in widget implementations. Should be called by
// implementations of the remove method on Container, to dissociate a child from
// the container.
func (w *CWidget) Unparent() {
	if w.parent != nil {
		if parent, ok := w.parent.Self().(Container); ok {
			if f := w.Emit(SignalUnparent, w, w.parent); f == cenums.EVENT_PASS {
				parent.Remove(w)
			}
		}
	}
}

func (w *CWidget) Map() {
	if !w.IsMapped() {
		region := w.GetRegion()
		// w.LockDraw()
		_ = memphis.MakeConfigureSurface(
			w.ObjectID(),
			region.Origin(),
			region.Size(),
			w.GetThemeRequest().Content.Normal,
		)
		w.SetFlags(enums.MAPPED)
		// w.UnlockDraw()
		w.Emit(SignalMap, w, region)
	}
}

func (w *CWidget) Unmap() {
	if w.IsMapped() {
		memphis.RemoveSurface(w.ObjectID())
		w.UnsetFlags(enums.MAPPED)
		w.Emit(SignalUnmap, w)
	}
}

func (w *CWidget) IsMapped() (mapped bool) {
	_, err := memphis.GetSurface(w.ObjectID())
	mapped = err == nil && w.HasFlags(enums.MAPPED)
	return
}

// Show flags a widget to be displayed. Any widget that isn't shown will not
// appear on the screen. If you want to show all the widgets in a container,
// it's easier to call ShowAll on the container, instead of
// individually showing the widgets. Remember that you have to show the
// containers containing a widget, in addition to the widget itself, before
// it will appear onscreen. When a toplevel container is shown, it is
// immediately realized and mapped; other shown widgets are realized and
// mapped when their toplevel container is realized and mapped.
func (w *CWidget) Show() {
	if w.HasFlags(enums.APP_PAINTABLE) {
		if !w.HasFlags(enums.VISIBLE) {
			w.SetFlags(enums.VISIBLE)
			w.Emit(SignalShow, w)
			w.Invalidate()
		}
	}
}

// Hide reverses the effects of Show, causing the widget to be hidden (invisible
// to the user).
func (w *CWidget) Hide() {
	if w.HasFlags(enums.APP_PAINTABLE) {
		if w.HasFlags(enums.VISIBLE) {
			w.closeTooltip()
			w.Unmap()
			w.UnsetFlags(enums.VISIBLE)
			w.Emit(SignalHide, w)
			w.Invalidate()
		}
	}
}

// GetRegion returns the current origin and allocation in a Region type, taking
// any positive SizeRequest set.
func (w *CWidget) GetRegion() (region ptypes.Region) {
	region = w.CObject.GetRegion()
	req := w.SizeRequest()
	if req.W > 0 && region.W > req.W {
		region.W = req.W
	}
	if req.H > 0 && region.H > req.H {
		region.H = req.H
	}
	return
}

func (w *CWidget) LockDraw() {
	w.drawLock.Lock()
}

func (w *CWidget) UnlockDraw() {
	w.drawLock.Unlock()
}

func (w *CWidget) LockEvent() {
	w.eventLock.Lock()
}

func (w *CWidget) UnlockEvent() {
	w.eventLock.Unlock()
}

func (w *CWidget) RenderFrozen() bool {
	w.RLock()
	defer w.RUnlock()
	return w.renderFrozen > 0
}

func (w *CWidget) RenderFreeze() {
	if !w.RenderFrozen() {
		w.PassSignal(SignalInvalidate, SignalDraw)
	}
	w.Lock()
	w.renderFrozen += 1
	w.Unlock()
}

func (w *CWidget) RenderThaw() {
	if w.RenderFrozen() {
		w.Lock()
		w.renderFrozen -= 1
		w.Unlock()
	}
	if !w.RenderFrozen() {
		w.ResumeSignal(SignalInvalidate)
		w.Resize()
		w.ResumeSignal(SignalDraw)
	}
}

// Installs an accelerator for this widget in accel_group that causes
// accel_signal to be emitted if the accelerator is activated. The
// accel_group needs to be added to the widget's toplevel via
// WindowAddAccelGroup, and the signal must be of type G_RUN_ACTION.
// Accelerators added through this function are not user changeable during
// runtime. If you want to support accelerators that can be changed by the
// user, use AccelMapAddEntry and SetAccelPath or
// MenuItemSetAccelPath instead.
// Parameters:
// 	widget	widget to install an accelerator on
// 	accelSignal	widget signal to emit on accelerator activation
// 	accelGroup	accel group for this widget, added to its toplevel
// 	accelKey	GDK keyval of the accelerator
// 	accelMods	modifier key combination of the accelerator
// 	accelFlags	flag accelerators, e.g. GTK_ACCEL_VISIBLE
//
// Method stub, unimplemented
func (w *CWidget) AddAccelerator(accelSignal string, accelGroup AccelGroup, accelKey int, accelMods enums.ModifierType, accelFlags enums.AccelFlags) {
}

// Removes an accelerator from widget , previously installed with
// AddAccelerator.
// Parameters:
// 	widget	widget to install an accelerator on
// 	accelGroup	accel group for this widget
// 	accelKey	GDK keyval of the accelerator
// 	accelMods	modifier key combination of the accelerator
// 	returns	whether an accelerator was installed and could be removed
//
// Method stub, unimplemented
func (w *CWidget) RemoveAccelerator(accelGroup AccelGroup, accelKey int, accelMods enums.ModifierType) (value bool) {
	return false
}

// Given an accelerator group, accel_group , and an accelerator path,
// accel_path , sets up an accelerator in accel_group so whenever the key
// binding that is defined for accel_path is pressed, widget will be
// activated. This removes any accelerators (for any accelerator group)
// installed by previous calls to SetAccelPath. Associating
// accelerators with paths allows them to be modified by the user and the
// modifications to be saved for future use. (See AccelMapSave.) This
// function is a low level function that would most likely be used by a menu
// creation system like UIManager. If you use UIManager, setting up
// accelerator paths will be done automatically. Even when you you aren't
// using UIManager, if you only want to set up accelerators on menu items
// MenuItemSetAccelPath provides a somewhat more convenient
// interface. Note that accel_path string will be stored in a GQuark.
// Therefore, if you pass a static string, you can save some memory by
// interning it first with g_intern_static_string.
// Parameters:
// 	accelPath	path used to look up the accelerator.
// 	accelGroup	a AccelGroup.
//
// Method stub, unimplemented
func (w *CWidget) SetAccelPath(accelPath string, accelGroup AccelGroup) {}

// Determines whether an accelerator that activates the signal identified by
// signal_id can currently be activated. This is done by emitting the
// can-activate-accel signal on widget ; if the signal isn't overridden
// by a handler or in a derived widget, then the default check is that the
// widget must be sensitive, and the widget and all its ancestors mapped.
// Parameters:
// 	signalId	the ID of a signal installed on widget
//
// Returns:
// 	TRUE if the accelerator can be activated.
//
// Method stub, unimplemented
func (w *CWidget) CanActivateAccel(signalId int) (value bool) {
	return false
}

// For widgets that can be "activated" (buttons, menu items, etc.) this
// function activates them. Activation is what happens when you press Enter
// on a widget during key navigation. If widget isn't activatable, the
// function returns FALSE.
// Returns:
// 	TRUE if the widget was activatable
func (w *CWidget) Activate() (value bool) {
	w.closeTooltip()
	if w.IsVisible() && w.IsSensitive() {
		if f := w.Emit(SignalActivate, w); f == cenums.EVENT_STOP {
			value = true
		}
	}
	return
}

// Move the Widget to the given container, removing itself first from any other
// container that was currently holding it. This method emits a reparent signal
// initially and if the listeners return EVENT_PAS, the change is applied
// Parameters:
// 	newParent	a Container to move the widget into
//
// Emits: SignalReparent, Argv=[Widget instance, new parent]
func (w *CWidget) Reparent(parent Widget) {
	if r := w.Emit(SignalReparent, w, parent); r == cenums.EVENT_PASS {
		if pc, ok := parent.Self().(Container); ok {
			if w.parent != nil {
				if pc, ok := w.parent.Self().(Container); ok {
					pc.Remove(w)
				}
			}
			pc.Add(w)
		}
	}
}

// Determines if the widget is the focus widget within its toplevel. (This
// does not mean that the HAS_FOCUS flag is necessarily set; HAS_FOCUS will
// only be set if the toplevel widget additionally has the global input
// focus.)
// Returns:
// 	TRUE if the widget is the focus widget.
// Returns TRUE if the Widget instance is currently the focus of it's parent
// Window, FALSE otherwise
func (w *CWidget) IsFocus() (value bool) {
	if window := w.GetWindow(); window != nil {
		if w.CanFocus() {
			if focused := window.GetFocus(); focused != nil {
				if focused.ObjectID() == focused.ObjectID() {
					return true
				}
			}
		}
	}
	return false
}

// Causes widget to have the keyboard focus for the Window it's inside.
// widget must be a focusable widget, such as a Entry; something like
// Frame won't work. More precisely, it must have the GTK_CAN_FOCUS flag
// set. Use SetCanFocus to modify that flag. The widget also
// needs to be realized and mapped. This is indicated by the related signals.
// Grabbing the focus immediately after creating the widget will likely fail
// and cause critical warnings.
// If the Widget instance CanFocus() then take the focus of the associated
// Window. Any previously focused Widget will emit a lost-focus signal and the
// newly focused Widget will emit a gained-focus signal. This method emits a
// grab-focus signal initially and if the listeners return EVENT_PASS, the
// changes are applied
//
// Emits: SignalGrabFocus, Argv=[Widget instance]
// Emits: SignalLostFocus, Argv=[Previous focus Widget instance], From=Previous focus Widget instance
// Emits: SignalGainedFocus, Argv=[Widget instance, previous focus Widget instance]
func (w *CWidget) GrabFocus() {
	if w.CanFocus() && w.IsVisible() && w.IsSensitive() {
		if r := w.Emit(SignalGrabFocus, w); r == cenums.EVENT_PASS {
			if tl := w.GetWindow(); tl != nil {
				if focused := tl.GetFocus(); focused != nil {
					if focused.ObjectID() != w.ObjectID() {
						if err := focused.SetProperty(PropertyHasFocus, false); err != nil {
							focused.LogErr(err)
						}
						focused.UnsetState(enums.StateSelected)
						focused.Emit(SignalLostFocus)
						focused.Invalidate()
					}
				}

				tl.SetFocus(w)
				if err := w.SetProperty(PropertyHasFocus, true); err != nil {
					w.LogErr(err)
				}
				w.SetState(enums.StateSelected)
				w.Emit(SignalGainedFocus)
				w.Invalidate()
			}
		}
	} else {
		w.LogError("cannot grab focus: can't focus, invisible or insensitive")
	}
}

// Causes widget to become the default widget. widget must have the
// GTK_CAN_DEFAULT flag set; typically you have to set this flag yourself by
// calling SetCanDefault (widget , TRUE). The default widget is
// activated when the user presses Enter in a window. Default widgets must be
// activatable, that is, Activate should affect them.
func (w *CWidget) GrabDefault() {}

// Sets the sensitivity of a widget. A widget is sensitive if the user can
// interact with it. Insensitive widgets are "grayed out" and the user can't
// interact with them. Insensitive widgets are known as "inactive",
// "disabled", or "ghosted" in some other toolkits.
// Parameters:
// 	sensitive	TRUE to make the widget sensitive
// Emits: SignalSetSensitive, Argv=[Widget instance, given sensitive bool]
func (w *CWidget) SetSensitive(sensitive bool) {
	if f := w.Emit(SignalSetSensitive, w, sensitive); f == cenums.EVENT_PASS {
		w.closeTooltip()
		if !sensitive {
			w.SetState(enums.StateInsensitive)
		} else {
			w.UnsetState(enums.StateInsensitive)
		}
		if err := w.SetBoolProperty(PropertySensitive, sensitive); err != nil {
			w.LogErr(err)
		}
	}
}

// CssFullPath returns a CSS selector rule for the Widget which includes the
// parent hierarchy in type form.
func (w *CWidget) CssFullPath() (selector string) {
	selector = w.CObject.CssSelector()
	p := w.GetParent()
	for {
		if p != nil && p.ObjectID() != w.ObjectID() {
			selector = p.CssSelector() + " > " + selector
			np := p.GetParent()
			if np != nil && np.ObjectID() == p.ObjectID() {
				break
			}
			p = np
		} else {
			break
		}
	}
	return
}

func (w *CWidget) CssState() (state enums.StateType) {
	if w.IsDrawable() && w.IsVisible() {
		if w.IsSensitive() {
			if w.IsFocus() {
				w.RLock()
				state = w.state.Set(enums.StateSelected)
				w.RUnlock()
			} else {
				w.RLock()
				state = state.Set(enums.StateNormal)
				w.RUnlock()
			}
		} else {
			w.RLock()
			state = w.state.Set(enums.StateInsensitive)
			w.RUnlock()
		}
	}
	return
}

// This function is useful only when implementing subclasses of Container.
// Sets the container as the parent of widget , and takes care of some
// details such as updating the state and style of the child to reflect its
// new location. The opposite function is Unparent.
// Parameters:
// 	parent	parent container
func (w *CWidget) SetParent(parent Widget) {
	if f := w.Emit(SignalSetParent, w, w.parent, parent); f == cenums.EVENT_PASS {
		w.closeTooltip()
		if w.HasFlags(enums.PARENT_SENSITIVE) && w.parent != nil {
			_ = parent.Disconnect(SignalLostFocus, WidgetLostFocusHandle)
			_ = parent.Disconnect(SignalGainedFocus, WidgetGainedFocusHandle)
		}
		if err := w.SetStructProperty(PropertyParent, parent); err != nil {
			w.LogErr(err)
		} else if parent != nil && w.ObjectID() != parent.ObjectID() {
			if w.HasFlags(enums.PARENT_SENSITIVE) {
				parent.Connect(SignalLostFocus, WidgetLostFocusHandle, w.lostFocus)
				parent.Connect(SignalGainedFocus, WidgetGainedFocusHandle, w.gainedFocus)
			}
		}
	}
}

// Gets widget's parent window.
// Returns:
// 	the parent window of widget.
func (w *CWidget) GetParentWindow() (value Window) {
	var ok bool
	p := w.GetParent()
	for p != nil {
		if value, ok = p.(Window); ok {
			return
		}
		p = p.GetParent()
	}
	return
}

// Sets the event mask (see EventMask) for a widget. The event mask
// determines which events a widget will receive. Keep in mind that different
// widgets have different default event masks, and by changing the event mask
// you may disrupt a widget's functionality, so be careful. This function
// must be called while a widget is unrealized. Consider
// AddEvents for widgets that are already realized, or if you
// want to preserve the existing event mask. This function can't be used with
// GTK_NO_WINDOW widgets; to get events on those widgets, place them inside a
// EventBox and receive events on the event box.
// Parameters:
// 	events	event mask
func (w *CWidget) SetEvents(events cdk.EventMask) {
	if err := w.SetStructProperty(PropertyEvents, events); err != nil {
		w.LogErr(err)
	}
}

// Adds the events in the bitfield events to the event mask for widget . See
// SetEvents for details.
// Parameters:
// 	events	an event mask, see EventMask
func (w *CWidget) AddEvents(events cdk.EventMask) {}

// This function returns the topmost widget in the container hierarchy widget
// is a part of. If widget has no parent widgets, it will be returned as the
// topmost widget. No reference will be added to the returned widget; it
// should not be unreferenced. Note the difference in behavior vs.
// GetAncestor; GetAncestor (widget,
// GTK_TYPE_WINDOW) would return NULL if widget wasn't inside a toplevel
// window, and if the window was inside a Window-derived widget which was
// in turn inside the toplevel Window. While the second case may seem
// unlikely, it actually happens when a Plug is embedded inside a
// Socket within the same application. To reliably find the toplevel
// Window, use GetToplevel and check if the TOPLEVEL flags
// is set on the result.
// Returns:
// 	the topmost ancestor of widget , or widget itself if there's no
// 	ancestor.
func (w *CWidget) GetToplevel() (value Widget) {
	return nil
}

// Gets the first ancestor of widget with type widget_type . For example,
// GetAncestor (widget, GTK_TYPE_BOX) gets the first Box
// that's an ancestor of widget . No reference will be added to the returned
// widget; it should not be unreferenced. See note about checking for a
// toplevel Window in the docs for GetToplevel. Note that
// unlike IsAncestor, GetAncestor considers
// widget to be an ancestor of itself.
// Parameters:
// 	widgetType	ancestor type
// Returns:
// 	the ancestor widget, or NULL if not found.
func (w *CWidget) GetAncestor(widgetType cdk.CTypeTag) (value Widget) {
	return nil
}

// Returns the event mask for the widget (a bitfield containing flags from
// the EventMask enumeration). These are the events that the widget will
// receive.
// Returns:
// 	event mask for widget
func (w *CWidget) GetEvents() (value cdk.EventMask) {
	if v, err := w.GetStructProperty(PropertyEvents); err != nil {
		w.LogErr(err)
	} else {
		var ok bool
		if value, ok = v.(cdk.EventMask); !ok {
			w.LogError("value stored in %v property is not of cdk.EventMask type: %v (%T)", PropertyEvents, v, v)
		}
	}
	return
}

// Obtains the location of the mouse pointer in widget coordinates. Widget
// coordinates are a bit odd; for historical reasons, they are defined as
// widget->window coordinates for widgets that are not GTK_NO_WINDOW widgets,
// and are relative to widget->allocation.x , widget->allocation.y for
// widgets that are GTK_NO_WINDOW widgets.
// Parameters:
// 	x	return location for the X coordinate, or NULL.
// 	y	return location for the Y coordinate, or NULL.
func (w *CWidget) GetPointer(x int, y int) {}

// Determines whether widget is somewhere inside ancestor , possibly with
// intermediate containers.
// Returns:
// 	TRUE if ancestor contains widget as a child, grandchild, great
// 	grandchild, etc.
func (w *CWidget) IsAncestor(ancestor Widget) (value bool) {
	return false
}

// Translate coordinates relative to src_widget 's allocation to coordinates
// relative to dest_widget 's allocations. In order to perform this
// operation, both widgets must be realized, and must share a common
// toplevel.
// Parameters:
// 	srcX	X position relative to src_widget
//
// 	srcY	Y position relative to src_widget
//
// 	destX	location to store X position relative to dest_widget
// .
// 	destY	location to store Y position relative to dest_widget
// .
// Returns:
// 	FALSE if either widget was not realized, or there was no common
// 	ancestor. In this case, nothing is stored in *dest_x and
// 	*dest_y . Otherwise TRUE.
func (w *CWidget) TranslateCoordinates(destWidget Widget, srcX int, srcY int, destX int, destY int) (value bool) {
	return false
}

// Utility function; intended to be connected to the delete-event
// signal on a Window. The function calls Hide on its
// argument, then returns TRUE. If connected to ::delete-event, the result is
// that clicking the close button for a window (on the window frame, top
// right corner usually) will hide but not destroy the window. By default,
// CTK destroys windows when ::delete-event is received.
// Returns:
// 	TRUE
func (w *CWidget) HideOnDelete() (value bool) {
	return false
}

// Sets the reading direction on a particular widget. This direction controls
// the primary direction for widgets containing text, and also the direction
// in which the children of a container are packed. The ability to set the
// direction is present in order so that correct localization into languages
// with right-to-left reading directions can be done. Generally, applications
// will let the default reading direction present, except for containers
// where the containers are arranged in an order that is explicitely visual
// rather than logical (such as buttons for text justification). If the
// direction is set to GTK_TEXT_DIR_NONE, then the value set by
// SetDefaultDirection will be used.
// Parameters:
// 	dir	the new direction
func (w *CWidget) SetDirection(dir enums.TextDirection) {}

// Gets the reading direction for a particular widget. See
// SetDirection.
// Returns:
// 	the reading direction for the widget.
func (w *CWidget) GetDirection() (value enums.TextDirection) {
	return enums.TextDirLtr
}

// Sets the default reading direction for widgets where the direction has not
// been explicitly set by SetDirection.
// Parameters:
// 	dir	the new default direction. This cannot be
// GTK_TEXT_DIR_NONE.
func (w *CWidget) SetDefaultDirection(dir enums.TextDirection) {}

// Obtains the current default reading direction. See
// SetDefaultDirection.
// Returns:
// 	the current default direction.
func (w *CWidget) GetDefaultDirection() (value enums.TextDirection) {
	return enums.TextDirLtr
}

// Obtains the full path to widget . The path is simply the name of a widget
// and all its parents in the container hierarchy, separated by periods. The
// name of a widget comes from GetName. Paths are used to apply
// styles to a widget in gtkrc configuration files. Widget names are the type
// of the widget by default (e.g. "Button") or can be set to an
// application-specific value with SetName. By setting the name
// of a widget, you allow users or theme authors to apply styles to that
// specific widget in their gtkrc file. path_reversed_p fills in the path in
// reverse order, i.e. starting with widget 's name instead of starting with
// the name of widget 's outermost ancestor.
// Parameters:
// 	pathLength	location to store length of the path, or NULL.
// 	path	location to store allocated path string, or NULL.
// 	pathReversed	location to store allocated reverse path string, or NULL.
func (w *CWidget) Path() (path string) {
	var parents []Widget
	parent := w.GetParent()
	for {
		if parent == nil {
			break
		}
		parents = append(parents, parent)
		grandparent := parent.GetParent()
		if grandparent != nil && grandparent.ObjectID() == parent.ObjectID() {
			break
		}
	}
	for i := 0; i < len(parents); i++ {
		parent := parents[i]
		parentName := parent.GetTypeTag().ClassName()
		if name := parent.GetName(); name != "" {
			parentName += "#" + name
		}
		if len(path) > 0 {
			path += ">"
		}
		path += parentName
	}
	return
}

// Same as Path, but always uses the name of a widget's type,
// never uses a custom name set with SetName.
// Parameters:
// 	pathLength	location to store the length of the class path, or NULL.
// 	path	location to store the class path as an allocated string, or NULL.
// 	pathReversed	location to store the reverse class path as an allocated
// string, or NULL.
func (w *CWidget) ClassPath(pathLength int, path string, pathReversed string) {}

// Obtains the composite name of a widget.
func (w *CWidget) GetCompositeName() (value string) {
	return ""
}

func (w *CWidget) PushCompositeChild(child Widget) {
	if f := w.Emit(SignalPushCompositeChild, w, child); f == cenums.EVENT_PASS {
		// log.DebugDF(1, "push composite child: %v", child.ObjectName())
		child.Map()
		child.SetParent(w)
		if window := w.GetWindow(); window != nil {
			if wc, ok := child.Self().(Container); ok {
				wc.SetWindow(window)
			} else {
				child.SetWindow(window)
			}
		}
		w.Lock()
		w.composites = append(w.composites, child)
		w.Unlock()
	}
}

func (w *CWidget) PopCompositeChild(child Widget) {
	if f := w.Emit(SignalPopCompositeChild, w, child); f == cenums.EVENT_PASS {
		// log.DebugDF(1, "pop composite child: %v", child.ObjectName())
		w.Lock()
		id := -1
		for idx, composite := range w.composites {
			if composite.ObjectID() == child.ObjectID() {
				id = idx
				break
			}
		}
		if id > -1 {
			count := len(w.composites)
			if id >= count-1 {
				w.composites = w.composites[:id]
			} else {
				w.composites = append(w.composites[:id], w.composites[id+1:]...)
			}
			child.SetWindow(nil)
		}
		w.Unlock()
	}
}

func (w *CWidget) GetCompositeChildren() (composites []Widget) {
	w.RLock()
	defer w.RUnlock()
	composites = append([]Widget{}, w.composites...)
	return
}

// Sets whether the application intends to draw on the widget in an
// expose-event handler. This is a hint to the widget and does not
// affect the behavior of the CTK core; many widgets ignore this flag
// entirely. For widgets that do pay attention to the flag, such as
// EventBox and Window, the effect is to suppress default themed
// drawing of the widget's background. (Children of the widget will still be
// drawn.) The application is then entirely responsible for drawing the
// widget background. Note that the background is still drawn when the widget
// is mapped. If this is not suitable (e.g. because you want to make a
// transparent window using an RGBA visual), you can work around this by
// doing:
// Parameters:
// 	appPaintable	TRUE if the application will paint on the widget
func (w *CWidget) SetAppPaintable(appPaintable bool) {
	w.SetFlags(enums.APP_PAINTABLE)
	if err := w.SetBoolProperty(PropertyAppPaintable, appPaintable); err != nil {
		w.LogErr(err)
	}
}

// Widgets are double buffered by default; you can use this function to turn
// off the buffering. "Double buffered" simply means that
// WindowBeginPaintRegion and WindowEndPaint are called
// automatically around expose events sent to the widget.
// WindowBeginPaint diverts all drawing to a widget's window to an
// offscreen buffer, and WindowEndPaint draws the buffer to the
// screen. The result is that users see the window update in one smooth step,
// and don't see individual graphics primitives being rendered. In very
// simple terms, double buffered widgets don't flicker, so you would only use
// this function to turn off double buffering if you had special needs and
// really knew what you were doing. Note: if you turn off double-buffering,
// you have to handle expose events, since even the clearing to the
// background color or pixmap will not happen automatically (as it is done in
// WindowBeginPaint).
// Parameters:
// 	doubleBuffered	TRUE to double-buffer a widget
func (w *CWidget) SetDoubleBuffered(doubleBuffered bool) {
	if err := w.SetBoolProperty(PropertyDoubleBuffered, doubleBuffered); err != nil {
		w.LogErr(err)
	}
}

// Sets whether the entire widget is queued for drawing when its size
// allocation changes. By default, this setting is TRUE and the entire widget
// is redrawn on every size change. If your widget leaves the upper left
// unchanged when made bigger, turning this setting off will improve
// performance. Note that for NO_WINDOW widgets setting this flag to FALSE
// turns off all allocation on resizing: the widget will not even redraw if
// its position changes; this is to allow containers that don't draw anything
// to avoid excess invalidations. If you set this flag on a NO_WINDOW widget
// that does draw on widget->window , you are responsible for invalidating
// both the old and new allocation of the widget when the widget is moved and
// responsible for invalidating regions newly when the widget increases size.
// Parameters:
// 	redrawOnAllocate	if TRUE, the entire widget will be redrawn
// when it is allocated to a new size. Otherwise, only the
// new portion of the widget will be redrawn.
func (w *CWidget) SetRedrawOnAllocate(redrawOnAllocate bool) {}

// Sets a widgets composite name. The widget must be a composite child of its
// parent; see PushCompositeChild.
// Parameters:
// 	name	the name to set
func (w *CWidget) SetCompositeName(name string) {}

// For widgets that support scrolling, sets the scroll adjustments and
// returns TRUE. For widgets that don't support scrolling, does nothing and
// returns FALSE. Widgets that don't support scrolling can be scrolled by
// placing them in a Viewport, which does support scrolling.
// Parameters:
// 	hadjustment	an adjustment for horizontal scrolling, or NULL.
// 	vadjustment	an adjustment for vertical scrolling, or NULL.
// Returns:
// 	TRUE if the widget supports scrolling
func (w *CWidget) SetScrollAdjustments(hadjustment Adjustment, vadjustment Adjustment) (value bool) {
	return false
}

// Emits a draw signal, primarily used to render canvases and cause end-user
// facing display updates. Signal listeners can draw to the Canvas and return
// EVENT_STOP to cause those changes to be composited upon the larger display
// canvas
//
// Emits: SignalDraw, Argv=[Object instance, canvas]
func (w *CWidget) Draw() cenums.EventFlag {
	if !w.IsFrozen() && w.IsDrawable() && w.GetInvalidated() {
		w.SetInvalidated(false)
		w.LockDraw()
		defer w.UnlockDraw()
		oid := w.ObjectID()
		theme := w.GetThemeRequest()
		if err := memphis.MakeConfigureSurface(oid, w.GetOrigin(), w.GetAllocation(), theme.Content.Normal); err != nil {
			w.LogErr(err)
		} else {
			if surface, err := memphis.GetSurface(oid); err != nil {
				w.LogErr(err)
			} else {
				return w.Emit(SignalDraw, w, surface)
			}
		}
	}
	return cenums.EVENT_PASS
}

func (w *CWidget) Resize() cenums.EventFlag {
	w.LockDraw()
	defer w.UnlockDraw()
	return w.CObject.Resize()
}

// Emits the mnemonic-activate signal. The default handler for this
// signal activates the widget if group_cycling is FALSE, and just grabs the
// focus if group_cycling is TRUE.
// Parameters:
// 	groupCycling	TRUE if there are other widgets with the same mnemonic
// Returns:
// 	TRUE if the signal has been handled
func (w *CWidget) MnemonicActivate(groupCycling bool) (value bool) {
	return false
}

// Very rarely-used function. This function is used to emit an expose event
// signals on a widget. This function is not normally used directly. The only
// time it is used is when propagating an expose event to a child NO_WINDOW
// widget, and that is normally done using ContainerPropagateExpose.
// If you want to force an area of a window to be redrawn, use
// WindowInvalidateRect or WindowInvalidateRegion. To cause
// the redraw to be done immediately, follow that call with a call to
// WindowProcessUpdates.
// Parameters:
// 	event	a expose Event
// Returns:
// 	return from the event signal emission (TRUE if the event was
// 	handled)
func (w *CWidget) SendExpose(event cdk.Event) (value int) {
	return 0
}

// Sends the focus change event to widget This function is not meant to be
// used by applications. The only time it should be used is when it is
// necessary for a Widget to assign focus to a widget that is semantically
// owned by the first widget even though it's not a direct child - for
// instance, a search entry in a floating window similar to the quick search
// in TreeView. An example of its usage is:
// Parameters:
// 	event	a Event of type GDK_FOCUS_CHANGE
// Returns:
// 	the return value from the event signal emission: TRUE if the
// 	event was handled, and FALSE otherwise
func (w *CWidget) SendFocusChange(event cdk.Event) (value bool) {
	return false
}

// This function is used by custom widget implementations; if you're writing
// an app, you'd use GrabFocus to move the focus to a
// particular widget, and ContainerSetFocusChain to change the focus
// tab order. So you may want to investigate those functions instead.
// ChildFocus is called by containers as the user moves around
// the window using keyboard shortcuts. direction indicates what kind of
// motion is taking place (up, down, left, right, tab forward, tab backward).
// ChildFocus emits the focus signal; widgets override
// the default handler for this signal in order to implement appropriate
// focus behavior. The default ::focus handler for a widget should return
// TRUE if moving in direction left the focus on a focusable location inside
// that widget, and FALSE if moving in direction moved the focus outside the
// widget. If returning TRUE, widgets normally call GrabFocus
// to place the focus accordingly; if returning FALSE, they don't modify the
// current focus location. This function replaces ContainerFocus from
// CTK 1.2. It was necessary to check that the child was visible, sensitive,
// and focusable before calling ContainerFocus.
// ChildFocus returns FALSE if the widget is not currently in a
// focusable state, so there's no need for those checks.
// Parameters:
// 	direction	direction of focus movement
// Returns:
// 	TRUE if focus ended up inside widget
func (w *CWidget) ChildFocus(direction enums.DirectionType) (value bool) {
	return false
}

// Emits a child-notify signal for the on widget . This is the analogue
// of g_object_notify for child properties.
// Parameters:
// 	childProperty	the name of a child property installed on the
// class of widget
// 's parent
func (w *CWidget) ChildNotify(childProperty string) {}

// Stops emission of child-notify signals on widget . The signals are
// queued until ThawChildNotify is called on widget . This is
// the analogue of g_object_freeze_notify for child properties.
func (w *CWidget) FreezeChildNotify() {}

// Gets the value set with SetChildVisible. If you feel a need
// to use this function, your code probably needs reorganization. This
// function is only useful for container implementations and never should be
// called by an application.
// Returns:
// 	TRUE if the widget is mapped with the parent.
func (w *CWidget) GetChildVisible() (value bool) {
	return false
}

// Returns the parent container of widget .
// Returns:
// 	the parent container of widget , or NULL.
func (w *CWidget) GetParent() (value Widget) {
	if v, err := w.GetStructProperty(PropertyParent); err != nil {
		w.LogErr(err)
	} else {
		var ok bool
		if value, ok = v.(Widget); !ok && v != nil {
			w.LogError("value stored as %v property is not of Widget type: %v (%T)", PropertyParent, v, v)
			value = nil
		}
	}
	return
}

func (w *CWidget) GetAllParents() (parents []Widget) {
	parent := w.GetParent()
	for parent != nil {
		parents = append(parents, parent)
		next := parent.GetParent()
		if next != nil && next.ObjectID() != parent.ObjectID() {
			parent = next
		} else {
			break
		}
	}
	return
}

// Get the Display for the toplevel window associated with this widget.
// This function can only be called after the widget has been added to a
// widget hierarchy with a Window at the top. In general, you should only
// create display specific resources when a widget has been realized, and you
// should free those resources when the widget is unrealized.
// Returns:
// 	the Display for the toplevel for this widget.
func (w *CWidget) GetDisplay() (value cdk.Display) {
	if window := w.GetWindow(); window != nil {
		value = window.GetDisplay()
	}
	return
}

// Returns the clipboard object for the given selection to be used with
// widget . widget must have a Display associated with it, so must be
// attached to a toplevel window.
// Parameters:
// 	selection	a Atom which identifies the clipboard
// to use. GDK_SELECTION_CLIPBOARD gives the
// default clipboard. Another common value
// is GDK_SELECTION_PRIMARY, which gives
// the primary X selection.
// Returns:
// 	the appropriate clipboard object. If no clipboard already
// 	exists, a new one will be created. Once a clipboard object has
// 	been created, it is persistent for all time.
// func (w *CWidget) GetClipboard(selection Atom) (value Clipboard) {
// 	return nil
// }

// Get the root window where this widget is located. This function can only
// be called after the widget has been added to a widget hierarchy with
// Window at the top. The root window is useful for such purposes as
// creating a popup Window associated with the window. In general, you
// should only create display specific resources when a widget has been
// realized, and you should free those resources when the widget is
// unrealized.
// Returns:
// 	the Window root window for the toplevel for this widget.
func (w *CWidget) GetRootWindow() (value Window) {
	return nil
}

// Get the Screen from the toplevel window associated with this widget.
// This function can only be called after the widget has been added to a
// widget hierarchy with a Window at the top. In general, you should only
// create screen specific resources when a widget has been realized, and you
// should free those resources when the widget is unrealized.
// Returns:
// 	the Screen for the toplevel for this widget.
func (w *CWidget) GetScreen() (value cdk.Display) {
	return nil
}

// Checks whether there is a Screen is associated with this widget. All
// toplevel widgets have an associated screen, and all widgets added into a
// hierarchy with a toplevel window at the top.
// Returns:
// 	TRUE if there is a Screen associcated with the widget.
func (w *CWidget) HasScreen() (value bool) {
	return false
}

// Gets the size request that was explicitly set for the widget using
// SetSizeRequest. A value of -1 stored in width or height
// indicates that that dimension has not been set explicitly and the natural
// requisition of the widget will be used intead. See
// SetSizeRequest. To get the size a widget will actually use,
// call SizeRequest instead of this function.
// Parameters:
// 	width	return location for width, or NULL.
// 	height	return location for height, or NULL.
func (w *CWidget) GetSizeRequest() (width, height int) {
	// w.RLock()
	// defer w.RUnlock()
	var err error
	if width, err = w.GetIntProperty(PropertyWidthRequest); err != nil {
		w.LogErr(err)
	}
	if height, err = w.GetIntProperty(PropertyHeightRequest); err != nil {
		w.LogErr(err)
	}
	return
}

// Returns the currently requested size
func (w *CWidget) SizeRequest() ptypes.Rectangle {
	return ptypes.MakeRectangle(w.GetSizeRequest())
}

// Sets the minimum size of a widget; that is, the widget's size request will
// be width by height . You can use this function to force a widget to be
// either larger or smaller than it normally would be. In most cases,
// WindowSetDefaultSize is a better choice for toplevel windows than
// this function; setting the default size will still allow users to shrink
// the window. Setting the size request will force them to leave the window
// at least as large as the size request. When dealing with window sizes,
// WindowSetGeometryHints can be a useful function as well. Note the
// inherent danger of setting any fixed size - themes, translations into
// other languages, different fonts, and user action can all change the
// appropriate size for a given widget. So, it's basically impossible to
// hardcode a size that will always be correct. The size request of a widget
// is the smallest size a widget can accept while still functioning well and
// drawing itself correctly. However in some strange cases a widget may be
// allocated less than its requested size, and in many cases a widget may be
// allocated more space than it requested. If the size request in a given
// direction is -1 (unset), then the "natural" size request of the widget
// will be used instead. Widgets can't actually be allocated a size less than
// 1 by 1, but you can pass 0,0 to this function to mean "as small as
// possible."
// Parameters:
// 	width	width widget should request, or -1 to unset
// 	height	height widget should request, or -1 to unset
//
// Emits: SignalSetSizeRequest, Argv=[Widget instance, given size]
func (w *CWidget) SetSizeRequest(width, height int) {
	if f := w.Emit(SignalSetSizeRequest, w, ptypes.MakeRectangle(width, height)); f == cenums.EVENT_PASS {
		if err := w.SetIntProperty(PropertyWidthRequest, width); err != nil {
			w.LogErr(err)
		}
		if err := w.SetIntProperty(PropertyHeightRequest, height); err != nil {
			w.LogErr(err)
		}
	}
}

// Sets the no-show-all property, which determines whether calls to
// ShowAll and HideAll will affect this widget.
// This is mostly for use in constructing widget hierarchies with externally
// controlled visibility, see UIManager.
// Parameters:
// 	noShowAll	the new value for the "no-show-all" property
func (w *CWidget) SetNoShowAll(noShowAll bool) {
	if err := w.SetBoolProperty(PropertyNoShowAll, noShowAll); err != nil {
		w.LogErr(err)
	}
}

// Returns the current value of the Widget:no-show-all property, which
// determines whether calls to ShowAll and
// HideAll will affect this widget.
// Returns:
// 	the current value of the "no-show-all" property.
func (w *CWidget) GetNoShowAll() (value bool) {
	var err error
	if value, err = w.GetBoolProperty(PropertyNoShowAll); err != nil {
		w.LogErr(err)
	}
	return
}

// Adds a widget to the list of mnemonic labels for this widget. (See
// ListMnemonicLabels). Note the list of mnemonic labels for
// the widget is cleared when the widget is destroyed, so the caller must
// make sure to update its internal state at this point as well, by using a
// connection to the destroy signal or a weak notifier.
func (w *CWidget) AddMnemonicLabel(label Widget) {}

// Removes a widget from the list of mnemonic labels for this widget. (See
// ListMnemonicLabels). The widget must have previously been
// added to the list with AddMnemonicLabel.
func (w *CWidget) RemoveMnemonicLabel(label Widget) {}

// Notifies the user about an input-related error on this widget. If the
// gtk-error-bell setting is TRUE, it calls WindowBeep,
// otherwise it does nothing. Note that the effect of WindowBeep can
// be configured in many ways, depending on the windowing backend and the
// desktop environment or window manager that is used.
func (w *CWidget) ErrorBell() {}

// This function should be called whenever keyboard navigation within a
// single widget hits a boundary. The function emits the keynav-failed
// signal on the widget and its return value should be interpreted in a way
// similar to the return value of ChildFocus: When TRUE is
// returned, stay in the widget, the failed keyboard navigation is Ok and/or
// there is nowhere we can/should move the focus to. When FALSE is returned,
// the caller should continue with keyboard navigation outside the widget,
// e.g. by calling ChildFocus on the widget's toplevel. The
// default ::keynav-failed handler returns TRUE for GTK_DIR_TAB_FORWARD and
// GTK_DIR_TAB_BACKWARD. For the other values of DirectionType, it looks
// at the gtk-keynav-cursor-only setting and returns FALSE if the
// setting is TRUE. This way the entire user interface becomes
// cursor-navigatable on input devices such as mobile phones which only have
// cursor keys but no tab key. Whenever the default handler returns TRUE, it
// also calls ErrorBell to notify the user of the failed
// keyboard navigation. A use case for providing an own implementation of
// ::keynav-failed (either by connecting to it or by overriding it) would be
// a row of Entry widgets where the user should be able to navigate the
// entire row with the cursor keys, as e.g. known from user interfaces that
// require entering license keys.
// Parameters:
// 	direction	direction of focus movement
// Returns:
// 	TRUE if stopping keyboard navigation is fine, FALSE if the
// 	emitting widget should try to handle the keyboard navigation
// 	attempt in its parent container(s).
func (w *CWidget) KeynavFailed(direction enums.DirectionType) (value bool) {
	return false
}

// Gets the contents of the tooltip for widget .
// Returns:
// 	the tooltip text, or NULL. You should free the returned string
// 	with g_free when done.
func (w *CWidget) GetTooltipMarkup() (value string) {
	var err error
	if value, err = w.GetStringProperty(PropertyTooltipMarkup); err != nil {
		w.LogErr(err)
	}
	return
}

// Sets markup as the contents of the tooltip, which is marked up with the
// Tango text markup language. This function will take care of setting
// Widget:has-tooltip to TRUE and of the default handler for the
// Widget::query-tooltip signal. See also the Widget:tooltip-markup
// property and TooltipSetMarkup.
// Parameters:
// 	markup	the contents of the tooltip for widget
// , or NULL.
func (w *CWidget) SetTooltipMarkup(markup string) {
	if err := w.SetStringProperty(PropertyTooltipMarkup, markup); err != nil {
		w.LogErr(err)
	}
}

// Gets the contents of the tooltip for widget .
// Returns:
// 	the tooltip text, or NULL. You should free the returned string
// 	with g_free when done.
func (w *CWidget) GetTooltipText() (value string) {
	var err error
	if value, err = w.GetStringProperty(PropertyTooltipText); err != nil {
		w.LogErr(err)
	}
	return
}

// Sets text as the contents of the tooltip. This function will take care of
// setting Widget:has-tooltip to TRUE and of the default handler for the
// Widget::query-tooltip signal. See also the Widget:tooltip-text
// property and TooltipSetText.
// Parameters:
// 	text	the contents of the tooltip for widget
//
func (w *CWidget) SetTooltipText(text string) {
	if err := w.SetStringProperty(PropertyTooltipText, text); err != nil {
		w.LogErr(err)
	}
	w.SetHasTooltip(len(text) > 0)
}

// Returns the Window of the current tooltip. This can be the Window
// created by default, or the custom tooltip window set using
// SetTooltipWindow.
// Returns:
// 	The Window of the current tooltip.
func (w *CWidget) GetTooltipWindow() (value Window) {
	w.RLock()
	defer w.RUnlock()
	if w.tooltipWindow == nil {
		w.tooltipWindow = w.newTooltipWindow()
	}
	return w.tooltipWindow
}

// Replaces the default, usually yellow, window used for displaying tooltips
// with custom_window . CTK will take care of showing and hiding
// custom_window at the right moment, to behave likewise as the default
// tooltip window. If custom_window is NULL, the default tooltip window will
// be used. If the custom window should have the default theming it needs to
// have the name "gtk-tooltip", see SetName.
// Parameters:
// 	customWindow	a Window, or NULL.
func (w *CWidget) SetTooltipWindow(customWindow Window) {
	if customWindow != nil {
		w.RLock()
		tooltipWindow := w.tooltipWindow
		w.RUnlock()
		if tooltipWindow != nil {
			w.Lock()
			tooltipWindow.Destroy()
			w.tooltipWindow = nil
			w.Unlock()
		}
		w.Lock()
		w.tooltipWindow = customWindow
		w.Unlock()
	} else {
		defaultTooltipWindow := w.newTooltipWindow()
		w.Lock()
		w.tooltipWindow = defaultTooltipWindow
		w.Unlock()
	}
}

// Returns the current value of the has-tooltip property. See
// Widget:has-tooltip for more information.
// Returns:
// 	current value of has-tooltip on widget .
func (w *CWidget) GetHasTooltip() (value bool) {
	var err error
	if value, err = w.GetBoolProperty(PropertyHasTooltip); err != nil {
		w.LogErr(err)
	}
	return
}

// Sets the has-tooltip property on widget to has_tooltip . See
// Widget:has-tooltip for more information.
// Parameters:
// 	hasTooltip	whether or not widget
// has a tooltip.
func (w *CWidget) SetHasTooltip(hasTooltip bool) {
	if err := w.SetBoolProperty(PropertyHasTooltip, hasTooltip); err != nil {
		w.LogErr(err)
	}
}

// Returns the widget's window if it is realized, NULL otherwise
// Returns:
// 	widget 's window.
// Returns the Window instance associated with this Widget instance, nil
// otherwise
func (w *CWidget) GetWindow() (window Window) {
	if v, err := w.GetStructProperty(PropertyWindow); err != nil {
		w.LogErr(err)
		return nil
	} else {
		var ok bool
		if window, ok = v.(Window); !ok && v != nil {
			w.LogError("value stored in %v property is not of Window type: %v (%T)", PropertyWindow, v, v)
			return nil
		}
	}
	if window == nil {
		window = w.GetParentWindow()
	}
	return
}

// Determines whether the application intends to draw on the widget in an
// expose-event handler. See SetAppPaintable
// Returns:
// 	TRUE if the widget is app paintable
func (w *CWidget) GetAppPaintable() (value bool) {
	var err error
	if value, err = w.GetBoolProperty(PropertyAppPaintable); err != nil {
		w.LogErr(err)
	}
	return
}

// Determines whether widget can be a default widget. See
// SetCanDefault.
// Returns:
// 	TRUE if widget can be a default widget, FALSE otherwise
func (w *CWidget) GetCanDefault() (value bool) {
	var err error
	if value, err = w.GetBoolProperty(PropertyCanDefault); err != nil {
		w.LogErr(err)
	}
	return
}

// Specifies whether widget can be a default widget. See
// GrabDefault for details about the meaning of "default".
// Parameters:
// 	canDefault	whether or not widget
// can be a default widget.
func (w *CWidget) SetCanDefault(canDefault bool) {
	if err := w.SetBoolProperty(PropertyCanDefault, canDefault); err != nil {
		w.LogErr(err)
	}
}

// Determines whether widget can own the input focus. See
// SetCanFocus.
// Returns:
// 	TRUE if widget can own the input focus, FALSE otherwise
func (w *CWidget) GetCanFocus() (value bool) {
	var err error
	if value, err = w.GetBoolProperty(PropertyCanFocus); err != nil {
		w.LogErr(err)
	}
	return
}

// Specifies whether widget can own the input focus. See
// GrabFocus for actually setting the input focus on a widget.
// Parameters:
// 	canFocus	whether or not widget
// can own the input focus.
func (w *CWidget) SetCanFocus(canFocus bool) {
	if err := w.SetBoolProperty(PropertyCanFocus, canFocus); err != nil {
		w.LogErr(err)
	}
}

// Determines whether widget has a Window of its own. See
// SetHasWindow.
// Returns:
// 	TRUE if widget has a window, FALSE otherwise
func (w *CWidget) GetHasWindow() (ok bool) {
	if value, err := w.GetStructProperty(PropertyWindow); err == nil && value != nil {
		_, ok = value.(Window)
	} else if err != nil {
		w.LogErr(err)
	}
	return
}

// Returns the widget's sensitivity (in the sense of returning the value that
// has been set using SetSensitive). The effective sensitivity
// of a widget is however determined by both its own and its parent widget's
// sensitivity. See IsSensitive.
// Returns:
// 	TRUE if the widget is sensitive
func (w *CWidget) GetSensitive() (value bool) {
	var err error
	if value, err = w.GetBoolProperty(PropertySensitive); err != nil {
		w.LogErr(err)
	}
	return
}

// Returns the widget's effective sensitivity, which means it is sensitive
// itself and also its parent widget is sensitive
// Returns:
// 	TRUE if the widget is effectively sensitive
func (w *CWidget) IsSensitive() bool {
	if w.HasState(enums.StateInsensitive) {
		return false
	}
	if parent := w.GetParent(); parent != nil {
		if parent.HasState(enums.StateInsensitive) {
			return false
		}
	}
	return true
}

// Determines whether the widget is visible. Note that this doesn't take into
// account whether the widget's parent is also visible or the widget is
// obscured in any way. See SetVisible.
// Returns:
// 	TRUE if the widget is visible
func (w *CWidget) GetVisible() (value bool) {
	var err error
	if value, err = w.GetBoolProperty(PropertyVisible); err != nil {
		w.LogErr(err)
	}
	return
}

// Sets the visibility state of widget . Note that setting this to TRUE
// doesn't mean the widget is actually viewable, see
// GetVisible. This function simply calls Show or
// Hide but is nicer to use when the visibility of the widget
// depends on some condition.
// Parameters:
// 	visible	whether the widget should be shown or not
func (w *CWidget) SetVisible(visible bool) {
	if err := w.SetBoolProperty(PropertyVisible, visible); err != nil {
		w.LogErr(err)
	}
}

// Determines whether widget is the current default widget within its
// toplevel. See SetCanDefault.
// Returns:
// 	TRUE if widget is the current default widget within its
// 	toplevel, FALSE otherwise
func (w *CWidget) HasDefault() (value bool) {
	return false
}

// Determines if the widget has the global input focus. See
// IsFocus for the difference between having the global input
// focus, and only having the focus within a toplevel.
// Returns:
// 	TRUE if the widget has the global input focus.
func (w *CWidget) HasFocus() (value bool) {
	var err error
	if value, err = w.GetBoolProperty(PropertyHasFocus); err != nil {
		w.LogErr(err)
	}
	return
}

// Determines whether the widget is currently grabbing events, so it is the
// only widget receiving input events (keyboard and mouse). See also
// GrabAdd.
// Returns:
// 	TRUE if the widget is in the grab_widgets stack
func (w *CWidget) HasGrab() (value bool) {
	return false
}

// Determines whether widget can be drawn to. A widget can be drawn to if it
// is mapped and visible.
// Returns:
// 	TRUE if widget is drawable, FALSE otherwise
// Returns TRUE if the APP_PAINTABLE flag is set, FALSE otherwise
func (w *CWidget) IsDrawable() (value bool) {
	return w.HasFlags(enums.APP_PAINTABLE)
}

// Determines whether widget is a toplevel widget. Currently only Window
// and Invisible are toplevel widgets. Toplevel widgets have no parent
// widget.
// Returns:
// 	TRUE if widget is a toplevel, FALSE otherwise
func (w *CWidget) IsToplevel() (value bool) {
	return w.HasFlags(enums.TOPLEVEL)
}

// Sets a widget's window. This function should only be used in a widget's
// Widget::realize implementation. The window passed is usually either
// new window created with WindowNew, or the window of its parent
// widget as returned by GetParentWindow. Widgets must
// indicate whether they will create their own Window by calling
// SetHasWindow. This is usually done in the widget's init
// function.
// Parameters:
// 	window	a Window
//
// Emits: SignalSetWindow, Argv=[Widget instance, given window]
//
// Locking: write [indirect]
func (w *CWidget) SetWindow(window Window) {
	if f := w.Emit(SignalSetWindow, w, window); f == cenums.EVENT_PASS {
		if err := w.SetStructProperty(PropertyWindow, window); err != nil {
			w.LogErr(err)
		} else if window != nil {
			for _, composite := range w.GetCompositeChildren() {
				if err := composite.SetStructProperty(PropertyWindow, window); err != nil {
					w.LogErr(err)
				}
			}
			window.ApplyStylesTo(w)
		}
	}
}

// Specifies whether widget will be treated as the default widget within its
// toplevel when it has the focus, even if another widget is the default. See
// GrabDefault for details about the meaning of "default".
// Parameters:
// 	receivesDefault	whether or not widget
// can be a default widget.
func (w *CWidget) SetReceivesDefault(receivesDefault bool) {
	if err := w.SetBoolProperty(PropertyReceivesDefault, receivesDefault); err != nil {
		w.LogErr(err)
	}
}

// Determines whether widget is alyways treated as default widget withing its
// toplevel when it has the focus, even if another widget is the default. See
// SetReceivesDefault.
// Returns:
// 	TRUE if widget acts as default widget when focussed, FALSE
// 	otherwise
func (w *CWidget) GetReceivesDefault() (value bool) {
	var err error
	if value, err = w.GetBoolProperty(PropertyReceivesDefault); err != nil {
		w.LogErr(err)
	}
	return
}

// Marks the widget as being realized. This function should only ever be
// called in a derived widget's "realize" or "unrealize" implementation.
// Parameters:
// 	realized	TRUE to mark the widget as realized
func (w *CWidget) SetRealized(realized bool) {}

// Determines whether widget is realized.
// Returns:
// 	TRUE if widget is realized, FALSE otherwise
func (w *CWidget) GetRealized() (value bool) {
	return false
}

// Marks the widget as being realized. This function should only ever be
// called in a derived widget's "map" or "unmap" implementation.
// Parameters:
// 	mapped	TRUE to mark the widget as mapped
func (w *CWidget) SetMapped(mapped bool) {}

// Whether the widget is mapped.
// Returns:
// 	TRUE if the widget is mapped, FALSE otherwise.
func (w *CWidget) GetMapped() (value bool) {
	return w.IsMapped()
}

func (w *CWidget) GetTheme() (theme paint.Theme) {
	theme = w.CObject.GetTheme()
	applyStyles := func(state enums.StateType, content, border paint.Style) (paint.Style, paint.Style) {
		var wColor, wBackgroundColor paint.Color
		var wBorderColor, wBorderBackgroundColor paint.Color
		var wBold, wBlink, wReverse, wUnderline, wDim, wItalic, wStrike bool
		wColor, _ = w.GetCssColor(CssPropertyColor, state)
		wBackgroundColor, _ = w.GetCssColor(CssPropertyBackgroundColor, state)
		wBorderColor, _ = w.GetCssColor(CssPropertyBorderColor, state)
		wBorderBackgroundColor, _ = w.GetCssColor(CssPropertyBorderBackgroundColor, state)
		wBold, _ = w.GetCssBool(CssPropertyBold, state)
		wBlink, _ = w.GetCssBool(CssPropertyBlink, state)
		wReverse, _ = w.GetCssBool(CssPropertyReverse, state)
		wUnderline, _ = w.GetCssBool(CssPropertyUnderline, state)
		wDim, _ = w.GetCssBool(CssPropertyDim, state)
		wItalic, _ = w.GetCssBool(CssPropertyItalic, state)
		wStrike, _ = w.GetCssBool(CssPropertyStrike, state)
		modContent := content.Foreground(wColor).Background(wBackgroundColor).Bold(wBold).Blink(wBlink).Reverse(wReverse).Underline(wUnderline).Dim(wDim).Italic(wItalic).Strike(wStrike)
		modBorder := border.Foreground(wBorderColor).Background(wBorderBackgroundColor).Bold(wBold).Blink(wBlink).Reverse(wReverse).Underline(wUnderline).Dim(wDim).Italic(wItalic).Strike(wStrike)
		return modContent, modBorder
	}
	theme.Content.Normal, theme.Border.Normal = applyStyles(enums.StateNormal, theme.Content.Normal, theme.Border.Normal)
	theme.Content.Active, theme.Border.Active = applyStyles(enums.StateActive, theme.Content.Active, theme.Border.Active)
	theme.Content.Selected, theme.Border.Selected = applyStyles(enums.StateSelected, theme.Content.Selected, theme.Border.Selected)
	theme.Content.Prelight, theme.Border.Prelight = applyStyles(enums.StatePrelight, theme.Content.Prelight, theme.Border.Prelight)
	theme.Content.Insensitive, theme.Border.Insensitive = applyStyles(enums.StateInsensitive, theme.Content.Insensitive, theme.Border.Insensitive)
	return
}

// GetThemeRequest returns the current theme, adjusted for Widget state and
// accounting for any PARENT_SENSITIVE conditions. This method is primarily
// useful in drawable Widget types during Invalidate() and Draw() stages of
// the Widget lifecycle. This method emits an initial get-theme-request signal
// with a pointer to the theme instance to be modified as there are no return
// values for signal listeners. If the signal listeners return EVENT_STOP, the
// theme instance is returned without modification. If the signal listeners
// return EVENT_PASS, this method will perform the changes to the theme instance
// mentioned above.
func (w *CWidget) GetThemeRequest() (theme paint.Theme) {
	theme = w.GetTheme()
	instance := &theme
	modified := theme
	if w.IsDrawable() && w.IsVisible() {
		if w.HasState(enums.StateInsensitive) {
			modified.Content.Normal = modified.Content.Insensitive
			modified.Border.Normal = modified.Border.Insensitive
		} else if w.HasState(enums.StateActive) {
			modified.Content.Normal = modified.Content.Active
			modified.Border.Normal = modified.Border.Active
		} else if w.HasState(enums.StateSelected) {
			modified.Content.Normal = modified.Content.Selected
			modified.Border.Normal = modified.Border.Selected
		} else if w.HasState(enums.StatePrelight) {
			modified.Content.Normal = modified.Content.Prelight
			modified.Border.Normal = modified.Border.Prelight
		}
		if f := w.Emit(SignalGetThemeRequest, instance, modified); f == cenums.EVENT_PASS {
			theme = modified
			w.LogTrace("get-theme-request changes applied")
		}
	}
	return
}

// Set the Theme for the Widget instance. This will also refresh the requested
// theme. A request theme is a transient theme, based on the actually set theme
// and adjusted for focus. If the given theme is equivalent to the current theme
// then no action is taken. After verifying that the given theme is different,
// this method emits a set-theme signal and if the listeners return EVENT_PASS,
// the changes are applied and the Widget.Invalidate() method is called
func (w *CWidget) SetTheme(theme paint.Theme) {
	// if theme.String() != w.GetTheme().String() {
	if f := w.Emit(SignalSetTheme, w, theme); f == cenums.EVENT_PASS {

		apply := func(state enums.StateType, cs, bs paint.Style) {
			fg, bg, attr := cs.Decompose()
			if prop := w.GetCssProperty(CssPropertyColor, state); prop != nil {
				if err := prop.Set(fg); err != nil {
					w.LogErr(err)
				}
			}
			if prop := w.GetCssProperty(CssPropertyBackgroundColor, state); prop != nil {
				if err := prop.Set(bg); err != nil {
					w.LogErr(err)
				}
			}
			if prop := w.GetCssProperty(CssPropertyBold, state); prop != nil {
				if err := prop.Set(attr.IsBold()); err != nil {
					w.LogErr(err)
				}
			}
			if prop := w.GetCssProperty(CssPropertyBlink, state); prop != nil {
				if err := prop.Set(attr.IsBlink()); err != nil {
					w.LogErr(err)
				}
			}
			if prop := w.GetCssProperty(CssPropertyReverse, state); prop != nil {
				if err := prop.Set(attr.IsReverse()); err != nil {
					w.LogErr(err)
				}
			}
			if prop := w.GetCssProperty(CssPropertyUnderline, state); prop != nil {
				if err := prop.Set(attr.IsUnderline()); err != nil {
					w.LogErr(err)
				}
			}
			if prop := w.GetCssProperty(CssPropertyDim, state); prop != nil {
				if err := prop.Set(attr.IsDim()); err != nil {
					w.LogErr(err)
				}
			}
			if prop := w.GetCssProperty(CssPropertyItalic, state); prop != nil {
				if err := prop.Set(attr.IsItalic()); err != nil {
					w.LogErr(err)
				}
			}
			if prop := w.GetCssProperty(CssPropertyStrike, state); prop != nil {
				if err := prop.Set(attr.IsStrike()); err != nil {
					w.LogErr(err)
				}
			}
			fg, bg, _ = bs.Decompose()
			if prop := w.GetCssProperty(CssPropertyBorderColor, state); prop != nil {
				if err := prop.Set(fg); err != nil {
					w.LogErr(err)
				}
			}
			if prop := w.GetCssProperty(CssPropertyBorderBackgroundColor, state); prop != nil {
				if err := prop.Set(bg); err != nil {
					w.LogErr(err)
				}
			}
			return
		}

		w.CObject.SetTheme(theme)

		apply(enums.StateNormal, theme.Content.Normal, theme.Border.Normal)
		apply(enums.StateActive, theme.Content.Active, theme.Border.Active)
		apply(enums.StateSelected, theme.Content.Selected, theme.Border.Selected)
		apply(enums.StatePrelight, theme.Content.Prelight, theme.Border.Prelight)
		apply(enums.StateInsensitive, theme.Content.Insensitive, theme.Border.Insensitive)

		w.Invalidate()
	}
	// }
}

// Returns the widget's state. See SetState.
// Returns:
// 	the state of the widget.
func (w *CWidget) GetState() (value enums.StateType) {
	w.flagsLock.RLock()
	defer w.flagsLock.RUnlock()
	return w.state
}

// SetState used in widget implementations. Sets the state of a widget
// (insensitive, prelighted, etc). Usually you should set the state using
// wrapper functions such as SetSensitive. If the state given is StateNone,
// this will reset the state flags to only StateNormal. This method emits a
// set-state signal initially and if the listeners return EVENT_PASS, the change
// is applied.
//
// Parameters:
// 	state	new state for widget
func (w *CWidget) SetState(state enums.StateType) {
	if f := w.Emit(SignalSetState, w, state); f == cenums.EVENT_PASS {
		w.flagsLock.Lock()
		defer w.flagsLock.Unlock()
		if state == enums.StateNone {
			w.state = enums.StateNormal
		} else {
			w.state = w.state.Set(state)
		}
	}
}

// HasState returns TRUE if the Widget has the given StateType, FALSE otherwise.
// Passing StateNone will always return FALSE.
func (w *CWidget) HasState(s enums.StateType) bool {
	if s == enums.StateNone {
		return false
	}
	w.flagsLock.RLock()
	defer w.flagsLock.RUnlock()
	return w.state.Has(s)
}

// UnsetState clears the given state bitmask from the Widget instance. If the
// given state is StateNone, no action is taken. This method emits an
// unset-state signal initially and if the listeners return EVENT_PASS, the
// change is applied.
func (w *CWidget) UnsetState(state enums.StateType) {
	if state == enums.StateNone {
		return
	}
	if f := w.Emit(SignalUnsetState, w, state); f == cenums.EVENT_PASS {
		w.flagsLock.Lock()
		defer w.flagsLock.Unlock()
		w.state = w.state.Clear(state)
	}
}

// Returns the current flags for the Widget instance
func (w *CWidget) GetFlags() enums.WidgetFlags {
	w.flagsLock.RLock()
	defer w.flagsLock.RUnlock()
	return w.flags
}

// Returns TRUE if the Widget instance has the given flag, FALSE otherwise
func (w *CWidget) HasFlags(f enums.WidgetFlags) bool {
	w.flagsLock.RLock()
	defer w.flagsLock.RUnlock()
	return w.flags.Has(f)
}

// Removes the given flags from the Widget instance. This method emits an
// unset-flags signal initially and if the listeners return EVENT_PASS, the
// change is applied
//
// Emits: SignalUnsetFlags, Argv=[Widget instance, given flags to unset]
func (w *CWidget) UnsetFlags(v enums.WidgetFlags) {
	if f := w.Emit(SignalUnsetFlags, w, v); f == cenums.EVENT_PASS {
		w.flagsLock.Lock()
		defer w.flagsLock.Unlock()
		w.flags = w.flags.Clear(v)
	}
}

// Sets the given flags on the Widget instance. This method emits a set-flags
// signal initially and if the listeners return EVENT_PASS, the change is
// applied
//
// Emits: SignalSetFlags, Argv=[Widget instance, given flags to set]
func (w *CWidget) SetFlags(v enums.WidgetFlags) {
	if f := w.Emit(SignalSetFlags, w, w.flags, v); f == cenums.EVENT_PASS {
		w.flagsLock.Lock()
		defer w.flagsLock.Unlock()
		w.flags = w.flags.Set(v)
	}
}

// If the Widget instance is PARENT_SENSITIVE and one of it's parents are the
// focus for the associated Window, return TRUE and FALSE otherwise
func (w *CWidget) IsParentFocused() bool {
	if w.HasFlags(enums.PARENT_SENSITIVE) {
		var lastParent Widget
		parent, _ := w.GetParent().(Widget)
		for parent != nil {
			if parent.IsFocus() {
				return true
			}
			if !parent.HasFlags(enums.PARENT_SENSITIVE) || parent.IsToplevel() {
				// don't recurse
				break
			}
			lastParent = parent
			if parent, _ = parent.GetParent().(Widget); parent != nil {
				if lastParent.ObjectID() == parent.ObjectID() {
					break // stop
				}
			}
		}
	}
	return false
}

// Returns TRUE if the Widget instance or it's parent are the current focus of
// the associated Window
func (w *CWidget) IsFocused() bool {
	return w.IsFocus() || w.IsParentFocused()
}

// Returns TRUE if the Widget instance has the CAN_FOCUS flag, FALSE otherwise
func (w *CWidget) CanFocus() bool {
	return w.HasFlags(enums.CAN_FOCUS)
}

// Returns TRUE if the Widget instance CanDefault() and the HAS_DEFAULT flag is
// set, returns FALSE otherwise
func (w *CWidget) IsDefault() bool {
	if w.CanDefault() {
		return w.HasFlags(enums.HAS_DEFAULT)
	}
	return false
}

// Returns TRUE if the Widget instance IsSensitive() and the CAN_DEFAULT flag is
// set, returns FALSE otherwise
func (w *CWidget) CanDefault() bool {
	if w.IsSensitive() {
		return w.HasFlags(enums.CAN_DEFAULT)
	}
	return false
}

// Returns TRUE if the VISIBLE flag is set, FALSE otherwise
func (w *CWidget) IsVisible() bool {
	return w.HasFlags(enums.VISIBLE)
}

func (w *CWidget) HasEventFocus() bool {
	oid := w.ObjectID()
	window := w.GetWindow()
	w.RLock()
	if window != nil {
		if ef := window.GetEventFocus(); ef != nil {
			if ef.ObjectID() == oid {
				w.RUnlock()
				return true
			}
		}
	}
	w.RUnlock()
	return false
}

// GrabEventFocus will attempt to set the Widget as the Window event focus
// handler. This method emits a grab-event-focus signal and if the listeners all
// return EVENT_PASS, the changes are applied.
//
// Note that this method needs to be implemented within each Drawable that can
// be focused because of the golang interface system losing the concrete struct
// when a Widget interface reference is passed as a generic interface{}
// argument.
func (w *CWidget) GrabEventFocus() {
	if window := w.GetWindow(); window != nil {
		if f := w.Emit(SignalGrabEventFocus, w, window); f == cenums.EVENT_PASS {
			window.SetEventFocus(w)
		}
	} else {
		w.LogError("cannot grab focus: missing window")
	}
}

func (w *CWidget) ReleaseEventFocus() {
	if window := w.GetWindow(); window != nil {
		if ef := window.GetEventFocus(); ef != nil {
			if wef, ok := ef.Self().(Widget); ok && wef.ObjectID() == w.ObjectID() {
				if f := w.Emit(SignalReleaseEventFocus, w, window); f == cenums.EVENT_PASS {
					window.SetEventFocus(nil)
				}
			}
		}
	}
}

// Returns the top-most parent in the Widget instance's parent hierarchy.
// Returns nil if the Widget has no parent container
func (w *CWidget) GetTopParent() (parent Widget) {
	if parent = w.GetParent(); parent != nil {
		var last Widget
		for {
			if parent = parent.GetParent(); parent != nil && parent.ObjectID() != last.ObjectID() {
				last = parent
			} else {
				parent = last
				break
			}
		}
	}
	return
}

// A wrapper around the Object.GetObjectAt() method, only returning Widget
// instance types or nil otherwise
func (w *CWidget) GetWidgetAt(p *ptypes.Point2I) Widget {
	if o := w.CObject.GetObjectAt(p); o != nil {
		if ow, ok := o.Self().(Widget); ok && ow.IsVisible() {
			return ow
		}
	}
	return nil
}

func (w *CWidget) RequestDrawAndShow() {
	if window := w.GetWindow(); window != nil {
		window.RequestDrawAndShow()
	}
}

func (w *CWidget) RequestDrawAndSync() {
	if window := w.GetWindow(); window != nil {
		window.RequestDrawAndSync()
	}
}

func (w *CWidget) Invalidate() cenums.EventFlag {
	if rv := w.CObject.Invalidate(); rv == cenums.EVENT_PASS {
		parent := w.GetParent()
		for parent != nil {
			parent.SetInvalidated(true)
			if next := parent.GetParent(); next != nil && next.ObjectID() != parent.ObjectID() {
				parent = next
			} else {
				parent = nil
			}
		}
		return cenums.EVENT_STOP
	}
	return cenums.EVENT_PASS
}

func (w *CWidget) lostFocus(_ []interface{}, _ ...interface{}) cenums.EventFlag {
	if w.IsDrawable() && w.IsVisible() {
		w.UnsetState(enums.StateSelected)
		w.Invalidate()
		w.LogDebug("lost focus")
	}
	return cenums.EVENT_PASS
}

func (w *CWidget) gainedFocus(_ []interface{}, _ ...interface{}) cenums.EventFlag {
	if w.IsDrawable() && w.IsVisible() && w.IsSensitive() {
		w.SetState(enums.StateSelected)
		w.Invalidate()
		w.LogDebug("gained focus")
	}
	return cenums.EVENT_PASS
}

func (w *CWidget) enter(_ []interface{}, argv ...interface{}) cenums.EventFlag {
	if w.IsDrawable() && w.IsVisible() {
		if w.IsSensitive() {
			w.SetState(enums.StatePrelight)
			if w.GetHasTooltip() {
				w.openTooltip()
			}
		}
		w.Invalidate()
		w.LogTrace("mouse enter - %v", w.ObjectInfo())
		return cenums.EVENT_STOP
	}
	return cenums.EVENT_PASS
}

func (w *CWidget) leave(_ []interface{}, _ ...interface{}) cenums.EventFlag {
	if w.IsDrawable() && w.IsVisible() {
		if w.HasState(enums.StatePrelight) {
			w.UnsetState(enums.StatePrelight)
		}
		w.closeTooltip()
		w.Invalidate()
		w.LogTrace("mouse leave - %v", w.ObjectInfo())
		return cenums.EVENT_STOP
	}
	return cenums.EVENT_PASS
}

// Whether the application will paint directly on the widget.
// Flags: Read / Write
// Default value: FALSE
const PropertyAppPaintable cdk.Property = "app-paintable"

// Whether the widget can be the default widget.
// Flags: Read / Write
// Default value: FALSE
const PropertyCanDefault cdk.Property = "can-default"

// Whether the widget can accept the input focus.
// Flags: Read / Write
// Default value: FALSE
const PropertyCanFocus cdk.Property = "can-focus"

// Whether the widget is part of a composite widget.
// Flags: Read
// Default value: FALSE
const PropertyCompositeChild cdk.Property = "composite-child"

// Whether or not the widget is double buffered.
// Flags: Read / Write
// Default value: TRUE
const PropertyDoubleBuffered cdk.Property = "double-buffered"

// The event mask that decides what kind of Events this widget gets.
// Flags: Read / Write
// Default value: GDK_STRUCTURE_MASK
const PropertyEvents cdk.Property = "events"

// The mask that decides what kind of extension events this widget gets.
// Flags: Read / Write
// Default value: GDK_EXTENSION_EVENTS_NONE
const PropertyExtensionEvents cdk.Property = "extension-events"

// Whether the widget is the default widget.
// Flags: Read / Write
// Default value: FALSE
const PropertyHasDefault cdk.Property = "has-default"

// Whether the widget has the input focus.
// Flags: Read / Write
// Default value: FALSE
const PropertyHasFocus cdk.Property = "has-focus"

// Enables or disables the emission of query-tooltip on widget . A
// value of TRUE indicates that widget can have a tooltip, in this case the
// widget will be queried using query-tooltip to determine whether it
// will provide a tooltip or not. Note that setting this property to TRUE for
// the first time will change the event masks of the Windows of this
// widget to include leave-notify and motion-notify events. This cannot and
// will not be undone when the property is set to FALSE again.
// Flags: Read / Write
// Default value: FALSE
const PropertyHasTooltip cdk.Property = "has-tooltip"

// Override for height request of the widget, or -1 if natural request should
// be used.
// Flags: Read / Write
// Allowed values: >= -1
// Default value: -1
const PropertyHeightRequest cdk.Property = "height-request"

// Whether the widget is the focus widget within the toplevel.
// Flags: Read / Write
// Default value: FALSE
const PropertyIsFocus cdk.Property = "is-focus"

// The name of the widget.
// Flags: Read / Write
// Default value: NULL
const PropertyName cdk.Property = "name"

// Whether ShowAll should not affect this widget.
// Flags: Read / Write
// Default value: FALSE
const PropertyNoShowAll cdk.Property = "no-show-all"

// The parent widget of this widget. Must be a Container widget.
// Flags: Read / Write
const PropertyParent cdk.Property = "parent"

// If TRUE, the widget will receive the default action when it is focused.
// Flags: Read / Write
// Default value: FALSE
const PropertyReceivesDefault cdk.Property = "receives-default"

// Whether the widget responds to input.
// Flags: Read / Write
// Default value: TRUE
const PropertySensitive cdk.Property = "sensitive"

// The style of the widget, which contains information about how it will look
// (colors etc).
// Flags: Read / Write
const PropertyStyle cdk.Property = "style"

// Sets the text of tooltip to be the given string, which is marked up with
// the Tango text markup language. Also see TooltipSetMarkup. This is
// a convenience property which will take care of getting the tooltip shown
// if the given string is not NULL: has-tooltip will automatically be
// set to TRUE and there will be taken care of query-tooltip in the
// default signal handler.
// Flags: Read / Write
// Default value: NULL
const PropertyTooltipMarkup cdk.Property = "tooltip-markup"

// Sets the text of tooltip to be the given string. Also see
// TooltipSetText. This is a convenience property which will take
// care of getting the tooltip shown if the given string is not NULL:
// has-tooltip will automatically be set to TRUE and there will be
// taken care of query-tooltip in the default signal handler.
// Flags: Read / Write
// Default value: NULL
const PropertyTooltipText cdk.Property = "tooltip-text"

// Whether the widget is visible.
// Flags: Read / Write
// Default value: FALSE
const PropertyVisible cdk.Property = "visible"

// Override for width request of the widget, or -1 if natural request should
// be used.
// Flags: Read / Write
// Allowed values: >= -1
// Default value: -1
const PropertyWidthRequest cdk.Property = "width-request"

// The widget's window if it is realized, NULL otherwise.
// Flags: Read
const PropertyWindow cdk.Property = "window"

const SignalAccelClosuresChanged cdk.Signal = "accel-closures-changed"

// Determines whether an accelerator that activates the signal identified by
// signal_id can currently be activated. This signal is present to allow
// applications and derived widgets to override the default Widget
// handling for determining whether an accelerator can be activated.
const SignalCanActivateAccel cdk.Signal = "can-activate-accel"

// The ::child-notify signal is emitted for each changed on an object. The
// signal's detail holds the property name.
// Listener function arguments:
//      pspec GParamSpec        the GParamSpec of the changed child property
const SignalChildNotify cdk.Signal = "child-notify"

// The ::client-event will be emitted when the widget 's window receives a
// message (via a ClientMessage event) from another application.
const SignalClientEvent cdk.Signal = "client-event"

// The ::composited-changed signal is emitted when the composited status of
// widget s screen changes. See ScreenIsComposited.
const SignalCompositedChanged cdk.Signal = "composited-changed"

// The ::configure-event signal will be emitted when the size, position or
// stacking of the widget 's window has changed. To receive this signal, the
// Window associated to the widget needs to enable the GDK_STRUCTURE_MASK
// mask. GDK will enable this mask automatically for all new windows.
const SignalConfigureEvent cdk.Signal = "configure-event"

// Emitted when a redirected window belonging to widget gets drawn into. The
// region/area members of the event shows what area of the redirected
// drawable was drawn into.
const SignalDamageEvent cdk.Signal = "damage-event"

// The ::delete-event signal is emitted if a user requests that a toplevel
// window is closed. The default handler for this signal destroys the window.
// Connecting HideOnDelete to this signal will cause the
// window to be hidden instead, so that it can later be shown again without
// reconstructing it.
const SignalDeleteEvent cdk.Signal = "delete-event"

// The ::destroy-event signal is emitted when a Window is destroyed. You
// rarely get this signal, because most widgets disconnect themselves from
// their window before they destroy it, so no widget owns the window at
// destroy time. To receive this signal, the Window associated to the
// widget needs to enable the GDK_STRUCTURE_MASK mask. GDK will enable this
// mask automatically for all new windows.
const SignalDestroyEvent cdk.Signal = "destroy-event"

// The ::direction-changed signal is emitted when the text direction of a
// widget changes.
// Listener function arguments:
//      previousDirection TextDirection the previous text direction of widget
const SignalDirectionChanged cdk.Signal = "direction-changed"

// The ::drag-begin signal is emitted on the drag source when a drag is
// started. A typical reason to connect to this signal is to set up a custom
// drag icon with DragSourceSetIcon. Note that some widgets set up a
// drag icon in the default handler of this signal, so you may have to use
// g_signal_connect_after to override what the default handler did.
// Listener function arguments:
//      dragContext DragContext the drag context
const SignalDragBegin cdk.Signal = "drag-begin"

// The ::drag-data-delete signal is emitted on the drag source when a drag
// with the action GDK_ACTION_MOVE is successfully completed. The signal
// handler is responsible for deleting the data that has been dropped. What
// "delete" means depends on the context of the drag operation.
// Listener function arguments:
//      dragContext DragContext the drag context
const SignalDragDataDelete cdk.Signal = "drag-data-delete"

// The ::drag-data-get signal is emitted on the drag source when the drop
// site requests the data which is dragged. It is the responsibility of the
// signal handler to fill data with the data in the format which is indicated
// by info . See SelectionDataSet and SelectionDataSetText.
// Listener function arguments:
//      dragContext DragContext the drag context
//      data SelectionData      the GtkSelectionData to be filled with the dragged data
//      info int        the info that has been registered with the target in the GtkTargetList
//      time int        the timestamp at which the data was requested
const SignalDragDataGet cdk.Signal = "drag-data-get"

// The ::drag-data-received signal is emitted on the drop site when the
// dragged data has been received. If the data was received in order to
// determine whether the drop will be accepted, the handler is expected to
// call DragStatus and not finish the drag. If the data was received
// in response to a drag-drop signal (and this is the last target to be
// received), the handler for this signal is expected to process the received
// data and then call DragFinish, setting the success parameter
// depending on whether the data was processed successfully. The handler may
// inspect and modify drag_context->action before calling DragFinish,
// e.g. to implement GDK_ACTION_ASK as shown in the following example:
// Listener function arguments:
//      dragContext DragContext the drag context
//      x int   where the drop happened
//      y int   where the drop happened
//      data SelectionData      the received data
//      info int        the info that has been registered with the target in the GtkTargetList
//      time int        the timestamp at which the data was received
const SignalDragDataReceived cdk.Signal = "drag-data-received"

// The ::drag-drop signal is emitted on the drop site when the user drops the
// data onto the widget. The signal handler must determine whether the cursor
// position is in a drop zone or not. If it is not in a drop zone, it returns
// FALSE and no further processing is necessary. Otherwise, the handler
// returns TRUE. In this case, the handler must ensure that DragFinish
// is called to let the source know that the drop is done. The call to
// DragFinish can be done either directly or in a
// drag-data-received handler which gets triggered by calling
// DragGetData to receive the data for one or more of the supported
// targets.
const SignalDragDrop cdk.Signal = "drag-drop"

// The ::drag-end signal is emitted on the drag source when a drag is
// finished. A typical reason to connect to this signal is to undo things
// done in drag-begin.
// Listener function arguments:
//      dragContext DragContext the drag context
const SignalDragEnd cdk.Signal = "drag-end"

// The ::drag-failed signal is emitted on the drag source when a drag has
// failed. The signal handler may hook custom code to handle a failed DND
// operation based on the type of error, it returns TRUE is the failure has
// been already handled (not showing the default "drag operation failed"
// animation), otherwise it returns FALSE.
const SignalDragFailed cdk.Signal = "drag-failed"

// The ::drag-leave signal is emitted on the drop site when the cursor leaves
// the widget. A typical reason to connect to this signal is to undo things
// done in drag-motion, e.g. undo highlighting with
// DragUnhighlight
// Listener function arguments:
//      dragContext DragContext the drag context
//      time int        the timestamp of the motion event
const SignalDragLeave cdk.Signal = "drag-leave"

// The drag-motion signal is emitted on the drop site when the user moves the
// cursor over the widget during a drag. The signal handler must determine
// whether the cursor position is in a drop zone or not. If it is not in a
// drop zone, it returns FALSE and no further processing is necessary.
// Otherwise, the handler returns TRUE. In this case, the handler is
// responsible for providing the necessary information for displaying
// feedback to the user, by calling DragStatus. If the decision
// whether the drop will be accepted or rejected can't be made based solely
// on the cursor position and the type of the data, the handler may inspect
// the dragged data by calling DragGetData and defer the
// DragStatus call to the drag-data-received handler. Note that
// you cannot not pass GTK_DEST_DEFAULT_DROP, GTK_DEST_DEFAULT_MOTION or
// GTK_DEST_DEFAULT_ALL to DragDestSet when using the drag-motion
// signal that way. Also note that there is no drag-enter signal. The drag
// receiver has to keep track of whether he has received any drag-motion
// signals since the last drag-leave and if not, treat the drag-motion
// signal as an "enter" signal. Upon an "enter", the handler will typically
// highlight the drop site with DragHighlight.
const SignalDragMotion cdk.Signal = "drag-motion"

// The ::enter-notify-event will be emitted when the pointer enters the
// widget 's window. To receive this signal, the Window associated to the
// widget needs to enable the GDK_ENTER_NOTIFY_MASK mask. This signal will be
// sent to the grab widget if there is one.
const SignalEnterNotifyEvent cdk.Signal = "enter-notify-event"

// The CTK main loop will emit three signals for each GDK event delivered to
// a widget: one generic ::event signal, another, more specific, signal that
// matches the type of event delivered (e.g. key-press-event) and
// finally a generic event-after signal.
const SignalEvent cdk.Signal = "event"

// After the emission of the event signal and (optionally) the second
// more specific signal, ::event-after will be emitted regardless of the
// previous two signals handlers return values.
// Listener function arguments:
//      event Event     the GdkEvent which triggered this signal
const SignalEventAfter cdk.Signal = "event-after"

// The ::expose-event signal is emitted when an area of a previously obscured
// Window is made visible and needs to be redrawn. GTK_NO_WINDOW widgets
// will get a synthesized event from their parent widget. To receive this
// signal, the Window associated to the widget needs to enable the
// GDK_EXPOSURE_MASK mask. Note that the ::expose-event signal has been
// replaced by a ::draw signal in CTK 3. The CTK 3 migration guide for
// hints on how to port from ::expose-event to ::draw.
const SignalExposeEvent cdk.Signal = "expose-event"

const SignalFocus cdk.Signal = "focus"

// The ::focus-in-event signal will be emitted when the keyboard focus enters
// the widget 's window. To receive this signal, the Window associated to
// the widget needs to enable the GDK_FOCUS_CHANGE_MASK mask.
const SignalFocusInEvent cdk.Signal = "focus-in-event"

// The ::focus-out-event signal will be emitted when the keyboard focus
// leaves the widget 's window. To receive this signal, the Window
// associated to the widget needs to enable the GDK_FOCUS_CHANGE_MASK mask.
const SignalFocusOutEvent cdk.Signal = "focus-out-event"

// Emitted when a pointer or keyboard grab on a window belonging to widget
// gets broken. On X11, this happens when the grab window becomes unviewable
// (i.e. it or one of its ancestors is unmapped), or if the same application
// grabs the pointer or keyboard again.
const SignalGrabBrokenEvent cdk.Signal = "grab-broken-event"

const SignalGrabFocus cdk.Signal = "grab-focus"

// The ::grab-notify signal is emitted when a widget becomes shadowed by a
// CTK grab (not a pointer or keyboard grab) on another widget, or when it
// becomes unshadowed due to a grab being removed. A widget is shadowed by a
// GrabAdd when the topmost grab widget in the grab stack of its
// window group is not its ancestor.
// Listener function arguments:
//      wasGrabbed bool FALSE if the widget becomes shadowed, TRUE if it becomes unshadowed
const SignalGrabNotify cdk.Signal = "grab-notify"

const SignalHide cdk.Signal = "hide"

// The ::hierarchy-changed signal is emitted when the anchored state of a
// widget changes. A widget is anchored when its toplevel ancestor is a
// Window. This signal is emitted when a widget changes from un-anchored
// to anchored or vice-versa.
const SignalHierarchyChanged cdk.Signal = "hierarchy-changed"

// The ::key-press-event signal is emitted when a key is pressed. To receive
// this signal, the Window associated to the widget needs to enable the
// GDK_KEY_PRESS_MASK mask. This signal will be sent to the grab widget if
// there is one.
const SignalKeyPressEvent cdk.Signal = "key-press-event"

// The ::key-release-event signal is emitted when a key is pressed. To
// receive this signal, the Window associated to the widget needs to
// enable the GDK_KEY_RELEASE_MASK mask. This signal will be sent to the grab
// widget if there is one.
const SignalKeyReleaseEvent cdk.Signal = "key-release-event"

// Gets emitted if keyboard navigation fails. See KeynavFailed
// for details.
const SignalKeynavFailed cdk.Signal = "keynav-failed"

// The ::leave-notify-event will be emitted when the pointer leaves the
// widget 's window. To receive this signal, the Window associated to the
// widget needs to enable the GDK_LEAVE_NOTIFY_MASK mask. This signal will be
// sent to the grab widget if there is one.
const SignalLeaveNotifyEvent cdk.Signal = "leave-notify-event"

const SignalMap cdk.Signal = "map"

// The ::map-event signal will be emitted when the widget 's window is
// mapped. A window is mapped when it becomes visible on the screen. To
// receive this signal, the Window associated to the widget needs to
// enable the GDK_STRUCTURE_MASK mask. GDK will enable this mask
// automatically for all new windows.
const SignalMapEvent cdk.Signal = "map-event"

const SignalMnemonicActivate cdk.Signal = "mnemonic-activate"

// The ::motion-notify-event signal is emitted when the pointer moves over
// the widget's Window. To receive this signal, the Window associated
// to the widget needs to enable the GDK_POINTER_MOTION_MASK mask. This
// signal will be sent to the grab widget if there is one.
const SignalMotionNotifyEvent cdk.Signal = "motion-notify-event"

// Listener function arguments:
//      direction DirectionType
const SignalMoveFocus cdk.Signal = "move-focus"

// The ::no-expose-event will be emitted when the widget 's window is drawn
// as a copy of another Drawable (with DrawDrawable or
// WindowCopyArea) which was completely unobscured. If the source
// window was partially obscured EventExpose events will be generated for
// those areas.
const SignalNoExposeEvent cdk.Signal = "no-expose-event"

// The ::parent-set signal is emitted when a new parent has been set on a
// widget.
const SignalParentSet cdk.Signal = "parent-set"

// This signal gets emitted whenever a widget should pop up a context menu.
// This usually happens through the standard key binding mechanism; by
// pressing a certain key while a widget is focused, the user can cause the
// widget to pop up a menu. For example, the Entry widget creates a menu
// with clipboard commands. See the section called Implement
// Widget::popup_menu for an example of how to use this signal.
const SignalPopupMenu cdk.Signal = "popup-menu"

// The ::property-notify-event signal will be emitted when a property on the
// widget 's window has been changed or deleted. To receive this signal, the
// Window associated to the widget needs to enable the
// GDK_PROPERTY_CHANGE_MASK mask.
const SignalPropertyNotifyEvent cdk.Signal = "property-notify-event"

// To receive this signal the Window associated to the widget needs to
// enable the GDK_PROXIMITY_IN_MASK mask. This signal will be sent to the
// grab widget if there is one.
const SignalProximityInEvent cdk.Signal = "proximity-in-event"

// To receive this signal the Window associated to the widget needs to
// enable the GDK_PROXIMITY_OUT_MASK mask. This signal will be sent to the
// grab widget if there is one.
const SignalProximityOutEvent cdk.Signal = "proximity-out-event"

// Emitted when has-tooltip is TRUE and the gtk-tooltip-timeout
// has expired with the cursor hovering "above" widget ; or emitted when
// widget got focus in keyboard mode. Using the given coordinates, the signal
// handler should determine whether a tooltip should be shown for widget . If
// this is the case TRUE should be returned, FALSE otherwise. Note that if
// keyboard_mode is TRUE, the values of x and y are undefined and should not
// be used. The signal handler is free to manipulate tooltip with the
// therefore destined function calls.
const SignalQueryTooltip cdk.Signal = "query-tooltip"

const SignalRealize cdk.Signal = "realize"

// The ::screen-changed signal gets emitted when the screen of a widget has
// changed.
// Listener function arguments:
//      previousScreen Screen   the previous screen, or NULL if the widget was not associated with a screen before.
const SignalScreenChanged cdk.Signal = "screen-changed"

// The ::scroll-event signal is emitted when a button in the 4 to 7 range is
// pressed. Wheel mice are usually configured to generate button press events
// for buttons 4 and 5 when the wheel is turned. To receive this signal, the
// Window associated to the widget needs to enable the
// GDK_BUTTON_PRESS_MASK mask. This signal will be sent to the grab widget if
// there is one.
const SignalScrollEvent cdk.Signal = "scroll-event"

// The ::selection-clear-event signal will be emitted when the the widget 's
// window has lost ownership of a selection.
const SignalSelectionClearEvent cdk.Signal = "selection-clear-event"

// Listener function arguments:
//      data SelectionData
//      info int
//      time int
const SignalSelectionGet cdk.Signal = "selection-get"

const SignalSelectionNotifyEvent cdk.Signal = "selection-notify-event"

// Listener function arguments:
//      data SelectionData
//      time int
const SignalSelectionReceived cdk.Signal = "selection-received"

// The ::selection-request-event signal will be emitted when another client
// requests ownership of the selection owned by the widget 's window.
const SignalSelectionRequestEvent cdk.Signal = "selection-request-event"

const SignalShow cdk.Signal = "show"

const SignalShowHelp cdk.Signal = "show-help"

// Listener function arguments:
//      allocation Rectangle
const SignalSizeAllocate cdk.Signal = "size-allocate"

// Listener function arguments:
//      requisition Requisition
const SignalSizeRequest cdk.Signal = "size-request"

// The ::state-changed signal is emitted when the widget state changes. See
// GetState.
// Listener function arguments:
//      state StateType the previous state
const SignalStateChanged cdk.Signal = "state-changed"

// The ::style-set signal is emitted when a new style has been set on a
// widget. Note that style-modifying functions like ModifyBase
// also cause this signal to be emitted.
// Listener function arguments:
//      previousStyle Style     the previous style, or NULL if the widget just got its initial style.
const SignalStyleSet cdk.Signal = "style-set"

const SignalUnmap cdk.Signal = "unmap"

// The ::unmap-event signal will be emitted when the widget 's window is
// unmapped. A window is unmapped when it becomes invisible on the screen. To
// receive this signal, the Window associated to the widget needs to
// enable the GDK_STRUCTURE_MASK mask. GDK will enable this mask
// automatically for all new windows.
const SignalUnmapEvent cdk.Signal = "unmap-event"

const SignalUnrealize cdk.Signal = "unrealize"

// The ::visibility-notify-event will be emitted when the widget 's window is
// obscured or unobscured. To receive this signal the Window associated to
// the widget needs to enable the GDK_VISIBILITY_NOTIFY_MASK mask.
const SignalVisibilityNotifyEvent cdk.Signal = "visibility-notify-event"

// The ::window-state-event will be emitted when the state of the toplevel
// window associated to the widget changes. To receive this signal the
// Window associated to the widget needs to enable the GDK_STRUCTURE_MASK
// mask. GDK will enable this mask automatically for all new windows.
const SignalWindowStateEvent cdk.Signal = "window-state-event"

const SignalGetThemeRequest = "get-theme-request"

const WidgetLostFocusHandle = "widget-lost-focus-handler"

const WidgetGainedFocusHandle = "widget-gained-focus-handler"

const WidgetEnterHandle = "widget-enter-handler"

const WidgetLeaveHandle = "widget-leave-handler"

const WidgetActivateHandle = "widget-activate-handler"

const WidgetTooltipWindowEventHandle = "widget-tooltip-window-event-handler"

const CssPropertyClass cdk.Property = "class"

const CssPropertyWidth cdk.Property = "width"

const CssPropertyHeight cdk.Property = "height"

const CssPropertyColor cdk.Property = "color"

const CssPropertyBackgroundColor cdk.Property = "background-color"

const CssPropertyBorder cdk.Property = "border"

const CssPropertyBorderColor cdk.Property = "border-color"

const CssPropertyBorderBackgroundColor cdk.Property = "border-background-color"

const CssPropertyBold cdk.Property = "bold"

const CssPropertyBlink cdk.Property = "blink"

const CssPropertyReverse cdk.Property = "reverse"

const CssPropertyUnderline cdk.Property = "underline"

const CssPropertyDim cdk.Property = "dim"

const CssPropertyItalic cdk.Property = "italic"

const CssPropertyStrike cdk.Property = "strike"