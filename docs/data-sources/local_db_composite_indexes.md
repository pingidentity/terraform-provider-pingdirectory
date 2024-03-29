---
page_title: "pingdirectory_local_db_composite_indexes Data Source - terraform-provider-pingdirectory"
subcategory: "Local Db Composite Index"
description: |-
  Lists Local Db Composite Index objects in the server configuration.
---

# pingdirectory_local_db_composite_indexes (Data Source)

Lists Local Db Composite Index objects in the server configuration.

Local DB Composite Indexes may be used to define an index based on a filter pattern and an optional base DN pattern.

## Example Usage

```terraform
data "pingdirectory_local_db_composite_indexes" "list" {
  backend_name = "MyBackend"
}
```

## Documentation
See the [PingDirectory documentation](https://docs.pingidentity.com/r/en-us/pingdirectory-93/pd_ds_composite_indexes)

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `backend_name` (String) Name of the parent Backend

### Optional

- `filter` (String) SCIM filter used when searching the configuration.

### Read-Only

- `id` (String) The ID of this resource.
- `ids` (Set of String) Local Db Composite Index IDs found in the configuration

