# Jsonnet - project build description language for BrewKit

# Context and Problem Statement

The project build description language for BrewKit should support a powerful templating mechanism, imports (for splitting into files) and have internal validation.

# Decision Drivers

- The configuration language must have a powerful templating mechanism
- The language should either be easy to learn or be a mockup of one of the popular languages (YAML. JSON).

# Considered Options

- YAML
- JSON
- HCL
- JSONNET

# Decision Outcome

Decision selected: "Use Jsonnet". Jsonnet offers great flexibility in templating Json documents and has a Go version that can be used as a library.

# Pros and Cons of the Options

## YAML

- Good, popular description format
- Good, considered more human-friendly than JSON
- Good, supports templating via anchors
- Bad, has many problems related to its design (typing, nesting)
- Bad, not enough templating tools to describe the build of a complex project

## JSON

- Good, standard in describing schemas, configs, etc.
- Bad, no templating => not suitable for describing a complex project assembly

## HCL - [github](https://github.com/hashicorp/hcl)

A configuration description language from HashiCorp. Mostly used for [Terraform](https://www.terraform.io/) configuration.

- Good, has built-in validation (you can set field optionality, etc.).
- Good, is a new language for declarative style description
- Bad, not enough templating tools to describe the build of a complex project (lack of variables, etc.).

## Jsonnet - [Jsonnet](https://jsonnet.org/)

A superset of JSON augmented with functions, variables and object-oriented approach. Implemented by Google.
Used for description in AML projects, actively used by [Grafana Labs](https://grafana.github.io/grafonnet-lib/).

- Good, full-featured language with support for f-i, variables, etc.
- Good, output is simple JSON that can be viewed and simply parsed.
- Good, strict type system (as in JSON)
- Good, allows you to build a custom project build process
- Good, powerful templateizer, file import - you can build a project of great complexity.
- Good, has a full-featured Go implementation that can be used as a library
- Bad, little known and few examples on the Internet (ChatGPT knows the language very well:) )
- Bad, when describing a large project - you can get a bulky description file.