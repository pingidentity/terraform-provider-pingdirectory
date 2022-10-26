# pingdirectory-terraform-poc
This repository contains a POC provider that manages users in PingDirectory (or any other LDAP server). This isn't really a use case that makes sense for Terraform but it's just for getting familiar with writing a provider.

It's built with [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework).

# Testing locally
## Prerequisites
The following must be installed locally to run the provider
- Go 1.18+
- Terraform

To run the example in this repository, you will also need
- Docker and Docker Compose

### Installing required Go modules
Run the following commands to install the required Go modules locally
- `go get github.com/hashicorp/terraform-plugin-framework@latest`
- `go get github.com/go-ldap/ldap/v3@latest`
- `go get github.com/hashicorp/terraform-plugin-log@latest`

Then tidy the modules from the root of the repository:
`go mod tidy`

## Preparing your Terraform environment for running locally-built providers
By default Terraform will attempt to pull providers from remote registries. Update the `~/.terraformrc` file to allow using this POC provider locally.

First, find the GOBIN path where Go installs your binaries. Your path may vary depending on how your Go environment variables are configured.
```
$ go env GOBIN
/Users/<Username>/go/bin
```

If the GOBIN go environment variable is not set, use the default path, /Users/<Username>/go/bin. Create a `~/.terraformrc` file with the following contentes. Change the <PATH> value to the value returned from `go env GOBIN`.

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
Run `go install .` or `make install` to install the provider locally.

## Running an example
### Starting the PingDirectory server
Start a PingDirectory server running locally with the provided docker-compose.yaml file. Change to the `docker-compose` directory and run `docker compose up`. The server will take a couple minutes to become ready. When you see
```
docker-compose-pingdirectory-1  | Replication will not be configured.
docker-compose-pingdirectory-1  | Setting Server to Available
```
in the terminal, the server is ready to receive requests.

### Running Terraform
Change to the `examples/simple-example` directory. The `main.tf` file in this directory defines the Terraform configuration.

Run `terraform plan` to view what changes will be made by Terraform. Run `terraform apply` to apply them.

You can verify the user is created on the PingDirectory server:
```
docker exec -ti docker-compose-pingdirectory-1 ldapsearch --baseDN ou=people,dc=example,dc=com --searchscope sub "(objectClass=person)" dn description
```

You can make changes to the user and use `terraform apply` to apply them, and use the above command to view those changes in PingDirectory.

Run `terraform destroy` to destroy any users managed by Terraform.

