data "aws_iam_policy_document" "read_only" {
  version = "2012-10-17"
  statement {
    effect = "Allow"
    actions = [
      "dynamodb:GetItem",
      "dynamodb:Query",
      "dynamodb:Scan",
    ]
    resources = [
      aws_dynamodb_table.this.arn,
    ]
  }
}

data "aws_iam_policy_document" "write_only" {
  version = "2012-10-17"
  statement {
    effect = "Allow"
    actions = [
      "dynamodb:PutItem",
      "dynamodb:UpdateItem",
      "dynamodb:DeleteItem",
    ]
    resources = [
      aws_dynamodb_table.this.arn,
    ]
  }
}

data "aws_iam_policy_document" "read_write" {
  source_policy_documents = [
    data.aws_iam_policy_document.read_only.json,
    data.aws_iam_policy_document.write_only.json,
  ]
}
