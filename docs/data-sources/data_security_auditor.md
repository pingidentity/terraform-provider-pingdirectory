---
page_title: "pingdirectory_data_security_auditor Data Source - terraform-provider-pingdirectory"
subcategory: "Data Security Auditor"
description: |-
  Describes a Data Security Auditor.
---

# pingdirectory_data_security_auditor (Data Source)

Describes a Data Security Auditor.

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

data "pingdirectory_data_security_auditor" "myDataSecurityAuditor" {
  id = "MyDataSecurityAuditor"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) Name of this object.

### Read-Only

- `account_expiration_warning_interval` (String) If set, the auditor will report all users with account expiration times are in the future, but are within the specified length of time away from the current time.
- `audit_backend` (Set of String) Specifies which backends the data security auditor may be applied to. By default, the data security auditors will audit entries in all backend types that support data auditing (Local DB, LDIF, and Config File Handler).
- `audit_severity` (String) Specifies the severity of events to include in the report.
- `enabled` (Boolean) Indicates whether the Data Security Auditor is enabled for use.
- `extension_argument` (Set of String) The set of arguments used to customize the behavior for the Third Party Data Security Auditor. Each configuration property should be given in the form 'name=value'.
- `extension_class` (String) The fully-qualified name of the Java class providing the logic for the Third Party Data Security Auditor.
- `filter` (Set of String) The filter to use to identify entries that should be reported. Multiple filters may be configured, and each reported entry will indicate which of these filter(s) matched that entry.
- `idle_account_error_interval` (String) The length of time to use as the error interval for idle accounts. If the length of time since a user last authenticated is greater than the error interval, then an error will be generated for that account. If no error interval is defined, then only the warning interval will be used.
- `idle_account_warning_interval` (String) The length of time to use as the warning interval for idle accounts. If the length of time since a user last authenticated is greater than the warning interval but less than the error interval (or if it is greater than the warning interval and no error interval is defined), then a warning will be generated for that account.
- `include_attribute` (Set of String) Specifies the attributes from the audited entries that should be included detailed reports. By default, no attributes are included.
- `include_privilege` (Set of String) If defined, only entries with the specified privileges will be reported. By default, entries with any privilege assigned will be reported.
- `maximum_idle_time` (String) If set, users that have not authenticated for more than the specified time will be reported even if idle account lockout is not configured. Note that users may only be reported if the last login time tracking is enabled.
- `never_logged_in_account_error_interval` (String) The length of time to use as the error interval for accounts that do not appear to have authenticated. If this is not specified, then the never-logged-in warning interval will be used. The idle account warning and error intervals will be used if no never-logged-in interval is configured.
- `never_logged_in_account_warning_interval` (String) The length of time to use as the warning interval for accounts that do not appear to have authenticated. If this is not specified, then the idle account warning interval will be used.
- `password_evaluation_age` (String) If set, the auditor will report all users with passwords older than the specified value even if password expiration is not enabled.
- `report_file` (String) Specifies the name of the detailed report file.
- `type` (String) The type of Data Security Auditor resource. Options are ['expired-password', 'idle-account', 'disabled-account', 'weakly-encoded-password', 'privilege', 'account-usability-issues', 'locked-account', 'filter', 'account-validity-window', 'multiple-password', 'deprecated-password-storage-scheme', 'nonexistent-password-policy', 'access-control', 'third-party']
- `weak_crypt_encoding` (Set of String) Reporting on users with passwords encoded using the Crypt Password Storage scheme may be further limited by selecting one or more encoding mechanisms that are considered weak.
- `weak_password_storage_scheme` (Set of String) The password storage schemes that are considered weak. Users with any of the specified password storage schemes will be included in the report.
