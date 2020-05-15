// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const planAction = "plan"

// Plan represents the plugin configuration for plan information.
type Plan struct {
	// terraform file or directory to plan
	Directory string
	// If set, a plan will be generated to destroy all resources managed by the given configuration and state. i.e. "-destroy"
	Destroy bool
	// Return detailed exit codes when the command exits. i.e. "-detailed-exitcode"
	DetailedExitCode bool
	// ask for input for variables if not directly set. i.e. "-input=true"
	Input bool
	// the state file when locking is supported. i.e. -lock=true
	Lock bool
	// duration to retry a state lock. i.e. "-lock-timeout=0s"
	LockTimeout time.Duration
	// Specifies the depth of modules to show in the output. i.e. "-module-depth=n"
	ModuleDepth int
	// if specified, output won't contain any color. i.e. "-no-color"
	NoColor bool
	// Write a plan file to the given path. i.e. "-out=path"
	Out string
	// limit the number of parallel resource operations. i.e. "-parallelism=n"
	Parallelism int
	// update state prior to checking for differences. i.e. "-refresh=true"
	Refresh bool
	// path to read and save state (unless state-out is specified). i.e. "-state=path"
	State string
	// resource to target. i.e. "-target=resource"
	Target string
	// set a variable in the Terraform configuration. i.e. "-var 'foo=bar'"
	Var []string
	// set variables in the Terraform configuration from a file. i.e. "-var-file=foo"
	VarFile []string
}

// Command formats and outputs the Plan command from
// the provided configuration to plan to resources.
func (p *Plan) Command(dir string) *exec.Cmd {
	logrus.Trace("creating terraform plan command from plugin configuration")

	// variable to store flags for command
	var flags []string

	// check if Destroy is provided
	if p.Destroy {
		// add flag for Destroy from provided plan command
		flags = append(flags, "-destroy")
	}

	// check if DetailedExitCode is provided
	if p.DetailedExitCode {
		// add flag for DetailedExitCode from provided plan command
		flags = append(flags, "-detailed-exitcode")
	}

	// check if Input is provided
	if p.Input {
		// add flag for Input from provided plan command
		flags = append(flags, "-input=true")
	}

	// check if Lock is provided
	if p.Lock {
		// add flag for Lock from provided plan command
		flags = append(flags, "-lock=true")
	}

	// check if LockTimeout is provided
	if p.LockTimeout > 0 {
		// add flag for LockTimeout from provided plan command
		flags = append(flags, fmt.Sprintf("-lock-timeout=%s", p.LockTimeout))
	}

	// check if ModuleDepth is provided
	if p.ModuleDepth > 0 {
		// add flag for ModuleDepth from provided plan command
		flags = append(flags, fmt.Sprintf("-module-depth=%d", p.ModuleDepth))
	}

	// check if NoColor is provided
	if p.NoColor {
		// add flag for NoColor from provided plan command
		flags = append(flags, "-no-color")
	}

	// check if Out is provided
	if len(p.Out) > 0 {
		// add flag for Out from provided plan command
		flags = append(flags, fmt.Sprintf("-out=%s", p.Out))
	}

	// check if Parallelism is provided
	if p.Parallelism > 0 {
		// add flag for Parallelism from provided plan command
		flags = append(flags, fmt.Sprintf("-parallelism=%d", p.Parallelism))
	}

	// check if Refresh is provided
	if p.Refresh {
		// add flag for Refresh from provided plan command
		flags = append(flags, "-refresh=true")
	}

	// check if State is provided
	if len(p.State) > 0 {
		// add flag for State from provided plan command
		flags = append(flags, fmt.Sprintf("-state=%s", p.State))
	}

	// check if Target is provided
	if len(p.Target) > 0 {
		// add flag for Target from provided plan command
		flags = append(flags, fmt.Sprintf("-target=%s", p.Target))
	}

	// check if Var is provided
	if len(p.Var) > 0 {
		var vars string
		for _, v := range p.Var {
			vars += fmt.Sprintf(" %s", v)
		}

		// add flag for Var from provided plan command
		flags = append(flags, fmt.Sprintf("-var=\"%s\"", strings.TrimPrefix(vars, " ")))
	}

	// check if Var is provided
	if len(p.VarFile) > 0 {
		var files string
		for _, v := range p.VarFile {
			files += fmt.Sprintf("-var-file=%s ", v)
		}

		// add flag for Var from provided plan command
		flags = append(flags, strings.TrimPrefix(files, " "))
	}

	// add the required dir param
	flags = append(flags, dir)

	return exec.Command(_terraform, append([]string{planAction}, flags...)...)
}

// Exec formats and runs the commands for planning Terraform.
func (p *Plan) Exec() error {
	logrus.Trace("running plan with provided configuration")

	// create the plan command for the file
	cmd := p.Command(p.Directory)

	// run the plan command for the file
	err := execCmd(cmd)
	if err != nil {
		return err
	}

	return nil
}

// Validate verifies the Delete is properly configured.
func (p *Plan) Validate() error {
	logrus.Trace("validating plan plugin configuration")

	if len(p.Directory) == 0 {
		logrus.Warn("terrafrom plan will run in current dir")

		// set the directory to run in current dir
		p.Directory = "."
	}

	return nil
}
