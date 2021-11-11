package glade

import (
	"fmt"
	"io/ioutil"

	"github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/cdk/lib/ptypes"
	"github.com/urfave/cli/v2"

	"github.com/go-curses/cdk"
	cpaths "github.com/go-curses/cdk/lib/paths"
	"github.com/go-curses/cdk/log"
	"github.com/go-curses/ctk"
)

var CliCommand = &cli.Command{
	Name:        "glade",
	Usage:       "preview glade interfaces",
	Description: "load the given .glade file and preview in CTK",
	Action:      glade,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "window",
			Aliases:     []string{"W"},
			Value:       "main-window",
			DefaultText: "main-window",
			Usage:       "specify a named (glade \"id\" attribute) window to preview",
		},
		&cli.StringFlag{
			Name:        "dialog",
			Aliases:     []string{"D"},
			Value:       "main-dialog",
			DefaultText: "main-dialog",
			Usage:       "specify a named (glade \"id\" attribute) dialog to preview",
		},
		&cli.BoolFlag{
			Name:    "no-dialog-transient",
			Aliases: []string{"n"},
			Value:   false,
			Usage:   "when rendering ctk.Dialog types, do not set the transient for to a default window and use the dialog itself as a top-level window",
		},
	},
}

func glade(ctx *cli.Context) error {
	_ = log.StartRestart()
	argc := ctx.Args().Len()
	if argc == 0 {
		return cli.Exit("Error, missing glade interface file\nsee: go-ctk glade --help", 1)
	}
	gladeFile := ctx.Args().Get(0)
	if !cpaths.IsFile(gladeFile) {
		fmt.Printf("file not found: %v\n", gladeFile)
	} else if gladeFile[len(gladeFile)-6:] != ".glade" {
		fmt.Printf("not a .glade interface file: %v\n", gladeFile)
	}
	app := cdk.NewApp("ctk-glade", "", "", "", "ctk-glade", "CTK Glade", "/dev/tty", func(d cdk.Display) error {
		return ProcessFile(ctx, gladeFile, d)
	})
	if err := app.Run([]string{"ctk-glade"}); err != nil {
		log.Fatal(err)
	}
	return nil
}

func ProcessFile(ctx *cli.Context, path string, dm cdk.Display) (err error) {
	var bytes []byte
	if bytes, err = ioutil.ReadFile(path); err != nil {
		return fmt.Errorf("error reading glade file: %v", err)
	}
	contents := string(bytes)
	builder := ctk.NewBuilder()
	if _, err := builder.LoadFromString(contents); err != nil {
		return err
	}
	for _, bi := range builder.GetWidgetsBuiltByType(ctk.TypeButton) {
		if button, ok := bi.(ctk.Button); ok {
			button.Connect(ctk.SignalActivate, fmt.Sprintf("glade.activate"), func(data []interface{}, argv ...interface{}) enums.EventFlag {
				dm.AddQuitHandler(button.ObjectName(), func() {
					if label, err := button.GetStringProperty(ctk.PropertyLabel); err == nil && label != "" {
						fmt.Printf("button pressed: %v \"%v\"\n", button.ObjectName(), label)
					} else {
						fmt.Printf("button pressed: %v \"%v\"\n", button.ObjectName(), button.GetLabel())
					}
				})
				dm.RequestQuit()
				return enums.EVENT_STOP
			})
		}
	}
	window := ctx.String("window")
	dialog := ctx.String("dialog")
	if dialog != "" {
		if do := builder.GetWidget(dialog); do != nil {
			return setupUi(builder, do, dm)
		}
	}
	if window != "" {
		if wo := builder.GetWidget(window); wo != nil {
			return setupUi(builder, wo, dm)
		}
	}
	return fmt.Errorf("auto-window selection not implemented yet")
}

func setupUi(builder ctk.Builder, widget interface{}, dm cdk.Display) error {
	if widget != nil {
		if dialog, ok := widget.(ctk.Dialog); ok {
			return setupUiDialog(dialog, dm)
		}
		if window, ok := widget.(ctk.Window); ok {
			window.Show()
			dm.CaptureCtrlC()
			dm.SetActiveWindow(window)
			return nil
		}
		return fmt.Errorf("widget is not a window or dialog")
	}
	return fmt.Errorf("widget is nil")
}

func setupUiDialog(dialog ctk.Dialog, dm cdk.Display) error {
	window := dialog.GetTransientFor()
	if window != nil {
		dm.SetActiveWindow(window)
	} else {
		dm.SetActiveWindow(dialog)
	}
	if display := dm.Screen(); display != nil {
		dw, dh := display.Size()
		sr := ptypes.NewRectangle(dw/2, dh/2)
		sr.Clamp(20, 10, dw, dh)
		dialog.SetSizeRequest(sr.W, sr.H)
		dialog.SetTransientFor(window)
	}
	dialog.Show()
	dialog.LogInfo("starting Run()")
	response := dialog.Run()
	go func() {
		select {
		case r := <-response:
			dialog.Destroy()
			_ = dialog.DestroyObject()
			dm.AddQuitHandler("dialog-response", func() {
				fmt.Printf("dialog response: %v\n", r)
			})
			dm.RequestQuit()
		}
	}()
	return nil
}
