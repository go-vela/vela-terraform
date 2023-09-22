// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os/exec"
	"reflect"
	"testing"

	"github.com/Masterminds/semver/v3"
)

func TestTerraform_FMT_Command(t *testing.T) {
	v, _ := semver.NewVersion("1.0.0")
	// setup types
	f := &FMT{
		Check:     true,
		Diff:      true,
		Directory: "foobar/",
		List:      false,
		Write:     false,
		Version:   v,
	}

	// nolint: gosec // ignore G204
	want := exec.Command(
		_terraform,
		fmt.Sprintf("-chdir=%s", f.Directory),
		fmtAction,
		fmt.Sprintf("-list=%t", f.List),
		fmt.Sprintf("-write=%t", f.Write),
		fmt.Sprintf("-diff=%t", f.Diff),
		fmt.Sprintf("-check=%t", f.Check),
	)

	got := f.Command()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Command is %v, want %v", got, want)
	}
}

func TestTerraform_FMT_Command_tf13(t *testing.T) {
	v, _ := semver.NewVersion("0.13.0")
	// setup types
	f := &FMT{
		Check:     true,
		Diff:      true,
		Directory: "foobar/",
		List:      false,
		Write:     false,
		Version:   v,
	}

	// nolint: gosec // ignore G204
	want := exec.Command(
		_terraform,
		fmtAction,
		fmt.Sprintf("-list=%t", f.List),
		fmt.Sprintf("-write=%t", f.Write),
		fmt.Sprintf("-diff=%t", f.Diff),
		fmt.Sprintf("-check=%t", f.Check),
		fmt.Sprintf(f.Directory),
	)

	got := f.Command()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Command is %v, want %v", got, want)
	}
}

func TestTerraform_FMT_Exec_Error(t *testing.T) {
	v, _ := semver.NewVersion("1.0.0")
	// setup types
	f := &FMT{
		Directory: "foobar/",
		Version:   v,
	}

	err := f.Exec()
	if err == nil {
		t.Errorf("Exec should have returned err")
	}
}

func TestTerraform_FMT_Validate(t *testing.T) {
	v, _ := semver.NewVersion("1.0.0")
	// setup types
	tests := []struct {
		fmt *FMT
	}{
		{
			fmt: &FMT{Directory: "foobar/", Version: v},
		},
		{
			fmt: &FMT{Directory: "foobar.tf", Version: v},
		},
		{
			fmt: &FMT{Directory: "", Version: v},
		},
	}

	// run test
	for _, test := range tests {
		err := test.fmt.Validate()
		if err != nil {
			t.Errorf("Validate returned err: %v", err)
		}
	}
}
