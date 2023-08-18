resource "pingdirectory_conjur_authentication_method" "myConjurAuthenticationMethod" {
  name     = "MyConjurAuthenticationMethod"
  username = "myusername"
  password = "mypassword"
}
