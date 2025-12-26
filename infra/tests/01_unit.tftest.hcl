run "prepare" {
  module {
    source = "./tests/utils/inputs"
  }
}

run "verify_plan" {
  command = plan

  variables {
    telegram_bot_api_token = run.prepare.token
    prefix                 = run.prepare.prefix
  }

  # Configs
  assert {
    condition     = module.telegram_bot_api_token.name == "/${run.prepare.prefix}telegram/bot/api_token"
    error_message = "module telegram_bot_api_token should create the SSM parameter /${run.prepare.prefix}telegram/bot/api_token"
  }
  assert {
    condition     = module.telegram_bot_cache_poweron.name == "/${run.prepare.prefix}telegram/bot/cache/poweron"
    error_message = "module telegram_bot_cache_poweron should create the SSM parameter /${run.prepare.prefix}telegram/bot/cache/poweron"
  }

  # Queues
  assert {
    condition     = module.telegram_bot_api.sqs_queue_name == "${run.prepare.prefix}telegram-bot-webhook.fifo"
    error_message = "module telegram_bot_api should create the SQS queue ${run.prepare.prefix}telegram-bot-webhook.fifo"
  }
  assert {
    condition     = module.telegram_bot_queue_mux.sqs_queue_name == "${run.prepare.prefix}telegram-bot-mux.fifo"
    error_message = "module telegram_bot_queue_mux should create the SQS queue ${run.prepare.prefix}telegram-bot-mux.fifo"
  }
  assert {
    condition     = module.telegram_bot_queue_alerting.sqs_queue_name == "${run.prepare.prefix}telegram-bot-alerting.fifo"
    error_message = "module telegram_bot_queue_alerting should create the SQS queue ${run.prepare.prefix}telegram-bot-alerting.fifo"
  }
  assert {
    condition     = module.telegram_bot_queue_cmd_poweron.sqs_queue_name == "${run.prepare.prefix}telegram-bot-cmd-poweron.fifo"
    error_message = "module telegram_bot_queue_cmd_poweron should create the SQS queue ${run.prepare.prefix}telegram-bot-cmd-poweron.fifo"
  }

  # Functions
  assert {
    condition     = module.telegram_bot_handler_mux.lambda_function_name == "${run.prepare.prefix}telegram-bot-mux"
    error_message = "module telegram_bot_handler_mux should create the Lambda function ${run.prepare.prefix}telegram-bot-mux"
  }
  assert {
    condition     = module.telegram_bot_cmd_poweron.lambda_function_name == "${run.prepare.prefix}telegram-bot-cmd-poweron"
    error_message = "module telegram_bot_cmd_poweron should create the Lambda function ${run.prepare.prefix}telegram-bot-cmd-poweron"
  }
}

run "verify_module_api_positive" {
  command = plan

  module {
    source = "./modules/api"
  }

  variables {
    api_name = "${run.prepare.prefix}api-positive"
    sqs_queue = {
      name = "${run.prepare.prefix}api.fifo"
      arn  = "arn:aws:sqs:000000000000:us-east-1:${run.prepare.prefix}api.fifo"
    }
  }

  assert {
    condition     = output.sqs_queue_name == "${run.prepare.prefix}api-positive-webhook.fifo"
    error_message = "module api_positive should create the SQS queue ${run.prepare.prefix}api-positive-webhook.fifo"
  }
}

run "verify_module_api_negative" {
  command = plan

  module {
    source = "./modules/api"
  }

  variables {
    api_name = "${run.prepare.prefix}api-negative"
    sqs_queue = {
      name = "${run.prepare.prefix}my-queue.fifo"
      arn  = "arn:aws:sqs:000000000000:us-east-1:${run.prepare.prefix}unknown-queue.fifo"
    }
  }

  expect_failures = [var.sqs_queue]
}

run "verify_module_queue_positive" {
  command = plan

  module {
    source = "./modules/queue"
  }

  variables {
    queue_name = "${run.prepare.prefix}queue-positive"
  }

  assert {
    condition     = output.sqs_queue_name == "${run.prepare.prefix}queue-positive.fifo"
    error_message = "module queue_positive should create the SQS queue ${run.prepare.prefix}queue-positive.fifo"
  }
}

run "verify_module_queue_negative" {
  command = plan

  module {
    source = "./modules/queue"
  }

  variables {
    queue_name               = "${run.prepare.prefix}queue-negative"
    enable_dead_letter_queue = true
  }

  expect_failures = [var.dead_letter_queue_arn]
}

run "verify_module_kv_positive" {
  command = plan

  module {
    source = "./modules/kv"
  }

  variables {
    name  = "${run.prepare.prefix}kv-positive"
    value = "test"
  }

  assert {
    condition     = output.name == "${run.prepare.prefix}kv-positive"
    error_message = "module kv_positive should create the SSM parameter ${run.prepare.prefix}kv-positive"
  }
}

run "verify_module_handler_positive" {
  command = plan

  module {
    source = "./modules/handler"
  }

  variables {
    function_name                  = "${run.prepare.prefix}handler-positive"
    reserved_concurrent_executions = 1
    source_path                    = "../../apps/mux"
    sqs_queue_arn                  = "arn:aws:sqs:000000000000:us-east-1:${run.prepare.prefix}queue.fifo"
    sqs_batch_size                 = 10
    role_policies                  = [["{}"]]
  }

  assert {
    condition     = output.lambda_function_name == "${run.prepare.prefix}handler-positive"
    error_message = "module handler_positive should create the Lambda function ${run.prepare.prefix}handler-positive"
  }
}

run "verify_module_handler_negative" {
  command = plan

  module {
    source = "./modules/handler"
  }

  variables {
    function_name                  = "${run.prepare.prefix}handler-negative"
    reserved_concurrent_executions = 1
    source_path                    = "../../apps/mux"
    sqs_queue_arn                  = "arn:aws:sqs:000000000000:us-east-1:${run.prepare.prefix}queue.fifo"
    sqs_batch_size                 = 10
    role_policies                  = [["1"], ["2"], ["3"], ["4"], ["5"], ["6"], ["7"], ["8"], ["9"]]
  }

  expect_failures = [var.role_policies]
}

run "verify_module_alerting_positive" {
  command = plan

  module {
    source = "./modules/alerting"
  }

  variables {
    name                           = "${run.prepare.prefix}alerting-positive"
    reserved_concurrent_executions = 1
    role_policies                  = [["{}"]]
  }
}

run "verify_module_alerting_negative" {
  command = plan

  module {
    source = "./modules/alerting"
  }

  variables {
    name                           = "${run.prepare.prefix}alerting-negative"
    reserved_concurrent_executions = 1
    role_policies                  = [["1"], ["2"], ["3"], ["4"], ["5"], ["6"], ["7"], ["8"], ["9"]]
  }

  expect_failures = [var.role_policies]
}
