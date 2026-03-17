package provider

import (
	"fmt"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestTermsOfUseRoleAssignmentCreate(t *testing.T) {
	t.Parallel()

	// A timestamp is added to the name to avoid failure due to previous
	// test failures
	time := time.Now().Unix()
	roleName := fmt.Sprintf("TestTermsOfUseRoleAssignment Role %d", time)
	testUser := os.Getenv("ACC_TEST_SYSTEM_ROLE_OWNER")
	itShopName := "General access shop shelf"
	termsOfUse := "Test TermsOfUse 2"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "tilgangsportalen_system_role" "test" {
					name              = "%s"
					product_category  = "TBD"
					system_role_owner = "%s"
					approval_level    = "L2"
					description       = "Terraform acceptance test role for terms of use assignment."
					it_shop_name      = "%s"
				} 

				resource "tilgangsportalen_terms_of_use_role_assignment" "test" {
					role_name = tilgangsportalen_system_role.test.name
					terms_of_use_identifier = "%s"
				}
				`, roleName, testUser, itShopName, termsOfUse),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tilgangsportalen_terms_of_use_role_assignment.test", "role_name", roleName),
					resource.TestCheckResourceAttr("tilgangsportalen_terms_of_use_role_assignment.test", "terms_of_use_identifier", termsOfUse),
					resource.TestCheckResourceAttrSet("tilgangsportalen_terms_of_use_role_assignment.test", "id"),
					resource.TestCheckResourceAttrSet("tilgangsportalen_terms_of_use_role_assignment.test", "terms_of_use_description"),
				),
			},
		},
	})
}

func TestTermsOfUseRoleAssignmentUpdate(t *testing.T) {
	t.Parallel()

	// A timestamp is added to the name to avoid failure due to previous
	// test failures
	time := time.Now().Unix()
	roleName := fmt.Sprintf("TestTermsOfUseRoleAssignmentUpdate Role %d", time)
	testUser := os.Getenv("ACC_TEST_SYSTEM_ROLE_OWNER")
	itShopName := "General access shop shelf"
	termsOfUse1 := "Test TermsOfUse 2"
	termsOfUse2 := "Test TermsOfUse"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			// Test initial creation
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "tilgangsportalen_system_role" "test" {
					name              = "%s"
					product_category  = "TBD"
					system_role_owner = "%s"
					approval_level    = "L2"
					description       = "Terraform acceptance test role for terms of use assignment."
					it_shop_name      = "%s"
				} 

				resource "tilgangsportalen_terms_of_use_role_assignment" "test" {
					role_name = tilgangsportalen_system_role.test.name
					terms_of_use_identifier = "%s"
				}
				`, roleName, testUser, itShopName, termsOfUse1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tilgangsportalen_terms_of_use_role_assignment.test", "role_name", roleName),
					resource.TestCheckResourceAttr("tilgangsportalen_terms_of_use_role_assignment.test", "terms_of_use_identifier", termsOfUse1),
					resource.TestCheckResourceAttrSet("tilgangsportalen_terms_of_use_role_assignment.test", "id"),
					resource.TestCheckResourceAttrSet("tilgangsportalen_terms_of_use_role_assignment.test", "terms_of_use_description"),
				),
			},
			// Test update terms of use in place
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "tilgangsportalen_system_role" "test" {
					name              = "%s"
					product_category  = "TBD"
					system_role_owner = "%s"
					approval_level    = "L2"
					description       = "Terraform acceptance test role for terms of use assignment."
					it_shop_name      = "%s"
				}

				resource "tilgangsportalen_terms_of_use_role_assignment" "test" {
					role_name = tilgangsportalen_system_role.test.name
					terms_of_use_identifier = "%s"
				}
				`, roleName, testUser, itShopName, termsOfUse2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tilgangsportalen_terms_of_use_role_assignment.test", "role_name", roleName),
					resource.TestCheckResourceAttr("tilgangsportalen_terms_of_use_role_assignment.test", "terms_of_use_identifier", termsOfUse2),
					resource.TestCheckResourceAttrSet("tilgangsportalen_terms_of_use_role_assignment.test", "id"),
					resource.TestCheckResourceAttrSet("tilgangsportalen_terms_of_use_role_assignment.test", "terms_of_use_description"),
				),
			},
		},
	})
}

func TestTermsOfUseRoleAssignmentRoleRename(t *testing.T) {
	t.Parallel()

	// A timestamp is added to the name to avoid failure due to previous
	// test failures
	time := time.Now().Unix()
	roleName := fmt.Sprintf("TestTermsOfUseRoleAssignmentRoleRename Role %d", time)
	newRoleName := fmt.Sprintf("TestTermsOfUseRoleAssignmentRoleRename New Role %d", time)
	testUser := os.Getenv("ACC_TEST_SYSTEM_ROLE_OWNER")
	itShopName := "General access shop shelf"
	termsOfUse := "Test TermsOfUse"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			// Test initial creation
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "tilgangsportalen_system_role" "test" {
					name              = "%s"
					product_category  = "TBD"
					system_role_owner = "%s"
					approval_level    = "L2"
					description       = "Terraform acceptance test role for terms of use assignment."
					it_shop_name      = "%s"
				} 

				resource "tilgangsportalen_terms_of_use_role_assignment" "test" {
					role_name = tilgangsportalen_system_role.test.name
					terms_of_use_identifier = "%s"
				}
				`, roleName, testUser, itShopName, termsOfUse),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tilgangsportalen_terms_of_use_role_assignment.test", "role_name", roleName),
					resource.TestCheckResourceAttr("tilgangsportalen_terms_of_use_role_assignment.test", "terms_of_use_identifier", termsOfUse),
				),
			},
			// Test role rename (should trigger replace)
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "tilgangsportalen_system_role" "test" {
					name              = "%s"
					product_category  = "TBD"
					system_role_owner = "%s"
					approval_level    = "L2"
					description       = "Terraform acceptance test role for terms of use assignment."
					it_shop_name      = "%s"
				} 

				resource "tilgangsportalen_terms_of_use_role_assignment" "test" {
					role_name = tilgangsportalen_system_role.test.name
					terms_of_use_identifier = "%s"
				}
				`, newRoleName, testUser, itShopName, termsOfUse),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tilgangsportalen_terms_of_use_role_assignment.test", "role_name", newRoleName),
					resource.TestCheckResourceAttr("tilgangsportalen_terms_of_use_role_assignment.test", "terms_of_use_identifier", termsOfUse),
				),
			},
		},
	})
}

func TestTermsOfUseRoleAssignmentImport(t *testing.T) {
	t.Parallel()

	// A timestamp is added to the name to avoid failure due to previous
	// test failures
	time := time.Now().Unix()
	roleName := fmt.Sprintf("TestTermsOfUseRoleAssignmentImport Role %d", time)
	testUser := os.Getenv("ACC_TEST_SYSTEM_ROLE_OWNER")
	itShopName := "General access shop shelf"
	termsOfUse := "Test TermsOfUse"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			// Test initial creation
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "tilgangsportalen_system_role" "test" {
					name              = "%s"
					product_category  = "TBD"
					system_role_owner = "%s"
					approval_level    = "L2"
					description       = "Terraform acceptance test role for terms of use assignment."
					it_shop_name      = "%s"
				} 

				resource "tilgangsportalen_terms_of_use_role_assignment" "test" {
					role_name = tilgangsportalen_system_role.test.name
					terms_of_use_identifier = "%s"
				}
				`, roleName, testUser, itShopName, termsOfUse),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tilgangsportalen_terms_of_use_role_assignment.test", "role_name", roleName),
					resource.TestCheckResourceAttr("tilgangsportalen_terms_of_use_role_assignment.test", "terms_of_use_identifier", termsOfUse),
					resource.TestCheckResourceAttrSet("tilgangsportalen_terms_of_use_role_assignment.test", "id"),
					resource.TestCheckResourceAttrSet("tilgangsportalen_terms_of_use_role_assignment.test", "terms_of_use_description"),
				),
			},
			// Test import only
			{
				ImportState:       true,
				ResourceName:      "tilgangsportalen_terms_of_use_role_assignment.test",
				ImportStateVerify: true,
				ImportStateId:     fmt.Sprintf("%s|%s", roleName, termsOfUse),
			},
		},
	})
}

func TestTermsOfUseRoleAssignmentImportError(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			// Test import with invalid format
			{
				Config: providerConfig + `
				resource "tilgangsportalen_terms_of_use_role_assignment" "test" {
					role_name = "dummy"
					terms_of_use_identifier = "dummy"
				}
				`,
				ImportState:   true,
				ResourceName:  "tilgangsportalen_terms_of_use_role_assignment.test",
				ImportStateId: "invalid-format",
				ExpectError:   regexp.MustCompile("Invalid Import ID"),
			},
		},
	})
}

func TestTermsOfUseRoleAssignmentAlreadyAssigned(t *testing.T) {
	t.Parallel()

	// A timestamp is added to the name to avoid failure due to previous
	// test failures
	time := time.Now().Unix()
	roleName := fmt.Sprintf("TestTermsOfUseRoleAssignmentAlreadyAssigned Role %d", time)
	testUser := os.Getenv("ACC_TEST_SYSTEM_ROLE_OWNER")
	itShopName := "General access shop shelf"
	termsOfUse := "Test TermsOfUse 2"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			// First assignment - should succeed
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "tilgangsportalen_system_role" "test" {
					name              = "%s"
					product_category  = "TBD"
					system_role_owner = "%s"
					approval_level    = "L2"
					description       = "Terraform acceptance test role for terms of use assignment."
					it_shop_name      = "%s"
				} 

				resource "tilgangsportalen_terms_of_use_role_assignment" "test" {
					role_name = tilgangsportalen_system_role.test.name
					terms_of_use_identifier = "%s"
				}
				`, roleName, testUser, itShopName, termsOfUse),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tilgangsportalen_terms_of_use_role_assignment.test", "role_name", roleName),
					resource.TestCheckResourceAttr("tilgangsportalen_terms_of_use_role_assignment.test", "terms_of_use_identifier", termsOfUse),
					resource.TestCheckResourceAttrSet("tilgangsportalen_terms_of_use_role_assignment.test", "id"),
					resource.TestCheckResourceAttrSet("tilgangsportalen_terms_of_use_role_assignment.test", "terms_of_use_description"),
				),
			},
			// Assign the same terms of use with a new resource - should not fail even if already assigned
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "tilgangsportalen_system_role" "test" {
					name              = "%s"
					product_category  = "TBD"
					system_role_owner = "%s"
					approval_level    = "L2"
					description       = "Terraform acceptance test role for terms of use assignment."
					it_shop_name      = "%s"
				} 

				resource "tilgangsportalen_terms_of_use_role_assignment" "test" {
					role_name = tilgangsportalen_system_role.test.name
					terms_of_use_identifier = "%s"
				}

				resource "tilgangsportalen_terms_of_use_role_assignment" "test_same" {
					role_name = tilgangsportalen_system_role.test.name
					terms_of_use_identifier = "%s"
				}
				`, roleName, testUser, itShopName, termsOfUse, termsOfUse),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tilgangsportalen_terms_of_use_role_assignment.test", "role_name", roleName),
					resource.TestCheckResourceAttr("tilgangsportalen_terms_of_use_role_assignment.test", "terms_of_use_identifier", termsOfUse),
					resource.TestCheckResourceAttrSet("tilgangsportalen_terms_of_use_role_assignment.test", "id"),
					resource.TestCheckResourceAttrSet("tilgangsportalen_terms_of_use_role_assignment.test", "terms_of_use_description"),
					resource.TestCheckResourceAttr("tilgangsportalen_terms_of_use_role_assignment.test_same", "role_name", roleName),
					resource.TestCheckResourceAttr("tilgangsportalen_terms_of_use_role_assignment.test_same", "terms_of_use_identifier", termsOfUse),
					resource.TestCheckResourceAttrSet("tilgangsportalen_terms_of_use_role_assignment.test_same", "id"),
					resource.TestCheckResourceAttrSet("tilgangsportalen_terms_of_use_role_assignment.test_same", "terms_of_use_description"),
				),
			},
			// Remove the second assignment without destroying the resource. This allows the test to be cleaned up without errors
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "tilgangsportalen_system_role" "test" {
					name              = "%s"
					product_category  = "TBD"
					system_role_owner = "%s"
					approval_level    = "L2"
					description       = "Terraform acceptance test role for terms of use assignment."
					it_shop_name      = "%s"
				} 

				resource "tilgangsportalen_terms_of_use_role_assignment" "test" {
					role_name = tilgangsportalen_system_role.test.name
					terms_of_use_identifier = "%s"
				}

				removed {
				  from = tilgangsportalen_terms_of_use_role_assignment.test_same
				  lifecycle {
				    destroy = false
				  }
				}
				`, roleName, testUser, itShopName, termsOfUse),
			},
		},
	})
}

func TestTermsOfUseRoleAssignmentNonExistentRole(t *testing.T) {
	t.Parallel()

	time := time.Now().Unix()
	nonExistentRole := fmt.Sprintf("NonExistentRole%d", time)
	termsOfUse := "Test TermsOfUse"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			// Test assignment to non-existent role
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "tilgangsportalen_terms_of_use_role_assignment" "test" {
					role_name = "%s"
					terms_of_use_identifier = "%s"
				}
				`, nonExistentRole, termsOfUse),
				ExpectError: regexp.MustCompile("Role doesn't exist or you don't have permissions against it|Unable to assign Terms of Use"),
			},
		},
	})
}
