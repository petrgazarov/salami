resource "aws_vpc_endpoint" "EcrDkrVpcEndpoint" {
  vpc_id            = aws_vpc.MainVpc.id
  service_name      = "com.amazonaws.${var.aws_region}.ecr.dkr"
  vpc_endpoint_type = "Interface"
  subnet_ids        = [aws_subnet.PrivateSubnetA.id, aws_subnet.PrivateSubnetB.id]
  private_dns_enabled = true
  security_group_ids = [aws_security_group.EcrVpcEndpointSG.id]

  policy = <<POLICY
  {
    "Version": "2012-10-17",
    "Statement": [
      {
        "Effect": "Allow",
        "Principal": {
          "AWS": ["${aws_iam_role.ServerEcsExecutionRole.arn}", "${aws_iam_role.PythonExecEcsExecutionRole.arn}"]
        },
        "Action": [
          "ecr:BatchCheckLayerAvailability",
          "ecr:GetDownloadUrlForLayer",
          "ecr:BatchGetImage"
        ],
        "Resource": "*"
      }
    ]
  }
  POLICY
}

resource "aws_vpc_endpoint" "EcrApiVpcEndpoint" {
  vpc_id            = aws_vpc.MainVpc.id
  service_name      = "com.amazonaws.${var.aws_region}.ecr.api"
  vpc_endpoint_type = "Interface"
  subnet_ids        = [aws_subnet.PrivateSubnetA.id, aws_subnet.PrivateSubnetB.id]
  private_dns_enabled = true
  security_group_ids  = [aws_security_group.EcrVpcEndpointSG.id]

  policy = <<POLICY
  {
    "Version": "2012-10-17",
    "Statement": [
      {
        "Effect": "Allow",
        "Principal": {
          "AWS": ["${aws_iam_role.ServerEcsExecutionRole.arn}", "${aws_iam_role.PythonExecEcsExecutionRole.arn}"]
        },
        "Action": [
          "ecr:GetAuthorizationToken",
          "ecr:BatchCheckLayerAvailability",
          "ecr:GetDownloadUrlForLayer",
          "ecr:BatchGetImage"
        ],
        "Resource": "*"
      }
    ]
  }
  POLICY
}

resource "aws_vpc_endpoint" "S3VpcEndpoint" {
  vpc_id              = aws_vpc.MainVpc.id
  service_name        = "com.amazonaws.${var.aws_region}.s3"
  vpc_endpoint_type   = "Gateway"
  route_table_ids     = [aws_route_table.PrivateRouteTable.id]

  policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": "*",
      "Action": "*",
      "Resource": "arn:aws:s3:::prod-${var.aws_region}-starport-layer-bucket/*"
    }
  ]
}
POLICY
}

resource "aws_security_group" "CloudWatchLogsVpcEndpointSG" {
  name        = "cloudwatch-logs-vpc-endpoint-sg"
  description = "Security Group for CloudWatch Logs VPC Endpoint"
  vpc_id      = aws_vpc.MainVpc.id

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_vpc_endpoint" "CloudWatchLogsVpcEndpoint" {
  vpc_id            = aws_vpc.MainVpc.id
  service_name      = "com.amazonaws.${var.aws_region}.logs"
  vpc_endpoint_type = "Interface"
  subnet_ids        = [aws_subnet.PrivateSubnetA.id, aws_subnet.PrivateSubnetB.id]
  private_dns_enabled = true
  security_group_ids = [aws_security_group.CloudWatchLogsVpcEndpointSG.id]

  policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {"AWS": "*"},
      "Action": [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents",
        "logs:DescribeLogStreams"
      ],
      "Resource": "*"
    }
  ]
}
POLICY
}