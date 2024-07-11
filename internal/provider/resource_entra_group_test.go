package provider

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestCreateNewEntraGroup(t *testing.T) {
	t.Parallel()

	// A timestamp is added to the name to avoid failure due to previous
	// test failures
	time := time.Now().Unix()
	name := fmt.Sprintf("[Group] Test-Create_New_Entra_Group %d", time)
	newName := name + " new name"
	alias := "TestCreateNewEntraGroup_alias"

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
				`, name, alias),
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
				`, newName, alias),
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
