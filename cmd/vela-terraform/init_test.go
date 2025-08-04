// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os/exec"
	"slices"
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

	//nolint:gosec // ignore G204
	want := exec.CommandContext(
		t.Context(),
		_terraform,
		fmt.Sprintf("-chdir=%s", i.Directory),
		initAction,
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

	got := i.Command(t.Context())
	if got.Path != want.Path {
		t.Errorf("Command path is %v, want %v", got.Path, want.Path)
	}

	if !slices.Equal(got.Args, want.Args) {
		t.Errorf("Command args is %v, want %v", got.Args, want.Args)
	}
}

func TestTerraform_Init_Exec_Error(t *testing.T) {
	// setup types
	i := &Init{
		Directory:   "foobar/",
		InitOptions: &InitOptions{},
	}

	err := i.Exec(t.Context())
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
