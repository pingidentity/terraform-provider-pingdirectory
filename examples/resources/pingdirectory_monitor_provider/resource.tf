# resource "pingdirectory_default_monitor_provider" "myMonitorProvider" {
#   name        = "MyMonitorProvider"
#   type        = "general"
#   description = "My general monitor entry resource provider"
#   enabled     = false
# }

resource "pingdirectory_default_monitor_provider" "mine" {
  name        = "General Monitor Entry"
  description = "My general monitor entry resource provider"
  enabled     = true
}
