---
page_title: "pingdirectory_root_dn Data Source - terraform-provider-pingdirectory"
subcategory: "Root Dn"
description: |-
  Describes a Root Dn.
---

# pingdirectory_root_dn (Data Source)

Describes a Root Dn.

The Root DN configuration contains all the Root DN Users defined in the Directory Server. In addition, it also defines the default set of privileges that Root DN Users automatically inherit.

## Example Usage

```terraform
data "pingdirectory_root_dn" "myRootDn" {
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `default_root_privilege_name` (Set of String) Specifies the names of the privileges that root users will be granted by default.
- `id` (String) The ID of this resource.
- `type` (String) The type of Root DN resource. Options are ['root-dn']

