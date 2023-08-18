data "pingdirectory_prometheus_monitor_attribute_metric" "myPrometheusMonitorAttributeMetric" {
  http_servlet_extension_name = "MyHttpServletExtension"
  metric_name                 = "myMetric"
}
