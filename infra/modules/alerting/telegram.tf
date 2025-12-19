resource "aws_lambda_function" "telegram" {
  filename                       = data.archive_file.lambda_zip.output_path
  function_name                  = var.name
  role                           = aws_iam_role.lambda_execution.arn
  handler                        = "bootstrap"
  source_code_hash               = data.archive_file.lambda_zip.output_base64sha256
  runtime                        = "provided.al2023"
  timeout                        = 10
  memory_size                    = 128
  reserved_concurrent_executions = 10

  environment {
    variables = {
      TELEGRAM_CHAT_ID            = var.telegram_chat_id
      SSM_PARAM_TELEGRAM_APITOKEN = var.ssm_param_telegram_api_token
    }
  }

  tracing_config {
    mode = "PassThrough"
  }

  tags = var.tags
}

resource "aws_cloudwatch_log_group" "lambda" {
  name              = "/aws/lambda/${var.name}"
  retention_in_days = 7
  tags              = var.tags
}


resource "aws_lambda_permission" "telegram_sns" {
  statement_id  = "SNS"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.telegram.function_name
  principal     = "sns.amazonaws.com"
  source_arn    = aws_sns_topic.this.arn
}

resource "aws_iam_role" "lambda_execution" {
  name = "${var.name}-func-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })
  tags = var.tags
}

resource "aws_iam_role_policy_attachment" "lambda_basic_execution" {
  role       = aws_iam_role.lambda_execution.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_role_policy" "fallback_to_sns" {
  name   = "${var.name}-fallback-to-sns-policy"
  role   = aws_iam_role.lambda_execution.name
  policy = data.aws_iam_policy_document.fallback_to_sns.json
}

resource "aws_iam_policy" "role_policies" {
  count = length(var.role_policies)

  name   = "${var.name}-role-policy-${count.index}"
  policy = data.aws_iam_policy_document.role_policies[count.index].json
  tags   = var.tags
}

resource "aws_iam_role_policy_attachment" "role_policies" {
  count = length(var.role_policies)

  role       = aws_iam_role.lambda_execution.name
  policy_arn = aws_iam_policy.role_policies[count.index].arn
}
