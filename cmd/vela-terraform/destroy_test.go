// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os/exec"
	"reflect"
	"testing"
	"time"

	"github.com/Masterminds/semver/v3"
)

func TestTerraform_Destroy_Command(t *testing.T) {
	v, _ := semver.NewVersion("1.0.0")
	// setup types
	d := &Destroy{
		AutoApprove: true,
		Backup:      "backup/",
		Directory:   "foobar/",
		Lock:        true,
		LockTimeout: 1 * time.Second,
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

	// nolint: gosec // ignore G204
	want := exec.Command(
		_terraform,
		fmt.Sprintf("-chdir=%s", d.Directory),
		destroyAction,
		"-auto-approve",
		fmt.Sprintf("-backup=%s", d.Backup),
		"-lock=true",
		fmt.Sprintf("-lock-timeout=%s", d.LockTimeout),
		"-no-color",
		fmt.Sprintf("-parallelism=%d", d.Parallelism),
		"-refresh=true",
		fmt.Sprintf("-state=%s", d.State),
		fmt.Sprintf("-state-out=%s", d.StateOut),
		fmt.Sprintf("-target=%s", d.Target),
		fmt.Sprintf("-var=%s", d.Vars[0]),
		fmt.Sprintf("-var=%s", d.Vars[1]),
		fmt.Sprintf("-var-file=%s", d.VarFiles[0]),
		fmt.Sprintf("-var-file=%s", d.VarFiles[1]),
	)

	got := d.Command()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Command is %v, want %v", got, want)
	}
}

func TestTerraform_Destroy_Command_tf13(t *testing.T) {
	v, _ := semver.NewVersion("0.13.0")
	// setup types
	d := &Destroy{
		AutoApprove: true,
		Backup:      "backup/",
		Directory:   "foobar/",
		Lock:        true,
		LockTimeout: 1 * time.Second,
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

	// nolint: gosec // ignore G204
	want := exec.Command(
		_terraform,
		destroyAction,
		"-auto-approve",
		fmt.Sprintf("-backup=%s", d.Backup),
		"-lock=true",
		fmt.Sprintf("-lock-timeout=%s", d.LockTimeout),
		"-no-color",
		fmt.Sprintf("-parallelism=%d", d.Parallelism),
		"-refresh=true",
		fmt.Sprintf("-state=%s", d.State),
		fmt.Sprintf("-state-out=%s", d.StateOut),
		fmt.Sprintf("-target=%s", d.Target),
		fmt.Sprintf("-var=%s", d.Vars[0]),
		fmt.Sprintf("-var=%s", d.Vars[1]),
		fmt.Sprintf("-var-file=%s", d.VarFiles[0]),
		fmt.Sprintf("-var-file=%s", d.VarFiles[1]),
		fmt.Sprintf(d.Directory),
	)

	got := d.Command()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Command is %v, want %v", got, want)
	}
}

func TestTerraform_Destroy_Exec_Error(t *testing.T) {
	v, _ := semver.NewVersion("1.0.0")
	// setup types
	d := &Destroy{
		Directory: "foobar/",
		Version:   v,
	}

	err := d.Exec()
	if err == nil {
		t.Errorf("Exec should have returned err")
	}
}

func TestTerraform_Destroy_Validate(t *testing.T) {
	v, _ := semver.NewVersion("1.0.0")
	// setup types
	tests := []struct {
		destroy *Destroy
	}{
		{
			destroy: &Destroy{Directory: "foobar/", Version: v},
		},
		{
			destroy: &Destroy{Directory: "foobar.tf", Version: v},
		},
		{
			destroy: &Destroy{Directory: "", Version: v},
		},
	}

	// run test
	for _, test := range tests {
		err := test.destroy.Validate()
		if err != nil {
			t.Errorf("Validate returned err: %v", err)
		}
	}
}
