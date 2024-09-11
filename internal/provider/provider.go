// Package provider implements the Tilgangportalen provider
package provider

import (
	"context"
	"fmt"
	"os"
	"strings"
	"terraform-provider-tilgangsportalen/internal/tilgangsportalapi"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure TilgangsportalenProvider satisfies various provider interfaces.
var _ provider.Provider = &TilgangsportalenProvider{}

// TilgangsportalenProvider defines the provider implementation.
type TilgangsportalenProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// TilgangsportalenProviderModel describes the provider data model.
type TilgangsportalenProviderModel struct {
	HostURL  types.String `tfsdk:"hosturl"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

// Metadata returns the provider type name.
func (p *TilgangsportalenProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "tilgangsportalen"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *TilgangsportalenProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This provider is used to create roles and security groups in Tilgangsportalen",
		Attributes: map[string]schema.Attribute{
			"hosturl": schema.StringAttribute{
				MarkdownDescription: "Tilgangsportalen host url",
				Optional:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Tilgangsportalen username",
				Optional:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Tilgangsportalen password",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

// Configure sets up the Tilgangsportalen provider client
func (p *TilgangsportalenProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {

	// Check environment variables
	urlFormatString := strings.Trim(os.Getenv("TILGANGSPORTALEN_URL"), "\"")
	password := strings.Trim(os.Getenv("TILGANGSPORTALEN_PASSWORD"), "\"")
	username := strings.Trim(os.Getenv("TILGANGSPORTALEN_USERNAME"), "\"")
	var data TilgangsportalenProviderModel

	tflog.Info(ctx, "Configuring Tilgangsportalen client")
	tflog.Debug(ctx, "Got HostURL from environment variable: %s")

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.HostURL.ValueString() != "" {
		tflog.Debug(ctx, fmt.Sprintf("HostURL overridden by Terraform configuration variable: %s", data.HostURL.ValueString()))
		urlFormatString = strings.Trim(data.HostURL.ValueString(), "\"")
	}
	if data.Password.ValueString() != "" {
		password = strings.Trim(data.Password.ValueString(), "\"")
	}
	if data.Username.ValueString() != "" {
		username = strings.Trim(data.Username.ValueString(), "\"")
	}
	ctx = tflog.SetField(ctx, "tilgangsportalen_host", urlFormatString)
	ctx = tflog.SetField(ctx, "tilgangsportalen_username", username)
	ctx = tflog.SetField(ctx, "tilgangsportalen_password", password)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "tilgangsportalen_password")

	tflog.Info(ctx, "Creating Tilgangsportalen client")

	client, err := tilgangsportalapi.NewClient(urlFormatString, username, password)
	resp.DataSourceData = client
	resp.ResourceData = client

	if err != nil {
		resp.Diagnostics.AddError(err.Error(), "Error when creating new client. Make sure variables for host, username, and password are correct.")
	}
}

// Resources defines the resources implemented in the provider.
func (p *TilgangsportalenProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		CreateNewSystemRole,
		CreateNewEntraGroup,
		CreateNewEntraGroupRoleAssignment,
	}
}

// DataSources defines the data sources implemented in the provider.
func (p *TilgangsportalenProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewSystemRolesDataSource,
		NewEntraGroupsDataSource,
		NewEntraGroupsForRoleDataSource,
		NewSystemRoleDataSource,
	}
}

// New is a helper function
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &TilgangsportalenProvider{
			version: version,
		}
	}
}
