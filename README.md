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
