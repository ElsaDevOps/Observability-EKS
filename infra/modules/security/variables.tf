variable "oidc_arn" {
  type        = string
  description = "The OIDC ARN"
}

variable "issuer_url" {
  type        = string
  description = "The OIDC issuer URL"
}

variable "role_name" {
  type        = string
  description = "The name of the IRSA role"
}

variable "policy_json" {
  type        = string
  description = "The permissions policy for roles"
}

variable "namespace" {
  type        = string
  description = "The namespace for the roles"
}

variable "service_account_name" {
  type        = string
  description = "The name of the ServiceAccount for each role"
}
