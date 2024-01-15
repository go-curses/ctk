[![Go-Curses](https://go-curses.org/media/curses-logo-banner.png)](https://go-curses.org)

[![Made with Go](https://img.shields.io/badge/go-v1.16+-blue.svg)](https://golang.org)
[![Go documentation](https://img.shields.io/badge/godoc-reference-blue.svg)](https://pkg.go.dev/github.com/go-curses/ctk)

# CTK - Curses Tool Kit

Golang package to provide an advanced terminal user interface with a [GTK]
inspired API, built upon [CDK].

## Notice

This project should not be used for any purpose other than intellectual
curiosity. This status is reflected in the tagged versioning of this `trunk`
branch, v0.1.x, ie: entirely experimental and unfinished in any sense of the
word "done".

## Getting Started

CTK is a Go module and as such can be used in any of the typical Golang ways.

### Prerequisites

Go v1.16 (or later) is required in order to build and use the package. Beyond
that, there aren't any other dependencies. Visit: https://golang.org/doc/install
for installation instructions.

### Installing

CTK uses the Go mod system and is installed in any of the usual ways.

```
$ go get -u github.com/go-curses/ctk/...
```

### Programming Hello World in CTK

The following application will display a window with "Hello World" as the title, containing a button centered nicely
that when pressed will exit the application nicely.

```go
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
  app := cdk.NewApplication(
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
      // add a quit handler to say goodbye when the program exits
      d.AddQuitHandler(
        "hello-world-quit-handler",
        func() {
          // Note that the Display and other CTK things are no longer
          // functional at this point.
          fmt.Println("Hello World says Goodbye!")
          // Logging however still works.
          log.InfoF("Hello World logging goodbye!")
        },
      )
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
```

Compile the [hello-world.go] source file.

```
$ go build examples/hello-world/hello-world.go
```

View the command-line help:

```
$ ./hello-world -h
NAME:
   hello-world - hello-world

USAGE:
   hello-world [global options] command [command options] [arguments...]

VERSION:
   0.0.1

DESCRIPTION:
   the most basic CTK application

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h, --usage  display command-line usage information (default: false)
   --version            display the version (default: false)
```

Run the program:

```
$ ./hello-world
```

![hello-world screenshot]

Pressing the Curses! button will exit the program and print the quit message to
the terminal, which should have cleanly cleared the screen and restored the
terminal to shell control again.

### Commands

CTK includes a number of command programs for the purpose of enabling a better
developer experience or as a means of having a necessity in creating new
widgets and features.

#### go-dialog

[go-dialog] is a [dialog] replacement, fully implemented in CTK and takes
advantage of [Glade] interface files for implementing the user interface.

##### Installation

```
$ go install github.com/go-curses/ctk/cmd/go-dialog
```

##### Usage

```
$ ./go-dialog --help
NAME:
   go-dialog - display dialog boxes from shell scripts

USAGE:
   go-dialog [global options] command [command options] [arguments...]

VERSION:
   0.0.1

COMMANDS:
   msgbox   display a message with an OK button, each string following msgbox is a new line and concatenated into the message
   yesno    display a yes/no prompt with a message (see msgbox)
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --title value                     specify the dialog title text
   --print-maxsize                   print the width and height on stdout and exit (default: false)
   --back-title value                specify the window title text
   --help, -h, --usage               display command-line usage information (default: false)
   --version                         display the version (default: false)
```

##### Example

Display a message-box dialog (with the title of "Hello Dialog"), centered on top
of a full-screen window (with the title "Hello Window"), displaying a message of
"This is the message" and presenting a single button labelled "OK".

```
$ ./go-dialog \
   --back-title "Hello Window" \
   --title "Hello Dialog" \
   msgbox \
   "This is the message."
```

![go-dialog screenshot]

#### go-charmap

This is a simple character-set viewer called [go-charmap].

##### Installation

```
$ go install github.com/go-curses/ctk/cmd/go-charmap
```

##### Usage

```
$ ./go-charmap --help
NAME:
   go-charmap - View the details of a character

USAGE:
   go-charmap [integer]

VERSION:
   0.0.1

DESCRIPTION:
   Get informational details for a specific character given in integer form.

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h, --usage  display command-line usage information (default: false)
   --version            display the version (default: false)

```

##### Example

Display information on the current terminal settings.

```
$ go-charmap
```

![go-charmap screenshot]

Display information on the lowercase sigma character.

```
$ go-charmap 1010
```

![go-charmap 1010 screenshot]

#### go-ctk

This is a weird one. [go-ctk] can do a number of things related to working
on the CTK project, some useful for normal developers. Particularly the
`glade` option, so that's the one we're going to talk about here.

The `glade` option for go-ctk enables the developer to test their [Glade]
interface files outside of their normal codebase, as a means of validating
whether their problems are in their code or with how CTK is handling the
[Glade] files.

For those not aware, [Glade] is the formal way to design user interfaces
for actual GTK applications. Given that CTK is a curses implementation of
GTK, it supports loading widgets and such from Glade interface files.

There are of course all sorts of undocumented and unknown caveats to the
support of Glade and all it's configurable features at this time. It is
an objective of CTK to formally support Glade as much as sensibly possible.

##### Usage

```
$ ./go-ctk help glade
NAME:
   go-ctk glade - preview glade interfaces

USAGE:
   go-ctk glade [command options] [arguments...]

DESCRIPTION:
   load the given .glade file and preview in CTK

OPTIONS:
   --window value, -W value   specify a named (glade "id" attribute) window to preview (default: main-window)
   --dialog value, -D value   specify a named (glade "id" attribute) dialog to preview (default: main-dialog)
   --no-dialog-transient, -n  when rendering ctk.Dialog types, do not set the transient-for to a default window and use the dialog itself as a top-level window (default: false)
```

##### Example

There are two example glade files in the examples directory. This
example command will load the `builder-dialog.glade` file and if
all is working as it should, display a window with a dialog containing
some lorem ipsum text and two buttons, a yes and a no.

For this to work however, it's important to note that within glade
we need to set the `id` on the window and dialog objects. go-ctk
defaults to "main-window" and "main-dialog", which the example files
use. That's the purpose to the `--window` and `--dialog` command-line
options.

```
$ ./go-ctk glade ./examples/builder-dialog.glade
```

![go-ctk glade screenshot]

Pressing either of the buttons will exit the viewer and print out which
button was pressed.

```
button pressed: ctk-button-2194440b-86b2-4386-a5a1-984bfec916e6 "Yes"
```

## Running the unit tests

Normal go testing mechanics work.

```
$ go test -v
=== RUN   TestAdjustment

  ... (per-test output, trimmed for brevity) ...

--- PASS: TestWidget (0.00s)
PASS
ok  	github.com/go-curses/ctk	(0.018s)
```

## Makefile

Included is a Makefile that has a number of useful build targets, as shown in
the usage help text.

```
$ make help
usage: make [target]

qa targets:
  vet         - run go vet command
  test        - perform all available tests
  cover       - perform all available tests with coverage report

cleanup targets:
  clean       - cleans package and built files
  clean-logs  - cleans *.log from the project

go.mod helpers:
  local       - add go.mod local CDK package replacements
  unlocal     - remove go.mod local CDK package replacements

build targets:
  examples    - builds all examples
  build       - build the go-ctk command
  build-all   - build all commands
  dev         - build demo-app with profiling

run targets:
  run         - run the dev build (sanely handle crashes)
  profile.cpu - run the dev build and profile CPU
  profile.mem - run the dev build and profile memory
```

## Versioning

The current API is unstable and subject to change dramatically.

## License

This project is licensed under the Apache License, Version 2.0 - see the
[LICENSE.md] file for details.

## Acknowledgments

* Thanks to [TCell] for providing a solid and robust platform to build upon
* Thanks to the [GTK Team] for developing and maintaining the [GTK] API that
  CTK is modeled after

[CDK]: https://github.com/go-curses/cdk
[go-ctk]: https://github.com/go-curses/ctk/blob/trunk/cmd/go-ctk
[go-charmap]: https://github.com/go-curses/ctk/blob/trunk/cmd/go-charmap
[go-dialog]: https://github.com/go-curses/ctk/blob/trunk/cmd/go-dialog
[hello-world.go]: https://github.com/go-curses/ctk/blob/trunk/examples/hello-world.go
[dialog]: https://invisible-island.net/dialog/
[kckrinke]: https://github.com/kckrinke
[LICENSE.md]: https://github.com/go-curses/ctk/blob/trunk/LICENSE.md
[TCell]: https://github.com/gdamore/tcell
[GTK Team]: https://www.gtk.org/development.php#Team
[GTK]: https://docs.gtk.org/gtk3/
[Glade]: https://glade.gnome.org/
[dialog]: https://invisible-island.net/dialog/
[go-dialog screenshot]: https://go-curses.org/screenshots/go-dialog--msgbox.png
[go-charmap screenshot]: https://go-curses.org/screenshots/go-charmap.png
[go-charmap 1010 screenshot]: https://go-curses.org/screenshots/go-charmap--1010.png
[go-ctk glade screenshot]: https://go-curses.org/screenshots/go-ctk--glade.png
[hello-world screenshot]: https://go-curses.org/screenshots/hello-world.png
