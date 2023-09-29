resource "aws_ecs_cluster" "EcsCluster" {
  name = "cluster"
}

resource "aws_servicediscovery_private_dns_namespace" "EcsPrivateDnsNamespace" {
  name        = var.local_dns_namespace_name
  description = "Private namespace for ECS cluster"
  vpc = aws_vpc.MainVpc.id
}

resource "aws_ecs_service" "ServerEcsService" {
  name            = "server"
  cluster         = aws_ecs_cluster.EcsCluster.id
  task_definition = aws_ecs_task_definition.ServerTaskDefinition.arn
  desired_count   = 1
  launch_type     = "FARGATE"

  network_configuration {
    assign_public_ip = true
    subnets          = [aws_subnet.PublicSubnetA.id, aws_subnet.PublicSubnetB.id]
    security_groups  = [aws_security_group.ServerEcsSecurityGroup.id]
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.ServerTargetGroup.arn
    container_name   = var.server_container_name
    container_port   = var.container_port
  }

  deployment_controller {
    type = "ECS"
  }

  deployment_circuit_breaker {
    enable   = true
    rollback = true
  }

  wait_for_steady_state = true
}

resource "aws_iam_role_policy_attachment" "PythonExecEcsExecutionRolePolicyAttachment1" {
  role       = aws_iam_role.PythonExecEcsExecutionRole.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

resource "aws_iam_role_policy_attachment" "PythonExecEcsExecutionRolePolicyAttachment2" {
  role       = aws_iam_role.PythonExecEcsExecutionRole.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
}

resource "aws_servicediscovery_service" "PythonExecEcsServiceDiscovery" {
  name = var.python_exec_local_service_name

  dns_config {
    namespace_id = aws_servicediscovery_private_dns_namespace.EcsPrivateDnsNamespace.id
    dns_records {
      ttl  = 10
      type = "A"
    }
  }
}

resource "aws_security_group" "ServerEcsSecurityGroup" {
  name        = "server-ecs-security-group"
  description = "Security group for Server ECS service"
  vpc_id      = aws_vpc.MainVpc.id

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port       = var.container_port
    to_port         = var.container_port
    protocol        = "tcp"
    security_groups = [aws_security_group.ALBSecurityGroup.id]
  }
}

resource "aws_ecs_service" "PythonExecEcsService" {
  name            = "server"
  cluster         = aws_ecs_cluster.EcsCluster.id
  task_definition = aws_ecs_task_definition.PythonExecTaskDefinition.arn
  desired_count   = 1
  launch_type     = "FARGATE"

  network_configuration {
    assign_public_ip = false
    subnets          = [aws_ec2_subnet.PrivateSubnetA.id, aws_ec2_subnet.PrivateSubnetB.id]
    security_groups  = [aws_ec2_security_group.PythonExecEcsSecurityGroup.id]
  }

  deployment_controller {
    type = "ECS"
  }

  deployment_circuit_breaker {
    enable   = true
    rollback = true
  }

  wait_for_steady_state = true
}

resource "aws_security_group" "PythonExecEcsSecurityGroup" {
  name        = "python-exec-ecs-security-group"
  description = "Security group for python exec ECS service"
  vpc_id      = aws_vpc.MainVpc.id

  egress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port       = var.container_port
    to_port         = var.container_port
    protocol        = "tcp"
    security_groups = [aws_security_group.ServerEcsSecurityGroup.id]
  }
}