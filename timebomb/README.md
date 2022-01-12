# buildtag/timebomb

[![pkg.go.dev][gopkg-badge]][gopkg]

`timebomb` timebomb finds suspicious build tags such as below.

* a build tag which has "go" prefix and it corresponds to a version which does not start development yet

## Install

You can get `timebomb` by `go install` command (Go 1.16 and higher).

```bash
$ go install github.com/gostaticanalysis/buildtag/timebomb/cmd/timebomb@latest
```

## How to use

`timebomb` run with `go vet` as below when Go is 1.12 and higher.

```bash
$ go vet -vettool=$(which timebomb) ./...
```

## Analyze with golang.org/x/tools/go/analysis

You can use [timebomb.Analyzer](https://pkg.go.dev/github.com/gostaticanalysis/buildtag/timebomb/#Analyzer) with [unitchecker](https://golang.org/x/tools/go/analysis/unitchecker).

<!-- links -->
[gopkg]: https://pkg.go.dev/github.com/gostaticanalysis/buildtag/timebomb
[gopkg-badge]: https://pkg.go.dev/badge/github.com/gostaticanalysis/buildtag/timebomb?status.svg
