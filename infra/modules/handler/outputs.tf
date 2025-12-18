output "lambda_function_name" {
  description = "Name of the Lambda function"
  value       = aws_lambda_function.this.function_name
}

output "lambda_function_arn" {
  description = "ARN of the Lambda function"
  value       = aws_lambda_function.this.arn
}

output "lambda_function_invoke_arn" {
  description = "ARN to be used for invoking Lambda function"
  value       = aws_lambda_function.this.invoke_arn
}

output "sqs_event_source_mapping_id" {
  description = "ID of the SQS event source mapping"
  value       = aws_lambda_event_source_mapping.sqs.id
}

output "sqs_event_source_mapping_uuid" {
  description = "UUID of the SQS event source mapping"
  value       = aws_lambda_event_source_mapping.sqs.uuid
}

output "lambda_function_qualified_arn" {
  description = "Qualified ARN of the Lambda function"
  value       = aws_lambda_function.this.qualified_arn
}

output "lambda_function_version" {
  description = "Latest published version of the Lambda function"
  value       = aws_lambda_function.this.version
}

output "lambda_role_arn" {
  description = "ARN of the IAM role used by the Lambda function"
  value       = aws_iam_role.lambda_execution.arn
}

output "lambda_role_name" {
  description = "Name of the IAM role used by the Lambda function"
  value       = aws_iam_role.lambda_execution.name
}
