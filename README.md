# PingDirectory Terraform provider
This repository contains a Terraform provider that manages PingDirectory configuration.

It's built with [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework).

# Running the provider locally
## Prerequisites
The following must be installed locally to run the provider
- Go 1.18+
- Terraform

To run the example in this repository, you will also need
- Docker and Docker Compose

### Installing the PingDirectory Go client
The PingDirectory provider relies on the [PingDirectory Go Client](https://gitlab.corp.pingidentity.com/henryrecker/pingdata-config-api-go-client), which is generated from an OpenAPI specification YAML file.

The `go.mod` file lists this module as `github.com/pingidentity/pingdata-config-api-go-client`, but uses `replace` to point it to a local path rather than an actual downloaded module.

```
replace github.com/pingidentity/pingdata-config-api-go-client v0.0.0 => ../pingdata-config-api-go-client
```

Because this `replace` path points to `../pingdata-config-api-go-client`, you will need to clone the client repo and place it alongside this repo in your filesystem.

### Installing required Go modules
Run the following commands to install the required Go modules locally
```
go get github.com/hashicorp/terraform-plugin-framework@latest
go get github.com/go-ldap/ldap/v3@latest
go get github.com/hashicorp/terraform-plugin-log@latest
```

Then tidy the modules from the root of the repository:

`go mod tidy`

## Preparing your Terraform environment for running locally-built providers
By default Terraform will attempt to pull providers from remote registries. Update the `~/.terraformrc` file to allow using this provider locally.

First, find the GOBIN path where Go installs your binaries. Your path may vary depending on how your Go environment variables are configured.
```
$ go env GOBIN
/Users/<Username>/go/bin
```

If the GOBIN go environment variable is not set, use the default path, /Users/\<Username\>/go/bin. Create a `~/.terraformrc` file with the following contents. Change the \<PATH\> value to the value returned from `go env GOBIN`.

```
provider_installation {
  dev_overrides {
    "pingidentity.com/terraform/pingdirectory" = "<PATH>"
  }
  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```

## Install the provider
Run `make install` (or just `make`) to install the provider locally.

## Running an example
### Starting the PingDirectory server
Start a PingDirectory server running locally with the provided docker-compose.yaml file. Change to the `docker-compose` directory and run `docker compose up`. The server will take a couple minutes to become ready. When you see
```
docker-compose-pingdirectory-1  | Replication will not be configured.
docker-compose-pingdirectory-1  | Setting Server to Available
```
in the terminal, the server is ready to receive requests.

### Running Terraform
Change to the `examples/location-example` directory. The `main.tf` file in this directory defines the Terraform configuration.

Run `terraform plan` to view what changes will be made by Terraform. Run `terraform apply` to apply them.

You can verify the location is created on the PingDirectory server:

```
docker exec -ti docker-compose-pingdirectory-1 dsconfig list-locations
```

```
docker exec -ti docker-compose-pingdirectory-1 dsconfig get-location-prop --location-name Drangleic --property description
```

You can make changes to the location and use `terraform apply` to apply them, and use the above commands to view those changes in PingDirectory.

Run `terraform destroy` to destroy any objects managed by Terraform.

## Debugging with VSCode
You can attach a debugger to the provider with VSCode. The `.vscode/launch.json` file defines the debug configuration.

To debug the provider, go to Run->Start Debugging. Then, open the Debug Console and wait for a message like this:

```
Provider started. To attach Terraform CLI, set the TF_REATTACH_PROVIDERS environment variable with the following:

	TF_REATTACH_PROVIDERS='{"pingidentity.com/terraform/pingdirectory":{"Protocol":"grpc","ProtocolVersion":6,"Pid":53173,"Test":true,"Addr":{"Network":"unix","String":"/var/folders/m8/hpzxbdws7rdgj3cc21vrb1jw0000gn/T/plugin3225934397"}}}'
```

You can then use this to attach the debugger to command-line terraform commands by pasting this line before each command.

```
$ TF_REATTACH_PROVIDERS='{"pingidentity.com/terraform/pingdirectory":{"Protocol":"grpc","ProtocolVersion":6,"Pid":53173,"Test":true,"Addr":{"Network":"unix","String":"/var/folders/m8/hpzxbdws7rdgj3cc21vrb1jw0000gn/T/plugin3225934397"}}}' terraform apply
```

Note that the `TF_REATTACH_PROVIDERS` variable changes each time you run the debugger. You will need to copy it each time you start a new debugger.

# Running acceptance tests
Acceptance tests for the provider use a local PingDirectory instance running in Docker. The following `make` targets will help with running acceptance tests:

- `make testacc`: Runs the acceptance tests, with the assumption that a local PingDirectory instance is available
- `make starttestcontainer`: Starts a PingDirectory Docker container and waits for it to become ready
- `make removetestcontainer`: Stops and removes the PingDirectory Docker container used for testing
- `make testacccomplete`: Starts the PingDirectory Docker container, waits for it to become ready, runs the acceptance tests, and then removes the Docker container. This is good for running the tests from scratch, but you will have to wait for the container startup each time. If you plan on running the tests multiple times and don't mind reusing the same server, then it is better to use the previous three targets individually
