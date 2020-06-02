// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	app := cli.NewApp()

	// Plugin Information

	app.Name = "vela-terraform"
	app.HelpName = "vela-terraform"
	app.Usage = "Vela Terraform plugin for running Terraform"
	app.Copyright = "Copyright (c) 2020 Target Brands, Inc. All rights reserved."
	app.Authors = []*cli.Author{
		{
			Name:  "Vela Admins",
			Email: "vela@target.com",
		},
	}

	// Plugin Metadata

	app.Compiled = time.Now()
	app.Action = run

	// Plugin Flags

	app.Flags = []cli.Flag{
		&cli.StringSliceFlag{
			Name:    "action",
			Usage:   "the action to have terraform perform",
			EnvVars: []string{"PARAMETER_ACTION"},
		},
		&cli.BoolFlag{
			Name:    "auto-approve",
			Usage:   "an interactive approval of of running command",
			EnvVars: []string{"PARAMETER_AUTO_APPROVE"},
		},
		&cli.StringFlag{
			Name:    "back-up",
			Usage:   "path to backup the existing state file",
			EnvVars: []string{"PARAMETER_BACK_UP"},
		},
		&cli.BoolFlag{
			Name:    "check",
			Usage:   "validate if the input is formatted",
			EnvVars: []string{"PARAMETER_CHECK"},
		},
		&cli.BoolFlag{
			Name:    "check-variables",
			Usage:   "command will check whether all required variables have been specified",
			EnvVars: []string{"PARAMETER_CHECK_VARIABLES"},
		},
		&cli.BoolFlag{
			Name:    "destroy",
			Usage:   "will be generated to destroy all resources managed by the given configuration and state",
			EnvVars: []string{"PARAMETER_DESTROY"},
		},
		&cli.BoolFlag{
			Name:    "detailed-exit-code",
			Usage:   "return detailed exit codes when the command exits",
			EnvVars: []string{"PARAMETER_DETAILED_EXIT_CODE"},
		},
		&cli.BoolFlag{
			Name:    "diff",
			Usage:   "diffs of formatting changes",
			EnvVars: []string{"PARAMETER_DIFF"},
		},
		&cli.StringFlag{
			Name:    "directory",
			Usage:   "the directory for action to be performed on",
			EnvVars: []string{"PARAMETER_DIRECTORY"},
		},
		&cli.BoolFlag{
			Name:    "input",
			Usage:   "ask for input for variables if not directly set",
			EnvVars: []string{"PARAMETER_INPUT"},
		},

		&cli.BoolFlag{
			Name:    "lock",
			Usage:   "the state file when locking is supported",
			EnvVars: []string{"PARAMETER_LOCK"},
		},
		&cli.DurationFlag{
			Name:    "lock-timeout",
			Usage:   "duration to retry a state lock",
			EnvVars: []string{"PARAMETER_LOCK_TIMEOUT"},
		},
		&cli.IntFlag{
			Name:    "module-depth",
			Usage:   "specifies the depth of modules to show in the output",
			EnvVars: []string{"PARAMETER_MODULE_DEPTH"},
		},
		&cli.BoolFlag{
			Name:    "no-color",
			Usage:   "the state file when locking is supported",
			EnvVars: []string{"PARAMETER_NO_COLOR"},
		},
		&cli.StringFlag{
			Name:    "out",
			Usage:   "write a plan file to the given path",
			EnvVars: []string{"PARAMETER_OUT"},
		},
		&cli.IntFlag{
			Name:    "parallelism",
			Usage:   "number of concurrent operations as Terraform walks its graph",
			EnvVars: []string{"PARAMETER_PARALLELISM"},
		},
		&cli.BoolFlag{
			Name:    "refresh",
			Usage:   "update state prior to checking for differences/cal",
			EnvVars: []string{"PARAMETER_REFRESH"},
		},
		&cli.StringFlag{
			Name:    "state",
			Usage:   "path to read and save state",
			EnvVars: []string{"PARAMETER_STATE"},
		},
		&cli.StringFlag{
			Name:    "state-out",
			Usage:   "path to write state to that is different than state",
			EnvVars: []string{"PARAMETER_STATE_OUT"},
		},
		&cli.StringFlag{
			Name:    "target",
			Usage:   "resource to target",
			EnvVars: []string{"PARAMETER_TARGET"},
		},
		&cli.StringSliceFlag{
			Name:    "vars",
			Usage:   "a map of variables to pass to the Terraform `plan` and `apply` commands. Each value is passed as a `<key>=<value>` option",
			EnvVars: []string{"PARAMETER_VARS"},
		},
		&cli.StringSliceFlag{
			Name:    "var_files",
			Usage:   "a list of var files to use. Each value is passed as -var-file=<value>",
			EnvVars: []string{"PARAMETER_VAR_FILES"},
		},
		&cli.StringFlag{
			EnvVars: []string{"PARAMETER_VERSION", "VELA_TERRAFORM_VERSION", "TERRAFORM_VERSION"},
			Name:    "version",
			Usage:   "set terraform version for plugin",
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}

// run executes the plugin based off the configuration provided.
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
	version := c.String("version")

	// check if a custom terraform version was requested
	if len(version) > 0 {
		// attempt to install the custom terraform version
		err := install(version, os.Getenv("PLUGIN_TERRAFORM_VERSION"))
		if err != nil {
			return err
		}
	}

	// create the plugin
	plugin := Plugin{
		// Apply configuration
		Apply: &Apply{
			AutoApprove: c.Bool("auto-approve"),
			Backup:      c.String("backup"),
			Directory:   c.String("directory"),
			Lock:        c.Bool("lock"),
			LockTimeout: c.Duration("lock-timeout"),
			NoColor:     c.Bool("no-color"),
			Parallelism: c.Int("parallelism"),
			Refresh:     c.Bool("refresh"),
			State:       c.String("state"),
			StateOut:    c.String("state-out"),
			Target:      c.String("target"),
			Vars:        c.StringSlice("vars"),
			VarFiles:    c.StringSlice("var-files"),
		},
		// Config configuration
		Config: Config{
			Action: c.String("action"),
			Netrc: &Netrc{
				Login:    c.String("netrc.username"),
				Machine:  c.String("netrc.machine"),
				Password: c.String("netrc.password"),
			},
		},
		// Destroy configuration
		Destroy: &Destroy{
			AutoApprove: c.Bool("auto-approve"),
			Backup:      c.String("backup"),
			Directory:   c.String("directory"),
			Lock:        c.Bool("lock"),
			LockTimeout: c.Duration("lock-timeout"),
			NoColor:     c.Bool("no-color"),
			Parallelism: c.Int("parallelism"),
			Refresh:     c.Bool("refresh"),
			State:       c.String("state"),
			StateOut:    c.String("state-out"),
			Target:      c.String("target"),
			Vars:        c.StringSlice("vars"),
			VarFiles:    c.StringSlice("var-files"),
		},
		// FMT configuration
		FMT: &FMT{
			Check:     c.Bool("check"),
			Diff:      c.Bool("diff"),
			Directory: c.String("directory"),
			List:      c.Bool("list"),
			Write:     c.Bool("write"),
		},
		// Plan configuration
		Plan: &Plan{
			Destroy:          c.Bool("destroy"),
			DetailedExitCode: c.Bool("detailed-exit-code"),
			Directory:        c.String("directory"),
			Input:            c.Bool("input"),
			Lock:             c.Bool("lock"),
			LockTimeout:      c.Duration("lock-timeout"),
			ModuleDepth:      c.Int("module-depth"),
			NoColor:          c.Bool("no-color"),
			Parallelism:      c.Int("parallelism"),
			Refresh:          c.Bool("refresh"),
			State:            c.String("state"),
			Target:           c.String("target"),
			Vars:             c.StringSlice("vars"),
			VarFiles:         c.StringSlice("var-files"),
		},
		// Validation configuration
		Validation: &Validation{
			CheckVariables: c.Bool("check-variables"),
			Directory:      c.String("directory"),
			NoColor:        c.Bool("no-color"),
			Vars:           c.StringSlice("vars"),
			VarFiles:       c.StringSlice("var-files"),
		},
	}

	return plugin.Exec()
}
