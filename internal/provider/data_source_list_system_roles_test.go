package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestNewSystemRolesDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
				data "tilgangsportalen_system_roles" "all_roles" {}
				`,
				Check: resource.ComposeTestCheckFunc(),
			},
		},
	})
}
