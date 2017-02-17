package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"text/template"

	"github.com/kardianos/osext"
	"github.com/spf13/pflag"
)

var execDir string

func init() {
	dir, err := osext.ExecutableFolder()
	if err != nil {
		fmt.Printf("Exec dir error: %s\n", err)
		os.Exit(1)
	}
	execDir = dir
}

var usage = template.Must(template.New("usage").Parse(`Usage:
	ww command [...args]

Available Commands:
{{- range $i, $cmd := . }}
	{{ $cmd }}
{{- end }}
`))

func main() {
	pflag.SetInterspersed(false)
	pflag.Usage = func() {
		usage.Execute(os.Stderr, allCmds())
	}
	pflag.Parse()

	args := pflag.Args()
	fmt.Printf("%#v\n", args)
	if len(args) == 0 {
		pflag.Usage()
		os.Exit(1)
	}

	cmd := filepath.Join(execDir, "ww-"+args[0])
	if _, err := os.Stat(cmd); err != nil {
		fmt.Printf("Command read error: %s\n", err)
		os.Exit(1)
	}

	args = args[1:]

	ex := exec.Command(cmd, args...)
	ex.Stdin = os.Stdin
	ex.Stdout = os.Stdout
	ex.Stderr = os.Stderr

	if err := ex.Run(); err != nil {
		fmt.Printf("working wheatley run error: %s\n", err)
		os.Exit(exitCode(err))
	}
}

func allCmds() []string {
	fis, err := ioutil.ReadDir(execDir)
	if err != nil {
		fmt.Printf("error reading dir %s: %s\n", execDir, err)
		os.Exit(1)
	}

	cmds := []string{}

	for _, fi := range fis {
		if strings.HasPrefix(fi.Name(), "ww-") {
			cmds = append(cmds, fi.Name()[3:])
		}
	}

	return cmds
}

// exitCode attempts to return the exit code from an
// exec error
// since there is no portable way to do this, we assume
// a unix system
func exitCode(err error) int {
	if err == nil {
		return -1
	}

	osErr, ok := err.(*exec.ExitError)

	if !ok {
		return -1
	}

	// assume unix system
	return osErr.Sys().(syscall.WaitStatus).ExitStatus()
}
