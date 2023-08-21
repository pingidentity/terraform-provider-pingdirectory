resource "pingdirectory_log_rotation_policy" "myLogRotationPolicy" {
  name              = "MyLogRotationPolicy"
  type              = "time-limit"
  description       = "Time limit before rotating logs"
  rotation_interval = "2 w"
}
