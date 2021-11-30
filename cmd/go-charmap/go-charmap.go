package main

import (
	"fmt"
	"os"
	"strconv"
	"unicode"
	"unicode/utf8"

	"github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/cdk/lib/paint"
	"github.com/go-curses/cdk/lib/ptypes"
	"golang.org/x/text/unicode/runenames"

	"github.com/go-curses/cdk"
	"github.com/go-curses/cdk/env"
	cstrings "github.com/go-curses/cdk/lib/strings"
	"github.com/go-curses/cdk/log"
	"github.com/go-curses/ctk"
)

const (
	AppName    = "go-charmap"
	AppUsage   = "View the details of a character"
	AppDesc    = "Get informational details for a specific character given in integer form."
	AppVersion = "0.0.1"
	AppTag     = "charmap"
	AppTitle   = "Character Map"
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
}

func main() {
	app := ctk.NewApplication(
		AppName,
		AppUsage,
		AppDesc,
		AppVersion,
		AppTag,
		AppTitle,
		"/dev/tty",
	)
	app.Connect(cdk.SignalStartup, "go-charmap-startup-handler", func(_ []interface{}, argv ...interface{}) enums.EventFlag {
		if _, d, _, _, _, ok := ctk.ArgvApplicationSignalStartup(argv...); ok {
			if err := setup(d); err != nil {
				app.LogErr(err)
				return enums.EVENT_STOP
			}
			app.NotifyStartupComplete()
			return enums.EVENT_PASS
		}
		return enums.EVENT_STOP
	})
	// app.CLI().Commands = nil
	app.CLI().UsageText = "go-charmap [integer]"
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func setup(d cdk.Display) error {
	theme := paint.DefaultColorTheme
	theme.Content.Normal = theme.Content.Selected.Dim(false)
	theme.Border.Normal = theme.Border.Selected.Dim(false)
	d.CaptureCtrlC()
	w := ctk.NewWindowWithTitle(AppTitle)
	w.SetTheme(theme)
	w.Connect(ctk.SignalEventKey, "escape-quit", func(data []interface{}, argv ...interface{}) enums.EventFlag {
		if evt, ok := argv[1].(cdk.Event); ok {
			switch e := evt.(type) {
			case *cdk.EventKey:
				if e.Key() == cdk.KeyEscape {
					w.LogInfo("window caught escape key, quitting now")
					d.RequestQuit()
				}
			}
			return enums.EVENT_STOP
		}
		return enums.EVENT_PASS
	})
	align := ctk.NewAlignment(0.5, 0.5, 0.0, 0.0)
	align.SetTheme(theme)
	frame := ctk.NewFrame("Character Details")
	frame.SetTheme(theme)
	label := ctk.NewLabel("")
	label.SetJustify(enums.JUSTIFY_LEFT)
	label.SetAlignment(0.5, 0.5)
	label.SetTheme(theme)
	label.SetText("loading...")
	label.SetPadding(1, 1)
	_ = label.SetBoolProperty(cdk.PropertyDebug, false)
	// frame.SetSizeRequest(41, 8)
	// label.SetSizeRequest(30, 6)
	_ = frame.SetBoolProperty(cdk.PropertyDebug, false)
	label.Show()
	frame.Add(label)
	frame.Show()
	align.Add(frame)
	align.Show()
	w.GetVBox().PackStart(align, true, true, 0)
	ctx := d.App().GetContext()
	args := ctx.Args().Slice()
	var message string
	if len(args) > 0 {
		if num, err := strconv.Atoi(args[0]); err != nil {
			message = fmt.Sprintf("invalid argument: %v", args[1])
		} else {
			r, w := utf8.DecodeRune([]byte(string(rune(num))))
			name := runenames.Name(r)
			message += name + "\n"
			message += "Unicode: " + fmt.Sprintf("%U (%d)", r, num) + "\n"
			message += "Entity: &#" + fmt.Sprintf("%d", num) + ";\n"
			message += "Print: "
			if unicode.IsGraphic(r) {
				message += fmt.Sprintf("%c", r)
				if w > 1 {
					for i := 1; i < w; i++ {
						message += "_"
					}
				}
			} else {
				message += fmt.Sprintf("%x", r)
			}
			message += "\n"
			message += "Width: " + fmt.Sprintf("%d", w) + "\n"
		}
	} else {
		message = fmt.Sprintf(
			"Character Set: %s\nDisplay Size: %s\nDisplay Colors: %d\nTerminal: %s",
			d.Screen().CharacterSet(),
			ptypes.MakeRectangle(d.Screen().Size()),
			d.Screen().Colors(),
			env.Get("TERM", "(unset)"),
		)
	}
	label.SetText(message)
	w.ShowAll()
	d.SetActiveWindow(w)
	d.RequestDraw()
	d.RequestShow()
	return nil
}
