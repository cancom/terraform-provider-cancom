resource "cancom_object_storage_bucket" "bucket" {
  bucket_name        = "test-bucket"
  availability_class = "multiDc"
  description        = "Test bucket"
  ip_whitelist       = ["10.0.0.0/16"]
}
