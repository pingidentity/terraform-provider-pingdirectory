---
page_title: "pingdirectory_access_control_handler Data Source - terraform-provider-pingdirectory"
subcategory: "Access Control Handler"
description: |-
  Describes a Access Control Handler.
---

# pingdirectory_access_control_handler (Data Source)

Describes a Access Control Handler.

Access Control Handlers manage the application-wide access control. The Directory Server access control handler is defined through an extensible interface, so that alternate implementations can be created. Only one access control handler may be active in the server at any given time.

## Example Usage

```terraform
data "pingdirectory_access_control_handler" "myAccessControlHandler" {
}
```

## Documentation
See the [PingDirectory documentation](https://docs.pingidentity.com/r/en-us/pingdirectory-93/pd_sec_define_global_acis)

<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `allowed_bind_control` (Set of String) Specifies a set of controls that clients should be allowed to include in bind requests. As bind requests are evaluated as the unauthenticated user, any controls included in this set will be permitted for any bind attempt. If you wish to grant permission for any bind controls not listed here, then the allowed-bind-control-oid property may be used to accomplish that.
- `allowed_bind_control_oid` (Set of String) Specifies the OIDs of any additional controls (not covered by the allowed-bind-control property) that should be permitted in bind requests.
- `enabled` (Boolean) Indicates whether this Access Control Handler is enabled. If set to FALSE, then no access control is enforced, and any client (including unauthenticated or anonymous clients) could be allowed to perform any operation if not subject to other restrictions, such as those enforced by the privilege subsystem.
- `global_aci` (Set of String) Defines global access control rules.
- `id` (String) The ID of this resource.
- `type` (String) The type of Access Control Handler resource. Options are ['dsee-compat']

