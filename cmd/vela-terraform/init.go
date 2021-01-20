// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const initAction = "init"

type (
	// Init represents the plugin configuration for init information.
	Init struct {
		// terraform file or directory to init
		Directory string
		// init for initialize a new or existing Terraform working directory
		InitOptions *InitOptions
		// raw input of init options provided for plugin
		RawInit string
	}

	// InitOptions represents the plugin configuration for options for init.
	InitOptions struct {
		// Configure the backend for this configuration i.e. "-backend=true"
		Backend bool
		// This is merged with what is in the configuration file i.e. "-backend-config=path"
		BackendConfigs []string
		// Suppress prompts about copying state data i.e. "-force-copy"
		ForceCopy bool
		// Copy the contents of the given module into the target directory before initialization. "-from-module=SOURCE"
		FromModule string
		// Download any modules for this configuration i.e. "-get=true"
		Get bool
		// Download any missing plugins for this configuration i.e. "-get-plugins=true"
		GetPlugins bool
		// ask for input for variables if not directly set. i.e. "-input=true"
		Input bool
		// the state file when locking is supported. i.e. -lock=true
		Lock bool
		// duration to retry a state lock. i.e. "-lock-timeout=0s"
		LockTimeout time.Duration
		// if specified, output won't contain any color. i.e. "-no-color"
		NoColor bool
		// Directory containing plugin binaries. This overrides all default search paths for plugins i.e. "-plugin-dir"
		PluginDirs []string
		// Reconfigure the backend, ignoring any saved configuration i.e. "-reconfigure"
		Reconfigure bool
		// install the latest version allowed within configured constraints i.e. "-upgrade=false"
		Upgrade bool
		// Verify the authenticity and integrity of automatically downloaded plugins i.e. "-verify-plugins=true"
		VerifyPlugins bool
	}
)

// Command formats and outputs the Init command from
// the provided configuration to init to resources.
func (i *Init) Command(dir string) *exec.Cmd {
	logrus.Trace("creating terraform init command from plugin configuration")

	// variable to store flags for command
	var flags []string

	// check if Backend is provided
	if i.InitOptions.Backend {
		// add flag for Backend from provided init command
		flags = append(flags, "-backend=true")
	}

	// check if BackendConfigs is provided
	if len(i.InitOptions.BackendConfigs) > 0 {
		var configs string
		for _, v := range i.InitOptions.BackendConfigs {
			configs += fmt.Sprintf("-backend-config=%s ", v)
		}

		// add flag for BackendConfigs from provided init command
		flags = append(flags, strings.TrimSuffix(configs, " "))
	}

	// check if ForceCopy is provided
	if i.InitOptions.ForceCopy {
		// add flag for ForceCopy from provided init command
		flags = append(flags, "-force-copy")
	}

	// check if FromModule is provided
	if len(i.InitOptions.FromModule) > 0 {
		// add flag for FromModule from provided init command
		flags = append(flags, fmt.Sprintf("-from-module=%s", i.InitOptions.FromModule))
	}

	// check if Get is provided
	if i.InitOptions.Get {
		// add flag for Get from provided init command
		flags = append(flags, "-get=true")
	}

	// check if GetPlugins is provided
	if i.InitOptions.GetPlugins {
		// add flag for GetPlugins from provided init command
		flags = append(flags, "-get-plugins=true")
	}

	// check if Input is provided
	if i.InitOptions.Input {
		// add flag for Input from provided init command
		flags = append(flags, "-input=true")
	}

	// check if Lock is provided
	if i.InitOptions.Lock {
		// add flag for Lock from provided init command
		flags = append(flags, "-lock=true")
	}

	// check if LockTimeout is provided
	if i.InitOptions.LockTimeout > 0 {
		// add flag for LockTimeout from provided init command
		flags = append(flags, fmt.Sprintf("-lock-timeout=%s", i.InitOptions.LockTimeout))
	}

	// check if NoColor is provided
	if i.InitOptions.NoColor {
		// add flag for NoColor from provided init command
		flags = append(flags, "-no-color")
	}

	// check if PluginDirs is provided
	if len(i.InitOptions.PluginDirs) > 0 {
		var configs string
		for _, v := range i.InitOptions.PluginDirs {
			configs += fmt.Sprintf("-plugin-dir=%s ", v)
		}

		// add flag for PluginDirs from provided init command
		flags = append(flags, strings.TrimSuffix(configs, " "))
	}

	// check if Reconfigure is provided
	if i.InitOptions.Reconfigure {
		// add flag for Reconfigure from provided init command
		flags = append(flags, "-reconfigure")
	}

	// check if Upgrade is provided
	if i.InitOptions.Upgrade {
		// add flag for Upgrade from provided init command
		flags = append(flags, "-upgrade=false")
	}

	// check if VerifyPlugins is provided
	if i.InitOptions.VerifyPlugins {
		// add flag for VerifyPlugins from provided init command
		flags = append(flags, "-verify-plugins=true")
	}

	// add the required dir param
	flags = append(flags, dir)

	return exec.Command(_terraform, append([]string{initAction}, flags...)...)
}

// Exec formats and runs the commands for initing Terraform.
func (i *Init) Exec() error {
	logrus.Trace("running init with provided configuration")

	// create the init command for the file
	cmd := i.Command(i.Directory)

	// run the init command for the file
	err := execCmd(cmd)
	if err != nil {
		return err
	}

	return nil
}

// Validate verifies the Init is properly configured.
func (i *Init) Validate() error {
	logrus.Trace("validating plan plugin configuration")

	if strings.EqualFold(i.Directory, ".") {
		logrus.Warn("terrafrom init will run in current dir")
	}

	return nil
}

// Unmarshal captures the provided properties and
// serializes them into their expected form.
func (i *Init) Unmarshal() error {
	logrus.Trace("unmarshaling init options")

	i.InitOptions = &InitOptions{}

	// check if any options were passed
	if len(i.RawInit) > 0 {
		// cast raw properties into bytes
		bytes := []byte(i.RawInit)

		// serialize raw properties into expected Props type
		err := json.Unmarshal(bytes, &i.InitOptions)
		if err != nil {
			return err
		}
	}

	return nil
}
