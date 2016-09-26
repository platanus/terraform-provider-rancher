# Rancher Terraform Provider

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
