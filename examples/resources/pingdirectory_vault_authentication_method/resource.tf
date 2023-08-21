resource "pingdirectory_vault_authentication_method" "myVaultAuthenticationMethod" {
  name               = "MyVaultAuthenticationMethod"
  type               = "static-token"
  vault_access_token = "myExampleToken"
}
