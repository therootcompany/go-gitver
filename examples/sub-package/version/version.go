//go:generate go run -mod=vendor git.rootprojects.org/root/go-gitver --package version

package version

var (
	commit  = "0000000"
	version = "0.0.0-pre0+0000000"
	date    = "0000-00-00T00:00:00+0000"
)
