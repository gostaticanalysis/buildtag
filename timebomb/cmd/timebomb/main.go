package main

import (
	"github.com/gostaticanalysis/buildtag/timebomb"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(timebomb.Analyzer) }
