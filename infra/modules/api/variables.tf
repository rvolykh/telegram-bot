variable "api_name" {
  description = "Name of the API Gateway REST API"
  type        = string
}

variable "sqs_queue" {
  description = "SQS queue to route traffic to"
  type = object({
    name = string
    arn  = string
  })

  validation {
    condition     = var.sqs_queue.name != null
    error_message = "sqs_queue.name is required"
  }
  validation {
    condition     = var.sqs_queue.arn != null
    error_message = "sqs_queue.arn is required"
  }
  validation {
    condition     = strcontains(var.sqs_queue.arn, var.sqs_queue.name)
    error_message = "sqs_queue.arn must contain the sqs_queue.name"
  }
}

variable "ip_allowlist" {
  description = "IP addresses to allow access to the API"
  type        = list(string)
  default     = []
}

variable "tags" {
  description = "Tags to apply to the API Gateway and SQS resources"
  type        = map(string)
  default     = {}
}
