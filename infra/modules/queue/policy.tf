resource "aws_sqs_queue_redrive_allow_policy" "dlq" {
  count = length(var.dead_letter_queue_source_arns) > 0 ? 1 : 0

  queue_url = aws_sqs_queue.this.id
  redrive_allow_policy = jsonencode({
    redrivePermission = "byQueue",
    sourceQueueArns   = var.dead_letter_queue_source_arns
  })
}
