// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
)

type (
	// Plugin represents the plugin instance to be executed.
	Plugin struct {
		// Apply arguments loaded for the plugin
		Apply *Apply
		// config arguments loaded for the plugin
		Config *Config
		// Destroy arguments loaded for the plugin
		Destroy *Destroy
		// InitOptions arguments loaded for the plugin
		Init *Init
		// FMT arguments loaded for the plugin
		FMT *FMT
		// Plan arguments loaded for the plugin
		Plan *Plan
		// Validation arguments loaded for the plugin
		Validation *Validation
	}
)

var (
	// ErrInvalidAction defines the error type when the
	// Action provided to the Plugin is unsupported.
	ErrInvalidAction = errors.New("invalid action provided")
)

// Exec formats and runs the commands for running Terraform commands.
func (p *Plugin) Exec() error {
	logrus.Debug("running plugin with provided configuration")

	// write the .netrc file with Github credentials
	err := p.Config.Write()
	if err != nil {
		return err
	}

	// output terraform version for troubleshooting
	err = execCmd(versionCmd())
	if err != nil {
		return err
	}

	// initialize a new or existing Terraform working directory
	err = p.Init.Exec()
	if err != nil {
		return err
	}

	// retrieve terraform modules for actions
	err = execCmd(getCmd(p.Init.Directory))
	if err != nil {
		return err
	}

	// configure the terraform environment
	err = env()
	if err != nil {
		return err
	}

	// execute action specific configuration
	switch p.Config.Action {
	case applyAction:
		// execute apply action
		return p.Apply.Exec()
	case destroyAction:
		// execute destroy action
		return p.Destroy.Exec()
	case fmtAction:
		// execute fmt action
		return p.FMT.Exec()
	case planAction:
		// execute plan action
		return p.Plan.Exec()
	case validationAction:
		// execute validate action
		return p.Validation.Exec()
	default:
		return fmt.Errorf(
			"%w: %s (Valid actions: %s, %s, %s, %s, %s)",
			ErrInvalidAction,
			p.Config.Action,
			applyAction,
			destroyAction,
			fmtAction,
			planAction,
			validationAction,
		)
	}
}

// Validate verifies the plugin is properly configured.
func (p *Plugin) Validate() error {
	logrus.Debug("validating plugin configuration")

	// validate config configuration
	err := p.Config.Validate()
	if err != nil {
		return err
	}

	// when user adds additional init config
	// unmarshal it into the init command
	err = p.Init.Unmarshal()
	if err != nil {
		return err
	}

	// validate init action
	err = p.Init.Validate()
	if err != nil {
		return err
	}

	// validate action specific configuration
	switch p.Config.Action {
	case applyAction:
		// validate apply action
		return p.Apply.Validate()
	case destroyAction:
		// validate destroy action
		return p.Destroy.Validate()
	case initAction:
		// validate fmt action
		return p.Init.Validate()
	case fmtAction:
		// validate fmt action
		return p.FMT.Validate()
	case planAction:
		// validate plan action
		return p.Plan.Validate()
	case validationAction:
		// validate validate action
		return p.Validation.Validate()
	default:
		return fmt.Errorf(
			"%w: %s (Valid actions: %s, %s, %s, %s, %s, %s)",
			ErrInvalidAction,
			p.Config.Action,
			applyAction,
			destroyAction,
			initAction,
			fmtAction,
			planAction,
			validationAction,
		)
	}
}
