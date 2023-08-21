resource "pingdirectory_otp_delivery_mechanism" "myOtpDeliveryMechanism" {
  name           = "MyOtpDeliveryMechanism"
  type           = "email"
  sender_address = "sender@example.com"
  enabled        = true
}
