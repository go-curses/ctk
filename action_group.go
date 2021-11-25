package ctk

import (
	"github.com/go-curses/cdk"
	"github.com/go-curses/ctk/lib/enums"
)

// CDK type-tag for ActionGroup objects
const TypeActionGroup cdk.CTypeTag = "ctk-action-group"

func init() {
	_ = cdk.TypesManager.AddType(TypeActionGroup, func() interface{} { return MakeActionGroup() })
}

// ActionGroup Hierarchy:
//	Object
//	  +- ActionGroup
type ActionGroup interface {
	Object
	Buildable

	Init() (already bool)
	GetName() (value string)
	GetSensitive() (value bool)
	SetSensitive(sensitive bool)
	GetVisible() (value bool)
	SetVisible(visible bool)
	GetAction(actionName string) (value Action)
	ListActions() (value []Action)
	AddAction(action Action)
	AddActionWithAccel(action Action, accelerator string)
	RemoveAction(action Action)
	AddActions(entries []ActionEntry, nEntries int, userData interface{})
	AddActionsFull(entries []ActionEntry, nEntries int, userData interface{}, destroy GDestroyNotify)
	AddToggleActions(entries []ToggleActionEntry, nEntries int, userData interface{})
	AddToggleActionsFull(entries []ToggleActionEntry, nEntries int, userData interface{}, destroy GDestroyNotify)
	AddRadioActions(entries []RadioActionEntry, nEntries int, value int, onChange enums.GCallback, userData interface{})
	AddRadioActionsFull(entries []RadioActionEntry, nEntries int, value int, onChange enums.GCallback, userData interface{}, destroy GDestroyNotify)
	SetTranslateFunc(fn TranslateFunc, data interface{}, notify GDestroyNotify)
	SetTranslationDomain(domain string)
	TranslateString(string string) (value string)
}

// The CActionGroup structure implements the ActionGroup interface and is
// exported to facilitate type embedding with custom implementations. No member
// variables are exported as the interface methods are the only intended means
// of interacting with ActionGroup objects.
type CActionGroup struct {
	CObject
}

// MakeActionGroup is used by the Buildable system to construct a new ActionGroup.
func MakeActionGroup() *CActionGroup {
	return NewActionGroup("")
}

// NewActionGroup is the constructor for new ActionGroup instances.
func NewActionGroup(name string) (value *CActionGroup) {
	a := new(CActionGroup)
	a.Init()
	return a
}

// Init initializes an ActionGroup object. This must be called at least once to
// set up the necessary defaults and allocate any memory structures. Calling
// this more than once is safe though unnecessary. Only the first call will
// result in any effect upon the ActionGroup instance. Init is used in the
// NewActionGroup constructor and only necessary when implementing a derivative
// ActionGroup type.
func (a *CActionGroup) Init() (already bool) {
	if a.InitTypeItem(TypeActionGroup, a) {
		return true
	}
	a.CObject.Init()
	return false
}

// Gets the name of the action group.
// Parameters:
// 	actionGroup	the action group
// Returns:
// 	the name of the action group.
func (a *CActionGroup) GetName() (value string) {
	var err error
	if value, err = a.GetStringProperty(PropertyName); err != nil {
		a.LogErr(err)
	}
	return
}

// Returns TRUE if the group is sensitive. The constituent actions can only
// be logically sensitive (see ActionIsSensitive) if they are
// sensitive (see ActionGetSensitive) and their group is sensitive.
// Parameters:
// 	actionGroup	the action group
// Returns:
// 	TRUE if the group is sensitive.
func (a *CActionGroup) GetSensitive() (value bool) {
	var err error
	if value, err = a.GetBoolProperty(PropertySensitive); err != nil {
		a.LogErr(err)
	}
	return
}

// Changes the sensitivity of action_group
// Parameters:
// 	actionGroup	the action group
// 	sensitive	new sensitivity
func (a *CActionGroup) SetSensitive(sensitive bool) {
	if err := a.SetBoolProperty(PropertySensitive, sensitive); err != nil {
		a.LogErr(err)
	}
}

// Returns TRUE if the group is visible. The constituent actions can only be
// logically visible (see ActionIsVisible) if they are visible (see
// ActionGetVisible) and their group is visible.
// Parameters:
// 	actionGroup	the action group
// Returns:
// 	TRUE if the group is visible.
func (a *CActionGroup) GetVisible() (value bool) {
	var err error
	if value, err = a.GetBoolProperty(PropertyVisible); err != nil {
		a.LogErr(err)
	}
	return
}

// Changes the visible of action_group .
// Parameters:
// 	actionGroup	the action group
// 	visible	new visiblity
func (a *CActionGroup) SetVisible(visible bool) {
	if err := a.SetBoolProperty(PropertyVisible, visible); err != nil {
		a.LogErr(err)
	}
}

// Looks up an action in the action group by name.
// Parameters:
// 	actionGroup	the action group
// 	actionName	the name of the action
// Returns:
// 	the action, or NULL if no action by that name exists.
// 	[transfer none]
func (a *CActionGroup) GetAction(actionName string) (value Action) {
	return nil
}

// Lists the actions in the action group.
// Parameters:
// 	actionGroup	the action group
// Returns:
// 	an allocated list of the action objects in the action group.
// 	[element-type Action][transfer container]
func (a *CActionGroup) ListActions() (value []Action) {
	return nil
}

// Adds an action object to the action group. Note that this function does
// not set up the accel path of the action, which can lead to problems if a
// user tries to modify the accelerator of a menuitem associated with the
// action. Therefore you must either set the accel path yourself with
// ActionSetAccelPath, or use AddActionWithAccel
// (..., NULL).
// Parameters:
// 	actionGroup	the action group
// 	action	an action
func (a *CActionGroup) AddAction(action Action) {}

// Adds an action object to the action group and sets up the accelerator. If
// accelerator is NULL, attempts to use the accelerator associated with the
// stock_id of the action. Accel paths are set to
// <Actions>/group-name/action-name.
// Parameters:
// 	actionGroup	the action group
// 	action	the action to add
// 	accelerator	the accelerator for the action, in
// the format understood by AcceleratorParse, or "" for no accelerator, or
// NULL to use the stock accelerator.
func (a *CActionGroup) AddActionWithAccel(action Action, accelerator string) {}

// Removes an action object from the action group.
// Parameters:
// 	actionGroup	the action group
// 	action	an action
func (a *CActionGroup) RemoveAction(action Action) {}

// This is a convenience function to create a number of actions and add them
// to the action group. The "activate" signals of the actions are connected
// to the callbacks and their accel paths are set to
// <Actions>/group-name/action-name.
// Parameters:
// 	actionGroup	the action group
// 	entries	an array of action descriptions
// 	nEntries	the number of entries
// 	userData	data to pass to the action callbacks
func (a *CActionGroup) AddActions(entries []ActionEntry, nEntries int, userData interface{}) {}

// This variant of AddActions adds a GDestroyNotify
// callback for user_data .
// Parameters:
// 	actionGroup	the action group
// 	entries	an array of action descriptions
// 	nEntries	the number of entries
// 	userData	data to pass to the action callbacks
// 	destroy	destroy notification callback for user_data
//
func (a *CActionGroup) AddActionsFull(entries []ActionEntry, nEntries int, userData interface{}, destroy GDestroyNotify) {
}

// This is a convenience function to create a number of toggle actions and
// add them to the action group. The "activate" signals of the actions are
// connected to the callbacks and their accel paths are set to
// <Actions>/group-name/action-name.
// Parameters:
// 	actionGroup	the action group
// 	entries	an array of toggle action descriptions
// 	nEntries	the number of entries
// 	userData	data to pass to the action callbacks
func (a *CActionGroup) AddToggleActions(entries []ToggleActionEntry, nEntries int, userData interface{}) {
}

// This variant of AddToggleActions adds a
// GDestroyNotify callback for user_data .
// Parameters:
// 	actionGroup	the action group
// 	entries	an array of toggle action descriptions
// 	nEntries	the number of entries
// 	userData	data to pass to the action callbacks
// 	destroy	destroy notification callback for user_data
//
func (a *CActionGroup) AddToggleActionsFull(entries []ToggleActionEntry, nEntries int, userData interface{}, destroy GDestroyNotify) {
}

// This is a convenience routine to create a group of radio actions and add
// them to the action group. The "changed" signal of the first radio action
// is connected to the on_change callback and the accel paths of the actions
// are set to <Actions>/group-name/action-name.
// Parameters:
// 	actionGroup	the action group
// 	entries	an array of radio action descriptions
// 	nEntries	the number of entries
// 	value	the value of the action to activate initially, or -1 if
// no action should be activated
// 	onChange	the callback to connect to the changed signal
// 	userData	data to pass to the action callbacks
func (a *CActionGroup) AddRadioActions(entries []RadioActionEntry, nEntries int, value int, onChange enums.GCallback, userData interface{}) {
}

// This variant of AddRadioActions adds a GDestroyNotify
// callback for user_data .
// Parameters:
// 	actionGroup	the action group
// 	entries	an array of radio action descriptions
// 	nEntries	the number of entries
// 	value	the value of the action to activate initially, or -1 if
// no action should be activated
// 	onChange	the callback to connect to the changed signal
// 	userData	data to pass to the action callbacks
// 	destroy	destroy notification callback for user_data
//
func (a *CActionGroup) AddRadioActionsFull(entries []RadioActionEntry, nEntries int, value int, onChange enums.GCallback, userData interface{}, destroy GDestroyNotify) {
}

// Sets a function to be used for translating the label and tooltip of
// ActionGroupEntrys added by AddActions. If you're
// using gettext, it is enough to set the translation domain with
// SetTranslationDomain.
// Parameters:
// 	func	a TranslateFunc
// 	data	data to be passed to func
// and notify
//
// 	notify	a GDestroyNotify function to be called when action_group
// is
// destroyed and when the translation function is changed again
func (a *CActionGroup) SetTranslateFunc(fn TranslateFunc, data interface{}, notify GDestroyNotify) {}

// Sets the translation domain and uses g_dgettext for translating the
// label and tooltip of ActionEntrys added by
// AddActions. If you're not using gettext for
// localization, see SetTranslateFunc.
// Parameters:
// 	domain	the translation domain to use for g_dgettext calls
func (a *CActionGroup) SetTranslationDomain(domain string) {}

// Translates a string using the specified translate_func. This is mainly
// intended for language bindings.
// Parameters:
// 	string	a string
// Returns:
// 	the translation of string
func (a *CActionGroup) TranslateString(string string) (value string) {
	return ""
}

// A name for the action group.
// Flags: Read / Write / Construct Only
// Default value: NULL
// const PropertyName cdk.Property = "name"

// Whether the action group is enabled.
// Flags: Read / Write
// Default value: TRUE
// const PropertySensitive cdk.Property = "sensitive"

// Whether the action group is visible.
// Flags: Read / Write
// Default value: TRUE
// const PropertyVisible cdk.Property = "visible"

// The ::connect-proxy signal is emitted after connecting a proxy to an
// action in the group. Note that the proxy may have been connected to a
// different action before. This is intended for simple customizations for
// which a custom action class would be too clumsy, e.g. showing tooltips for
// menuitems in the statusbar. UIManager proxies the signal and provides
// global notification just before any action is connected to a proxy, which
// is probably more convenient to use.
// Listener function arguments:
// 	action Action	the action
// 	proxy Widget	the proxy
const SignalConnectProxy cdk.Signal = "connect-proxy"

// The ::disconnect-proxy signal is emitted after disconnecting a proxy from
// an action in the group. UIManager proxies the signal and provides
// global notification just before any action is connected to a proxy, which
// is probably more convenient to use.
// Listener function arguments:
// 	action Action	the action
// 	proxy Widget	the proxy
const SignalDisconnectProxy cdk.Signal = "disconnect-proxy"

// The ::post-activate signal is emitted just after the action in the
// action_group is activated This is intended for UIManager to proxy the
// signal and provide global notification just after any action is activated.
// Listener function arguments:
// 	action Action	the action
const SignalPostActivate cdk.Signal = "post-activate"

// The ::pre-activate signal is emitted just before the action in the
// action_group is activated This is intended for UIManager to proxy the
// signal and provide global notification just before any action is
// activated.
// Listener function arguments:
// 	action Action	the action
const SignalPreActivate cdk.Signal = "pre-activate"

type GDestroyNotify = func(data interface{})

type TranslateFunc = func(messageId string) (translated string)
