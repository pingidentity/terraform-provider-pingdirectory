resource "pingdirectory_default_alarm_manager" "myAlarmManager" {
  default_gauge_alert_level = "critical-only"
  generated_alert_types     = ["standard"]
}
