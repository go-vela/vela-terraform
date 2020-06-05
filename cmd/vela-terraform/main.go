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
		&cli.StringFlag{
			EnvVars: []string{"PARAMETER_LOG_LEVEL", "VELA_LOG_LEVEL", "KUBERNETES_LOG_LEVEL"},
			Name:    "log.level",
			Usage:   "set log level - options: (trace|debug|info|warn|error|fatal|panic)",
			Value:   "info",
		},
		&cli.StringFlag{
			EnvVars: []string{"PARAMETER_VERSION", "VELA_TERRAFORM_VERSION", "TERRAFORM_VERSION"},
			Name:    "version",
			Usage:   "set terraform version for plugin",
		},

		// Apply flags
		&cli.BoolFlag{
			Name:    "apply.auto-approve",
			Usage:   "skip interactive approval of running command",
			EnvVars: []string{"PARAMETER_AUTO_APPROVE"},
		},
		&cli.StringFlag{
			Name:    "apply.back-up",
			Usage:   "path to backup the existing state file",
			EnvVars: []string{"PARAMETER_BACK_UP"},
		},
		&cli.StringFlag{
			Name:    "apply.directory",
			Usage:   "the directory for action to be performed on",
			EnvVars: []string{"PARAMETER_DIRECTORY"},
		},
		&cli.BoolFlag{
			Name:    "apply.lock",
			Usage:   "lock the state file when locking is supported",
			EnvVars: []string{"PARAMETER_LOCK"},
		},
		&cli.DurationFlag{
			Name:    "apply.lock-timeout",
			Usage:   "duration to retry a state lock",
			EnvVars: []string{"PARAMETER_LOCK_TIMEOUT"},
		},
		&cli.BoolFlag{
			Name:    "apply.no-color",
			Usage:   "disables colors in output",
			EnvVars: []string{"PARAMETER_NO_COLOR"},
		},
		&cli.IntFlag{
			Name:    "apply.parallelism",
			Usage:   "number of concurrent operations as Terraform walks its graph",
			EnvVars: []string{"PARAMETER_PARALLELISM"},
		},
		&cli.BoolFlag{
			Name:    "apply.refresh",
			Usage:   "update state prior to checking for differences",
			EnvVars: []string{"PARAMETER_REFRESH"},
		},
		&cli.StringFlag{
			Name:    "apply.state",
			Usage:   "path to read and save state",
			EnvVars: []string{"PARAMETER_STATE"},
		},
		&cli.StringFlag{
			Name:    "apply.state-out",
			Usage:   "path to write updated state file",
			EnvVars: []string{"PARAMETER_STATE_OUT"},
		},
		&cli.StringFlag{
			Name:    "apply.target",
			Usage:   "resource to target",
			EnvVars: []string{"PARAMETER_TARGET"},
		},
		&cli.StringSliceFlag{
			Name:    "apply.vars",
			Usage:   "a map of variables to pass to the Terraform `plan` and `apply` commands. Each value is passed as a `<key>=<value>` option",
			EnvVars: []string{"PARAMETER_VARS"},
		},
		&cli.StringSliceFlag{
			Name:    "apply.var_files",
			Usage:   "a list of var files to use. Each value is passed as -var-file=<value>",
			EnvVars: []string{"PARAMETER_VAR_FILES"},
		},

		// Config flags
		&cli.StringFlag{
			Name:    "config.action",
			Usage:   "the action to have terraform perform",
			EnvVars: []string{"PARAMETER_ACTION"},
		},
		&cli.StringFlag{
			EnvVars: []string{"PARAMETER_NETRC_MACHINE", "VELA_NETRC_MACHINE"},
			Name:    "config.netrc.machine",
			Usage:   "remote machine name to communicate with",
		},
		&cli.StringFlag{
			EnvVars: []string{"PARAMETER_NETRC_USERNAME", "VELA_NETRC_USERNAME", "GIT_USERNAME"},
			Name:    "config.netrc.username",
			Usage:   "user name for communication with the remote machine",
		},
		&cli.StringFlag{
			EnvVars: []string{"PARAMETER_NETRC_PASSWORD", "VELA_NETRC_PASSWORD", "GIT_PASSWORD"},
			Name:    "config.netrc.password",
			Usage:   "password for communication with the remote machine",
		},

		// Destroy flags
		&cli.BoolFlag{
			Name:    "destroy.auto-approve",
			Usage:   "skip interactive approval of running command",
			EnvVars: []string{"PARAMETER_AUTO_APPROVE"},
		},
		&cli.StringFlag{
			Name:    "destroy.back-up",
			Usage:   "path to backup the existing state file",
			EnvVars: []string{"PARAMETER_BACK_UP"},
		},
		&cli.StringFlag{
			Name:    "destroy.directory",
			Usage:   "the directory for action to be performed on",
			EnvVars: []string{"PARAMETER_DIRECTORY"},
		},
		&cli.BoolFlag{
			Name:    "destroy.lock",
			Usage:   "lock the state file when locking is supported",
			EnvVars: []string{"PARAMETER_LOCK"},
		},
		&cli.DurationFlag{
			Name:    "destroy.lock-timeout",
			Usage:   "duration to retry a state lock",
			EnvVars: []string{"PARAMETER_LOCK_TIMEOUT"},
		},
		&cli.BoolFlag{
			Name:    "destroy.no-color",
			Usage:   "disables colors in output",
			EnvVars: []string{"PARAMETER_NO_COLOR"},
		},
		&cli.IntFlag{
			Name:    "destroy.parallelism",
			Usage:   "number of concurrent operations as Terraform walks its graph",
			EnvVars: []string{"PARAMETER_PARALLELISM"},
		},
		&cli.BoolFlag{
			Name:    "destroy.refresh",
			Usage:   "update state prior to checking for differences",
			EnvVars: []string{"PARAMETER_REFRESH"},
		},
		&cli.StringFlag{
			Name:    "destroy.state",
			Usage:   "path to read and save state",
			EnvVars: []string{"PARAMETER_STATE"},
		},
		&cli.StringFlag{
			Name:    "destroy.state-out",
			Usage:   "path to write updated state file",
			EnvVars: []string{"PARAMETER_STATE_OUT"},
		},
		&cli.StringFlag{
			Name:    "destroy.target",
			Usage:   "resource to target",
			EnvVars: []string{"PARAMETER_TARGET"},
		},
		&cli.StringSliceFlag{
			Name:    "destroy.vars",
			Usage:   "a map of variables to pass to the Terraform each value is passed as a `<key>=<value>` option",
			EnvVars: []string{"PARAMETER_VARS"},
		},
		&cli.StringSliceFlag{
			Name:    "destroy.var_files",
			Usage:   "a list of var files to use. Each value is passed as -var-file=<value>",
			EnvVars: []string{"PARAMETER_VAR_FILES"},
		},

		// FMT flags
		&cli.BoolFlag{
			Name:    "fmt.check",
			Usage:   "validate if the input is formatted",
			EnvVars: []string{"PARAMETER_CHECK"},
		},
		&cli.BoolFlag{
			Name:    "fmt.diff",
			Usage:   "diffs of formatting changes",
			EnvVars: []string{"PARAMETER_DIFF"},
		},
		&cli.StringFlag{
			Name:    "fmt.directory",
			Usage:   "the directory for action to be performed on",
			EnvVars: []string{"PARAMETER_DIRECTORY"},
		},
		&cli.BoolFlag{
			Name:    "fmt.list",
			Usage:   "list files whose formatting differs",
			EnvVars: []string{"PARAMETER_LIST"},
		},
		&cli.BoolFlag{
			Name:    "fmt.write",
			Usage:   "write result to source file instead of STDOUT",
			EnvVars: []string{"PARAMETER_WRITE"},
		},

		// InitOptions flags
		&cli.StringFlag{
			EnvVars: []string{"PARAMETER_INIT_OPTIONS"},
			Name:    "init.options",
			Usage:   "properties to set on terraform init action",
		},

		// Plan flags
		&cli.BoolFlag{
			Name:    "plan.destroy",
			Usage:   "destroy all resources managed by the given configuration and state",
			EnvVars: []string{"PARAMETER_DESTROY"},
		},
		&cli.BoolFlag{
			Name:    "plan.detailed-exit-code",
			Usage:   "return detailed exit codes when the command exits",
			EnvVars: []string{"PARAMETER_DETAILED_EXIT_CODE"},
		},
		&cli.StringFlag{
			Name:    "plan.directory",
			Usage:   "the directory for action to be performed on",
			EnvVars: []string{"PARAMETER_DIRECTORY"},
		},
		&cli.BoolFlag{
			Name:    "plan.lock",
			Usage:   "lock the state file when locking is supported",
			EnvVars: []string{"PARAMETER_LOCK"},
		},
		&cli.DurationFlag{
			Name:    "plan.lock-timeout",
			Usage:   "duration to retry a state lock",
			EnvVars: []string{"PARAMETER_LOCK_TIMEOUT"},
		},
		&cli.IntFlag{
			Name:    "module-depth",
			Usage:   "specifies the depth of modules to show in the output",
			EnvVars: []string{"PARAMETER_MODULE_DEPTH"},
		},
		&cli.BoolFlag{
			Name:    "plan.no-color",
			Usage:   "disables colors in output",
			EnvVars: []string{"PARAMETER_NO_COLOR"},
		},
		&cli.IntFlag{
			Name:    "plan.parallelism",
			Usage:   "number of concurrent operations as Terraform walks its graph",
			EnvVars: []string{"PARAMETER_PARALLELISM"},
		},
		&cli.BoolFlag{
			Name:    "plan.refresh",
			Usage:   "update state prior to checking for differences",
			EnvVars: []string{"PARAMETER_REFRESH"},
		},
		&cli.StringFlag{
			Name:    "plan.state",
			Usage:   "path to read and save state",
			EnvVars: []string{"PARAMETER_STATE"},
		},
		&cli.StringFlag{
			Name:    "plan.state-out",
			Usage:   "path to write updated state file",
			EnvVars: []string{"PARAMETER_STATE_OUT"},
		},
		&cli.StringFlag{
			Name:    "plan.target",
			Usage:   "resource to target",
			EnvVars: []string{"PARAMETER_TARGET"},
		},
		&cli.StringSliceFlag{
			Name:    "plan.vars",
			Usage:   "a map of variables to pass to the Terraform each value is passed as a `<key>=<value>` option",
			EnvVars: []string{"PARAMETER_VARS"},
		},
		&cli.StringSliceFlag{
			Name:    "plan.var_files",
			Usage:   "a list of var files to use. Each value is passed as -var-file=<value>",
			EnvVars: []string{"PARAMETER_VAR_FILES"},
		},

		// Validation flags
		&cli.BoolFlag{
			Name:    "validation.check-variables",
			Usage:   "command will check whether all required variables have been specified",
			EnvVars: []string{"PARAMETER_CHECK_VARIABLES"},
		},
		&cli.StringFlag{
			Name:    "validation.directory",
			Usage:   "the directory for action to be performed on",
			EnvVars: []string{"PARAMETER_DIRECTORY"},
		},
		&cli.BoolFlag{
			Name:    "validation.no-color",
			Usage:   "disables colors in output",
			EnvVars: []string{"PARAMETER_NO_COLOR"},
		},
		&cli.StringSliceFlag{
			Name:    "validation.vars",
			Usage:   "a map of variables to pass to the Terraform `plan` and `apply` commands. Each value is passed as a `<key>=<value>` option",
			EnvVars: []string{"PARAMETER_VARS"},
		},
		&cli.StringSliceFlag{
			Name:    "validation.var_files",
			Usage:   "a list of var files to use. Each value is passed as -var-file=<value>",
			EnvVars: []string{"PARAMETER_VAR_FILES"},
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
	p := Plugin{
		// Apply configuration
		Apply: &Apply{
			AutoApprove: c.Bool("apply.auto-approve"),
			Backup:      c.String("apply.backup"),
			Directory:   c.String("apply.directory"),
			Lock:        c.Bool("apply.lock"),
			LockTimeout: c.Duration("apply.lock-timeout"),
			NoColor:     c.Bool("apply.no-color"),
			Parallelism: c.Int("apply.parallelism"),
			Refresh:     c.Bool("apply.refresh"),
			State:       c.String("apply.state"),
			StateOut:    c.String("apply.state-out"),
			Target:      c.String("apply.target"),
			Vars:        c.StringSlice("apply.vars"),
			VarFiles:    c.StringSlice("apply.var-files"),
		},
		// Config configuration
		Config: &Config{
			Action: c.String("config.action"),
			Netrc: &Netrc{
				Login:    c.String("config.netrc.username"),
				Machine:  c.String("config.netrc.machine"),
				Password: c.String("config.netrc.password"),
			},
		},
		// Destroy configuration
		Destroy: &Destroy{
			AutoApprove: c.Bool("destroy.auto-approve"),
			Backup:      c.String("destroy.backup"),
			Directory:   c.String("destroy.directory"),
			Lock:        c.Bool("destroy.lock"),
			LockTimeout: c.Duration("destroy.lock-timeout"),
			NoColor:     c.Bool("destroy.no-color"),
			Parallelism: c.Int("destroy.parallelism"),
			Refresh:     c.Bool("destroy.refresh"),
			State:       c.String("destroy.state"),
			StateOut:    c.String("destroy.state-out"),
			Target:      c.String("destroy.target"),
			Vars:        c.StringSlice("destroy.vars"),
			VarFiles:    c.StringSlice("destroy.var-files"),
		},
		// FMT configuration
		FMT: &FMT{
			Check:     c.Bool("fmt.check"),
			Diff:      c.Bool("fmt.diff"),
			Directory: c.String("fmt.directory"),
			List:      c.Bool("fmt.list"),
			Write:     c.Bool("fmt.write"),
		},
		// Config configuration
		InitOptions: &InitOptions{
			RawInit: c.String("init.options"),
		},
		// Plan configuration
		Plan: &Plan{
			Destroy:          c.Bool("plan.destroy"),
			DetailedExitCode: c.Bool("plan.detailed-exit-code"),
			Directory:        c.String("plan.directory"),
			Input:            c.Bool("plan.input"),
			Lock:             c.Bool("plan.lock"),
			LockTimeout:      c.Duration("plan.lock-timeout"),
			ModuleDepth:      c.Int("plan.module-depth"),
			NoColor:          c.Bool("plan.no-color"),
			Parallelism:      c.Int("plan.parallelism"),
			Refresh:          c.Bool("plan.refresh"),
			State:            c.String("plan.state"),
			Target:           c.String("plan.target"),
			Vars:             c.StringSlice("plan.vars"),
			VarFiles:         c.StringSlice("plan.var-files"),
		},
		// Validation configuration
		Validation: &Validation{
			CheckVariables: c.Bool("validation.check-variables"),
			Directory:      c.String("validation.directory"),
			NoColor:        c.Bool("validation.no-color"),
			Vars:           c.StringSlice("validation.vars"),
			VarFiles:       c.StringSlice("validation.var-files"),
		},
	}

	// validate the plugin
	err := p.Validate()
	if err != nil {
		return err
	}

	return p.Exec()
}
