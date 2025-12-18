variable "function_name" {
  description = "Name of the Lambda function"
  type        = string
}

variable "source_path" {
  description = "Path to the directory containing main.go (the module will build the binary and create a ZIP)"
  type        = string
}

variable "environment_variables" {
  description = "Environment variables for the Lambda function"
  type        = map(string)
  default     = {}
}

variable "timeout" {
  description = "Lambda function timeout in seconds"
  type        = number
  default     = 3
}

variable "memory_size" {
  description = "Amount of memory in MB your Lambda function can use"
  type        = number
  default     = 128
}

variable "tags" {
  description = "Tags to apply to the Lambda function"
  type        = map(string)
  default     = {}
}

variable "sqs_queue_arn" {
  description = "ARN of the SQS FIFO queue to subscribe to"
  type        = string
}

variable "sqs_batch_size" {
  description = "Maximum number of records in each batch that Lambda pulls from SQS"
  type        = number
  default     = 1
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
