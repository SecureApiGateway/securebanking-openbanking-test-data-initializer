## securebanking-test-data-initializer
A service that creates a Payment service user with some data populated in RS for TPP consent.

## Requirements

- [go 1.16](https://golang.org/doc/install)
- configure [gopath](https://golang.org/doc/gopath_code.html#GOPATH)
- [pact](https://github.com/pact-foundation/pact-go#installation-on-nix)

## Program configuration variables (environment program)
The Test data initializer has been built to create a Payment service user and its corresponding data 

The initializer application provides a default configuration yaml file (properties values to run de application), the default configuration yaml file is loaded using the [viper library](https://github.com/spf13/viper),
the initializer application supports a personalized configuration file (as a profile) that it can be personalized for each required environment following the below rules:
- Path of environment file: `config/viper`
- Pattern environment file name: `viper-${environment-profile.viper_config}-configuration.yaml`
- Format configuration file (extension file): `yaml`

**Example:** `viper-default-configuration.yaml`
> You will find the provides default configuration yaml file in `config/viper` as an example.

> :warning: The initializer application only supports one configuration yaml file by application instance. It's recommended copy the provided default configuration yaml file and change its values.

> :memo: All the variables/properties values provides by the configuration file can be overwritten using environment variables or a kubernetes config map.
> ```shell
> go build -o setup \
> env ENVIRONMENT.VERBOSE=true ./setup
> ```

#### ConfigMap for variables mount example

```
apiVersion: v1
kind: ConfigMap
metadata:
  name: initializer-config
data:
  HOSTS.IG_FQDN: obdemo.dev.forgerock.financial
  HOSTS.IDENTITY_PLATFORM_FQDN: iam.dev.forgerock.financial
  ...
          
```

**Check the variables/properties in [Configuration variables](#configuration-variables) section.**

### The application configuration file
The configuration file is loaded from the path `config/viper` following the pattern `viper- + ${environment-name} + -configuration`
where the environment/profile can be specified by environment variable, passing that environment variable to the program.
**Examples**
```shell
go build -o setup
```
```shell
env ENVIRONMENT.VIPER_CONFIG=MY-ENVIRONMENT-PROFILE-VIPER_CONFIG ./setup
```
> The application will attempt to load the configuration file `viper-MY-ENVIRONMENT-PROFILE-VIPER_CONFIG-configuration.yaml`

**Other example**
```shell
env ENVIRONMENT.VIPER_CONFIG=MY-ENVIRONMENT-PROFILE-VIPER_CONFIG ENVIRONMENT.VERBOSE=true ... ./setup
```

#### Configuration variables
**Environment variables**
There are a variables used before load the configuration file and these variables can change the behaviour of the application.
- Behaviour variables:
  - `ENVIRONMENT.VERBOSE`
  - `ENVIRONMENT.VIPER_CONFIG`
  - `ENVIRONMENT.STRICT`
  - `ENVIRONMENT.ONLY_CONFIG`

<details>
<summary>Variables table</summary>
<!-- always an empty line before table -->

| variable                                     | Default value                        | Description                                                                                                                                                                                                       |
|----------------------------------------------|--------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `ENVIRONMENT.VERBOSE`                        | false                                | Log level (verbose=true means debug mode)                                                                                                                                                                         |
| `ENVIRONMENT.VIPER_CONFIG`                   | default                              | The profile that contains the configuration to be overwritten from system env                                                                                                                                     |
| `ENVIRONMENT.ONLY_CONFIG`                    | false                                | Prints the configuration and exiting the program, to review the properties                                                                                                                                        |
| `ENVIRONMENT.STRICT`                         | false                                | true = strict mode on, otherwise off, will exit if go resty returns an error in STRICT mode enabled, be it client error, server error or other. Turning off STRICT mode will simply warn of client/server errors. |
| `ENVIRONMENT.TYPE`                           | CDK                                  | values: CDK, CDM or FIDC,  to identify the kind of identity platform                                                                                                                                              |
| `ENVIRONMENT.PATHS.CONFIG_BASE_DIRECTORY`    | config/defaults/                     | Base configuration root path folder for data files and templates to populate them into identity platform                                                                                                          |
| `ENVIRONMENT.PATHS.CONFIG_SECURE_BANKING`    | config/defaults/secure-open-banking/ | Base configuration path folder for specific secure open banking data files and templates to populate them into identity platform                                                                                  |
| `ENVIRONMENT.PATHS.CONFIG_IDENTITY_PLATFORM` | config/defaults/identity-platform/   | Base configuration path folder for generic data files and templates to populate them into identity platform                                                                                                       |
</details>

**Host variables**
<details>
<summary>Table</summary>
<!-- always an empty line before table -->

| Environment variable           | default                        | description                                  |
|--------------------------------|--------------------------------|----------------------------------------------|
| `HOSTS.IDENTITY_PLATFORM_FQDN` | iam.dev.forgerock.financial    | Identity platform Full Qualified Domain Name |
| `HOSTS.RS_FQDN`                | rs.dev.forgerock.financial     | RS Full Qualified Domain Name                |           |
| `HOSTS.SCHEME`                 | https                          | URI scheme, Syntax part of a generic URI     |
</details>


**Users variables**
<details>
<summary>Table</summary>
<!-- always an empty line before table -->

| Environment variable       | default                        | description                                                               |
|----------------------------|--------------------------------|---------------------------------------------------------------------------|
| `USERS.FR_PLATFORM_ADMIN_USERNAME` | amadmin                        | Identity platform Username with admin grants (must exist previously)      |
| `USERS.FR_PLATFORM_ADMIN_PASSWORD` | add-here-the-user-password     | Identity platform User password with admin grants (must exist previously) |
| `USERS.PSU_USERNAME`       | add-here-the-psu-user-name     | Psu Username to (It will be created)                                      |
| `USERS.PSU_PASSWORD`       | add-here-the-psu-user-password | Psu user password (It will be created)                                    |
</details>


**Namespace variables**

| Environment variable | default                                              | description                                                     |
|----------------------|------------------------------------------------------|-----------------------------------------------------------------|
| `NAMESPACE         | ns-env | Developer namespace to populate PSU data |


## Kubernetes ConfigMap
You can override all identity platform configuration files with config predefined within a kubernetes config map.

> :warning: If a path variable as is set to the default relative path of `config/defaults/` then default pre-baked configuration json objects will be used and not your mounted ConfigMap

### ConfigMap for identity platform files mount example

```
spec:
  volumes:
  - name: ob-defaults-objects
    configMap:
      name: openbanking-objects
  containers:
  - name: init-container
    env:
    - name: ENVIRONMENT.PATHS.CONFIG_BASE_DIRECTORY
      value: /opt/config/
    volumeMounts:
    - mountPath: /opt/config/
      name: ob-managed-objects
      readOnly: true
    - name: ENVIRONMENT.PATHS.CONFIG_SECURE_BANKING
      value: /opt/config/secure-open-banking/
    volumeMounts:
    - mountPath: /opt/config/secure-open-banking/
      name: ob-managed-objects
      readOnly: true
    - name: ENVIRONMENT.PATHS.CONFIG_IDENTITY_PLATFORM
      value: /opt/config/identity-platform/
    volumeMounts:
    - mountPath: /opt/config/identity-platform/
      name: ob-managed-objects
      readOnly: true      
```

## Running tests
The tests run against a mockserver which is supplied by [Pact](https://docs.pact.io/). It is used specifically to test internal logic rather than to verify the provider contract.
running the `make test-ci` target will download the required binaries to be able to run the pact tests. this target is used for github actions but can work locally too (if you do not have the pact bonaries installed)

## Temporary Patch
Creation of PSU user on AM and populate the user data to RS service for each environment.
- For functional test purposes @See /rs folder.

### Commands
| Command             | description                                                                                                          |
|---------------------|----------------------------------------------------------------------------------------------------------------------|
| `go mod tidy`       | add missing and remove unused modules                                                                                |
| `go build -o setup` | compiles the packages named by the import paths, along with their dependencies, but it does not install the results. |
| `go run`            | compiles and runs the named main Go package                                                                          |
| `./setup`           | run the compiled program                                                                                             |

> For more information about go command `go help`
