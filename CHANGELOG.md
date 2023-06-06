# v0.7.0 (Unreleased)

### FEATURES
* **New Resource:** `pingdirectory_root_dse_backend ` (#206)
* **New Resource:** `pingdirectory_search_entry_criteria ` (#207)
* **New Resource:** `pingdirectory_search_reference_criteria ` (#208)
* **New Resource:** `pingdirectory_server_group` (#209)
* **New Resource:** `pingdirectory_soft_delete_policy ` (#210)
* **New Resource:** `pingdirectory_token_claim_validation` (#211)
* **New Resource:** `pingdirectory_trusted_ccertificate` (#217)


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
