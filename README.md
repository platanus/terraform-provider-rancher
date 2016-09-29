# Rancher Terraform Provider [![Build Status](https://travis-ci.org/platanus/terraform-provider-rancher.svg?branch=master)](https://travis-ci.org/platanus/terraform-provider-rancher) [![GitHub version](https://badge.fury.io/gh/platanus%2Fterraform-provider-rancher.svg)](https://badge.fury.io/gh/platanus%2Fterraform-provider-rancher)

## Provider

The Rancher provider is used to interact with the
resources supported by Rancher. The provider needs to be configured
with the proper credentials before it can be used.

#### Example Usage

```hcl
# Configure the Rancher provider
provider "rancher" {
     api_url = "http://rancher.my-domain.com/v1"
  access_key = "${var.rancher_access_key}"
  secret_key = "${var.rancher_secret_key}"
}
```

#### Argument Reference

The following arguments are supported:

* `api_url` - (Required) Rancher API url. It must be provided, but it can also be sourced from the `CATTLE_URL` environment variable.
* `access_key` - (Required) Rancher API access key. It must be provided, but it can also be sourced from the `CATTLE_ACCESS_KEY` environment variable.
* `secret_key` - (Required) Rancher API access key. It must be provided, but it can also be sourced from the `CATTLE_SECRET_KEY` environment variable.

## Resources

- [Environment](#environment)
- [Registration Token](#registration-token)
- [Registry](#registry)
- [Registry Credential](#registry-credential)
- [Stack](#stack)

### Environment

Provides a Rancher Environment resource. This can be used to create and manage environments on rancher.

#### Example Usage

```hcl
# Create a new Rancher environment
resource "rancher_environment" "default" {
  name = "staging"
  description = "The staging environment"
  orchestration = "cattle"
}
```

#### Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the environment.
* `description` - (Optional) An environment description.
* `orchestration` - (Optional) Must be one of **cattle**, **swarm**, **mesos** or **kubernetes**. Defaults to **cattle**.

#### Attributes Reference

The following attributes are exported:

* `id` - The ID of the environment.
* `name` - The name of the environment.
* `description` - The description of the environment.
* `orchestration` - The orchestration engine for the environment.

### Registration Token

Provides a Rancher Registration Token resource. This can be used to create registration tokens for rancher environments and retrieve their information.

#### Example Usage

```hcl
# Create a new Rancher registration token
resource "rancher_registration_token" "default" {
  name = "staging_token"
  description = "Registration token for the staging environment"
  environment_id = "${rancher_environment.default.id}"
}
```

#### Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the registration token.
* `description` - (Optional) A registration token description.
* `environment_id` - (Required) The ID of the environment to create the token for.

#### Attributes Reference

The following attributes are exported:

* `id` - The ID of the environment.
* `name` - The name of the registration token.
* `description` - The description of the registration token.
* `environment_id` - The ID of the environment to create the token for.
* `registration_url` - The URL to use to register new nodes to the environment.
* `token` - The token to use to register new nodes to the environment.

### Registry

Provides a Rancher Registy resource. This can be used to create registries for rancher environments and retrieve their information.

#### Example Usage

```hcl
# Create a new Rancher registry
resource "rancher_registry" "dockerhub" {
  name = "dockerhub"
  description = "DockerHub Registry"
  environment_id = "${rancher_environment.default.id}"
  server_address = "index.dockerhub.io"
}
```

#### Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the registry.
* `description` - (Optional) A registry description.
* `environment_id` - (Required) The ID of the environment to create the registry for.
* `server_address` - (Required) The server address for the registry.

#### Attributes Reference

The following attributes are exported:

* `id` - The ID of the registry.
* `name` - (Required) The name of the registry.
* `description` - (Optional) The registry description.
* `environment_id` - (Required) The ID of the environment to create the registry for.
* `server_address` - (Required) The server address for the registry.

### Registry Credential

Provides a Rancher Registy Credential resource. This can be used to create registry credentials for rancher environments and retrieve their information.

#### Example Usage

```hcl
# Create a new Rancher registry
resource "rancher_registry_credential" "dockerhub" {
  name = "dockerhub"
  description = "DockerHub Registry Credential"
  registry_id = "${rancher_registry.dockerhub.id}"
  email = "myself@company.com"
  public_value = "myself"
  secret_value = "mypass"
}
```

#### Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the registry credential.
* `description` - (Optional) A registry credential description.
* `registry_id` - (Required) The ID of the registry to create the credential for.
* `email` - (Required) The email of the account.
* `public_value` - (Required) The public value (user name) of the account.
* `secret_value` - (Required) The secret value (password) of the account.

#### Attributes Reference

The following attributes are exported:

* `id` - The ID of the registry credential.
* `name` - (Required) The name of the registry credential.
* `description` - (Optional) The registry credential description.
* `registry_id` - (Required) The ID of the registry to create the credential for.
* `email` - (Required) The email of the account.
* `public_value` - (Required) The public value (user name) of the account.
* `secret_value` - (Required) The secret value (password) of the account.

### Stack

Provides a Rancher Stack resource. This can be used to create and manage stacks on rancher.

#### Example Usage

```hcl
# Create a new empty Rancher stack
resource "rancher_stack" "external-dns" {
  name = "route53"
  description = "Route53 stack"
  environment_id = "${rancher_environment.default.id}"
  catalog_id = "library:route53:7"
  scope = "system"
  environment {
    AWS_ACCESS_KEY = "MYKEY"
    AWS_SECRET_KEY = "MYSECRET"
    AWS_REGION = "eu-central-1"
    TTL = "60"
    ROOT_DOMAIN = "example.com"
    ROUTE53_ZONE_ID = ""
    HEALTH_CHECK_INTERVAL = "15"
  }
}
```

#### Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the stack.
* `description` - (Optional) A stack description.
* `environment_id` - (Required) The ID of the environment to create the stack for.
* `docker_compose` - (Optional) The `docker-compose.yml` content to apply for the stack.
* `rancher_compose` - (Optional) The `rancher-compose.yml` content to apply for the stack.
* `environment` - (Optional) The environment to apply to interpret the docker-compose and rancher-compose files.
* `catalog_id` - (Optional) The catalog ID to link this stack to. When provided, `docker_compose` and `rancher_compose` will be retrieved from the catalog unless they are overridden.
* `scope` - (Optional) The scope to attach the stack to. Must be one of **user** or **system**. Defaults to **user**.
* `start_on_create` - (Optional) Whether to start the stack automatically.

#### Attributes Reference

The following attributes are exported:

* `id` - The ID of the stack.
* `name` - The name of the stack.
* `description` - The description of the stack.
* `environment_id` - (Required) The ID of the environment to create the stack for.
* `docker_compose` - (Optional) The `docker-compose.yml` content to apply for the stack.
* `rancher_compose` - (Optional) The `rancher-compose.yml` content to apply for the stack.
* `environment` - (Optional) The environment to apply to interpret the docker-compose and rancher-compose files.
* `catalog_id` - (Optional) The catalog ID to link this stack to. When provided, `docker_compose` and `rancher_compose` will be retrieved from the catalog unless they are overridden.
* `scope` - (Optional) The scope to attach the stack to. Must be one of **user** or **system**. Defaults to **user**.
* `start_on_create` - (Optional) Whether to start the stack automatically.
