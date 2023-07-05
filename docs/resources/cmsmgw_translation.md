---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "CANCOM: cmsmgw_translation Resource"
subcategory: ""
description: |-
  Represents a translation resource
---

# cancom_cmsmgw_gateway (Resource)

## Example Usage

```terraform
resource "cancom_cmsmgw_translation" "<name of resource>" {
  mgw_id = <gateway-Id>
  name_tag = "<name>"
  customer_ip = "10.0.0.10"
  dns_zone = "" or "int.cc-mase.com"  //others are possible later
}

resource "cancom_cmsmgw_translation" "translation_tf_01" {
  mgw_id = cancom_cmsmgw_gateway.gw-tf-03.id
  name_tag = "testtranslation01"
  customer_ip = "10.0.0.10"
  dns_zone = ""
}
```


<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `customer_ip` (String)
- `mgw_id` (String)

### Optional

- `dns_zone` (String)
- `name_tag` (String)

### Read-Only

- `id` (String) The ID of this resource.
- `spark_ip` (String)