variable "local_dns_namespace_name" {
  type = string
  default = "local"
}

variable "domain_name" {
  type = string
  default = "mydomain.ai"
}

variable "aws_account_id" {
  type = string
  default = "123456789012"
}

variable "aws_region" {
  type = string
  default = "us-west-1"
}

variable "server_container_name" {
  description = "Server container name"
  type        = string
  default     = "server-container"
}

variable "python_exec_container_name" {
  type = string
  default = "python-exec-container"
}

variable "container_port" {
  type    = number
  default = 8000
}

variable "python_exec_local_service_name" {
  type = string
  default = "python-exec"
}