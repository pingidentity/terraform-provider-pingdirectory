---
page_title: "pingdirectory_delegated_admin_resource_rights Data Source - terraform-provider-pingdirectory"
subcategory: "Delegated Admin Resource Rights"
description: |-
  Describes a Delegated Admin Resource Rights.
---

# pingdirectory_delegated_admin_resource_rights (Data Source)

Describes a Delegated Admin Resource Rights.

Delegated Admin Resource Rights give a user, or group of users, authority to manage a specific resource type through the Delegated Admin API.

## Example Usage

```terraform
data "pingdirectory_delegated_admin_resource_rights" "myDelegatedAdminResourceRights" {
  delegated_admin_rights_name = "MyDelegatedAdminRights"
  rest_resource_type          = "myRestResourceType"
}
```

## Documentation
See the [PingDirectory documentation](https://docs.pingidentity.com/r/en-us/pingdirectory-93/pd_da_config_delegated_admin)

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `delegated_admin_rights_name` (String) Name of the parent Delegated Admin Rights
- `rest_resource_type` (String) Specifies the resource type applicable to these Delegated Admin Resource Rights.

### Read-Only

- `admin_permission` (Set of String) Specifies administrator(s) permissions.
- `admin_scope` (String) Specifies the scope of these Delegated Admin Resource Rights.
- `description` (String) A description for this Delegated Admin Resource Rights
- `enabled` (Boolean) Indicates whether these Delegated Admin Resource Rights are enabled.
- `id` (String) The ID of this resource.
- `resource_subtree` (Set of String) Specifies subtrees within the search base whose entries can be managed by the administrator(s). The admin-scope must be set to resources-in-specific-subtrees.
- `resources_in_group` (Set of String) Specifies groups whose members can be managed by the administrator(s). The admin-scope must be set to resources-in-specific-groups.
- `type` (String) The type of Delegated Admin Resource Rights resource. Options are ['delegated-admin-resource-rights']

