resource "cancom_ipam_network" "<name of resource>" {
  supernet_id = "<id of the supernet this network shall be assigned from>"
  name_tag    = "<name>"
  request     = "<size of the network as netmask>"
  host_assign = "<true or false, decides if ip-assignments are possible from the created network>"
  description = "<a description>"
}

# example with supernet created within the same stack
resource "cancom_ipam_network" "tfnetwork01" {
  supernet_id = cancom_ipam_supernet.tfsupernet01.id
  name_tag    = "network-from-tf"
  request     = "/25"
  host_assign = true
  description = "erstellt über terraform"
}
