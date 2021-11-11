// +build example_hello_world

package main

import (
	"os"

	"github.com/go-curses/cdk"
	"github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/cdk/log"
	"github.com/go-curses/ctk"
)

func main() {
	// Construct a new CDK application
	app := cdk.NewApp(
		// program binary name
		"hello-world",
		// usage summary
		"Simple Hello World example for CTK",
		// description
		"A simple terminal program written using the Curses Tool Kit",
		// because versioning is important
		"0.0.1",
		// used in logs, internal debugging, etc
		"helloWorld",
		// used where human-readable titles are necessary
		"Hello World",
		// the TTY device to use, /dev/tty is the default
		"/dev/tty",
		// initialize the user-interface
		func(d cdk.Display) error {
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
			button := ctk.NewButtonWithLabel("Curses<u><i>!</i></u>")
			button.SetUseMarkup(true)    // enable markup in the label
			button.SetSizeRequest(11, 3) // request a certain size
			// make the button quit the application when activated by connecting
			// a handler to the button's activate signal
			button.Connect(
				ctk.SignalActivate,
				"hello-button-handle",
				func(data []interface{}, argv ...interface{}) enums.EventFlag {
					d.RequestQuit() // ask the display to exit nicely
					return enums.EVENT_STOP
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
			// tell CDK that this window is the foreground window
			d.SetActiveWindow(w)
			// no errors to report, nil to proceed
			return nil
		},
	)
	// run the application, handing over the command-line arguments received
	if err := app.Run(os.Args); err != nil {
		// doesn't have to be a Fatal exit
		log.Fatal(err)
	}
	// end of program
}
