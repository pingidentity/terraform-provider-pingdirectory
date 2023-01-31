# Development Environment Setup
## Prerequisites
The following applications must be installed locally to run the provider:
- [Terraform](https://www.terraform.io/downloads.html) 1.0+ (to run acceptance tests)
- [Go](https://golang.org/doc/install) 1.18+ (to build the provider plugin)

To run the example in this repository, you will also need
- [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/install/)

## File system overview

Before starting, take a moment to familiarize yourself with the structure of this repository, [found here](filelayout.md)

## Preparing your Terraform environment to run locally-built providers
By default, Terraform attempts to pull providers from remote registries. This behavior can be overwritten by modifying the `~/.terraformrc` file to enable the local use of this provider. This configuration can be used to test in-development changes to the provider.

First, find the **GOBIN** path where Go installs your binaries. Your path may vary depending on how your Go environment variables are configured.

```sh
$ go env GOBIN
/Users/<Username>/go/bin
```

If the GOBIN go environment variable is not set, use the default path, **/Users/\<Username\>/go/bin**. Create a `~/.terraformrc` file with the following contents, changing the \<PATH\> value to the value returned from `go env GOBIN`.

```text
provider_installation {
  dev_overrides {
    "pingidentity/pingdirectory" = "<PATH>"
  }
  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this next line, Terraform will _only_ use
  # the dev_overrides block, meaning no other providers will be available.
  direct {}
}
```

### Running a locally-built PingDirectory Go client
The PingDirectory Terraform provider relies on the [PingDirectory Go Client](https://github.com/pingidentity/pingdirectory-go-client).

If changes are needed in the Go client, the `replace` command in the `go.mod` file can be used to point to a modified local Go client while testing.

```txt
replace github.com/pingidentity/pingdirectory-go-client v0.0.0 with ==> ./pingdirectory-go-client
```

In this example above, the `replace` path points to `../pingdirectory-go-client`, meaning you would need to clone the client repo and place it alongside this repo in your filesystem.

## Install the provider
Run `go mod tidy` to get any required dependencies.

Run `make install` (or just `make`) to install the provider locally.

## Running acceptance tests
Acceptance tests for the provider use a local PingDirectory instance running in Docker. The following `make` targets will help with running acceptance tests:

- `make testacc`: Runs the acceptance tests, with the assumption that a local PingDirectory instance is available
- `make starttestcontainer`: Starts a PingDirectory Docker container and waits for it to become ready
- `make removetestcontainer`: Stops and removes the PingDirectory Docker container used for testing
- `make testacccomplete`: Starts the PingDirectory Docker container, waits for it to become ready, runs the acceptance tests, and then removes the Docker container. This option is good for running the tests from scratch and for use in automation, but you will have to wait for the container startup each time.
  
**Tip**: If you plan on running tests multiple times and do not mind reusing the same server, then it is recommended to use the first three options above to perform each step individually.

## Run an example
### Start the PingDirectory server
Start a PingDirectory server running locally with the provided **docker-compose.yaml** file. Change to the `docker-compose` directory and run `docker compose up`. (Alternatively, use the `make starttestcontainer` command from the previous section.) The server will take a couple of minutes to become ready. When you see the following output in the terminal, the server is ready to process requests:
```
docker-compose-pingdirectory-1  | Replication will not be configured.
docker-compose-pingdirectory-1  | Setting Server to Available
```

### Run Terraform
Change to the `examples/location-example` directory. The `main.tf` file in this directory defines the Terraform configuration.

Run `terraform plan` to view what changes will be made by Terraform. Run `terraform apply` to apply them.

You can verify the location is created on the PingDirectory server:

```sh
docker exec -ti docker-compose-pingdirectory-1 dsconfig list-locations
```

```sh
docker exec -ti docker-compose-pingdirectory-1 dsconfig get-location-prop --location-name Drangleic --property description
```

You can make changes to the location and use `terraform apply` to apply them, and use the above commands to view those changes in PingDirectory.

Run `terraform destroy` to destroy any objects managed by Terraform.

## Debugging with VSCode
You can attach a debugger to the provider with VSCode. The `.vscode/launch.json` file defines the debug configuration.

To debug the provider, navigate to **Run > Start Debugging**. Then, open the Debug Console and wait for a message like this:

```text
Provider started. To attach Terraform CLI, set the TF_REATTACH_PROVIDERS environment variable with the following:

	TF_REATTACH_PROVIDERS='{"registry.terraform.io/pingidentity/pingdirectory":{"Protocol":"grpc","ProtocolVersion":6,"Pid":53173,"Test":true,"Addr":{"Network":"unix","String":"/var/folders/m8/hpzxbdws7rdgj3cc21vrb1jw0000gn/T/plugin3225934397"}}}'
```

You can then use this to attach the debugger to command-line terraform commands by pasting this line before each command.

```sh
$ TF_REATTACH_PROVIDERS='{"registry.terraform.io/pingidentity/pingdirectory":{"Protocol":"grpc","ProtocolVersion":6,"Pid":53173,"Test":true,"Addr":{"Network":"unix","String":"/var/folders/m8/hpzxbdws7rdgj3cc21vrb1jw0000gn/T/plugin3225934397"}}}' terraform apply
```

**Note**: The `TF_REATTACH_PROVIDERS` variable changes each time you run the debugger. You will need to copy the output for use each time you start a new debugger.
