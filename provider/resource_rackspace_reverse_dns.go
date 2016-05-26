package provider

import "github.com/hashicorp/terraform/helper/schema"

func resourceRackspaceReverseDNS() *schema.Resource {
	return &schema.Resource{
		Create: resourceRackspaceReverseDNSCreate,
		Read:   resourceRackspaceReverseDNSRead,
		Update: resourceRackspaceReverseDNSUpdate,
		Delete: resourceRackspaceReverseDNSDelete,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "name for dns record",
			},
			"data": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ip address for dns record",
			},
			"ttl": &schema.Schema{
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "ttl for dns record",
			},
		},
	}
}

func resourceRackspaceReverseDNSCreate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceRackspaceReverseDNSRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceRackspaceReverseDNSUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceRackspaceReverseDNSDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
