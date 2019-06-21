# Example

Prints the version or a nice message

# Build

Typically the developer would perform these steps
and then commit the results (`go.mod`, `go.sum`, `vendor`).

However, since this is an example within the project directory,
that seemed a little redundant.

```bash
go mod tidy
go mod vendor
```

These are the instructions that someone cloning the repo might use.

```bash
go generate -mod=vendor ./...
go build -mod=vendor -o hello *.go
./hello
./hello --version
```

Note: If the source is distributed in a non-git tarball then
`generated-version.go` will not be output, and whatever
version info is in `package main` will remain as-is.

If you would prefer the build process to fail (i.e. in a CI/CD pipeline),
you can set the environment variable `GITVER_FAIL=true`.
