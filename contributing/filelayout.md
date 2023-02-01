# PingDirectory Terraform Provider Repository Layout

## Overall structure

This project is split into a few **Go** packages to make things easier to use. The **main.go** and **go.mod** files are at the top level, while the rest of the packages are in the **internal/** folder.

The provider package contains only the basic **provider.go** file. When adding resources, you will need to update the function near the bottom of this file to add the resource to the list that the provider can manage. For example:

```text
// Resources defines the resources implemented in the provider.
// Maintain alphabetical order for ease of management
func (p *pingdirectoryProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		config.NewGlobalConfigurationResource,
		config.NewLocationResource,
		config.NewReallyNeededResource,  ⇐ example new function in development
		config.NewRootDnResource, 
		serverinstance.NewAuthorizeServerInstanceResource,
		serverinstance.NewDirectoryServerInstanceResource,
		serverinstance.NewProxyServerInstanceResource,
		serverinstance.NewSyncServerInstanceResource,
		trustmanagerprovider.NewBlindTrustManagerProviderResource,
		trustmanagerprovider.NewJVMDefaultTrustManagerProviderResource,
		trustmanagerprovider.NewFileBasedTrustManagerProviderResource,
		trustmanagerprovider.NewThirdPartyTrustManagerProviderResource,
	}
```

Finally, the actual **configuration** object resources are contained under the **resource/** folder:
- **resource/config/** contains the configuration object resources

The **config** folder will have two types of resources:
- Resources that only support a single type, such as **Location** and **Global Configuration**.  These resources will be found directly in the **config** folder.
- Resources with an API that manages multiple types, such at the **Trust Manager**.  PingDirectory supports multiple types of providers (Blind Trust, Third Party, JVM Default, and File Based).  In these cases, each type is located in a separate sub-package folder that is named for the type.

## Acceptance Tests

Tests are under the **acctest** folder. The ***acctest.go*** file contains functions used across the acceptance tests. Tests for each resource are located in a separate file, such as ***location_resource_test.go***.

## Non-Go code structure

- **examples/**: This folder contains Terraform examples that can be used to try out the Provider
- **docker-compose/**:  This folder contains a ***docker-compose.yaml*** file that can be used to quickly set up a PingDirectory server for testing. This method provides a quick way to perform local testing.
- **.vscode/**: This folder contains the configuration for debugging with Visual Studio Code - see the *Debugging* section below.

A partial listing of the files and directories is here:

```text
├── .vscode                      ← debugging
├── contributing                 ← public documentation for development
├── docker-compose               ← stand up local PD instance
├── examples                     ← documentation examples folder
├── go.mod
├── internal
│   ├── acctest                  ← testing folder
│   ├── operations
│   │   └── operation.go
│   ├── provider
│   │   └── provider.go          ← add your modules to the Resources() block
│   ├── resource
│   │   └── config               ← provider resource functionality
│   │       ├── api_utils.go
│   │       ├── common.go
│   │       ├── global_configuration_resource.go
│   │       ├── location_resource.go
│   │       ├── root_dn_resource.go
│   │       ├── serverinstance
│   │       │   ├── authorize_server_instance_resource.go
│   │       │   ├── directory_server_instance_resource.go
│   │       │   ├── proxy_server_instance_resource.go
│   │       │   ├── server_instance_resource_common.go
│   │       │   └── sync_server_instance_resource.go
│   │       └── trustmanagerprovider
│   │           ├── blind_trust_manager_provider_resource.go
│   │           ├── file_based_trust_manager_provider_resource.go
│   │           ├── jvm_default_trust_manager_provider_resource.go
│   │           └── third_party_trust_manager_provider_resource.go
│   └── types                    ← utility type stuff
│       ├── conversion.go
│       ├── definitions.go
│       └── utils.go
└── main.go                      ← standard Go convention
```

After the resource is implemented, add support for it in **provider.go**. When this has been done, you can rebuild the provider with `go install .` and test using the new resource in Terraform.

## Debugging

### VSCode

If you want to debug and step through breakpoints using VSCode, you can use the debug configuration provided in this repository. The [development.md](development.md) file describes how to run the debugger in detail.

### Debugging with tflog output

The provider code includes many debug messages written with **tflog**, the logging package for the Terraform plugin framework.  This package can provide detail on the requests that are being sent and responses that are being returned from the configuration API. When debugging in VSCode, these messages will be written to the Debug Console. If you want to see these messages written to *stderr* without running the debugger, see the logging guide for the terraform plugin framework at https://developer.hashicorp.com/terraform/plugin/log/managing.

### Debugging with PingDirectory logs

If you want to look through the PingDirectory logs for configuration changes that have been made to the server, there is a **logs/config-audit.log** file in the server root that contains a history of all configuration changes made to the server. This file may be useful to see the API requests actually being applied to the server. In our PingDirectory Docker images, this file would be located at **/opt/out/instance/logs/config-audit.log**.
