resource "pingdirectory_default_mac_secret_key" "myMacSecretKey" {
  name                 = "MyKeyId"
  server_instance_name = "MyServerInstance"
  key_id               = "MyKeyId"
}
