//go:generate go run -mod=vendor git.rootprojects.org/root/go-gitver

package main

import (
	"bytes"
	"fmt"
	"go/format"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"
)

var exitCode int
var exactVer *regexp.Regexp
var gitVer *regexp.Regexp
var verFile = "generated-version.go"

var (
	GitRev       = "0000000"
	GitVersion   = "v0.0.0-pre0+g0000000"
	GitTimestamp = "0000-00-00T00:00:00+0000"
)

func init() {
	// exactly vX.Y.Z (go-compatible semver)
	exactVer = regexp.MustCompile(`^v\d+\.\d+\.\d+$`)

	// vX.Y.Z-n-g0000000 git post-release, semver prerelease
	// vX.Y.Z-dirty git post-release, semver prerelease
	gitVer = regexp.MustCompile(`^(v\d+\.\d+)\.(\d+)(-(\d+))?(-(g[0-9a-f]+))?(-(dirty))?`)
}

func main() {
	args := os.Args[1:]
	for i := range args {
		arg := args[i]
		if "-f" == arg || "--fail" == arg {
			exitCode = 1
		} else if "-V" == arg || "version" == arg || "-version" == arg || "--version" == arg {
			fmt.Println(GitRev)
			fmt.Println(GitVersion)
			fmt.Println(GitTimestamp)
			os.Exit(0)
		}
	}
	if "" != os.Getenv("GITVER_FAIL") && "false" != os.Getenv("GITVER_FAIL") {
		exitCode = 1
	}

	desc, err := gitDesc()
	if nil != err {
		log.Fatalf("Failed to get git version: %s\n", err)
		os.Exit(exitCode)
	}
	rev := gitRev()
	ver := semVer(desc)
	ts, err := gitTimestamp(desc)
	if nil != err {
		ts = time.Now()
	}

	v := struct {
		Timestamp string
		Version   string
		GitRev    string
	}{
		Timestamp: ts.Format(time.RFC3339),
		Version:   ver,
		GitRev:    rev,
	}

	// Create or overwrite the go file from template
	var buf bytes.Buffer
	if err := versionTpl.Execute(&buf, v); nil != err {
		panic(err)
	}

	// Format
	src, err := format.Source(buf.Bytes())
	if nil != err {
		panic(err)
	}

	// Write to disk (in the Current Working Directory)
	f, err := os.Create(verFile)
	if nil != err {
		panic(err)
	}
	if _, err := f.Write(src); nil != err {
		panic(err)
	}
	if err := f.Close(); nil != err {
		panic(err)
	}
}

func gitDesc() (string, error) {
	args := strings.Split("git describe --tags --dirty --always", " ")
	cmd := exec.Command(args[0], args[1:]...)
	out, err := cmd.CombinedOutput()
	if nil != err {
		// Don't panic, just carry on
		//out = []byte("v0.0.0-0-g0000000")
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func gitRev() string {
	args := strings.Split("git rev-parse HEAD", " ")
	cmd := exec.Command(args[0], args[1:]...)
	out, err := cmd.CombinedOutput()
	if nil != err {
		fmt.Fprintf(os.Stderr,
			"\nUnexpected Error\n\n"+
				"Please open an issue at https://git.rootprojects.org/root/go-gitver/issues/new \n"+
				"Please include the following:\n\n"+
				"Command: %s\n"+
				"Output: %s\n"+
				"Error: %s\n"+
				"\nPlease and Thank You.\n\n", strings.Join(args, " "), out, err)
		os.Exit(exitCode)
	}
	return strings.TrimSpace(string(out))
}

func semVer(desc string) string {
	if exactVer.MatchString(desc) {
		// v1.0.0
		return desc
	}

	if !gitVer.MatchString(desc) {
		return ""
	}

	// (v1.0).(0)(-(1))(-(g0000000))(-(dirty))
	vers := gitVer.FindStringSubmatch(desc)
	patch, err := strconv.Atoi(vers[2])
	if nil != err {
		fmt.Fprintf(os.Stderr,
			"\nUnexpected Error\n\n"+
				"Please open an issue at https://git.rootprojects.org/root/go-gitver/issues/new \n"+
				"Please include the following:\n\n"+
				"git description: %s\n"+
				"RegExp: %#v\n"+
				"Error: %s\n"+
				"\nPlease and Thank You.\n\n", desc, gitVer, err)
		os.Exit(exitCode)
	}

	// v1.0.1-pre1
	// v1.0.1-pre1+g0000000
	// v1.0.1-pre0+dirty
	// v1.0.1-pre0+g0000000-dirty
	if "" == vers[4] {
		vers[4] = "0"
	}
	ver := fmt.Sprintf("%s.%d-pre%s", vers[1], patch+1, vers[4])
	if "" != vers[6] || "dirty" == vers[8] {
		ver += "+"
		if "" != vers[6] {
			ver += vers[6]
			if "" != vers[8] {
				ver += "-"
			}
		}
		ver += vers[8]
	}

	return ver
}

func gitTimestamp(desc string) (time.Time, error) {
	args := []string{
		"git",
		"show", desc,
		"--format=%cd",
		"--date=format:%Y-%m-%dT%H:%M:%SZ%z",
		"--no-patch",
	}
	cmd := exec.Command(args[0], args[1:]...)
	out, err := cmd.CombinedOutput()
	if nil != err {
		// a dirty desc was probably used
		return time.Time{}, err
	}
	return time.Parse(time.RFC3339, strings.TrimSpace(string(out)))
}

var versionTpl = template.Must(template.New("").Parse(`// Code generated by go generate; DO NOT EDIT.
package main

func init() {
	GitRev = "{{ .GitRev }}"
	if "" != "{{ .Version }}" {
		GitVersion = "{{ .Version }}"
	}
	GitTimestamp = "{{ .Timestamp }}"
}
`))
