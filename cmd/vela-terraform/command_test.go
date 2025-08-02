// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os/exec"
	"reflect"
	"testing"
)

func TestTerraform_execCmd(t *testing.T) {
	// setup types
	e := exec.Command("echo", "hello")

	err := execCmd(e)
	if err != nil {
		t.Errorf("execCmd returned err: %v", err)
	}
}

func TestTerraform_versionCmd(t *testing.T) {
	// setup types
	want := exec.CommandContext(
		t.Context(),
		_terraform,
		"version",
	)

	got := versionCmd(t.Context())
	if got.Path != want.Path {
		t.Errorf("versionCmd path is %v, want %v", got.Path, want.Path)
	}

	if !reflect.DeepEqual(got.Args, want.Args) {
		t.Errorf("versionCmd args is %v, want %v", got.Args, want.Args)
	}
}
