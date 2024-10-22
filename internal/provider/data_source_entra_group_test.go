package provider

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Create Entra ID group and read values using the data source
func TestEntraGroupDataSource(t *testing.T) {
	t.Parallel()

	// A timestamp is added to the name to avoid failure due to previous
	// test failures
	time := time.Now().Unix()
	entraGroupName := fmt.Sprintf("[Ex] TestEntraGroupDataSource Group %d", time)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "tilgangsportalen_entra_group" "test_entra_group" {
					name              = "%s"
					description       = "Terraform acceptance test Entra ID group."
					inheritance_level = "User"
				} 

				data "tilgangsportalen_entra_group" "test_entra_group" {
					name = tilgangsportalen_entra_group.test_entra_group.name
				}
				`, entraGroupName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.tilgangsportalen_entra_group.test_entra_group", "name", entraGroupName),
					resource.TestCheckResourceAttr("data.tilgangsportalen_entra_group.test_entra_group", "description", "Terraform acceptance test Entra ID group."),
					resource.TestCheckResourceAttr("data.tilgangsportalen_entra_group.test_entra_group", "inheritance_level", "User"),
				),
			},
		},
	})
}

// Test that inheritance level is handled correctly when it's set to "Admin"
func TestEntraGroupDataSource_InheritanceLevelAdmin(t *testing.T) {
	t.Parallel()

	time := time.Now().Unix()
	entraGroupName := fmt.Sprintf("[Ex] TestEntraGroupDataSource_InheritanceLevelAdmin Group %d", time)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "tilgangsportalen_entra_group" "test_entra_group_admin" {
					name              = "%s"
					description       = "Terraform acceptance test Entra ID group."
					inheritance_level = "Admin"
				} 

				data "tilgangsportalen_entra_group" "test_entra_group_admin" {
					name = tilgangsportalen_entra_group.test_entra_group_admin.name
				}
				`, entraGroupName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.tilgangsportalen_entra_group.test_entra_group_admin", "name", entraGroupName),
					resource.TestCheckResourceAttr("data.tilgangsportalen_entra_group.test_entra_group_admin", "description", "Terraform acceptance test Entra ID group."),
					resource.TestCheckResourceAttr("data.tilgangsportalen_entra_group.test_entra_group_admin", "inheritance_level", "Admin"),
				),
			},
		},
	})
}
