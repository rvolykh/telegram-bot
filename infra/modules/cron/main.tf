resource "aws_scheduler_schedule" "this" {
  name       = var.name
  group_name = "default"

  flexible_time_window {
    mode                      = "FLEXIBLE"
    maximum_window_in_minutes = 10
  }

  schedule_expression = var.schedule_expression

  target {
    arn      = var.lambda_function_arn
    role_arn = aws_iam_role.scheduler.arn

    retry_policy {
      maximum_retry_attempts = 3
    }
  }
}
