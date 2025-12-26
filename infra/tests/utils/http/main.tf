terraform {
  required_providers {
    http = {
      source  = "hashicorp/http"
      version = "~> 3.5"
    }
  }
}

variable "url" {
  type = string
}

data "http" "request" {
  url    = var.url
  method = "POST"

  request_body = jsonencode({
    update_id = 42,
    message = {
      message_id = 42
      date       = 1767225600
      text       = "/test"
      chat = {
        id       = 42
        type     = "private"
        username = "test"
      }
      from = {
        id         = 42
        username   = "test"
        first_name = "Test"
        last_name  = "Test"
        is_bot     = true
      }
      entities = [
        {
          type   = "bot_command"
          offset = 0
          length = 5
        }
      ]
    }
  })
}
