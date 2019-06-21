//go:generate go run -mod=vendor git.rootprojects.org/root/go-gitver

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
