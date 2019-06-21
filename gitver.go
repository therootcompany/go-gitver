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

func init() {
	// exactly vX.Y.Z (go-compatible semver)
	exactVer = regexp.MustCompile(`^v\d+\.\d+\.\d+$`)

	// vX.Y.Z-n-g0000000 git post-release, semver prerelease
	// vX.Y.Z-dirty git post-release, semver prerelease
	gitVer = regexp.MustCompile(`^(v\d+\.\d+)\.(\d+)-((\d+)-)?(dirty|g)`)
}

func main() {
	args := os.Args[1:]
	for i := range args {
		arg := args[i]
		if "-f" == arg || "--fail" == arg {
			exitCode = 1
		}
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
		fmt.Println("badtimes", err) // TODO remove
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
	var ver string
	if exactVer.MatchString(desc) {
		// v1.0.0
		ver = desc
	} else if gitVer.MatchString(desc) {
		// ((v1.0).(0)-(1))
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
		ver = fmt.Sprintf("%s.%d-pre%s", vers[1], patch+1, vers[4])
	}
	return ver
}

func gitTimestamp(desc string) (time.Time, error) {
	args := strings.Split(fmt.Sprintf("git show %s --format=%%cd --date=format:%%Y-%%m-%%dT%%H:%%M:%%SZ%%z --no-patch", desc), " ")
	cmd := exec.Command(args[0], args[1:]...)
	out, err := cmd.CombinedOutput()
	fmt.Printf("args:\n%#v\n\n", args)
	fmt.Printf("in:\n%s\n\n", strings.Join(args, " "))
	fmt.Println("out:\n\n", string(out), err)
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
	GitVersion = "{{ .Version }}"
	GitTimestamp = "{{ .Timestamp }}"
}
`))
