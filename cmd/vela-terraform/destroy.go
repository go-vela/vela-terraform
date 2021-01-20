// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
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

const destroyAction = "destroy"

// Destroy represents the plugin configuration for destroy information.
type Destroy struct {
	// skip interactive approval of plan before applying. i.e. "-auto-approve"
	AutoApprove bool
	// path to backup the existing state file before modifying. i.e. "-backup=path "
	Backup string
	// terraform file or directory to destroy
	Directory string
	// the state file when locking is supported. i.e. -lock=true
	Lock bool
	// duration to retry a state lock. i.e. "-lock-timeout=0s"
	LockTimeout time.Duration
	// if specified, output won't contain any color. i.e. "-no-color"
	NoColor bool
	// limit the number of parallel resource operations. i.e. "-parallelism=n"
	Parallelism int
	// update state prior to checking for differences. i.e. "-refresh=true"
	Refresh bool
	// path to read and save state (unless state-out is specified). i.e. "-state=path"
	State string
	// path to write state to that is different than state. i.e. "-state-out=path"
	StateOut string
	// resource to target. i.e. "-target=resource"
	Target string
	// set a variable in the Terraform configuration. i.e. "-var 'foo=bar'"
	Vars []string
	// set variables in the Terraform configuration from a file. i.e. "-var-file=foo"
	VarFiles []string
}

// Command formats and outputs the Destroy command from
// the provided configuration to destroy to resources.
func (a *Destroy) Command(dir string) *exec.Cmd {
	logrus.Trace("creating terraform destroy command from plugin configuration")

	// variable to store flags for command
	var flags []string

	// check if AutoApprove is provided
	if a.AutoApprove {
		// add flag for AutoApprove from provided destroy command
		flags = append(flags, "-auto-approve")
	}

	// check if Backup is provided
	if len(a.Backup) > 0 {
		// add flag for Backup from provided destroy command
		flags = append(flags, fmt.Sprintf("-backup=%s", a.Backup))
	}

	// check if Lock is provided
	if a.Lock {
		// add flag for Lock from provided destroy command
		flags = append(flags, "-lock=true")
	}

	// check if LockTimeout is provided
	if a.LockTimeout > 0 {
		// add flag for LockTimeout from provided destroy command
		flags = append(flags, fmt.Sprintf("-lock-timeout=%s", a.LockTimeout))
	}

	// check if NoColor is provided
	if a.NoColor {
		// add flag for NoColor from provided destroy command
		flags = append(flags, "-no-color")
	}

	// check if Parallelism is provided
	if a.Parallelism > 0 {
		// add flag for Parallelism from provided destroy command
		flags = append(flags, fmt.Sprintf("-parallelism=%d", a.Parallelism))
	}

	// check if Refresh is provided
	if a.Refresh {
		// add flag for Refresh from provided destroy command
		flags = append(flags, "-refresh=true")
	}

	// check if State is provided
	if len(a.State) > 0 {
		// add flag for State from provided destroy command
		flags = append(flags, fmt.Sprintf("-state=%s", a.State))
	}

	// check if StateOut is provided
	if len(a.StateOut) > 0 {
		// add flag for StateOut from provided destroy command
		flags = append(flags, fmt.Sprintf("-state-out=%s", a.StateOut))
	}

	// check if Target is provided
	if len(a.Target) > 0 {
		// add flag for Target from provided destroy command
		flags = append(flags, fmt.Sprintf("-target=%s", a.Target))
	}

	// check if Vars is provided
	if len(a.Vars) > 0 {
		var vars string
		for _, v := range a.Vars {
			vars += fmt.Sprintf(" %s", v)
		}
		// add flag for Vars from provided destroy command
		flags = append(flags, fmt.Sprintf("-var=\"%s\"", strings.TrimPrefix(vars, " ")))
	}

	// check if VarFiles is provided
	if len(a.VarFiles) > 0 {
		var files string
		for _, v := range a.VarFiles {
			files += fmt.Sprintf("-var-file=%s ", v)
		}

		// add flag for VarFiles from provided destroy command
		flags = append(flags, strings.TrimSuffix(files, " "))
	}

	// add the required dir param
	flags = append(flags, dir)

	return exec.Command(_terraform, append([]string{destroyAction}, flags...)...)
}

// Exec formats and runs the commands for destroying resources with Terraform.
func (d *Destroy) Exec() error {
	logrus.Trace("running destroy with provided configuration")

	// create the destroy command for the file
	cmd := d.Command(d.Directory)

	// run the destroy command for the file
	err := execCmd(cmd)
	if err != nil {
		return err
	}

	return nil
}

// Validate verifies the Delete is properly configured.
func (d *Destroy) Validate() error {
	logrus.Trace("validating plan plugin configuration")

	if strings.EqualFold(d.Directory, ".") {
		logrus.Warn("terrafrom destroy will run in current dir")
	}

	return nil
}
