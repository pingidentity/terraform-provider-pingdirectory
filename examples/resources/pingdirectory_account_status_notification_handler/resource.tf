resource "pingdirectory_account_status_notification_handler" "myAccountStatusNotificationHandler" {
  name                                  = "MyAccountStatusNotificationHandler"
  type                                  = "smtp"
  send_message_without_end_user_address = false
  recipient_address                     = ["test@example.com", "users@example.com"]
  sender_address                        = "unboundid-notifications@example.com"
  message_subject                       = ["account-disabled:Your directory account has been disabled"]
  message_template_file                 = ["account-disabled:config/messages/account-disabled.template"]
  enabled                               = false
}
