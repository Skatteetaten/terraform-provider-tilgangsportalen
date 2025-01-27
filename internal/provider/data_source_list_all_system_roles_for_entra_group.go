package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-tilgangsportalen/internal/tilgangsportalapi"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ListAllSystemRolesForEntraGroup{}

// NewListAllSystemRolesForEntraGroup is a helper function
func NewListAllSystemRolesForEntraGroup() datasource.DataSource {
	return &ListAllSystemRolesForEntraGroup{}
}

// ListAllSystemRolesForEntraGroup defines the data source implementation.
type ListAllSystemRolesForEntraGroup struct {
	client *tilgangsportalapi.Client
}

// ListAllSystemRolesForEntraGroupModel describes the data source data model.
type ListAllSystemRolesForEntraGroupModel struct {
	GroupName types.String                     `tfsdk:"group_name"`
	Roles     []SingleSystemRoleWithOwnerModel `tfsdk:"roles"`
}

type SingleSystemRoleWithOwnerModel struct {
	DisplayName   types.String `tfsdk:"display_name"`
	L2Ident       types.String `tfsdk:"system_role_owner"`
	L2DisplayName types.String `tfsdk:"system_role_owner_display_name"`
}

// Metadata returns the resource type name.
func (d *ListAllSystemRolesForEntraGroup) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_system_roles_assigned_to_entra_group"
}

// Schema defines the schema for the resource.
func (d *ListAllSystemRolesForEntraGroup) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Lists all system roles assigned to an Entra Group. Also lists any roles assigned that the current user do not manage.",

		Attributes: map[string]schema.Attribute{
			"group_name": schema.StringAttribute{
				Description: "The name of the Entra group to list all System Roles for.",
				Required:    true,
			},
			"roles": schema.ListNestedAttribute{
				Description: "List of System Roles.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"display_name": schema.StringAttribute{
							Description: "String identifier of the System Role display name.",
							Computed:    true,
						},
						"system_role_owner": schema.StringAttribute{
							Description: "String identifier of the System Role owner.",
							Computed:    true,
						},
						"system_role_owner_display_name": schema.StringAttribute{
							Description: "String identifier of the System Role owner display name.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (d *ListAllSystemRolesForEntraGroup) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	ConfigureClientDataSource(ctx, req, resp, func(client *tilgangsportalapi.Client) {
		d.client = client
	})
}

// Read calls the API to get the latest data for the resource
func (d *ListAllSystemRolesForEntraGroup) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ListAllSystemRolesForEntraGroupModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Getting owned system roles
	response, err := d.client.ListAllSystemRolesForEntraGroup(data.GroupName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read all system roles for group, got error: %s", err))
		return
	}

	for _, role := range response.Roles {
		roleState := SingleSystemRoleWithOwnerModel{
			DisplayName:   types.StringValue(role.DisplayName),
			L2Ident:       types.StringValue(role.L2Ident),
			L2DisplayName: types.StringValue(role.L2DisplayName),
		}
		data.Roles = append(data.Roles, roleState)
	}

	tflog.Trace(ctx, "Found list of owned system roles.")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
