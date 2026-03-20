locals {
  lambda_zip_path = "${path.module}/../lambda.zip"
}

resource "aws_lambda_function" "dover" {
  function_name    = "dover"
  filename         = local.lambda_zip_path
  source_code_hash = fileexists(local.lambda_zip_path) ? filebase64sha256(local.lambda_zip_path) : null
  handler          = "bootstrap"
  runtime          = "provided.al2023"
  architectures    = ["arm64"]
  role             = aws_iam_role.lambda.arn
  timeout          = 30

  environment {
    variables = {
      DISCORD_WEBHOOK_URL    = var.discord_webhook_url
      OCTOPUS_API_URL        = var.octopus_api_url
      OCTOPUS_ACCOUNT_NUMBER = var.octopus_account_number
      OCTOPUS_SECRET_NAME    = aws_secretsmanager_secret.octopus_credentials.name
    }
  }

  tracing_config {
    mode = "Active"
  }

  depends_on = [
    aws_cloudwatch_log_group.lambda,
  ]
}
