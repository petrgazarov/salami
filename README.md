<p align="center">
  <img src="docs/images/salami-icon.svg" alt="Salami Icon" width="100px" height="100px">
</p>

# Salami

Infrastructure As Natural Language

## 🤔 What is this?

Salami is a declarative domain-specific language for cloud infrastructure, centered around natural language descriptions. Salami compiler uses GPT4 to convert the natural language to Terraform files. You can think of Salami as writing documentation for each cloud resource object, and letting the compiler take care of converting that to IaC (Infrastructure as Code).

## 🎨 DSL Design

Salami files are mostly comprised of natural language, with several special constructs:

1. **Blocks** - multiline blocks of text that each represent either a resource or a variable; delimited by double newlines.
2. **Constructor functions** - functions that are used to specify the nature of the block; start with an `@` symbol.
3. **Variable references** - references to variables that are defined in the program; delimited by curly braces.
4. **Resource references** - references to resources that are defined in the program; start with a dollar sign.

<p align="center">
  <img src="docs/images/salami-example.png">
  <br>
  <small>Example Salami code with 3 blocks: VPC resource, Security Group resource and the <code>container_port</code> variable.</small>
</p>

`.sami` is the extension for Salami files. For more examples, see the `examples` directory.

## 🚀 Getting Started

### Installation

```bash
brew tap petrgazarov/salami
brew install salami
```

### Usage

From the root of your project, run:

```bash
salami compile
```

The root of your project should contain `salami.yaml` file with the following structure:

```yaml
compiler:
  target:
    platform: terraform
  llm:
    provider: openai
    model: gpt4
    api_key: ${OPENAI_API_KEY}
  source_dir: <source directory> # relative to the root of your project
  target_dir: <target directory> # relative to the root of your project
```

The `${OPENAI_API_KEY}` expression directs the compiler to read the api key from the `$OPENAI_API_KEY` environment variable.

### VS Code Extension

It's recommended to install the [VS Code extension](https://marketplace.visualstudio.com/items?itemName=Salami.salami). It provides highlighting functionality for the `.sami` files.
