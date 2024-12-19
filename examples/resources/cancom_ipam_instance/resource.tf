resource "cancom_ipam_instance" "tfcreated01" {
  name_tag          = "created-by-terraform"
  description       = "Instance created by terraform"
  managed_by        = "CUSTOMER"
  release_wait_time = "350"
}
