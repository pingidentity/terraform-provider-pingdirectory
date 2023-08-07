---
page_title: "pingdirectory_local_db_indexes Data Source - terraform-provider-pingdirectory"
subcategory: "Local Db Index"
description: |-
  Lists Local Db Index objects in the server configuration.
---

# pingdirectory_local_db_indexes (Data Source)

Lists Local Db Index objects in the server configuration.

## Example Usage

```terraform
terraform {
  required_version = ">=1.1"
  required_providers {
    pingdirectory = {
      version = "~> 0.3.0"
      source  = "pingidentity/pingdirectory"
    }
  }
}

provider "pingdirectory" {
  username   = "cn=administrator"
  password   = "2FederateM0re"
  https_host = "https://localhost:1443"
  # Warning: The insecure_trust_all_tls attribute configures the provider to trust any certificate presented by the PingDirectory server.
  # It should not be used in production. If you need to specify trusted CA certificates, use the
  # ca_certificate_pem_files attribute to point to any number of trusted CA certificate files
  # in PEM format. If you do not specify certificates, the host's default root CA set will be used.
  # Example:
  # ca_certificate_pem_files = ["/example/path/to/cacert1.pem", "/example/path/to/cacert2.pem"]
  insecure_trust_all_tls = true
  product_version        = "9.3.0.0"
}

data "pingdirectory_local_db_indexes" "list" {
  backend_name = "MyBackend"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `backend_name` (String) Name of the parent Backend

### Optional

- `filter` (String) SCIM filter used when searching the configuration.

### Read-Only

- `id` (String) The ID of this resource.
- `ids` (Set of String) Local Db Index IDs found in the configuration
