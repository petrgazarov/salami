resource_type = "ec2.VpcEndpoint"
logical_name = "EcrDkrVpcEndpoint"
text = (
    "VPC endpoint type: Interface\nService name: com.amazonaws.{aws_region}.ecr.dkr\nPrivate DNS enabled: True\n"
    "Security group: EcrVpcEndpointSG\nPolicy: A JSON policy that allows two AWS principals - CumuliEcsExecutionRole "
    "and PythonExecEcsExecutionRole - to perform three actions on ECR resources. The actions are "
    '"ecr:BatchCheckLayerAvailability", "ecr:GetDownloadUrlForLayer", and "ecr:BatchGetImage". The resources on which '
    'these actions can be performed are defined by the ARN pattern "arn:aws:ecr:{aws_region}:{aws_account_id}:repository/*", '
    "which includes all repositories in the specified AWS region and account."
)
pulumi_aws_imports = ["ec2"]
other_packages = ["json"]
properties = []
uses = [
    "CumuliVpc",
    "PrivateSubnetA",
    "PrivateSubnetB",
    "EcrVpcEndpointSG",
    "CumuliEcsExecutionRole",
    "PythonExecEcsExecutionRole",
]
exports = {}
