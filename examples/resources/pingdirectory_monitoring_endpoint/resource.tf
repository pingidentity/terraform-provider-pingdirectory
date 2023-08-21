resource "pingdirectory_monitoring_endpoint" "myMonitoringEndpoint" {
  name     = "MyMonitoringEndpoint"
  hostname = "localhost"
  enabled  = false
}
