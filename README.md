# buidtag

`buildtag` prints buildtags which are used in specified packages.

## Install

```
$ go install github.com/gostaticanalysis/buildtag/cmd/buildtag@latest
```

## Usage

```
$ buildtag golang.org/x/net/context
!go1.7 /Users/tenntenn/go/pkg/mod/golang.org/x/net@v0.0.0-20220107192237-5cfca573fb4d/context/context_test.go
go1.7 /Users/tenntenn/go/pkg/mod/golang.org/x/net@v0.0.0-20220107192237-5cfca573fb4d/context/go17.go
go1.9 /Users/tenntenn/go/pkg/mod/golang.org/x/net@v0.0.0-20220107192237-5cfca573fb4d/context/go19.go
!go1.7 /Users/tenntenn/go/pkg/mod/golang.org/x/net@v0.0.0-20220107192237-5cfca573fb4d/context/pre_go17.go
!go1.9 /Users/tenntenn/go/pkg/mod/golang.org/x/net@v0.0.0-20220107192237-5cfca573fb4d/context/pre_go19.go
```

```
$ buildtag -tag "^go1\.*" golang.org/x/net/context
/Users/tenntenn/go/pkg/mod/golang.org/x/net@v0.0.0-20220111093109-d55c255bac03/context/go17.go:go1.7
/Users/tenntenn/go/pkg/mod/golang.org/x/net@v0.0.0-20220111093109-d55c255bac03/context/go19.go:go1.9
```

```
$ buildtag -f="{{.Constraint}}" golang.org/x/net/context | sort | uniq
!go1.7
!go1.9
go1.7
go1.9
```
