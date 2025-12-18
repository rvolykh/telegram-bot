module "telegram_bot_queue_cmd_poweron" {
  source = "./modules/queue"

  queue_name               = "telegram-bot-cmd-poweron"
  enable_dead_letter_queue = true
  dead_letter_queue_arn    = module.telegram_bot_queue_alerting.sqs_queue_arn
}

module "telegram_bot_cmd_poweron" {
  source = "./modules/handler"

  function_name = "telegram-bot-cmd-poweron"
  source_path   = "${path.root}/../apps/poweron"

  sqs_batch_size = 10
  sqs_queue_arn  = module.telegram_bot_queue_cmd_poweron.sqs_queue_arn

  environment_variables = {
    SSM_PARAM_TELEGRAM_APITOKEN = module.telegram_bot_api_token.name
    SSM_PARAM_POWERON_CACHE     = module.telegram_bot_cache_poweron.name
  }

  role_policies = [
    // policy 0
    [
      module.telegram_bot_api_token.policy_document_read_only,
      module.telegram_bot_cache_poweron.policy_document_read_write,
    ],
    //
  ]
}
