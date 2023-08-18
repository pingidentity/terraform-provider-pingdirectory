resource "pingdirectory_log_publisher" "myLogPublisher" {
  name                   = "MyLogPublisher"
  type                   = "syslog-json-audit"
  syslog_external_server = "example.com"
  enabled                = false
}
