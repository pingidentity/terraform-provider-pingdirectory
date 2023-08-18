resource "pingdirectory_default_global_configuration" "myGlobalConfiguration" {
  location              = "Docker"
  encrypt_data          = true
  sensitive_attribute   = ["Delivered One-Time Password", "TOTP Shared Secret"]
  tracked_application   = ["Requests by Root Users"]
  result_code_map       = "Sun DS Compatible Behavior"
  disabled_privilege    = ["jmx-write", "jmx-read"]
  maximum_shutdown_time = "4 m"
}
