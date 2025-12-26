module "telegram_bot_queue_alerting" {
  source = "./modules/queue"

  queue_name = "${var.prefix}telegram-bot-alerting"

  enable_dead_letter_queue = false
  dead_letter_queue_source_arns = [
    module.telegram_bot_queue_mux.sqs_queue_arn,
    module.telegram_bot_queue_cmd_poweron.sqs_queue_arn,
  ]
}

resource "aws_cloudwatch_metric_alarm" "non_empty_dlq" {
  alarm_name          = "${var.prefix}telegram-bot-non-empty-dlq"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = 1
  period              = 5 * 60
  threshold           = 1 # 1 message
  metric_name         = "ApproximateNumberOfMessagesVisible"
  namespace           = "AWS/SQS"
  statistic           = "Average"

  dimensions = {
    QueueName = module.telegram_bot_queue_alerting.sqs_queue_name
  }

  alarm_actions = [
    module.telegram_bot_alerting.sns_topic_arn,
  ]
  ok_actions = [
    module.telegram_bot_alerting.sns_topic_arn,
  ]
}

module "telegram_bot_alerting" {
  source = "./modules/alerting"

  name                           = "${var.prefix}telegram-bot-alerting"
  reserved_concurrent_executions = -1
  emails                         = var.alerting_emails
  telegram_chat_id               = var.alerting_telegram_chat_id
  ssm_param_telegram_api_token   = module.telegram_bot_api_token.name
  role_policies = [
    // policy 0
    [
      module.telegram_bot_api_token.policy_document_read_only,
    ]
    //
  ]
}
