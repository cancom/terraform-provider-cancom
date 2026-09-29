resource "cancom_iam_service_user_session" "session" {
  service_user = cancom_iam_service_user.su.principal
  comment      = "CI/CD ephemeral runner session"
  reroll_days  = 7
}
