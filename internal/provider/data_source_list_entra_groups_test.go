package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestEntraGroupsDataSource(t *testing.T) {

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
				data "tilgangsportalen_entra_groups" "all_groups" {}
				`,
				Check: resource.ComposeTestCheckFunc(),
			},
		},
	})
}
