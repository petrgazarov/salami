# Example: simple S3 bucket

Creates an S3 bucket.

## Running the example

### Prerequisites

To run this example, you need:

- `terraform` installed
- `salami` installed (follow installation instructions in the [README](../../README.md))
- AWS credentials (optional, to deploy the infrastructure)

### Steps

1. Clone this repository
2. `cd` into the `examples/simple_s3_bucket` directory
3. Run `salami compile` to run the compiler
4. Optionally, `cd` into the `examples/simple_s3_bucket/terraform` directory and run the `terraform init` and `terraform apply` commands to deploy to AWS.

### FYI

1. `salami compile` command examines `salami-lock.toml` and the source `.sami` files to determine the changeset. To force a complete regeneration, delete the `salami-lock.toml` file and rerun the compiler.

2. If timeout error is raised, try setting `compiler.llm.max_concurrent` config to a lower number. This slows down the compilation process, but reduces the likelihood of timeouts from OpenAI.
