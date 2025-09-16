package provider

import (
	"fmt"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestTermsOfUseRoleAssignmentDataSource(t *testing.T) {
	t.Parallel()

	// Create a unique role name and get test environment variables
	time := time.Now().Unix()
	roleName := fmt.Sprintf("TestTermsOfUseRoleAssignmentDataSource Role %d", time)
	testUser := os.Getenv("ACC_TEST_SYSTEM_ROLE_OWNER")
	itShopName := "General access shop shelf"
	termsOfUse := "Test TermsOfUse"

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
					description       = "Terraform acceptance test role for Terms of Use data source testing."
					it_shop_name      = "%s"
				}

				resource "tilgangsportalen_terms_of_use_role_assignment" "test" {
					role_name = tilgangsportalen_system_role.test.name
					terms_of_use_identifier = "%s"
				}

				data "tilgangsportalen_terms_of_use_role_assignment" "test" {
					role_name = tilgangsportalen_system_role.test.name
					depends_on = [tilgangsportalen_terms_of_use_role_assignment.test]
				}
				`, roleName, testUser, itShopName, termsOfUse),
				Check: resource.ComposeTestCheckFunc(
					// Check that the data source returns the expected values
					resource.TestCheckResourceAttr("data.tilgangsportalen_terms_of_use_role_assignment.test", "role_name", roleName),
					resource.TestCheckResourceAttr("data.tilgangsportalen_terms_of_use_role_assignment.test", "terms_of_use_identifier", termsOfUse),
					resource.TestCheckResourceAttrSet("data.tilgangsportalen_terms_of_use_role_assignment.test", "id"),
					resource.TestCheckResourceAttrSet("data.tilgangsportalen_terms_of_use_role_assignment.test", "terms_of_use_description"),
					resource.TestCheckResourceAttrSet("data.tilgangsportalen_terms_of_use_role_assignment.test", "terms_of_use_uid"),

					// Verify the ID format is correct
					resource.TestCheckResourceAttr("data.tilgangsportalen_terms_of_use_role_assignment.test", "id", fmt.Sprintf("%s|%s", roleName, termsOfUse)),

					// Cross-check with the resource values
					resource.TestCheckResourceAttrPair("data.tilgangsportalen_terms_of_use_role_assignment.test", "role_name", "tilgangsportalen_terms_of_use_role_assignment.test", "role_name"),
					resource.TestCheckResourceAttrPair("data.tilgangsportalen_terms_of_use_role_assignment.test", "terms_of_use_identifier", "tilgangsportalen_terms_of_use_role_assignment.test", "terms_of_use_identifier"),
					resource.TestCheckResourceAttrPair("data.tilgangsportalen_terms_of_use_role_assignment.test", "terms_of_use_description", "tilgangsportalen_terms_of_use_role_assignment.test", "terms_of_use_description"),
				),
			},
		},
	})
}

func TestTermsOfUseRoleAssignmentDataSourceNoAssignment(t *testing.T) {
	t.Parallel()

	// Test that the data source returns an error when no Terms of Use is assigned
	time := time.Now().Unix()
	roleName := fmt.Sprintf("TestTermsOfUseRoleAssignmentDataSourceNoAssignment Role %d", time)
	testUser := os.Getenv("ACC_TEST_SYSTEM_ROLE_OWNER")
	itShopName := "General access shop shelf"

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
					description       = "Terraform acceptance test role with no Terms of Use assignment."
					it_shop_name      = "%s"
				}

				data "tilgangsportalen_terms_of_use_role_assignment" "test" {
					role_name = tilgangsportalen_system_role.test.name
				}
				`, roleName, testUser, itShopName),
				ExpectError: regexp.MustCompile("No Terms of Use Assignment Found"),
			},
		},
	})
}

func TestTermsOfUseRoleAssignmentDataSourceNonExistentRole(t *testing.T) {
	t.Parallel()

	// Test that the data source returns an error for non-existent roles
	time := time.Now().Unix()
	roleName := fmt.Sprintf("NonExistentRole%d", time)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
				data "tilgangsportalen_terms_of_use_role_assignment" "test" {
					role_name = "%s"
				}
				`, roleName),
				ExpectError: regexp.MustCompile("API Error|No Terms of Use Assignment Found"),
			},
		},
	})
}

func TestTermsOfUseRoleAssignmentDataSourceAfterUpdate(t *testing.T) {
	t.Parallel()

	// Test that the data source reflects updates to assignments
	time := time.Now().Unix()
	roleName := fmt.Sprintf("TestTermsOfUseRoleAssignmentDataSourceUpdate Role %d", time)
	testUser := os.Getenv("ACC_TEST_SYSTEM_ROLE_OWNER")
	itShopName := "General access shop shelf"
	termsOfUse1 := "Test TermsOfUse"
	termsOfUse2 := "Test TermsOfUse 2"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			// Step 1: Create assignment with first Terms of Use
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "tilgangsportalen_system_role" "test" {
					name              = "%s"
					product_category  = "TBD"
					system_role_owner = "%s"
					approval_level    = "L2"
					description       = "Terraform acceptance test role for Terms of Use data source update testing."
					it_shop_name      = "%s"
				}

				resource "tilgangsportalen_terms_of_use_role_assignment" "test" {
					role_name = tilgangsportalen_system_role.test.name
					terms_of_use_identifier = "%s"
				}

				data "tilgangsportalen_terms_of_use_role_assignment" "test" {
					role_name = tilgangsportalen_system_role.test.name
					depends_on = [tilgangsportalen_terms_of_use_role_assignment.test]
				}
				`, roleName, testUser, itShopName, termsOfUse1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.tilgangsportalen_terms_of_use_role_assignment.test", "terms_of_use_identifier", termsOfUse1),
					resource.TestCheckResourceAttr("data.tilgangsportalen_terms_of_use_role_assignment.test", "id", fmt.Sprintf("%s|%s", roleName, termsOfUse1)),
				),
			},
			// Step 2: Update assignment to second Terms of Use
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "tilgangsportalen_system_role" "test" {
					name              = "%s"
					product_category  = "TBD"
					system_role_owner = "%s"
					approval_level    = "L2"
					description       = "Terraform acceptance test role for Terms of Use data source update testing."
					it_shop_name      = "%s"
				}

				resource "tilgangsportalen_terms_of_use_role_assignment" "test" {
					role_name = tilgangsportalen_system_role.test.name
					terms_of_use_identifier = "%s"
				}

				data "tilgangsportalen_terms_of_use_role_assignment" "test" {
					role_name = tilgangsportalen_system_role.test.name
					depends_on = [tilgangsportalen_terms_of_use_role_assignment.test]
				}
				`, roleName, testUser, itShopName, termsOfUse2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.tilgangsportalen_terms_of_use_role_assignment.test", "terms_of_use_identifier", termsOfUse2),
					resource.TestCheckResourceAttr("data.tilgangsportalen_terms_of_use_role_assignment.test", "id", fmt.Sprintf("%s|%s", roleName, termsOfUse2)),
				),
			},
		},
	})
}

func TestTermsOfUseRoleAssignmentDataSourceImport(t *testing.T) {
	t.Parallel()

	// Test that the data source works with imported assignments
	time := time.Now().Unix()
	roleName := fmt.Sprintf("TestTermsOfUseRoleAssignmentDataSourceImport Role %d", time)
	testUser := os.Getenv("ACC_TEST_SYSTEM_ROLE_OWNER")
	itShopName := "General access shop shelf"
	termsOfUse := "Test TermsOfUse"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			// Step 1: Create assignment through resource
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "tilgangsportalen_system_role" "test" {
					name              = "%s"
					product_category  = "TBD"
					system_role_owner = "%s"
					approval_level    = "L2"
					description       = "Terraform acceptance test role for Terms of Use data source import testing."
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
			// Step 2: Test import and verify data source can read imported assignment
			{
				ResourceName:      "tilgangsportalen_terms_of_use_role_assignment.test",
				ImportState:       true,
				ImportStateId:     fmt.Sprintf("%s|%s", roleName, termsOfUse),
				ImportStateVerify: true,
			},
			// Step 3: Use data source with imported resource
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "tilgangsportalen_system_role" "test" {
					name              = "%s"
					product_category  = "TBD"
					system_role_owner = "%s"
					approval_level    = "L2"
					description       = "Terraform acceptance test role for Terms of Use data source import testing."
					it_shop_name      = "%s"
				}

				resource "tilgangsportalen_terms_of_use_role_assignment" "test" {
					role_name = tilgangsportalen_system_role.test.name
					terms_of_use_identifier = "%s"
				}

				data "tilgangsportalen_terms_of_use_role_assignment" "test" {
					role_name = tilgangsportalen_system_role.test.name
					depends_on = [tilgangsportalen_terms_of_use_role_assignment.test]
				}
				`, roleName, testUser, itShopName, termsOfUse),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.tilgangsportalen_terms_of_use_role_assignment.test", "role_name", roleName),
					resource.TestCheckResourceAttr("data.tilgangsportalen_terms_of_use_role_assignment.test", "terms_of_use_identifier", termsOfUse),
					resource.TestCheckResourceAttrSet("data.tilgangsportalen_terms_of_use_role_assignment.test", "terms_of_use_description"),
					resource.TestCheckResourceAttrSet("data.tilgangsportalen_terms_of_use_role_assignment.test", "terms_of_use_uid"),
				),
			},
		},
	})
}
