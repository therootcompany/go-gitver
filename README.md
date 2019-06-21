# git-version.go

Use git tags to add semver to your go package.

>     Goal: Either use an exact version like v1.0.0
>           or translate the git version like v1.0.0-4-g0000000
>           to a semver like v1.0.1-pre4+g0000000
>
>           Fail gracefully when git repo isn't available.

# Demo

```bash
go run git.rootprojects.org/root/go-gitver
```

# QuickStart

Add this to the top of your main file:

```go
//go:generate go run -mod=vendor git.rootprojects.org/root/go-gitver

```

Add a file that imports go-gitver (for versioning)

```go
// build +tools

package example

import _ "git.rootprojects.org/root/go-gitver"
```

Change you build instructions to be something like this:

```bash
go mod vendor
go generate -mod=vendor ./...
go build -mod=vendor -o example cmd/example/*.go
```

You don't have to use `mod vendor`, but I highly recommend it.

# Options

```
version   print version and exit
--fail    will cause non-zero exit status on failure
```

ENVs

```
# Alias for --fail
GITVER_FAIL=true
```

For example:

```go
//go:generate go run -mod=vendor git.rootprojects.org/root/go-gitver --fail

```

```bash
go run -mod=vendor git.rootprojects.org/root/go-gitver version
```

# Usage

See `examples/basic`

1. Create a `tools` package in your project
2. Guard it against regular builds with `// build +tools`
3. Include `_ "git.rootprojects.org/root/go-gitver"` in the imports
4. Declare `var GitRev, GitVersion, GitTimestamp string` in your `package main`
5. Include `//go:generate go run -mod=vendor git.rootprojects.org/root/go-gitver` as well

`tools/tools.go`:

```go
// build +tools

// This is a dummy package for build tooling
package tools

import (
	_ "git.rootprojects.org/root/go-gitver"
)
```

`main.go`:

```go
//go:generate go run git.rootprojects.org/root/go-gitver --fail

package main

import "fmt"

var (
	GitRev       = "0000000"
	GitVersion   = "v0.0.0-pre0+0000000"
	GitTimestamp = "0000-00-00T00:00:00+0000"
)

func main() {
  fmt.Println(GitRev)
  fmt.Println(GitVersion)
  fmt.Println(GitTimestamp)
}
```

If you're using `go mod vendor` (which I highly recommend that you do),
you'd modify the `go:generate` ever so slightly:

```go
//go:generate go run -mod=vendor git.rootprojects.org/root/go-gitver --fail
```

The only reason I didn't do that in the example is that I'd be included
the repository in itself and that would be... weird.

# Why a tools package?

>     import "git.rootprojects.org/root/go-gitver" is a program, not an importable package

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
