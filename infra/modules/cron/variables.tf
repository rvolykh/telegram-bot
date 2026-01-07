variable "name" {
  description = "Name to use for resources"
  type        = string
}

variable "schedule_expression" {
  description = "Schedule"
  type        = string
}

variable "lambda_function_arn" {
  description = "ARN of the Lambda function to invoke"
  type        = string
}

variable "maximum_retry_attempts" {
  description = "Maximum number of retry attempts to make before the request fails"
  type        = number
  default     = 3
}

variable "tags" {
  description = "Tags to apply to resources"
  type        = map(string)
  default     = {}
}
