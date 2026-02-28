variable "cluster_name" {
  type        = string
  description = "The name of the EKS cluster."
}

variable "subnet_ids" {
  type        = list(string)
  description = "The subnet IDs for public and private subnets"
}

variable "node_instance_types" {
  type        = list(string)
  description = "The instance type of the node group"
  default     = ["t3.medium"]
}

variable "node_desired_size" {
  type        = number
  description = "The desired size of the node group"
  default     = 1
}

variable "node_min_size" {
  type        = number
  description = "The minimum size of the node group"
  default     = 1

}

variable "node_max_size" {
  type        = number
  description = "The maxiumum size of the node group"
  default     = 2
}

variable "gha_role_arn" {
  type        = string
  description = "The arn of the github actions role"
}
