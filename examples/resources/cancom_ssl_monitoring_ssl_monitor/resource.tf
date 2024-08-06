resource "cancom_ssl_monitoring_ssl_monitor" "test_ssl" {
  domain_name          = "example.com"
  comment              = "This is a test"
  minimum_grade        = "A"
  contact_email_cancom = "test@example.de"
  is_managed_by_cancom = true
  protocol             = "http"
  port                 = 443
}
