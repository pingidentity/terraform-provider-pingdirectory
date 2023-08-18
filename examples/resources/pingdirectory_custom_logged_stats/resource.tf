resource "pingdirectory_custom_logged_stats" "myCustomLoggedStats" {
  name                = "MyCustomLoggedStats"
  plugin_name         = "JSON Stats Logger"
  monitor_objectclass = "ds-memory-usage-monitor-entry"
  attribute_to_log    = ["total-bytes-used-by-memory-consumers"]
  statistic_type      = ["raw"]
}
