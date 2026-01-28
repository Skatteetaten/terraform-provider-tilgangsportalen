package provider

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// create system role and read values using data source
func TestSystemRoleDataSource(t *testing.T) {
	t.Parallel()

	// A timestamp is added to the name to avoid failure due to previous
	// test failures
	time := time.Now().Unix()
	roleName := fmt.Sprintf("TestSystemRoleDataSource Role %d", time)
	roleOwner := os.Getenv("ACC_TEST_SYSTEM_ROLE_OWNER")
	itShopName := "General access shop shelf"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "tilgangsportalen_system_role" "test_role_data_source" {
					name              = "%s"
					product_category  = "TBD"
					system_role_owner = "%s"
					approval_level    = "L2"
					description       = "Terraform acceptance test role for assignment."
					it_shop_name      = "%s"

				} 

				data "tilgangsportalen_system_role" "system_role" {
					name = tilgangsportalen_system_role.test_role_data_source.name
				}
				`, roleName, roleOwner, itShopName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.tilgangsportalen_system_role.system_role", "name", roleName),
					resource.TestCheckResourceAttr("data.tilgangsportalen_system_role.system_role", "system_role_owner", roleOwner),
					resource.TestCheckResourceAttr("data.tilgangsportalen_system_role.system_role", "approval_level", "L2"),
					resource.TestCheckResourceAttr("data.tilgangsportalen_system_role.system_role", "description", "Terraform acceptance test role for assignment."),
					resource.TestCheckResourceAttrSet("data.tilgangsportalen_system_role.system_role", "object_id"),
				),
			},
		},
	})
}

func TestSystemRoleDataSourceLookupByObjectID(t *testing.T) {
	t.Parallel()

	// A timestamp is added to the name to avoid failure due to previous
	// test failures
	time := time.Now().Unix()
	roleName := fmt.Sprintf("TestSystemRoleDataSourceLookupByObjectID Role %d", time)
	roleOwner := os.Getenv("ACC_TEST_SYSTEM_ROLE_OWNER")
	itShopName := "General access shop shelf"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "tilgangsportalen_system_role" "test_role_data_source" {
					name              = "%s"
					product_category  = "TBD"
					system_role_owner = "%s"
					approval_level    = "L2"
					description       = "Terraform acceptance test role for assignment."
					it_shop_name      = "%s"
				} 

				data "tilgangsportalen_system_role" "system_role" {
					object_id = tilgangsportalen_system_role.test_role_data_source.object_id
				}
				`, roleName, roleOwner, itShopName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.tilgangsportalen_system_role.system_role", "name", roleName),
					resource.TestCheckResourceAttr("data.tilgangsportalen_system_role.system_role", "system_role_owner", roleOwner),
					resource.TestCheckResourceAttr("data.tilgangsportalen_system_role.system_role", "approval_level", "L2"),
					resource.TestCheckResourceAttr("data.tilgangsportalen_system_role.system_role", "description", "Terraform acceptance test role for assignment."),
					resource.TestCheckResourceAttrSet("data.tilgangsportalen_system_role.system_role", "object_id"),
				),
			},
		},
	})
}

// create system role and read values using data source
func TestSystemRoleDataSourceWithL3(t *testing.T) {
	t.Parallel()

	// A timestamp is added to the name to avoid failure due to previous
	// test failures
	time := time.Now().Unix()
	roleName := fmt.Sprintf("TestSystemRoleDataSourceWithL3 Role %d", time)
	roleOwner := os.Getenv("ACC_TEST_SYSTEM_ROLE_OWNER")
	securityOwner := os.Getenv("ACC_TEST_SYSTEM_ROLE_SECURITY_OWNER")
	itShopName := "General access shop shelf"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "tilgangsportalen_system_role" "test_role_data_source" {
				name                       = "%s"
				product_category           = "TBD"
				system_role_owner          = "%s"
				system_role_security_owner = "%s"
				approval_level             = "L3"
				description                = "Terraform acceptance test role for assignment."
				it_shop_name               = "%s"
				}

				data "tilgangsportalen_system_role" "system_role" {
					name = tilgangsportalen_system_role.test_role_data_source.name
				}
				`, roleName, roleOwner, securityOwner, itShopName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.tilgangsportalen_system_role.system_role", "name", roleName),
					resource.TestCheckResourceAttr("data.tilgangsportalen_system_role.system_role", "system_role_owner", roleOwner),
					resource.TestCheckResourceAttr("data.tilgangsportalen_system_role.system_role", "approval_level", "L3"),
					resource.TestCheckResourceAttr("data.tilgangsportalen_system_role.system_role", "description", "Terraform acceptance test role for assignment."),
					resource.TestCheckResourceAttrSet("data.tilgangsportalen_system_role.system_role", "object_id"),
				),
			},
		},
	})
}
