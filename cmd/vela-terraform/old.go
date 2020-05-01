package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/sirupsen/logrus"
)

// TODO: This file and all it's function will be deleted in future implementations. Code needed to move to make use of new design

// CopyTfEnv creates copies of TF_VAR_ to lowercase
func CopyTfEnv() {
	tfVar := regexp.MustCompile(`^TF_VAR_.*$`)
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if tfVar.MatchString(pair[0]) {
			name := strings.Split(pair[0], "TF_VAR_")
			os.Setenv(fmt.Sprintf("TF_VAR_%s", strings.ToLower(name[1])), pair[1])
		}
	}
}

func assumeRole(roleArn string) {
	client := sts.New(session.New())
	duration := time.Hour * 1
	stsProvider := &stscreds.AssumeRoleProvider{
		Client:          client,
		Duration:        duration,
		RoleARN:         roleArn,
		RoleSessionName: "vela",
	}

	value, err := credentials.NewCredentials(stsProvider).Get()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Error assuming role!")
	}
	os.Setenv("AWS_ACCESS_KEY_ID", value.AccessKeyID)
	os.Setenv("AWS_SECRET_ACCESS_KEY", value.SecretAccessKey)
	os.Setenv("AWS_SESSION_TOKEN", value.SessionToken)
}

func deleteCache() *exec.Cmd {
	return exec.Command(
		"rm",
		"-rf",
		".terraform",
	)
}

func getModules() *exec.Cmd {
	return exec.Command(
		"terraform",
		"get",
	)
}

func initCommand(config InitOptions) *exec.Cmd {
	args := []string{
		"init",
	}

	for _, v := range config.BackendConfig {
		args = append(args, fmt.Sprintf("-backend-config=%s", v))
	}

	// True is default in TF
	if config.Lock != nil {
		args = append(args, fmt.Sprintf("-lock=%t", *config.Lock))
	}

	// "0s" is default in TF
	if config.LockTimeout != "" {
		args = append(args, fmt.Sprintf("-lock-timeout=%s", config.LockTimeout))
	}

	// Fail Terraform execution on prompt
	args = append(args, "-input=false")

	return exec.Command(
		"terraform",
		args...,
	)
}

func installCaCert(cacert string) *exec.Cmd {
	ioutil.WriteFile("/usr/local/share/ca-certificates/ca_cert.crt", []byte(cacert), 0644)
	return exec.Command(
		"update-ca-certificates",
	)
}

func trace(cmd *exec.Cmd) {
	fmt.Println("$", strings.Join(cmd.Args, " "))
}

func tfApply(config Config) *exec.Cmd {
	args := []string{
		"apply",
	}
	for _, v := range config.Targets {
		args = append(args, "--target", fmt.Sprintf("%s", v))
	}
	if config.Parallelism > 0 {
		args = append(args, fmt.Sprintf("-parallelism=%d", config.Parallelism))
	}
	if config.InitOptions.Lock != nil {
		args = append(args, fmt.Sprintf("-lock=%t", *config.InitOptions.Lock))
	}
	if config.InitOptions.LockTimeout != "" {
		args = append(args, fmt.Sprintf("-lock-timeout=%s", config.InitOptions.LockTimeout))
	}
	args = append(args, "plan.tfout")
	return exec.Command(
		"terraform",
		args...,
	)
}

func tfDestroy(config Config) *exec.Cmd {
	args := []string{
		"destroy",
	}
	for _, v := range config.Targets {
		args = append(args, fmt.Sprintf("-target=%s", v))
	}
	args = append(args, varFiles(config.VarFiles)...)
	args = append(args, vars(config.Vars)...)
	if config.Parallelism > 0 {
		args = append(args, fmt.Sprintf("-parallelism=%d", config.Parallelism))
	}
	if config.InitOptions.Lock != nil {
		args = append(args, fmt.Sprintf("-lock=%t", *config.InitOptions.Lock))
	}
	if config.InitOptions.LockTimeout != "" {
		args = append(args, fmt.Sprintf("-lock-timeout=%s", config.InitOptions.LockTimeout))
	}
	args = append(args, "-force")
	return exec.Command(
		"terraform",
		args...,
	)
}

func tfPlan(config Config, destroy bool) *exec.Cmd {
	args := []string{
		"plan",
	}

	if destroy {
		args = append(args, "-destroy")
	} else {
		args = append(args, "-out=plan.tfout")
	}

	for _, v := range config.Targets {
		args = append(args, "--target", fmt.Sprintf("%s", v))
	}
	args = append(args, varFiles(config.VarFiles)...)
	args = append(args, vars(config.Vars)...)
	if config.Parallelism > 0 {
		args = append(args, fmt.Sprintf("-parallelism=%d", config.Parallelism))
	}
	if config.InitOptions.Lock != nil {
		args = append(args, fmt.Sprintf("-lock=%t", *config.InitOptions.Lock))
	}
	if config.InitOptions.LockTimeout != "" {
		args = append(args, fmt.Sprintf("-lock-timeout=%s", config.InitOptions.LockTimeout))
	}
	return exec.Command(
		"terraform",
		args...,
	)
}

func tfValidate(config Config) *exec.Cmd {
	args := []string{
		"validate",
	}
	for _, v := range config.VarFiles {
		args = append(args, fmt.Sprintf("-var-file=%s", v))
	}
	for k, v := range config.Vars {
		args = append(args, "-var", fmt.Sprintf("%s=%s", k, v))
	}
	return exec.Command(
		"terraform",
		args...,
	)
}

func tfFmt(config Config) *exec.Cmd {
	args := []string{
		"fmt",
	}
	if config.FmtOptions.List != nil {
		args = append(args, fmt.Sprintf("-list=%t", *config.FmtOptions.List))
	}
	if config.FmtOptions.Write != nil {
		args = append(args, fmt.Sprintf("-write=%t", *config.FmtOptions.Write))
	}
	if config.FmtOptions.Diff != nil {
		args = append(args, fmt.Sprintf("-diff=%t", *config.FmtOptions.Diff))
	}
	if config.FmtOptions.Check != nil {
		args = append(args, fmt.Sprintf("-check=%t", *config.FmtOptions.Check))
	}
	return exec.Command(
		"terraform",
		args...,
	)
}

func vars(vs map[string]string) []string {
	var args []string
	for k, v := range vs {
		args = append(args, "-var", fmt.Sprintf("%s=%s", k, v))
	}
	return args
}

func varFiles(vfs []string) []string {
	var args []string
	for _, v := range vfs {
		args = append(args, fmt.Sprintf("-var-file=%s", v))
	}
	return args
}

func writeNetrc(machine, login, password string) error {
	if machine == "" {
		return nil
	}
	out := fmt.Sprintf(
		netrcFile,
		machine,
		login,
		password,
	)

	home := "/root"
	u, err := user.Current()
	if err == nil {
		home = u.HomeDir
	}
	path := filepath.Join(home, ".netrc")
	return ioutil.WriteFile(path, []byte(out), 0600)
}

const netrcFile = `
machine %s
login %s
password %s
`
