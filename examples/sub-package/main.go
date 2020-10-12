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
