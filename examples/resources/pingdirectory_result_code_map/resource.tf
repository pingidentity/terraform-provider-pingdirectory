resource "pingdirectory_result_code_map" "myResultCodeMap" {
  name                          = "MyResultCodeMap"
  description                   = "mapping my codes"
  bind_missing_user_result_code = 59
  server_error_result_code      = 81
}
