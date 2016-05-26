package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"rackspace": testAccProvider,
	}
}

func TestServerCreate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckRackspaceServerConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rackspace_server.foo", "name", "foo"),
					resource.TestCheckResourceAttr("rackspace_server.foo", "metadata", "foo"),
				),
			},
		},
	})
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("RACKSPACE_TOKEN"); v == "" {
		t.Fatal("RACKSPACE_TOKEN must be set for acceptance tests")
	}
}

var testAccCheckRackspaceServerConfig_basic = fmt.Sprintf(`
provider "rackspace" {
	token = "token1"
	tenant_id = "133919392"
	region = "lon"

}

resource "rackspace_server" "foo" {
	name        = "foo"
	image_ref   = "00000-000000-000000-00000-000000"
	flavor_ref = "2"
	metadata = {
		env = "production"
	}
}
`)
