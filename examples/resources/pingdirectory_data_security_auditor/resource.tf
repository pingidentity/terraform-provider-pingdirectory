resource "pingdirectory_data_security_auditor" "myDataSecurityAuditor" {
  name = "MyDataSecurityAuditor"
  type = "expired-password"
}
