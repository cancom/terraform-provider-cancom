
resource "cancom_ipam_supernet" "<name of this resource>" {
  instance_id = "<id of the managed instance this supernet shall belong to>"
  name_tag = "<name>"
  supernet_cidr = "<the prefix in cidr notation>"
  description = "<a description>"
}

# example with id from instance within same stack
resource "cancom_ipam_supernet" "tfsupernet01" {
  instance_id = cancom_ipam_instance.tfcreated01.id
  name_tag = "supernet-from-tf"
  supernet_cidr = "10.100.0.0/16"
  description = "supernet created by terraform"
}

