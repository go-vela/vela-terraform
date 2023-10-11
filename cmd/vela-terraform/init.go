// SPDX-License-Identifier: Apache-2.0

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
		Backend bool `json:"backend,omitempty"`
		// This is merged with what is in the configuration file i.e. "-backend-config=path"
		BackendConfigs []string `json:"backend_configs,omitempty"`
		// Suppress prompts about copying state data i.e. "-force-copy"
		ForceCopy bool `json:"force_copy,omitempty"`
		// Copy the contents of the given module into the target
		// directory before initialization. "-from-module=SOURCE"
		FromModule string `json:"from_module,omitempty"`
		// Download any modules for this configuration i.e. "-get=true"
		Get bool `json:"get,omitempty"`
		// Download any missing plugins for this configuration i.e. "-get-plugins=true"
		GetPlugins bool `json:"get_plugins,omitempty"`
		// ask for input for variables if not directly set. i.e. "-input=true"
		Input bool `json:"input,omitempty"`
		// the state file when locking is supported. i.e. -lock=true
		Lock bool `json:"lock,omitempty"`
		// duration to retry a state lock. i.e. "-lock-timeout=0s"
		LockTimeout time.Duration `json:"lock_timeout,omitempty"`
		// if specified, output won't contain any color. i.e. "-no-color"
		NoColor bool `json:"no_color,omitempty"`
		// Directory containing plugin binaries.
		// This overrides all default search paths for plugins i.e. "-plugin-dir"
		PluginDirs []string `json:"plugin_dirs,omitempty"`
		// Reconfigure the backend, ignoring any saved configuration i.e. "-reconfigure"
		Reconfigure bool `json:"reconfigure,omitempty"`
		// install the latest version allowed within configured constraints i.e. "-upgrade=false"
		Upgrade bool `json:"upgrade,omitempty"`
		// Verify the authenticity and integrity of automatically
		// downloaded plugins i.e. "-verify-plugins=true"
		VerifyPlugins bool `json:"verify_plugins,omitempty"`
	}
)

// Command formats and outputs the Init command from
// the provided configuration to init to resources.
func (i *Init) Command() *exec.Cmd {
	logrus.Trace("creating terraform init command from plugin configuration")

	// global Variables
	var globalFlags []string

	// variable to store flags for command
	var flags []string

	// check if Directory is provided
	if i.Directory != "." {
		globalFlags = append(flags, fmt.Sprintf("-chdir=%s", i.Directory))
	}

	// check if Backend is provided
	if i.InitOptions.Backend {
		// add flag for Backend from provided init command
		flags = append(flags, "-backend=true")
	}

	// check if BackendConfigs is provided
	if len(i.InitOptions.BackendConfigs) > 0 {
		for _, v := range i.InitOptions.BackendConfigs {
			// add flag for BackendConfigs from provided init command
			flags = append(flags, fmt.Sprintf(`-backend-config=%s`, v))
		}
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
		for _, v := range i.InitOptions.PluginDirs {
			// add flag for PluginDirs from provided init command
			flags = append(flags, fmt.Sprintf(`-plugin-dir=%s`, v))
		}
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

	globalFlags = append(globalFlags, initAction)

	// nolint: gosec // ignore G204
	return exec.Command(_terraform, append(globalFlags, flags...)...)
}

// Exec formats and runs the commands for initing Terraform.
func (i *Init) Exec() error {
	logrus.Trace("running init with provided configuration")

	// create the init command for the file
	cmd := i.Command()

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
		logrus.Warn("terraform init will run in current dir")
	}

	return nil
}

// Unmarshal captures the provided properties and
// serializes them into their expected form.
func (i *Init) Unmarshal() error {
	logrus.Trace("unmarshalling init options")

	i.InitOptions = &InitOptions{}

	// check if any options were passed
	if len(i.RawInit) > 0 {
		// cast raw properties into bytes
		bytes := []byte(i.RawInit)

		// serialize raw properties into expected Props type
		err := json.Unmarshal(bytes, &i.InitOptions)
		if err != nil {
			return fmt.Errorf("failed to unmarshal init options: %w", err)
		}
	}

	return nil
}
