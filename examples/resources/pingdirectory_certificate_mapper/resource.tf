resource "pingdirectory_certificate_mapper" "myCertificateMapper" {
  name    = "MyCertificateMapper"
  type    = "subject-equals-dn"
  enabled = false
}
