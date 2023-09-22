// SPDX-License-Identifier: Apache-2.0

package main

import (
	"testing"
	"time"
)

func TestTerraform_Plugin_Validate(t *testing.T) {
	// setup tests
	tests := []struct {
		plugin *Plugin
		want   *error
	}{
		{ // test success for apply action
			plugin: &Plugin{
				Apply: &Apply{
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
				},
				Config: &Config{
					Action: "apply",
					Netrc: &Netrc{
						Machine:  "machine.example.com",
						Login:    "octocat",
						Password: "foobar",
					},
				},
				Destroy: &Destroy{},
				Init: &Init{
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
				},
				FMT:        &FMT{},
				Plan:       &Plan{},
				Validation: &Validation{},
			},
			want: nil,
		},
		{ // test success for destroy action
			plugin: &Plugin{
				Apply: &Apply{},
				Config: &Config{
					Action: "apply",
					Netrc: &Netrc{
						Machine:  "machine.example.com",
						Login:    "octocat",
						Password: "foobar",
					},
				},
				Destroy: &Destroy{
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
				},
				Init: &Init{
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
				},
				FMT:        &FMT{},
				Plan:       &Plan{},
				Validation: &Validation{},
			},
			want: nil,
		},
		{ // test success for fmt action
			plugin: &Plugin{
				Apply: &Apply{},
				Config: &Config{
					Action: "apply",
					Netrc: &Netrc{
						Machine:  "machine.example.com",
						Login:    "octocat",
						Password: "foobar",
					},
				},
				Destroy: &Destroy{},
				Init: &Init{
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
				},
				FMT: &FMT{
					Check:     true,
					Diff:      true,
					Directory: "foobar/",
					List:      false,
					Write:     false,
				},
				Plan:       &Plan{},
				Validation: &Validation{},
			},
			want: nil,
		},
		{ // test success for plan action
			plugin: &Plugin{
				Apply: &Apply{},
				Config: &Config{
					Action: "apply",
					Netrc: &Netrc{
						Machine:  "machine.example.com",
						Login:    "octocat",
						Password: "foobar",
					},
				},
				Destroy: &Destroy{},
				Init: &Init{
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
				},
				FMT: &FMT{},
				Plan: &Plan{
					Destroy:          true,
					DetailedExitCode: true,
					Directory:        "foobar/",
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
					Vars:             []string{"foo=bar", "bar=foo"},
					VarFiles:         []string{"vars1.tf", "vars2.tf"},
				},
				Validation: &Validation{},
			},
			want: nil,
		},
		{ // test success for validation action
			plugin: &Plugin{
				Apply: &Apply{},
				Config: &Config{
					Action: "apply",
					Netrc: &Netrc{
						Machine:  "machine.example.com",
						Login:    "octocat",
						Password: "foobar",
					},
				},
				Destroy: &Destroy{},
				Init: &Init{
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
				},
				FMT:  &FMT{},
				Plan: &Plan{},
				Validation: &Validation{
					CheckVariables: false,
					Directory:      "foobar/",
					NoColor:        true,
					Vars:           []string{"foo=bar", "bar=foo"},
					VarFiles:       []string{"vars1.tf", "vars2.tf"},
				},
			},
			want: nil,
		},
	}

	// run tests
	for _, test := range tests {
		err := test.plugin.Validate()
		if err != nil {
			t.Errorf("Validate returned err: %v", err)
		}
	}
}
