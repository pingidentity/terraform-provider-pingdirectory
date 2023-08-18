resource "pingdirectory_trusted_certificate" "myTrustedCertificate" {
  name        = "MyTrustedCertificate"
  certificate = file("${path.module}/cert.pem")
}
