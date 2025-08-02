// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os/exec"
	"reflect"
	"testing"

	"github.com/Masterminds/semver/v3"
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

	//nolint:gosec // ignore G204
	want := exec.CommandContext(
		t.Context(),
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

	got := v.Command(t.Context())
	if got.Path != want.Path {
		t.Errorf("Command path is %v, want %v", got.Path, want.Path)
	}

	if !reflect.DeepEqual(got.Args, want.Args) {
		t.Errorf("Command args is %v, want %v", got.Args, want.Args)
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

	//nolint:gosec // ignore G204
	want := exec.CommandContext(
		t.Context(),
		_terraform,
		validationAction,
		fmt.Sprintf("-check-variables=%t", v.CheckVariables),
		"-no-color",
		fmt.Sprintf("-var=%s", v.Vars[0]),
		fmt.Sprintf("-var=%s", v.Vars[1]),
		fmt.Sprintf("-var-file=%s", v.VarFiles[0]),
		fmt.Sprintf("-var-file=%s", v.VarFiles[1]),
		fmt.Sprint(v.Directory),
	)

	got := v.Command(t.Context())
	if got.Path != want.Path {
		t.Errorf("Command path is %v, want %v", got.Path, want.Path)
	}

	if !reflect.DeepEqual(got.Args, want.Args) {
		t.Errorf("Command args is %v, want %v", got.Args, want.Args)
	}
}

func TestTerraform_Validation_Exec(t *testing.T) {
	ver, _ := semver.NewVersion("1.0.0")
	// setup types
	v := &Validation{
		Directory: "foobar/",
		Version:   ver,
	}

	err := v.Exec(t.Context())
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
