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
	"github.com/go-curses/cdk"
)

const TypeToggleAction cdk.CTypeTag = "ctk-toggle-action"

func init() {
	_ = cdk.TypesManager.AddType(TypeToggleAction, func() interface{} { return MakeToggleAction() })
}

// ToggleAction Hierarchy:
//	Object
//	  +- Action
//	    +- ToggleAction
//	      +- RadioAction
type ToggleAction interface {
	Action

	Toggled()
	SetActive(isActive bool)
	GetActive() (value bool)
	SetDrawAsRadio(drawAsRadio bool)
	GetDrawAsRadio() (value bool)
}

var _ ToggleAction = (*CToggleAction)(nil)

// The CToggleAction structure implements the ToggleAction interface and is
// exported to facilitate type embedding with custom implementations. No member
// variables are exported as the interface methods are the only intended means
// of interacting with ToggleAction objects.
type CToggleAction struct {
	CAction
}

// MakeToggleAction is used by the Buildable system to construct a new ToggleAction.
func MakeToggleAction() ToggleAction {
	return NewToggleAction("", "", "", "")
}

// NewToggleAction is the constructor for new ToggleAction instances.
func NewToggleAction(name string, label string, tooltip string, stockId string) (value ToggleAction) {
	t := new(CToggleAction)
	t.Init()
	return t
}

// Init initializes an ToggleAction object. This must be called at least once to
// set up the necessary defaults and allocate any memory structures. Calling
// this more than once is safe though unnecessary. Only the first call will
// result in any effect upon the ToggleAction instance. Init is used in the
// NewToggleAction constructor and only necessary when implementing a derivative
// ToggleAction type.
func (t *CToggleAction) Init() (already bool) {
	if t.InitTypeItem(TypeToggleAction, t) {
		return true
	}
	t.CAction.Init()
	_ = t.InstallProperty(PropertyActive, cdk.BoolProperty, true, false)
	_ = t.InstallProperty(PropertyDrawAsRadio, cdk.BoolProperty, true, false)
	return false
}

// Toggled emits the "toggled" signal on the toggle action.
//
// Parameters:
// 	action	the action object
func (t *CToggleAction) Toggled() {
	t.Emit(SignalToggled, t)
}

// SetActive updates the checked state on the toggle action.
//
// Parameters:
// 	isActive	whether the action should be checked or not
func (t *CToggleAction) SetActive(isActive bool) {
	if err := t.SetBoolProperty(PropertyActive, isActive); err != nil {
		t.LogErr(err)
	}
}

// GetActive returns the checked state of the toggle action.
//
// Parameters:
// 	action	the action object
func (t *CToggleAction) GetActive() (value bool) {
	var err error
	if value, err = t.GetBoolProperty(PropertyActive); err != nil {
		t.LogErr(err)
	}
	return
}

// SetDrawAsRadio updates whether the action should have proxies like a radio
// action.
//
// Parameters:
// 	drawAsRadio	whether the action should have proxies like a radio action
func (t *CToggleAction) SetDrawAsRadio(drawAsRadio bool) {
	if err := t.SetBoolProperty(PropertyDrawAsRadio, drawAsRadio); err != nil {
		t.LogErr(err)
	}
}

// GetDrawAsRadio returns whether the action should have proxies like a radio
// action.
//
// Parameters:
// 	action	the action object
func (t *CToggleAction) GetDrawAsRadio() (value bool) {
	var err error
	if value, err = t.GetBoolProperty(PropertyDrawAsRadio); err != nil {
		t.LogErr(err)
	}
	return
}

// If the toggle action should be active in or not.
// Flags: Read / Write
// Default value: FALSE
const PropertyActive cdk.Property = "active"

// Whether the proxies for this action look like radio action proxies. This
// is an appearance property and thus only applies if
// “use-action-appearance” is TRUE.
// Flags: Read / Write
// Default value: FALSE
const PropertyDrawAsRadio cdk.Property = "draw-as-radio"

const SignalToggled cdk.Signal = "toggled"