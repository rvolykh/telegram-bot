variable "name" {
  description = "Name of the parameter"
  type        = string
}

variable "value" {
  description = "Value of the parameter"
  type        = string
  sensitive   = true
}
