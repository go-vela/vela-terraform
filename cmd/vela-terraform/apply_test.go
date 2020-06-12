// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
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

func TestTerraform_Apply_Command(t *testing.T) {
	// setup types
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
	}

	want := exec.Command(
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
		fmt.Sprintf("-var=\"%s %s\"", a.Vars[0], a.Vars[1]),
		fmt.Sprintf("-var-file=%s -var-file=%s", a.VarFiles[0], a.VarFiles[1]),
		a.Directory,
	)

	got := a.Command("foobar/")
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Command is %v, want %v", got, want)
	}
}

func TestTerraform_Apply_Exec_Error(t *testing.T) {
	// setup types
	a := &Apply{
		Directory: "foobar/",
	}

	err := a.Exec()
	if err == nil {
		t.Errorf("Exec should have returned err")
	}
}

func TestTerraform_Apply_Validate(t *testing.T) {
	// setup types
	tests := []struct {
		apply *Apply
	}{
		{
			apply: &Apply{Directory: "foobar/"},
		},
		{
			apply: &Apply{Directory: "foobar.tf"},
		},
		{
			apply: &Apply{Directory: ""},
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
