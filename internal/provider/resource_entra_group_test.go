package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestCreateNewEntraGroup(t *testing.T) {

	name := "[Test] test new entra group"
	newName := "[Test] test new entra group new name"
	alias := "test_new_entra_group_alias"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "tilgangsportalen_entra_group" "test" {
					name = "%s"
					alias = "%s"
					description = "terraform provider acceptance test"
					inheritance_level = "User"
				}
				`,name, alias),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tilgangsportalen_entra_group.test", "name", name),
					resource.TestCheckResourceAttr("tilgangsportalen_entra_group.test", "alias", alias),
					resource.TestCheckResourceAttr("tilgangsportalen_entra_group.test", "description", "terraform provider acceptance test"),
					resource.TestCheckResourceAttr("tilgangsportalen_entra_group.test", "inheritance_level", "User"),
				),
			},
			// test import to state using ImportStateCheckFunc
			{
				ImportState:             true,
				ResourceName:            "tilgangsportalen_entra_group.test",
				ImportStateVerifyIgnore: []string{"alias", "description", "inheritance_level"},
			},
			// test update name
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "tilgangsportalen_entra_group" "test" {
					name = "%s"
					alias = "%s"
					description = "terraform provider acceptance test"
					inheritance_level = "User"
				}
				`,newName,alias),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tilgangsportalen_entra_group.test", "name", newName),
					resource.TestCheckResourceAttr("tilgangsportalen_entra_group.test", "alias", alias),
					resource.TestCheckResourceAttr("tilgangsportalen_entra_group.test", "description", "terraform provider acceptance test"),
					resource.TestCheckResourceAttr("tilgangsportalen_entra_group.test", "inheritance_level", "User"),
				),
			},
		},
	})
}
