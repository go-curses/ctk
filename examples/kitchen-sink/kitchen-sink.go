// +build example_kitchen_sink

package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/go-curses/cdk"
	"github.com/go-curses/cdk/lib/enums"
	cstrings "github.com/go-curses/cdk/lib/strings"
	"github.com/go-curses/cdk/log"
	"github.com/go-curses/ctk"
)

const (
	APP_NAME    = "kitchen-sink"
	APP_USAGE   = "kitchen-sink"
	APP_DESC    = "an example CLI application demonstrating various CTK widgets"
	APP_VERSION = "0.0.1"
	APP_TAG     = "kitchensink"
	APP_TITLE   = "Kitchen Sink"
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
	Debug                     = false
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
}

func main() {
	app := cdk.NewApp(APP_NAME, APP_USAGE, APP_DESC, APP_VERSION, APP_TAG, APP_TITLE, "/dev/tty", setupUi)
	app.AddFlag(&cli.BoolFlag{
		Name:    "Debug",
		Aliases: []string{"d"},
	})
	// app.AddFlag(&cli.StringFlag{
	// 	Name:  "example-flag",
	// 	Value: "testing",
	// })
	// app.AddCommand(&cli.Command{
	// 	Name:  "demo-cmd",
	// 	Usage: "demonstrate custom commands",
	// 	Action: func(c *cli.Context) error {
	// 		log.InfoF("demo-cmd command action")
	// 		return nil
	// 	},
	// })
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

var (
	contentBox  ctk.HBox
	knownPages  []ctk.Alignment
	currentPage int
	pageBox0    ctk.Alignment
	pageBox1    ctk.Alignment
	pageBox2    ctk.Alignment
	actionBox   ctk.HButtonBox
	buttonNext  ctk.Button
	buttonPrev  ctk.Button
	actionNote  ctk.Label
)

func setupUi(d cdk.Display) error {
	if d.App().GetContext().Bool("Debug") {
		log.DebugF("enabling Debug")
		Debug = true
	}
	// note that screen is captured at this time!
	d.CaptureCtrlC()
	w := ctk.NewWindowWithTitle(APP_TITLE)
	d.SetActiveWindow(w)
	w.Show()
	if err := setupDruidUi(d, w); err != nil {
		return err
	}
	if err := setupPage0(d); err != nil {
		return err
	}
	if err := setupPage1(d); err != nil {
		return err
	}
	if err := setupPage2(d); err != nil {
		return err
	}
	switchPage(0)
	return nil
}

func setupDruidUi(d cdk.Display, w ctk.Window) error {
	wVbox := w.GetVBox()
	vbox := ctk.NewVBox(false, 0)
	vbox.SetName("main")
	wVbox.PackStart(vbox, true, true, 0)
	vbox.Show()
	// content area is top row
	contentBox = ctk.NewHBox(true, 0)
	contentBox.SetName("content")
	contentBox.SetBoolProperty(cdk.PropertyDebug, Debug)
	contentBox.Show()
	vbox.PackStart(contentBox, true, true, 0)
	// bottom area is nav buttons, starting with a container for them
	actionBox = ctk.NewHButtonBox(false, 0)
	actionBox.SetName("action")
	actionBox.SetBoolProperty(cdk.PropertyDebug, Debug)
	actionBox.Show()
	vbox.PackEnd(actionBox, false, false, 0)
	actionBox.SetSizeRequest(-1, 3)
	// back button
	buttonPrev = newButton("previous", "_Back", handlePrevious)
	buttonPrev.SetBoolProperty(cdk.PropertyDebug, Debug)
	buttonPrev.Show()
	actionBox.PackStart(buttonPrev, false, false, 0)
	// informational text area
	var err error
	actionNote, err = ctk.NewLabelWithMarkup("Curses<u><i>!</i></u>")
	if err != nil {
		log.ErrorF("failed to set action markup: %v", err)
	}
	actionNote.SetName("note")
	actionNote.SetBoolProperty(cdk.PropertyDebug, Debug)
	actionNote.SetAlignment(0.5, 0.5)
	actionNote.SetJustify(enums.JUSTIFY_RIGHT)
	actionNote.Show()
	actionBox.PackEnd(actionNote, true, true, 0)
	// forward button
	buttonNext = newButton("next", "_Next", handleNext)
	buttonNext.SetBoolProperty(cdk.PropertyDebug, Debug)
	actionBox.PackStart(buttonNext, false, false, 0)
	buttonNext.Show()
	return nil
}

func switchPage(id int) {
	numKnownPages := len(knownPages)
	if numKnownPages > 0 && id < numKnownPages {
		for _, child := range contentBox.GetChildren() {
			contentBox.Remove(child)
		}
		log.InfoF("known page: [%d] %v", id, knownPages[id].ObjectName())
		contentBox.PackStart(knownPages[id], true, true, 0)
		contentBox.ShowAll()
		currentPage = id
	}
	if currentPage == 0 {
		// start
		buttonPrev.SetLabel("Back")
		buttonPrev.Hide()
		if numKnownPages > 1 {
			buttonNext.GrabFocus()
			buttonNext.SetLabel("_Next")
			buttonNext.Show()
		} else {
			buttonNext.Hide()
		}
	} else if currentPage < numKnownPages-1 {
		// middle
		buttonNext.SetLabel("_Next")
		buttonNext.Show()
		buttonPrev.SetLabel("_Back")
		buttonPrev.Show()
	} else {
		// end
		buttonPrev.SetLabel("_Back")
		buttonPrev.Show()
		buttonNext.SetLabel("_Quit")
		buttonNext.Show()
	}
}

func handleNext(data []interface{}, argv ...interface{}) enums.EventFlag {
	log.InfoF("pressed next")
	numKnownPages := len(knownPages)
	if currentPage+1 < numKnownPages {
		switchPage(currentPage + 1)
	} else {
		log.InfoF("end of known pages, quitting")
		cdk.GetDefaultDisplay().RequestQuit()
	}
	return enums.EVENT_STOP
}

func handlePrevious(data []interface{}, argv ...interface{}) enums.EventFlag {
	log.InfoF("pressed previous")
	if currentPage-1 > -1 {
		switchPage(currentPage - 1)
	} else {
		log.InfoF("start of known pages")
	}
	return enums.EVENT_STOP
}

const (
	WelcomeMarkup = "Welcome to the Curses<u><i>!</i></u> kitchen sink application."
	Page1Markup   = "Lorem <i>ipsum</i> dolor sit amet, consectetur adipiscing elit. Vestibulum tincidunt orci a quam dignissim mattis. Nulla volutpat egestas nibh vitae facilisis. Nam dictum risus a nisl suscipit, in luctus felis facilisis. Sed et ante pellentesque, vehicula dui vel, dictum eros. Duis convallis sem vitae tellus feugiat rhoncus. Curabitur risus lectus, elementum id molestie vel, gravida fermentum libero. In aliquet massa eu tellus pulvinar, in scelerisque <i>ipsum</i> ultricies. Quisque elementum nulla vitae condimentum venenatis. Vestibulum vitae lectus sit amet <i>ipsum</i> congue semper ornare tempus magna. Aliquam varius, eros eget ultrices auctor, lacus nibh blandit purus, sed rhoncus erat ex sed enim."
)

func setupPage0(d cdk.Display) error {
	if pageBox0 == nil {
		pageBox0 = ctk.MakeAlignment()
		pageBox0.SetBoolProperty(cdk.PropertyDebug, Debug)
		pageBox0.SetName("pg0")
		pageBox0.Set(0.5, 0.5, 0.0, 0.0)
		pageBox0.Show()
		pageBox0.SetBoolProperty(cdk.PropertyDebug, Debug)
	}
	if pageBox0.GetChild() == nil {
		if welcome, err := ctk.NewLabelWithMarkup(WelcomeMarkup); err != nil {
			return err
		} else {
			welcome.SetName("pg0welcome")
			welcome.SetSizeRequest(20, -1)
			welcome.SetLineWrapMode(enums.WRAP_WORD)
			welcome.SetJustify(enums.JUSTIFY_CENTER)
			welcome.SetAlignment(0.5, 0.5)
			welcome.SetBoolProperty(cdk.PropertyDebug, Debug)
			welcome.Show()
			welcome.SetBoolProperty(cdk.PropertyDebug, Debug)
			pageBox0.Add(welcome)
			knownPages = append(knownPages, pageBox0)
		}
	}
	return nil
}

func setupPage1(d cdk.Display) error {
	if pageBox1 == nil {
		pageBox1 = ctk.MakeAlignment()
		pageBox1.SetBoolProperty(cdk.PropertyDebug, Debug)
		pageBox1.SetName("pg1")
		pageBox1.Set(0.5, 0.5, 0.0, 0.0)
		pageBox1.Show()
	}
	if pageBox1.GetChild() == nil {
		if content, err := ctk.NewLabelWithMarkup(Page1Markup); err != nil {
			return err
		} else {
			content.SetName("pg1content")
			// content.SetSizeRequest(30, -1)
			content.SetLineWrapMode(enums.WRAP_WORD)
			content.SetJustify(enums.JUSTIFY_LEFT)
			content.SetAlignment(0.5, 0.5)
			content.SetBoolProperty(cdk.PropertyDebug, Debug)
			contentBox.Connect(
				ctk.SignalResize,
				fmt.Sprintf("%s.resize", content.ObjectName()),
				func(data []interface{}, argv ...interface{}) enums.EventFlag {
					if len(argv) > 0 {
						if localBox, ok := argv[0].(ctk.HBox); ok {
							alloc := localBox.GetAllocation()
							if alloc.H > 0 && alloc.W > 0 {
								content.SetMaxWidthChars(alloc.W)
								content.LogInfo("updating max chars")
							}
						}
					}
					return enums.EVENT_PASS
				},
			)
			content.Show()
			pageBox1.Add(content)
			knownPages = append(knownPages, pageBox1)
		}
	}
	return nil
}

func setupPage2(d cdk.Display) error {
	if pageBox2 == nil {
		pageBox2 = ctk.MakeAlignment()
		pageBox2.SetBoolProperty(cdk.PropertyDebug, Debug)
		pageBox2.SetName("pg1")
		pageBox2.Set(0.5, 0.5, 0.0, 0.0)
		pageBox2.Show()
	}
	if pageBox2.GetChild() == nil {
		if content, err := ctk.NewLabelWithMarkup(Page1Markup); err != nil {
			return err
		} else {
			scroll := ctk.NewScrolledViewport()
			scroll.SetPolicy(ctk.PolicyAutomatic, ctk.PolicyAutomatic)
			scroll.SetSizeRequest(40, 10)
			content.SetSizeRequest(50, -1)
			content.SetName("pg2content")
			content.SetLineWrapMode(enums.WRAP_WORD)
			content.SetJustify(enums.JUSTIFY_LEFT)
			content.SetAlignment(0.5, 0.5)
			content.SetBoolProperty(cdk.PropertyDebug, Debug)
			// scroll.Connect(
			// 	ctk.SignalResize,
			// 	cdk.Signal(fmt.Sprintf("%s.resize", content.ObjectName())),
			// 	func(data []interface{}, argv ...interface{}) enums.EventFlag {
			// 		if len(argv) > 0 {
			// 			if localBox, ok := argv[0].(ctk.ScrolledViewport); ok {
			// 				alloc := localBox.GetSizeRequest()
			// 				if alloc.H > 0 && alloc.W > 0 {
			// 					content.SetMaxWidthChars(alloc.W)
			// 					content.LogInfo("updating max chars")
			// 				}
			// 			}
			// 		}
			// 		// size := scroll.GetAllocation()
			// 		// content.SetMaxWidthChars(size.W)
			// 		return enums.EVENT_PASS
			// 	},
			// )
			content.Show()

			scroll.Add(content)
			scroll.Show()
			pageBox2.Add(scroll)
			knownPages = append(knownPages, pageBox2)
		}
	}
	return nil
}

// var (
// IPSUM_LONG_PLAIN = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vestibulum tincidunt orci a quam dignissim mattis. Nulla volutpat egestas nibh vitae facilisis. Nam dictum risus a nisl suscipit, in luctus felis facilisis. Sed et ante pellentesque, vehicula dui vel, dictum eros. Duis convallis sem vitae tellus feugiat rhoncus. Curabitur risus lectus, elementum id molestie vel, gravida fermentum libero. In aliquet massa eu tellus pulvinar, in scelerisque ipsum ultricies. Quisque elementum nulla vitae condimentum venenatis. Vestibulum vitae lectus sit amet ipsum congue semper ornare tempus magna. Aliquam varius, eros eget ultrices auctor, lacus nibh blandit purus, sed rhoncus erat ex sed enim."
// IPSUM_LONG_MARKUP = "Lorem <i>ipsum</i> dolor sit amet, consectetur adipiscing elit. Vestibulum tincidunt orci a quam dignissim mattis. Nulla volutpat egestas nibh vitae facilisis. Nam dictum risus a nisl suscipit, in luctus felis facilisis. Sed et ante pellentesque, vehicula dui vel, dictum eros. Duis convallis sem vitae tellus feugiat rhoncus. Curabitur risus lectus, elementum id molestie vel, gravida fermentum libero. In aliquet massa eu tellus pulvinar, in scelerisque <i>ipsum</i> ultricies. Quisque elementum nulla vitae condimentum venenatis. Vestibulum vitae lectus sit amet <i>ipsum</i> congue semper ornare tempus magna. Aliquam varius, eros eget ultrices auctor, lacus nibh blandit purus, sed rhoncus erat ex sed enim."
// )

func newButton(name string, label string, fn cdk.SignalListenerFn) ctk.Button {
	b := ctk.NewButtonWithLabel("")
	if child := b.GetChild(); child != nil {
		if l, ok := child.(ctk.Label); ok {
			l.SetMarkup(label)
			l.SetEllipsize(true)
		}
	}
	b.SetName(name)
	b.SetUseUnderline(true)
	b.SetSensitive(true)
	if Debug {
		b.SetBoolProperty(cdk.PropertyDebug, true)
	}
	b.Connect(
		ctk.SignalActivate,
		fmt.Sprintf("%s.activate", name),
		fn,
	)
	return b
}
