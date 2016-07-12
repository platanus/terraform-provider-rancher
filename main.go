package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/platanus/terraform-provider-rancher/rancher"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: rancher.Provider,
	})
}
