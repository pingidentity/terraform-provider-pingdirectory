resource "pingdirectory_password_storage_scheme" "myPasswordStorageScheme" {
  name                     = "MyPasswordStorageScheme"
  type                     = "argon2d"
  enabled                  = false
  iteration_count          = 10
  parallelism_factor       = 1
  memory_usage_kb          = 16
  salt_length_bytes        = 16
  derived_key_length_bytes = 16
}
