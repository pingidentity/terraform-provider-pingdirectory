resource "pingdirectory_json_attribute_constraints" "myJsonAttributeConstraints" {
  attribute_type       = "ubidEntitlement"
  description          = "ubidEntitlement attribute constraint"
  allow_unnamed_fields = false
}
