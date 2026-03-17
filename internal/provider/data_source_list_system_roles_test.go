package provider

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSystemRolesDataSource(t *testing.T) {
	t.Parallel()

	// A timestamp is added to the name to avoid failure due to previous test failures
	time := time.Now().Unix()
	roleName := fmt.Sprintf("TestSystemRolesDataSource Role %d", time)
	roleOwner := os.Getenv("ACC_TEST_SYSTEM_ROLE_OWNER")
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

				data "tilgangsportalen_system_roles" "all_roles" {
					depends_on = [tilgangsportalen_system_role.test_role_data_source]
				}
				`, roleName, roleOwner, itShopName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.tilgangsportalen_system_roles.all_roles", "roles.#"),             // Check that the roles list exists
					resource.TestCheckResourceAttrSet("data.tilgangsportalen_system_roles.all_roles", "roles.0.displayname"), // Check that the first role has a display name
					resource.TestCheckResourceAttrSet("data.tilgangsportalen_system_roles.all_roles", "roles.0.object_id"),   // Check that the first role has an object ID

					// Check that the created role is in the list
					resource.TestCheckTypeSetElemNestedAttrs("data.tilgangsportalen_system_roles.all_roles", "roles.*",
						map[string]string{
							"displayname": roleName,
						}),

					// Check that object_id is not empty
					resource.TestCheckResourceAttrWith("data.tilgangsportalen_system_roles.all_roles", "roles.0.object_id",
						func(value string) error {
							if value == "" { // Check that object_id is not empty
								return fmt.Errorf("object_id should not be empty")
							}
							return nil
						}),
				),
			},
		},
	})
}
