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
$ buildtag -tags go1.7 golang.org/x/net/context
go1.7 /Users/tenntenn/go/pkg/mod/golang.org/x/net@v0.0.0-20220107192237-5cfca573fb4d/context/go17.go
!go1.9 /Users/tenntenn/go/pkg/mod/golang.org/x/net@v0.0.0-20220107192237-5cfca573fb4d/context/pre_go19.go
```

```
$ buildtag -f="{{.Constraint}}" golang.org/x/net/context | sort | uniq
!go1.7
!go1.9
go1.7
go1.9
```
