package timebomb

import (
	"strings"

	"github.com/gostaticanalysis/buildtag"
	"github.com/tenntenn/goplayground"
	"golang.org/x/mod/semver"
	"golang.org/x/tools/go/analysis"
)

const doc = "timebomb is ..."

// latest version of Go (including dev version)
var (
	latest  string
	release string
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

type hasReleaseTag struct{}

func (*hasReleaseTag) AFact() {}

func run(pass *analysis.Pass) (interface{}, error) {
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

	for i := range entries {
		pass.Reportf(entries[i].Pos, "%s: the build constraint is suspicious because latest version (including dev version) is %s", entries[i].Constraint, release)
	}

	return nil, nil
}
