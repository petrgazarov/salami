# Salami

<img src="salami-icon.svg" alt="Salami Icon">

## Development Setup

See `Makefile` for available commands.

## Object Types

- **Resource**
  Supported decorators:
  ```python
  @exports(key1: value1, key2: value2)
  @uses(OtherResource1, OtherResource2)
  ```
  Required fields - must be the first two fields:
  ```
  Resource type: aws.[Service].[ResourceType]
  Logical name: UniqueLogicalNameForResource
  ```
  Example source:
  ```
  @exports(name: exported-name)
  @uses(CumuliEcsCluster, CumuliServerTaskDefinition, PublicSubnetA, PublicSubnetB, CumuliServerEcsSecurityGroup, CumuliServerTargetGroup)
  Resource type: aws.ecs.Service
  Logical name: CumuliServerEcsService
  Name: cumuli-server
  Desired count: 1
  Launch type: FARGATE
  ---
  Network configuration:
    Assigned public IP.
    Subnets: PublicSubnetA and PublicSubnetB.
    Security group: CumuliServerEcsSecurityGroup.
  Load balancers:
    Target group: CumuliServerTargetGroup.
    Container name: {server_container_name}.
    Port: {container_port}.
  Deployment:
    ECS type deployment controller.
    Deployment circuit breaker: enabled with rollback.
    Wait for steady state: True
  ```
  Example target code (pulumi python):
  ```python
  import pulumi
  import pulumi_aws as aws

  cumuli_server_ecs_service = aws.ecs.Service(
      "CumuliServerEcsService",
      name="cumuli-server",
      cluster=cumuli_ecs_cluster.arn,
      task_definition=cumuli_server_task_definition.arn,
      desired_count=1,
      launch_type="FARGATE",
      network_configuration={
          "assign_public_ip": True,
          "subnets": [public_subnet_a.id, public_subnet_b.id],
          "security_groups": [cumuli_server_ecs_security_group.id],
      },
      load_balancers=[
          {
              "target_group_arn": cumuli_server_target_group.arn,
              "container_name": server_container_name,
              "container_port": container_port,
          }
      ],
      deployment_controller={"type": "ECS"},
      deployment_circuit_breaker={"enable": True, "rollback": True},
      wait_for_steady_state=True,
  )

  pulumi.export("ecs-service-name", cumuli_server_ecs_service.name)
  ```
- **Variable**
  Supported decorators:
  ```python
  @variable(type)
  ```
  Fields:
  ```
  Description: Description string (optional)
  Name: variable-name
  Value: default-value (optional)
  ```
  Source example:
  ```
  @variable(string)
  Description: Server container name
  Name: server_container_name
  Value: cumuli-server-container
  ```
