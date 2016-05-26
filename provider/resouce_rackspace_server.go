package provider

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pythonandchips/terraform-provider-rackspace/rackspace"
)

func resourceRackspaceServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceRackspaceServerCreate,
		Read:   resourceRackspaceServerRead,
		Update: resourceRackspaceServerUpdate,
		Delete: resourceRackspaceServerDelete,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The server name",
			},
			"image_ref": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The image reference for the desired image for your server instance.",
			},
			"block_device_mapping_v2": &schema.Schema{
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "The container of bootable volume details",
				Elem:        &schema.Schema{},
			},
			"flavor_ref": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The flavor reference for the desired flavor for your server instance",
			},
			"config_drive": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Enables metadata injection in a server through a configuration drive. To enable a configuration drive, specify true. Otherwise, specify false",
			},
			"key_name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the key pair used to authenticate by using key-based authentication instead of password- based authentication",
			},
			"os_dcf": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The disk configuration value. The image auto_disk_config metadata key set will affect the value you can choose to set the server OS- DCF:diskConfig. If an image has auto_disk_config value of disabled, you cannot create a server from that image when specifying OS-DCF:diskConfig value of AUTO. Valid values are: AUTO:The server is built with a single partition which is the size of the target flavor disk. The file system is automatically adjusted to fit the entire partition. This keeps things simple and automated. AUTO is valid only for images and servers with a single partition that use the EXT3 file system. This is the default setting for applicable Rackspace base images. MANUAL:The server is built using the partition scheme and file system of the source image. If the target flavor disk is larger, the remaining disk space is left unpartitioned. This enables images to have non-EXT3 file systems, multiple partitions, and so on, and it enables you to manage the disk configuration",
			},
			"metadata": &schema.Schema{
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Metadata key and value pairs. The maximum size of each metadata key and value is 255 bytes each",
			},
			"personality": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The array of personality files for the server",
				Elem:        &schema.Schema{},
			},
			"user_data": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Data used with config_drive for configuring a server",
			},
			"networks": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The array of networks attached to the server. By default, the server instance is provisioned with all isolated networks for the tenant. You can specify multiple NICs on the server. Optionally, you can create one or more NICs on the server. To provision the server instance with a NIC for a Nova- network network, specify the UUID in the uuid attribute in a networks object. To provision the server instance with a NIC for a Neutron network, specify the UUID in the port attribute in a networks object",
				Elem:        &schema.Schema{},
			},
			//returned data
			"admin_pass": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{},
			},
			"accessIPv4": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"progress": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceRackspaceServerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(rackspace.RackspaceClient)

	createServerRequest := rackspace.CreateServerRequest{
		Name:      d.Get("name").(string),
		ImageRef:  d.Get("image_ref").(string),
		FlavorRef: d.Get("flavor_ref").(string),
	}

	if keyName, ok := d.GetOk("key_name"); ok {
		createServerRequest.KeyName = keyName.(string)
	}

	log.Printf("[DEBUG] Creating service: %#v", createServerRequest)

	resp, err := client.CreateServer(createServerRequest)
	if err != nil {
		return err
	}

	d.SetId(resp.Id)
	d.Set("admin_pass", resp.AdminPass)

	log.Printf("[INFO] Server ID: %s", d.Id())

	_, err = waitForServerReady(d, meta)
	return resourceRackspaceServerRead(d, meta)
}

func resourceRackspaceServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(rackspace.RackspaceClient)
	id := d.Id()
	resp, err := client.ReadServer(id)
	if err != nil {
		return nil
	}

	d.Set("status", resp.Status)
	d.Set("progress", resp.Progress)
	d.Set("accessIPv4", resp.AccessIPv4)

	return nil
}

func resourceRackspaceServerUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceRackspaceServerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(rackspace.RackspaceClient)
	err := client.DestroyServer(d.Id())
	if err != nil {
		return fmt.Errorf("Error deleting server: %s", err)
	}
	return nil
}

func waitForServerReady(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	stateChange := &resource.StateChangeConf{
		Pending:        []string{"BUILD"},
		Target:         []string{"status"},
		Refresh:        newServerStateRefreshFunc(d, meta),
		Timeout:        60 * time.Minute,
		Delay:          10 * time.Second,
		MinTimeout:     3 * time.Second,
		NotFoundChecks: 60,
	}
	return stateChange.WaitForState()
}

func newServerStateRefreshFunc(d *schema.ResourceData, meta interface{}) resource.StateRefreshFunc {
	client := meta.(rackspace.RackspaceClient)
	return func() (interface{}, string, error) {
		id := d.Id()
		err := resourceRackspaceServerRead(d, meta)
		if err != nil {
			return nil, "", err
		}
		if status, ok := d.GetOk("status"); ok {
			resp, _ := client.ReadServer(id)
			log.Printf("[DEBUG] server status: %s", status.(string))
			log.Printf("[DEBUG] server progress: %s", d.Get("progress").(string))
			return resp, status.(string), nil

		}
		return nil, "", nil
	}
}
