resource "pingdirectory_default_cipher_secret_key" "myCipherSecretKey" {
  name                 = "MyKeyId"
  server_instance_name = "MyServerInstance"
  key_id               = "MyKeyId"
}
