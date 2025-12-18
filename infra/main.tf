module "telegram_bot_queue_mux" {
  source = "./modules/queue"

  queue_name = "telegram-bot-mux"

  enable_dead_letter_queue = true
  dead_letter_queue_arn    = module.telegram_bot_queue_alerting.sqs_queue_arn
}

resource "aws_cloudwatch_metric_alarm" "mux_command_rate" {
  alarm_name          = "telegram-bot-mux-command-rate"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = 3  # over 3 minutes
  period              = 60 # 1 minute
  threshold           = 10 # 10 messages
  metric_name         = "NumberOfMessagesReceived"
  namespace           = "AWS/SQS"
  statistic           = "Sum"

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

module "telegram_bot_api" {
  source = "./modules/api"

  api_name = "telegram-bot"
  sqs_queue = {
    name = module.telegram_bot_queue_mux.sqs_queue_name
    arn  = module.telegram_bot_queue_mux.sqs_queue_arn
  }
}

module "telegram_bot_handler_mux" {
  source = "./modules/handler"

  function_name = "telegram-bot-mux"
  source_path   = "${path.root}/../apps/mux"

  sqs_batch_size = 10
  sqs_queue_arn  = module.telegram_bot_queue_mux.sqs_queue_arn

  environment_variables = {
    SQS_COMMAND_POWERON_QUEUE_URL = module.telegram_bot_queue_cmd_poweron.sqs_queue_url
  }

  role_policies = [
    // policy 0
    [
      module.telegram_bot_queue_cmd_poweron.policy_document_send_message,
    ]
    //
  ]
}
