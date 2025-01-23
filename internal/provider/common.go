package provider

import (
	"context"
	"fmt"

	"terraform-provider-tilgangsportalen/internal/tilgangsportalapi"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// CommonConfigureClient is a helper function to configure the client for both data sources and resources.
func CommonConfigureClient(ctx context.Context, providerData interface{}, diagnostics *diag.Diagnostics, setClient func(*tilgangsportalapi.Client)) {
	// Prevent panic if the provider has not been configured.
	if providerData == nil {
		return
	}

	client, ok := providerData.(*tilgangsportalapi.Client)
	if !ok {
		diagnostics.AddError(
			"Unexpected Configure Type",
			fmt.Sprintf("Expected *tilgangsportalapi.Client, got: %T. Please report this issue to the provider developers.",
				providerData),
		)
		return
	}

	setClient(client)
}

// ConfigureClientDataSource is a helper function to configure the client for a data source.
func ConfigureClientDataSource(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse, setClient func(*tilgangsportalapi.Client)) {
	CommonConfigureClient(ctx, req.ProviderData, &resp.Diagnostics, setClient)
}

// ConfigureClientResource is a helper function to configure the client for a resource.
func ConfigureClientResource(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse, setClient func(*tilgangsportalapi.Client)) {
	CommonConfigureClient(ctx, req.ProviderData, &resp.Diagnostics, setClient)
}
