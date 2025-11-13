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
	entraGroupName := fmt.Sprintf("[APPTEST] TestEntraGroupDataSource Group %d", time)
	description := "APPTEST"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "tilgangsportalen_entra_group" "test_entra_group" {
					name              = "%s"
					description       = "%s"
					inheritance_level = "User"
				}

				data "tilgangsportalen_entra_group" "test_entra_group" {
					name = tilgangsportalen_entra_group.test_entra_group.name
				}
				`, entraGroupName, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.tilgangsportalen_entra_group.test_entra_group", "name", entraGroupName),
					resource.TestCheckResourceAttr("data.tilgangsportalen_entra_group.test_entra_group", "description", description),
					resource.TestCheckResourceAttr("data.tilgangsportalen_entra_group.test_entra_group", "inheritance_level", "User"),
					resource.TestCheckResourceAttrSet("data.tilgangsportalen_entra_group.test_entra_group", "object_id"),

					// Cross-check with the resource values
					resource.TestCheckResourceAttrPair("data.tilgangsportalen_entra_group.test_entra_group", "object_id", "tilgangsportalen_entra_group.test_entra_group", "object_id"),
				),
			},
		},
	})
}

// Test that inheritance level is handled correctly when it's set to "Admin"
func TestEntraGroupDataSource_InheritanceLevelAdmin(t *testing.T) {
	t.Parallel()

	time := time.Now().Unix()
	entraGroupName := fmt.Sprintf("[APPTEST] TestEntraGroupDataSource_InheritanceLevelAdmin Group %d", time)
	description := "APPTEST"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "tilgangsportalen_entra_group" "test_entra_group_admin" {
					name              = "%s"
					description       = "%s"
					inheritance_level = "Admin"
				}

				data "tilgangsportalen_entra_group" "test_entra_group_admin" {
					name = tilgangsportalen_entra_group.test_entra_group_admin.name
				}
				`, entraGroupName, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.tilgangsportalen_entra_group.test_entra_group_admin", "name", entraGroupName),
					resource.TestCheckResourceAttr("data.tilgangsportalen_entra_group.test_entra_group_admin", "description", description),
					resource.TestCheckResourceAttr("data.tilgangsportalen_entra_group.test_entra_group_admin", "inheritance_level", "Admin"),
					resource.TestCheckResourceAttrSet("data.tilgangsportalen_entra_group.test_entra_group_admin", "object_id"),
				),
			},
		},
	})
}
