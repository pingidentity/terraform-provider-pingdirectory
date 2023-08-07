resource "pingdirectory_default_connection_handler" "defaultHttpsConnectionHandler" {
  name                      = "HTTPS Connection Handler"
  web_application_extension = []
  http_servlet_extension    = ["Delegated Admin", "Available or Degraded State", "Available State", "Configuration", "Consent", "Directory REST API", "Instance Root File Servlet", "SCIM2", "Swagger UI"]
}
