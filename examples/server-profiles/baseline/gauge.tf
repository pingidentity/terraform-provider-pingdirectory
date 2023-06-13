resource "pingdirectory_default_gauge" "defaultCpuUsageGauge" {
  type    = "numeric"
  id      = "CPU Usage (Percent)"
  enabled = false
}

resource "pingdirectory_default_gauge" "defaultLicenseExpirationGauge" {
  type    = "numeric"
  id      = "License Expiration (Days)"
  enabled = false
}

resource "pingdirectory_default_gauge" "defaultAvailableFileDescriptorsGauge" {
  type    = "numeric"
  id      = "Available File Descriptors"
  enabled = false
}

resource "pingdirectory_default_log_publisher" "defaultDataRecoveryLog" {
  type = "file-based-audit"
  id      = "Data Recovery Log"
  enabled = false
}
