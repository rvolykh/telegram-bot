variable "queue_name" {
  description = "Name of the SQS queue"
  type        = string
}

variable "visibility_timeout_seconds" {
  description = "Visibility timeout in seconds"
  type        = number
  default     = 30
}

variable "message_retention_seconds" {
  description = "Message retention in seconds"
  type        = number
  default     = 4 * 24 * 60 * 60 # 4 days
}

variable "enable_dead_letter_queue" {
  description = "Enable dead letter queue"
  type        = bool
  default     = false
}

variable "dead_letter_queue_arn" {
  description = "Name of the dead letter queue"
  type        = string
  default     = null

  validation {
    condition     = var.enable_dead_letter_queue ? var.dead_letter_queue_arn != null : true
    error_message = "dead_letter_queue_arn is required when enable_dead_letter_queue is true"
  }
}

variable "dead_letter_queue_source_arns" {
  description = "ARNs of the source queues to allow redriving"
  type        = list(string)
  default     = []

  validation {
    condition     = var.enable_dead_letter_queue ? length(var.dead_letter_queue_source_arns) == 0 : true
    error_message = "dead_letter_queue_source_arns must be empty when enable_dead_letter_queue is true"
  }
}

variable "tags" {
  description = "Tags to apply to the SQS queue"
  type        = map(string)
  default     = {}
}
