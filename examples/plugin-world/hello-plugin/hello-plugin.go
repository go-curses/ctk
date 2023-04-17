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

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/go-curses/cdk"
	cenums "github.com/go-curses/cdk/lib/enums"
	cstrings "github.com/go-curses/cdk/lib/strings"
	"github.com/go-curses/cdk/lib/sync"
	"github.com/go-curses/cdk/log"
	"github.com/go-curses/ctk"
)

var (
	CtkName        = "hello-plugin"
	CtkUsage       = "demonstrates a CTK application plugin"
	CtkDescription = "This program can be built and loaded as Go plugin using ctk.NewApplicationFromPlugin, or built as a stand-alone binary."
	CtkVersion     = "0.0.1"
	CtkTag         = "hello.plugin"
	CtkTitle       = "Hello Plugin"
	CtkTtyPath     = "/dev/tty"
)

// Build Configuration Flags
// setting these will enable command line flags and their corresponding features
// use `go build -v -ldflags="-X 'main.IncludeLogFullPaths=false'"`
var (
	IncludeProfiling          = "false"
	IncludeLogFile            = "false"
	IncludeLogFormat          = "false"
	IncludeLogFullPaths       = "false"
	IncludeLogLevel           = "false"
	IncludeLogLevels          = "false"
	IncludeLogTimestamps      = "false"
	IncludeLogTimestampFormat = "false"
	IncludeLogOutput          = "false"
)

func init() {
	cdk.Build.Profiling = cstrings.IsTrue(IncludeProfiling)
	cdk.Build.LogFile = cstrings.IsTrue(IncludeLogFile)
	cdk.Build.LogFormat = cstrings.IsTrue(IncludeLogFormat)
	cdk.Build.LogFullPaths = cstrings.IsTrue(IncludeLogFullPaths)
	cdk.Build.LogLevel = cstrings.IsTrue(IncludeLogLevel)
	cdk.Build.LogLevels = cstrings.IsTrue(IncludeLogLevels)
	cdk.Build.LogTimestamps = cstrings.IsTrue(IncludeLogTimestamps)
	cdk.Build.LogTimestampFormat = cstrings.IsTrue(IncludeLogTimestampFormat)
	cdk.Build.LogOutput = cstrings.IsTrue(IncludeLogOutput)
	cdk.Init()
}

func main() {
	app := ctk.NewApplication(
		CtkName,
		CtkUsage,
		CtkDescription,
		CtkVersion,
		CtkTag,
		CtkTitle,
		CtkTtyPath,
	)
	app.Connect(cdk.SignalStartup, "main-startup-handler", ctk.WithArgvApplicationSignalStartup(CtkStartup))
	app.Connect(cdk.SignalShutdown, "main-shutdown-handler", ctk.WithArgvNoneWithFlagsSignal(CtkShutdown))
	CtkInit(app)
	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func CtkInit(app ctk.Application) {
	app.LogInfo("application initialized")
}

func CtkStartup(app ctk.Application, d cdk.Display, ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup) cenums.EventFlag {
	// tell the display to listen for CTRL+C and interrupt gracefully
	d.CaptureCtrlC()
	// create a new window, give it a human-readable title
	w := ctk.NewWindowWithTitle("Hello World")
	// get the vertical box for the content area of the window
	vbox := w.GetVBox()
	// here is where we add other widgets and such to implement the
	// desired user interface elements, in this case we want a nice
	// button in the middle of the window. One way to do this is to
	// use an Alignment widget to place the button neatly for us.
	align := ctk.MakeAlignment()
	// the alignment scales are from 0 (left) to 1 (right) with the 0.5
	// being centered
	align.Set(0.5, 0.5, 0.0, 0.0)
	// a nice button for us to press
	button := ctk.NewButtonWithLabel("_Curses<u><i>!</i></u>")
	button.SetUseMarkup(true)    // enable markup in the label
	button.SetUseUnderline(true) // enable mnemonics
	button.SetSizeRequest(11, 3) // request a certain size
	// make the button quit the application when activated by connecting
	// a handler to the button's activate signal
	button.Connect(
		ctk.SignalActivate,
		"hello-button-handle",
		func(data []interface{}, argv ...interface{}) cenums.EventFlag {
			d.RequestQuit() // ask the display to exit nicely
			return cenums.EVENT_STOP
		},
	)
	align.Add(button) // add the button to the alignment
	// finally adding the alignment to the window's content area by
	// packing them into the window's vertical box
	vbox.PackStart(align, true /*expand*/, true /*fill*/, 0 /*padding*/)
	// tell CTK that the window and its contents are to be drawn upon
	// the terminal display, this effectively calls Show() on the vbox,
	// alignment and button
	w.ShowAll()
	// now that the button is visible, grab the focus
	button.GrabFocus()
	// notify that startup has completed
	app.NotifyStartupComplete()
	return cenums.EVENT_PASS
}

func CtkShutdown() cenums.EventFlag {
	// Note that the Display and other CTK things are no longer functioning at
	// this point.
	fmt.Println("Hello World says Goodbye!")
	// Logging however still works.
	log.InfoF("Hello World logging goodbye!")
	return cenums.EVENT_PASS
}