---
page_title: "pingdirectory_constructed_attributes Data Source - terraform-provider-pingdirectory"
subcategory: "Constructed Attribute"
description: |-
  Lists Constructed Attribute objects in the server configuration.
---

# pingdirectory_constructed_attributes (Data Source)

Lists Constructed Attribute objects in the server configuration.

A constructed attribute constructs values for an attribute by using a combination of fixed text and values of other attributes from the entry. Note that just creating one of these objects will not have any effect. The object must be referenced from another configuration object such as a Delegated Admin Resource Type.

## Example Usage

```terraform
data "pingdirectory_constructed_attributes" "list" {
}
```

## Documentation
See the [PingDirectory documentation](https://docs.pingidentity.com/r/en-us/pingdirectory-93/pd_da_config_attr_search_pingdir_server)

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `filter` (String) SCIM filter used when searching the configuration.

### Read-Only

- `id` (String) The ID of this resource.
- `ids` (Set of String) Constructed Attribute IDs found in the configuration

