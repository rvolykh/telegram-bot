variable "aws_account_id" {
  type        = string
  description = "The AWS account ID"
  sensitive   = true
}

variable "aws_region" {
  type        = string
  description = "The AWS region"
}

variable "github_repository" {
  type        = string
  description = "The GitHub repository"
  default     = "rvolykh/telegram-bot"
}

variable "environment" {
  type        = string
  description = "The environment"
  default     = "sandbox"
}
