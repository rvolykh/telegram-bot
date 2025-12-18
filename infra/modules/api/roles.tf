resource "aws_iam_role" "api_gateway_sqs" {
  name = "${var.api_name}-api-gateway-sqs-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "apigateway.amazonaws.com"
        }
      }
    ]
  })

  tags = var.tags
}

resource "aws_iam_role_policy" "api_gateway_sqs" {
  name = "${var.api_name}-api-gateway-sqs-policy"
  role = aws_iam_role.api_gateway_sqs.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "sqs:SendMessage",
        ]
        Resource = var.sqs_queue.arn
      }
    ]
  })
}

resource "aws_iam_role" "logging" {
  name = "${var.api_name}-logging-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "apigateway.amazonaws.com"
        }
      }
    ]
  })

  tags = var.tags
}

resource "aws_iam_role_policy" "logging" {
  name = "${var.api_name}-logging-policy"
  role = aws_iam_role.logging.id

  policy = data.aws_iam_policy_document.cloudwatch_permissions.json
}
