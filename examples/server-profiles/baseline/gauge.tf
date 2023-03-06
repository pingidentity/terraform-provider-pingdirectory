resource "pingdirectory_default_numeric_gauge" "defaultCpuUsageGauge" {
  id      = "CPU Usage (Percent)"
  enabled = false
}

resource "pingdirectory_default_numeric_gauge" "defaultLicenseExpirationGauge" {
  id      = "License Expiration (Days)"
  enabled = false
}

resource "pingdirectory_default_numeric_gauge" "defaultAvailableFileDescriptorsGauge" {
  id      = "Available File Descriptors"
  enabled = false
}

resource "pingdirectory_default_file_based_audit_log_publisher" "defaultDataRecoveryLog" {
  id      = "Data Recovery Log"
  enabled = false
}
