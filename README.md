# Salami

<p align="center">
  <img src="docs/images/salami-icon.svg" alt="Salami Icon" width="100px" height="100px">
</p>

Salami is a declarative, domain-specific language for cloud infrastructure that is based on natural language. The compiler uses LLM to convert natural language into structured code. Currently, the supported target is Terraform.

## DSL Design

Salami programs have several special constructs:

1. **Blocks** - multiline blocks of text that each represent either a `resource` or a `variable`.
2. **Constructor functions** - functional expressions that are used to specify the nature of the block.
3. **Variable references** - references to variables that are defined in the program; delimited by curly braces.
4. **Resource references** - references to resources that are defined in the program; start with a dollar sign.

<p align="center">
  <img height="200" src="docs/images/salami-example.png">
  <br>
  <i>Example Salami code with 3 blocks: VPC resource, Security Group resource and a variable.</i>
</p>
