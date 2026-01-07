resource "aws_iam_role" "scheduler" {
  name = "${var.name}-scheduler-role"

  assume_role_policy = data.aws_iam_policy_document.assume_role_policy.json

  tags = var.tags
}

resource "aws_iam_role_policy" "lambda_invoke" {
  name = "${var.name}-lambda-invoke"
  role = aws_iam_role.scheduler.id

  policy = data.aws_iam_policy_document.role_permissions.json
}
