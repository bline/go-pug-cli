package main

import (
	"errors"
	"github.com/Joker/jade"
	"github.com/urfave/cli"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	app := cli.NewApp()
	app.Name = "go-pug-cli"
	app.Usage = "convert pug files to html/template files"
	app.Authors = []cli.Author{
		{
			Name:  "Scott Beck",
			Email: "scottbeck@gmail.com",
		},
	}
	app.ArgsUsage = "src dest"
	app.Before = func(c *cli.Context) error {
		args := c.Args()
		if len(args) != 2 {
			return errors.New("must specify src and dest")
		}
		return nil
	}
	app.Flags = []cli.Flag{
		cli.BoolFlag{Name: "pretty", Destination: &jade.PrettyOutput, Usage: "pretty print output html"},
		cli.StringFlag{Name: "indent_str", Destination: &jade.OutputIndent, Usage: "string to use for indenting. use quotes: \"    \" for 4 spaces etc"},
		cli.StringFlag{Name: "right_delim", Destination: &jade.RightDelim, Usage: "set the start delimiter for output template/html template"},
		cli.StringFlag{Name: "left_delim", Destination: &jade.LeftDelim, Usage: "set the start delimiter for output template/html template"},
	}
	app.Action = handler
	app.Run(os.Args)
	return
}

func handler(c *cli.Context) error {
	args := c.Args()
	fromPath := args[0]
	toPath := args[1]
	pathSep := string(os.PathSeparator)
	var npathArr []string
	fromPathSecs := strings.Split(fromPath, pathSep)
	toPathSecs := strings.Split(toPath, pathSep)
	for i, sec := range fromPathSecs {
		npathArr = append(npathArr, sec)
		if i > len(toPathSecs) {
			break
		}
		if sec != toPathSecs[i] {
			break
		}
	}
	prefix := ""
	if strings.HasPrefix(fromPath, pathSep) {
		prefix = pathSep
	}
	log.Printf("FromPath: %v; ToPath: %v", fromPath, toPath)
	basePath := prefix + strings.Join(npathArr, pathSep)
	log.Printf("basePath: %v", basePath)
	var errs []error
	filepath.Walk(fromPath, filepath.WalkFunc(func(path string, f os.FileInfo, err error) error {
		if f.IsDir() {
			return nil
		}
		destPath := strings.Replace(path, basePath, toPath, 1)
		dirName := filepath.Base(destPath)
		if _, err = os.Stat(dirName); os.IsNotExist(err) {
			err = os.MkdirAll(dirName, 0755)
			if err != nil {
				errs = append(errs, err)
				log.Printf("mkdir failed: %v", err)
				return nil
			}
		}
		tplStr, err := jade.ParseFile(path)
		if err != nil {
			errs = append(errs, err)
			log.Printf("parse error: %v", err)
			return nil
		}
		err = ioutil.WriteFile(destPath, []byte(tplStr), f.Mode())
		if err != nil {
			errs = append(errs, err)
			log.Printf("parse error: %v", err)
		}
		return nil
	}))
	if len(errs) > 0 {
		log.Fatalf("Had %d errors", len(errs))
	}
	return nil
}
