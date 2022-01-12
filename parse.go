package buildtag

import (
	"bytes"
	"go/ast"
	"go/build/constraint"
	goparser "go/parser"
	"go/token"
	"go/types"
	"os"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/packages"
)

type parser struct {
	fset         *token.FileSet
	files        []*ast.File
	otherFiles   []string
	ignoredFiles []string
	pkg          *types.Package
}

func newParserFromPass(pass *analysis.Pass) *parser {
	return &parser{
		fset:         pass.Fset,
		files:        pass.Files,
		otherFiles:   pass.OtherFiles,
		ignoredFiles: pass.IgnoredFiles,
		pkg:          pass.Pkg,
	}
}

func newParserFromPkg(pkg *packages.Package) *parser {
	return &parser{
		fset:         pkg.Fset,
		files:        pkg.Syntax,
		otherFiles:   pkg.OtherFiles,
		ignoredFiles: pkg.IgnoredFiles,
		pkg:          pkg.Types,
	}
}

func (p *parser) parse(info *Info) error {
	for _, f := range p.files {
		expr, pos, err := p.parseFile(f)
		if err != nil {
			return err
		}

		fname := p.fset.File(f.Pos()).Name()
		info.add(p.pkg, fname, pos, expr)
	}

	for _, fname := range p.otherFiles {
		expr, pos, err := p.parseOtherFile(fname)
		if err != nil {
			return err
		}

		info.add(p.pkg, fname, pos, expr)
	}

	for _, fname := range p.ignoredFiles {
		if strings.HasSuffix(fname, ".go") {
			f, err := goparser.ParseFile(p.fset, fname, nil, goparser.ParseComments)
			if err != nil {
				continue
			}

			expr, pos, err := p.parseFile(f)
			if err != nil {
				return err
			}

			info.add(p.pkg, fname, pos, expr)

		} else {
			expr, pos, err := p.parseOtherFile(fname)
			if err != nil {
				return err
			}

			info.add(p.pkg, fname, pos, expr)
		}
	}

	return nil
}

func (p *parser) parseFile(file *ast.File) (constraint.Expr, token.Pos, error) {
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
func (p *parser) readFile(filename string) ([]byte, *token.File, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, nil, err
	}
	tf := p.fset.AddFile(filename, -1, len(content))
	tf.SetLinesForContent(content)
	return content, tf, nil
}

func (p *parser) parseOtherFile(fname string) (constraint.Expr, token.Pos, error) {
	content, tf, err := p.readFile(fname)
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
