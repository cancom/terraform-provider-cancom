resource "cancom_s3_bucket" "bucket" {
  name               = "test-bucket"
  availability_class = "multiDc"
}
