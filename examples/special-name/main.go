//go:generate go run -mod=vendor git.rootprojects.org/root/go-gitver/v2 --package version --outfile ./version/zversion.go

package main

import (
	"flag"
	"fmt"

	"example.com/hello/version"
)

func main() {
	showVersion := flag.Bool("version", false, "Print version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Println(version.Commit())
		fmt.Println(version.Version())
		fmt.Println(version.Date())
		return
	}

	fmt.Println("Hello, World!")
}
