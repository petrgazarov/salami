resource "aws_ecr_repository" "ServerRepository" {
  name                 = "server"
  image_tag_mutability = "MUTABLE"
}

resource "aws_ecr_lifecycle_policy" "ServerRepoLifecyclePolicy" {
  repository = aws_ecr_repository.ServerRepository.name

  policy = <<EOF
{
  "rules": [
    {
      "rulePriority": 1,
      "description": "Expire images beyond count 10",
      "selection": {
        "tagStatus": "untagged",
        "countType": "imageCountMoreThan",
        "countNumber": 10
      },
      "action": {
        "type": "expire"
      }
    }
  ]
}
EOF
}

resource "aws_iam_role" "ServerEcsExecutionRole" {
  name = "server-ecs-execution-role"
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

resource "aws_ecr_repository_policy" "ServerRepositoryPolicy" {
  repository = aws_ecr_repository.ServerRepository.name

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "AWS": "${aws_iam_role.ServerEcsExecutionRole.arn}"
      },
      "Action": [
        "ecr:GetDownloadUrlForLayer",
        "ecr:BatchGetImage",
        "ecr:BatchCheckLayerAvailability"
      ]
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "ServerRepositoryPolicyAttachment1" {
  role       = aws_iam_role.ServerEcsExecutionRole.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

resource "aws_iam_role_policy_attachment" "ServerRepositoryPolicyAttachment2" {
  role       = aws_iam_role.ServerEcsExecutionRole.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
}

resource "aws_ecr_repository" "PythonExecRepository" {
  name                 = "python-exec"
  image_tag_mutability = "MUTABLE"
}

resource "aws_ecr_lifecycle_policy" "PythonExecRepoLifecyclePolicy" {
  repository = aws_ecr_repository.PythonExecRepository.name

  policy = <<EOF
{
  "rules": [
    {
      "rulePriority": 1,
      "description": "Expire images beyond count 10",
      "selection": {
        "tagStatus": "untagged",
        "countType": "imageCountMoreThan",
        "countNumber": 10
      },
      "action": {
        "type": "expire"
      }
    }
  ]
}
EOF
}

resource "aws_iam_role" "PythonExecEcsExecutionRole" {
  name = "python-exec-ecs-execution-role"
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

resource "aws_ecr_repository_policy" "PythonExecRepositoryPolicy" {
  repository = aws_ecr_repository.PythonExecRepository.name

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "AWS": "${aws_iam_role.PythonExecEcsExecutionRole.arn}"
      },
      "Action": [
        "ecr:GetDownloadUrlForLayer",
        "ecr:BatchGetImage",
        "ecr:BatchCheckLayerAvailability"
      ]
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "PythonExecEcsExecutionRolePolicyAttachment1" {
  role       = aws_iam_role.PythonExecEcsExecutionRole.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

resource "aws_iam_role_policy_attachment" "PythonExecEcsExecutionRolePolicyAttachment2" {
  role       = aws_iam_role.PythonExecEcsExecutionRole.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
}