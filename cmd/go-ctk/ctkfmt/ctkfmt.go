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

package ctkfmt

// TODO: fmt messes up interface for: ctype_item.go
// TODO: add `go fmt` to the mix, perhaps vet etc too?
// TODO: implement a "list all go source files with valid interface/struct pairs"
// TODO: support multiple interface/struct pairs per file

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/urfave/cli/v2"

	cpaths "github.com/go-curses/cdk/lib/paths"
	cstrings "github.com/go-curses/cdk/lib/strings"
	log "github.com/go-curses/cdk/log"
)

var CliCommand = &cli.Command{
	Name:        "fmt",
	Usage:       "reformat a Go-CTK source file",
	Description: "scan the source file for an interface and matching concrete type, along with any exported methods, finally rewriting the interface to match the exported methods",
	Action:      ctkFmt,
}

func ctkFmt(c *cli.Context) error {
	if path := c.String("log-file"); path != "" {
		_ = os.Setenv("GO_CDK_LOG_FULL_PATHS", "true")
		_ = os.Setenv("GO_CDK_LOG_TIMESTAMPS", "true")
		_ = os.Setenv("GO_CDK_LOG_LEVEL", "debug")
		_ = os.Setenv("GO_CDK_LOG_OUTPUT", "file")
		_ = os.Setenv("GO_CDK_LOG_FILE", path)
	} else {
		_ = os.Setenv("GO_CDK_LOG_OUTPUT", "stdout")
		_ = os.Setenv("GO_CDK_LOG_LEVEL", "error")
	}
	_ = log.StartRestart()
	argc := c.Args().Len()
	if argc == 0 {
		return cli.Exit("Error, missing go source files to fmt.\nsee: go-ctk fmt --help", 1)
	}
	for i := 0; i < argc; i++ {
		arg := c.Args().Get(i)
		if !cpaths.IsFile(arg) {
			fmt.Printf("file not found: %v\n", arg)
			continue
		}
		if arg[len(arg)-3:] != ".go" {
			fmt.Printf("not a .go source file: %v\n", arg)
			continue
		}
		if err := ProcessFile(arg); err != nil {
			fmt.Printf("error processing file: %v - %v\n", arg, err)
		}
	}
	return nil
}

var (
	rxFindInterface = regexp.MustCompile(`(?ms)^\s*type (\w+) interface \{\s*$`)
	rxFindStruct    = regexp.MustCompile(`(?ms)^\s*type (\w+) struct \{\s*$`)
)

func ProcessFile(path string) (err error) {
	var bytes []byte
	if bytes, err = ioutil.ReadFile(path); err != nil {
		return
	}
	var facade, concrete string
	var methods []string
	var contents, actual string
	actual = string(bytes)
	contents = actual
	if rxFindInterface.MatchString(contents) {
		mi := rxFindInterface.FindStringSubmatch(contents)
		if len(mi) == 2 {
			facade = mi[1]
			if rxFindStruct.MatchString(contents) {
				ms := rxFindStruct.FindStringSubmatch(contents)
				if len(ms) == 2 {
					concrete = ms[1]
					rxFindExported := regexp.MustCompile(`(?ms)^\s*func \(.+?\b\Q` + concrete + `\E\) ([A-Z].+?) \{\}??\s*$`)
					if rxFindExported.MatchString(contents) {
						me := rxFindExported.FindAllStringSubmatch(contents, -1)
						if len(me) >= 1 {
							for _, method := range me {
								methods = append(methods, method[1])
							}
						}
					}
				}
			}
		}
	}
	if facade == "" {
		return fmt.Errorf("%v is missing an interface declaration\n", path)
	}
	if concrete == "" {
		return fmt.Errorf("%v is missing a matching struct declaration\n", path)
	}
	if len(methods) == 0 {
		return fmt.Errorf("%v is missing exported methods", path)
	}
	rx := regexp.MustCompile(`(?ms)^type ` + facade + ` interface \{\s*\r??\n(.+?)\r??\n\s*\}\s*$`)
	if rx.MatchString(contents) {
		m := rx.FindStringSubmatch(contents)
		if len(m) > 1 {
			var save []string
			var found bool
			for _, line := range strings.Split(m[1], "\n") {
				if cstrings.IsEmpty(line) {
					found = true
					break
				}
				save = append(save, strings.TrimSpace(line))
			}
			current := "type " + facade + " interface {\n"
			if found {
				for _, saved := range save {
					current += "\t" + saved + "\n"
				}
				if len(save) > 0 {
					current += "\n"
				}
			}
			for _, method := range methods {
				current += "\t" + method + "\n"
			}
			current += "}\n"
			contents = rx.ReplaceAllString(contents, current)
		}
	}
	if contents != actual {
		if stat, err := os.Stat(path); err != nil {
			fmt.Printf("error: %v - %v\n", path, err)
		} else {
			if err := ioutil.WriteFile(path, []byte(contents), stat.Mode().Perm()); err != nil {
				fmt.Printf("error: %v - %v\n", path, err)
			} else {
				fmt.Printf("updated: %v\n", path)
			}
		}
	} else {
		fmt.Printf("skipped: %v\n", path)
	}
	return nil
}