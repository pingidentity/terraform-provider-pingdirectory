---
page_title: "pingdirectory_alarm_manager Data Source - terraform-provider-pingdirectory"
subcategory: "Alarm Manager"
description: |-
  Describes a Alarm Manager.
---

# pingdirectory_alarm_manager (Data Source)

Describes a Alarm Manager.

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

data "pingdirectory_alarm_manager" "myAlarmManager" {
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `default_gauge_alert_level` (String) Specifies the level at which alerts are sent for alarms raised by the Alarm Manager.
- `generated_alert_types` (Set of String) Indicates what kind of alert types should be generated.
- `id` (String) Name of this object.
- `suppressed_alarm` (Set of String) Specifies the names of the alarm alert types that should be suppressed. If the condition that triggers an alarm in this list occurs, then the alarm will not be raised and no alerts will be generated. Only a subset of alarms can be suppressed in this way. Alarms triggered by a gauge can be disabled by disabling the gauge.
