output "cluster_name" {
  description = "the name of the cluster"
  value       = aws_eks_cluster.main.name
}

output "cluster_endpoint" {
  description = "The cluster endpoint"
  value       = aws_eks_cluster.main.endpoint
}

output "oidc_provider_arn" {
  description = "The ARN for the OIDC provider"
  value       = aws_iam_openid_connect_provider.eks.arn
}

output "oidc_issuer_url" {
  description = "The issuer URL from OIDC"
  value       = aws_iam_openid_connect_provider.eks.url
}
