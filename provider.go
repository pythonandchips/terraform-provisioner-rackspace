package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/pythonandchips/terraform-provider-rackspace/provider"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: provider.Provider,
	})
}
