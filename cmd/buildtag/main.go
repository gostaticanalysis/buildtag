package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/gostaticanalysis/buildtag"
)

var (
	flagTags   string
	flagFormat string
)

func init() {
	flag.StringVar(&flagTags, "tags", "", "comma separated tags")
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

	ok := createEvalFunc()
	patterns := args
	info, err := buildtag.Load(patterns...)
	if err != nil {
		return err
	}

	tmpl, err := template.New("buildtag").Parse(flagFormat)
	if err != nil {
		return err
	}

	for _, e := range info.Matches(ok) {
		if err := tmpl.Execute(os.Stdout, e); err != nil {
			return err
		}
		fmt.Println()
	}

	return nil
}

func createEvalFunc() func(tag string) bool {
	if flagTags == "" {
		return nil
	}

	tags := strings.Split(flagTags, ",")
	for i := range tags {
		tags[i] = strings.TrimSpace(tags[i])
	}

	return func(tag string) bool {
		for i := range tags {
			if tag == tags[i] {
				return true
			}
		}
		return false
	}
}
