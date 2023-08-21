resource "pingdirectory_cipher_stream_provider" "myCipherStreamProvider" {
  name                   = "MyCipherStreamProvider"
  type                   = "amazon-key-management-service"
  kms_encryption_key_arn = "my_example_Encryption_Key"
  enabled                = false
}
