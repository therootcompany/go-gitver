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
		fmt.Println(version.commit)
		fmt.Println(version.version)
		fmt.Println(version.date)
		return
	}

	fmt.Println("Hello, World!")
}
