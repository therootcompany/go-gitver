# Example

Prints the version or a nice message

# Doesn't have a separate tools package

This is just like `examples/sub-package` except that its `//go:generate` is in `main.go`
and it outputs `./version/zversion.go` instead of `xversion.go`.

See `examples/basic` for more details.

# Demo

```bash
go mod tidy
go mod vendor
```

```bash
go generate -mod=vendor ./...
go build -mod=vendor -o hello *.go
./hello
./hello --version
```
