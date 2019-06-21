# Example

Prints the version or a nice message

# Doesn't have a separate tools package

This is just like `examples/basic`,
but it uses a normal file with a build tag
rather than a tools package.

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
