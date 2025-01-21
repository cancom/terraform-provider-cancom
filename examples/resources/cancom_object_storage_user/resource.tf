resource "cancom_object_storage_user" "user" {
  username     = "svc-testuser"
  description  = "This is a test user"
  permissions  = jsonencode({
    Statement = [
        {
            Effect = "Allow"
            Action = "*"
            Resource = "*"
        }
    ]
  })
}
