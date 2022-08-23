## Description

This plugin enables you to run [Terraform](https://www.terraform.io/) against [providers](https://www.terraform.io/docs/providers/index.html) in a Vela pipeline.

Source Code: https://github.com/go-vela/vela-terraform

Registry: https://hub.docker.com/r/target/vela-terraform

## Usage

> **NOTE:**
>
> Users should refrain from using latest as the tag for the Docker image.
>
> It is recommended to use a semantically versioned tag instead.

Sample of adding installing terraform version:

```yaml
steps:
  - name: apply
    image: target/vela-terraform:latest
    pull: always
    parameters:
      action: apply
      auto_approve: true # Required for versions of Terraform 0.12.x
      version: 0.11.7
```

Sample of adding init options to Terraform configuration:

```yaml
steps:
  - name: apply
    image: target/vela-terraform:latest
    pull: always
    parameters:
      action: apply
      auto_approve: true # Required for versions of Terraform 0.12.x
      init_options:
        get_plugins: true
```

Sample of applying Terraform configuration:

```yaml
steps:
  - name: apply
    image: target/vela-terraform:latest
    pull: always
    parameters:
      action: apply
      auto_approve: true # Required for versions of Terraform 0.12.x
```

Sample of destroying Terraform configuration:

```yaml
steps:
  - name: destroy
    image: target/vela-terraform:latest
    pull: always
    parameters:
      action: destroy
      auto_approve: true # Required for versions of Terraform 0.12.x
```

Sample of formatting Terraform configuration files:

```yaml
steps:
  - name: fmt
    image: target/vela-terraform:latest
    pull: always
    parameters:
      action: fmt
```

Sample of planning Terraform configuration:

```yaml
steps:
  - name: plan
    image: target/vela-terraform:latest
    pull: always
    parameters:
      action: plan
```

Sample of validating Terraform configuration:

```yaml
steps:
  - name: validate
    image: target/vela-terraform:latest
    pull: always
    parameters:
      action: validate
```

## Secrets

> **NOTE:** Users should refrain from configuring sensitive information in your pipeline in plain text.

### Internal

Users can use [Vela internal secrets](https://go-vela.github.io/docs/tour/secrets/) to substitute these sensitive values at runtime:

```diff
steps:
  - name: apply
    image: target/vela-terraform:latest
    pull: always
+   secrets: [ terraform_username, terraform_password ]
    parameters:
      action: apply
      auto_approve: true # Required for versions of Terraform 0.12.x
-     username: octocat
-     password: superSecretPassword
```

> This example will add the secrets to the `apply` step as environment variables:
>
> * `TERRAFORM_USERNAME=<value>`
> * `TERRAFORM_PASSWORD=<value>`

### External

The plugin accepts the following files for authentication:

| Parameter  | Volume Configuration                                                      |
| ---------- | ------------------------------------------------------------------------- |
| `password` | `/vela/parameters/terraform/password`, `/vela/secrets/terraform/password` |
| `username` | `/vela/parameters/terraform/username`, `/vela/secrets/terraform/username` |

Users can use [Vela external secrets](https://go-vela.github.io/docs/tour/secrets/) to substitute these sensitive values at runtime:

```diff
steps:
  - name: apply
    image: target/vela-terraform:latest
    pull: always
    parameters:
      action: apply
      auto_approve: true # Required for versions of Terraform 0.12.x
-     username: octocat
-     password: superSecretPassword
```

> This example will read the secret values in the volume stored at `/vela/secrets/`

## Parameters

> **NOTE:**
>
> The plugin supports reading all parameters via environment variables or files.
>
> Any values set from a file take precedence over values set from the environment.
>
> Terraform commands will be invoked in the current directory by default.

The following parameters are used to configure the image:

| Name           | Description                                 | Required | Default         | Environment Variables                                                 |
| -------------- | ------------------------------------------- | -------- | --------------- | --------------------------------------------------------------------- |
| `action`       | action to perform with Terraform            | `true`   | `N/A`           | `PARAMETER_ACTION`<br>`TERRAFORM_ACTION`                              |
| `init_options` | options to use for Terraform init operation | `false`  | `N/A`           | `PARAMETER_INIT_OPTIONS`<br>`TERRAFORM_INIT_OPTIONS`                  |
| `log_level`    | set the log level for the plugin            | `true`   | `info`          | `PARAMETER_LOG_LEVEL`<br>`TERRAFORM_LOG_LEVEL`                        |
| `machine`      | netrc machine name to communicate with      | `true`   | `github.com`    | `PARAMETER_MACHINE`<br>`TERRAFORM_MACHINE`<br>`VELA_NETRC_MACHINE`    |
| `password`     | netrc password for authentication           | `true`   | **set by Vela** | `PARAMETER_PASSWORD`<br>`TERRAFORM_PASSWORD`<br>`VELA_NETRC_PASSWORD` |
| `username`     | netrc user name for authentication          | `true`   | **set by Vela** | `PARAMETER_USERNAME`<br>`TERRAFORM_USERNAME`<br>`VELA_NETRC_USERNAME` |
| `version`      | set the Terraform CLI version               | `true`   | `1.2.7`         | `PARAMETER_VERSION`<br>`TERRAFORM_VERSION`                            |

The following parameters can be used within the `init_options` to configure the image:

| Name              | Description                                                                           | Required | Default |
| ----------------- | ------------------------------------------------------------------------------------- | -------- | ------- |
| `backend`         | configure the backend for this configuration                                          | `true`   | `N/A`   |
| `backend_configs` | this is merged with what is in the configuration file                                 | `true`   | `N/A`   |
| `force_copy`      | suppress prompts about copying state data                                             | `true`   | `N/A`   |
| `from_module`     | copy the contents of the given module into the target directory before initialization | `true`   | `N/A`   |
| `get`             | download any modules for this configuration                                           | `true`   | `N/A`   |
| `get_plugins`     | download any missing plugins for this configuration                                   | `true`   | `N/A`   |
| `input`           | ask for input for variables if not directly set                                       | `true`   | `N/A`   |
| `lock`            | lock the state file when locking is supported                                         | `false`  | `N/A`   |
| `lock_timeout`    | duration to retry a state lock                                                        | `false`  | `N/A`   |
| `no_color`        | disables colors in output                                                             | `false`  | `N/A`   |
| `plugin_dirs`     | directory containing plugin binaries; overrides all default search paths for plugins  | `false`  | `N/A`   |
| `reconfigure`     | reconfigure the backend, ignoring any saved configuration                             | `false`  | `N/A`   |
| `upgrade`         | install the latest version allowed within configured constraints                      | `false`  | `N/A`   |
| `verify_plugins`  | verify the authenticity and integrity of automatically downloaded plugins             | `false`  | `N/A`   |

#### Apply

The following parameters are used to configure the `apply` action:

_Command uses Terraform CLI command defaults if not overridden in config._

| Name           | Description                                                   | Required | Default | Environment Variables                                |
| -------------- | ------------------------------------------------------------- | -------- | ------- | ---------------------------------------------------- |
| `auto_approve` | skip interactive approval of applying resources               | `false`  | `false` | `PARAMETER_AUTO_APPROVE`<br>`TERRAFORM_AUTO_APPROVE` |
| `backup`       | path to backup the existing state file                        | `false`  | `N/A`   | `PARAMETER_BACKUP`<br>`TERRAFORM_BACKUP`             |
| `directory`    | the directory containing Terraform files to apply             | `false`  | `.`     | `PARAMETER_DIRECTORY`<br>`TERRAFORM_DIRECTORY`       |
| `lock`         | lock the state file when locking is supported                 | `false`  | `false` | `PARAMETER_LOCK`<br>`TERRAFORM_LOCK`                 |
| `lock_timeout` | duration to retry a state lock                                | `false`  | `N/A`   | `PARAMETER_LOCK_TIMEOUT`<br>`TERRAFORM_LOCK_TIMEOUT` |
| `no_color`     | disables colors in output                                     | `false`  | `false` | `PARAMETER_NO_COLOR`<br>`TERRAFORM_NO_COLOR`         |
| `parallelism`  | number of concurrent operations as Terraform walks its graph  | `false`  | `N/A`   | `PARAMETER_PARALLELISM`<br>`TERRAFORM_PARALLELISM`   |
| `refresh`      | update state prior to checking for differences                | `false`  | `false` | `PARAMETER_REFRESH`<br>`TERRAFORM_REFRESH`           |
| `state`        | path to read and save state                                   | `false`  | `N/A`   | `PARAMETER_STATE`<br>`TERRAFORM_STATE`               |
| `state_out`    | path to write updated state file                              | `false`  | `N/A`   | `PARAMETER_STATE_OUT`<br>`TERRAFORM_STATE_OUT`       |
| `target`       | resource to target                                            | `false`  | `N/A`   | `PARAMETER_TARGET`<br>`TERRAFORM_TARGET`             |
| `vars`         | a map of variables to pass to the Terraform (`<key>=<value>`) | `false`  | `N/A`   | `PARAMETER_VARS`<br>`TERRAFORM_VARS`                 |
| `var_files`    | a list of var files to use                                    | `false`  | `N/A`   | `PARAMETER_VAR_FILES`<br>`TERRAFORM_VAR_FILES`       |

#### Destroy

The following parameters are used to configure the `destroy` action:

_Command uses Terraform CLI command defaults if not overridden in config._

| Name           | Description                                                   | Required | Default | Environment Variables                                |
| -------------- | ------------------------------------------------------------- | -------- | ------- | ---------------------------------------------------- |
| `auto_approve` | skip interactive approval of destroying resources             | `false`  | `false` | `PARAMETER_AUTO_APPROVE`<br>`TERRAFORM_AUTO_APPROVE` |
| `backup`       | path to backup the existing state file                        | `false`  | `N/A`   | `PARAMETER_BACKUP`<br>`TERRAFORM_BACKUP`             |
| `directory`    | the directory containing Terraform files to destroy           | `false`  | `.`     | `PARAMETER_DIRECTORY`<br>`TERRAFORM_DIRECTORY`       |
| `lock`         | lock the state file when locking is supported                 | `false`  | `false` | `PARAMETER_LOCK`<br>`TERRAFORM_LOCK`                 |
| `lock_timeout` | duration to retry a state lock                                | `false`  | `N/A`   | `PARAMETER_LOCK_TIMEOUT`<br>`TERRAFORM_LOCK_TIMEOUT` |
| `no_color`     | disables colors in output                                     | `false`  | `false` | `PARAMETER_NO_COLOR`<br>`TERRAFORM_NO_COLOR`         |
| `parallelism`  | number of concurrent operations as Terraform walks its graph  | `false`  | `N/A`   | `PARAMETER_PARALLELISM`<br>`TERRAFORM_PARALLELISM`   |
| `refresh`      | update state prior to checking for differences                | `false`  | `false` | `PARAMETER_REFRESH`<br>`TERRAFORM_REFRESH`           |
| `state`        | path to read and save state                                   | `false`  | `N/A`   | `PARAMETER_STATE`<br>`TERRAFORM_STATE`               |
| `state_out`    | path to write updated state file                              | `false`  | `N/A`   | `PARAMETER_STATE_OUT`<br>`TERRAFORM_STATE_OUT`       |
| `target`       | resource to target                                            | `false`  | `N/A`   | `PARAMETER_TARGET`<br>`TERRAFORM_TARGET`             |
| `vars`         | a map of variables to pass to the Terraform (`<key>=<value>`) | `false`  | `N/A`   | `PARAMETER_VARS`<br>`TERRAFORM_VARS`                 |
| `var_files`    | a list of var files to use                                    | `false`  | `N/A`   | `PARAMETER_VAR_FILES`<br>`TERRAFORM_VAR_FILES`       |

#### Format

The following parameters are used to configure the `fmt` action:

_Command uses Terraform CLI command defaults if not overridden in config._

| Name        | Description                                        | Required | Default | Environment Variables                          |
| ----------- | -------------------------------------------------- | -------- | ------- | ---------------------------------------------- |
| `check`     | validate if the input is formatted                 | `false`  | `false` | `PARAMETER_CHECK`<br>`TERRAFORM_CHECK`         |
| `diff`      | diffs of formatting changes                        | `false`  | `false` | `PARAMETER_DIFF`<br>`TERRAFORM_DIFF`           |
| `directory` | the directory containing Terraform files to format | `false`  | `.`     | `PARAMETER_DIRECTORY`<br>`TERRAFORM_DIRECTORY` |
| `list`      | list files whose formatting differs                | `false`  | `false` | `PARAMETER_LIST`<br>`TERRAFORM_LIST`           |
| `write`     | write result to source file instead of STDOUT      | `false`  | `false` | `PARAMETER_WRITE`<br>`TERRAFORM_WRITE`         |

#### Plan

The following parameters are used to configure the `plan` action:

_Command uses Terraform CLI command defaults if not overridden in config._

| Name                 | Description                                                        | Required | Default | Environment Variables                                            |
| -------------------- | ------------------------------------------------------------------ | -------- | ------- | ---------------------------------------------------------------- |
| `destroy`            | destroy all resources managed by the given configuration and state | `false`  | `false` | `PARAMETER_DESTROY`<br>`TERRAFORM_DESTROY`                       |
| `detailed_exit_code` | return detailed exit codes when the command exits                  | `false`  | `false` | `PARAMETER_DETAILED_EXIT_CODE`<br>`TERRAFORM_DETAILED_EXIT_CODE` |
| `directory`          | the directory containing Terraform files to plan                   | `false`  | `.`     | `PARAMETER_DIRECTORY`<br>`TERRAFORM_DIRECTORY`                   |
| `input`              | ask for input for variables if not directly set                    | `false`  | `false` | `PARAMETER_INPUT`<br>`TERRAFORM_INPUT`                           |
| `lock`               | lock the state file when locking is supported                      | `false`  | `false` | `PARAMETER_LOCK`<br>`TERRAFORM_LOCK`                             |
| `lock_timeout`       | duration to retry a state lock                                     | `false`  | `N/A`   | `PARAMETER_LOCK_TIMEOUT`<br>`TERRAFORM_LOCK_TIMEOUT`             |
| `module_depth`       | specifies the depth of modules to show in the output               | `false`  | `N/A`   | `PARAMETER_MODULE_DEPTH`<br>`TERRAFORM_MODULE_DEPTH`             |
| `no_color`           | disables colors in output                                          | `false`  | `false` | `PARAMETER_NO_COLOR`<br>`TERRAFORM_NO_COLOR`                     |
| `parallelism`        | number of concurrent operations as Terraform walks its graph       | `false`  | `N/A`   | `PARAMETER_PARALLELISM`<br>`TERRAFORM_PARALLELISM`               |
| `refresh`            | update state prior to checking for differences                     | `false`  | `false` | `PARAMETER_REFRESH`<br>`TERRAFORM_REFRESH`                       |
| `state`              | path to read and save state                                        | `false`  | `N/A`   | `PARAMETER_STATE`<br>`TERRAFORM_STATE`                           |
| `target`             | resource to target                                                 | `false`  | `N/A`   | `PARAMETER_TARGET`<br>`TERRAFORM_TARGET`                         |
| `vars`               | a map of variables to pass to the Terraform (`<key>=<value>`)      | `false`  | `N/A`   | `PARAMETER_VARS`<br>`TERRAFORM_VARS`                             |
| `var_files`          | a list of var files to use                                         | `false`  | `N/A`   | `PARAMETER_VAR_FILES`<br>`TERRAFORM_VAR_FILES`                   |

#### Validate

The following parameters are used to configure the `validate` action:

_Command uses Terraform CLI command defaults if not overridden in config._

| Name              | Description                                                           | Required | Default | Environment Variables                                      |
| ----------------- | --------------------------------------------------------------------- | -------- | ------- | ---------------------------------------------------------- |
| `check_variables` | command will check whether all required variables have been specified | `false`  | `false` | `PARAMETER_CHECK_VARIABLES`<br>`TERRAFORM_CHECK_VARIABLES` |
| `directory`       | the directory containing Terraform files to validate                  | `false`  | `.`     | `PARAMETER_DIRECTORY`<br>`TERRAFORM_DIRECTORY`             |
| `no_color`        | disables colors in output                                             | `false`  | `false` | `PARAMETER_NO_COLOR`<br>`TERRAFORM_NO_COLOR`               |
| `vars`            | a map of variables to pass to the Terraform (`<key>=<value>`)         | `false`  | `N/A`   | `PARAMETER_VARS`<br>`TERRAFORM_VARS`                       |
| `var_files`       | a list of var files to use                                            | `false`  | `N/A`   | `PARAMETER_VAR_FILES`<br>`TERRAFORM_VAR_FILES`             |

## Template

COMING SOON!

## Troubleshooting

You can start troubleshooting this plugin by tuning the level of logs being displayed:

```diff
steps:
  - name: apply
    image: target/vela-terraform:latest
    pull: always
    parameters:
      action: apply
      auto_approve: true
+     log_level: trace
```

You can also instruct the Terraform CLI to output verbose logging:

```diff
steps:
  - name: apply
    image: target/vela-terraform:latest
    pull: always
+   environment:
+     TF_LOG: TRACE
    parameters:
      action: apply
      auto_approve: true
```

Below are a list of common problems and how to solve them:
