resource "cancom_cmsmgw_gateway" "<resource-name>" {
  customer        = "<ACCT Number>"
  name            = "<name>"
  size            = "<small, medium, large"
  nat_translation = false
  description     = "<a description for this gateway>"
}
