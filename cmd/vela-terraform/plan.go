// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"github.com/sirupsen/logrus"
)

const planAction = "plan"

// Plan represents the plugin configuration for plan information.
type Plan struct {
}

// Exec formats and runs the commands for removing artifacts in Artifactory.
func (p *Plan) Exec() error {
	logrus.Trace("running delete with provided configuration")

	return nil
}

// Validate verifies the Delete is properly configured.
func (p *Plan) Validate() error {
	logrus.Trace("validating delete plugin configuration")

	return nil
}
