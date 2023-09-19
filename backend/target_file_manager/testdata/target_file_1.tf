resource "aws_ecs_service" "ServerEcsService" {
  name            = "server"
  cluster         = aws_ecs_cluster.EcsCluster.name
  task_definition = aws_ecs_task_definition.ServerTaskDefinition.arn
  desired_count   = 1
  launch_type     = "FARGATE"

  network_configuration {
    assign_public_ip = true
    subnets = [
      aws_subnet.PublicSubnetA.id,
      aws_subnet.PublicSubnetB.id
    ]
    security_groups = [aws_security_group.ServerEcsSecurityGroup.id]
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

output "server-ecs-service-name" {
  value = aws_ecs_service.ServerEcsService.name
}