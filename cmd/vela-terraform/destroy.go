// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"

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
	Version  *semver.Version
}

// Command formats and outputs the Destroy command from
// the provided configuration to destroy to resources.
func (d *Destroy) Command() *exec.Cmd {
	logrus.Trace("creating terraform destroy command from plugin configuration")

	// global Variables
	var globalFlags []string

	// variable to store flags for command
	var flags []string

	// check if Directory is provided
	if d.Directory != "." && SupportsChdir(d.Version) {
		globalFlags = append(flags, fmt.Sprintf("-chdir=%s", d.Directory))
	}

	// check if AutoApprove is provided
	if d.AutoApprove {
		// add flag for AutoApprove from provided destroy command
		flags = append(flags, "-auto-approve")
	}

	// check if Backup is provided
	if len(d.Backup) > 0 {
		// add flag for Backup from provided destroy command
		flags = append(flags, fmt.Sprintf("-backup=%s", d.Backup))
	}

	// check if Lock is provided
	if d.Lock {
		// add flag for Lock from provided destroy command
		flags = append(flags, "-lock=true")
	}

	// check if LockTimeout is provided
	if d.LockTimeout > 0 {
		// add flag for LockTimeout from provided destroy command
		flags = append(flags, fmt.Sprintf("-lock-timeout=%s", d.LockTimeout))
	}

	// check if NoColor is provided
	if d.NoColor {
		// add flag for NoColor from provided destroy command
		flags = append(flags, "-no-color")
	}

	// check if Parallelism is provided
	if d.Parallelism > 0 {
		// add flag for Parallelism from provided destroy command
		flags = append(flags, fmt.Sprintf("-parallelism=%d", d.Parallelism))
	}

	// check if Refresh is provided
	if d.Refresh {
		// add flag for Refresh from provided destroy command
		flags = append(flags, "-refresh=true")
	}

	// check if State is provided
	if len(d.State) > 0 {
		// add flag for State from provided destroy command
		flags = append(flags, fmt.Sprintf("-state=%s", d.State))
	}

	// check if StateOut is provided
	if len(d.StateOut) > 0 {
		// add flag for StateOut from provided destroy command
		flags = append(flags, fmt.Sprintf("-state-out=%s", d.StateOut))
	}

	// check if Target is provided
	if len(d.Target) > 0 {
		// add flag for Target from provided destroy command
		flags = append(flags, fmt.Sprintf("-target=%s", d.Target))
	}

	// check if Vars is provided
	if len(d.Vars) > 0 {
		for _, v := range d.Vars {
			// add flag for Vars from provided command
			flags = append(flags, fmt.Sprintf(`-var=%s`, v))
		}
	}

	// check if VarFiles is provided
	if len(d.VarFiles) > 0 {
		for _, v := range d.VarFiles {
			// add flag for VarFiles from provided command
			flags = append(flags, fmt.Sprintf(`-var-file=%s`, v))
		}
	}

	// check if Directory is provided and terraform version doesn't support chdir
	if d.Directory != "." && !SupportsChdir(d.Version) {
		flags = append(flags, d.Directory)
	}

	globalFlags = append(globalFlags, destroyAction)

	// nolint: gosec // ignore G204
	return exec.Command(_terraform, append(globalFlags, flags...)...)
}

// Exec formats and runs the commands for destroying resources with Terraform.
func (d *Destroy) Exec() error {
	logrus.Trace("running destroy with provided configuration")

	// create the destroy command for the file
	cmd := d.Command()

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
		logrus.Warn("terraform destroy will run in current dir")
	}

	return nil
}
