---
page_title: "pingdirectory_local_db_vlv_index Data Source - terraform-provider-pingdirectory"
subcategory: "Local Db Vlv Index"
description: |-
  Describes a Local Db Vlv Index.
---

# pingdirectory_local_db_vlv_index (Data Source)

Describes a Local Db Vlv Index.

Local DB VLV Indexes are used to store information about a specific search request that makes it possible to efficiently process them using the VLV control.

## Example Usage

```terraform
data "pingdirectory_local_db_vlv_index" "myLocalDbVlvIndex" {
  backend_name = "MyBackend"
  name         = "myLocalDbVlvIndex"
}
```

## Documentation
See the [PingDirectory documentation](https://docs.pingidentity.com/r/en-us/pingdirectory-93/pd_ds_local_db_vlv_indexes)

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `backend_name` (String) Name of the parent Backend
- `name` (String) Specifies a unique name for this VLV index.

### Read-Only

- `base_dn` (String) Specifies the base DN used in the search query that is being indexed.
- `cache_mode` (String) Specifies the cache mode that should be used when accessing the records in the database for this index.
- `filter` (String) Specifies the LDAP filter used in the query that is being indexed.
- `id` (String) The ID of this resource.
- `max_block_size` (Number) Specifies the number of entry IDs to store in a single sorted set before it must be split.
- `scope` (String) Specifies the LDAP scope of the query that is being indexed.
- `sort_order` (String) Specifies the names of the attributes that are used to sort the entries for the query being indexed.
- `type` (String) The type of Local DB VLV Index resource. Options are ['local-db-vlv-index']

