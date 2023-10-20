# Example: simple S3 bucket

This example creates an S3 bucket.

## Running the example

### Prerequisites

To run this example, you need:

- `terraform` installed
- `salami` installed (follow installation instructions in the [README](../../README.md))
- AWS credentials (optional, if you want to deploy the infrastructure)

### Steps

1. Clone this repository
2. `cd` into the `examples/simple_s3_bucket` directory
3. Run `salami compile` to compile the Salami descriptions into Terraform code
4. Optionally, `cd` into the `examples/simple_s3_bucket/terraform` directory and run the usual `terraform init`, `terraform plan` and `terraform apply` commands to deploy the infrastructure. Make sure to pass the AWS credentials to Terraform (Salami does not generate the `provider` block for you).

### FYI

1. Note that `salami compile` will examine the salami-lock.toml file and the source `.sami` files, and determine which Salami objects have changed since the last compilation. To force a complete recompilation, delete the `salami-lock.toml` file. Or, you can change source `.sami` files and `salami compile` will recompile only the changed objects.

2. Occassionally, OpenAI API delays responses significantly. If `salami compile` is stuck for a long time, try again later. The total compilation time varies significantly depending on the number of objects and their complexity. This project takes me about 10 seconds to compile from scratch.