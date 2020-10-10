//go:generate go run -mod=vendor git.rootprojects.org/root/go-gitver --fail

package main

import (
	"flag"
	"fmt"
)

var (
	commit  = "0000000"
	version = "0.0.0-pre0+0000000"
	date    = "0000-00-00T00:00:00+0000"
)

func main() {
	showVersion := flag.Bool("version", false, "Print version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Println(commit)
		fmt.Println(version)
		fmt.Println(date)
		return
	}

	fmt.Println("Hello, World!")
}
