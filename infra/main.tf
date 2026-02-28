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


module "eks" {
  source = "./modules/eks"

  cluster_name        = local.cluster_name
  subnet_ids          = values(module.vpc.private_subnet_id)
  node_instance_types = var.node_instance_types
  node_desired_size   = var.node_desired_size
  node_min_size       = var.node_min_size
  node_max_size       = var.node_max_size
  gha_role_arn        = var.gha_role_arn
}

module "aws_lb_controller_irsa" {
  source = "./modules/security"

  namespace            = "kube-system"
  service_account_name = "aws-load-balancer-controller"
  policy_json          = file("${path.module}/policies/aws-lb-controller.json")
  oidc_arn             = module.eks.oidc_provider_arn
  issuer_url           = module.eks.oidc_issuer_url
  role_name            = "lb-controller-irsa"
}

module "cert_manager" {
  source = "./modules/security"

  namespace            = "cert-manager"
  service_account_name = "cert-manager"
  policy_json          = file("${path.module}/policies/cert-manager.json")
  oidc_arn             = module.eks.oidc_provider_arn
  issuer_url           = module.eks.oidc_issuer_url
  role_name            = "cert-manager-irsa"
}

module "external_dns" {
  source = "./modules/security"

  namespace            = "external-dns"
  service_account_name = "external-dns"
  policy_json          = file("${path.module}/policies/external-dns.json")
  oidc_arn             = module.eks.oidc_provider_arn
  issuer_url           = module.eks.oidc_issuer_url
  role_name            = "external-dns-irsa"
}

resource "aws_ssm_parameter" "lb_controller_role_arn" {
  name   = "/${var.environment}/platform/lb-controller-role-arn"
  type   = "SecureString"
  value  = module.aws_lb_controller_irsa.role_arn
  key_id = aws_kms_key.encryption_key.arn
}

resource "aws_ssm_parameter" "external_dns_role_arn" {
  name   = "/${var.environment}/platform/external-dns-role-arn"
  type   = "SecureString"
  value  = module.external_dns.role_arn
  key_id = aws_kms_key.encryption_key.arn
}

resource "aws_ssm_parameter" "cert_manager_role_arn" {
  name   = "/${var.environment}/platform/cert-manager-role-arn"
  type   = "SecureString"
  value  = module.cert_manager.role_arn
  key_id = aws_kms_key.encryption_key.arn
}

resource "aws_ssm_parameter" "vpc_id" {
  name   = "/${var.environment}/platform/vpc-id"
  type   = "SecureString"
  value  = module.vpc.vpc_id
  key_id = aws_kms_key.encryption_key.id
}

resource "aws_kms_key" "encryption_key" {
  # checkov:skip=CKV2_AWS_64: Default key policy grants root account access, IAM policies control usage
  description             = "A symmetric encryption KMS key"
  enable_key_rotation     = true
  deletion_window_in_days = 20
}
