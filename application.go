// Copyright 2021  The CDK Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use file except in compliance with the License.
// You may obtain a copy of the license at
//
//    http://www.apache.org/licenses/LICENSE-2.0
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

	Init() (already bool)
	AccelMap() (accelMap AccelMap)
	AccelGroup() (accelGroup AccelGroup)
}

type CApplication struct {
	cdk.CApplication

	accelGroup AccelGroup
	accelMap   AccelMap
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
