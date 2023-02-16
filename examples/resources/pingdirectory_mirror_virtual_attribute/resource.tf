terraform {
  required_providers {
    pingdirectory = {
      source = "pingidentity/pingdirectory"
    }
  }
}

provider "pingdirectory" {
  username   = "cn=administrator"
  password   = "2FederateM0re"
  https_host = "https://localhost:1443"
}

resource "pingdirectory_mirror_virtual_attribute" "myMirrorVirtualAttribute" {
  id               = "MyMirrorVirtualAttribute"
  source_attribute = "mail"
  enabled          = true
  attribute_type   = "name"
}
