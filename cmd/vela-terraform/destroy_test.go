// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"
	"os/exec"
	"reflect"
	"testing"
	"time"
)

func TestTerraform_Destroy_Command(t *testing.T) {
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
	}

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
		fmt.Sprintf("-var=%s -var=%s", d.Vars[0], d.Vars[1]),
		fmt.Sprintf("-var-file=%s -var-file=%s", d.VarFiles[0], d.VarFiles[1]),
		d.Directory,
	)

	got := d.Command("foobar/")
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Command is %v, want %v", got, want)
	}
}

func TestTerraform_Destroy_Exec_Error(t *testing.T) {
	// setup types
	d := &Destroy{
		Directory: "foobar/",
	}

	err := d.Exec()
	if err == nil {
		t.Errorf("Exec should have returned err")
	}
}

func TestTerraform_Destroy_Validate(t *testing.T) {
	// setup types
	tests := []struct {
		destroy *Destroy
	}{
		{
			destroy: &Destroy{Directory: "foobar/"},
		},
		{
			destroy: &Destroy{Directory: "foobar.tf"},
		},
		{
			destroy: &Destroy{Directory: ""},
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
