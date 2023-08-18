resource "pingdirectory_gauge" "myGauge" {
  name              = "MyGauge"
  type              = "indicator"
  gauge_data_source = "Replication Connection Status"
  enabled           = false
}
