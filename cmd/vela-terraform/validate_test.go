// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"
	"os/exec"
	"reflect"
	"testing"
)

func TestTerraform_Validation_Command(t *testing.T) {
	// setup types
	v := &Validation{
		CheckVariables: true,
		Directory:      "foobar/",
		NoColor:        true,
		Vars:           []string{"foo=bar", "bar=foo"},
		VarFiles:       []string{"vars1.tf", "vars2.tf"},
	}

	want := exec.Command(
		_terraform,
		validationAction,
		fmt.Sprintf("-check-variables=%t", v.CheckVariables),
		"-no-color",
		fmt.Sprintf("-var=\"%s %s\"", v.Vars[0], v.Vars[1]),
		fmt.Sprintf("-var-file=%s -var-file=%s", v.VarFiles[0], v.VarFiles[1]),
		v.Directory,
	)

	got := v.Command("foobar/")
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Command is %v, want %v", got, want)
	}
}

func TestTerraform_Validation_Exec(t *testing.T) {
	// setup types
	v := &Validation{
		Directory: "foobar/",
	}

	err := v.Exec()
	if err == nil {
		t.Errorf("Exec should have returned err")
	}
}

func TestTerraform_Validation_Validate(t *testing.T) {
	// setup types
	tests := []struct {
		validation *Validation
	}{
		{
			validation: &Validation{Directory: "foobar/"},
		},
		{
			validation: &Validation{Directory: "foobar.tf"},
		},
		{
			validation: &Validation{Directory: ""},
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
