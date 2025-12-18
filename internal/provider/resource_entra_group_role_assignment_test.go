package provider

import (
	"fmt"
	"os"
	"regexp"
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
	groupName := fmt.Sprintf("[APPTEST] TestEntraGroupRoleAssignment %d", time)
	newGroupName := fmt.Sprintf("[APPTEST] TestEntraGroupRoleAssignment New Name %d", time)
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
					description = "APPTEST"
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
					description = "APPTEST"
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
					description = "APPTEST"
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

func TestEntraGroupRoleAssignmentWrongCasingRole(t *testing.T) {
	t.Parallel()

	time := time.Now().Unix()
	roleName := fmt.Sprintf("TestEntraGroupRoleAssignmentWrongCasingRole Role %d", time)
	roleNameWrongCasing := fmt.Sprintf("TestEntraGroupRoleAssignmentWrongCasingRole role %d", time)
	testUser := os.Getenv("ACC_TEST_SYSTEM_ROLE_OWNER")
	groupName := fmt.Sprintf("[APPTEST] TestEntraGroupRoleAssignmentWrongCasingRole %d", time)
	itShopName := "General access shop shelf"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			// Create role and group with correct casing
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "tilgangsportalen_system_role" "test" {
					name              = "%s"
					product_category  = "TBD"
					system_role_owner = "%s"
					approval_level    = "L2"
					description       = "Terraform acceptance test role for casing validation."
					it_shop_name      = "%s"
				}

				resource "tilgangsportalen_entra_group" "test" {
					name = "%s"
					description = "APPTEST"
					inheritance_level = "User"
				}
				`, roleName, testUser, itShopName, groupName),
			},
			// Try to create assignment with wrong casing on role name - should fail in Create
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "tilgangsportalen_system_role" "test" {
					name              = "%s"
					product_category  = "TBD"
					system_role_owner = "%s"
					approval_level    = "L2"
					description       = "Terraform acceptance test role for casing validation."
					it_shop_name      = "%s"
				}

				resource "tilgangsportalen_entra_group" "test" {
					name = "%s"
					description = "APPTEST"
					inheritance_level = "User"
				}

				resource "tilgangsportalen_entra_group_role_assignment" "test" {
					role_name = "%s"
					entra_group = tilgangsportalen_entra_group.test.name
					force = true
				}
				`, roleName, testUser, itShopName, groupName, roleNameWrongCasing),
				ExpectError: regexp.MustCompile(`Role Name Casing Mismatch`),
			},
		},
	})
}

func TestEntraGroupRoleAssignmentWrongCasingGroup(t *testing.T) {
	t.Parallel()

	time := time.Now().Unix()
	roleName := fmt.Sprintf("TestEntraGroupRoleAssignmentWrongCasingGroup Role %d", time)
	testUser := os.Getenv("ACC_TEST_SYSTEM_ROLE_OWNER")
	groupName := fmt.Sprintf("[APPTEST] TestEntraGroupRoleAssignmentWrongCasingGroup %d", time)
	groupNameWrongCasing := fmt.Sprintf("[APPTEST] Testentragrouproleassignmentwrongcasinggroup %d", time)
	itShopName := "General access shop shelf"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			// Create role and group with correct casing
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "tilgangsportalen_system_role" "test" {
					name              = "%s"
					product_category  = "TBD"
					system_role_owner = "%s"
					approval_level    = "L2"
					description       = "Terraform acceptance test group for casing validation."
					it_shop_name      = "%s"
				}

				resource "tilgangsportalen_entra_group" "test" {
					name = "%s"
					description = "APPTEST"
					inheritance_level = "User"
				}
				`, roleName, testUser, itShopName, groupName),
			},
			// Try to create assignment with wrong casing on group name - should fail in Create
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "tilgangsportalen_system_role" "test" {
					name              = "%s"
					product_category  = "TBD"
					system_role_owner = "%s"
					approval_level    = "L2"
					description       = "Terraform acceptance test group for casing validation."
					it_shop_name      = "%s"
				}

				resource "tilgangsportalen_entra_group" "test" {
					name = "%s"
					description = "APPTEST"
					inheritance_level = "User"
				}

				resource "tilgangsportalen_entra_group_role_assignment" "test" {
					role_name = tilgangsportalen_system_role.test.name
					entra_group = "%s"
					force = true
				}
				`, roleName, testUser, itShopName, groupName, groupNameWrongCasing),
				ExpectError: regexp.MustCompile(`Entra Group Name Casing Mismatch`),
			},
		},
	})
}
