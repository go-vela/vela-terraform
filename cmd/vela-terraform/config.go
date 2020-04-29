// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"github.com/sirupsen/logrus"
)

type (
	// Config holds input parameters for the plugin
	Config struct {
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
	}

	// Netrc is credentials for cloning
	Netrc struct {
		Machine  string
		Login    string
		Password string
	}
)

// New creates an Artifactory client for managing artifacts.
func (c *Config) New() error {
	logrus.Trace("creating new Terraform client from plugin configuration")

	return nil
}

// Validate verifies the Config is properly configured.
func (c *Config) Validate() error {
	logrus.Trace("validating config plugin configuration")

	return nil
}
