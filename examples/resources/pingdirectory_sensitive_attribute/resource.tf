resource "pingdirectory_sensitive_attribute" "mySensitiveAttribute" {
  name           = "MySensitiveAttribute"
  attribute_type = ["userPassword", "pwdHistory"]
}
