package provider

import (
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// Set providerconfig and get client varables from environment
var (
	providerConfig = fmt.Sprintf(`
provider "tilgangsportalen" {
	// Configuration options
	hosturl  = "%s"
	username = "%s"
	password = "%s"
}
`, os.Getenv("TF_VAR_TILGANGSPORTALEN_URL"), os.Getenv("TF_VAR_TILGANGSPORTALEN_USERNAME"), os.Getenv("TF_VAR_TILGANGSPORTALEN_PASSWORD"))

	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"tilgangsportalen": providerserver.NewProtocol6WithError(New("test")()),
	}
)
