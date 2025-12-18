resource "aws_sqs_queue" "webhook" {
  name                        = "${var.api_name}-webhook.fifo"
  fifo_queue                  = true
  content_based_deduplication = true

  tags = var.tags
}

resource "aws_api_gateway_rest_api" "this" {
  name = var.api_name

  endpoint_configuration {
    types = ["REGIONAL"]
  }

  tags = var.tags
}

resource "aws_api_gateway_rest_api_policy" "this" {
  rest_api_id = aws_api_gateway_rest_api.this.id
  policy      = data.aws_iam_policy_document.api_gateway_policy.json
}

resource "aws_api_gateway_resource" "webhook" {
  rest_api_id = aws_api_gateway_rest_api.this.id
  parent_id   = aws_api_gateway_rest_api.this.root_resource_id
  path_part   = "webhook"
}

resource "aws_api_gateway_method" "webhook" {
  rest_api_id   = aws_api_gateway_rest_api.this.id
  resource_id   = aws_api_gateway_resource.webhook.id
  http_method   = "POST"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "sqs" {
  rest_api_id = aws_api_gateway_rest_api.this.id
  resource_id = aws_api_gateway_resource.webhook.id
  http_method = aws_api_gateway_method.webhook.http_method

  type                    = "AWS"
  integration_http_method = "POST"
  uri                     = "arn:aws:apigateway:${data.aws_region.current.region}:sqs:path/${data.aws_caller_identity.current.account_id}/${var.sqs_queue.name}"
  credentials             = aws_iam_role.api_gateway_sqs.arn

  passthrough_behavior = "NEVER"

  request_templates = {
    "application/json" = <<EOT
#set($dedupId = $input.json('$.update_id'))
#set($groupId = $input.json('$.message.chat.id'))
Action=SendMessage&MessageBody=$util.urlEncode($input.body)&MessageGroupId=$groupId&MessageDeduplicationId=$dedupId
EOT
  }

  request_parameters = {
    "integration.request.header.Content-Type" = "'application/x-www-form-urlencoded'"
  }
}

resource "aws_api_gateway_method_response" "webhook_200" {
  rest_api_id = aws_api_gateway_rest_api.this.id
  resource_id = aws_api_gateway_resource.webhook.id
  http_method = aws_api_gateway_method.webhook.http_method
  status_code = "200"

  response_models = {
    "application/json" = "Empty"
  }

  depends_on = [
    // Wait for the SQS integration to be created
    aws_api_gateway_integration.sqs,
  ]
}

resource "aws_api_gateway_integration_response" "webhook_200" {
  rest_api_id = aws_api_gateway_rest_api.this.id
  resource_id = aws_api_gateway_resource.webhook.id
  http_method = aws_api_gateway_method.webhook.http_method
  status_code = aws_api_gateway_method_response.webhook_200.status_code

  response_templates = {
    "application/json" = "{\"status\":\"ok\"}"
  }

  depends_on = [
    // Wait for the SQS integration to be created
    aws_api_gateway_integration.sqs,
  ]
}

resource "aws_api_gateway_deployment" "this" {
  rest_api_id = aws_api_gateway_rest_api.this.id

  triggers = {
    redeployment = sha1(jsonencode([
      aws_api_gateway_resource.webhook,
      aws_api_gateway_method.webhook,
      aws_api_gateway_integration.sqs,
      aws_api_gateway_integration_response.webhook_200,
      aws_api_gateway_rest_api_policy.this,
    ]))
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_api_gateway_stage" "this" {
  deployment_id = aws_api_gateway_deployment.this.id
  rest_api_id   = aws_api_gateway_rest_api.this.id
  stage_name    = local.stage_v1

  access_log_settings {
    destination_arn = aws_cloudwatch_log_group.access_logs.arn
    format          = "{ \"requestId\":\"$context.requestId\", \"extendedRequestId\":\"$context.extendedRequestId\", \"ip\":\"$context.identity.sourceIp\", \"caller\":\"$context.identity.caller\", \"user\":\"$context.identity.user\", \"requestTime\":\"$context.requestTime\", \"httpMethod\":\"$context.httpMethod\", \"resourcePath\":\"$context.resourcePath\", \"status\":\"$context.status\", \"protocol\":\"$context.protocol\", \"responseLength\":\"$context.responseLength\" }"
  }

  tags = var.tags

  depends_on = [
    // Wait for the CloudWatch log group to be created
    aws_cloudwatch_log_group.stage_v1,
    // Wait for the API Gateway account to be created
    aws_api_gateway_account.this,
  ]
}

resource "aws_api_gateway_method_settings" "this" {
  rest_api_id = aws_api_gateway_rest_api.this.id
  stage_name  = aws_api_gateway_stage.this.stage_name
  method_path = "*/*"

  settings {
    logging_level          = "INFO"
    metrics_enabled        = true
    throttling_rate_limit  = 20
    throttling_burst_limit = 10
  }
}
