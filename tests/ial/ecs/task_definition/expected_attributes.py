resource_type = "ecs.TaskDefinition"
logical_name = "CumuliServerTaskDefinition"
text = (
    "Name: cumuli-server-task-definition\nFamily: cumuli-server\nCpu: 256\nMemory: 512\n"
    "Network mode: awsvpc\nTask role: CumuliServerTaskRole\n"
    "Execution role: CumuliEcsExecutionRole\nContainer definition:\n"
    "  Name: {server_container_name}\n  Image: CumuliServerRepository url with the 'latest' tag\n"
    "  Memory: 512\n  Cpu: 256\n  Essential: True\n  Port mappings: Only container port {container_port} is mapped to port {container_port}. Host port is auto-assigned by ECS.\n"
    "  Environment:\n    OPENAI_API_KEY: taken from environment\n    ASSUMED_ROLE_SECRET_TOKEN: taken from environment\n"
    "    PYTHON_EXEC_URL: The value is enterpolated as {python_exec_local_service_name}.{local_dns_namespace_name}:{container_port}.\n"
    "  Log configuration: awslogs log driver, CumuliServerLogGroup log group name, and {aws_region} AWS region. The stream prefix is set to ecs."
)
pulumi_aws_imports = ["ecs"]
other_packages = ["os", "json"]
properties = ["Requires FARGATE compatibility"]
uses = [
    "CumuliServerTaskRole",
    "CumuliEcsExecutionRole",
    "CumuliServerRepository",
    "CumuliServerLogGroup",
]
exports = {}
