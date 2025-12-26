run "prepare" {
  module {
    source = "./tests/utils/inputs"
  }
}

variables {
  telegram_bot_api_token = run.prepare.token
  prefix                 = run.prepare.prefix
  api_ip_allowlist       = ["${run.prepare.my_ip}/32"]
}

override_resource {
  # We don't want to replace currently applied global account configuration for Lambda
  target = module.telegram_bot_api.aws_api_gateway_account.this
}

run "apply" {
  command = apply
}

run "verify" {
  module {
    source = "./tests/utils/http"
  }

  variables {
    url = run.apply.api_gateway_url
  }

  assert {
    condition     = data.http.request.status_code == 200
    error_message = "HTTP request to Webhook URL should return 200"
  }
}
