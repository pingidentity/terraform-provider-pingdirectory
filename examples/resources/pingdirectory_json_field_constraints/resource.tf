resource "pingdirectory_json_attribute_constraints" "myJsonAttributeConstraints" {
  attribute_type       = "ubidEntitlement"
  description          = "ubidEntitlement attribute constraint"
  allow_unnamed_fields = false
}

resource "pingdirectory_json_field_constraints" "myJsonFieldConstraints" {
  json_attribute_constraints_name = pingdirectory_json_attribute_constraints.myJsonAttributeConstraints.attribute_type
  json_field                      = "id"
  value_type                      = "string"
}