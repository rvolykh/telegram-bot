resource "aws_lambda_function" "this" {
  filename                       = data.archive_file.lambda_zip.output_path
  function_name                  = var.function_name
  role                           = aws_iam_role.lambda_execution.arn
  handler                        = "bootstrap"
  source_code_hash               = data.archive_file.lambda_zip.output_base64sha256
  runtime                        = "provided.al2023"
  timeout                        = var.timeout
  memory_size                    = var.memory_size
  kms_key_arn                    = "arn:aws:kms:${data.aws_region.current.region}:${data.aws_caller_identity.current.account_id}:alias/aws/lambda"
  reserved_concurrent_executions = 10

  environment {
    variables = var.environment_variables
  }

  tracing_config {
    mode = "PassThrough"
  }

  tags = var.tags
}

resource "aws_lambda_event_source_mapping" "sqs" {
  event_source_arn = var.sqs_queue_arn
  function_name    = aws_lambda_function.this.arn
  enabled          = true

  batch_size = var.sqs_batch_size
}
