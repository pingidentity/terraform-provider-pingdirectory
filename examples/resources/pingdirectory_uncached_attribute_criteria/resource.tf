resource "pingdirectory_uncached_attribute_criteria" "myUncachedAttributeCriteria" {
  name        = "MyUncachedAttributeCriteria"
  type        = "default"
  description = "My default uncached attribute criteria"
  enabled     = false
}
