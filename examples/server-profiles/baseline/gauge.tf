resource "pingdirectory_default_gauge" "defaultCpuUsageGauge" {
  name    = "CPU Usage (Percent)"
  enabled = false
}

resource "pingdirectory_default_gauge" "defaultLicenseExpirationGauge" {
  name    = "License Expiration (Days)"
  enabled = false
}

resource "pingdirectory_default_gauge" "defaultAvailableFileDescriptorsGauge" {
  name    = "Available File Descriptors"
  enabled = false
}

resource "pingdirectory_default_log_publisher" "defaultDataRecoveryLog" {
  name    = "Data Recovery Log"
  enabled = false
}
