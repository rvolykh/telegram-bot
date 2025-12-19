variable "name" {
  description = "Name of to use on provisioned resources"
  type        = string
}

variable "reserved_concurrent_executions" {
  description = "Reserved concurrent executions for the Lambda function"
  type        = number
}

variable "emails" {
  description = "Emails to subscribe to the SNS topic"
  type        = list(string)
  default     = []
  sensitive   = true
}

variable "telegram_chat_id" {
  description = "Telegram chat ID to send alerts to"
  type        = string
  default     = ""
  sensitive   = true
}

variable "ssm_param_telegram_api_token" {
  description = "SSM parameter name for the Telegram API token"
  type        = string
  default     = ""
  sensitive   = true
}

variable "role_policies" {
  description = "Policy documents for the IAM role"
  type        = list(list(string))
  default     = []

  validation {
    condition     = length(var.role_policies) <= 8
    error_message = "role_policies must be less than 9 policies"
  }
}

variable "tags" {
  description = "Tags to apply to apply to provisioned resources"
  type        = map(string)
  default     = {}
}
