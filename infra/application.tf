resource "aws_servicecatalogappregistry_application" "telegram_bot" {
  provider    = aws.application
  name        = "telegram-bot"
  description = "Telegram Bot"
}
