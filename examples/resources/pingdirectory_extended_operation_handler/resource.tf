resource "pingdirectory_extended_operation_handler" "myExtendedOperationHandler" {
  name    = "MyExtendedOperationHandler"
  type    = "cancel"
  enabled = false
}
