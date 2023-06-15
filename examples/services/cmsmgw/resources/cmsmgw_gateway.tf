

resource "cancom_cmsmgw_gateway" {
	name = "Test"
	size = "small"

	nat             = true
	bastion_linux   = false
	bastion_windows = false
}
