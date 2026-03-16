package provider

import (
	"context"
	"terraform-provider-tilgangsportalen/internal/tilgangsportalapi"

	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &EntraGroupDataSource{}
var _ datasource.DataSourceWithConfigValidators = &EntraGroupDataSource{}

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
			"entitlement_uid": schema.StringAttribute{
				Description: "The entitlement UID of the Entra Group. At least one of `entitlement_uid` and `name` must be specified. This attribute will be empty if not set in the configuration.",
				Optional:    true,
			},
			"object_id": schema.StringAttribute{
				Description: "The object ID of the Entra Group",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The display name of the Entra Group. At least one of `name` and `entitlement_uid` must be specified. This attribute will be empty if not set in the configuration.",
				Optional:    true,
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

// ConfigValidators ensures at least one group lookup key is provided.
func (d *EntraGroupDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.AtLeastOneOf(
			path.MatchRoot("entitlement_uid"),
			path.MatchRoot("name"),
		),
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

	var (
		group *tilgangsportalapi.EntraGroup
		err   error
	)

	// Prefer EntitlementUID lookup when available, otherwise fall back to display name.
	if !data.EntitlementUID.IsNull() && data.EntitlementUID.ValueString() != "" {
		group, err = d.client.GetEntraGroupByID(data.EntitlementUID.ValueString(), false)
	} else {
		group, err = d.client.GetEntraGroup(data.DisplayName.ValueString(), false)
	}
	if err != nil {
		resp.Diagnostics.AddError("failed to get Entra ID group", err.Error())
		return
	}

	// Check if we need to wait for the object_id to be set for this group
	waitForObjectId := checkIfGroupWillBeCreatedInEntra(d.client, group.DisplayName, group.Description)

	// If we need to wait for the object_id to be set, wait for the value to be returned by the API
	if waitForObjectId && group.EntraIDOID == "" {
		if !data.EntitlementUID.IsNull() && data.EntitlementUID.ValueString() != "" {
			group, err = d.client.GetEntraGroupByID(data.EntitlementUID.ValueString(), waitForObjectId)
		} else {
			group, err = d.client.GetEntraGroup(data.DisplayName.ValueString(), waitForObjectId)
		}
		if err != nil {
			resp.Diagnostics.AddError("failed to get Entra ID group", err.Error())
			return
		}
	}

	// Set the resource data
	data.Id = types.StringValue(group.DisplayName)
	data.EntitlementUID = types.StringValue(group.EntitlementUID)
	data.DisplayName = types.StringValue(group.DisplayName)
	data.Description = types.StringValue(group.Description)
	data.InheritanceLevel = types.StringValue(group.InheritanceLevel)
	data.EntraIDOID = types.StringValue(group.EntraIDOID)

	// Write the resource data
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
