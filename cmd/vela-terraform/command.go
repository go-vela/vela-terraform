// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
)

// execCmd is a helper function to
// run the provided command.
func execCmd(e *exec.Cmd) error {
	logrus.Tracef("executing cmd %s", strings.Join(e.Args, " "))

	// set command stdout to OS stdout
	e.Stdout = os.Stdout
	// set command stderr to OS stderr
	e.Stderr = os.Stderr

	// output "trace" string for command
	fmt.Println("$", strings.Join(e.Args, " "))

	return e.Run()
}

// getCmd is a helper function to retrieve
// the terraform modules required for the files.
func getCmd(dir string) *exec.Cmd {
	// default cmd and arg
	name := "terraform"

	var args []string

	// check if Directory is provided
	if dir != "." {
		args = append(args, fmt.Sprintf("-chdir=%s", dir))
	}

	args = append(args, "get")

	return exec.Command(
		name,
		args...,
	)
}

// versionCmd is a helper function to output
// the client and server version information.
func versionCmd() *exec.Cmd {
	logrus.Trace("creating terraform version command")

	// variable to store flags for command
	var flags []string

	// add flag for version kubectl command
	flags = append(flags, "version")

	return exec.Command(_terraform, flags...)
}
