package provider

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestNewEntraGroupsForRoleDataSource(t *testing.T) {
	t.Parallel()

	// A timestamp is added to the name to avoid failure due to previous
	// test failures
	time := time.Now().Unix()
	roleName := fmt.Sprintf("TestNewEntraGroupsForRoleDataSource Role %d", time)
	testUser := os.Getenv("ACC_TEST_SYSTEM_ROLE_OWNER")
	itShopName := "General access shop shelf"

	
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
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
					name = "[Test] group to be assigned to role"
					alias = "group_to_be_assigned_to_role_test"
					description = "terraform provider acceptance test"
					inheritance_level = "User"
				}

				resource "tilgangsportalen_entra_group_role_assignment" "test_role_assignment_data_source" {
					role_name = tilgangsportalen_system_role.test_role_assignment_data_source.name
					entra_group = tilgangsportalen_entra_group.test_role_assignment_data_source.name
				}

				data "tilgangsportalen_entra_groups_assigned_to_role" "groups_assigned_role" {
					role_name = tilgangsportalen_system_role.test_role_assignment_data_source.name
				}
				`, roleName, testUser,itShopName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.tilgangsportalen_entra_groups_assigned_to_role.groups_assigned_role", "role_name", roleName),
				),
			},
		},
	})
}
