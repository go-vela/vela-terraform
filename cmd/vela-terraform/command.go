// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

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

	// name of cmd
	name := "terraform"

	// default arg
	args := []string{
		"get",
	}

	// set working directory if provided
	if len(dir) > 0 {
		args = append(args, fmt.Sprintf("-chdir=%s", dir))
	}

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
