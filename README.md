# [Go GitVer](https://git.rootprojects.org/root/go-gitver)

Use **git tags** to add (GoReleaser-compatible) [**semver**](https://semver.org/)
to your go package in under 150
[lines of code](https://git.rootprojects.org/root/go-gitver/src/branch/master/gitver/gitver.go).

```txt
Goals:

      1. Use an exact `git tag` version, like v1.0.0, when clean
      2. Translate the `git describe` version  (v1.0.0-4-g0000000)
	     to semver (1.0.1-pre4+g0000000) in between releases
      3. Note when `dirty` (and have build timestamp)

      Fail gracefully when git repo isn't available.
```

# GoDoc

See <https://pkg.go.dev/git.rootprojects.org/root/go-gitver/v2>.

# How it works

1. You define the fallback version and version printing in `main.go`:

```go
//go:generate go run git.rootprojects.org/root/go-gitver/v2

package main

import (
	"fmt"
	"strings"
)

var (
	commit  = "0000000"
	version = "0.0.0-pre0+0000000"
	date    = "0000-00-00T00:00:00+0000"
)

func main() {
	if (len(os.Args) > 1 && "version" == strings.TrimLeft(os.Args[1], "-")) {
		fmt.Printf("Foobar v%s (%s) %s\n", version, commit[:7], date)
	}
	// ...
}
```

2. You `go generate` or `go run git.rootprojects.org/root/go-gitver/v2` to generate `xversion.go`:

```go
package main

func init() {
    commit  = "0921ed1e"
    version = "1.1.2"
    date    = "2019-07-01T02:32:58-06:00"
}
```

# Demo

Generate an `xversion.go` file:

```bash
go run git.rootprojects.org/root/go-gitver/v2
cat xversion.go
```

```go
// Code generated by go generate; DO NOT EDIT.
package main

func init() {
	commit  = "6dace8255b52e123297a44629bc32c015add310a"
	version = "1.1.4-pre2+g6dace82"
	date    = "2020-07-16T20:48:15-06:00"
}
```

<small>**Note**: The file is named `xversion.go` by default so that the
generated file's `init()` will come later, and thus take priority, over
most other files.</small>

See `go-gitver`s self-generated version:

```bash
go run git.rootprojects.org/root/go-gitver/v2 version
```

```txt
6dace8255b52e123297a44629bc32c015add310a
v1.1.4-pre2+g6dace82
2020-07-16T20:48:15-06:00
```

# QuickStart

Add this to the top of your main file, so that it runs with `go generate`:

```go
//go:generate go run -mod=vendor git.rootprojects.org/root/go-gitver/v2

```

Add a file that imports go-gitver (for versioning)

```go
// +build tools

package example

import _ "git.rootprojects.org/root/go-gitver/v2"
```

Change you build instructions to be something like this:

```bash
go mod vendor
go generate -mod=vendor ./...
go build -mod=vendor -o example cmd/example/*.go
```

You don't have to use `-mod=vendor`, but I highly recommend it (just `go mod tidy; go mod vendor` to start).

# Options

```txt
version           print version and exit
--fail            exit with non-zero status code on failure
--package <name>  will set the package name
--outfile <name>  will replace `xversion.go` with the given file path
```

ENVs

```bash
# Alias for --fail
GITVER_FAIL=true
```

For example:

```go
//go:generate go run -mod=vendor git.rootprojects.org/root/go-gitver/v2 --fail

```

```bash
go run -mod=vendor git.rootprojects.org/root/go-gitver/v2 version
```

# Usage

See `examples/basic`

1. Create a `tools` package in your project
2. Guard it against regular builds with `// +build tools`
3. Include `_ "git.rootprojects.org/root/go-gitver/v2"` in the imports
4. Declare `var commit, version, date string` in your `package main`
5. Include `//go:generate go run -mod=vendor git.rootprojects.org/root/go-gitver/v2` as well

`tools/tools.go`:

```go
// +build tools

// This is a dummy package for build tooling
package tools

import (
	_ "git.rootprojects.org/root/go-gitver/v2"
)
```

`main.go`:

```go
//go:generate go run git.rootprojects.org/root/go-gitver/v2 --fail

package main

import "fmt"

var (
	commit  = "0000000"
	version = "0.0.0-pre0+0000000"
	date    = "0000-00-00T00:00:00+0000"
)

func main() {
  fmt.Println(commit)
  fmt.Println(version)
  fmt.Println(date)
}
```

If you're using `go mod vendor` (which I highly recommend that you do),
you'd modify the `go:generate` ever so slightly:

```go
//go:generate go run -mod=vendor git.rootprojects.org/root/go-gitver/v2 --fail
```

The only reason I didn't do that in the example is that I'd be included
the repository in itself and that would be... weird.

# Why a tools package?

> import "git.rootprojects.org/root/go-gitver/v2" is a program, not an importable package

Having a tools package with a build tag that you don't use is a nice way to add exact
versions of a command package used for tooling to your `go.mod` with `go mod tidy`,
without getting the error above.

# git: behind the curtain

These are the commands that are used under the hood to produce the versions.

Shows the git tag + description. Assumes that you're using the semver format `v1.0.0` for your base tags.

```bash
git describe --tags --dirty --always
# v1.0.0
# v1.0.0-1-g0000000
# v1.0.0-dirty
```

Show the commit date (when the commit made it into the current tree).
Internally we use the current date when the working tree is dirty.

```bash
git show v1.0.0-1-g0000000 --format=%cd --date=format:%Y-%m-%dT%H:%M:%SZ%z --no-patch
# 2010-01-01T20:30:00Z-0600
# fatal: ambiguous argument 'v1.0.0-1-g0000000-dirty': unknown revision or path not in the working tree.
```

Shows the most recent commit.

```bash
git rev-parse HEAD
# 0000000000000000000000000000000000000000
```

# Errors

### cannot find package "."

```txt
package git.rootprojects.org/root/go-gitver/v2: cannot find package "." in:
	/Users/me/go-example/vendor/git.rootprojects.org/root/go-gitver/v2
cmd/example/example.go:1: running "go": exit status 1
```

You forgot to update deps and re-vendor:

```bash
go mod tidy
go mod vendor
```
