output "webhook_url" {
  description = "Webhook URL for Telegram bot (POST /webhook)"
  value       = "${aws_api_gateway_stage.this.invoke_url}/webhook"
  sensitive   = true
}

output "sqs_queue_url" {
  description = "URL of the SQS FIFO queue"
  value       = aws_sqs_queue.webhook.url
}

output "sqs_queue_arn" {
  description = "ARN of the SQS FIFO queue"
  value       = aws_sqs_queue.webhook.arn
}

output "sqs_queue_name" {
  description = "Name of the SQS FIFO queue"
  value       = aws_sqs_queue.webhook.name
}
