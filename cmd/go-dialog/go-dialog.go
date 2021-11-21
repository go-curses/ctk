package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/cdk/lib/ptypes"
	"github.com/urfave/cli/v2"

	"github.com/go-curses/cdk"
	cstrings "github.com/go-curses/cdk/lib/strings"
	"github.com/go-curses/cdk/log"
	"github.com/go-curses/ctk"
)

const (
	AppName    = "go-dialog"
	AppUsage   = "go-dialog [options] command [command options]"
	AppDesc    = "display dialog boxes from shell scripts"
	AppVersion = "0.0.1"
	AppTag     = "go-dialog"
	AppTitle   = "go-dialog"
)

//go:embed dialog-msgbox.glade
var gladeMsgBox string

//go:embed dialog-yesno.glade
var gladeYesNo string

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
	app := cdk.NewApp(
		AppName, AppUsage,
		AppDesc, AppVersion,
		AppTag, AppTitle,
		"/dev/tty",
		setupUserInterface,
	)
	app.AddFlag(&cli.StringFlag{
		Name:  "back-title",
		Usage: "specify the window title text",
		Value: "",
	})
	app.AddFlag(&cli.StringFlag{
		Name:  "title",
		Usage: "specify the dialog title text",
		Value: "",
	})
	app.AddFlag(&cli.BoolFlag{
		Name:  "print-maxsize",
		Usage: "print the width and height on stdout and exit",
		Value: false,
	})
	app.AddCommand(&cli.Command{
		Name:      "msgbox",
		Usage:     "display a message with an OK button, each string following msgbox is a new line and concatenated into the message",
		ArgsUsage: "[message lines]",
		Action:    app.CliActionFn,
	})
	app.AddCommand(&cli.Command{
		Name:      "yesno",
		Usage:     "display a yes/no prompt with a message (see msgbox)",
		ArgsUsage: "[message lines]",
		Action:    app.CliActionFn,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "default",
				Usage:       "specify which button is focused initially",
				Value:       "yes",
				DefaultText: "yes",
			},
		},
	})
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func setupUserInterface(dm cdk.Display) error {
	ctx := dm.App().GetContext()
	if ctx.Bool("print-maxsize") {
		if display := dm.Screen(); display != nil {
			w, h := display.Size()
			dm.AddQuitHandler("print-maxsize", func() {
				fmt.Printf("%v %v\n", w, h)
			})
			dm.RequestQuit()
			return nil
		}
	}
	dm.LogInfo("setting up user interface")
	dm.CaptureCtrlC()

	builder := ctk.NewBuilder()
	var proceed bool
	switch ctx.Command.Name {
	case "msgbox":
		if err := setupUiMsgbox(ctx, builder, dm); err != nil {
			return err
		}
		proceed = true
	case "yesno":
		if err := setupUiYesNo(ctx, builder, dm); err != nil {
			return err
		}
		proceed = true
	case "":
		dm.AddQuitHandler("see-help", func() {
			fmt.Printf("see: %v --help\n", dm.App().Name())
		})
		dm.RequestQuit()
		return nil
	default:
		return fmt.Errorf("invalid command: %v", ctx.Command.Name)
	}
	if proceed {
		if err := startupUiDialog(ctx, builder, dm); err != nil {
			return err
		}
		dm.LogInfo("user interface set up complete")
		return nil
	}
	return fmt.Errorf("error intializing user interface")
}

func startupUiDialog(ctx *cli.Context, builder ctk.Builder, dm cdk.Display) error {
	backTitle := ctx.String("back-title")
	title := ctx.String("title")
	window := getWindow(builder)
	dialog := getDialog(builder)
	if window != nil {
		window.Show()
		window.SetTitle(backTitle)
		dm.SetActiveWindow(window)
		if dialog != nil {
			if display := dm.Screen(); display != nil {
				dw, dh := display.Size()
				if dw > 22 && dh > 12 {
					sr := ptypes.NewRectangle(dw/3, dh/3)
					sr.Clamp(20, 10, dw, dh)
					dialog.SetSizeRequest(sr.W, sr.H)
				}
				dialog.SetTransientFor(window)
				dialog.SetTitle(title)
			}
			dialog.Show()
			dialog.LogInfo("starting Run()")
			defBtn := ctx.String("default")
			switch strings.ToLower(defBtn) {
			case "no", "ctk-no":
				if no := builder.GetWidget("yesno-no"); no != nil {
					if nw, ok := no.(ctk.Widget); ok {
						nw.GrabFocus()
					}
				}
			case "yes", "ctk-yes":
				fallthrough
			default:
				if yes := builder.GetWidget("yesno-yes"); yes != nil {
					if yw, ok := yes.(ctk.Widget); ok {
						yw.GrabFocus()
					}
				}
			}
			dm.RequestDraw()
			dm.RequestShow()
			response := dialog.Run()
			go func() {
				select {
				case r := <-response:
					dialog.Destroy()
					_ = dialog.DestroyObject()
					switch ctx.Command.Name {
					case "yesno":
						dm.AddQuitHandler("dialog-response", func() {
							fmt.Printf("%v\n", r)
						})
					}
					dm.RequestQuit()
				}
			}()
		} else {
			builder.LogError("missing main-dialog")
		}
	} else {
		builder.LogError("missing main-window")
	}
	return nil
}

func setupUiMsgbox(ctx *cli.Context, builder ctk.Builder, dm cdk.Display) error {
	builder.AddNamedSignalHandler("msgbox-ok", func(data []interface{}, argv ...interface{}) enums.EventFlag {
		if dialog := getDialog(builder); dialog != nil {
			dialog.Response(ctk.ResponseOk)
		} else {
			builder.LogError("msgbox-ok missing main-dialog")
		}
		return enums.EVENT_STOP
	})
	if tmpl, err := template.New("msgbox").Parse(gladeMsgBox); err != nil || tmpl == nil {
		dm.LogErr(err)
	} else {
		content := ""
		if ctx.Args().Len() < 1 {
			return fmt.Errorf("msgbox missing message to display")
		}
		for i := 0; i < ctx.Args().Len(); i++ {
			if content != "" {
				content += "\n"
			}
			content += fmt.Sprintf("%v", ctx.Args().Get(i))
		}

		buff := new(bytes.Buffer)
		data := struct {
			Message string
		}{
			Message: content,
		}
		if err := tmpl.Execute(buff, data); err == nil {
			xml := string(buff.Bytes())
			var err error
			if _, err = builder.LoadFromString(xml); err != nil {
				return err
			}
		}
	}
	return nil
}

func setupUiYesNo(ctx *cli.Context, builder ctk.Builder, dm cdk.Display) error {
	builder.AddNamedSignalHandler("yesno-yes", func(data []interface{}, argv ...interface{}) enums.EventFlag {
		if dialog := getDialog(builder); dialog != nil {
			dialog.Response(ctk.ResponseYes)
		} else {
			builder.LogError("yesno-yes missing main-dialog")
		}
		return enums.EVENT_STOP
	})
	builder.AddNamedSignalHandler("yesno-no", func(data []interface{}, argv ...interface{}) enums.EventFlag {
		if dialog := getDialog(builder); dialog != nil {
			dialog.Response(ctk.ResponseNo)
		} else {
			builder.LogError("yesno-no missing main-dialog")
		}
		return enums.EVENT_STOP
	})
	if tmpl, err := template.New("yesno").Parse(gladeYesNo); err != nil || tmpl == nil {
		dm.LogErr(err)
	} else {
		content := ""
		if ctx.Args().Len() < 1 {
			return fmt.Errorf("yesno missing message to display")
		}
		for i := 0; i < ctx.Args().Len(); i++ {
			if content != "" {
				content += "\n"
			}
			content += fmt.Sprintf("%v", ctx.Args().Get(i))
		}
		buff := new(bytes.Buffer)
		data := struct {
			Message string
		}{
			Message: content,
		}
		if err := tmpl.Execute(buff, data); err == nil {
			xml := string(buff.Bytes())
			var err error
			if _, err = builder.LoadFromString(xml); err != nil {
				return err
			}
		}
	}
	return nil
}

func getWindow(builder ctk.Builder) (window ctk.Window) {
	if mw := builder.GetWidget("main-window"); mw != nil {
		var ok bool
		if window, ok = mw.(ctk.Window); !ok {
			builder.LogError("main-window widget is not of ctk.Window type: %v (%T)", mw, mw)
		}
	} else {
		builder.LogError("missing main-window widget")
	}
	return
}

func getDialog(builder ctk.Builder) (dialog ctk.Dialog) {
	if md := builder.GetWidget("main-dialog"); md != nil {
		var ok bool
		if dialog, ok = md.(ctk.Dialog); !ok {
			builder.LogError("main-dialog widget is not of ctk.Dialog type: %v (%T)", md, md)
		}
	} else {
		builder.LogError("missing main-dialog widget")
	}
	return
}
