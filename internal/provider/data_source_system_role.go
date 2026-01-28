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
var _ datasource.DataSource = &SystemRoleDataSource{}
var _ datasource.DataSourceWithConfigValidators = &SystemRoleDataSource{}

// NewSystemRoleDataSource is a helper function
func NewSystemRoleDataSource() datasource.DataSource {
	return &SystemRoleDataSource{}
}

// defines the data source implementation.
type SystemRoleDataSource struct {
	client *tilgangsportalapi.Client
}

// Metadata returns the resource type name.
func (d *SystemRoleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_system_role"
}

// Schema defines the schema for the resource.
func (d *SystemRoleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Schema for the SystemRole data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the system role.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the system role. Either name or object_id must be specified.",
				Optional:    true,
				Computed:    true,
			},
			"object_id": schema.StringAttribute{
				Description: "The object ID of the system role. Either name or object_id must be specified.",
				Optional:    true,
				Computed:    true,
			},
			"system_role_owner": schema.StringAttribute{
				Description: "The system role owner.",
				Computed:    true,
			},
			"system_role_security_owner": schema.StringAttribute{
				Description: "The system role security owner.",
				Computed:    true,
			},
			"approval_level": schema.StringAttribute{
				Description: "The approval level of the system role.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the system role.",
				Computed:    true,
			},
			"product_category": schema.StringAttribute{
				Description: "The product category of the system role.",
				Computed:    true,
			},
			"it_shop_name": schema.StringAttribute{
				Description: "The IT shop name of the system role.",
				Computed:    true,
			},
		},
	}
}

// ConfigValidators returns validators for the data source configuration
func (d *SystemRoleDataSource) ConfigValidators(ctx context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(
			path.MatchRoot("name"),
			path.MatchRoot("object_id"),
		),
	}
}

// Configure adds the provider configured client to the resource.
func (d *SystemRoleDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	ConfigureClientDataSource(ctx, req, resp, func(client *tilgangsportalapi.Client) {
		d.client = client
	})
}

// Read the resource data.
func (d *SystemRoleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SystemRoleModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Find the role by name or object_id
	role, err := d.findSystemRole(&data)

	if err != nil {
		resp.Diagnostics.AddError("Failed to get system role", err.Error())
		return
	}
	if role == nil {
		resp.Diagnostics.AddError("System role not found", "No matching system role found")
		return
	}

	// Set the resource data
	data.ID = types.StringValue(role.Name)
	data.Name = types.StringValue(role.Name)
	data.ObjectID = types.StringValue(role.ObjectID)
	data.SystemRoleOwner = types.StringValue(role.L2Ident)
	data.SystemRoleSecurityOwner = types.StringValue(role.L3Ident)
	data.ApprovalLevel = types.StringValue(role.ApprovalLevel)
	data.Description = types.StringValue(role.Description)
	data.ProductCategory = types.StringValue(role.ProductCategory)
	data.ItShopName = types.StringValue(role.ItShopName)

	// Write the resource data
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// findSystemRole finds a system role by name or object_id
func (d *SystemRoleDataSource) findSystemRole(data *SystemRoleModel) (*tilgangsportalapi.SystemRole, error) {
	// Lookup by name
	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		return d.client.GetSystemRole(data.Name.ValueString())
	}

	// Lookup by object_id
	if !data.ObjectID.IsNull() && !data.ObjectID.IsUnknown() {
		return d.client.GetSystemRoleID(data.ObjectID.ValueString())
	}

	return nil, nil
}
