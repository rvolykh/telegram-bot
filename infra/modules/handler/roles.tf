resource "aws_iam_role" "lambda_execution" {
  name = "${var.function_name}-func-role"

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
}

resource "aws_iam_role_policy_attachment" "lambda_basic_execution" {
  role       = aws_iam_role.lambda_execution.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_role_policy" "sqs" {
  name = "${var.function_name}-sqs-policy"
  role = aws_iam_role.lambda_execution.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "sqs:ReceiveMessage",
          "sqs:DeleteMessage",
          "sqs:GetQueueAttributes"
        ]
        Resource = var.sqs_queue_arn
      }
    ]
  })
}

resource "aws_iam_policy" "role_policies" {
  count = length(var.role_policies)

  name   = "${var.function_name}-role-policy-${count.index}"
  policy = data.aws_iam_policy_document.role_policies[count.index].json
}

resource "aws_iam_role_policy_attachment" "role_policies" {
  count = length(var.role_policies)

  role       = aws_iam_role.lambda_execution.name
  policy_arn = aws_iam_policy.role_policies[count.index].arn
}
