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
var _ datasource.DataSource = &EntraGroupsDataSource{}

// NewEntraGroupsDataSource is a helper function
func NewEntraGroupsDataSource() datasource.DataSource {
	return &EntraGroupsDataSource{}
}

// EntraGroupsDataSource defines the data source implementation.
type EntraGroupsDataSource struct {
	client *tilgangsportalapi.Client
}

// EntraGroupsDataSourceModel describes the data source data model.
type EntraGroupsDataSourceModel struct {
	Groups []SingleEntraGroupModel `tfsdk:"groups"`
}

// SingleEntraGroupModel is used when reading the API response when listing
// Entra groups
type SingleEntraGroupModel struct {
	GroupName types.String `tfsdk:"displayname"`
}

// Metadata returns the resource type name.
func (d *EntraGroupsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_entra_groups"
}

// Schema for the Entra group datasource 
func (d *EntraGroupsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Schema for the EntraGroups data source",

		Attributes: map[string]schema.Attribute{
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

// Configure adds the provider configured client to the resource.
func (d *EntraGroupsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Read calls the API to get the latest data for the resource
func (d *EntraGroupsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data EntraGroupsDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := d.client.ListEntraGroups()
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

	tflog.Debug(ctx, "Read an Entra ID group data source.")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}
