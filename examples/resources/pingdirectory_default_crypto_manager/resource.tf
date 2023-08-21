resource "pingdirectory_default_crypto_manager" "myCryptoManager" {
  mac_key_length    = 256
  cipher_key_length = 256
  ssl_cert_nickname = "ssl-certificate-alias"
}
