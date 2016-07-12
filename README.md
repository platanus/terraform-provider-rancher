# Rancher Terraform Provider

## Provider

The Rancher provider is used to interact with the
resources supported by Rancher. The provider needs to be configured
with the proper credentials before it can be used.

#### Example Usage

```
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

```
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
