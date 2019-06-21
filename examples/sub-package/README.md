# Example

Prints the version or a nice message

# Put version in its own packge

This is just like `examples/basic`,
but it uses the package name `version` instead of `main`.

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
