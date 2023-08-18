resource "pingdirectory_notification_manager" "myNotificationManager" {
  name                 = "MyNotificationManager"
  extension_class      = "com.example.MyClass"
  enabled              = true
  subscription_base_dn = "ou=subscriptionbase,dc=example,dc=com"
}
