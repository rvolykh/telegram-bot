variable "name" {
  description = "Name of the DynamoDB table"
  type        = string
}

variable "hash_key_name" {
  description = "Name of the hash key"
  type        = string
}

variable "hash_key_type" {
  description = "Type of the hash key"
  type        = string

  validation {
    condition     = contains(["S", "N", "B"], var.hash_key_type)
    error_message = "hash_key_type must be one of: S, N, B"
  }
}

variable "tags" {
  description = "Tags to apply to apply to provisioned resources"
  type        = map(string)
  default     = {}
}
