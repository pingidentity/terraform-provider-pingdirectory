resource "pingdirectory_key_manager_provider" "myKeyManagerProvider" {
  name           = "MyKeyManagerProvider"
  type           = "file-based"
  enabled        = false
  key_store_file = "/tmp/key-store-file"
}
