// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"github.com/sirupsen/logrus"
)

const validateAction = "validate"

// Validate represents the plugin configuration for validate information.
type Validate struct {
}

// Exec formats and runs the commands for removing artifacts in Artifactory.
func (v *Validate) Exec() error {
	logrus.Trace("running validate with provided configuration")

	return nil
}

// Validate verifies the Delete is properly configured.
func (v *Validate) Validate() error {
	logrus.Trace("validating validate plugin configuration")

	return nil
}
