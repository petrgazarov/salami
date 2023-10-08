<p align="center">
  <img src="docs/images/salami-icon.svg" alt="Salami Icon" width="100px" height="100px">
</p>

# Salami

Infrastructure As Natural Language

## 🤔 What is this?

Salami is a declarative domain-specific language for cloud infrastructure, centered around natural language descriptions. You can think of Salami as writing documentation for each cloud resource object, and letting the compiler take care of converting that to IaC. The compiler uses LLM to convert natural language to IaC, with Terraform currently as the supported target.

## 🎨 DSL Design

Salami files are mostly comprised of natural language, with several special constructs:

1. **Blocks** - multiline blocks of text that each represent either a resource or a variable; delimited by double newlines.
2. **Constructor functions** - functions that are used to specify the nature of the block; start with an `@` symbol.
3. **Variable references** - references to variables that are defined in the program; delimited by curly braces.
4. **Resource references** - references to resources that are defined in the program; start with a dollar sign.

<p align="center">
  <img src="docs/images/salami-example.png">
  <br>
  <i>Example Salami code with 3 blocks: VPC resource, Security Group resource and the `container_port` variable.</i>
</p>

For more examples, see the `examples` directory.
