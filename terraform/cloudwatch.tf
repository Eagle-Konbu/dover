resource "aws_cloudwatch_log_group" "lambda" {
  name              = "/aws/lambda/dover"
  retention_in_days = var.log_retention_days
}
