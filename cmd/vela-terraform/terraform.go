// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/hashicorp/go-version"
	install "github.com/hashicorp/hc-install"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/hc-install/src"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

const (
	_installDir = "/bin"
	_terraform  = _installDir + "/" + "terraform"
)

func installBinary(ctx context.Context, customVer, defaultVer string) error {
	// use custom filesystem which enables us to test
	a := &afero.Afero{
		Fs: appFS,
	}

	// check if the custom version matches the default version
	if strings.EqualFold(customVer, defaultVer) {
		// the terraform versions match so no action is required
		return nil
	}

	// parse the custom version
	v, err := version.NewVersion(customVer)
	if err != nil {
		return err
	}

	logrus.Infof("custom terraform version requested: %s", customVer)

	logrus.Debugf("custom version does not match default: %s", defaultVer)

	// rename the old terraform binary since we can't overwrite it for now
	err = a.Rename(_terraform, fmt.Sprintf("%s.default", _terraform))
	if err != nil {
		return err
	}

	// use hc-install to install the custom version
	installer := install.NewInstaller()
	_, err = installer.Install(ctx, []src.Installable{
		&releases.ExactVersion{
			Product:    product.Terraform,
			Version:    v,
			InstallDir: _installDir,
		},
	})

	if err != nil {
		return err
	}

	return nil
}

// sets up environment for terraform.
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
