<a href="https://terraform.io">
    <img src=".github/tf.png" alt="Terraform logo" title="Terraform" align="left" height="50" />
</a>

# Terraform Provider for CANCOM Cloud

The CANCOM Terraform Provider allows managing resources within CANCOM Managed Services Cloud.

- Website: https://www.terraform.io
- Documentation: https://registry.terraform.io/providers/cancom/cancom/latest/docs

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.13.x or higher
-	[Go](https://golang.org/doc/install) 1.18 (to build the provider plugin)

## Building The Provider

Clone repository to: `$GOPATH/src/github.com/cancom/terraform-provider-cancom`

```sh
$ mkdir -p $GOPATH/src/github.com/cancom; cd $GOPATH/src/github.com/cancom
$ git clone git@github.com:cancom/terraform-provider-cancom
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/fastly/terraform-provider-cancom
$ make build
```
## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.18+ is *required*).

To compile the provider, run `make build`. This will build the provider and put the provider binary in a local `bin` directory.

```sh
$ make build
...
```

Alongside the newly built binary a file called `developer_overrides.tfrc` will be created.  The `make build` target will communicate
back details for setting the `TF_CLI_CONFIG_FILE` environment variable that will enable Terraform to use your locally built provider binary.

* HashiCorp - [Development Overrides for Provider developers](https://www.terraform.io/docs/cli/config/config-file.html#development-overrides-for-provider-developers).

> **NOTE**: If you have issues seeing any behaviours from code changes you've made to the provider, then it might be the terraform CLI is getting confused by which provider binary it should be using. Check inside the `./bin/` directory to see if there are multiple providers with different commit hashes (e.g. `terraform-provider-cancom_v1.1.0-3-fgcd28ca1`) and delete them first before running `make build`. This should help the Terraform CLI resolve to the correct binary.

## Building The Documentation

The documentation is built from components (go templates) stored in the `templates` folder.
Building the documentation copies the full markdown into the `docs` folder, ready for deployment to Hashicorp.

> NOTE: you'll need the [`tfplugindocs`](https://github.com/hashicorp/terraform-plugin-docs) tool for generating the Markdown to be deployed to Hashicorp. For more information on generating documentation, refer to https://www.terraform.io/docs/registry/providers/docs.html

* To validate the `/template` directory structure:
```
make validate-docs
```

* To build the `/docs` documentation Markdown files:
```
make generate-docs
```

* To view the documentation:
Paste `/docs` Markdown file content into https://registry.terraform.io/tools/doc-preview

## Contributing

Refer to [CONTRIBUTING.md](./CONTRIBUTING.md)
