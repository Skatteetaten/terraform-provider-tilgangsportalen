package provider

import (
	"context"
	"fmt"
	"strings"
	"terraform-provider-tilgangsportalen/internal/tilgangsportalapi"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &TermsOfUseRoleAssignmentDataSource{}

// NewTermsOfUseRoleAssignmentDataSource is a helper function
func NewTermsOfUseRoleAssignmentDataSource() datasource.DataSource {
	return &TermsOfUseRoleAssignmentDataSource{}
}

// TermsOfUseRoleAssignmentDataSource defines the data source implementation.
type TermsOfUseRoleAssignmentDataSource struct {
	client *tilgangsportalapi.Client
}

// TermsOfUseRoleAssignmentDataSourceModel describes the data source data model.
type TermsOfUseRoleAssignmentDataSourceModel struct {
	ID                    types.String `tfsdk:"id"`
	RoleName              types.String `tfsdk:"role_name"`
	TermsOfUseIdentifier  types.String `tfsdk:"terms_of_use_identifier"`
	TermsOfUseDescription types.String `tfsdk:"terms_of_use_description"`
	TermsOfUseUID         types.String `tfsdk:"terms_of_use_uid"`
}

// Metadata returns the data source type name.
func (d *TermsOfUseRoleAssignmentDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_terms_of_use_role_assignment"
}

// Schema defines the schema for the data source.
func (d *TermsOfUseRoleAssignmentDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Retrieves information about Terms of Use assignments for System Roles in Tilgangsportalen.\n\n" +
			"This data source allows you to look up which Terms of Use document is currently assigned to a specific " +
			"System Role, along with detailed information about the Terms of Use document itself.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				MarkdownDescription: "The unique identifier for this Terms of Use assignment " +
					"formatted as `RoleName|TermsOfUseIdentifier`.",
			},
			"role_name": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "The name of the System Role to look up Terms of Use assignments for. " +
					"This must be an existing System Role in Tilgangsportalen. If no Terms of Use is assigned " +
					"to this role, the data source will return an error.",
			},
			"terms_of_use_identifier": schema.StringAttribute{
				Computed: true,
				MarkdownDescription: "The unique identifier of the Terms of Use document assigned to the role. " +
					"This identifier uniquely identifies the Terms of Use document within Tilgangsportalen " +
					"and can be used to reference the same document in other configurations.",
			},
			"terms_of_use_description": schema.StringAttribute{
				Computed: true,
				MarkdownDescription: "The human-readable description of the assigned Terms of Use document. " +
					"This provides detailed information about what the Terms of Use document contains, " +
					"helping administrators understand the legal agreement users must accept.",
			},
			"terms_of_use_uid": schema.StringAttribute{
				Computed: true,
				MarkdownDescription: "The unique system identifier (UID) of the Terms of Use document. " +
					"This is an internal system identifier used by Tilgangsportalen for tracking and " +
					"referencing the Terms of Use document in various system operations.",
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *TermsOfUseRoleAssignmentDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	ConfigureClientDataSource(ctx, req, resp, func(client *tilgangsportalapi.Client) {
		d.client = client
	})
}

// Read retrieves the Terms of Use assignment data from the API.
func (d *TermsOfUseRoleAssignmentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data TermsOfUseRoleAssignmentDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	roleName := data.RoleName.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Looking up Terms of Use assignment for role: %s", roleName))

	// Get the Terms of Use assignment from the API
	termsOfUse, err := d.client.GetTermsOfUseForRole(roleName)
	if err != nil {
		if strings.Contains(err.Error(), "no terms of use is assigned to role") {
			resp.Diagnostics.AddError(
				"No Terms of Use Assignment Found",
				fmt.Sprintf("No Terms of Use is currently assigned to role '%s'. Ensure the role exists and has a Terms of Use assignment before using this data source.", roleName),
			)
			return
		}
		resp.Diagnostics.AddError(
			"API Error",
			fmt.Sprintf("Unable to retrieve Terms of Use assignment for role '%s': %s", roleName, err),
		)
		return
	}

	if termsOfUse == nil {
		resp.Diagnostics.AddError(
			"No Terms of Use Assignment Found",
			fmt.Sprintf("No Terms of Use is currently assigned to role '%s'.", roleName),
		)
		return
	}

	// Populate the data model
	data.ID = types.StringValue(fmt.Sprintf("%s|%s", roleName, termsOfUse.TermsOfUseIdentifier))
	data.RoleName = types.StringValue(roleName)
	data.TermsOfUseIdentifier = types.StringValue(termsOfUse.TermsOfUseIdentifier)
	data.TermsOfUseDescription = types.StringValue(termsOfUse.Description)
	data.TermsOfUseUID = types.StringValue(termsOfUse.TermsOfUseUID)

	tflog.Debug(ctx, fmt.Sprintf("Successfully retrieved Terms of Use assignment: %s -> %s", roleName, termsOfUse.TermsOfUseIdentifier))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
