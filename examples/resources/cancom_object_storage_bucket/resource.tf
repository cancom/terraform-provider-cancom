resource "cancom_object_storage_bucket" "bucket" {
  bucket_name        = "test-bucket"
  availability_class = "multiDc"
}
