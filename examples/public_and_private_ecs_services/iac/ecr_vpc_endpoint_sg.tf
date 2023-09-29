resource "aws_security_group" "EcrVpcEndpointSG" {
  name        = "ecr-vpc-endpoint-sg"
  description = "Security Group for ECR VPC Endpoint"
  vpc_id      = aws_vpc.MainVpc.id

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port       = 443
    to_port         = 443
    protocol        = "tcp"
    security_groups = [aws_security_group.ServerEcsSecurityGroup.id, aws_security_group.PythonExecEcsSecurityGroup.id]
  }
}