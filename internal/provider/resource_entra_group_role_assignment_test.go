package provider

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestCreateNewEntraGroupRoleAssignment(t *testing.T) {
	t.Parallel()

	// A timestamp is added to the name to avoid failure due to previous
	// test failures
	time := time.Now().Unix()
	roleName := fmt.Sprintf("TestCreateNewEntraGroupRoleAssignment Role %d", time)
	testUser := os.Getenv("ACC_TEST_SYSTEM_ROLE_OWNER")
	groupName := fmt.Sprintf("[Group] TestCreateNewEntraGroupRoleAssignment %d", time)
	itShopName := "General access shop shelf"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				ResourceName: "tilgangsportalen_entra_group_role_assignment.test_role_group_assignment",
				Config: providerConfig + fmt.Sprintf(`
				resource "tilgangsportalen_system_role" "test_role_group_assignment" {
					name              = "%s"
					product_category  = "TBD"
					system_role_owner = "%s"
					approval_level    = "L2"
					description       = "Terraform acceptance test role for assignment."
					it_shop_name      = "%s"
				} 

				resource "tilgangsportalen_entra_group" "test_role_group_assignment" {
					name = "%s"
					alias = "group_acceptance_test_for_assignment"
					description = "terraform provider acceptance test"
					inheritance_level = "User"
				}

				resource "tilgangsportalen_entra_group_role_assignment" "test_role_group_assignment" {
					role_name = tilgangsportalen_system_role.test_role_group_assignment.name
					entra_group = tilgangsportalen_entra_group.test_role_group_assignment.name
				}
				`, roleName, testUser, itShopName, groupName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tilgangsportalen_entra_group_role_assignment.test_role_group_assignment", "role_name", roleName),
					resource.TestCheckResourceAttr("tilgangsportalen_entra_group_role_assignment.test_role_group_assignment", "entra_group", groupName),
				),
			},
			// test import to state using ImportStateCheckFunc
			{
				ImportState:             true,
				ResourceName:            "tilgangsportalen_entra_group_role_assignment.test_role_group_assignment",
				ImportStateVerifyIgnore: []string{"alias", "description", "inheritance_level"},
			},
		},
	})
}
