package buildtag

import (
	"fmt"
	"sort"

	"golang.org/x/tools/go/packages"
)

func Load(patterns ...string) (*Info, error) {
	cfg := &packages.Config{
		Mode: packages.NeedTypesInfo | packages.NeedName | packages.NeedModule | packages.NeedTypes | packages.NeedSyntax | packages.NeedFiles,
	}

	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		return nil, fmt.Errorf("buildtag.Load: %w", err)
	}

	var info Info
	for _, pkg := range pkgs {
		p := newParserFromPkg(pkg)
		if err := p.parse(&info); err != nil {
			return nil, fmt.Errorf("buildtag.Load: %w", err)
		}
	}

	sort.Slice(info.Entries, func(i, j int) bool {
		return info.Entries[i].File < info.Entries[j].File
	})

	return &info, nil
}
