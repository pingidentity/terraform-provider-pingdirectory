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

resource "pingdirectory_file_based_access_log_publisher" "myFileBasedAccessLogPublisher" {
  id                   = "MyFileBasedAccessLogPublisher"
  log_file             = "logs/example.log"
  log_file_permissions = "600"
  rotation_policy      = ["Size Limit Rotation Policy"]
  retention_policy     = ["File Count Retention Policy"]
  asynchronous         = true
  enabled              = false
}
