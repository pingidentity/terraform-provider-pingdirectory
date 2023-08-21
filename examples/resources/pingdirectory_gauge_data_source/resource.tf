resource "pingdirectory_gauge_data_source" "myGaugeDataSource" {
  name                = "MyGaugeDataSource"
  type                = "indicator"
  monitor_objectclass = "ds-host-system-disk-monitor-entry"
  monitor_attribute   = "pct-busy"
}
