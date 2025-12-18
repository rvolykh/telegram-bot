resource "aws_iam_openid_connect_provider" "github" {
  url = "https://token.actions.githubusercontent.com"

  client_id_list = [
    "sts.amazonaws.com",
  ]
}

resource "aws_iam_role" "github" {
  name        = "github-telegram-bot"
  path        = "/bootstrap/"
  description = "Role for GitHub Actions to access AWS resources"

  assume_role_policy = data.aws_iam_policy_document.github_trust_policy.json
}

resource "aws_iam_policy" "github_policy_1" {
  name   = "github-telegram-bot-1"
  policy = data.aws_iam_policy_document.github_policy_1.json
}

resource "aws_iam_role_policy_attachment" "github_policy_1" {
  role       = aws_iam_role.github.name
  policy_arn = aws_iam_policy.github_policy_1.arn
}

resource "aws_iam_policy" "github_policy_2" {
  name   = "github-telegram-bot-2"
  policy = data.aws_iam_policy_document.github_policy_2.json
}

resource "aws_iam_role_policy_attachment" "github_policy_2" {
  role       = aws_iam_role.github.name
  policy_arn = aws_iam_policy.github_policy_2.arn
}

resource "aws_iam_policy" "github_policy_3" {
  name   = "github-telegram-bot-3"
  policy = data.aws_iam_policy_document.github_policy_3.json
}

resource "aws_iam_role_policy_attachment" "github_policy_3" {
  role       = aws_iam_role.github.name
  policy_arn = aws_iam_policy.github_policy_3.arn
}

resource "aws_iam_policy" "github_policy_4" {
  name   = "github-telegram-bot-4"
  policy = data.aws_iam_policy_document.github_policy_4.json
}

resource "aws_iam_role_policy_attachment" "github_policy_4" {
  role       = aws_iam_role.github.name
  policy_arn = aws_iam_policy.github_policy_4.arn
}

resource "aws_iam_policy" "github_policy_5" {
  name   = "github-telegram-bot-5"
  policy = data.aws_iam_policy_document.github_policy_5.json
}

resource "aws_iam_role_policy_attachment" "github_policy_5" {
  role       = aws_iam_role.github.name
  policy_arn = aws_iam_policy.github_policy_5.arn
}

resource "aws_iam_policy" "github_policy_6" {
  name   = "github-telegram-bot-6"
  policy = data.aws_iam_policy_document.github_policy_6.json
}

resource "aws_iam_role_policy_attachment" "github_policy_6" {
  role       = aws_iam_role.github.name
  policy_arn = aws_iam_policy.github_policy_6.arn
}

resource "aws_iam_policy" "github_policy_7" {
  name   = "github-telegram-bot-7"
  policy = data.aws_iam_policy_document.github_policy_7.json
}

resource "aws_iam_role_policy_attachment" "github_policy_7" {
  role       = aws_iam_role.github.name
  policy_arn = aws_iam_policy.github_policy_7.arn
}
