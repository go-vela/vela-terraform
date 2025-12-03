// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os/exec"
	"slices"
	"testing"
	"time"

	"github.com/Masterminds/semver/v3"
)

func TestTerraform_Apply_Command(t *testing.T) {
	// setup types
	v, err := semver.NewVersion("1.0.0")
	if err != nil {
		t.Error(err)
	}

	a := &Apply{
		AutoApprove: true,
		Backup:      "backup/",
		Directory:   "foobar/",
		Lock:        true,
		LockTimeout: 1 * time.Second,
		Input:       true,
		NoColor:     true,
		Parallelism: 1,
		Refresh:     true,
		State:       "state.tf",
		StateOut:    "stateout.tf",
		Target:      "target.tf",
		Vars:        []string{"foo=bar", "bar=foo"},
		VarFiles:    []string{"vars1.tf", "vars2.tf"},
		Version:     v,
	}

	//nolint:gosec // ignore G204
	want := exec.CommandContext(
		t.Context(),
		_terraform,
		fmt.Sprintf("-chdir=%s", a.Directory),
		applyAction,
		"-auto-approve",
		fmt.Sprintf("-backup=%s", a.Backup),
		"-lock=true",
		fmt.Sprintf("-lock-timeout=%s", a.LockTimeout),
		"-input=true",
		"-no-color",
		fmt.Sprintf("-parallelism=%d", a.Parallelism),
		"-refresh=true",
		fmt.Sprintf("-state=%s", a.State),
		fmt.Sprintf("-state-out=%s", a.StateOut),
		fmt.Sprintf("-target=%s", a.Target),
		fmt.Sprintf("-var-file=%s", a.VarFiles[0]),
		fmt.Sprintf("-var-file=%s", a.VarFiles[1]),
		fmt.Sprintf("-var=%s", a.Vars[0]),
		fmt.Sprintf("-var=%s", a.Vars[1]),
	)

	got := a.Command(t.Context())
	if got.Path != want.Path {
		t.Errorf("Command path is %v, want %v", got.Path, want.Path)
	}

	if !slices.Equal(got.Args, want.Args) {
		t.Errorf("Command args is %v, want %v", got.Args, want.Args)
	}
}

func TestTerraform_Apply_Command_tf13(t *testing.T) {
	// setup types
	v, err := semver.NewVersion("0.13.0")
	if err != nil {
		t.Error(err)
	}

	a := &Apply{
		AutoApprove: true,
		Backup:      "backup/",
		Directory:   "foobar/",
		Lock:        true,
		LockTimeout: 1 * time.Second,
		Input:       true,
		NoColor:     true,
		Parallelism: 1,
		Refresh:     true,
		State:       "state.tf",
		StateOut:    "stateout.tf",
		Target:      "target.tf",
		Vars:        []string{"foo=bar", "bar=foo"},
		VarFiles:    []string{"vars1.tf", "vars2.tf"},
		Version:     v,
	}

	//nolint:gosec // ignore G204
	want := exec.CommandContext(
		t.Context(),
		_terraform,
		applyAction,
		"-auto-approve",
		fmt.Sprintf("-backup=%s", a.Backup),
		"-lock=true",
		fmt.Sprintf("-lock-timeout=%s", a.LockTimeout),
		"-input=true",
		"-no-color",
		fmt.Sprintf("-parallelism=%d", a.Parallelism),
		"-refresh=true",
		fmt.Sprintf("-state=%s", a.State),
		fmt.Sprintf("-state-out=%s", a.StateOut),
		fmt.Sprintf("-target=%s", a.Target),
		fmt.Sprintf("-var-file=%s", a.VarFiles[0]),
		fmt.Sprintf("-var-file=%s", a.VarFiles[1]),
		fmt.Sprintf("-var=%s", a.Vars[0]),
		fmt.Sprintf("-var=%s", a.Vars[1]),
		fmt.Sprint(a.Directory),
	)

	got := a.Command(t.Context())
	if got.Path != want.Path {
		t.Errorf("Command path is %v, want %v", got.Path, want.Path)
	}

	if !slices.Equal(got.Args, want.Args) {
		t.Errorf("Command args is %v, want %v", got.Args, want.Args)
	}
}

func TestTerraform_Apply_Exec_Error(t *testing.T) {
	v, _ := semver.NewVersion("1.0.0")
	// setup types
	a := &Apply{
		Directory: "foobar/",
		Version:   v,
	}

	err := a.Exec(t.Context())
	if err == nil {
		t.Errorf("Exec should have returned err")
	}
}

func TestTerraform_Apply_Validate(t *testing.T) {
	v, _ := semver.NewVersion("1.0.0")
	// setup types
	tests := []struct {
		apply *Apply
	}{
		{
			apply: &Apply{Directory: "foobar/", Version: v},
		},
		{
			apply: &Apply{Directory: "foobar.tf", Version: v},
		},
		{
			apply: &Apply{Directory: "", Version: v},
		},
	}

	// run test
	for _, test := range tests {
		err := test.apply.Validate()
		if err != nil {
			t.Errorf("Validate returned err: %v", err)
		}
	}
}
