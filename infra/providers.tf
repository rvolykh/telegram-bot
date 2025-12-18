terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.0"
    }
    archive = {
      source  = "hashicorp/archive"
      version = "~> 2.0"
    }
    null = {
      source  = "hashicorp/null"
      version = "~> 3.0"
    }
  }
}

provider "aws" {
  alias = "application"

  default_tags {
    tags = {
      "ManagedBy" = "infra"
    }
  }
}

provider "aws" {
  default_tags {
    tags = merge({
      "ManagedBy" = "infra"
    }, aws_servicecatalogappregistry_application.telegram_bot.application_tag)
  }
}
