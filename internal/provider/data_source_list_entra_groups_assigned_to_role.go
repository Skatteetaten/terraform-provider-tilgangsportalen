// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"terraform-provider-tilgangsportalen/internal/tilgangsportalapi"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &EntraGroupsForRoleDataSource{}

func NewEntraGroupsForRoleDataSource() datasource.DataSource {
	return &EntraGroupsForRoleDataSource{}
}

// EntraGroupsForRoleDataSource defines the data source implementation.
type EntraGroupsForRoleDataSource struct {
	client *tilgangsportalapi.Client
}

// EntraGroupsForRoleDataSource describes the data source data model.
type EntraGroupsForRoleDataSourceModel struct {
	RoleName types.String            `tfsdk:"role_name"`
	Groups   []SingleEntraGroupModel `tfsdk:"groups"`
}

func (d *EntraGroupsForRoleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_entra_groups_assigned_to_role"
}

func (d *EntraGroupsForRoleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Schema for the EntraGroupsForRole data source",

		Attributes: map[string]schema.Attribute{
			"role_name": schema.StringAttribute{
				Description: "The name of the role to list Entra groups for.",
				Required:    true,
			},
			"groups": schema.ListNestedAttribute{
				Description: "List of Entra groups.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"displayname": schema.StringAttribute{
							Description: "String identifier of the Entra group display name.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *EntraGroupsForRoleDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*tilgangsportalapi.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *tilgangsportalapi.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *EntraGroupsForRoleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data EntraGroupsForRoleDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := d.client.ListEntraGroupsForRole(data.RoleName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read example, got error: %s", err))
		return
	}

	for _, group := range response.EntraGroups {
		groupState := SingleEntraGroupModel{
			GroupName: types.StringValue(group.DisplayName),
		}
		data.Groups = append(data.Groups, groupState)
	}

	tflog.Trace(ctx, "Read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
