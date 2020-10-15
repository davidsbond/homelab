package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func main() {
	version, err := getVersion()
	if err != nil {
		fmt.Println("failed to determine version", err)
		os.Exit(1)
	}

	binaries, err := getEntrypoints()
	if err != nil {
		fmt.Println("failed to determine entrypoints", err)
		os.Exit(1)
	}

	if len(binaries) == 0 {
		fmt.Println("no entrypoints found")
		os.Exit(0)
	}

	// Disable CGO.
	environ := append(os.Environ(), "CGO_ENABLED=0")
	for _, binary := range binaries {
		if err := compile(version, binary, environ); err != nil {
			fmt.Println("no entrypoints found")
			os.Exit(0)
		}
	}
}

func getEntrypoints() ([]string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// Return all `cmd/*/main.go` files to build.
	pattern := filepath.Join(wd, "cmd", "*", "main.go")
	return filepath.Glob(pattern)
}

func getVersion() (string, error) {
	buf := bytes.NewBuffer([]byte{})

	// Extract version from `git describe`
	cmd := exec.Command("git", "describe", "--tags", "--always")
	cmd.Stdout = buf
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	return strings.TrimSuffix(buf.String(), "\n"), err
}

func compile(version, entrypoint string, environ []string) error {
	dir := filepath.Dir(entrypoint)
	base := filepath.Base(dir)
	output := filepath.Join("bin", base)

	// If a `DESCRIPTION` file exists use it to set the application description.
	desc, err := ioutil.ReadFile(filepath.Join(dir, "DESCRIPTION"))
	switch {
	case errors.Is(err, os.ErrNotExist):
		desc = []byte("A command-line application")
	case err != nil:
		return err
	}

	ldFlags := map[string]string{
		"pkg.dsb.dev/environment.Version":                version,
		"pkg.dsb.dev/environment.compiled":               strconv.FormatInt(time.Now().Unix(), 10),
		"pkg.dsb.dev/environment.ApplicationName":        base,
		"pkg.dsb.dev/environment.ApplicationDescription": string(bytes.TrimSuffix(desc, []byte{'\n'})),
	}

	ldFlagsStr := "-s -w "
	for k, v := range ldFlags {
		ldFlagsStr += fmt.Sprintf(`-X "%s=%s" `, k, v)
	}

	cmd := exec.Command("go", "build",
		"-o", output,
		"-ldflags", ldFlagsStr,
		dir)

	cmd.Env = environ
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
