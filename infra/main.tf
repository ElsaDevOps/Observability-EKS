module "vpc" {
  source = "./modules/vpc"

  cidr_blockvpc           = local.network_config.cidr_vpc
  cidr_public_subnet_web  = local.network_config.cidr_public_subnet_web
  cidr_private_subnet_app = local.network_config.cidr_private_subnet_app
  availability_zones      = local.network_config.availability_zones
  project_name            = var.project_name

  public_subnet_tags = {
    "kubernetes.io/cluster/${local.cluster_name}" = "owned"
    "kubernetes.io/role/elb"                      = "1"
  }

  private_subnet_tags = {
    "kubernetes.io/cluster/${local.cluster_name}" = "owned"
    "kubernetes.io/role/internal-elb"             = "1"
  }

}