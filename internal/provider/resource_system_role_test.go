package provider

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestCreateNewSystemRole(t *testing.T) {
	t.Parallel()

	// A timestamp is added to the name to avoid failure due to previous
	// test failures
	timeNow := time.Now().Unix()
	name := fmt.Sprintf("Test-Create_New_System_Role Role %d", timeNow)
	roleOwner := os.Getenv("ACC_TEST_SYSTEM_ROLE_OWNER")
	productCategory := "TBD"
	approvalLevel := "L2"
	itShopName := "General access shop shelf"

	// Testing special characters and line breaks in the group description.
	// The number of "\" characters changes in the expected description due to the use of GO´s Raw string literals in the input description.
	description := `<<-EOT
	Terraform_'acceptance'\n-_øåæ\tØÅÆ
	'!#$x%&/
	()[]{}'?!=(a){b}[c]@^*<>:,;.|
	\"test\" \t\"r\"ole
	EOT`
	expectedDescription := "Terraform_'acceptance'\\n-_øåæ\\tØÅÆ\n'!#$x%&/\n()[]{}'?!=(a){b}[c]@^*<>:,;.|\n\\\"test\\\" \\t\\\"r\\\"ole\n"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`		  
				resource "tilgangsportalen_system_role" "test_role" {
					name              = "%s"
					product_category  = "%s"
					system_role_owner = "%s"
					approval_level    = "%s"
					description       = %s
					it_shop_name      = "%s"
				} 

				data "tilgangsportalen_system_roles" "all_roles" {}
				`, name, productCategory, roleOwner, approvalLevel, description, itShopName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "name", name),
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "product_category", productCategory),
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "system_role_owner", roleOwner),
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "approval_level", approvalLevel),
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "description", expectedDescription),
					resource.TestCheckResourceAttrSet("tilgangsportalen_system_role.test_role", "object_id"),
				),
			},
		},
	})
}

func TestUpdateSystemRoleDescription(t *testing.T) {
	t.Parallel()

	timeNow := time.Now().Unix()
	name := fmt.Sprintf("Test-Update_Description Role %d", timeNow)
	roleOwner := os.Getenv("ACC_TEST_SYSTEM_ROLE_OWNER")
	productCategory := "TBD"
	approvalLevel := "L2"
	itShopName := "General access shop shelf"

	// Testing special characters and line breaks in the group description.
	// The number of "\" characters changes in the expected description due to the use of GO´s Raw string literals in the input description.
	description := `<<-EOT
	Terraform_'acceptance'\n-_øåæ\tØÅÆ
	'!#$x%&/
	()[]{}'?!=(a){b}[c]@^*<>:,;.|
	\"test\" \t\"r\"ole
	EOT`
	expectedDescription := "Terraform_'acceptance'\\n-_øåæ\\tØÅÆ\n'!#$x%&/\n()[]{}'?!=(a){b}[c]@^*<>:,;.|\n\\\"test\\\" \\t\\\"r\\\"ole\n"

	newDescription := `<<-EOT
	Terraform_new_'acceptance'\n-_øåæ\tØÅÆ
	'!#$x%&/
	()[]{}'?!=(a){b}[c]@^*<>:,;.|
	\"test\" \t\"r\"ole
	EOT`
	expectedNewDescription := "Terraform_new_'acceptance'\\n-_øåæ\\tØÅÆ\n'!#$x%&/\n()[]{}'?!=(a){b}[c]@^*<>:,;.|\n\\\"test\\\" \\t\\\"r\\\"ole\n"

	var originalObjectID string

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`		  
				resource "tilgangsportalen_system_role" "test_role" {
					name              = "%s"
					product_category  = "%s"
					system_role_owner = "%s"
					approval_level    = "%s"
					description       = %s
					it_shop_name      = "%s"
				} 
				`, name, productCategory, roleOwner, approvalLevel, description, itShopName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "name", name),
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "product_category", productCategory),
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "system_role_owner", roleOwner),
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "approval_level", approvalLevel),
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "description", expectedDescription),
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "it_shop_name", itShopName),
					resource.TestCheckResourceAttrSet("tilgangsportalen_system_role.test_role", "object_id"),

					// Retrieve and store the original ObjectID for later comparison
					func(s *terraform.State) error {
						rs := s.RootModule().Resources["tilgangsportalen_system_role.test_role"]
						originalObjectID = rs.Primary.Attributes["object_id"]
						return nil
					},
				),
			},
			// test update description only (name stays the same)
			{
				Config: providerConfig + fmt.Sprintf(`		  
				resource "tilgangsportalen_system_role" "test_role" {
					name              = "%s"
					product_category  = "%s"
					system_role_owner = "%s"
					approval_level    = "%s"
					description       = %s
					it_shop_name      = "%s"
				} 
				`, name, productCategory, roleOwner, approvalLevel, newDescription, itShopName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "name", name),
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "description", expectedNewDescription),

					// Verify ObjectID has not changed
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources["tilgangsportalen_system_role.test_role"]
						if !ok {
							return fmt.Errorf("not found: %s", "tilgangsportalen_system_role.test_role")
						}
						updatedObjectID := rs.Primary.Attributes["object_id"]
						if originalObjectID != updatedObjectID {
							return fmt.Errorf("object_id changed after description update: original %s, updated %s", originalObjectID, updatedObjectID)
						}
						return nil
					},
				),
			},
		},
	})
}

func TestUpdateSystemRoleName(t *testing.T) {
	t.Parallel()

	// A timestamp is added to the name to avoid failure due to previous
	// test failures
	timeNow := time.Now().Unix()
	name := fmt.Sprintf("Test-Update_System_Role_Name Role %d", timeNow)
	newName := name + " new name"
	roleOwner := os.Getenv("ACC_TEST_SYSTEM_ROLE_OWNER")
	productCategory := "TBD"
	approvalLevel := "L2"
	itShopName := "General access shop shelf"

	description := "Description for role rename test"

	var originalObjectID string

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`		  
				resource "tilgangsportalen_system_role" "test_role" {
					name              = "%s"
					product_category  = "%s"
					system_role_owner = "%s"
					approval_level    = "%s"
					description       = "%s"
					it_shop_name      = "%s"
				} 
				`, name, productCategory, roleOwner, approvalLevel, description, itShopName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "name", name),
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "product_category", productCategory),
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "system_role_owner", roleOwner),
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "approval_level", approvalLevel),
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "description", description),
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "it_shop_name", itShopName),
					resource.TestCheckResourceAttrSet("tilgangsportalen_system_role.test_role", "object_id"),

					// Retrieve and store the original ObjectID for later comparison
					func(s *terraform.State) error {
						rs := s.RootModule().Resources["tilgangsportalen_system_role.test_role"]
						originalObjectID = rs.Primary.Attributes["object_id"]
						return nil
					},
				),
			},
			// test update name
			{
				Config: providerConfig + fmt.Sprintf(`		  
				resource "tilgangsportalen_system_role" "test_role" {
					name              = "%s"
					product_category  = "%s"
					system_role_owner = "%s"
					approval_level    = "%s"
					description       = "%s"
					it_shop_name      = "%s"
				} 
				`, newName, productCategory, roleOwner, approvalLevel, description, itShopName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "name", newName),
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "product_category", productCategory),
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "system_role_owner", roleOwner),
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "approval_level", approvalLevel),
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "description", description),
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "it_shop_name", itShopName),
					resource.TestCheckResourceAttrSet("tilgangsportalen_system_role.test_role", "object_id"),

					// Verify that ObjectID has NOT changed after rename.
					// The RenameSystemRole API preserves the ObjectID.
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources["tilgangsportalen_system_role.test_role"]
						if !ok {
							return fmt.Errorf("not found: %s", "tilgangsportalen_system_role.test_role")
						}
						updatedObjectID := rs.Primary.Attributes["object_id"]
						if updatedObjectID == "" {
							return fmt.Errorf("object_id is empty after rename")
						}
						if originalObjectID != updatedObjectID {
							return fmt.Errorf("object_id changed after rename: original %s, updated %s", originalObjectID, updatedObjectID)
						}
						return nil
					},
				),
			},
		},
	})
}

func TestCreateNewSystemRoleWithApprovalLevelL3(t *testing.T) {
	t.Parallel()

	// A timestamp is added to the name to avoid failure due to previous
	// test failures
	timeNow := time.Now().Unix()
	name := fmt.Sprintf("Test-Create_New_System_Role_L3 Role %d", timeNow)
	roleOwner := os.Getenv("ACC_TEST_SYSTEM_ROLE_OWNER")
	securityOwner := os.Getenv("ACC_TEST_SYSTEM_ROLE_SECURITY_OWNER")
	productCategory := "TBD"
	approvalLevel := "L3"

	description := "Test role with approval level L3"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`		  
				resource "tilgangsportalen_system_role" "test_role" {
					name                        = "%s"
					product_category            = "%s"
					system_role_owner           = "%s"
					system_role_security_owner  = "%s"
					approval_level              = "%s"
					description                 = "%s"
				} 
				`, name, productCategory, roleOwner, securityOwner, approvalLevel, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "name", name),
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "product_category", productCategory),
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "system_role_owner", roleOwner),
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "system_role_security_owner", securityOwner),
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "approval_level", approvalLevel),
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "description", description),
					resource.TestCheckResourceAttrSet("tilgangsportalen_system_role.test_role", "object_id"),
				),
			},
			// test import to state using ImportStateCheckFunc
			{
				ImportState:             true,
				ResourceName:            "tilgangsportalen_system_role.test_role",
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"it_shop_name"},
			},
		},
	})
}
