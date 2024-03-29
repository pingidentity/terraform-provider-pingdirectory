---
page_title: "pingdirectory_replication_assurance_policies Data Source - terraform-provider-pingdirectory"
subcategory: "Replication Assurance Policy"
description: |-
  Lists Replication Assurance Policy objects in the server configuration.
---

# pingdirectory_replication_assurance_policies (Data Source)

Lists Replication Assurance Policy objects in the server configuration.

A Replication Assurance Policy is used to specify the local and remote replication assurance levels and a timeout to use for update operations. Optionally, request and connection criteria can be configured in the policy to allow matching a policy to requests that satisfy such criteria.

## Example Usage

```terraform
data "pingdirectory_replication_assurance_policies" "list" {
}
```

## Documentation
See the [PingDirectory documentation](https://docs.pingidentity.com/r/en-us/pingdirectory-93/pd_ds_config_assured_replication)

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `filter` (String) SCIM filter used when searching the configuration.

### Read-Only

- `id` (String) The ID of this resource.
- `ids` (Set of String) Replication Assurance Policy IDs found in the configuration

