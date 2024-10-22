package provider

import (
	"context"
	"terraform-provider-tilgangsportalen/internal/tilgangsportalapi"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &EntraGroupDataSource{}

// Helper function
func NewEntraGroupDataSource() datasource.DataSource {
	return &EntraGroupDataSource{}
}

// Defines the data source implementation
type EntraGroupDataSource struct {
	client *tilgangsportalapi.Client
}

// Metadata returns the resource type name.
func (d *EntraGroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_entra_group"
}

// Schema defines the schema for the resource.
func (d *EntraGroupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Schema for the Entra Group data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the Entra Group",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The display name of the Entra Group.",
				Required:    true,
			},
			"alias": schema.StringAttribute{
				Description: "The alias for the Entra Group",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the Entra Group.",
				Computed:    true,
			},
			"inheritance_level": schema.StringAttribute{
				Description: "The inheritance level of the Entra Group.",
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (d *EntraGroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	ConfigureClientDataSource(ctx, req, resp, func(client *tilgangsportalapi.Client) {
		d.client = client
	})
}

// Read the resource data.
func (d *EntraGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data EntraGroupModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get the entra id group from the API
	group, err := d.client.GetEntraGroup(data.DisplayName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to get Entra ID group", err.Error())
		return
	}

	// Set the resource data
	data.Id = types.StringValue(group.DisplayName)
	data.DisplayName = types.StringValue(group.DisplayName)
	data.Description = types.StringValue(group.Description)
	data.InheritanceLevel = types.StringValue(group.InheritanceLevel)

	// Write the resource data
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
