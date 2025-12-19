import {
  to = aws_s3_bucket.state
  id = "telegram-bot-bootstrap-${var.aws_account_id}"
}

import {
  to = aws_s3_bucket_public_access_block.state
  id = "telegram-bot-bootstrap-${var.aws_account_id}"
}

import {
  to = aws_s3_bucket_versioning.state
  id = "telegram-bot-bootstrap-${var.aws_account_id}"
}

import {
  to = aws_s3_bucket_lifecycle_configuration.state
  id = "telegram-bot-bootstrap-${var.aws_account_id}"
}

import {
  to = aws_s3_bucket_notification.state
  id = "telegram-bot-bootstrap-${var.aws_account_id}"
}

import {
  to = aws_s3_bucket_server_side_encryption_configuration.state
  id = "telegram-bot-bootstrap-${var.aws_account_id}"
}

import {
  to = aws_iam_openid_connect_provider.github
  id = "arn:aws:iam::${var.aws_account_id}:oidc-provider/token.actions.githubusercontent.com"
}

import {
  to = aws_iam_role.github
  id = "github-telegram-bot"
}

import {
  to = aws_iam_policy.github_policy_1
  id = "arn:aws:iam::${var.aws_account_id}:policy/github-telegram-bot-1"
}

import {
  to = aws_iam_role_policy_attachment.github_policy_1
  id = "github-telegram-bot/arn:aws:iam::${var.aws_account_id}:policy/github-telegram-bot-1"
}

import {
  to = aws_iam_policy.github_policy_2
  id = "arn:aws:iam::${var.aws_account_id}:policy/github-telegram-bot-2"
}

import {
  to = aws_iam_role_policy_attachment.github_policy_2
  id = "github-telegram-bot/arn:aws:iam::${var.aws_account_id}:policy/github-telegram-bot-2"
}

import {
  to = aws_iam_policy.github_policy_3
  id = "arn:aws:iam::${var.aws_account_id}:policy/github-telegram-bot-3"
}

import {
  to = aws_iam_role_policy_attachment.github_policy_3
  id = "github-telegram-bot/arn:aws:iam::${var.aws_account_id}:policy/github-telegram-bot-3"
}


import {
  to = aws_iam_policy.github_policy_4
  id = "arn:aws:iam::${var.aws_account_id}:policy/github-telegram-bot-4"
}

import {
  to = aws_iam_role_policy_attachment.github_policy_4
  id = "github-telegram-bot/arn:aws:iam::${var.aws_account_id}:policy/github-telegram-bot-4"
}

import {
  to = aws_iam_policy.github_policy_5
  id = "arn:aws:iam::${var.aws_account_id}:policy/github-telegram-bot-5"
}

import {
  to = aws_iam_role_policy_attachment.github_policy_5
  id = "github-telegram-bot/arn:aws:iam::${var.aws_account_id}:policy/github-telegram-bot-5"
}

import {
  to = aws_iam_policy.github_policy_6
  id = "arn:aws:iam::${var.aws_account_id}:policy/github-telegram-bot-6"
}

import {
  to = aws_iam_role_policy_attachment.github_policy_6
  id = "github-telegram-bot/arn:aws:iam::${var.aws_account_id}:policy/github-telegram-bot-6"
}

import {
  to = aws_iam_policy.github_policy_7
  id = "arn:aws:iam::${var.aws_account_id}:policy/github-telegram-bot-7"
}

import {
  to = aws_iam_role_policy_attachment.github_policy_7
  id = "github-telegram-bot/arn:aws:iam::${var.aws_account_id}:policy/github-telegram-bot-7"
}
