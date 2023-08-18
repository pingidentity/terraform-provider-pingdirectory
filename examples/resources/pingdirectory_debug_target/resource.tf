resource "pingdirectory_debug_target" "myDebugTarget" {
  log_publisher_name = "File-Based Debug Logger"
  debug_scope        = "com.example.MyClass"
  debug_level        = "all"
}
