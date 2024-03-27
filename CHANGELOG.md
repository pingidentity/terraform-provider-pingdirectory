# v1.2.0 (Unreleased)
### ENHANCEMENTS
* Added validators for enum string attributes (#337)
* Allow `product_version` values that are not explicitly supported as long as the major-minor version is supported. For example, version `10.0.0.10` would be allowed, but version `20.0.0.0` would not be allowed. (#349)
* Added support for PingDirectory patch versions (#350)
  * `9.1.0.3` - `9.1.0.4`
  * `9.2.0.2` - `9.2.0.5`
  * `9.3.0.1` - `9.3.0.5`
  * `10.0.0.1` - `10.0.0.2`

# v1.1.0 February 1, 2024
### FEATURES
* Added support for PingDirectory version `10.0.0.0` (#331)

# v1.0.2 January 18, 2024

### DEPENDENCIES
* Bump `golang.org/x/crypto` from `0.16.0` to `0.17.0` (#324)
* Bump `github.com/cloudflare/circl` from `1.3.3` to `1.3.7` (#325)
* Bump `github.com/hashicorp/terraform-plugin-framework` from `1.3.5` to `1.5.0` (#327)
* Bump `github.com/hashicorp/terraform-plugin-go` from `0.18.0` to `0.20.0` (#327)
 
# v1.0.1 October 26, 2023

### DEPENDENCIES
* Bump `golang.org/x/net` from `0.14.0` to `0.17.0` (#307)
* Bump `google.golang.org/grpc` from `1.57.0` to `1.57.1` (#310)

# v1.0.0 September 6, 2023

### ENHANCEMENTS
* Mark attributes as required for specific values of the `type` attribute (#289)
* Add default values based on the `type` attribute value, and add default values to schema for empty sets (#291)
* Set computed attributes that don't apply for a given `type` value to null

### BUG FIXES
* Fix mismatched state issues with certain string attributes (#290)
* Reset values that are null in the configuration to default values (#292)
* Fix attributes that required replacement for default resources (#293)

### DOCUMENTATION UPDATES
* Various documentation improvements (#298)
  * Add links to PingDirectory online docs
  * Add config object synopsis to doc pages
  * Add note to doc pages for default resources

### DEPENDENCIES
* Bump `github.com/hashicorp/terraform-plugin-framework-validators` from `0.11.0` to `0.12.0` (#294)
* Bump `github.com/hashicorp/terraform-plugin-testing` from `1.4.0` to `1.5.1` (#295)
* Bump `github.com/terraform-linters/tflint` from `0.47.0` to `0.48.0` (#296)

# v0.11.0 August 25, 2023

### ENHANCEMENTS
* Various improvements to the `type` attribute (#274, #276)
  * Added config validators that depend on specific `type` values
  * Add computed `type` values for resources that have a single valid `type` value
  * Make the `type` value read-only for default resources
  * Separate descriptions for attributes that differ by `type`
* Build in PingDirectory defaults into resource schema where possible (#275)
* Reduce the number of computed attributes that appear in plan output after initial create (#280)
* Add User-Agent header with provider version to HTTP requests (#282, #287)

### BUG FIXES
* Avoid returning errors for 404 responses in resource read and delete methods for non-default resources (#277)
* Require replacement for changes to attributes that are read-only after creation (#281)

### DOCUMENTATION UPDATES
* Remove unnecessary provider config from resource and data source examples (#283)
* Various doc fixes (#286)

### DEPENDENCIES
* Bump `github.com/hashicorp/terraform-plugin-framework` from `1.3.4` to `1.3.5` (#285)

# v0.10.0 August 4, 2023

### FEATURES
* Add a corresponding data source for each resource (#267)
* Add data sources for listing existing config objects (#268)

### ENHANCEMENTS
* Make the `id` attribute read-only across all resources and data sources. Use the `name` attribute instead for resources that previously depended on the `id` attribute being set by the user (#269)

### DEPENDENCIES
* Bump `github.com/hashicorp/terraform-plugin-framework` from `1.3.2` to `1.3.3` (#265)
* Bump `github.com/hashicorp/terraform-plugin-testing` from `1.3.0` to `1.4.0` (#266)

# v0.9.0 July 14, 2023

### FEATURES
* Support PingDirectory version `9.2.0.1` (#257)
* Support PingDirectory version `9.3.0.0` (#258)

### BUG FIXES
* Fix some attributes that were incorrectly version-restricted (#259)
* Fix formatted properties causing errors on create (#261)

### DEPENDENCIES
* Bump `github.com/terraform-linters/tflint` from `0.46.1` to `0.47.0` (#235)
* Bump `github.com/hashicorp/terraform-plugin-framework` from `1.3.1` to `1.3.2` (#254)
* Bump `github.com/hashicorp/terraform-plugin-go` from `0.16.0` to `0.18.0` (#256)
* Bump `github.com/hashicorp/terraform-plugin-docs` from `0.15.0` to `0.16.0` (#260)
* Update to Go `1.20` (#262)

# v0.8.0 June 30 2023

### FEATURES
* **New Resource:** `pingdirectory_custom_logged_stats` (#240)
* **New Resource:** `pingdirectory_prometheus_monitor_attribute_metric` (#241)
* **New Resource:** `pingdirectory_velocity_context_provider` (#242)
* **New Resource:** `pingdirectory_passphrase_provider` (#243)
* **New Resource:** `pingdirectory_notification_manager` (#244)
* **New Resource:** `pingdirectory_vault_authentication_method` (#245)
* **New Resource:** `pingdirectory_obscured_value` (#246)
* **New Resource:** `pingdirectory_license` (#247)
* **New Resource:** `pingdirectory_delegated_admin_correlated_rest_resource` (#248)
* **New Resource:** `pingdirectory_sensitive_attribute` (#249)
* **New Resource:** `pingdirectory_password_validator` (#250)
* **New Resource:** `pingdirectory_password_generator` (#251)
* **New Resource:** `pingdirectory_replication_server` (#252)

### ENHANCEMENTS
* Added config validators for `pingdirectory_delegated_admin_rights`, `pingdirectory_rest_resource_type`, and `pingdirectory_web_application_extension` (#253)

### DEPENDENCIES
* `Bump github.com/hashicorp/terraform-plugin-framework from 1.3.0 to 1.3.1` (#236)
* `Bump github.com/golangci/golangci-lint from 1.53.2 to 1.53.3` (#237)
* `Bump github.com/hashicorp/terraform-plugin-testing from 1.2.0 to 1.3.0` (#238)
* `Bump github.com/hashicorp/terraform-plugin-go from 0.15.0 to 0.16.0` (#239)

# v0.7.0 June 16 2023

### FEATURES
* Collapsed subresources into single parent resources, differentiated by a new `type` attribute. For example, the `pingdirectory_indicator_gauge` and `pingdirectory_numeric_gauge` resources are now combined into a single `pingdirectory_gauge` resource, with a `type` attribute that can be set to either `indicator` or `numeric`. See the documentation for more information. (#221)
* **New Resource:** `pingdirectory_root_dse_backend ` (#206)
* **New Resource:** `pingdirectory_search_entry_criteria ` (#207)
* **New Resource:** `pingdirectory_search_reference_criteria ` (#208)
* **New Resource:** `pingdirectory_server_group` (#209)
* **New Resource:** `pingdirectory_soft_delete_policy ` (#210)
* **New Resource:** `pingdirectory_token_claim_validation` (#211)
* **New Resource:** `pingdirectory_uncached_attribute_criteria` (#212)
* **New Resource:** `pingdirectory_uncached_entry_criteria` (#213)
* **New Resource:** `pingdirectory_pass_through_authentication_handler` (#214)
* **New Resource:** `pingdirectory_trusted_certificate` (#217)
* **New Resource:** `pingdirectory_matching_rule` (#224)
* **New Resource:** `pingdirectory_conjur_authentication_method` (#225)
* **New Resource:** `pingdirectory_inter_server_authentication_method` (#226)
* **New Resource:** `pingdirectory_key_pair` (#227)
* **New Resource:** `pingdirectory_mac_secret_key` (#228)
* **New Resource:** `pingdirectory_otp_delivery_mechanism` (#229)
* **New Resource:** `pingdirectory_password_storage_scheme` (#230)
* **New Resource:** `pingdirectory_server_instance_listener` (#231)
* **New Resource:** `pingdirectory_replication_domain` (#232)
* **New Resource:** `pingdirectory_cipher_secret_key` (#233)

### BUG FIXES
* Fixed an issue where certain default_ resources would fail on first apply (#226)

# v0.6.0 May 23 2023

### FEATURES
* **New Resource:** `pingdirectory_http_configuration` (#161)
* **New Resource:** `pingdirectory_dn_map` (#162)
* **New Resource:** `pingdirectory_result_code_map` (#163)
* **New Resource:** `pingdirectory_attribute_syntax` (#164)
* **New Resource:** `pingdirectory_crypto_manager` (#165)
* **New Resource:** `pingdirectory_azure_authentication_method` (#166)
* **New Resource:** `pingdirectory_log_field_syntax` (#167)
* **New Resource:** `pingdirectory_change_subscription_handler` (#168)
* **New Resource:** `pingdirectory_log_field_behavior` (#169)
* **New Resource:** `pingdirectory_log_field_mapping` (#170)
* **New Resource:** `pingdirectory_log_file_rotation_listener` (#171)
* **New Resource:** `pingdirectory_log_retention_policy` (#172)
* **New Resource:** `pingdirectory_log_rotation_policy` (#173)
* **New Resource:** `pingdirectory_alarm_manager` (#174)
* **New Resource:** `pingdirectory_alert_handler` (#175)
* **New Resource:** `pingdirectory_change_subscription` (#176)
* **New Resource:** `pingdirectory_monitor_provider` (#177)
* **New Resource:** `pingdirectory_replication_assurance_policy` (#178)
* **New Resource:** `pingdirectory_velocity_template_loader` (#179)
* **New Resource:** `pingdirectory_cipher_secret_key` (#180)
* **New Resource:** `pingdirectory_work_queue` (#181)
* **New Resource:** `pingdirectory_client_connection_policy` (#182)
* **New Resource:** `pingdirectory_constructed_attribute` (#183)
* **New Resource:** `pingdirectory_correlated_ldap_data_view` (#184)
* **New Resource:** `pingdirectory_data_security_auditor` (#185)
* **New Resource:** `pingdirectory_delegated_admin_attribute_category` (#186)
* **New Resource:** `pingdirectory_extended_operation_handler` (#188)
* **New Resource:** `pingdirectory_failure_lockout_action` (#189)
* **New Resource:** `pingdirectory_json_attribute_constraints` (#190)
* **New Resource:** `pingdirectory_group_implementation` (#191)
* **New Resource:** `pingdirectory_key_manager_provider` (#192)
* **New Resource:** `pingdirectory_ldap_sdk_debug_logger` (#193)
* **New Resource:** `pingdirectory_ldap_correlation_attribute_pair` (#194)
* **New Resource:** `pingdirectory_local_db_composite_index` (#196)
* **New Resource:** `pingdirectory_local_db_vlv_index` (#197)
* **New Resource:** `pingdirectory_json_field_constraints` (#198)
* **New Resource:** `pingdirectory_oauth_token_handler` (#199)
* **New Resource:** `pingdirectory_password_policy` (#200)
* **New Resource:** `pingdirectory_synchronization_provider` (#201)
* **New Resource:** `pingdirectory_scim_subattribute` (#202)
* **New Resource:** `pingdirectory_plugin_root` (#203)
* **New Resource:** `pingdirectory_result_criteria` (#204)

# v0.5.0 April 28 2023

### FEATURES
* **New Resource:** `pingdirectory_certificate_mapper` (#142)
* **New Resource:** `pingdirectory_sasl_mechanism_handler` (#146)
* **New Resource:** `pingdirectory_monitoring_endpoint` (#147)
* **New Resource:** `pingdirectory_recurring_task_chain` (#148)
* **New Resource:** `pingdirectory_scim_schema` (#149)
* **New Resource:** `pingdirectory_scim_resource_type` (#150)
* **New Resource:** `pingdirectory_scim_attribute` (#151)
* **New Resource:** `pingdirectory_scim_attribute_mapping` (#153)
* **New Resource:** `pingdirectory_id_token_validator` (#154)
* **New Resource:** `pingdirectory_web_application_extension` (#154)
* **New Resource:** `pingdirectory_entry_cache` (#158)
* **New Resource:** `pingdirectory_gauge_data_source` (#159)

### DOCUMENTATION UPDATES
* `Move default_ resource documentation pages to correct subcategories` (#144)

### DEPENDENCIES
* `Bump github.com/pingidentity/pingdirectory-go-client from v9200.0.0 to v9200.5.0` (#153)
* `Bump github.com/bflad/tfproviderlint from 0.28.1 to 0.29.0` (#140)
* `Bump github.com/terraform-linters/tflint from 0.46.0 to 0.46.1` (#152)

# v0.4.0 April 14 2023

### ENHANCEMENTS
* `Support PingDirectory versions 9.1.0.1 and 9.1.0.2` (#136)
* `Add config validators to resources to handle constraints between attributes` (#132)
### DOCUMENTATION UPDATES
* `Add provider version requirement to HCL examples` (#119)
* `Update registry resource documentation hierarchy` (#135)
### DEPENDENCIES
* `Use terraform-plugin-testing for acceptance tests rather than sdkv2` (#128)
* `Bump github.com/golangci/golangci-lint from 1.52.0 to 1.52.2` (#125)
* `Bump github.com/hashicorp/terraform-plugin-framework from 1.1.1 to 1.2.0` (#124)
* `Bump github.com/hashicorp/terraform-plugin-go from 0.14.3 to 0.15.0` (#123)
* `Bump github.com/hashicorp/terraform-plugin-sdk/v2 from 2.25.0 to 2.26.1` (#122)
* `Bump actions/setup-go from 3 to 4` (#121)
* `Bump github.com/terraform-linters/tflint from 0.45.0 to 0.46.0` (#134)
