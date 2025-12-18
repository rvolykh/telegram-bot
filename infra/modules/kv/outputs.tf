output "name" {
  description = "Name of the parameter"
  value       = aws_ssm_parameter.this.name
}

output "arn" {
  description = "ARN of the parameter"
  value       = aws_ssm_parameter.this.arn
}

output "policy_document_read_only" {
  description = "IAM policy document for read-only access"
  value       = data.aws_iam_policy_document.read_only.json
}

output "policy_document_write_only" {
  description = "IAM policy document for write-only access"
  value       = data.aws_iam_policy_document.write_only.json
}

output "policy_document_read_write" {
  description = "IAM policy document for read-write access"
  value       = data.aws_iam_policy_document.read_write.json
}
