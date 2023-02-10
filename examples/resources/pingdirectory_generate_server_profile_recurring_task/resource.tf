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

resource "pingdirectory_generate_server_profile_recurring_task" "myGenerateServerProfileRecurringTask" {
  id                            = "MyGenerateServerProfileRecurringTask"
  profile_directory             = "/opt/out/instance"
  retain_previous_profile_count = 10
}