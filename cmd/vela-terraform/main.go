// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"encoding/json"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
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
			Name:    "actions",
			Usage:   "a list of actions to have terraform perform",
			EnvVars: []string{"PARAMETER_ACTIONS"},
		},
		&cli.StringFlag{
			Name:    "ca_cert",
			Usage:   "ca cert to add to your environment to allow terraform to use internal/private resources",
			EnvVars: []string{"PARAMETER_CA_CERT"},
		},
		&cli.StringFlag{
			Name:  "env-file",
			Usage: "source env file",
		},
		&cli.StringFlag{
			Name:    "init_options",
			Usage:   "options for the init command. See https://www.terraform.io/docs/commands/init.html",
			EnvVars: []string{"PARAMETER_INIT_OPTIONS"},
		},
		&cli.StringFlag{
			Name:    "fmt_options",
			Usage:   "options for the fmt command. See https://www.terraform.io/docs/commands/fmt.html",
			EnvVars: []string{"PARAMETER_FMT_OPTIONS"},
		},
		&cli.IntFlag{
			Name:    "parallelism",
			Usage:   "The number of concurrent operations as Terraform walks its graph",
			EnvVars: []string{"PARAMETER_PARALLELISM"},
		},
		&cli.StringFlag{
			Name:    "netrc.machine",
			Usage:   "netrc machine",
			EnvVars: []string{"VELA_NETRC_MACHINE"},
		},
		&cli.StringFlag{
			Name:    "netrc.username",
			Usage:   "netrc username",
			EnvVars: []string{"VELA_NETRC_USERNAME"},
		},
		&cli.StringFlag{
			Name:    "netrc.password",
			Usage:   "netrc password",
			EnvVars: []string{"VELA_NETRC_PASSWORD"},
		},
		&cli.StringFlag{
			Name:    "role_arn_to_assume",
			Usage:   "A role to assume before running the terraform commands",
			EnvVars: []string{"PARAMETER_ROLE_ARN_TO_ASSUME"},
		},
		&cli.StringFlag{
			Name:    "root_dir",
			Usage:   "The root directory where the terraform files live. When unset, the top level directory will be assumed",
			EnvVars: []string{"PARAMETER_ROOT_DIR"},
		},
		&cli.StringFlag{
			Name:    "secrets",
			Usage:   "a map of secrets to pass to the Terraform `plan` and `apply` commands. Each value is passed as a `<key>=<ENV>` option",
			EnvVars: []string{"PARAMETER_SECRETS"},
		},
		&cli.BoolFlag{
			Name:    "sensitive",
			Usage:   "whether or not to suppress terraform commands to stdout",
			EnvVars: []string{"PARAMETER_SENSITIVE"},
		},
		&cli.StringSliceFlag{
			Name:    "targets",
			Usage:   "targets to run apply or plan on",
			EnvVars: []string{"PARAMETER_TARGETS"},
		},
		&cli.StringFlag{
			Name:    "vars",
			Usage:   "a map of variables to pass to the Terraform `plan` and `apply` commands. Each value is passed as a `<key>=<value>` option",
			EnvVars: []string{"PARAMETER_VARS"},
		},
		&cli.StringSliceFlag{
			Name:    "var_files",
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

	if c.String("env-file") != "" {
		_ = godotenv.Load(c.String("env-file"))
	}

	var vars map[string]string
	if c.String("vars") != "" {
		if err := json.Unmarshal([]byte(c.String("vars")), &vars); err != nil {
			panic(err)
		}
	}
	var secrets map[string]string
	if c.String("secrets") != "" {
		if err := json.Unmarshal([]byte(c.String("secrets")), &secrets); err != nil {
			panic(err)
		}
	}

	initOptions := InitOptions{}
	json.Unmarshal([]byte(c.String("init_options")), &initOptions)
	fmtOptions := FmtOptions{}
	json.Unmarshal([]byte(c.String("fmt_options")), &fmtOptions)

	plugin := Plugin{
		Config: Config{
			Actions:     c.StringSlice("actions"),
			Vars:        vars,
			Secrets:     secrets,
			InitOptions: initOptions,
			FmtOptions:  fmtOptions,
			Cacert:      c.String("ca_cert"),
			Sensitive:   c.Bool("sensitive"),
			RoleARN:     c.String("role_arn_to_assume"),
			RootDir:     c.String("root_dir"),
			Parallelism: c.Int("parallelism"),
			Targets:     c.StringSlice("targets"),
			VarFiles:    c.StringSlice("var_files"),
		},
		Netrc: Netrc{
			Login:    c.String("netrc.username"),
			Machine:  c.String("netrc.machine"),
			Password: c.String("netrc.password"),
		},
	}

	return plugin.Exec()
}
