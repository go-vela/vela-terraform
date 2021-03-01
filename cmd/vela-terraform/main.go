// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-vela/vela-terraform/version"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	_ "github.com/joho/godotenv/autoload"
)

// nolint: funlen // ignore function length due to comments and flags
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
	app := cli.NewApp()

	// Plugin Information

	app.Name = "vela-terraform"
	app.HelpName = "vela-terraform"
	app.Usage = "Vela Terraform plugin for running Terraform"
	app.Copyright = "Copyright (c) 2021 Target Brands, Inc. All rights reserved."
	app.Authors = []*cli.Author{
		{
			Name:  "Vela Admins",
			Email: "vela@target.com",
		},
	}

	// Plugin Metadata

	app.Action = run
	app.Compiled = time.Now()
	app.Version = v.Semantic()

	// Plugin Flags

	app.Flags = []cli.Flag{

		&cli.BoolFlag{
			EnvVars:  []string{"PARAMETER_AUTO_APPROVE", "TERRAFORM_AUTO_APPROVE"},
			FilePath: "/vela/parameters/terraform/auto_approve,/vela/secrets/terraform/auto_approve",
			Name:     "auto_approve",
			Usage:    "skip interactive approval of running command",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_BACKUP", "TERRAFORM_BACKUP"},
			FilePath: "/vela/parameters/terraform/backup,/vela/secrets/terraform/backup",
			Name:     "backup",
			Usage:    "path to backup the existing state file",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_DIRECTORY", "TERRAFORM_DIRECTORY"},
			FilePath: "/vela/parameters/terraform/directory,/vela/secrets/terraform/directory",
			Name:     "directory",
			Usage:    "the directory for action to be performed on",
			Value:    ".",
		},
		&cli.BoolFlag{
			EnvVars:  []string{"PARAMETER_LOCK", "TERRAFORM_LOCK"},
			FilePath: "/vela/parameters/terraform/lock,/vela/secrets/terraform/lock",
			Name:     "lock",
			Usage:    "lock the state file when locking is supported",
		},
		&cli.DurationFlag{
			EnvVars:  []string{"PARAMETER_LOCK_TIMEOUT", "TERRAFORM_LOCK_TIMEOUT"},
			FilePath: "/vela/parameters/terraform/lock_timeout,/vela/secrets/terraform/lock_timeout",
			Name:     "lock_timeout",
			Usage:    "duration to retry a state lock",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_LOG_LEVEL", "TERRAFORM_LOG_LEVEL"},
			FilePath: "/vela/parameters/terraform/log_level,/vela/secrets/terraform/log_level",
			Name:     "log.level",
			Usage:    "set log level - options: (trace|debug|info|warn|error|fatal|panic)",
			Value:    "info",
		},
		&cli.BoolFlag{
			EnvVars:  []string{"PARAMETER_NO_COLOR", "TERRAFORM_NO_COLOR"},
			FilePath: "/vela/parameters/terraform/no_color,/vela/secrets/terraform/no_color",
			Name:     "no_color",
			Usage:    "disables colors in output",
		},
		&cli.IntFlag{
			EnvVars:  []string{"PARAMETER_PARALLELISM", "TERRAFORM_PARALLELISM"},
			FilePath: "/vela/parameters/terraform/parallelism,/vela/secrets/terraform/parallelism",
			Name:     "parallelism",
			Usage:    "number of concurrent operations as Terraform walks its graph",
		},
		&cli.BoolFlag{
			EnvVars:  []string{"PARAMETER_REFRESH", "TERRAFORM_REFRESH"},
			FilePath: "/vela/parameters/terraform/refresh,/vela/secrets/terraform/refresh",
			Name:     "refresh",
			Usage:    "update state prior to checking for differences",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_STATE", "TERRAFORM_STATE"},
			FilePath: "/vela/parameters/terraform/state,/vela/secrets/terraform/state",
			Name:     "state",
			Usage:    "path to read and save state",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_STATE_OUT", "TERRAFORM_STATE_OUT"},
			FilePath: "/vela/parameters/terraform/state_out,/vela/secrets/terraform/state_out",
			Name:     "state_out",
			Usage:    "path to write updated state file",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_TARGET", "TERRAFORM_TARGET"},
			FilePath: "/vela/parameters/terraform/target,/vela/secrets/terraform/target",
			Name:     "target",
			Usage:    "resource to target",
		},
		&cli.StringSliceFlag{
			EnvVars:  []string{"PARAMETER_VARS", "TERRAFORM_VARS"},
			FilePath: "/vela/parameters/terraform/vars,/vela/secrets/terraform/vars",
			Name:     "vars",
			Usage:    "a map of variables to pass to the Terraform (`<key>=<value>`)",
		},
		&cli.StringSliceFlag{
			EnvVars:  []string{"PARAMETER_VAR_FILES", "TERRAFORM_VAR_FILES"},
			FilePath: "/vela/parameters/terraform/var_files,/vela/secrets/terraform/var_files",
			Name:     "var_files",
			Usage:    "a list of var files to use",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_VERSION", "TERRAFORM_VERSION"},
			FilePath: "/vela/parameters/terraform/version,/vela/secrets/terraform/version",
			Name:     "terraform.version",
			Usage:    "set terraform version for plugin",
		},

		// Config Flags

		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_ACTION", "TERRAFORM_ACTION"},
			FilePath: "/vela/parameters/terraform/action,/vela/secrets/terraform/action",
			Name:     "config.action",
			Usage:    "the action to have terraform perform",
		},

		// FMT Flags

		&cli.BoolFlag{
			EnvVars:  []string{"PARAMETER_CHECK", "TERRAFORM_CHECK"},
			FilePath: "/vela/parameters/terraform/check,/vela/secrets/terraform/check",
			Name:     "fmt.check",
			Usage:    "validate if the input is formatted",
		},
		&cli.BoolFlag{
			EnvVars:  []string{"PARAMETER_DIFF", "TERRAFORM_DIFF"},
			FilePath: "/vela/parameters/terraform/diff,/vela/secrets/terraform/diff",
			Name:     "fmt.diff",
			Usage:    "diffs of formatting changes",
		},
		&cli.BoolFlag{
			EnvVars:  []string{"PARAMETER_LIST", "TERRAFORM_LIST"},
			FilePath: "/vela/parameters/terraform/list,/vela/secrets/terraform/list",
			Name:     "fmt.list",
			Usage:    "list files whose formatting differs",
		},
		&cli.BoolFlag{
			EnvVars:  []string{"PARAMETER_WRITE", "TERRAFORM_WRITE"},
			FilePath: "/vela/parameters/terraform/write,/vela/secrets/terraform/write",
			Name:     "fmt.write",
			Usage:    "write result to source file instead of STDOUT",
		},

		// InitOptions Flags

		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_INIT_OPTIONS", "TERRAFORM_INIT_OPTIONS"},
			FilePath: "/vela/parameters/terraform/init_options,/vela/secrets/terraform/init_options",
			Name:     "init.options",
			Usage:    "properties to set on terraform init action",
		},

		// Netrc Flags

		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_MACHINE", "TERRAFORM_MACHINE", "VELA_NETRC_MACHINE"},
			FilePath: "/vela/parameters/terraform/machine,/vela/secrets/terraform/machine",
			Name:     "netrc.machine",
			Usage:    "remote machine name to communicate with",
			Value:    "github.com",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_USERNAME", "TERRAFORM_USERNAME", "VELA_NETRC_USERNAME"},
			FilePath: "/vela/parameters/terraform/username,/vela/secrets/terraform/username",
			Name:     "netrc.username",
			Usage:    "user name for communication with the remote machine",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_PASSWORD", "TERRAFORM_PASSWORD", "VELA_NETRC_PASSWORD"},
			FilePath: "/vela/parameters/terraform/password,/vela/secrets/terraform/password",
			Name:     "netrc.password",
			Usage:    "password for communication with the remote machine",
		},

		// Plan Flags

		&cli.BoolFlag{
			EnvVars:  []string{"PARAMETER_DESTROY", "TERRAFORM_DESTROY"},
			FilePath: "/vela/parameters/terraform/destroy,/vela/secrets/terraform/destroy",
			Name:     "plan.destroy",
			Usage:    "destroy all resources managed by the given configuration and state",
		},
		// nolint: lll // skip line length due to long parameter name
		&cli.BoolFlag{
			EnvVars:  []string{"PARAMETER_DETAILED_EXIT_CODE", "TERRAFORM_DETAILED_EXIT_CODE"},
			FilePath: "/vela/parameters/terraform/detailed_exit_code,/vela/secrets/terraform/detailed_exit_code",
			Name:     "plan.detailed_exit_code",
			Usage:    "return detailed exit codes when the command exits",
		},
		&cli.BoolFlag{
			EnvVars:  []string{"PARAMETER_INPUT", "TERRAFORM_INPUT"},
			FilePath: "/vela/parameters/terraform/input,/vela/secrets/terraform/input",
			Name:     "plan.input",
			Usage:    "ask for input for variables if not directly set",
		},
		&cli.IntFlag{
			EnvVars:  []string{"PARAMETER_MODULE_DEPTH", "TERRAFORM_MODULE_DEPTH"},
			FilePath: "/vela/parameters/terraform/module_depth,/vela/secrets/terraform/module_depth",
			Name:     "plan.module_depth",
			Usage:    "specifies the depth of modules to show in the output",
		},

		// Validation Flags

		&cli.BoolFlag{
			EnvVars:  []string{"PARAMETER_CHECK_VARIABLES", "TERRAFORM_CHECK_VARIABLES"},
			FilePath: "/vela/parameters/terraform/check_variables,/vela/secrets/terraform/check_variables",
			Name:     "validation.check_variables",
			Usage:    "command will check whether all required variables have been specified",
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}

// run executes the plugin based off the configuration provided.
//
// nolint: funlen // ignore function length due to comments and flags
func run(c *cli.Context) error {
	// set the log level for the plugin
	switch c.String("log.level") {
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
		"docs":     "https://go-vela.github.io/docs/plugins/registry/terraform",
		"registry": "https://hub.docker.com/r/target/vela-terraform",
	}).Info("Vela Terraform Plugin")

	// capture custom terraform version requested
	version := c.String("terraform.version")

	// check if a custom terraform version was requested
	if len(version) > 0 {
		// attempt to install the custom terraform version
		err := install(version, os.Getenv("PLUGIN_TERRAFORM_VERSION"))
		if err != nil {
			return err
		}
	}

	// create the plugin
	p := Plugin{
		// Apply configuration
		Apply: &Apply{
			AutoApprove: c.Bool("auto_approve"),
			Backup:      c.String("backup"),
			Directory:   c.String("directory"),
			Lock:        c.Bool("lock"),
			LockTimeout: c.Duration("lock_timeout"),
			NoColor:     c.Bool("no_color"),
			Parallelism: c.Int("parallelism"),
			Refresh:     c.Bool("refresh"),
			State:       c.String("state"),
			StateOut:    c.String("state_out"),
			Target:      c.String("target"),
			Vars:        c.StringSlice("vars"),
			VarFiles:    c.StringSlice("var_files"),
		},
		// Config configuration
		Config: &Config{
			Action: c.String("config.action"),
			Netrc: &Netrc{
				Login:    c.String("netrc.username"),
				Machine:  c.String("netrc.machine"),
				Password: c.String("netrc.password"),
			},
		},
		// Destroy configuration
		Destroy: &Destroy{
			AutoApprove: c.Bool("auto_approve"),
			Backup:      c.String("backup"),
			Directory:   c.String("directory"),
			Lock:        c.Bool("lock"),
			LockTimeout: c.Duration("lock_timeout"),
			NoColor:     c.Bool("no_color"),
			Parallelism: c.Int("parallelism"),
			Refresh:     c.Bool("refresh"),
			State:       c.String("state"),
			StateOut:    c.String("state_out"),
			Target:      c.String("target"),
			Vars:        c.StringSlice("vars"),
			VarFiles:    c.StringSlice("var_files"),
		},
		// FMT configuration
		FMT: &FMT{
			Check:     c.Bool("fmt.check"),
			Diff:      c.Bool("fmt.diff"),
			Directory: c.String("directory"),
			List:      c.Bool("fmt.list"),
			Write:     c.Bool("fmt.write"),
		},
		// InitOptions configuration
		Init: &Init{
			Directory: c.String("directory"),
			RawInit:   c.String("init.options"),
		},
		// Plan configuration
		Plan: &Plan{
			Destroy:          c.Bool("plan.destroy"),
			DetailedExitCode: c.Bool("plan.detailed_exit_code"),
			Directory:        c.String("directory"),
			Input:            c.Bool("plan.input"),
			Lock:             c.Bool("lock"),
			LockTimeout:      c.Duration("lock_timeout"),
			ModuleDepth:      c.Int("plan.module_depth"),
			NoColor:          c.Bool("no_color"),
			Parallelism:      c.Int("parallelism"),
			Refresh:          c.Bool("refresh"),
			State:            c.String("state"),
			Target:           c.String("target"),
			Vars:             c.StringSlice("vars"),
			VarFiles:         c.StringSlice("var_files"),
		},
		// Validation configuration
		Validation: &Validation{
			CheckVariables: c.Bool("validation.check_variables"),
			Directory:      c.String("directory"),
			NoColor:        c.Bool("no_color"),
			Vars:           c.StringSlice("vars"),
			VarFiles:       c.StringSlice("var_files"),
		},
	}

	// validate the plugin
	err := p.Validate()
	if err != nil {
		return err
	}

	return p.Exec()
}
