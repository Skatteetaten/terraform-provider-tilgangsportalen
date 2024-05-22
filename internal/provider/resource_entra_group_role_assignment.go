package provider

import (
	"context"
	"fmt"
	"strings"
	"terraform-provider-tilgangsportalen/internal/tilgangsportalapi"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &NewEntraGroupRoleAssignmentResource{}
var _ resource.ResourceWithImportState = &NewEntraGroupRoleAssignmentResource{}

// CreateNewEntraGroupRoleAssignment is a helper function
func CreateNewEntraGroupRoleAssignment() resource.Resource {
	return &NewEntraGroupRoleAssignmentResource{}
}

// NewEntraGroupRoleAssignmentResource defines the resource implementation.
type NewEntraGroupRoleAssignmentResource struct {
	client *tilgangsportalapi.Client
}

// EntraGroupRoleAssignmentModel defines the resource data model.
type EntraGroupRoleAssignmentModel struct {
	ID         types.String `tfsdk:"id"`
	RoleName   types.String `tfsdk:"role_name"`
	EntraGroup types.String `tfsdk:"entra_group"`
}

// Metadata returns the resource type name.
func (r *NewEntraGroupRoleAssignmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_entra_group_role_assignment"
}

// Schema defines the schema for the resource.
func (r *NewEntraGroupRoleAssignmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This resource is used to create assignments between Entra Groups and System Roles in Tilgangsportalen",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Identifier for the Entra Group System Role assignment. Currently, as we do not get a unique ID we can use from the API, ID is set by combining the role name and the Entra group name, with a pipe symbol as separator: RoleName|EntraGroupName",
				// Plan modifier to import id from previous state to avoid "know after apply" message
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"role_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the Role to assign the Entra Group to",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"entra_group": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the Entra Group to assign to the Role",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *NewEntraGroupRoleAssignmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*tilgangsportalapi.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *tilgangsportalapi.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

// Create a new role assignment resource
func (r *NewEntraGroupRoleAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data EntraGroupRoleAssignmentModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	roleAssignment := tilgangsportalapi.EntraGroupRoleAssignment{
		RoleName:   data.RoleName.ValueString(),
		EntraGroup: data.EntraGroup.ValueString(),
	}
	
	_, err := r.client.AssignEntraGroupToRole(roleAssignment)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Entra Group %s assignment to system Role %s, got error: %s", data.EntraGroup, data.RoleName, err))
		return
	}

	// Setting role ID to be equal the new role name combining RoleName and GroupName with a pipe
	data.ID = types.StringValue(fmt.Sprintf("%s|%s", data.RoleName.ValueString(), data.EntraGroup.ValueString()))
	tflog.Trace(ctx,fmt.Sprintf("resource ID %s|%s added to resource", data.RoleName, data.EntraGroup))

	tflog.Debug(ctx, fmt.Sprintf("Successfully created Entra Group %s assignment to system Role %s. Role assignment ID set to %s", data.EntraGroup, data.RoleName, data.ID))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read calls the API to get the latest data for the resource
func (r *NewEntraGroupRoleAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data EntraGroupRoleAssignmentModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx,"Read Entra Group Role Assignment to see if the role is still assigned to the group")

	response, err := r.client.ListEntraGroupsForRole(data.RoleName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list Entra Goups for System Role %s, got error: %s", data.RoleName, err))
		return
	}

	// Check if the group is assigned to the role
	assigned := 0

	for _, group := range response.EntraGroups {

		tflog.Debug(ctx,fmt.Sprintf("Checking if group %s matches the group from resource state %s", group.DisplayName, data.EntraGroup))

		if group.DisplayName == data.EntraGroup.ValueString() {
			assigned = 1
		}
	}

	// if group is not assigned, then remove it from the state
	if assigned == 0 {
		tflog.Debug(ctx,fmt.Sprintf("Group %s is not assigned to role %s, removing from state", data.EntraGroup.ValueString(), data.RoleName.ValueString()))
		resp.State.RemoveResource(ctx)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update - We have no update API method for role assignments. All changes will
// require replacement. 
func (r *NewEntraGroupRoleAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data EntraGroupRoleAssignmentModel
	
	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete a role assignment resource
func (r *NewEntraGroupRoleAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data EntraGroupRoleAssignmentModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	role := tilgangsportalapi.EntraGroupRoleAssignment{
		RoleName:   data.RoleName.ValueString(),
		EntraGroup: data.EntraGroup.ValueString(),
	}

	_, err := r.client.RemoveEntraGroupFromRole(role)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to Remove Entra Group %s from Role %s, got error: %s", data.EntraGroup, data.RoleName, err))
		return
	}

	// Removal from state is handled by the framework
}

// ImportState imports a role assignment to state
func (r *NewEntraGroupRoleAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, "|")

	tflog.Debug(ctx,fmt.Sprintf("Importing Entra Group Role Assignment with ID %s to state", req.ID))

	if len(idParts) != 2 {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Invalid ID %s, expected RoleName|EntraGroupName", req.ID))
		return
	}

	data := EntraGroupRoleAssignmentModel{
		RoleName:   types.StringValue(idParts[0]),
		EntraGroup: types.StringValue(idParts[1]),
		ID:         types.StringValue(req.ID),
	}

	// Check if the group is assigned to the role
	assigned, err := r.client.CheckIfGroupIsAssignedToRole(data.EntraGroup.ValueString(), data.RoleName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to check if Entra Group %s is assigned to Role %s, got error: %s", data.EntraGroup.ValueString(), data.RoleName.ValueString(), err))
		return
	}

	if !assigned {
		resp.Diagnostics.AddError("Assignment not found", fmt.Sprintf("Group Assignment %s not found on role %s", data.EntraGroup.ValueString(), data.RoleName.ValueString()))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}
