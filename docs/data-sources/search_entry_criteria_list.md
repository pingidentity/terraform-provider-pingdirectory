---
page_title: "pingdirectory_search_entry_criteria_list Data Source - terraform-provider-pingdirectory"
subcategory: "Search Entry Criteria"
description: |-
  Lists Search Entry Criteria objects in the server configuration.
---

# pingdirectory_search_entry_criteria_list (Data Source)

Lists Search Entry Criteria objects in the server configuration.

Search Entry Criteria define sets of criteria for grouping and describing search result entries based on a number of properties, including properties of the associated client connection and operation request, the entry location and contents, and included controls.

## Example Usage

```terraform
data "pingdirectory_search_entry_criteria_list" "list" {
}
```

## Documentation
See the [PingDirectory documentation](https://docs.pingidentity.com/r/en-us/pingdirectory-93/pd_sec_search_entry_criteria)

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `filter` (String) SCIM filter used when searching the configuration.

### Read-Only

- `id` (String) The ID of this resource.
- `objects` (Set of Object) Search Entry Criteria objects found in the configuration (see [below for nested schema](#nestedatt--objects))

<a id="nestedatt--objects"></a>
### Nested Schema for `objects`

Read-Only:

- `id` (String)
- `type` (String)

