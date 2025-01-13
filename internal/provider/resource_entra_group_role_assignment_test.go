package provider

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestEntraGroupRoleAssignment(t *testing.T) {
	t.Parallel()

	// A timestamp is added to the name to avoid failure due to previous
	// test failures
	time := time.Now().Unix()
	roleName := fmt.Sprintf("TestEntraGroupRoleAssignment Role %d", time)
	newRoleName := fmt.Sprintf("TestEntraGroupRoleAssignment New Name Role %d", time)
	testUser := os.Getenv("ACC_TEST_SYSTEM_ROLE_OWNER")
	groupName := fmt.Sprintf("[Group] TestEntraGroupRoleAssignment %d", time)
	newGroupName := fmt.Sprintf("[Group] TestEntraGroupRoleAssignment New Name %d", time)
	itShopName := "General access shop shelf"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				ResourceName: "tilgangsportalen_entra_group_role_assignment.test",
				Config: providerConfig + fmt.Sprintf(`
				resource "tilgangsportalen_system_role" "test" {
					name              = "%s"
					product_category  = "TBD"
					system_role_owner = "%s"
					approval_level    = "L2"
					description       = "Terraform acceptance test role for assignment."
					it_shop_name      = "%s"
				} 

				resource "tilgangsportalen_entra_group" "test" {
					name = "%s"
					description = "terraform provider acceptance test"
					inheritance_level = "User"
				}

				resource "tilgangsportalen_entra_group_role_assignment" "test" {
					role_name = tilgangsportalen_system_role.test.name
					entra_group = tilgangsportalen_entra_group.test.name
					force = true
				}
				`, roleName, testUser, itShopName, groupName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tilgangsportalen_entra_group_role_assignment.test", "role_name", roleName),
					resource.TestCheckResourceAttr("tilgangsportalen_entra_group_role_assignment.test", "entra_group", groupName),
				),
			},
			{
				ImportState:             true,
				ResourceName:            "tilgangsportalen_entra_group_role_assignment.test",
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force"},
			},
			// Test rename system role while group is assigned
			{
				Config: providerConfig + fmt.Sprintf(`		  
				resource "tilgangsportalen_system_role" "test" {
					name              = "%s"
					product_category  = "TBD"
					system_role_owner = "%s"
					approval_level    = "L2"
					description       = "Terraform acceptance test role for assignment."
					it_shop_name      = "%s"
				} 
				resource "tilgangsportalen_entra_group" "test" {
					name = "%s"
					description = "terraform provider acceptance test"
					inheritance_level = "User"
				}

				resource "tilgangsportalen_entra_group_role_assignment" "test" {
					role_name = tilgangsportalen_system_role.test.name
					entra_group = tilgangsportalen_entra_group.test.name
					force = true
				}
				`, newRoleName, testUser, itShopName, groupName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tilgangsportalen_entra_group_role_assignment.test", "role_name", newRoleName),
					resource.TestCheckResourceAttr("tilgangsportalen_entra_group_role_assignment.test", "entra_group", groupName),
				),
			},
			// Test rename group while assigned to role
			{
				Config: providerConfig + fmt.Sprintf(`		  
				resource "tilgangsportalen_system_role" "test" {
					name              = "%s"
					product_category  = "TBD"
					system_role_owner = "%s"
					approval_level    = "L2"
					description       = "Terraform acceptance test role for assignment."
					it_shop_name      = "%s"
				} 
				resource "tilgangsportalen_entra_group" "test" {
					name = "%s"
					description = "terraform provider acceptance test"
					inheritance_level = "User"
				}

				resource "tilgangsportalen_entra_group_role_assignment" "test" {
					role_name = tilgangsportalen_system_role.test.name
					entra_group = tilgangsportalen_entra_group.test.name
					force = true
				}
				`, newRoleName, testUser, itShopName, newGroupName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tilgangsportalen_entra_group_role_assignment.test", "role_name", newRoleName),
					resource.TestCheckResourceAttr("tilgangsportalen_entra_group_role_assignment.test", "entra_group", newGroupName),
				),
			},
		},
	})
}
