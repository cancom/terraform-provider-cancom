
resource "cancom_ipam_host" "host" {
	network_crn = "crn:KmBRNWt87RV26jJfZMeJus::ip-management:network:F4tfeFwx2i7HFjNrFop49z"
   	name_tag = "hostname.cancom.de"
   	qualifier = "<qualifier, optional>"
   	description = "<description for this host-assignment>"
}

