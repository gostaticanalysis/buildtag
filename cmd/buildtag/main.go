package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"text/template"

	"github.com/gostaticanalysis/buildtag"
)

var (
	flagTag    string
	flagFormat string
)

func init() {
	flag.StringVar(&flagTag, "tag", "", "a regular expression of tag")
	flag.StringVar(&flagFormat, "f", "{{.File}}:{{.Constraint}}", "output format")
	flag.Parse()
}

func main() {
	if err := run(flag.Args()); err != nil {
		fmt.Fprintln(os.Stderr, "buildtag:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		return errors.New("package patterns must be specified")
	}

	patterns := args
	info, err := buildtag.Load(patterns...)
	if err != nil {
		return err
	}

	tmpl, err := template.New("buildtag").Parse(flagFormat)
	if err != nil {
		return err
	}

	reg, err := regexp.Compile(flagTag)
	if err != nil {
		return err
	}

	for _, e := range info.FindByRegexp(reg) {
		if err := tmpl.Execute(os.Stdout, e); err != nil {
			return err
		}
		fmt.Println()
	}

	return nil
}
