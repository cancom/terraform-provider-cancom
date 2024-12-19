resource "cancom_cmsmgw_connection" "<name of resource>" {
  mgw_id                   = "<id of the gateway>"
  customer_primary_gw_ip   = "<ip of the customer gateway"
  name_tag                 = "<name>"
  customer_secondary_gw_ip = "" // not used currently
  connection_profile       = "AZURE, CANCOOM-V01-A, CANCOM-V01-B, CANCOM-V01-C"
  customer_networks        = ["cidr-1", "cidr-2", "cidr-n"]
  cancom_networks          = ["100.97.0.0/16"]
  ipsec_psk_a              = "123456abcdefg123456" // provide a secret or empty to let the service create one
  ipsec_psk_b              = ""                    // not used currently
}

resource "cancom_cmsmgw_connection" "conn_tf-01" {
  mgw_id                   = cancom_cmsmgw_gateway.gw-tf-03.id
  customer_primary_gw_ip   = "192.168.3.1"
  name_tag                 = "testconnection01"
  customer_secondary_gw_ip = ""
  connection_profile       = "AZURE"
  customer_networks        = ["10.0.0.0/16", "10.10.0.0/16"]
  cancom_networks          = ["100.97.0.0/16"]
  ipsec_psk_a              = "123456abcdefg123456"
  ipsec_psk_b              = ""
}
