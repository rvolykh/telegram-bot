data "aws_iam_policy_document" "role_policies" {
  count = length(var.role_policies)

  source_policy_documents = var.role_policies[count.index]
}
