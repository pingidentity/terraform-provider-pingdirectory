resource "pingdirectory_web_application_extension" "myWebApplicationExtension" {
  name                    = "MyWebApplicationExtension"
  type                    = "generic"
  base_context_path       = "/myexamplepath"
  document_root_directory = "/my/directory/path"
}
