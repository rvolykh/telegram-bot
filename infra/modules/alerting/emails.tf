resource "aws_sns_topic" "emails" {
  name              = "${var.name}-fallback-to-emails"
  kms_master_key_id = "alias/aws/sns"
  tags              = var.tags
}

resource "aws_sns_topic_subscription" "emails" {
  count = length(var.emails)

  topic_arn = aws_sns_topic.emails.arn
  protocol  = "email"
  endpoint  = var.emails[count.index]
}
