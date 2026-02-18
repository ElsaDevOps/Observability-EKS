resource "aws_iam_role" "irsa_role" {
  name               = var.role_name
  assume_role_policy = data.aws_iam_policy_document.trust.json
}

data "aws_iam_policy_document" "trust" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRoleWithWebIdentity"]

    principals {
      type        = "Federated"
      identifiers = [var.oidc_arn]
    }

    condition {
      test     = "StringEquals"
      variable = "${replace(var.issuer_url, "https://", "")}:sub"
      values   = ["system:serviceaccount:${var.namespace}:${var.service_account_name}"]
    }

    condition {
      test     = "StringEquals"
      variable = "${replace(var.issuer_url, "https://", "")}:aud"
      values   = ["sts.amazonaws.com"]
    }
  }
}


resource "aws_iam_policy" "permissions_policy" {
  name   = "${var.role_name}-policy"
  policy = var.policy_json
}

resource "aws_iam_role_policy_attachment" "irsa_eks" {
  role       = aws_iam_role.irsa_role.name
  policy_arn = aws_iam_policy.permissions_policy.arn
}
