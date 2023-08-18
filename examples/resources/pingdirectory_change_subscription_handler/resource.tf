resource "pingdirectory_change_subscription_handler" "myChangeSubscriptionHandler" {
  name         = "MyChangeSubscriptionHandler"
  type         = "groovy-scripted"
  script_class = "com.example.myscriptclass"
  enabled      = false
}
