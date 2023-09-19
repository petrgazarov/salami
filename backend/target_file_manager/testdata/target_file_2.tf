provider "aws" {
  region = "us-east-2"
}

variable "server_container_name" {
  description = "Server container name"
  type        = string
  default     = "cumuli-server-container"
}