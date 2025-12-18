output "sqs_queue_url" {
  description = "URL of the SQS queue"
  value       = aws_sqs_queue.this.url
}

output "sqs_queue_arn" {
  description = "ARN of the SQS queue"
  value       = aws_sqs_queue.this.arn
}

output "sqs_queue_name" {
  description = "Name of the SQS queue"
  value       = aws_sqs_queue.this.name
}

output "policy_document_send_message" {
  description = "IAM policy document for sending messages to the SQS queue"
  value       = data.aws_iam_policy_document.send_message.json
}
