// SPDX-License-Identifier: Apache-2.0

package main

import (
	"testing"

	"github.com/spf13/afero"
)

func TestTerraform_Config_Write(t *testing.T) {
	// setup filesystem
	appFS = afero.NewMemMapFs()

	// setup types
	c := &Config{
		Netrc: &Netrc{
			Machine:  "machine.com",
			Login:    "octocat",
			Password: "mypassword",
		},
	}

	err := c.Write()
	if err != nil {
		t.Errorf("Write returned err: %v", err)
	}
}

func TestKubernetes_Config_Write_Error(t *testing.T) {
	// setup filesystem
	appFS = afero.NewReadOnlyFs(afero.NewMemMapFs())

	// setup types
	c := &Config{
		Netrc: &Netrc{
			Machine:  "machine.com",
			Login:    "octocat",
			Password: "mypassword",
		},
	}

	err := c.Write()
	if err == nil {
		t.Errorf("Write should have returned err")
	}
}

func TestTerraform_Config_Validate(t *testing.T) {
	// setup types
	c := &Config{
		Action: "apply",
	}

	err := c.Validate()
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}

func TestTerraform_Config_Validate_Error(t *testing.T) {
	// setup types
	c := &Config{}

	err := c.Validate()
	if err == nil {
		t.Errorf("Write should have returned err")
	}
}
