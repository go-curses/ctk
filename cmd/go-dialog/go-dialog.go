package main

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/urfave/cli/v2"

	"github.com/go-curses/cdk"
	cenums "github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/cdk/lib/ptypes"
	cstrings "github.com/go-curses/cdk/lib/strings"
	"github.com/go-curses/cdk/lib/sync"
	"github.com/go-curses/cdk/log"
	"github.com/go-curses/ctk"
	"github.com/go-curses/ctk/lib/enums"
)

const (
	AppName    = "go-dialog"
	AppUsage   = "display dialog boxes from shell scripts"
	AppDesc    = "go-dialog is another version of (c)dialog, just written in Go and using the Curses Tool Kit"
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
	app := ctk.NewApplication(
		AppName, AppUsage,
		AppDesc, AppVersion,
		AppTag, AppTitle,
		"/dev/tty",
	)
	app.Connect(
		cdk.SignalStartup, "go-dialog-startup-handler",
		ctk.WithArgvApplicationSignalStartup(
			func(app ctk.Application, display cdk.Display, ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup) cenums.EventFlag {
				if err := setupUserInterface(app, display); err != nil {
					app.LogErr(err)
					return cenums.EVENT_STOP
				}
				return cenums.EVENT_PASS
			},
		),
	)
	app.AddFlag(&cli.StringFlag{
		Name:  "back-title",
		Usage: "specify the window title text",
	})
	app.AddFlag(&cli.StringFlag{
		Name:  "title",
		Usage: "specify the dialog title text",
	})
	app.AddFlag(&cli.BoolFlag{
		Name:        "print-maxsize",
		Usage:       "print the width and height on stdout and exit",
		DefaultText: "",
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

func setupUserInterface(app ctk.Application, d cdk.Display) error {
	ctx := app.GetContext()
	if ctx.Bool("print-maxsize") {
		if display := d.Screen(); display != nil {
			w, h := display.Size()
			app.Connect(cdk.SignalShutdown, "print-maxsize", func(_ []interface{}, _ ...interface{}) cenums.EventFlag {
				fmt.Printf("%v %v\n", w, h)
				return cenums.EVENT_PASS
			})
			d.RequestQuit()
			return nil
		}
	}
	d.LogInfo("setting up user interface")
	d.CaptureCtrlC()
	builder := ctk.NewBuilder()
	var proceed bool
	switch ctx.Command.Name {
	case "msgbox":
		if err := setupUiMsgbox(ctx, builder, app, d); err != nil {
			return err
		}
		proceed = true
	case "yesno":
		if err := setupUiYesNo(ctx, builder, app, d); err != nil {
			return err
		}
		proceed = true
	case "":
		app.Connect(cdk.SignalShutdown, "see-help", func(_ []interface{}, _ ...interface{}) cenums.EventFlag {
			fmt.Printf("see: %v --help\n", app.Name())
			return cenums.EVENT_PASS
		})
		d.RequestQuit()
		return nil
	default:
		return fmt.Errorf("invalid command: %v", ctx.Command.Name)
	}
	if proceed {
		if err := startupUiDialog(ctx, builder, app, d); err != nil {
			return err
		}
		d.LogInfo("user interface set up complete")
		return nil
	}
	return fmt.Errorf("error intializing user interface")
}

func startupUiDialog(ctx *cli.Context, builder ctk.Builder, app ctk.Application, display cdk.Display) error {
	app.NotifyStartupComplete()
	backTitle := ctx.String("back-title")
	title := ctx.String("title")
	window := getWindow(builder)
	dialog := getDialog(builder)
	if window != nil {
		window.SetTitle(backTitle)
		window.Show()
		if dialog != nil {
			dialog.Show()
			dw, dh := display.Screen().Size()
			if dw > 22 && dh > 12 {
				sr := ptypes.NewRectangle(dw/3, dh/3)
				sr.Clamp(20, 10, dw, dh)
				dialog.SetSizeRequest(sr.W, sr.H)
			}
			dialog.SetTransientFor(window)
			dialog.SetParent(window)
			if err := dialog.ImportStylesFromString(window.ExportStylesToString()); err != nil {
				dialog.LogErr(err)
			}
			dialog.SetTitle(title)
			dialog.LogInfo("starting Run()")
			defBtn := ctx.String("default")
			switch strings.ToLower(defBtn) {
			case "no", "ctk-no":
				if no := builder.GetWidget("main-yesno-no"); no != nil {
					if nw, ok := no.(ctk.Sensitive); ok {
						nw.GrabFocus()
					}
				}
			default:
				if yes := builder.GetWidget("main-yesno-yes"); yes != nil {
					if yw, ok := yes.(ctk.Sensitive); ok {
						yw.GrabFocus()
					}
				} else if okay := builder.GetWidget("main-msgbox-ok"); okay != nil {
					if yw, ok := okay.(ctk.Sensitive); ok {
						yw.GrabFocus()
					}
				}
			}
			display.RequestDraw()
			display.RequestShow()
			response := dialog.Run()
			cdk.Go(func() {
				select {
				case r := <-response:
					dialog.Destroy()
					_ = dialog.DestroyObject()
					switch ctx.Command.Name {
					case "yesno":
						app.Connect(cdk.SignalShutdown, "dialog-response", func(_ []interface{}, _ ...interface{}) cenums.EventFlag {
							fmt.Printf("%v\n", r)
							return cenums.EVENT_PASS
						})
					}
					display.RequestQuit()
				}
			})
		} else {
			builder.LogError("missing main-dialog")
		}
	} else {
		builder.LogError("missing main-window")
	}
	return nil
}

func setupUiMsgbox(ctx *cli.Context, builder ctk.Builder, app ctk.Application, dm cdk.Display) error {
	builder.AddNamedSignalHandler("msgbox-ok", func(data []interface{}, argv ...interface{}) cenums.EventFlag {
		if dialog := getDialog(builder); dialog != nil {
			dialog.Response(enums.ResponseOk)
		} else {
			builder.LogError("msgbox-ok missing main-dialog")
		}
		return cenums.EVENT_STOP
	})
	if tmpl, err := template.New("msgbox").Parse(gladeMsgBox); err != nil {
		return err
	} else if tmpl != nil {
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

func setupUiYesNo(ctx *cli.Context, builder ctk.Builder, app ctk.Application, dm cdk.Display) error {
	builder.AddNamedSignalHandler("yesno-yes", func(data []interface{}, argv ...interface{}) cenums.EventFlag {
		if dialog := getDialog(builder); dialog != nil {
			dialog.Response(enums.ResponseYes)
		} else {
			builder.LogError("yesno-yes missing main-dialog")
		}
		return cenums.EVENT_STOP
	})
	builder.AddNamedSignalHandler("yesno-no", func(data []interface{}, argv ...interface{}) cenums.EventFlag {
		if dialog := getDialog(builder); dialog != nil {
			dialog.Response(enums.ResponseNo)
		} else {
			builder.LogError("yesno-no missing main-dialog")
		}
		return cenums.EVENT_STOP
	})
	if tmpl, err := template.New("yesno").Parse(gladeYesNo); err != nil {
		return err
	} else if tmpl != nil {
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
