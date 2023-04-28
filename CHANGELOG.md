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
