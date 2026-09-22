variable "discord_webhook_url" {
  description = "Discord Webhook URL for sending notifications"
  type        = string
  sensitive   = true
}

variable "octopus_api_url" {
  description = "Octopus Energy GraphQL API endpoint"
  type        = string
  default     = "https://api.oejp-kraken.energy/v1/graphql/"
}

variable "octopus_account_number" {
  description = "Octopus Energy account number"
  type        = string
  sensitive   = true
}

variable "schedule_expression" {
  description = "EventBridge schedule expression for triggering Lambda"
  type        = string
  default     = "cron(0 9 * * ? *)"
}

variable "log_retention_days" {
  description = "CloudWatch Logs retention period in days"
  type        = number
  default     = 14
}
