package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDedicatedServerBringYourOwnImage_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCredentials(t)
			testAccPreCheckDedicatedServer(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDedicatedServerBringYourOwnImageConfig("basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_dedicated_server_byoi.server_install", "status", "done"),
				),
			},
		},
	})
}

func testAccDedicatedServerBringYourOwnImageConfig(config string) string {
	dedicated_server := os.Getenv("OVH_DEDICATED_SERVER")
	url := acctest.RandomWithPrefix(test_prefix)
	imgtype := acctest.RandomWithPrefix(test_prefix)

	// sshKey := os.Getenv("OVH_SSH_KEY")
	// if sshKey == "" {
	// 	sshKey = "ssh-ed25519 AAAAC3NzaC1yyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyy"
	// }

	return fmt.Sprintf(
		testAccDedicatedServerBringYourOwnImageConfig_Basic,
		dedicated_server,
		url,
		imgtype,
	)
}

const testAccDedicatedServerBringYourOwnImageConfig_Basic = `
resource "ovh_dedicated_server_byoi" "server_install" {
  service_name  = "%s"
  url           = "%s"
  type          = "%s"
}
`
