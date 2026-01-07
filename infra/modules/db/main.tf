resource "aws_dynamodb_table" "this" {
  name = var.name

  billing_mode   = "PROVISIONED"
  read_capacity  = 5
  write_capacity = 5
  hash_key       = var.hash_key_name

  attribute {
    name = var.hash_key_name
    type = var.hash_key_type
  }

  tags = var.tags
}
