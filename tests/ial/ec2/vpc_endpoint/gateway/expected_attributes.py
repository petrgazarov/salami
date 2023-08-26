resource_type = "ec2.VpcEndpoint"
logical_name = "S3VpcEndpoint"
text = (
    "VPC endpoint type: Gateway\nService name: com.amazonaws.{aws_region}.s3\n"
    "Policy: A JSON policy that allows all principals to perform all actions on resources defined "
    'by the ARN pattern "arn:aws:s3:::prod-{aws_region}-starport-layer-bucket/*". This includes all objects in the specified S3 bucket.'
)
pulumi_aws_imports = ["ec2"]
other_packages = ["json"]
properties = []
uses = ["CumuliVpc", "PrivateRouteTable"]
exports = {}
