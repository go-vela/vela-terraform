// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/Masterminds/semver/v3"

	"github.com/sirupsen/logrus"
)

const fmtAction = "fmt"

// FMT represents the plugin configuration for fmt information.
type FMT struct {
	// Check if the input is formatted
	Check bool
	// Display diffs of formatting changes
	Diff bool
	// terraform file or directory to fmt
	Directory string
	// List files whose formatting differs
	List bool
	// Write result to source file instead of STDOUT
	Write   bool
	Version *semver.Version
}

// Command formats and outputs the FMT command from
// the provided configuration to fmt to resources.
func (f *FMT) Command() *exec.Cmd {
	logrus.Trace("creating terraform fmt command from plugin configuration")

	// global Variables
	var globalFlags []string

	// variable to store flags for command
	var flags []string

	// check if Directory is provided
	if f.Directory != "." && SupportsChdir(f.Version) {
		globalFlags = append(flags, fmt.Sprintf("-chdir=%s", f.Directory))
	}

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

	// check if Directory is provided and terraform version doesn't support chdir
	if f.Directory != "." && !SupportsChdir(f.Version) {
		flags = append(flags, f.Directory)
	}

	globalFlags = append(globalFlags, fmtAction)
	return exec.Command(_terraform, append(globalFlags, flags...)...)
}

// Exec formats and runs the commands for formatting Terraform files.
func (f *FMT) Exec() error {
	logrus.Trace("running fmt with provided configuration")

	// create the fmt command for the file
	cmd := f.Command()

	// run the fmt command for the file
	err := execCmd(cmd)
	if err != nil {
		return err
	}

	return nil
}

// Validate verifies the Delete is properly configured.
func (f *FMT) Validate() error {
	logrus.Trace("validating plan plugin configuration")

	if strings.EqualFold(f.Directory, ".") {
		logrus.Warn("terraform fmt will run in current dir")
	}

	return nil
}
