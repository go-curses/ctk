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
	"github.com/go-curses/cdk/lib/enums"
)

const TypeApplication cdk.CTypeTag = "ctk-application"

func init() {
	_ = cdk.TypesManager.AddType(TypeApplication, nil)
}

// Application Hierarchy:
//	Object
//	  +- Application
//
// An Application is the CTK replacement for cdk.Application.
type Application interface {
	cdk.Application

	AccelMap() (accelMap AccelMap)
	AccelGroup() (accelGroup AccelGroup)
}

var _ Application = (*CApplication)(nil)

type CApplication struct {
	cdk.CApplication

	accelGroup AccelGroup
	accelMap   AccelMap
	windows    []Window
}

func NewApplication(name, usage, description, version, tag, title, ttyPath string) (app Application) {
	app = &CApplication{}
	app.Init()
	app.Reconfigure(name, usage, description, version, tag, title, ttyPath)
	return
}

func (app *CApplication) Init() (already bool) {
	if app.InitTypeItem(TypeApplication, app) {
		return true
	}
	app.CApplication.Init()
	app.accelMap = &CAccelMap{}
	app.accelMap.Init()
	app.accelGroup = NewAccelGroup()
	app.Connect(cdk.SignalSetupDisplay, ApplicationSetupDisplayHandle, func(_ []interface{}, argv ...interface{}) enums.EventFlag {
		if display, ok := argv[0].(cdk.Display); ok {
			if !display.Handled(cdk.SignalFocusedWindow, ApplicationFocusedWindowHandle) {
				display.Connect(cdk.SignalFocusedWindow, ApplicationFocusedWindowHandle, app.displayWindowsChanged)
			}
			if !display.Handled(cdk.SignalMappedWindow, ApplicationMappedWindowHandle) {
				display.Connect(cdk.SignalMappedWindow, ApplicationMappedWindowHandle, app.displayWindowsChanged)
			}
			if !display.Handled(cdk.SignalUnmappedWindow, ApplicationUnmappedWindowHandle) {
				display.Connect(cdk.SignalUnmappedWindow, ApplicationUnmappedWindowHandle, app.displayWindowsChanged)
			}
		}
		return enums.EVENT_PASS
	})
	return false
}

func (app *CApplication) AccelMap() (accelMap AccelMap) {
	app.RLock()
	defer app.RUnlock()
	accelMap = app.accelMap
	return
}

func (app *CApplication) AccelGroup() (accelGroup AccelGroup) {
	app.RLock()
	defer app.RUnlock()
	accelGroup = app.accelGroup
	return
}

func (app *CApplication) GetWindows() (windows []Window) {
	app.RLock()
	for _, w := range app.windows {
		windows = append(windows, w)
	}
	app.RUnlock()
	return
}

func (app *CApplication) displayWindowsChanged(_ []interface{}, argv ...interface{}) enums.EventFlag {
	if display := app.Display(); display != nil {
		app.Lock()
		windows := display.GetWindows()
		ctkWindows := []Window{}
		for _, w := range windows {
			if ww, ok := w.Self().(Window); ok {
				ctkWindows = append(ctkWindows, ww)
			}
		}
		app.windows = ctkWindows
		app.Unlock()
	}
	return enums.EVENT_PASS
}

const ApplicationSetupDisplayHandle = "application-setup-display-handler"
const ApplicationFocusedWindowHandle = "application-focused-window-handler"
const ApplicationMappedWindowHandle = "application-mapped-window-handler"
const ApplicationUnmappedWindowHandle = "application-unmapped-window-handler"