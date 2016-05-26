package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/pythonandchips/terraform-provider-rackspace/rackspace"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("RACKSPACE_API", nil),
				Description: "Your rackspace api token",
			},
			"tenant_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "rackspace tenant id",
				DefaultFunc: schema.EnvDefaultFunc("RACKSPACE_TENANT", nil),
			},
			"region": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "location to build servers",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"rackspace_server": resourceRackspaceServer(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	client := rackspace.NewRackspaceClient(
		d.Get("token").(string),
		d.Get("tenant_id").(string),
		d.Get("region").(string),
	)
	return client, nil
}
