package provider

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestListAllSystemRolesForEntraGroupDataSource(t *testing.T) {
	t.Parallel()

	// A timestamp is added to the name to avoid failure due to previous test failures
	time := time.Now().Unix()
	roleName := fmt.Sprintf("TestListAllSystemRolesForEntraGroup Role %d", time)
	groupName := fmt.Sprintf("[Test] group to be assigned TestListAllSystemRolesForEntraGroup %d", time)
	testUser := os.Getenv("ACC_TEST_SYSTEM_ROLE_OWNER")
	itShopName := "General access shop shelf"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
                resource "tilgangsportalen_entra_group" "test_role_assignment_data_source" {
                    name = "%s"
                    description = "terraform provider acceptance test"
                    inheritance_level = "User"
                }

                data "tilgangsportalen_system_roles_assigned_to_entra_group" "roles_for_group" {
                    group_name = tilgangsportalen_entra_group.test_role_assignment_data_source.name
					
                }
                `, groupName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.tilgangsportalen_system_roles_assigned_to_entra_group.roles_for_group", "group_name", groupName),
					// confirm that roles for the group are empty
					resource.TestCheckResourceAttr("data.tilgangsportalen_system_roles_assigned_to_entra_group.roles_for_group", "roles.#", "0"),
				),
			},
			// add group role assignment and check that the role is listed
			{
				Config: providerConfig + fmt.Sprintf(`
                resource "tilgangsportalen_system_role" "test_role_assignment_data_source" {
                    name              = "%s"
                    product_category  = "TBD"
                    system_role_owner = "%s"
                    approval_level    = "L2"
                    description       = "Terraform acceptance test role for assignment."
                    it_shop_name      = "%s"
                } 

                resource "tilgangsportalen_entra_group" "test_role_assignment_data_source" {
                    name = "%s"
                    description = "terraform provider acceptance test"
                    inheritance_level = "User"
                }

                resource "tilgangsportalen_entra_group_role_assignment" "test_role_assignment_data_source" {
                    role_name = tilgangsportalen_system_role.test_role_assignment_data_source.name
                    entra_group = tilgangsportalen_entra_group.test_role_assignment_data_source.name
                }

                data "tilgangsportalen_system_roles_assigned_to_entra_group" "roles_for_group" {
                    group_name = tilgangsportalen_entra_group.test_role_assignment_data_source.name
					
					depends_on = [tilgangsportalen_entra_group_role_assignment.test_role_assignment_data_source]
                }
                `, roleName, testUser, itShopName, groupName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.tilgangsportalen_system_roles_assigned_to_entra_group.roles_for_group", "group_name", groupName),
					resource.TestCheckResourceAttr("data.tilgangsportalen_system_roles_assigned_to_entra_group.roles_for_group", "roles.#", "1"),
					resource.TestCheckResourceAttr("data.tilgangsportalen_system_roles_assigned_to_entra_group.roles_for_group", "roles.0.display_name", roleName),
					resource.TestCheckResourceAttr("data.tilgangsportalen_system_roles_assigned_to_entra_group.roles_for_group", "roles.0.system_role_owner", testUser),
					resource.TestCheckResourceAttrSet("data.tilgangsportalen_system_roles_assigned_to_entra_group.roles_for_group", "roles.0.system_role_owner_display_name"),
				),
			},
		},
	})
}
