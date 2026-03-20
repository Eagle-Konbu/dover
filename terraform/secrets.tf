resource "aws_secretsmanager_secret" "octopus_credentials" {
  name = "dover/octopus-credentials"
}
