---
page_title: "pingdirectory_virtual_attribute Data Source - terraform-provider-pingdirectory"
subcategory: "Virtual Attribute"
description: |-
  Describes a Virtual Attribute.
---

# pingdirectory_virtual_attribute (Data Source)

Describes a Virtual Attribute.

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

data "pingdirectory_virtual_attribute" "myVirtualAttribute" {
  id = "MyVirtualAttribute"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) Name of this object.

### Read-Only

- `allow_index_conflicts` (Boolean) Indicates whether the server should allow creating or altering this virtual attribute definition even if it conflicts with one or more indexes defined in the server.
- `allow_retrieving_membership` (Boolean) Indicates whether to handle requests that request all values for the virtual attribute.
- `attribute_type` (String) Specifies the attribute type for the attribute whose values are to be dynamically assigned by the virtual attribute.
- `base_dn` (Set of String) Specifies the base DNs for the branches containing entries that are eligible to use this virtual attribute.
- `bypass_access_control_for_searches` (Boolean) Indicates whether searches performed by this virtual attribute provider should be exempted from access control restrictions.
- `client_connection_policy` (Set of String) Specifies a set of client connection policies for which this Virtual Attribute should be generated. If this is undefined, then this Virtual Attribute will always be generated. If it is associated with one or more client connection policies, then this Virtual Attribute will be generated only for operations requested by clients assigned to one of those client connection policies.
- `conflict_behavior` (String) Specifies the behavior that the server is to exhibit for entries that already contain one or more real values for the associated attribute.
- `description` (String) A description for this Virtual Attribute
- `direct_memberships_only` (Boolean) Specifies whether to only include groups in which the user is directly associated with and the membership maybe modified via the group entry. Groups in which the user's membership is derived dynamically or through nested groups will not be included.
- `enabled` (Boolean) Indicates whether the Virtual Attribute is enabled for use.
- `exclude_operational_attributes` (Boolean) Indicates whether all operational attributes should be excluded from the generated checksum.
- `excluded_attribute` (Set of String) Specifies the attributes that should be excluded from the checksum calculation.
- `extension_argument` (Set of String) The set of arguments used to customize the behavior for the Third Party Virtual Attribute. Each configuration property should be given in the form 'name=value'.
- `extension_class` (String) The fully-qualified name of the Java class providing the logic for the Third Party Virtual Attribute.
- `filter` (Set of String) Specifies the search filters to be applied against entries to determine if the virtual attribute is to be generated for those entries.
- `group_dn` (Set of String) Specifies the DNs of the groups whose members can be eligible to use this virtual attribute.
- `include_milliseconds` (Boolean) Indicates whether the current time includes millisecond precision.
- `included_group_filter` (String) A search filter that will be used to identify which groups should be included in the values of the virtual attribute. With no value defined (which is the default behavior), all groups that contain the target user will be included.
- `join_attribute` (Set of String) An optional set of the names of the attributes to include with joined entries.
- `join_base_dn_type` (String) Specifies how server should determine the base DN for the internal searches used to identify joined entries.
- `join_custom_base_dn` (String) The fixed, administrator-specified base DN for the internal searches used to identify joined entries.
- `join_dn_attribute` (String) The attribute in related entries whose set of values must contain the DN of the search result entry to be joined with that entry.
- `join_filter` (String) An optional filter that specifies additional criteria for identifying joined entries. If a join-filter value is specified, then only entries matching that filter (in addition to satisfying the other join criteria) will be joined with the search result entry.
- `join_match_all` (Boolean) Indicates whether joined entries will be required to have all values for the source attribute, or only at least one of its values.
- `join_scope` (String) The scope for searches used to identify joined entries.
- `join_size_limit` (Number) The maximum number of entries that may be joined with the source entry, which also corresponds to the maximum number of values that the virtual attribute provider will generate for an entry.
- `join_source_attribute` (String) The attribute containing the value(s) in the source entry to use to identify related entries.
- `join_target_attribute` (String) The attribute in target entries whose value(s) match values of the source attribute in the source entry.
- `multiple_virtual_attribute_evaluation_order_index` (Number) Specifies the order in which virtual attribute definitions for the same attribute type will be evaluated when generating values for an entry.
- `multiple_virtual_attribute_merge_behavior` (String) Specifies the behavior that will be exhibited for cases in which multiple virtual attribute definitions apply to the same multivalued attribute type. This will be ignored for single-valued attribute types.
- `reference_search_base_dn` (Set of String) The base DN that will be used when searching for references to the target entry. If no reference search base DN is specified, the default behavior will be to search below all public naming contexts configured in the server.
- `referenced_by_attribute` (Set of String) The name or OID of an attribute type whose values will be searched for references to the target entry. The attribute type must be defined in the server schema, must have a syntax of either "distinguished name" or "name and optional UID", and must be indexed for equality.
- `require_explicit_request_by_name` (Boolean) Indicates whether attributes of this type must be explicitly included by name in the list of requested attributes. Note that this will only apply to virtual attributes which are associated with an attribute type that is operational. It will be ignored for virtual attributes associated with a non-operational attribute type.
- `return_utc_time` (Boolean) Indicates whether to return current time in UTC.
- `rewrite_search_filters` (String) Search filters that include Is Member Of Virtual Attribute searches on dynamic groups can be updated to include the dynamic group filter in the search filter itself. This can allow the backend to more efficiently process the search filter by using attribute indexes sooner in the search processing.
- `script_argument` (Set of String) The set of arguments used to customize the behavior for the Scripted Virtual Attribute. Each configuration property should be given in the form 'name=value'.
- `script_class` (String) The fully-qualified name of the Groovy class providing the logic for the Groovy Scripted Virtual Attribute.
- `sequence_number_attribute` (String) Specifies the name or OID of the attribute which contains the sequence number from which unique identifiers are generated. The attribute should have Integer syntax or a String syntax permitting integer values. If this property is modified then the filter property should be updated accordingly so that only entries containing the sequence number attribute are eligible to have a value generated for this virtual attribute.
- `source_attribute` (String) Specifies the source attribute containing the values to use for this virtual attribute.
- `source_entry_dn_attribute` (String) Specifies the attribute containing the DN of another entry from which to obtain the source attribute providing the values for this virtual attribute.
- `source_entry_dn_map` (String) Specifies a DN map that will be used to identify the entry from which to obtain the source attribute providing the values for this virtual attribute.
- `type` (String) The type of Virtual Attribute resource. Options are ['mirror', 'entry-checksum', 'member-of-server-group', 'constructed', 'is-member-of', 'custom', 'num-subordinates', 'reverse-dn-join', 'identify-references', 'user-defined', 'current-time', 'short-unique-id', 'entry-dn', 'has-subordinates', 'equality-join', 'groovy-scripted', 'instance-name', 'replication-state-detail', 'member', 'password-policy-state-json', 'subschema-subentry', 'dn-join', 'third-party']
- `value` (Set of String) Specifies the values to be included in the virtual attribute.
- `value_pattern` (Set of String) Specifies a pattern for constructing the virtual attribute value using fixed text and attribute values from the entry.
