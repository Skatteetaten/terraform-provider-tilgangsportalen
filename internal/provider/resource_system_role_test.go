package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestCreateNewSystemRole(t *testing.T) {
	testUser := os.Getenv("TF_VAR_TEST_USER")
	name := "system role acceptance test"
	newName := "system role acceptance test new name"
	productCategory := "TBD"
	approvalLevel := "L2"
	description := "Terraform acceptance test role."
	newDescription := "Terraform acceptance test role new description."
	itShopName := "Access shop shelf"
	
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
				`,name,productCategory,testUser,approvalLevel,description,itShopName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "name", name),
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "product_category", productCategory),
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "system_role_owner", testUser),
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "approval_level", approvalLevel),
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "description", description),
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "it_shop_name", itShopName),
				),
			},
			// test import to state using ImportStateCheckFunc
			{
				ImportState:             true,
				ResourceName:            "tilgangsportalen_system_role.test_role",
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
					description       = "%s"
					it_shop_name      = "%s"
				} 
				`,newName,productCategory,testUser,approvalLevel,description,itShopName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "name", newName),
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "product_category", productCategory),
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "system_role_owner", testUser),
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "approval_level", approvalLevel),
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "description", description),
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "it_shop_name", itShopName),
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
					description       = "%s"
					it_shop_name      = "%s"
				} 
				`,newName,productCategory,testUser,approvalLevel,newDescription,itShopName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "name", newName),
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "product_category", productCategory),
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "system_role_owner", testUser),
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "approval_level", approvalLevel),
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "description", newDescription),
					resource.TestCheckResourceAttr("tilgangsportalen_system_role.test_role", "it_shop_name", itShopName),
				),
			},
		},
	})
}
