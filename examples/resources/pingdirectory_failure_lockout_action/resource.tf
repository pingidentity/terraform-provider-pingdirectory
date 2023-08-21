resource "pingdirectory_failure_lockout_action" "myFailureLockoutAction" {
  name  = "MyFailureLockoutAction"
  type  = "delay-bind-response"
  delay = "1 s"
}
