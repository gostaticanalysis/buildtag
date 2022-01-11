package buildtag

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/build/constraint"
	"go/parser"
	"go/token"
	"os"
	"sort"
	"strings"

	"golang.org/x/tools/go/packages"
)

func Load(patterns ...string) (*Info, error) {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedModule | packages.NeedTypes | packages.NeedSyntax | packages.NeedFiles,
	}

	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		return nil, fmt.Errorf("buildtag.Load: %w", err)
	}

	var info Info
	for _, pkg := range pkgs {
		if err := parsePkg(&info, pkg); err != nil {
			return nil, fmt.Errorf("buildtag.Load: %w", err)
		}
	}

	sort.Slice(info.Entries, func(i, j int) bool {
		return info.Entries[i].File < info.Entries[j].File
	})

	return &info, nil
}

func parsePkg(info *Info, pkg *packages.Package) error {
	for _, f := range pkg.Syntax {
		expr, pos, err := parseFile(f)
		if err != nil {
			return err
		}

		fname := pkg.Fset.File(f.Pos()).Name()
		info.add(pkg, fname, pos, expr)
	}

	for _, fname := range pkg.OtherFiles {
		expr, pos, err := parseOtherFile(pkg.Fset, fname)
		if err != nil {
			return err
		}

		info.add(pkg, fname, pos, expr)
	}

	for _, fname := range pkg.IgnoredFiles {
		if strings.HasSuffix(fname, ".go") {
			f, err := parser.ParseFile(pkg.Fset, fname, nil, parser.ParseComments)
			if err != nil {
				continue
			}

			expr, pos, err := parseFile(f)
			if err != nil {
				return err
			}

			info.add(pkg, fname, pos, expr)

		} else {
			expr, pos, err := parseOtherFile(pkg.Fset, fname)
			if err != nil {
				return err
			}

			info.add(pkg, fname, pos, expr)
		}
	}

	return nil
}

func parseFile(file *ast.File) (constraint.Expr, token.Pos, error) {
	for _, cg := range file.Comments {
		for _, c := range cg.List {
			if constraint.IsGoBuild(c.Text) {
				expr, err := constraint.Parse(c.Text)
				if err != nil {
					return nil, token.NoPos, err
				}
				return expr, c.Pos(), nil
			}
		}
	}

	return nil, token.NoPos, nil
}

// copy from golang.org/x/tools/go/analysis/passes/internal/analysisutil/util.go
func readFile(fset *token.FileSet, filename string) ([]byte, *token.File, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, nil, err
	}
	tf := fset.AddFile(filename, -1, len(content))
	tf.SetLinesForContent(content)
	return content, tf, nil
}

func parseOtherFile(fset *token.FileSet, fname string) (constraint.Expr, token.Pos, error) {
	content, tf, err := readFile(fset, fname)
	if err != nil {
		return nil, token.NoPos, err
	}

	full := content
	for len(content) != 0 {
		i := bytes.Index(content, []byte("\n"))
		if i < 0 {
			i = len(content)
		} else {
			i++
		}

		offset := len(full) - len(content)
		line := strings.TrimSpace(string(content[:i]))
		content = content[i:]

		if line == "" {
			continue
		}

		if !strings.HasPrefix(line, "//") {
			break
		}

		if constraint.IsGoBuild(line) {
			expr, err := constraint.Parse(line)
			if err != nil {
				return nil, token.NoPos, err
			}

			return expr, tf.Pos(offset), nil
		}
	}

	return nil, token.NoPos, nil
}

type Entry struct {
	Package    *packages.Package
	File       string
	Pos        token.Pos
	Constraint constraint.Expr
}

type Info struct {
	Entries []Entry
}

func (info *Info) add(pkg *packages.Package, file string, pos token.Pos, expr constraint.Expr) {
	if expr == nil {
		return
	}

	info.Entries = append(info.Entries, Entry{
		Package:    pkg,
		File:       file,
		Pos:        pos,
		Constraint: expr,
	})
}

func (info *Info) Matches(ok func(tag string) bool) []Entry {
	var matches []Entry

	if ok == nil {
		matches = make([]Entry, len(info.Entries))
		copy(matches, info.Entries)
		return matches
	}

	for i := range info.Entries {
		if info.Entries[i].Constraint.Eval(ok) {
			matches = append(matches, info.Entries[i])
		}
	}

	return matches
}
