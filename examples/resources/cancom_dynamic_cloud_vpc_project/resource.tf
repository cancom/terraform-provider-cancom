resource "cancom_dynamic_cloud_vpc_project" "basic_usage_example" {
  name            = "test-cancom-terraform-provider"
  project_comment = "Test VPC Project created with the CANCOM Terraform provider"
  users           = ["crn:cancom::iam:user:example.name@example.domain"]
}
