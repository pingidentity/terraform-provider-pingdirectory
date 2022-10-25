terraform {
  required_providers {
    hashicups = {
      source = "pingidentity.com/terraform/pingdirectory"
    }
  }
}

provider "pingdirectory" {}

data "pingdirectory_users" "example" {}
