terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6"
    }
  }
}

provider "aws" {
  region = "us-west-2"
}

resource "aws_dynamodb_table" "visitor_counter" {
  name         = "homepage-visitor-table"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "page_id"

  attribute {
    name = "page_id"
    type = "S"
  }

  # backup protection, minimal cost for small table
  point_in_time_recovery {
    enabled = true
  }

  tags = {
    Name        = "Homepage Visitor Counter"
    Environment = "production"
    ManagedBy   = "terraform"
  }
}
