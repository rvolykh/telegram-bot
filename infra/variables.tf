variable "telegram_bot_api_token" {
  description = "Telegram bot API token"
  type        = string
  sensitive   = true
}

variable "prefix" {
  description = "Prefix to use for the resources"
  type        = string
  default     = ""
}

variable "api_ip_allowlist" {
  description = "IP addresses to allow access to the API"
  type        = list(string)
  default     = []
}

variable "alerting_telegram_chat_id" {
  description = "Telegram chat ID to send alerts to"
  type        = string
  default     = ""
  sensitive   = true
}

variable "alerting_emails" {
  description = "Emails to subscribe to the alerting SNS topic"
  type        = list(string)
  default     = []
  sensitive   = true
}
