resource "pingdirectory_http_servlet_cross_origin_policy" "myHttpServletCrossOriginPolicy" {
  name                 = "MyHttpServletCrossOriginPolicy"
  cors_allowed_headers = ["Accept, Access-Control-Request-Headers"]
}
