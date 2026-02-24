terraform {
  backend "s3" {
    bucket       = "terraform-state-kay8s"
    key          = "infra/terraform.tfstate"
    region       = "eu-west-2"
    encrypt      = true
    use_lockfile = true
  }
}
