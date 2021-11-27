package ctk

import (
	"github.com/go-curses/cdk"
	"github.com/go-curses/ctk/lib/enums"
)

const TypeAction cdk.CTypeTag = "ctk-action"

func init() {
	_ = cdk.TypesManager.AddType(TypeAction, func() interface{} { return MakeAction() })
}

// Action Hierarchy:
//	Object
//	  +- Action
//	    +- ToggleAction
//	    +- RecentAction
//
// Actions represent operations that the user can perform, along with some
// information how it should be presented in the interface. Each action
// provides methods to create icons, menu items and toolbar items
// representing itself. As well as the callback that is called when the
// action gets activated, the following also gets associated with the action:
// a name (not translated, for path lookup) a label (translated, for display)
// an accelerator whether label indicates a stock id a tooltip (optional,
// translated) a toolbar label (optional, shorter than label) The action will
// also have some state information: visible (shown/hidden) sensitive
// (enabled/disabled) Apart from regular actions, there are toggle actions,
// which can be toggled between two states and radio actions, of which only
// one in a group can be in the "active" state. Other actions can be
// implemented as Action subclasses. Each action can have one or more proxy
// menu item, toolbar button or other proxy widgets. Proxies mirror the state
// of the action (text label, tooltip, icon, visible, sensitive, etc), and
// should change when the action's state changes. When the proxy is
// activated, it should activate its action.
type Action interface {
	Object

	Init() (already bool)
	GetName() (value string)
	IsSensitive() (value bool)
	GetSensitive() (value bool)
	SetSensitive(sensitive bool)
	IsVisible() (value bool)
	GetVisible() (value bool)
	SetVisible(visible bool)
	Activate()
	CreateMenuItem() (value Widget)
	CreateToolItem() (value Widget)
	CreateMenu() (value Widget)
	GetProxies() (value []interface{})
	ConnectAccelerator()
	DisconnectAccelerator()
	UnblockActivate()
	GetAlwaysShowImage() (value bool)
	SetAlwaysShowImage(alwaysShow bool)
	GetAccelPath() (value string)
	SetAccelPath(accelPath string)
	GetAccelClosure() (value enums.GClosure)
	SetAccelGroup(accelGroup AccelGroup)
	SetLabel(label string)
	GetLabel() (value string)
	SetShortLabel(shortLabel string)
	GetShortLabel() (value string)
	SetTooltip(tooltip string)
	GetTooltip() (value string)
	SetStockId(stockId StockID)
	GetStockId() (value StockID)
	SetIcon(icon rune)
	GetIcon() (value rune)
	SetIconName(iconName string)
	GetIconName() (value string)
	SetVisibleHorizontal(visibleHorizontal bool)
	GetVisibleHorizontal() (value bool)
	SetVisibleVertical(visibleVertical bool)
	GetVisibleVertical() (value bool)
	SetIsImportant(isImportant bool)
	GetIsImportant() (value bool)
}

// The CAction structure implements the Action interface and is
// exported to facilitate type embedding with custom implementations. No member
// variables are exported as the interface methods are the only intended means
// of interacting with Action objects
type CAction struct {
	CObject
}

// Default constructor for Action objects
func MakeAction() Action {
	return NewAction("", "", "", "")
}

// Constructor for Action objects
func NewAction(name string, label string, tooltip string, stockId string) (value Action) {
	a := new(CAction)
	a.Init()
	return a
}

// Action object initialization. This must be called at least once to setup
// the necessary defaults and allocate any memory structures. Calling this more
// than once is safe though unnecessary. Only the first call will result in any
// effect upon the Action instance
func (a *CAction) Init() (already bool) {
	if a.InitTypeItem(TypeAction, a) {
		return true
	}
	a.CObject.Init()
	_ = a.InstallProperty(PropertyAlwaysShowImage, cdk.BoolProperty, true, false)
	_ = a.InstallProperty(PropertyIcon, cdk.StructProperty, true, nil)
	_ = a.InstallProperty(PropertyHideIfEmpty, cdk.BoolProperty, true, true)
	_ = a.InstallProperty(PropertyIconName, cdk.StringProperty, true, nil)
	_ = a.InstallProperty(PropertyIsImportant, cdk.BoolProperty, true, false)
	_ = a.InstallProperty(PropertyLabel, cdk.StringProperty, true, nil)
	_ = a.InstallProperty(PropertyName, cdk.StringProperty, true, nil)
	_ = a.InstallProperty(PropertySensitive, cdk.BoolProperty, true, true)
	_ = a.InstallProperty(PropertyShortLabel, cdk.StringProperty, true, nil)
	_ = a.InstallProperty(PropertyStockId, cdk.StructProperty, true, nil)
	_ = a.InstallProperty(PropertyTooltip, cdk.StringProperty, true, nil)
	_ = a.InstallProperty(PropertyVisible, cdk.BoolProperty, true, true)
	_ = a.InstallProperty(PropertyVisibleHorizontal, cdk.BoolProperty, true, true)
	_ = a.InstallProperty(PropertyVisibleOverflown, cdk.BoolProperty, true, true)
	_ = a.InstallProperty(PropertyVisibleVertical, cdk.BoolProperty, true, true)
	return false
}

// Returns the name of the action.
// Parameters:
// 	action	the action object
// Returns:
// 	the name of the action. The string belongs to CTK and should
// 	not be freed.
func (a *CAction) GetName() (value string) {
	var err error
	if value, err = a.GetStringProperty(PropertyName); err != nil {
		a.LogErr(err)
	}
	return
}

// Returns whether the action is effectively sensitive.
// Parameters:
// 	action	the action object
// Returns:
// 	TRUE if the action and its associated action group are both
// 	sensitive.
func (a *CAction) IsSensitive() (value bool) {
	return false
}

// Returns whether the action itself is sensitive. Note that this doesn't
// necessarily mean effective sensitivity. See IsSensitive for
// that.
// Parameters:
// 	action	the action object
// Returns:
// 	TRUE if the action itself is sensitive.
func (a *CAction) GetSensitive() (value bool) {
	var err error
	if value, err = a.GetBoolProperty(PropertySensitive); err != nil {
		a.LogErr(err)
	}
	return
}

// Sets the ::sensitive property of the action to sensitive . Note that this
// doesn't necessarily mean effective sensitivity. See
// IsSensitive for that.
// Parameters:
// 	action	the action object
// 	sensitive	TRUE to make the action sensitive
func (a *CAction) SetSensitive(sensitive bool) {
	if err := a.SetBoolProperty(PropertySensitive, sensitive); err != nil {
		a.LogErr(err)
	}
}

// Returns whether the action is effectively visible.
// Parameters:
// 	action	the action object
// Returns:
// 	TRUE if the action and its associated action group are both
// 	visible.
func (a *CAction) IsVisible() (value bool) {
	return false
}

// Returns whether the action itself is visible. Note that this doesn't
// necessarily mean effective visibility. See IsSensitive for
// that.
// Parameters:
// 	action	the action object
// Returns:
// 	TRUE if the action itself is visible.
func (a *CAction) GetVisible() (value bool) {
	var err error
	if value, err = a.GetBoolProperty(PropertyVisible); err != nil {
		a.LogErr(err)
	}
	return
}

// Sets the ::visible property of the action to visible . Note that this
// doesn't necessarily mean effective visibility. See IsVisible
// for that.
// Parameters:
// 	action	the action object
// 	visible	TRUE to make the action visible
func (a *CAction) SetVisible(visible bool) {
	if err := a.SetBoolProperty(PropertyVisible, visible); err != nil {
		a.LogErr(err)
	}
}

// Emits the "activate" signal on the specified action, if it isn't
// insensitive. This gets called by the proxy widgets when they get
// activated. It can also be used to manually activate an action.
// Parameters:
// 	action	the action object
func (a *CAction) Activate() {}

// Creates a menu item widget that proxies for the given action.
// Parameters:
// 	action	the action object
// Returns:
// 	a menu item connected to the action.
// 	[transfer none]
func (a *CAction) CreateMenuItem() (value Widget) {
	return nil
}

// Creates a toolbar item widget that proxies for the given action.
// Parameters:
// 	action	the action object
// Returns:
// 	a toolbar item connected to the action.
// 	[transfer none]
func (a *CAction) CreateToolItem() (value Widget) {
	return nil
}

// If action provides a Menu widget as a submenu for the menu item or the
// toolbar item it creates, this function returns an instance of that menu.
// Returns:
// 	the menu item provided by the action, or NULL.
// 	[transfer none]
func (a *CAction) CreateMenu() (value Widget) {
	return nil
}

// Returns the proxy widgets for an action. See also WidgetGetAction.
// Parameters:
// 	action	the action object
// Returns:
// 	a GSList of proxy widgets. The list is owned by CTK and must
// 	not be modified.
// 	[element-type Widget][transfer none]
func (a *CAction) GetProxies() (value []interface{}) {
	return nil
}

// Installs the accelerator for action if action has an accel path and group.
// See SetAccelPath and SetAccelGroup Since
// multiple proxies may independently trigger the installation of the
// accelerator, the action counts the number of times this function has been
// called and doesn't remove the accelerator until
// DisconnectAccelerator has been called as many times.
func (a *CAction) ConnectAccelerator() {}

// Undoes the effect of one call to ConnectAccelerator.
func (a *CAction) DisconnectAccelerator() {}

// Reenable activation signals from the action
func (a *CAction) UnblockActivate() {}

// Returns whether action 's menu item proxies will ignore the
// “gtk-menu-images” setting and always show their image, if available.
// Returns:
// 	TRUE if the menu item proxies will always show their image
func (a *CAction) GetAlwaysShowImage() (value bool) {
	var err error
	if value, err = a.GetBoolProperty(PropertyAlwaysShowImage); err != nil {
		a.LogErr(err)
	}
	return
}

// Sets whether action 's menu item proxies will ignore the
// “gtk-menu-images” setting and always show their image, if available.
// Use this if the menu item would be useless or hard to use without their
// image.
// Parameters:
// 	alwaysShow	TRUE if menuitem proxies should always show their image
func (a *CAction) SetAlwaysShowImage(alwaysShow bool) {
	if err := a.SetBoolProperty(PropertyAlwaysShowImage, alwaysShow); err != nil {
		a.LogErr(err)
	}
}

// Returns the accel path for this action.
// Parameters:
// 	action	the action object
// Returns:
// 	the accel path for this action, or NULL if none is set. The
// 	returned string is owned by CTK and must not be freed or
// 	modified.
func (a *CAction) GetAccelPath() (value string) {
	return ""
}

// Sets the accel path for this action. All proxy widgets associated with the
// action will have this accel path, so that their accelerators are
// consistent.
// Parameters:
// 	action	the action object
// 	accelPath	the accelerator path
func (a *CAction) SetAccelPath(accelPath string) {}

// Returns the accel closure for this action.
// Parameters:
// 	action	the action object
// Returns:
// 	the accel closure for this action.
func (a *CAction) GetAccelClosure() (value enums.GClosure) {
	return nil
}

// Sets the AccelGroup in which the accelerator for this action will be
// installed.
// Parameters:
// 	action	the action object
// 	accelGroup	a AccelGroup or NULL.
func (a *CAction) SetAccelGroup(accelGroup AccelGroup) {}

// Sets the label of action .
// Parameters:
// 	label	the label text to set
func (a *CAction) SetLabel(label string) {
	if err := a.SetStringProperty(PropertyLabel, label); err != nil {
		a.LogErr(err)
	}
}

// Gets the label text of action .
// Returns:
// 	the label text
func (a *CAction) GetLabel() (value string) {
	var err error
	if value, err = a.GetStringProperty(PropertyLabel); err != nil {
		a.LogErr(err)
	}
	return
}

// Sets a shorter label text on action .
// Parameters:
// 	shortLabel	the label text to set
func (a *CAction) SetShortLabel(shortLabel string) {
	if err := a.SetStringProperty(PropertyShortLabel, shortLabel); err != nil {
		a.LogErr(err)
	}
}

// Gets the short label text of action .
// Returns:
// 	the short label text.
func (a *CAction) GetShortLabel() (value string) {
	var err error
	if value, err = a.GetStringProperty(PropertyShortLabel); err != nil {
		a.LogErr(err)
	}
	return
}

// Sets the tooltip text on action
// Parameters:
// 	tooltip	the tooltip text
func (a *CAction) SetTooltip(tooltip string) {
	if err := a.SetStringProperty(PropertyTooltip, tooltip); err != nil {
		a.LogErr(err)
	}
}

// Gets the tooltip text of action .
// Returns:
// 	the tooltip text
func (a *CAction) GetTooltip() (value string) {
	var err error
	if value, err = a.GetStringProperty(PropertyTooltip); err != nil {
		a.LogErr(err)
	}
	return
}

// SetStockId updates the stock id on the action.
//
// Parameters:
// 	stockId	the stock id
func (a *CAction) SetStockId(stockId StockID) {
	if ValidStockId(stockId) {
		if err := a.SetStructProperty(PropertyStockId, stockId); err != nil {
			a.LogErr(err)
		}
	} else {
		a.LogError("unknown StockID: %v", stockId)
	}
}

// GetStockId returns the stock id of the action.
func (a *CAction) GetStockId() (value StockID) {
	var err error
	var v interface{}
	if v, err = a.GetStructProperty(PropertyStockId); err != nil {
		a.LogErr(err)
	} else {
		if val, ok := v.(StockID); ok {
			value = val
		} else {
			a.LogError("value stored in %v is not of StockID type: %v (%T)", PropertyStockId, v, v)
		}
	}
	return
}

// SetIcon updates the icon rune of the action.
//
// Parameters:
// 	icon	the rune to set
func (a *CAction) SetIcon(icon rune) {
	if err := a.SetStructProperty(PropertyIcon, icon); err != nil {
		a.LogErr(err)
	}
}

// GetIcon returns the icon rune of the action.
func (a *CAction) GetIcon() (value rune) {
	var err error
	var v interface{}
	if v, err = a.GetStructProperty(PropertyIcon); err != nil {
		a.LogErr(err)
	} else if val, ok := v.(rune); ok {
		value = val
	}
	return
}

// SetIconName updates the icon name on the action.
//
// Parameters:
// 	iconName	the icon name to set
func (a *CAction) SetIconName(iconName string) {
	if err := a.SetStringProperty(PropertyIconName, iconName); err != nil {
		a.LogErr(err)
	}
}

// GetIconName returns the icon name of the action.
func (a *CAction) GetIconName() (value string) {
	var err error
	if value, err = a.GetStringProperty(PropertyIconName); err != nil {
		a.LogErr(err)
	}
	return
}

// SetVisibleHorizontal updates whether action is visible when horizontal.
//
// Parameters:
// 	visibleHorizontal	whether the action is visible horizontally
func (a *CAction) SetVisibleHorizontal(visibleHorizontal bool) {
	if err := a.SetBoolProperty(PropertyVisibleHorizontal, visibleHorizontal); err != nil {
		a.LogErr(err)
	}
}

// GetVisibleHorizontal checks whether the action is visible when horizontal.
func (a *CAction) GetVisibleHorizontal() (value bool) {
	var err error
	if value, err = a.GetBoolProperty(PropertyVisibleHorizontal); err != nil {
		a.LogErr(err)
	}
	return
}

// SetVisibleVertical updates whether action is visible when vertical.
//
// Parameters:
// 	visibleVertical	whether the action is visible vertically
func (a *CAction) SetVisibleVertical(visibleVertical bool) {
	if err := a.SetBoolProperty(PropertyVisibleVertical, visibleVertical); err != nil {
		a.LogErr(err)
	}
}

// GetVisibleVertical checks whether the action is visible when vertical.
func (a *CAction) GetVisibleVertical() (value bool) {
	var err error
	if value, err = a.GetBoolProperty(PropertyVisibleVertical); err != nil {
		a.LogErr(err)
	}
	return
}

// SetIsImportant updates whether the action is important, this attribute is
// used primarily by toolbar items to decide whether to show a label or not.
//
// Parameters:
// 	isImportant	TRUE to make the action important
func (a *CAction) SetIsImportant(isImportant bool) {
	if err := a.SetBoolProperty(PropertyIsImportant, isImportant); err != nil {
		a.LogErr(err)
	}
}

// GetIsImportant returns whether action is important or not.
func (a *CAction) GetIsImportant() (value bool) {
	var err error
	if value, err = a.GetBoolProperty(PropertyIsImportant); err != nil {
		a.LogErr(err)
	}
	return
}

// If TRUE, the action's menu item proxies will ignore the
// “gtk-menu-images” setting and always show their image, if available.
// Use this property if the menu item would be useless or hard to use without
// their image.
// Flags: Read / Write / Construct
// Default value: FALSE
const PropertyAlwaysShowImage cdk.Property = "always-show-image"

// The Icon displayed in the Action. Note that the stock icon is
// preferred, if the “stock-id” property holds the id of an existing
// stock icon. This is an appearance property and thus only applies if
// “use-action-appearance” is TRUE.
// Flags: Read / Write
const PropertyActionIcon cdk.Property = "action-icon"

// When TRUE, empty menu proxies for this action are hidden.
// Flags: Read / Write
// Default value: TRUE
const PropertyHideIfEmpty cdk.Property = "hide-if-empty"

// The name of the icon from the icon theme. Note that the stock icon is
// preferred, if the “stock-id” property holds the id of an existing
// stock icon, and the GIcon is preferred if the “gicon” property is set.
// This is an appearance property and thus only applies if
// “use-action-appearance” is TRUE.
// Flags: Read / Write
// Default value: NULL
const PropertyActionIconName cdk.Property = "icon-name"

// Whether the action is considered important. When TRUE, toolitem proxies
// for this action show text in GTK_TOOLBAR_BOTH_HORIZ mode.
// Flags: Read / Write
// Default value: FALSE
const PropertyIsImportant cdk.Property = "is-important"

// The label used for menu items and buttons that activate this action. If
// the label is NULL, CTK uses the stock label specified via the stock-id
// property. This is an appearance property and thus only applies if
// “use-action-appearance” is TRUE.
// Flags: Read / Write
// Default value: NULL
const PropertyActionLabel cdk.Property = "label"

// A unique name for the action.
// Flags: Read / Write / Construct Only
// Default value: NULL
const PropertyActionName cdk.Property = "name"

// Whether the action is enabled.
// Flags: Read / Write
// Default value: TRUE
const PropertyActionSensitive cdk.Property = "sensitive"

// A shorter label that may be used on toolbar buttons. This is an appearance
// property and thus only applies if “use-action-appearance” is TRUE.
// Flags: Read / Write
// Default value: NULL
const PropertyShortLabel cdk.Property = "short-label"

// The stock icon displayed in widgets representing this action. This is an
// appearance property and thus only applies if “use-action-appearance”
// is TRUE.
// Flags: Read / Write
// Default value: NULL
const PropertyStockId cdk.Property = "stock-id"

// A tooltip for this action.
// Flags: Read / Write
// Default value: NULL
const PropertyTooltip cdk.Property = "tooltip"

// Whether the action is visible.
// Flags: Read / Write
// Default value: TRUE
const PropertyActionVisible cdk.Property = "visible"

// Whether the toolbar item is visible when the toolbar is in a horizontal
// orientation.
// Flags: Read / Write
// Default value: TRUE
const PropertyVisibleHorizontal cdk.Property = "visible-horizontal"

// When TRUE, toolitem proxies for this action are represented in the toolbar
// overflow menu.
// Flags: Read / Write
// Default value: TRUE
const PropertyVisibleOverflown cdk.Property = "visible-overflown"

// Whether the toolbar item is visible when the toolbar is in a vertical
// orientation.
// Flags: Read / Write
// Default value: TRUE
const PropertyVisibleVertical cdk.Property = "visible-vertical"

// The "activate" signal is emitted when the action is activated.
const SignalActionActivate cdk.Signal = "activate"
