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
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/go-curses/cdk"
	cstrings "github.com/go-curses/cdk/lib/strings"
	"github.com/go-curses/cdk/log"
	"github.com/go-curses/ctk/cmd/go-ctk/ctkfmt"
	"github.com/go-curses/ctk/cmd/go-ctk/glade"
	"github.com/go-curses/ctk/cmd/go-ctk/gtkdoc2ctk"
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
	app := &cli.App{
		Name:   "go-ctk",
		Usage:  "Curses Tool Kit",
		Action: action,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "log-file",
				Aliases:     []string{"log"},
				Usage:       "enable logging and specify log path to write to",
				Value:       os.TempDir() + string(os.PathSeparator) + "go-ctk.log",
				DefaultText: os.TempDir() + string(os.PathSeparator) + "go-ctk.log",
			},
		},
		Commands: []*cli.Command{
			ctkfmt.CliCommand,
			gtkdoc2ctk.CliCommand,
			glade.CliCommand,
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func action(c *cli.Context) error {
	fmt.Println("this should display useful details about the current terminal environment")
	return nil
}