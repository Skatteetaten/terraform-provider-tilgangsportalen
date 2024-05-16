package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestCreateNewEntraGroupRoleAssignment(t *testing.T) {

	roleName := "role acceptance test assignment"
	groupName := "[Test] group acceptance test assignment"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				ResourceName: "tilgangsportalen_entra_group_role_assignment.test_role_group_assignment",
				Config: providerConfig + fmt.Sprintf(`
				resource "tilgangsportalen_system_role" "test_role_group_assignment" {
					name              = "%s"
					product_category  = "TBD"
					system_role_owner = "M00245"
					approval_level    = "L2"
					description       = "Terraform acceptance test role for assignment."
					it_shop_name      = "Access shop shelf"
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
				`,roleName,groupName),
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
