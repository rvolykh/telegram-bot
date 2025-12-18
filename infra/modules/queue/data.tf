data "aws_iam_policy_document" "send_message" {
  version = "2012-10-17"
  statement {
    effect = "Allow"
    actions = [
      "sqs:SendMessage",
    ]
    resources = [
      aws_sqs_queue.this.arn,
    ]
  }
}
