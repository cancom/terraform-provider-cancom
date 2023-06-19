resource "cancom_ipam_host" "<name of resource>" {
	network_crn = "crn:KmBRNWt87RV26jJfZMeJus::ip-management:network:F4tfeFwx2i7HFjNrFop49z"
   	name_tag = "<name for tis assignment>"
   	qualifier = "<qualifier, optional, see documentation>"
   	description = "<description for this host-assignment>"
}


# the following examples creates two assignments using resources created within the same stack
resource "cancom_ipam_host" "host-assignment-01" {
   network_crn = cancom_ipam_network.tfnetwork01.id
   name_tag = "host-assignment-01"
   description = "assignment creaed by terraform"
}

resource "cancom_ipam_host" "host-assignment-02" {
   network_crn = cancom_ipam_network.tfnetwork01.id
   name_tag = "host-assignment-02"
   description = "assignment creaed by terraform"
}