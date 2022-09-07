resource "cancom_iam_user" "user" {
  name = "test"
}

resource "cancom_iam_policy" "test_policy_azure-app-management" {
  service   = "domdns"
  principal = cancom_iam_user.user.id
  policy {
    profile = "full-access"
  }
}
