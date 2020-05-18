// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"
	"os/exec"

	"github.com/sirupsen/logrus"
)

const fmtAction = "fmt"

// FMT represents the plugin configuration for fmt information.
type FMT struct {
	// terraform file or directory to fmt
	Directory string
	// List files whose formatting differs
	List bool
	// Write result to source file instead of STDOUT
	Write bool
	// Display diffs of formatting changes
	Diff bool
	// Check if the input is formatted
	Check bool
}

// Command formats and outputs the FMT command from
// the provided configuration to fmt to resources.
func (f *FMT) Command(dir string) *exec.Cmd {
	logrus.Trace("creating terraform fmt command from plugin configuration")

	// variable to store flags for command
	var flags []string

	// check if List is provided
	if !f.List {
		// add flag for List from provided fmt command
		flags = append(flags, fmt.Sprintf("-list=%t", f.List))
	}

	// check if Write is provided
	if !f.Write {
		// add flag for Write from provided fmt command
		flags = append(flags, fmt.Sprintf("-write=%t", f.Write))
	}

	// check if Diff is provided
	if f.Diff {
		// add flag for Diff from provided fmt command
		flags = append(flags, fmt.Sprintf("-diff=%t", f.Diff))
	}

	// check if Check is provided
	if f.Check {
		// add flag for Check from provided fmt command
		flags = append(flags, fmt.Sprintf("-check=%t", f.Check))
	}

	// add the required dir param
	flags = append(flags, dir)

	return exec.Command(_terraform, append([]string{fmtAction}, flags...)...)
}

// Exec formats and runs the commands for removing artifacts in Artifactory.
func (f *FMT) Exec() error {
	logrus.Trace("running fmt with provided configuration")

	// create the fmt command for the file
	cmd := f.Command(f.Directory)

	// run the fmt command for the file
	err := execCmd(cmd)
	if err != nil {
		return err
	}

	return nil
}

// Validate verifies the Delete is properly configured.
func (f *FMT) Validate() error {
	logrus.Trace("validating fmt plugin configuration")

	if len(f.Directory) == 0 {
		logrus.Warn("terrafrom fmt will run in current dir")

		// set the directory to run in current dir
		f.Directory = "."
	}

	return nil
}
