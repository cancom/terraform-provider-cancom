resource "cancom_s3_user" "user" {
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
