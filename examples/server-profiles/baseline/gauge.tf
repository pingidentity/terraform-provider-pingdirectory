resource "pingdirectory_default_gauge" "defaultCpuUsageGauge" {
  type    = "numeric"
  id      = "CPU Usage (Percent)"
  enabled = false
}

resource "pingdirectory_default_numeric_gauge" "defaultLicenseExpirationGauge" {
  type    = "numeric"
  id      = "License Expiration (Days)"
  enabled = false
}

resource "pingdirectory_default_numeric_gauge" "defaultAvailableFileDescriptorsGauge" {
  type    = "numeric"
  id      = "Available File Descriptors"
  enabled = false
}

resource "pingdirectory_default_file_based_audit_log_publisher" "defaultDataRecoveryLog" {
  type    = "numeric"
  id      = "Data Recovery Log"
  enabled = false
}
