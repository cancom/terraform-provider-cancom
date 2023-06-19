resource "cancom_ipam_instance" "tfcreated01" {
   name_tag ="created-by-terraform"
   description = "Intance created by teraform"
   managed_by = "CUSTOMER"
   release_wait_time = "350"
}