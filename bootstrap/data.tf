data "aws_iam_policy_document" "github_trust_policy" {
  statement {
    actions = ["sts:AssumeRoleWithWebIdentity"]
    effect  = "Allow"
    principals {
      type        = "Federated"
      identifiers = [aws_iam_openid_connect_provider.github.arn]
    }
    condition {
      test     = "StringEquals"
      variable = "token.actions.githubusercontent.com:aud"
      values   = ["sts.amazonaws.com"]
    }
    condition {
      test     = "StringEquals"
      variable = "token.actions.githubusercontent.com:sub"
      values   = ["repo:${var.github_repository}:environment:${var.environment}"]
    }
  }
}

data "aws_iam_policy_document" "github_policy_1" {
  statement {
    sid = "TerraformState"
    actions = [
      "s3:GetObject",
      "s3:PutObject",
      "s3:ListBucket",
    ]
    resources = [
      "${aws_s3_bucket.state.arn}/*",
      aws_s3_bucket.state.arn,
    ]
  }
}

data "aws_iam_policy_document" "github_policy_2" {
  statement {
    sid = "IAMGlobal"
    actions = [
      "iam:ListRoles",
      "iam:ListPolicies",
      "iam:SimulateCustomPolicy",
      "iam:GetContextKeysForCustomPolicy",
    ]
    resources = [
      "*",
    ]
  }

  statement {
    sid = "IAMRole"
    actions = [
      "iam:CreateRole",
      "iam:GetRole*",
      "iam:ListRole*",
      "iam:UpdateRole*",
      "iam:DeleteRole*",
      "iam:TagRole",
      "iam:UntagRole",
      "iam:ListAttachedRolePolicies",
      "iam:AttachRolePolicy",
      "iam:DetachRolePolicy",
      "iam:UpdateAssumeRolePolicy",
      "iam:PassRole",
      "iam:PutRole*",
      "iam:SimulatePrincipalPolicy",
      "iam:ListInstanceProfilesForRole",
    ]
    resources = [
      "*",
    ]
    condition {
      test     = "StringNotEquals"
      variable = "aws:ResourceTag/ManagedBy"
      values   = ["bootstrap"]
    }
  }

  statement {
    sid = "IAMPolicy"
    actions = [
      "iam:CreatePolicy*",
      "iam:GetPolicy*",
      "iam:ListPolicy*",
      "iam:UpdatePolicy*",
      "iam:DeletePolicy*",
      "iam:TagPolicy",
      "iam:UntagPolicy",
      "iam:ListEntitiesForPolicy",
      "iam:SetDefaultPolicyVersion",
    ]
    resources = [
      "*",
    ]
    condition {
      test     = "StringNotEquals"
      variable = "aws:ResourceTag/ManagedBy"
      values   = ["bootstrap"]
    }
  }

  statement {
    sid = "LambdaGlobal"
    actions = [
      "lambda:ListFunctions",
      "lambda:ListLayerVersions",
      "lambda:ListLayers",
      "lambda:ListEventSourceMappings",
      "lambda:ListCodeSigningConfigs",
      "lambda:ListCapacityProviders",
      "lambda:GetAccountSettings",
      "lambda:CreateCodeSigningConfig",
    ]
    resources = [
      "*",
    ]
  }

  statement {
    sid = "LambdaFunctions"
    actions = [
      "lambda:*",
    ]
    resources = [
      "*",
    ]
    condition {
      test     = "StringNotEquals"
      variable = "aws:ResourceTag/ManagedBy"
      values   = ["bootstrap"]
    }
  }

  statement {
    sid = "APIGW"
    actions = [
      "apigateway:*",
    ]
    resources = [
      "*",
    ]
  }
}

data "aws_iam_policy_document" "github_policy_3" {
  statement {
    sid = "SQS"
    actions = [
      "sqs:CreateQueue",
      "sqs:DeleteQueue",
      "sqs:GetQueueAttributes",
      "sqs:GetQueueUrl",
      "sqs:ListQueues",
      "sqs:ListQueueTags",
      "sqs:SetQueueAttributes",
      "sqs:TagQueue",
      "sqs:UntagQueue",
      "sqs:ListDeadLetterSourceQueues",
    ]
    resources = [
      "*",
    ]
    condition {
      test     = "StringNotEquals"
      variable = "aws:ResourceTag/ManagedBy"
      values   = ["bootstrap"]
    }
  }
}

data "aws_iam_policy_document" "github_policy_4" {
  statement {
    sid = "CloudWatchLogs"
    actions = [
      "logs:CreateLogGroup",
      "logs:DeleteLogGroup",
      "logs:DescribeLogGroups",
      "logs:PutRetentionPolicy",
      "logs:ListTags*",
      "logs:TagLogGroup",
      "logs:UntagLogGroup",
    ]
    resources = [
      "*",
    ]
    condition {
      test     = "StringNotEquals"
      variable = "aws:ResourceTag/ManagedBy"
      values   = ["bootstrap"]
    }
  }

  statement {
    sid = "CloudWatchMetrics"
    actions = [
      "cloudwatch:PutMetricAlarm",
      "cloudwatch:DeleteAlarms",
      "cloudwatch:DescribeAlarms",
      "cloudwatch:ListTagsForResource",
      "cloudwatch:TagResource",
      "cloudwatch:UntagResource",
    ]
    resources = [
      "*",
    ]
    condition {
      test     = "StringNotEquals"
      variable = "aws:ResourceTag/ManagedBy"
      values   = ["bootstrap"]
    }
  }

  statement {
    sid = "CloudWatchGlobal"
    actions = [
      "cloudwatch:ListMetrics",
      "cloudwatch:GetMetricStatistics",
      "cloudwatch:GetMetricData",
      "cloudwatch:DescribeAlarmHistory",
    ]
    resources = [
      "*",
    ]
  }
}

data "aws_iam_policy_document" "github_policy_5" {
  statement {
    sid = "SNS"
    actions = [
      "sns:CreateTopic",
      "sns:DeleteTopic",
      "sns:GetTopicAttributes",
      "sns:ListTopics",
      "sns:ListTagsForResource",
      "sns:SetTopicAttributes",
      "sns:TagResource",
      "sns:UntagResource",
      "sns:Subscribe",
      "sns:Unsubscribe",
      "sns:ListSubscriptions",
      "sns:ListSubscriptionsByTopic",
      "sns:GetSubscriptionAttributes",
      "sns:SetSubscriptionAttributes",
    ]
    resources = [
      "*",
    ]
    condition {
      test     = "StringNotEquals"
      variable = "aws:ResourceTag/ManagedBy"
      values   = ["bootstrap"]
    }
  }
}

data "aws_iam_policy_document" "github_policy_6" {
  statement {
    sid = "SSMParameterStore"
    actions = [
      "ssm:PutParameter",
      "ssm:DeleteParameter",
      "ssm:GetParameter",
      "ssm:GetParameters",
      "ssm:GetParametersByPath",
      "ssm:DescribeParameters",
      "ssm:ListTagsForResource",
      "ssm:AddTagsToResource",
      "ssm:RemoveTagsFromResource",
      "ssm:LabelParameterVersion",
    ]
    resources = [
      "*",
    ]
    condition {
      test     = "StringNotEquals"
      variable = "aws:ResourceTag/ManagedBy"
      values   = ["bootstrap"]
    }
  }
}

data "aws_iam_policy_document" "github_policy_7" {
  statement {
    sid = "ServiceCatalogAppRegistry"
    actions = [
      "servicecatalog:Get*",
      "servicecatalog:List*",
      "servicecatalog:Describe*",
      "servicecatalog:CreateApplication",
      "servicecatalog:DeleteApplication",
      "servicecatalog:UpdateApplication",
      "servicecatalog:AssociateAttributeGroup",
      "servicecatalog:DisassociateAttributeGroup",
      "servicecatalog:ListAssociatedAttributeGroups",
      "servicecatalog:AssociateResource",
      "servicecatalog:DisassociateResource",
      "servicecatalog:TagResource",
      "servicecatalog:UntagResource",
    ]
    resources = [
      "*",
    ]
    condition {
      test     = "StringNotEquals"
      variable = "aws:ResourceTag/ManagedBy"
      values   = ["bootstrap"]
    }
  }
}
