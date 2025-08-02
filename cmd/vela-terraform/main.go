// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/mail"
	"os"

	"github.com/Masterminds/semver/v3"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"

	_ "github.com/joho/godotenv/autoload"

	"github.com/go-vela/vela-terraform/version"
)

//nolint:funlen // ignore function length due to comments and flags
func main() {
	// capture application version information
	v := version.New()

	// serialize the version information as pretty JSON
	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		logrus.Fatal(err)
	}

	// output the version information to stdout
	fmt.Fprintf(os.Stdout, "%s\n", string(bytes))

	// create new CLI application
	app := &cli.Command{
		Name:      "vela-terraform",
		Usage:     "Vela Terraform plugin for running Terraform",
		Copyright: "Copyright 2020 Target Brands, Inc. All rights reserved.",
		Authors: []any{
			&mail.Address{
				Name:    "Vela Admins",
				Address: "vela@target.com",
			},
		},
		Action:  run,
		Version: v.Semantic(),
		Flags: []cli.Flag{

			&cli.BoolFlag{
				Name:  "auto_approve",
				Usage: "skip interactive approval of running command",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_AUTO_APPROVE"),
					cli.EnvVar("TERRAFORM_AUTO_APPROVE"),
					cli.File("/vela/parameters/terraform/auto_approve"),
					cli.File("/vela/secrets/terraform/auto_approve"),
				),
			},
			&cli.StringFlag{
				Name:  "backup",
				Usage: "path to backup the existing state file",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_BACKUP"),
					cli.EnvVar("TERRAFORM_BACKUP"),
					cli.File("/vela/parameters/terraform/backup"),
					cli.File("/vela/secrets/terraform/backup"),
				),
			},
			&cli.StringFlag{
				Name:  "directory",
				Value: ".",
				Usage: "the directory for action to be performed on",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_DIRECTORY"),
					cli.EnvVar("TERRAFORM_DIRECTORY"),
					cli.File("/vela/parameters/terraform/directory"),
					cli.File("/vela/secrets/terraform/directory"),
				),
			},
			&cli.BoolFlag{
				Name:  "lock",
				Usage: "lock the state file when locking is supported",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_LOCK"),
					cli.EnvVar("TERRAFORM_LOCK"),
					cli.File("/vela/parameters/terraform/lock"),
					cli.File("/vela/secrets/terraform/lock"),
				),
			},
			&cli.DurationFlag{
				Name:  "lock_timeout",
				Usage: "duration to retry a state lock",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_LOCK_TIMEOUT"),
					cli.EnvVar("TERRAFORM_LOCK_TIMEOUT"),
					cli.File("/vela/parameters/terraform/lock_timeout"),
					cli.File("/vela/secrets/terraform/lock_timeout"),
				),
			},
			&cli.StringFlag{
				Name:  "log.level",
				Value: "info",
				Usage: "set log level - options: (trace|debug|info|warn|error|fatal|panic)",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_LOG_LEVEL"),
					cli.EnvVar("TERRAFORM_LOG_LEVEL"),
					cli.File("/vela/parameters/terraform/log_level"),
					cli.File("/vela/secrets/terraform/log_level"),
				),
			},
			&cli.BoolFlag{
				Name:  "no_color",
				Usage: "disables colors in output",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_NO_COLOR"),
					cli.EnvVar("TERRAFORM_NO_COLOR"),
					cli.File("/vela/parameters/terraform/no_color"),
					cli.File("/vela/secrets/terraform/no_color"),
				),
			},
			&cli.IntFlag{
				Name:  "parallelism",
				Usage: "number of concurrent operations as Terraform walks its graph",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_PARALLELISM"),
					cli.EnvVar("TERRAFORM_PARALLELISM"),
					cli.File("/vela/parameters/terraform/parallelism"),
					cli.File("/vela/secrets/terraform/parallelism"),
				),
			},
			&cli.BoolFlag{
				Name:  "refresh",
				Usage: "update state prior to checking for differences",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_REFRESH"),
					cli.EnvVar("TERRAFORM_REFRESH"),
					cli.File("/vela/parameters/terraform/refresh"),
					cli.File("/vela/secrets/terraform/refresh"),
				),
			},
			&cli.StringFlag{
				Name:  "state",
				Usage: "path to read and save state",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_STATE"),
					cli.EnvVar("TERRAFORM_STATE"),
					cli.File("/vela/parameters/terraform/state"),
					cli.File("/vela/secrets/terraform/state"),
				),
			},
			&cli.StringFlag{
				Name:  "state_out",
				Usage: "path to write updated state file",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_STATE_OUT"),
					cli.EnvVar("TERRAFORM_STATE_OUT"),
					cli.File("/vela/parameters/terraform/state_out"),
					cli.File("/vela/secrets/terraform/state_out"),
				),
			},
			&cli.StringFlag{
				Name:  "target",
				Usage: "resource to target",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_TARGET"),
					cli.EnvVar("TERRAFORM_TARGET"),
					cli.File("/vela/parameters/terraform/target"),
					cli.File("/vela/secrets/terraform/target"),
				),
			},
			&cli.StringSliceFlag{
				Name:  "vars",
				Usage: "a map of variables to pass to the Terraform (`<key>=<value>`)",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_VARS"),
					cli.EnvVar("TERRAFORM_VARS"),
					cli.File("/vela/parameters/terraform/vars"),
					cli.File("/vela/secrets/terraform/vars"),
				),
			},
			&cli.StringSliceFlag{
				Name:  "var_files",
				Usage: "a list of var files to use",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_VAR_FILES"),
					cli.EnvVar("TERRAFORM_VAR_FILES"),
					cli.File("/vela/parameters/terraform/var_files"),
					cli.File("/vela/secrets/terraform/var_files"),
				),
			},
			&cli.StringFlag{
				Name:  "terraform.version",
				Usage: "set terraform version for plugin",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_VERSION"),
					cli.EnvVar("TERRAFORM_VERSION"),
					cli.EnvVar("PLUGIN_TERRAFORM_VERSION"),
					cli.File("/vela/parameters/terraform/version"),
					cli.File("/vela/secrets/terraform/version"),
				),
			},

			// Config Flags

			&cli.StringFlag{
				Name:  "config.action",
				Usage: "the action to have terraform perform",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_ACTION"),
					cli.EnvVar("TERRAFORM_ACTION"),
					cli.File("/vela/parameters/terraform/action"),
					cli.File("/vela/secrets/terraform/action"),
				),
			},

			// FMT Flags

			&cli.BoolFlag{
				Name:  "fmt.check",
				Usage: "validate if the input is formatted",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_CHECK"),
					cli.EnvVar("TERRAFORM_CHECK"),
					cli.File("/vela/parameters/terraform/check"),
					cli.File("/vela/secrets/terraform/check"),
				),
			},
			&cli.BoolFlag{
				Name:  "fmt.diff",
				Usage: "diffs of formatting changes",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_DIFF"),
					cli.EnvVar("TERRAFORM_DIFF"),
					cli.File("/vela/parameters/terraform/diff"),
					cli.File("/vela/secrets/terraform/diff"),
				),
			},
			&cli.BoolFlag{
				Name:  "fmt.list",
				Usage: "list files whose formatting differs",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_LIST"),
					cli.EnvVar("TERRAFORM_LIST"),
					cli.File("/vela/parameters/terraform/list"),
					cli.File("/vela/secrets/terraform/list"),
				),
			},
			&cli.BoolFlag{
				Name:  "fmt.write",
				Usage: "write result to source file instead of STDOUT",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_WRITE"),
					cli.EnvVar("TERRAFORM_WRITE"),
					cli.File("/vela/parameters/terraform/write"),
					cli.File("/vela/secrets/terraform/write"),
				),
			},

			// InitOptions Flags

			&cli.StringFlag{
				Name:  "init.options",
				Usage: "properties to set on terraform init action",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_INIT_OPTIONS"),
					cli.EnvVar("TERRAFORM_INIT_OPTIONS"),
					cli.File("/vela/parameters/terraform/init_options"),
					cli.File("/vela/secrets/terraform/init_options"),
				),
			},

			// Netrc Flags

			&cli.StringFlag{
				Name:  "netrc.machine",
				Value: "github.com",
				Usage: "remote machine name to communicate with",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_MACHINE"),
					cli.EnvVar("TERRAFORM_MACHINE"),
					cli.EnvVar("VELA_NETRC_MACHINE"),
					cli.File("/vela/parameters/terraform/machine"),
					cli.File("/vela/secrets/terraform/machine"),
				),
			},
			&cli.StringFlag{
				Name:  "netrc.username",
				Usage: "user name for communication with the remote machine",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_USERNAME"),
					cli.EnvVar("TERRAFORM_USERNAME"),
					cli.EnvVar("VELA_NETRC_USERNAME"),
					cli.File("/vela/parameters/terraform/username"),
					cli.File("/vela/secrets/terraform/username"),
				),
			},
			&cli.StringFlag{
				Name:  "netrc.password",
				Usage: "password for communication with the remote machine",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_PASSWORD"),
					cli.EnvVar("TERRAFORM_PASSWORD"),
					cli.EnvVar("VELA_NETRC_PASSWORD"),
					cli.File("/vela/parameters/terraform/password"),
					cli.File("/vela/secrets/terraform/password"),
				),
			},

			// Plan Flags

			&cli.BoolFlag{
				Name:  "plan.destroy",
				Usage: "destroy all resources managed by the given configuration and state",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_DESTROY"),
					cli.EnvVar("TERRAFORM_DESTROY"),
					cli.File("/vela/parameters/terraform/destroy"),
					cli.File("/vela/secrets/terraform/destroy"),
				),
			},
			&cli.BoolFlag{
				Name:  "plan.detailed_exit_code",
				Usage: "return detailed exit codes when the command exits",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_DETAILED_EXIT_CODE"),
					cli.EnvVar("TERRAFORM_DETAILED_EXIT_CODE"),
					cli.File("/vela/parameters/terraform/detailed_exit_code"),
					cli.File("/vela/secrets/terraform/detailed_exit_code"),
				),
			},
			&cli.BoolFlag{
				Name:  "plan.input",
				Usage: "ask for input for variables if not directly set",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_INPUT"),
					cli.EnvVar("TERRAFORM_INPUT"),
					cli.File("/vela/parameters/terraform/input"),
					cli.File("/vela/secrets/terraform/input"),
				),
			},
			&cli.IntFlag{
				Name:  "plan.module_depth",
				Usage: "specifies the depth of modules to show in the output",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_MODULE_DEPTH"),
					cli.EnvVar("TERRAFORM_MODULE_DEPTH"),
					cli.File("/vela/parameters/terraform/module_depth"),
					cli.File("/vela/secrets/terraform/module_depth"),
				),
			},

			// Validation Flags

			&cli.BoolFlag{
				Name:  "validation.check_variables",
				Usage: "command will check whether all required variables have been specified",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_CHECK_VARIABLES"),
					cli.EnvVar("TERRAFORM_CHECK_VARIABLES"),
					cli.File("/vela/parameters/terraform/check_variables"),
					cli.File("/vela/secrets/terraform/check_variables"),
				),
			},
		},
	}

	err = app.Run(context.Background(), os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}

// run executes the plugin based off the configuration provided.
func run(ctx context.Context, cmd *cli.Command) error {
	// set the log level for the plugin
	switch cmd.String("log.level") {
	case "t", "trace", "Trace", "TRACE":
		logrus.SetLevel(logrus.TraceLevel)
	case "d", "debug", "Debug", "DEBUG":
		logrus.SetLevel(logrus.DebugLevel)
	case "w", "warn", "Warn", "WARN":
		logrus.SetLevel(logrus.WarnLevel)
	case "e", "error", "Error", "ERROR":
		logrus.SetLevel(logrus.ErrorLevel)
	case "f", "fatal", "Fatal", "FATAL":
		logrus.SetLevel(logrus.FatalLevel)
	case "p", "panic", "Panic", "PANIC":
		logrus.SetLevel(logrus.PanicLevel)
	case "i", "info", "Info", "INFO":
		fallthrough
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}

	logrus.WithFields(logrus.Fields{
		"code":     "https://github.com/go-vela/vela-terraform",
		"docs":     "https://go-vela.github.io/docs/plugins/registry/pipeline/terraform",
		"registry": "https://hub.docker.com/r/target/vela-terraform",
	}).Info("Vela Terraform Plugin")

	// capture custom terraform tfVersion requested
	tfVersion := cmd.String("terraform.version")

	// attempt to install the custom terraform tfVersion if different from default
	err := installBinary(ctx, tfVersion, os.Getenv("PLUGIN_TERRAFORM_VERSION"))
	if err != nil {
		return err
	}

	tfSemVersion, err := semver.NewVersion(tfVersion)
	if err != nil {
		logrus.Errorf("Unable to parse terraform version")
		return err
	}

	// create the plugin
	p := Plugin{
		// Apply configuration
		Apply: &Apply{
			AutoApprove: cmd.Bool("auto_approve"),
			Backup:      cmd.String("backup"),
			Directory:   cmd.String("directory"),
			Lock:        cmd.Bool("lock"),
			LockTimeout: cmd.Duration("lock_timeout"),
			NoColor:     cmd.Bool("no_color"),
			Parallelism: cmd.Int("parallelism"),
			Refresh:     cmd.Bool("refresh"),
			State:       cmd.String("state"),
			StateOut:    cmd.String("state_out"),
			Target:      cmd.String("target"),
			Vars:        cmd.StringSlice("vars"),
			VarFiles:    cmd.StringSlice("var_files"),
			Version:     tfSemVersion,
		},
		// Config configuration
		Config: &Config{
			Action: cmd.String("config.action"),
			Netrc: &Netrc{
				Login:    cmd.String("netrc.username"),
				Machine:  cmd.String("netrc.machine"),
				Password: cmd.String("netrc.password"),
			},
		},
		// Destroy configuration
		Destroy: &Destroy{
			AutoApprove: cmd.Bool("auto_approve"),
			Backup:      cmd.String("backup"),
			Directory:   cmd.String("directory"),
			Lock:        cmd.Bool("lock"),
			LockTimeout: cmd.Duration("lock_timeout"),
			NoColor:     cmd.Bool("no_color"),
			Parallelism: cmd.Int("parallelism"),
			Refresh:     cmd.Bool("refresh"),
			State:       cmd.String("state"),
			StateOut:    cmd.String("state_out"),
			Target:      cmd.String("target"),
			Vars:        cmd.StringSlice("vars"),
			VarFiles:    cmd.StringSlice("var_files"),
			Version:     tfSemVersion,
		},
		// FMT configuration
		FMT: &FMT{
			Check:     cmd.Bool("fmt.check"),
			Diff:      cmd.Bool("fmt.diff"),
			Directory: cmd.String("directory"),
			List:      cmd.Bool("fmt.list"),
			Write:     cmd.Bool("fmt.write"),
			Version:   tfSemVersion,
		},
		// InitOptions configuration
		Init: &Init{
			Directory: cmd.String("directory"),
			RawInit:   cmd.String("init.options"),
		},
		// Plan configuration
		Plan: &Plan{
			Destroy:          cmd.Bool("plan.destroy"),
			DetailedExitCode: cmd.Bool("plan.detailed_exit_code"),
			Directory:        cmd.String("directory"),
			Input:            cmd.Bool("plan.input"),
			Lock:             cmd.Bool("lock"),
			LockTimeout:      cmd.Duration("lock_timeout"),
			ModuleDepth:      cmd.Int("plan.module_depth"),
			NoColor:          cmd.Bool("no_color"),
			Parallelism:      cmd.Int("parallelism"),
			Refresh:          cmd.Bool("refresh"),
			State:            cmd.String("state"),
			Out:              cmd.String("state_out"),
			Target:           cmd.String("target"),
			Vars:             cmd.StringSlice("vars"),
			VarFiles:         cmd.StringSlice("var_files"),
			Version:          tfSemVersion,
		},
		// Validation configuration
		Validation: &Validation{
			CheckVariables: cmd.Bool("validation.check_variables"),
			Directory:      cmd.String("directory"),
			NoColor:        cmd.Bool("no_color"),
			Vars:           cmd.StringSlice("vars"),
			VarFiles:       cmd.StringSlice("var_files"),
			Version:        tfSemVersion,
		},
	}

	// validate the plugin
	err = p.Validate()
	if err != nil {
		return err
	}

	return p.Exec(ctx)
}

func SupportsChdir(v *semver.Version) bool {
	if v.Major() >= 1 {
		return true
	}

	if v.Minor() >= 14 {
		return true
	}

	return false
}
