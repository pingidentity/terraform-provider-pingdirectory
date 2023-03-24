terraform {
  required_version = ">=1.1"
  required_providers {
    pingdirectory = {
      version = "~> 0.3.0"
      source  = "pingidentity/pingdirectory"
    }
  }
}

provider "pingdirectory" {
  username   = "cn=administrator"
  password   = "2FederateM0re"
  https_host = "https://localhost:1443"
  # Warning: The insecure_trust_all_tls attribute configures the provider to trust any certificate presented by the PingDirectory server.
  # It should not be used in production. If you need to specify trusted CA certificates, use the
  # ca_certificate_pem_files attribute to point to any number of trusted CA certificate files
  # in PEM format. If you do not specify certificates, the host's default root CA set will be used.
  # Example:
  # ca_certificate_pem_files = ["/example/path/to/cacert1.pem", "/example/path/to/cacert2.pem"]
  insecure_trust_all_tls = true
  product_version        = "9.2.0.0"
}

# Use "pingdirectory_default_smtp_account_status_notification_handler" if you are adopting existing configuration from the PingDirectory server into Terraform
resource "pingdirectory_smtp_account_status_notification_handler" "mySmtpAccountStatusNotificationHandler" {
  id                                    = "MySmtpAccountStatusNotificationHandler"
  send_message_without_end_user_address = false
  recipient_address                     = ["test@example.com", "users@example.com"]
  sender_address                        = "unboundid-notifications@example.com"
  message_subject                       = ["account-disabled:Your directory account has been disabled"]
  message_template_file                 = ["account-disabled:config/messages/account-disabled.template"]
  enabled                               = false
}
