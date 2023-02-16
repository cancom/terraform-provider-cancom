---
page_title: "CANCOM: dns_record"
description: |-
  Defines a DNS record
---

# cancom_dns_record (Resource)

## Example Usage

```terraform
resource "cancom_dns_record" {
  zone_name = "example.com"
  name      = "test.example.com"
  type      = "A"
  content   = "127.0.0.1"
}
```


<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `content` (String) Target of the record
- `name` (String) Name of the record
- `ttl` (Number) TTL, if not set, defaults to the zones TTL
- `type` (String) Type of the record (i.e. A, CNAME, ...)
- `zone_name` (String) Zone that this record belongs to

### Read-Only

- `comments` (String) You can optionally add comments to records to describe their intended usage
- `id` (String) The uuid of the record
- `last_change_date` (String) The date at which the record was last updated
- `zone_id` (String) The uuid of the zone