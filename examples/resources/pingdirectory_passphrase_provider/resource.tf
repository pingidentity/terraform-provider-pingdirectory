resource "pingdirectory_passphrase_provider" "myPassphraseProvider" {
  name                 = "MyPassphraseProvider"
  type                 = "environment-variable"
  environment_variable = "PASSPHRASE_ENV_VARIABLE"
  enabled              = true
}
