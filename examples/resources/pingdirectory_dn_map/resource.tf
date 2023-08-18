resource "pingdirectory_dn_map" "myDnMap" {
  name            = "MyDnMap"
  from_dn_pattern = "*,**,dc=com"
  to_dn_pattern   = "uid={givenname:/^(.)(.*)/$1/s}{sn:/^(.)(.*)/$1/s}{eid},{2},o=example"
}
