terraform {
  required_providers {
    pingdirectory = {
      source = "pingidentity/pingdirectory"
    }
  }
}

provider "pingdirectory" {
  username = "cn=administrator"
  password = "2FederateM0re"
  https_host = "https://localhost:1443"
}

// This set is approximately the minimum set required for you to be able to run
// 'dsconfig get-root-dn-prop' successfully.  If you remove any of these permissions, 
// you risk loss of access to the RootDN permission object.
resource "pingdirectory_root_dn" "myrootdn" {
  default_root_privilege_name = ["bypass-acl", "config-read", "config-write", "modify-acl", "privilege-change", "use-admin-session"]
}
