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
				resource "tilgangsportalen_entra_group" "test" {
					name = "%s"
					alias = "%s"
					description = %s
					inheritance_level = "User"
				}
				`, name, alias, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tilgangsportalen_entra_group.test", "name", name),
					resource.TestCheckResourceAttr("tilgangsportalen_entra_group.test", "alias", alias),
					resource.TestCheckResourceAttr("tilgangsportalen_entra_group.test", "description", expectedDescription),
					resource.TestCheckResourceAttr("tilgangsportalen_entra_group.test", "inheritance_level", "User"),
				),
			},
			// test import to state using ImportStateCheckFunc
			{
				ImportState:             true,
				ResourceName:            "tilgangsportalen_entra_group.test",
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"alias", "description", "inheritance_level"},
			},
			// test update name
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "tilgangsportalen_entra_group" "test" {
					name = "%s"
					alias = "%s"
					description = %s
					inheritance_level = "User"
				}
				`, newName, alias, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tilgangsportalen_entra_group.test", "name", newName),
					resource.TestCheckResourceAttr("tilgangsportalen_entra_group.test", "alias", alias),
					resource.TestCheckResourceAttr("tilgangsportalen_entra_group.test", "description", expectedDescription),
					resource.TestCheckResourceAttr("tilgangsportalen_entra_group.test", "inheritance_level", "User"),
				),
			},
		},
	})
}
