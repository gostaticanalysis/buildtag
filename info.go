package buildtag

import (
	"go/build/constraint"
	"go/token"
	"go/types"
	"regexp"
)

type Entry struct {
	Package    *types.Package
	File       string
	Pos        token.Pos
	Constraint constraint.Expr
}

type Info struct {
	Entries []Entry
}

func (info *Info) add(pkg *types.Package, file string, pos token.Pos, expr constraint.Expr) {
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

func (info *Info) Find(f func(constraint string) bool) []Entry {
	var founds []Entry

	for i := range info.Entries {
		if f(info.Entries[i].Constraint.String()) {
			founds = append(founds, info.Entries[i])
		}
	}

	return founds
}

func (info *Info) FindByRegexp(reg *regexp.Regexp) []Entry {
	var founds []Entry

	for i := range info.Entries {
		if reg.MatchString(info.Entries[i].Constraint.String()) {
			founds = append(founds, info.Entries[i])
		}
	}

	return founds
}
