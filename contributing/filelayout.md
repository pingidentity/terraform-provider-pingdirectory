# PingDirectory Terraform Provider Repository Layout

## Overall structure

This project is split into a few **Go** packages to make things easier to use. The **main.go** and **go.mod** files are at the top level, while the rest of the packages are in the **internal/** folder.

The provider package contains only the basic **provider.go** file.

Finally, the actual **configuration** object resources are contained under the **internal/resource/** folder:
- **internal/resource/config/** contains the configuration object resources

The **internal/resource/config** folder will have two types of resources:
- Resources that only support a single type, such as **Location** and **Global Configuration**.  These resources will be found directly in the **config** folder.
- Resources with an API that manages multiple types, such at the **Trust Manager**.  PingDirectory supports multiple types of trust manager provider (Blind Trust, Third Party, JVM Default, and File Based).  In these cases, each type is located in a separate sub-package folder that is named for the type.

## Acceptance Tests

Tests are under the **internal/acctest** folder. The ***acctest.go*** file contains functions used across the acceptance tests. Tests for each resource are located in a separate file, such as ***location_resource_test.go***.

## Other packages

- **internal/configvalidators**: Custom config validators
- **internal/operations**: PingDirectory operations
- **internal/tools**: Defines tools needed by the project but not required elsewhere in the code
- **internal/types**: Utilities for handling types
- **internal/version**: Utilities for handling PingDirectory versions

## Non-Go code structure

- **scripts/**: This folder contains utility scripts.
- **examples/**: This folder contains Terraform examples that can be used to try out the Provider.
- **docs/**: This folder contains documentation.
- **tempmlates/**: This folder contains templates for generating documentation.
- **docker-compose/**:  This folder contains a ***docker-compose.yaml*** file that can be used to quickly set up a PingDirectory server for testing. This method provides a quick way to perform local testing.
- **.vscode/**: This folder contains the configuration for debugging with Visual Studio Code - see the *Debugging* section below.

## Debugging

### VSCode

If you want to debug and step through breakpoints using VSCode, you can use the debug configuration provided in this repository. The [development.md](development.md) file describes how to run the debugger in detail.

### Debugging with tflog output

The provider code includes many debug messages written with **tflog**, the logging package for the Terraform plugin framework.  This package can provide detail on the requests that are being sent and responses that are being returned from the configuration API. When debugging in VSCode, these messages will be written to the Debug Console. If you want to see these messages written to *stderr* without running the debugger, see the logging guide for the terraform plugin framework at https://developer.hashicorp.com/terraform/plugin/log/managing.

### Debugging with PingDirectory logs

If you want to look through the PingDirectory logs for configuration changes that have been made to the server, there is a **logs/config-audit.log** file in the server root that contains a history of all configuration changes made to the server. This file may be useful to see the API requests actually being applied to the server. In our PingDirectory Docker images, this file would be located at **/opt/out/instance/logs/config-audit.log**.
