// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"
	"testing"

	"github.com/spf13/afero"
)

func TestTerraform_install(t *testing.T) {
	// setup filesystem
	appFS = afero.NewMemMapFs()

	// run test
	err := installBinary(t.Context(), "0.11.0", "0.11.0")
	if err != nil {
		t.Errorf("install returned err: %v", err)
	}
}

func TestTerraform_install_NoBinary(t *testing.T) {
	// setup filesystem
	appFS = afero.NewMemMapFs()

	// run test
	err := installBinary(t.Context(), "0.11.0", "0.12.0")
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
	err = installBinary(t.Context(), "0.11.0", "0.12.0")
	if err == nil {
		t.Errorf("install should have returned err")
	}
}

func TestTerraform_env(t *testing.T) {
	want := "abc123"
	up := "TF_VAR_CHEF_PRIVATE_KEY"
	low := "TF_VAR_chef_private_key"

	// setup env
	t.Setenv(up, want)

	// check env
	got := os.Getenv(low)
	if got == want {
		t.Errorf("os.Getenv should not be %v", got)
	}

	// run env
	err := env()
	if err != nil {
		t.Errorf("env returned err: %v", err)
	}

	// check new env for same value
	got = os.Getenv(low)
	if got != want {
		t.Errorf("os.Getenv is %v, want %v", got, want)
	}
}

func TestTerraform_env_err(t *testing.T) {
	// run env
	err := env()
	if err != nil {
		t.Errorf("env returned err: %v", err)
	}
}
