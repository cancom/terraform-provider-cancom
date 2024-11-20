resource "cancom_dynamic_cloud_vpc_project" "basic_usage_example_users" {
  name            = "test-cancom-terraform-provider-project-users"
  project_comment = "Test VPC Project users managed with the CANCOM Terraform provider"
}

resource "cancom_dynamic_cloud_vpc_project_users" "basic_usage_example_users" {
  vpc_project_id = cancom_dynamic_cloud_vpc_project.basic_usage_example.id
  users          = ["crn:cancom::iam:user:example.name@example.domain"]
}
