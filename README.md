# Salami

<p align="center">
  <img src="salami-icon.svg" alt="Salami Icon" width="100px" height="100px">
</p>

Salami is a tool for managing cloud infrastructure as natural language. Salami provides a minimalist declarative DSL, a compiler, and a CLI utility.
The compiler is designed to be target and LLM-agnostic. Initially, only Terraform target and GPT-4 LLM are supported.

## DSL Design

Salami programs are _mostly_ natural language, with several special constructs to make source more expressive and concise,
and provide additional information to the compiler. The special constructs are:

1. **Blocks** - multiline blocks of text that each represent either a `Resource` or a `Variable`.
2. **Decorators** - python-inspired functions that are used to specify special information about the block, such as dependencies and exports.
3. **Fields** - colon-separated key-value pairs that have special meaning within a Salami block.
4. **Variables** - variables can be referenced anywhere within natural language by using curly braces.

### Resource

Supported decorators:

| Decorator | Function signature                                | Arguments format                   | Example                                 | Required |
| --------- | ------------------------------------------------- | ---------------------------------- | --------------------------------------- | -------- |
| @exports  | @exports(property: export-name)                   | Comma-separated list of key-values | @exports(name: ecs-service-name)        | No       |
| @uses     | @uses(Resource1LogicalName, Resource2LogicalName) | Comma-separated list of strings    | @uses(EcsCluster, ServerTaskDefinition) | No       |

Fields:

| Field         | Value format                                                | Example value   | Required (y/n) |
| ------------- | ----------------------------------------------------------- | --------------- | -------------- |
| Resource type | A string in the format: [provider].[service].[ResourceType] | aws.ecs.Service | Yes            |
| Logical name  | Alphanumeric string, starts with a letter                   | EcsCluster      | Yes            |

Source block example:

```
@exports(name: ecs-service-name)
@uses(EcsCluster, ServerTaskDefinition, PublicSubnetA, PublicSubnetB, ServerEcsSecurityGroup, ServerTargetGroup)
Resource type: aws.ecs.Service
Logical name: ServerEcsService
Name: api-server
Desired count: 1
Launch type: FARGATE
---
Network configuration:
  Assigned public IP.
  Subnets: PublicSubnetA and PublicSubnetB.
  Security group: ServerEcsSecurityGroup.
Load balancers:
  Target group: ServerTargetGroup.
  Container name: {server_container_name}.
  Port: {container_port}.
Deployment:
  ECS type deployment controller.
  Deployment circuit breaker: enabled with rollback.
  Wait for steady state: True
```

### Variable

Supported decorators:

| Decorator | Function signature | Arguments format                        | Example           | Required |
| --------- | ------------------ | --------------------------------------- | ----------------- | -------- |
| @variable | @variable(type)    | Must be one of: string, number, boolean | @variable(string) | Yes      |

Fields:

| Field       | Value format                             | Example value                             | Required (y/n) |
| ----------- | ---------------------------------------- | ----------------------------------------- | -------------- |
| Name        | Alphanumeric string, underscores allowed | server_container_name                     | Yes            |
| Description | String                                   | Container name for the server ECS service | No             |
| Value       | String, number, or boolean               | server-container                          | No             |

Source block example:

```
@variable(string)
Description: Container name for the server ECS service
Name: server_container_name
Value: server-container
```
