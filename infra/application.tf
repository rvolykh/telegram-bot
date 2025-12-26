resource "aws_servicecatalogappregistry_application" "telegram_bot" {
  provider    = aws.application
  name        = "${var.prefix}telegram-bot"
  description = "${var.prefix}Telegram Bot"
}
