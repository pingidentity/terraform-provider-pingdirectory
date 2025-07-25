terraform {
  required_version = ">=1.1"
  required_providers {
    pingdirectory = {
      version = "~> 1.0.0"
      source  = "pingidentity/pingdirectory"
    }
  }
}

provider "pingdirectory" {
  username   = "cn=administrator"
  password   = "2FederateM0re"
  https_host = "https://localhost:1443"
  # Warning: The insecure_trust_all_tls attribute configures the provider to trust any certificate presented by the PingDirectory server.
  insecure_trust_all_tls = true
  product_version        = "10.2.0.0"
}

# Create a sample location
resource "pingdirectory_location" "myLocation" {
  name        = "MyLocation"
  description = "My description"
}

# Update the default global configuration to use the created location, and to enable encryption
resource "pingdirectory_default_global_configuration" "global" {
  location     = pingdirectory_location.myLocation.id
  encrypt_data = true
}
