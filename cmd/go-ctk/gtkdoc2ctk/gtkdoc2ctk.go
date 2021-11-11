package gtkdoc2ctk

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/urfave/cli/v2"

	cpaths "github.com/go-curses/cdk/lib/paths"
	cstrings "github.com/go-curses/cdk/lib/strings"
	log "github.com/go-curses/cdk/log"
)

var CliCommand = &cli.Command{
	Name:   "gtk2doc",
	Usage:  "generate CTK boilerplate source code from GTK2 documentation",
	Action: gtk2doc,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "page",
			Usage: "parse the GTK2 developers.gnome.org documentation and output go-ctk source code",
		},
		&cli.StringFlag{
			Name:  "path",
			Usage: "parse the GTK2 documentation and output go-ctk source code",
		},
		&cli.StringFlag{
			Name:        "package-name",
			Aliases:     []string{"package", "pkg"},
			Usage:       "specify the package name for the generated source code",
			Value:       "ctk",
			DefaultText: "ctk",
		},
		&cli.StringFlag{
			Name:  "output",
			Usage: "write generated source code to the file path given instead of stdout",
		},
		&cli.StringFlag{
			Name:  "output-path",
			Usage: "write generated source code to a new file at the path given instead of stdout",
		},
		&cli.BoolFlag{
			Name:    "force",
			Aliases: []string{"f"},
			Usage:   "overwrite output file if it exists",
		},
		&cli.BoolFlag{
			Name:    "include-deprecated",
			Aliases: []string{"D"},
			Usage:   "include deprecated methods",
		},
	},
}

func gtk2doc(c *cli.Context) error {
	if path := c.String("log-file"); path != "" {
		_ = os.Setenv("GO_CDK_LOG_FULL_PATHS", "true")
		_ = os.Setenv("GO_CDK_LOG_TIMESTAMPS", "true")
		_ = os.Setenv("GO_CDK_LOG_FORMAT", "pretty")
		_ = os.Setenv("GO_CDK_LOG_LEVEL", "debug")
		_ = os.Setenv("GO_CDK_LOG_OUTPUT", "file")
		_ = os.Setenv("GO_CDK_LOG_FILE", path)
	} else {
		_ = os.Setenv("GO_CDK_LOG_LEVEL", "error")
	}
	_ = log.StartRestart()
	if path := c.String("output-path"); path != "" {
		if !cpaths.IsDir(path) {
			return fmt.Errorf("--output-path does not exist or is not a directory")
		}
	}
	if page := c.String("page"); page != "" {
		return ParseGtk2DocUrl(page, c)
	}
	if path := c.String("path"); path != "" {
		return ParseGtk2DocFile(path, c)
	}
	return cli.Exit("missing --page or --path arguments", 1)
}

func ParseGtk2DocUrl(gtk2name string, c *cli.Context) (err error) {
	docUrl := ""
	if cstrings.IsUrl(gtk2name) {
		docUrl = gtk2name
	} else {
		docUrl = fmt.Sprintf("https://developer.gnome.org/gtk2/stable/%v.html", gtk2name)
	}
	response, err := http.Get(docUrl)
	if err != nil {
		return cli.Exit(err, 1)
	}
	defer func() { _ = response.Body.Close() }()
	snakedName := strcase.ToSnake(gtk2name)
	if strings.HasPrefix(snakedName, "gdk_") {
		snakedName = snakedName[:4]
	} else if strings.HasPrefix(snakedName, "gtk_") {
		snakedName = snakedName[:4]
	}
	content := ""
	for {
		var b []byte
		if _, err := response.Body.Read(b); err != nil {
			break
		} else {
			content += string(b)
		}
	}
	if src, err := ParseGtk2Doc(snakedName, content, c); err != nil {
		return cli.Exit(err, 1)
	} else {
		if code, err := GenerateCtkSource(c, src); err != nil {
			return cli.Exit(err, 1)
		} else {
			return ProduceOutput(src.Name, code, c)
		}
	}
}

func ParseGtk2DocFile(path string, c *cli.Context) (err error) {
	content, cErr := cpaths.ReadFile(path)
	if cErr != nil {
		return cli.Exit(cErr, 1)
	}
	bName := filepath.Base(path)
	if strings.HasSuffix(bName, ".html") {
		bName = bName[:len(bName)-5]
	}
	snakedName := strcase.ToSnake(bName)
	if strings.HasPrefix(snakedName, "gdk_") {
		snakedName = snakedName[4:]
	} else if strings.HasPrefix(snakedName, "gtk_") {
		snakedName = snakedName[4:]
	} else if strings.HasPrefix(snakedName, "gtk2-") {
		snakedName = snakedName[5:]
	}
	if src, err := ParseGtk2Doc(snakedName, content, c); err != nil {
		return cli.Exit(err, 1)
	} else {
		if code, err := GenerateCtkSource(c, src); err != nil {
			return cli.Exit(err, 1)
		} else {
			return ProduceOutput(src.Name, code, c)
		}
	}
}

func WriteToFile(c *cli.Context, targetFileName string, content []byte) (err error) {
	overwrote := false
	if cpaths.IsFile(targetFileName) {
		if !c.Bool("force") {
			return fmt.Errorf("output file exists: %v, use --force to overwrite", targetFileName)
		}
		if err = os.Remove(targetFileName); err != nil {
			return
		}
		overwrote = true
	}
	if err = ioutil.WriteFile(targetFileName, content, 0664); err != nil {
		return
	}
	if overwrote {
		fmt.Printf("overwrote: %v\n", targetFileName)
	} else {
		fmt.Printf("wrote: %v\n", targetFileName)
	}
	return
}

func ProduceOutput(name string, code string, c *cli.Context) (err error) {
	if targetFileName := c.String("output"); targetFileName != "" {
		return WriteToFile(c, targetFileName, []byte(code))
	} else if targetPath := c.String("output-path"); targetPath != "" {
		if cpaths.IsDir(targetPath) {
			targetFileName := targetPath + string(os.PathSeparator) + strcase.ToSnake(name) + ".go"
			return WriteToFile(c, targetFileName, []byte(code))
		} else {
			return fmt.Errorf("directory not found: %v", targetPath)
		}
	} else {
		fmt.Printf("%v", code)
	}
	return
}
