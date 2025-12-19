resource "aws_sns_topic" "this" {
  name              = var.name
  kms_master_key_id = "alias/aws/sns"
  tags              = var.tags
}

resource "aws_sns_topic_policy" "this" {
  arn    = aws_sns_topic.this.arn
  policy = data.aws_iam_policy_document.sns_topic_resource_policy.json
}

resource "aws_sns_topic_subscription" "primary_telegram" {
  topic_arn = aws_sns_topic.this.arn
  protocol  = "lambda"
  endpoint  = aws_lambda_function.telegram.arn
}

resource "aws_lambda_function_event_invoke_config" "secondary_emails" {
  function_name = aws_lambda_function.telegram.function_name

  maximum_event_age_in_seconds = 60
  maximum_retry_attempts       = 2

  destination_config {
    on_failure {
      destination = aws_sns_topic.emails.arn
    }
  }

  depends_on = [
    # Ensure Role has necessary permissions
    aws_iam_role_policy.fallback_to_sns,
  ]
}
