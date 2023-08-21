resource "pingdirectory_prometheus_monitor_attribute_metric" "myPrometheusMonitorAttributeMetric" {
  http_servlet_extension_name = "Prometheus Monitoring"
  metric_name                 = "mymetric"
  monitor_attribute_name      = "max-queue-size"
  monitor_object_class_name   = "ds-unboundid-work-queue-monitor-entry"
  metric_type                 = "numeric"
}
