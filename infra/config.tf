module "telegram_bot_api_token" {
  source = "./modules/kv"

  name  = "/${var.prefix}telegram/bot/api_token"
  value = var.telegram_bot_api_token
}

module "telegram_bot_cache_poweron" {
  source = "./modules/kv"

  name  = "/${var.prefix}telegram/bot/cache/poweron"
  value = "none"
}
