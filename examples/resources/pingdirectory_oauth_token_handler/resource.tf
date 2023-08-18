resource "pingdirectory_oauth_token_handler" "myOauthTokenHandler" {
  name         = "MyOauthTokenHandler"
  type         = "groovy-scripted"
  description  = "My groovy scripted OAuth token handler"
  script_class = "com.example"
}
