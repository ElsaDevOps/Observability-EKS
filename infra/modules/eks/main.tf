resource "aws_eks_cluster" "main" {
  # checkov:skip=CKV_AWS_39: want public and private enabled
  # checkov:skip=CKV_AWS_38: github actions runners don't have fixed IPs
  name = var.cluster_name

  access_config {
    authentication_mode = "API"
  }

  role_arn = aws_iam_role.cluster.arn
  version  = "1.31"

 vpc_config {
    endpoint_private_access = true
    endpoint_public_access  = true

    subnet_ids = var.subnet_ids

  }

  encryption_config {
    provider {
      key_arn = aws_kms_key.my_key.arn
    }
    resources = ["secrets"]
  }

  enabled_cluster_log_types = [
    "api",
    "audit",
    "authenticator",
    "controllerManager",
    "scheduler"
  ]

  depends_on = [ aws_iam_role_policy_attachment.cluster_AmazonEKSClusterPolicy, ]
}

resource "aws_kms_key" "my_key" {
  # checkov:skip=CKV2_AWS_64: Default key policy grants root account access, IAM policies control usage
  description             = "A symmetric encryption KMS key"
  enable_key_rotation     = true
  deletion_window_in_days = 20
}

resource "aws_iam_role" "cluster" {
  name = "${var.cluster_name}-cluster-role"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "sts:AssumeRole",
          "sts:TagSession"
        ]
        Effect = "Allow"
        Principal = {
          Service = "eks.amazonaws.com"
        }
      },
    ]
  })
}

resource "aws_iam_role_policy_attachment" "cluster_AmazonEKSClusterPolicy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"
  role       = aws_iam_role.cluster.name
}
