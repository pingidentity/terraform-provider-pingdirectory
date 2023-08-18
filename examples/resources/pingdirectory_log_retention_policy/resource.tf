resource "pingdirectory_log_retention_policy" "myLogRetentionPolicy" {
  name            = "MyLogRetentionPolicy"
  type            = "time-limit"
  description     = "Time limit for log retention"
  retain_duration = "1 w"
}
