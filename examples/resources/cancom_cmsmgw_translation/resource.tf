resource "cancom_cmsmgw_translation" "<name of resource>" {
  mgw_id      = "<gateway-Id>"
  name_tag    = "<name>"
  customer_ip = "10.0.0.10"
  dns_zone    = "int.cc-mase.com" // or "" // others are possible later
}

resource "cancom_cmsmgw_translation" "translation_tf_01" {
  mgw_id      = cancom_cmsmgw_gateway.gw-tf-03.id
  name_tag    = "testtranslation01"
  customer_ip = "10.0.0.10"
  dns_zone    = ""
}
