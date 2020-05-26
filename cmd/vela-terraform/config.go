// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"
	"os/user"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

// netrcFile represents an empty .netrc config file
const netrcFile = `
machine %s
login %s
password %s
`

type (
	// Config holds input parameters for the plugin
	Config struct {
		// TODO: delete this code
		Actions     []string
		Vars        map[string]string
		Secrets     map[string]string
		InitOptions InitOptions
		FmtOptions  FmtOptions
		Cacert      string
		Sensitive   bool
		RoleARN     string
		RootDir     string
		Parallelism int
		Targets     []string
		VarFiles    []string

		// action to perform with Terraform
		Action string
		// Netrc is credentials for cloning
		Netrc *Netrc
	}

	// Netrc is credentials for cloning
	Netrc struct {
		Machine  string
		Login    string
		Password string
	}
)

var appFS = afero.NewOsFs()

// Write creates a .netrc file with the credentials provided in the plugin environment.
func (c *Config) Write() error {
	logrus.Trace("writing .netrc credentials file")

	// use custom filesystem which enables us to test
	a := &afero.Afero{
		Fs: appFS,
	}

	// create the .netrc file from the provided configuration
	file := fmt.Sprintf(netrcFile, c.Netrc.Machine, c.Netrc.Login, c.Netrc.Password)

	// set default home directory for root user
	home := "/root"

	// capture current user running commands
	u, err := user.Current()
	if err == nil {
		// set home directory to current user
		home = u.HomeDir
	}

	return a.WriteFile(filepath.Join(home, ".netrc"), []byte(file), 0600)
}

// Validate verifies the Config is properly configured.
func (c *Config) Validate() error {
	logrus.Trace("validating config plugin configuration")

	// verify action is provided
	if len(c.Action) == 0 {
		return fmt.Errorf("no config action provided")
	}

	return nil
}
