package buildtag

import (
	"reflect"

	"golang.org/x/tools/go/analysis"
)

const doc = "buildtag collect buid constraints from source files"

var Analyzer = &analysis.Analyzer{
	Name:       "buildtag",
	Doc:        doc,
	Run:        run,
	ResultType: reflect.TypeOf((*Info)(nil)),
}

func run(pass *analysis.Pass) (interface{}, error) {

	var info Info
	p := newParserFromPass(pass)
	if err := p.parse(&info); err != nil {
		return nil, err
	}

	return &info, nil
}
