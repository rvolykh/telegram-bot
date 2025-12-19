resource "aws_sqs_queue" "this" {
  name                              = "${var.queue_name}.fifo"
  fifo_queue                        = true
  content_based_deduplication       = true
  kms_master_key_id                 = "alias/aws/sqs"
  kms_data_key_reuse_period_seconds = 3600
  visibility_timeout_seconds        = var.visibility_timeout_seconds
  message_retention_seconds         = var.message_retention_seconds
  tags                              = var.tags
}

resource "aws_sqs_queue_redrive_policy" "dlq" {
  count = var.enable_dead_letter_queue ? 1 : 0

  queue_url = aws_sqs_queue.this.id
  redrive_policy = jsonencode({
    deadLetterTargetArn = var.dead_letter_queue_arn
    maxReceiveCount     = 4
  })
}
