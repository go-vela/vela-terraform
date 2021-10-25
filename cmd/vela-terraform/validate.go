// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
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

const validationAction = "validate"

// Validation represents the plugin configuration for validate information.
type Validation struct {
	//  If set to true (default), the command will check whether all required variables have been specified. i.e. "-check-variables=true"
	CheckVariables bool
	// terraform file or directory to apply
	Directory string
	// if specified, output won't contain any color. i.e. "-no-color"
	NoColor bool
	// set a variable in the Terraform configuration. i.e. "-var 'foo=bar'"
	Vars []string
	// set variables in the Terraform configuration from a file. i.e. "-var-file=foo"
	VarFiles []string
	Version  *semver.Version
}

// Command formats and outputs the Validate command from
// the provided configuration to validate to resources.
func (v *Validation) Command() *exec.Cmd {
	logrus.Trace("creating terraform validate command from plugin configuration")

	// global Variables
	var globalFlags []string

	// variable to store flags for command
	var flags []string

	// check if Directory is provided
	if v.Directory != "." && SupportsChdir(v.Version) {
		globalFlags = append(flags, fmt.Sprintf("-chdir=%s", v.Directory))
	}

	// check if CheckVariables is provided
	if v.CheckVariables {
		// add flag for CheckVariables from provided validate command
		flags = append(flags, fmt.Sprintf("-check-variables=%t", v.CheckVariables))
	}

	// check if NoColor is provided
	if v.NoColor {
		// add flag for NoColor from provided validate command
		flags = append(flags, "-no-color")
	}

	// check if Vars is provided
	if len(v.Vars) > 0 {
		for _, v := range v.Vars {
			// add flag for Vars from provided command
			flags = append(flags, fmt.Sprintf(`-var=%s`, v))
		}
	}

	// check if VarFiles is provided
	if len(v.VarFiles) > 0 {
		for _, v := range v.VarFiles {
			// add flag for VarFiles from provided command
			flags = append(flags, fmt.Sprintf(`-var-file=%s`, v))
		}
	}

	// check if Directory is provided and terraform version doesn't support chdir
	if v.Directory != "." && !SupportsChdir(v.Version) {
		flags = append(flags, v.Directory)
	}

	globalFlags = append(globalFlags, validationAction)

	return exec.Command(_terraform, append(globalFlags, flags...)...)
}

// Exec formats and runs the commands for validating Terraform.
func (v *Validation) Exec() error {
	logrus.Trace("running validate with provided configuration")

	// create the validate command for the file
	cmd := v.Command()

	// run the validate command for the file
	err := execCmd(cmd)
	if err != nil {
		return err
	}

	return nil
}

// Validate verifies the Delete is properly configured.
func (v *Validation) Validate() error {
	logrus.Trace("validating plan plugin configuration")

	if strings.EqualFold(v.Directory, ".") {
		logrus.Warn("terraform validate will run in current dir")
	}

	return nil
}
