---
page_title: "pingdirectory_local_db_vlv_index Data Source - terraform-provider-pingdirectory"
subcategory: "Local Db Vlv Index"
description: |-
  Describes a Local Db Vlv Index.
---

# pingdirectory_local_db_vlv_index (Data Source)

Describes a Local Db Vlv Index.

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

data "pingdirectory_local_db_vlv_index" "myLocalDbVlvIndex" {
  backend_name = "MyBackend"
  name         = "myLocalDbVlvIndex"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `backend_name` (String) Name of the parent Backend
- `name` (String) Specifies a unique name for this VLV index.

### Read-Only

- `base_dn` (String) Specifies the base DN used in the search query that is being indexed.
- `cache_mode` (String) Specifies the cache mode that should be used when accessing the records in the database for this index.
- `filter` (String) Specifies the LDAP filter used in the query that is being indexed.
- `id` (String) Name of this object.
- `max_block_size` (Number) Specifies the number of entry IDs to store in a single sorted set before it must be split.
- `scope` (String) Specifies the LDAP scope of the query that is being indexed.
- `sort_order` (String) Specifies the names of the attributes that are used to sort the entries for the query being indexed.
