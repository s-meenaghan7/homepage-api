terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6"
    }
  }

  backend "s3" {
    bucket       = "s-meenaghan-terraform-state"
    key          = "homepage-api/terraform.tfstate"
    region       = "us-west-2"
    encrypt      = true
    use_lockfile = true
  }
}

provider "aws" {
  region = "us-west-2"

  default_tags {
    tags = {
      Application    = "homepage"
      awsApplication = "arn:aws:resource-groups:us-west-2:659077917555:group/homepage/06i43yh8tvfxxnq6a2fqezbcxe"
      ManagedBy      = "Terraform"
    }
  }
}
