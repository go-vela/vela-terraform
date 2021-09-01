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

func TestTerraform_Init_Command(t *testing.T) {
	// setup types
	i := &Init{
		Directory: "foobar/",
		InitOptions: &InitOptions{
			Backend:        true,
			BackendConfigs: []string{"config1.tf", "config1.tf"},
			ForceCopy:      true,
			FromModule:     "module",
			Get:            true,
			GetPlugins:     true,
			Lock:           true,
			LockTimeout:    1 * time.Second,
			Input:          true,
			NoColor:        true,
			PluginDirs:     []string{"plugin1/", "plugin2/"},
			Reconfigure:    true,
			Upgrade:        true,
			VerifyPlugins:  true,
		},
	}

	want := exec.Command(
		_terraform,
		initAction,
		fmt.Sprintf("-chdir=%s", i.Directory),
		"-backend=true",
		fmt.Sprintf("-backend-config=%s", i.InitOptions.BackendConfigs[0]),
		fmt.Sprintf("-backend-config=%s", i.InitOptions.BackendConfigs[1]),
		"-force-copy",
		fmt.Sprintf("-from-module=%s", i.InitOptions.FromModule),
		"-get=true",
		"-get-plugins=true",
		"-input=true",
		"-lock=true",
		fmt.Sprintf("-lock-timeout=%s", i.InitOptions.LockTimeout),
		"-no-color",
		fmt.Sprintf("-plugin-dir=%s", i.InitOptions.PluginDirs[0]),
		fmt.Sprintf("-plugin-dir=%s", i.InitOptions.PluginDirs[1]),
		"-reconfigure",
		"-upgrade=false",
		"-verify-plugins=true",
	)

	got := i.Command()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Command is %v, want %v", got, want)
	}
}

func TestTerraform_Init_Exec_Error(t *testing.T) {
	// setup types
	i := &Init{
		Directory:   "foobar/",
		InitOptions: &InitOptions{},
	}

	err := i.Exec()
	if err == nil {
		t.Errorf("Exec should have returned err")
	}
}

func TestTerraform_Init_Validate(t *testing.T) {
	// setup types
	tests := []struct {
		init *Init
	}{
		{
			init: &Init{Directory: "foobar/"},
		},
		{
			init: &Init{Directory: "foobar.tf"},
		},
		{
			init: &Init{Directory: ""},
		},
	}

	// run test
	for _, test := range tests {
		err := test.init.Validate()
		if err != nil {
			t.Errorf("Validate returned err: %v", err)
		}
	}
}
