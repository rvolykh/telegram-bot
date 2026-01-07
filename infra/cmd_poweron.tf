module "telegram_bot_queue_cmd_poweron" {
  source = "./modules/queue"

  queue_name               = "${var.prefix}telegram-bot-cmd-poweron"
  enable_dead_letter_queue = true
  dead_letter_queue_arn    = module.telegram_bot_queue_alerting.sqs_queue_arn
}

module "telegram_bot_db_poweron_subscriptions" {
  source = "./modules/db"

  name          = "${var.prefix}telegram-bot-poweron-subscriptions"
  hash_key_name = "ChatId"
  hash_key_type = "N"
}

module "telegram_bot_cmd_poweron" {
  source = "./modules/handler"

  function_name                  = "${var.prefix}telegram-bot-cmd-poweron"
  reserved_concurrent_executions = -1
  source_path                    = "${path.root}/../apps/poweron"

  timeout        = 30
  sqs_batch_size = 10
  sqs_queue_arn  = module.telegram_bot_queue_cmd_poweron.sqs_queue_arn

  environment_variables = {
    SSM_PARAM_TELEGRAM_APITOKEN  = module.telegram_bot_api_token.name
    SSM_PARAM_CACHE              = module.telegram_bot_cache_poweron.name
    DYNAMODB_TABLE_SUBSCRIPTIONS = module.telegram_bot_db_poweron_subscriptions.name
  }

  role_policies = [
    // policy 0
    [
      module.telegram_bot_api_token.policy_document_read_only,
      module.telegram_bot_cache_poweron.policy_document_read_write,
      module.telegram_bot_db_poweron_subscriptions.policy_document_read_write,
    ],
    //
  ]
}

module "telegram_bot_cron_poweron" {
  source = "./modules/cron"

  name                = "${var.prefix}telegram-bot-poweron-subscriptions"
  schedule_expression = "rate(20 minutes)"
  lambda_function_arn = module.telegram_bot_cmd_poweron.lambda_function_arn
}
