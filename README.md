<a href="https://terraform.io">
    <img src=".github/tf.png" alt="Terraform logo" title="Terraform" align="left" height="50" />
</a>

# Terraform Provider for CANCOM Cloud

The CANCOM Terraform Provider allows managing resources within CANCOM Managed Services Cloud.

- Website: https://www.terraform.io
- Documentation: https://registry.terraform.io/providers/cancom/cancom/latest/docs

_**Please note:** We take Terraform's security and our users' trust very seriously. If you believe you have found a security issue in the Terraform CANCOM Provider, please responsibly disclose it by contacting us at security@cancom.io._

## Usage Example

```hcl
# 1. Specify the version of the CANCOM Provider to use
terraform {
  required_providers {
    cancom = {
      source = "cancom/cancom"
      version = "0.0.1"
    }
  }
}

# 2. Configure the CANCOM Provider
# You can provide your API token via CANCOM_TOKEN environment variable, representing your CANCOM token.
# When using this method, you may omit the CANCOM provider block entirely.
provider "cancom" {
  token = "<token>"
}

# 3. Create a resource group
resource "cancom_dns_record" "test" {
  zone_name = "example.com"
  name      = "test.example.com"
  type      = "A"
  content   = "127.0.0.1"
}
```

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.14.x or higher
- [Go](https://golang.org/doc/install) 1.21 (to build the provider plugin)

## Building The Provider

Clone repository to: `$GOPATH/src/github.com/cancom/terraform-provider-cancom`

```sh
$ mkdir -p $GOPATH/src/github.com/cancom; cd $GOPATH/src/github.com/cancom
$ git clone git@github.com:cancom/terraform-provider-cancom
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/cancom/terraform-provider-cancom
$ make build
```

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.21+ is _required_).

To compile the provider, run `make build`. This will build the provider and put the provider binary in a local `bin` directory.

```sh
$ make build
...
```

Alongside the newly built binary a file called `developer_overrides.tfrc` will be created. The `make build` target will communicate
back details for setting the `TF_CLI_CONFIG_FILE` environment variable that will enable Terraform to use your locally built provider binary.

- HashiCorp - [Development Overrides for Provider developers](https://www.terraform.io/docs/cli/config/config-file.html#development-overrides-for-provider-developers).

> **NOTE**: If you have issues seeing any behaviours from code changes you've made to the provider, then it might be the terraform CLI is getting confused by which provider binary it should be using. Check inside the `./bin/` directory to see if there are multiple providers with different commit hashes (e.g. `terraform-provider-cancom_v1.1.0-3-fgcd28ca1`) and delete them first before running `make build`. This should help the Terraform CLI resolve to the correct binary.

## Building The Documentation

The documentation is built from components (go templates) stored in the `templates` folder.
Building the documentation copies the full markdown into the `docs` folder, ready for deployment to Hashicorp.

> NOTE: you'll need the [`tfplugindocs`](https://github.com/hashicorp/terraform-plugin-docs) tool for generating the Markdown to be deployed to Hashicorp. For more information on generating documentation, refer to https://www.terraform.io/docs/registry/providers/docs.html

- To validate the `/template` directory structure:

```
make validate-docs
```

- To build the `/docs` documentation Markdown files:

```
make generate-docs
```

- To view the documentation:
  Paste `/docs` Markdown file content into https://registry.terraform.io/tools/doc-preview

Note: The description and the subcategory of resource are now automatically generated from the resource's `Description` field within the go file, similar to how the schema's arguments are generated.
This allows you to define resources descriptions with `My Subcategory --- The description of my resource` and have the subcategory and description filled in correctly. It also works if you don't provide a subcategory, in that case, it will just populate the description.  
The only caveat is that you can't use additional `---` in your description.  
Imports and examples are automatically filled in if they are placed and named correctly in `examples/{type}/{resource_name}`

## Contributing

Refer to [CONTRIBUTING.md](./CONTRIBUTING.md)
