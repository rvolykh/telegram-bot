# GitHub Actions setup

The following GitHub environments are required:
- `approve`
- `sandbox`

## Setup `approve` GitHub environment

1. Create the environment with name - "`approve`".
2. Enable "Required reviewers" option.
3. In "Add up to 5 more reviewers" add at lest yourself.
4. Save protection rules.

This will allow us to review terraform plan before apply.

## Setup `sandbox` GitHub environment

1. Create the environment with name - "`sandbox`".
2. Add the following Environment secrets:
  - `TELEGRAM_BOT_API_TOKEN` - your's Telegram bot API token.
  - `ALERTING_TELEGRAM_CHAT_ID` - separate telegram chat id with your bot, where alerts will be send
                                  (I have not explored enough how to get chat id without actually sending the message to
                                  the bot. Therefor, at first I set empty value here, and once it was deployed, triggered
                                  the bot command via the chat and extracted chat id from the CloudWatch logs. And updated
                                  the secret value in GitHub Actions)
  - `ALERTING_EMAILS` - your email address (the email alerts are only sent when alert to telegram chat id has failed).
                        The value should be in the following format `["your-email@example.com"]`.
  - `AWS_ACCOUNT_ID` - your's AWS Account Number (e.g. `012345678901`)
  - `AWS_REGION` - AWS region where infrastructure should be deployed (e.g. `us-east-1`)
