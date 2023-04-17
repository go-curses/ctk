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
	"context"

	"github.com/go-curses/cdk"
	"github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/cdk/lib/sync"
	"github.com/go-curses/cdk/log"
)

type AppFn func(app Application)

func WithApp(initFn cdk.SignalListenerFn, action AppFn) func() {
	return func() {
		app := NewApplication(
			"AppName", "AppUsage",
			"AppDesc", "v0.0.0",
			"app-tag", "AppTitle",
			cdk.OffscreenTtyPath,
		)
		app.Connect(cdk.SignalStartup, "testing-withapp-init-fn-handler", initFn)
		defer func() {
			if app != nil {
				app.Destroy()
			}
			app = nil
		}()
		app.SetupDisplay()
		ctx, cancel := context.WithCancel(context.Background())
		if f := app.Emit(cdk.SignalStartup, app, app.Display(), ctx, cancel, &sync.WaitGroup{}); f == enums.EVENT_STOP {
			log.ErrorF("WithApp startup listeners requested EVENT_STOP")
		} else {
			action(app)
		}
	}
}

func TestingWithCtkWindow(_ []interface{}, argv ...interface{}) enums.EventFlag {
	if _, d, _, _, _, ok := ArgvApplicationSignalStartup(argv...); ok {
		w := NewWindowWithTitle(d.GetTitle())
		d.FocusWindow(w)
	}
	return enums.EVENT_PASS
}