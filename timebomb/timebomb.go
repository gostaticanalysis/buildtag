package timebomb

import (
	"go/token"
	"go/types"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gostaticanalysis/analysisutil"
	"github.com/gostaticanalysis/buildtag"
	"github.com/tenntenn/goplayground"
	"golang.org/x/mod/semver"
	"golang.org/x/tools/go/analysis"
)

const doc = "timebomb finds suspicious buid tags"

// latest version of Go (including dev version)
var (
	latest  string
	release string
)

var (
	flagDeps bool
)

func init() {
	cli := goplayground.Client{
		Backend: "gotip",
	}
	ver, err := cli.Version()
	if err != nil {
		panic("cannot get dev version from Go Playground")
	}

	if strings.HasPrefix(ver.Release, "go") {
		release = ver.Release
		latest = "v" + ver.Release[2:]
		return
	}

	// get released version
	cli.Backend = ""
	ver, err = cli.Version()
	if err != nil {
		panic("cannot get version from Go Playground")
	}

	if strings.HasPrefix(ver.Release, "go") {
		release = ver.Release
		latest = "v" + ver.Release[2:]
	} else {
		panic("cannot get version from Go Playground")
	}
}

var Analyzer = &analysis.Analyzer{
	Name: "timebomb",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		buildtag.Analyzer,
	},
	FactTypes: []analysis.Fact{
		(*hasReleaseTag)(nil),
	},
}

func init() {
	Analyzer.Flags.BoolVar(&flagDeps, "deps", true, "print diagnostics of dependent packages")
}

type constraint struct {
	Value    string
	Position token.Position
}

type hasReleaseTag struct {
	Constraints []constraint
}

func (*hasReleaseTag) AFact() {}

func run(pass *analysis.Pass) (interface{}, error) {

	// dependent packages
	if flagDeps {
		for _, fact := range pass.AllPackageFacts() {
			f, ok := fact.Fact.(*hasReleaseTag)
			if !ok {
				continue
			}

			for _, c := range f.Constraints {
				pos := importedPos(pass, fact.Package)
				pass.Reportf(pos, "%s has a suspicious build constraint (%s) in %s because latest version (including dev version) is %s", fact.Package, c.Value, filepath.Base(c.Position.String()), release)
			}
		}
	}

	info := pass.ResultOf[buildtag.Analyzer].(*buildtag.Info)

	entries := info.Find(func(tag string) bool {
		if !strings.HasPrefix(tag, "go") {
			return false
		}
		ver := "v" + tag[2:]
		return semver.Compare(ver, latest) >= 0
	})

	if len(entries) == 0 {
		return nil, nil
	}

	fact := &hasReleaseTag{
		Constraints: make([]constraint, len(entries)),
	}

	for i := range entries {
		fact.Constraints[i] = constraint{
			Position: pass.Fset.Position(entries[i].Pos),
			Value:    entries[i].Constraint.String(),
		}
		pass.Reportf(entries[i].Pos, "%s: the build constraint is suspicious because latest version (including dev version) is %s", entries[i].Constraint, release)
	}

	pass.ExportPackageFact(fact)

	return nil, nil
}

func importedPos(pass *analysis.Pass, pkg *types.Package) token.Pos {
	fs := pass.Files
	if len(fs) == 0 {
		return token.NoPos
	}

	for _, f := range fs {
		for _, i := range f.Imports {
			path, err := strconv.Unquote(i.Path.Value)
			if err != nil {
				continue
			}
			if analysisutil.RemoveVendor(path) == pkg.Path() {
				return i.Pos()
			}
		}
	}
	return token.NoPos
}
