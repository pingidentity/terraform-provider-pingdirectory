resource "pingdirectory_log_file_rotation_listener" "myLogFileRotationListener" {
  name             = "MyLogFileRotationListener"
  type             = "summarize"
  description      = "My summarize log file rotation listener"
  enabled          = true
  output_directory = "/tmp"
}
