data "aws_region" "current" {}
data "aws_caller_identity" "current" {}

data "aws_iam_policy_document" "role_policies" {
  count = length(var.role_policies)

  source_policy_documents = var.role_policies[count.index]
}

data "aws_iam_policy_document" "fallback_to_sns" {
  version = "2012-10-17"

  statement {
    effect = "Allow"
    actions = [
      "sns:Publish",
    ]
    resources = [
      aws_sns_topic.emails.arn,
    ]
  }
}

data "aws_iam_policy_document" "sns_topic_resource_policy" {
  version = "2012-10-17"

  statement {
    principals {
      type = "Service"
      identifiers = [
        "cloudwatch.amazonaws.com",
      ]
    }

    effect = "Allow"
    actions = [
      "sns:Publish",
    ]
    resources = [
      aws_sns_topic.this.arn,
    ]
  }
}
