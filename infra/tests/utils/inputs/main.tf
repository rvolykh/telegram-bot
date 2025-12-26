terraform {
  required_providers {
    random = {
      source  = "hashicorp/random"
      version = "~> 3.7"
    }
    http = {
      source  = "hashicorp/http"
      version = "~> 3.5"
    }
  }
}

# Resources

resource "random_pet" "prefix" {
  prefix = "test"
  length = 1
}

resource "random_uuid" "token" {}

# Data sources

data "http" "my_ip" {
  url = "https://ifconfig.me/ip"
}

# Outputs

output "prefix" {
  value = "${random_pet.prefix.id}-"
}

output "token" {
  value = "test-token-${random_uuid.token.result}"
}

output "my_ip" {
  value = chomp(data.http.my_ip.response_body)
}
