resource "aws_secretsmanager_secret" "octopus_credentials" {
  name = "dover/octopus-credentials"
}

resource "aws_secretsmanager_secret_version" "octopus_credentials" {
  secret_id = aws_secretsmanager_secret.octopus_credentials.id
  secret_string = jsonencode({
    email    = var.octopus_email
    password = var.octopus_password
  })
}
