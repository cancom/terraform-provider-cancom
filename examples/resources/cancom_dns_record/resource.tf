resource "cancom_dns_record" "record" {
  zone_name = "example.com"
  name      = "test.example.com"
  type      = "A"
  content   = "127.0.0.1"
}
