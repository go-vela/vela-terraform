// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"testing"

	"github.com/spf13/afero"
)

func TestTerraform_install(t *testing.T) {
	// setup filesystem
	appFS = afero.NewMemMapFs()

	// run test
	err := install("0.11.0", "0.11.0")
	if err != nil {
		t.Errorf("install returned err: %v", err)
	}
}

func TestTerraform_install_NoBinary(t *testing.T) {
	// setup filesystem
	appFS = afero.NewMemMapFs()

	// run test
	err := install("0.11.0", "0.12.0")
	if err == nil {
		t.Errorf("install should have returned err")
	}
}

func TestTerraform_install_NotWritable(t *testing.T) {
	// setup filesystem
	appFS = afero.NewMemMapFs()

	a := &afero.Afero{
		Fs: appFS,
	}

	// create binary file
	err := a.WriteFile(_terraform, []byte("!@#$%^&*()"), 0777)
	if err != nil {
		t.Errorf("Unable to write file %s: %v", _terraform, err)
	}

	// run test
	err = install("0.11.0", "0.12.0")
	if err == nil {
		t.Errorf("install should have returned err")
	}
}
