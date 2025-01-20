resource "cancom_object_storage_bucket" "bucket" {
  name               = "test-bucket"
  availability_class = "multiDc"
}
