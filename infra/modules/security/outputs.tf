output "role_arn" {
  description = "The IAM role ARN"
  value       = aws_iam_role.irsa_role.arn
}
