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
var _ datasource.DataSource = &SystemRolesDataSource{}

func NewSystemRolesDataSource() datasource.DataSource {
	return &SystemRolesDataSource{}
}

// SystemRolesDataSource defines the data source implementation.
type SystemRolesDataSource struct {
	client *tilgangsportalapi.Client
}

// SystemRolesDataSource describes the data source data model.
type SystemRolesDataSourceModel struct {
	Roles []SingleSystemRoleModel `tfsdk:"roles"`
}

// Can be later replaced by the model found in the System Role resource file
type SingleSystemRoleModel struct {
	RoleName types.String `tfsdk:"displayname"`
}

func (d *SystemRolesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_system_roles"
}

func (d *SystemRolesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Schema for the SystemRoles data source",

		Attributes: map[string]schema.Attribute{
			"roles": schema.ListNestedAttribute{
				Description: "List of System Roles.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"displayname": schema.StringAttribute{
							Description: "String identifier of the System Role display name.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *SystemRolesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SystemRolesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SystemRolesDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	
	// Getting owned system roles
	response, err := d.client.ListSystemRoles()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read system roles, got error: %s", err))
		return
	}

	for _, role := range response.Roles {
		roleState := SingleSystemRoleModel{
			RoleName: types.StringValue(role.DisplayName),
		}
		data.Roles = append(data.Roles, roleState)
	}
	
	tflog.Trace(ctx, "Found list of owned system roles.")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
