output "api_gateway_url" {
  description = "URL of the Telegram bot API Gateway"
  value       = module.telegram_bot_api.webhook_url
  sensitive   = true
}
