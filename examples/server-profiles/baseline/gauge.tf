resource "pingdirectory_default_gauge" "defaultCpuUsageGauge" {
  type    = "numeric"
  name    = "CPU Usage (Percent)"
  enabled = false
}

resource "pingdirectory_default_gauge" "defaultLicenseExpirationGauge" {
  type    = "numeric"
  name    = "License Expiration (Days)"
  enabled = false
}

resource "pingdirectory_default_gauge" "defaultAvailableFileDescriptorsGauge" {
  type    = "numeric"
  name    = "Available File Descriptors"
  enabled = false
}

resource "pingdirectory_default_log_publisher" "defaultDataRecoveryLog" {
  type    = "file-based-audit"
  name    = "Data Recovery Log"
  enabled = false
}
