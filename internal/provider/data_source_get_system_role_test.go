package provider

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// create system role and read values using data source
func TestSystemRoleDataSource(t *testing.T) {
	t.Parallel()

	// A timestamp is added to the name to avoid failure due to previous
	// test failures
	time := time.Now().Unix()
	roleName := fmt.Sprintf("TestNewSystemRoleDataSource Role %d", time)
	testUser := os.Getenv("ACC_TEST_SYSTEM_ROLE_OWNER")
	itShopName := "General access shop shelf"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "tilgangsportalen_system_role" "test_role_data_source" {
					name              = "%s"
					product_category  = "TBD"
					system_role_owner = "%s"
					approval_level    = "L2"
					description       = "Terraform acceptance test role for assignment."
					it_shop_name      = "%s"

				} 

				data "tilgangsportalen_system_role" "system_role" {
					name = tilgangsportalen_system_role.test_role_data_source.name
				}
				`, roleName, testUser, itShopName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.tilgangsportalen_system_role.system_role", "name", roleName),
					resource.TestCheckResourceAttr("data.tilgangsportalen_system_role.system_role", "system_role_owner", testUser),
					resource.TestCheckResourceAttr("data.tilgangsportalen_system_role.system_role", "approval_level", "L2"),
					resource.TestCheckResourceAttr("data.tilgangsportalen_system_role.system_role", "description", "Terraform acceptance test role for assignment."),
				),
			},
		},
	})
}
