package provider

import (
	"context"
	"fmt"

	"terraform-provider-tilgangsportalen/internal/tilgangsportalapi"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// ConfigureClientDataSource is a helper function to configure the client for a data source.
func ConfigureClientDataSource(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse, setClient func(*tilgangsportalapi.Client)) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		if resp != nil {
			return
		}
	}

	client, ok := req.ProviderData.(*tilgangsportalapi.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexcepted Configure Type",
			fmt.Sprintf("Expected *tilgangsportalapi.Client, got: %T. Please report this issue to the provider developers.",
				req.ProviderData),
		)
		return
	}

	setClient(client)
}

func ConfigureClientResource(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse, setClient func(*tilgangsportalapi.Client)) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		if resp != nil {
			return
		}
	}

	client, ok := req.ProviderData.(*tilgangsportalapi.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexcepted Configure Type",
			fmt.Sprintf("Expected *tilgangsportalapi.Client, got: %T. Please report this issue to the provider developers.",
				req.ProviderData),
		)
		return
	}

	setClient(client)
}
