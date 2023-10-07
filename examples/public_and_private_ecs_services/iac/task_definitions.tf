resource "aws_iam_role" "ServerTaskRole" {
  name = "server-task-role"
  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "ecs-tasks.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_iam_policy" "ServerAssumeRolePolicy" {
  name        = "server-assume-role-policy"
  description = "A policy that allows the ECS task to assume a role in users' accounts"
  policy      = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Effect": "Allow",
      "Resource": "arn:aws:iam::*:role/salami-assumed-role-v0.1-*"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "ServerTaskRolePolicyAttachment" {
  role       = aws_iam_role.ServerTaskRole.name
  policy_arn = aws_iam_policy.ServerAssumeRolePolicy.arn
}

resource "aws_ecs_task_definition" "ServerTaskDefinition" {
  family                = "server"
  cpu                   = "256"
  memory                = "512"
  network_mode          = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  task_role_arn         = aws_iam_role.ServerTaskRole.arn
  execution_role_arn    = aws_iam_role.ServerEcsExecutionRole.arn

  container_definitions = <<DEFINITION
  [
    {
      "name": "${var.server_container_name}",
      "image": "${aws_ecr_repository.ServerRepository.repository_url}:latest",
      "cpu": 256,
      "memory": 512,
      "essential": true,
      "portMappings": [
        {
          "containerPort": ${var.container_port},
          "protocol": "tcp"
        }
      ],
      "environment": [
        {
          "name": "OPENAI_API_KEY",
          "value": "${var.openai_api_key}"
        },
        {
          "name": "ASSUMED_ROLE_SECRET_TOKEN",
          "value": "${var.assumed_role_secret_token}"
        },
        {
          "name": "PYTHON_EXEC_URL",
          "value": "${var.python_exec_local_service_name}.${var.local_dns_namespace_name}:${var.container_port}"
        }
      ],
      "logConfiguration": {
        "logDriver": "awslogs",
        "options": {
          "awslogs-group": "${aws_cloudwatch_log_group.ServerLogGroup.name}",
          "awslogs-region": "${var.aws_region}",
          "awslogs-stream-prefix": "ecs"
        }
      }
    }
  ]
  DEFINITION
}

resource "aws_ecs_task_definition" "PythonExecTaskDefinition" {
  family                = "python-exec"
  cpu                   = "256"
  memory                = "512"
  network_mode          = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  execution_role_arn    = aws_iam_role.PythonExecEcsExecutionRole.arn

  container_definitions = <<DEFINITION
  [
    {
      "name": "${var.python_exec_container_name}",
      "image": "${aws_ecr_repository.PythonExecRepository.repository_url}:latest",
      "cpu": 256,
      "memory": 512,
      "essential": true,
      "portMappings": [
        {
          "containerPort": ${var.container_port},
          "protocol": "tcp"
        }
      ],
      "logConfiguration": {
        "logDriver": "awslogs",
        "options": {
          "awslogs-group": "${aws_cloudwatch_log_group.PythonExecLogGroup.name}",
          "awslogs-region": "${var.aws_region}",
          "awslogs-stream-prefix": "ecs"
        }
      }
    }
  ]
  DEFINITION
}