# Disable the default failed operations access logger
resource "pingdirectory_default_log_publisher" "defaultFileBasedAccessLogPublisher" {
  name    = "Failed Operations Access Logger"
  enabled = false
}

# Create a new custom file based access logger
resource "pingdirectory_log_publisher" "myNewFileBasedAccessLogPublisher" {
  type                 = "file-based-access"
  name                 = "MyNewFileBasedAccessLogPublisher"
  log_file             = "logs/example.log"
  log_file_permissions = "600"
  rotation_policy      = ["Size Limit Rotation Policy"]
  retention_policy     = ["File Count Retention Policy"]
  asynchronous         = true
  enabled              = false
}

# Enable the default JMX connection handler
resource "pingdirectory_default_connection_handler" "defaultJmxConnHandler" {
  name    = "JMX Connection Handler"
  enabled = true
}

# Create a new custom JMX connection handler
resource "pingdirectory_connection_handler" "myJmxConnHandler" {
  type        = "jmx"
  name        = "MyJmxConnHandler"
  enabled     = false
  listen_port = 8888
}
