resource "pingdirectory_entry_cache" "myEntryCache" {
  name        = "MyEntryCache"
  enabled     = true
  cache_level = 1
}
