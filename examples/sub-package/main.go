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
		fmt.Println(version.GitRev)
		fmt.Println(version.GitVersion)
		fmt.Println(version.GitTimestamp)
		return
	}

	fmt.Println("Hello, World!")
}
