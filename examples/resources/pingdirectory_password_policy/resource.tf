resource "pingdirectory_password_policy" "myPasswordPolicy" {
  name                            = "MyPasswordPolicy"
  password_attribute              = "userPassword"
  default_password_storage_scheme = ["Blowfish", "Salted SHA-384", "Salted SHA-512"]
}
