// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"github.com/sirupsen/logrus"
)

const applyAction = "apply"

// Apply represents the plugin configuration for apply information.
type Apply struct {
}

// Exec formats and runs the commands for removing artifacts in Artifactory.
func (a *Apply) Exec() error {
	logrus.Trace("running delete with provided configuration")

	return nil
}

// Validate verifies the Delete is properly configured.
func (a *Apply) Validate() error {
	logrus.Trace("validating delete plugin configuration")

	return nil
}
