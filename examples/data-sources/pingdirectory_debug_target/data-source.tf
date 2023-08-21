data "pingdirectory_debug_target" "myDebugTarget" {
  log_publisher_name = "MyLogPublisher"
  debug_scope        = "com.Example"
}
