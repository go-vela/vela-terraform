// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"
	"os/exec"
	"reflect"
	"testing"

	"github.com/Masterminds/semver"
)

func TestTerraform_Validation_Command(t *testing.T) {
	ver, _ := semver.NewVersion("1.0.0")
	// setup types
	v := &Validation{
		CheckVariables: true,
		Directory:      "foobar/",
		NoColor:        true,
		Vars:           []string{"foo=bar", "bar=foo"},
		VarFiles:       []string{"vars1.tf", "vars2.tf"},
		Version:        ver,
	}

	want := exec.Command(
		_terraform,
		fmt.Sprintf("-chdir=%s", v.Directory),
		validationAction,
		fmt.Sprintf("-check-variables=%t", v.CheckVariables),
		"-no-color",
		fmt.Sprintf("-var=%s", v.Vars[0]),
		fmt.Sprintf("-var=%s", v.Vars[1]),
		fmt.Sprintf("-var-file=%s", v.VarFiles[0]),
		fmt.Sprintf("-var-file=%s", v.VarFiles[1]),
	)

	got := v.Command()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Command is %v, want %v", got, want)
	}
}

func TestTerraform_Validation_Command_tf13(t *testing.T) {
	ver, _ := semver.NewVersion("0.13.0")
	// setup types
	v := &Validation{
		CheckVariables: true,
		Directory:      "foobar/",
		NoColor:        true,
		Vars:           []string{"foo=bar", "bar=foo"},
		VarFiles:       []string{"vars1.tf", "vars2.tf"},
		Version:        ver,
	}

	want := exec.Command(
		_terraform,
		validationAction,
		fmt.Sprintf("-check-variables=%t", v.CheckVariables),
		"-no-color",
		fmt.Sprintf("-var=%s", v.Vars[0]),
		fmt.Sprintf("-var=%s", v.Vars[1]),
		fmt.Sprintf("-var-file=%s", v.VarFiles[0]),
		fmt.Sprintf("-var-file=%s", v.VarFiles[1]),
		fmt.Sprintf(v.Directory),
	)

	got := v.Command()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Command is %v, want %v", got, want)
	}
}

func TestTerraform_Validation_Exec(t *testing.T) {
	ver, _ := semver.NewVersion("1.0.0")
	// setup types
	v := &Validation{
		Directory: "foobar/",
		Version:   ver,
	}

	err := v.Exec()
	if err == nil {
		t.Errorf("Exec should have returned err")
	}
}

func TestTerraform_Validation_Validate(t *testing.T) {
	ver, _ := semver.NewVersion("1.0.0")
	// setup types
	tests := []struct {
		validation *Validation
	}{
		{
			validation: &Validation{Directory: "foobar/", Version: ver},
		},
		{
			validation: &Validation{Directory: "foobar.tf", Version: ver},
		},
		{
			validation: &Validation{Directory: "", Version: ver},
		},
	}

	// run test
	for _, test := range tests {
		err := test.validation.Validate()
		if err != nil {
			t.Errorf("Validate returned err: %v", err)
		}
	}
}
