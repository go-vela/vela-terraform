// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"
)

type (
	// InitOptions include options for the Terraform's init command
	InitOptions struct {
		BackendConfig []string `json:"backend-config"`
		Lock          *bool    `json:"lock"`
		LockTimeout   string   `json:"lock-timeout"`
	}

	// FmtOptions fmt options for the Terraform's fmt command
	FmtOptions struct {
		List  *bool `json:"list"`
		Write *bool `json:"write"`
		Diff  *bool `json:"diff"`
		Check *bool `json:"check"`
	}

	// Plugin represents the plugin instance to be executed
	Plugin struct {
		// TODO: remove theses legacy params
		Config Config
		Netrc  Netrc

		// Apply arguments loaded for the plugin
		Apply *Apply
		// Plan arguments loaded for the plugin
		Plan *Plan
		// Validation arguments loaded for the plugin
		Validation *Validation
	}
)

// Exec formats and runs the commands for running Terraform commands.
func (p *Plugin) Exec() error {
	logrus.Debug("running plugin with provided configuration")

	if p.Config.RoleARN != "" {
		assumeRole(p.Config.RoleARN)
	}

	// writing the .netrc file with Github credentials in it.
	err := writeNetrc(p.Netrc.Machine, p.Netrc.Login, p.Netrc.Password)
	if err != nil {
		return err
	}

	var commands []*exec.Cmd

	commands = append(commands, exec.Command("terraform", "version"))

	CopyTfEnv()

	if p.Config.Cacert != "" {
		commands = append(commands, installCaCert(p.Config.Cacert))
	}

	commands = append(commands, deleteCache())
	commands = append(commands, initCommand(p.Config.InitOptions))
	commands = append(commands, getModules())

	// Add commands listed from Actions
	for _, action := range p.Config.Actions {
		switch action {
		case "fmt":
			commands = append(commands, tfFmt(p.Config))
		case "validate":
			commands = append(commands, tfValidate(p.Config))
		case "plan":
			commands = append(commands, tfPlan(p.Config, false))
		case "plan-destroy":
			commands = append(commands, tfPlan(p.Config, true))
		case "apply":
			commands = append(commands, tfApply(p.Config))
		case "destroy":
			commands = append(commands, tfDestroy(p.Config))
		default:
			return fmt.Errorf("valid actions are: fmt, validate, plan, apply, plan-destroy, destroy.  You provided %s", action)
		}
	}

	commands = append(commands, deleteCache())

	for _, c := range commands {
		if c.Dir == "" {
			wd, err := os.Getwd()
			if err == nil {
				c.Dir = wd
			}
		}
		if p.Config.RootDir != "" {
			c.Dir = c.Dir + "/" + p.Config.RootDir
		}
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		if !p.Config.Sensitive {
			trace(c)
		}

		err := c.Run()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Fatal("Failed to execute a command")
		}
		logrus.Debug("Command completed successfully")
	}

	return nil
}

// Validate verifies the plugin is properly configured.
func (p *Plugin) Validate() error {
	logrus.Debug("validating plugin configuration")

	return nil
}
