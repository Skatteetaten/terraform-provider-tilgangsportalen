package provider

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestCreateNewSystemRole(t *testing.T) {
	t.Parallel()

	// A timestamp is added to the name to avoid failure due to previous
	// test failures
	timeNow := time.Now().Unix()
	name := fmt.Sprintf("Test-Create_New_System_Role Role %d", timeNow)
	newName := name + " new name"
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
	newDescription := `<<-EOT
	Terraform_new_'acceptance'\n-_øåæ\tØÅÆ
	'!#$x%&/
	()[]{}'?!=(a){b}[c]@^*<>:,;.|
	\"test\" \t\"r\"ole
	EOT`
	expectedDescription := "Terraform_'acceptance'\\n-_øåæ\\tØÅÆ\n'!#$x%&/\n()[]{}'?!=(a){b}[c]@^*<>:,;.|\n\\\"test\\\" \\t\\\"r\\\"ole\n"
	expectedNewDescription := "Terraform_new_'acceptance'\\n-_øåæ\\tØÅÆ\n'!#$x%&/\n()[]{}'?!=(a){b}[c]@^*<>:,;.|\n\\\"test\\\" \\t\\\"r\\\"ole\n"
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
				),
			},
			// test import to state using ImportStateCheckFunc
			{
				ImportState:             true,
				ResourceName:            "tilgangsportalen_system_role.test_role",
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"it_shop_name"},
			},
			// test update name
			{
				Config: providerConfig + fmt.Sprintf(`		  
				resource "tilgangsportalen_system_role" "test_role" {
					name              = "%s"
					product_category  = "%s"
					system_role_owner = "%s"
					approval_level    = "%s"
					description       = %s
				} 
				`, newName, productCategory, roleOwner, approvalLevel, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "name", newName),
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "product_category", productCategory),
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "system_role_owner", roleOwner),
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "approval_level", approvalLevel),
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "description", expectedDescription),
				),
			},
			// test update description
			{
				Config: providerConfig + fmt.Sprintf(`		  
				resource "tilgangsportalen_system_role" "test_role" {
					name              = "%s"
					product_category  = "%s"
					system_role_owner = "%s"
					approval_level    = "%s"
					description       = %s
				} 
				`, newName, productCategory, roleOwner, approvalLevel, newDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "name", newName),
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "product_category", productCategory),
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "system_role_owner", roleOwner),
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "approval_level", approvalLevel),
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "description", expectedNewDescription),
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

				data "tilgangsportalen_system_roles" "all_roles" {}
				`, name, productCategory, roleOwner, securityOwner, approvalLevel, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "name", name),
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "product_category", productCategory),
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "system_role_owner", roleOwner),
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "system_role_security_owner", securityOwner),
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "approval_level", approvalLevel),
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "description", description),
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
