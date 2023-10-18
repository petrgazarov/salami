variable "local_dns_namespace_name" {
  type = string
  default = "local"
}

variable "aws_account_id" {
  type = string
}

variable "aws_region" {
  type = string
  default = "us-west-2"
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

variable "openai_api_key" {
  type = string
}

variable "assumed_role_secret_token" {
  type = string
}