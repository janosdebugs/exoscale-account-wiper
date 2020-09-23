variable "exoscale_key" {}
variable "exoscale_secret" {}

terraform {
  required_providers {
    exoscale = {
      source  = "terraform-providers/exoscale"
    }
  }
}

provider "exoscale" {
  key = var.exoscale_key
  secret = var.exoscale_secret
}

provider "aws" {
  access_key = var.exoscale_key
  secret_key = var.exoscale_secret
  region = "at-vie-1"
  endpoints {
    s3 = "https://sos-at-vie-1.exo.io"
    s3control = "https://sos-at-vie-1.exo.io"
  }
  skip_credentials_validation = true
  skip_get_ec2_platforms = true
  skip_metadata_api_check = true
  skip_region_validation = true
  skip_requesting_account_id = true
}

