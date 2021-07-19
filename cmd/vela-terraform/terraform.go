// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"

	getter "github.com/hashicorp/go-getter/v2"
)

const (
	_terraform = "/bin/terraform"
	_download  = "https://releases.hashicorp.com/terraform/%s/terraform_%s_linux_amd64.zip"
)

func install(customVer, defaultVer string) error {
	logrus.Infof("custom terraform version requested: %s", customVer)

	// use custom filesystem which enables us to test
	a := &afero.Afero{
		Fs: appFS,
	}

	// check if the custom version matches the default version
	if strings.EqualFold(customVer, defaultVer) {
		// the terraform versions match so no action is required
		return nil
	}

	logrus.Debugf("custom version does not match default: %s", defaultVer)
	// rename the old terraform binary since we can't overwrite it for now
	//
	// https://github.com/hashicorp/go-getter/issues/219
	err := a.Rename(_terraform, fmt.Sprintf("%s.default", _terraform))
	if err != nil {
		return err
	}

	// create the download URL to install terraform
	url := fmt.Sprintf(_download, customVer, customVer)

	logrus.Infof("downloading terraform version from: %s", url)
	// send the HTTP request to install terraform
	err = getter.GetFile(_terraform, url, []getter.ClientOption{}...)
	if err != nil {
		return err
	}

	logrus.Debugf("changing ownership of file: %s", _terraform)
	// ensure the terraform binary is executable
	err = a.Chmod(_terraform, 0700)
	if err != nil {
		return err
	}

	return nil
}

// sets up environment for terraform
func env() error {
	// regexp for TF_VAR_ terraform vars
	tfVar := regexp.MustCompile(`^TF_VAR_.*$`)

	// match terraform vars in environment
	for _, e := range os.Environ() {
		// split on value
		pair := strings.SplitN(e, "=", 2)

		// match on TF_VAR_*
		if tfVar.MatchString(pair[0]) {
			// pull out the name
			name := strings.Split(pair[0], "TF_VAR_")

			// lower case the terraform variable
			//   to accommodate cicd injection capitalization
			err := os.Setenv(fmt.Sprintf("TF_VAR_%s",
				strings.ToLower(name[1])), pair[1])
			if err != nil {
				return err
			}
		}
	}
	return nil
}
