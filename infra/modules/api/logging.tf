resource "aws_api_gateway_account" "this" {
  cloudwatch_role_arn = aws_iam_role.logging.arn

  depends_on = [
    // Wait for the logging role permissions to be attached
    aws_iam_role_policy.logging,
  ]
}

resource "aws_cloudwatch_log_group" "stage_v1" {
  name              = "API-Gateway-Execution-Logs_${aws_api_gateway_rest_api.this.id}/${local.stage_v1}"
  retention_in_days = 7

  tags = var.tags
}

resource "aws_cloudwatch_log_group" "access_logs" {
  name              = "${aws_api_gateway_rest_api.this.id}-access-logs"
  retention_in_days = 7

  tags = var.tags
}
