locals {
  lambda_zip_path = "${path.module}/../lambda.zip"
}

# --- Secrets Manager ---

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

# --- IAM ---

resource "aws_iam_role" "lambda" {
  name = "dover-lambda"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = "sts:AssumeRole"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })
}

resource "aws_iam_role_policy" "lambda" {
  name = "dover-lambda"
  role = aws_iam_role.lambda.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "logs:CreateLogStream",
          "logs:PutLogEvents",
        ]
        Resource = "${aws_cloudwatch_log_group.lambda.arn}:*"
      },
      {
        Effect = "Allow"
        Action = [
          "secretsmanager:GetSecretValue",
        ]
        Resource = aws_secretsmanager_secret.octopus_credentials.arn
      },
    ]
  })
}

# --- CloudWatch Logs ---

resource "aws_cloudwatch_log_group" "lambda" {
  name              = "/aws/lambda/dover"
  retention_in_days = var.log_retention_days
}

# --- Lambda ---

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

# --- EventBridge ---

resource "aws_cloudwatch_event_rule" "daily" {
  name                = "dover-daily"
  description         = "Trigger Dover Lambda once per day"
  schedule_expression = var.schedule_expression
}

resource "aws_cloudwatch_event_target" "lambda" {
  rule = aws_cloudwatch_event_rule.daily.name
  arn  = aws_lambda_function.dover.arn
}

resource "aws_lambda_permission" "eventbridge" {
  statement_id  = "AllowEventBridgeInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.dover.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.daily.arn
}
