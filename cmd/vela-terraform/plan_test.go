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

func TestTerraform_Plan_Command(t *testing.T) {
	// setup types
	p := &Plan{
		Directory:        "foobar/",
		Destroy:          true,
		DetailedExitCode: true,
		Input:            true,
		Lock:             true,
		LockTimeout:      1 * time.Second,
		ModuleDepth:      1,
		NoColor:          true,
		Out:              "/path/to/out.tf",
		Parallelism:      1,
		Refresh:          true,
		State:            "state.tf",
		Target:           "target.tf",
		Var:              []string{"foo=bar", "bar=foo"},
		VarFile:          "vars.tf",
	}

	want := exec.Command(
		_terraform,
		planAction,
		"-destroy",
		"-detailed-exitcode",
		"-input=true",
		"-lock=true",
		fmt.Sprintf("-lock-timeout=%s", p.LockTimeout),
		fmt.Sprintf("-module-depth=%d", p.ModuleDepth),
		"-no-color",
		fmt.Sprintf("-out=%s", p.Out),
		fmt.Sprintf("-parallelism=%d", p.Parallelism),
		"-refresh=true",
		fmt.Sprintf("-state=%s", p.State),
		fmt.Sprintf("-target=%s", p.Target),
		fmt.Sprintf("-var=\"%s %s\"", p.Var[0], p.Var[1]),
		fmt.Sprintf("-var-file=%s", p.VarFile),
		p.Directory,
	)

	got := p.Command("foobar/")
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Command is %v, want %v", got, want)
	}
}

func TestTerraform_Plan_Exec(t *testing.T) {
	// setup types
	v := &Validation{
		Directory: "foobar/",
	}

	err := v.Exec()
	if err == nil {
		t.Errorf("Exec should have returned err")
	}
}

func TestTerraform_Plan_Validate(t *testing.T) {
	// setup types
	tests := []struct {
		plan *Plan
	}{
		{
			plan: &Plan{Directory: "foobar/"},
		},
		{
			plan: &Plan{Directory: "foobar.tf"},
		},
		{
			plan: &Plan{Directory: ""},
		},
	}

	// run test
	for _, test := range tests {
		err := test.plan.Validate()
		if err != nil {
			t.Errorf("Plan returned err: %v", err)
		}
	}
}
