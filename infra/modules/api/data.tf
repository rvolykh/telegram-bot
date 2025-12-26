data "aws_region" "current" {}

data "aws_caller_identity" "current" {}

data "aws_iam_policy_document" "cloudwatch_permissions" {
  statement {
    effect = "Allow"
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:DescribeLogGroups",
      "logs:DescribeLogStreams",
      "logs:PutLogEvents",
      "logs:GetLogEvents",
      "logs:FilterLogEvents"
    ]
    resources = [
      "*",
    ]
  }
}

data "aws_iam_policy_document" "api_gateway_policy" {
  statement {
    effect    = "Allow"
    actions   = ["execute-api:Invoke"]
    resources = ["${aws_api_gateway_rest_api.this.execution_arn}/*/*/*"]

    principals {
      type        = "AWS"
      identifiers = ["*"]
    }
  }

  statement {
    effect    = "Deny"
    actions   = ["execute-api:Invoke"]
    resources = ["${aws_api_gateway_rest_api.this.execution_arn}/*/*/*"]

    principals {
      type        = "AWS"
      identifiers = ["*"]
    }

    # https://core.telegram.org/resources/cidr.txt
    condition {
      test     = "NotIpAddress"
      variable = "aws:SourceIp"
      values = concat([
        "91.108.56.0/22",
        "91.108.4.0/22",
        "91.108.8.0/22",
        "91.108.16.0/22",
        "91.108.12.0/22",
        "149.154.160.0/20",
        "91.105.192.0/23",
        "91.108.20.0/22",
        "185.76.151.0/24",
        "2001:b28:f23d::/48",
        "2001:b28:f23f::/48",
        "2001:67c:4e8::/48",
        "2001:b28:f23c::/48",
        "2a0a:f280::/32",
      ], var.ip_allowlist)
    }
  }
}
