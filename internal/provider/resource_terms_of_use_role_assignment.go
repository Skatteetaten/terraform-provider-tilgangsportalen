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
var _ resource.Resource = &NewTermsOfUseRoleAssignmentResource{}
var _ resource.ResourceWithImportState = &NewTermsOfUseRoleAssignmentResource{}

// CreateNewTermsOfUseRoleAssignment is a helper function
func CreateNewTermsOfUseRoleAssignment() resource.Resource {
	return &NewTermsOfUseRoleAssignmentResource{}
}

// NewTermsOfUseRoleAssignmentResource defines the resource implementation.
type NewTermsOfUseRoleAssignmentResource struct {
	client *tilgangsportalapi.Client
}

// TermsOfUseRoleAssignmentModel defines the resource data model.
type TermsOfUseRoleAssignmentModel struct {
	ID                    types.String `tfsdk:"id"`
	RoleName              types.String `tfsdk:"role_name"`
	TermsOfUseIdentifier  types.String `tfsdk:"terms_of_use_identifier"`
	TermsOfUseDescription types.String `tfsdk:"terms_of_use_description"`
}

// Metadata returns the resource type name.
func (r *NewTermsOfUseRoleAssignmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_terms_of_use_role_assignment"
}

// IDPlanModifier is a custom plan modifier that updates the ID when terms_of_use_identifier changes
type IDPlanModifier struct{}

func (m IDPlanModifier) Description(ctx context.Context) string {
	return "Updates the ID when terms_of_use_identifier changes"
}

func (m IDPlanModifier) MarkdownDescription(ctx context.Context) string {
	return "Updates the ID when terms_of_use_identifier changes"
}

func (m IDPlanModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// If this is a create operation (no prior state), let the create method handle ID generation
	if req.State.Raw.IsNull() {
		tflog.Debug(ctx, "IDPlanModifier: Create operation detected, letting create method handle ID generation")
		return
	}

	// Get the current state and planned values
	var state, plan TermsOfUseRoleAssignmentModel
	req.State.Get(ctx, &state)
	req.Plan.Get(ctx, &plan)

	// If terms_of_use_identifier is changing, compute new ID
	if !plan.TermsOfUseIdentifier.Equal(state.TermsOfUseIdentifier) {
		newID := fmt.Sprintf("%s|%s", plan.RoleName.ValueString(), plan.TermsOfUseIdentifier.ValueString())
		tflog.Debug(ctx, fmt.Sprintf("IDPlanModifier: Terms of use identifier changed from %s to %s, updating ID to %s",
			state.TermsOfUseIdentifier.ValueString(), plan.TermsOfUseIdentifier.ValueString(), newID))
		resp.PlanValue = types.StringValue(newID)
		return
	}

	// Otherwise, keep the current ID
	tflog.Debug(ctx, fmt.Sprintf("IDPlanModifier: No changes detected, keeping current ID: %s", req.PlanValue.ValueString()))
	resp.PlanValue = req.PlanValue
}

// NewIDPlanModifier creates a new IDPlanModifier
func NewIDPlanModifier() planmodifier.String {
	return IDPlanModifier{}
}

// Schema defines the schema for the resource.
func (r *NewTermsOfUseRoleAssignmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Manages Terms of Use assignments for System Roles in Tilgangsportalen.\n\n" +
			"This resource creates and manages the association between a System Role and Terms of Use object. " +
			"When users are granted access to a role with Terms of Use, they must accept the terms before gaining access.\n\n" +
			"**Important Notes:**\n\n" +
			"- Only one Terms of Use can be assigned to a role at a time.\n\n" +
			"- Assigning new Terms of Use will replace any existing assignment.\n\n" +
			"- If the referenced System Role is deleted, this assignment will be automatically removed from state.\n\n" +
			"- Changes to `role_name` will force recreation of the assignment.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				MarkdownDescription: "The unique identifier for this Terms of Use assignment " +
					"formatted as `RoleName|TermsOfUseIdentifier`",
				// Plan modifier to import id from previous state to avoid "know after apply" message
				PlanModifiers: []planmodifier.String{
					NewIDPlanModifier(),
				},
			},
			"role_name": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "The name of the System Role to assign Terms of Use to. " +
					"This must be an existing System Role in Tilgangsportalen. " +
					"Changing this value will force creation of a new resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"terms_of_use_identifier": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "The unique identifier of the Terms of Use object to assign to the role. " +
					"This identifier references a Terms of Use object that exists in Tilgangsportalen. " +
					"Changing this value will update the assignment to use the new Terms of Use " +
					"object, replacing any previously assigned terms. Only one Terms of Use can be assigned " +
					"to a role at a time.",
			},
			"terms_of_use_description": schema.StringAttribute{
				Computed: true,
				MarkdownDescription: "The human-readable description of the assigned Terms of Use object. " +
					"This field is automatically populated from Tilgangsportalen.",
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *NewTermsOfUseRoleAssignmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	ConfigureClientResource(ctx, req, resp, func(client *tilgangsportalapi.Client) {
		r.client = client
	})
}

// Create a new terms of use role assignment resource
func (r *NewTermsOfUseRoleAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TermsOfUseRoleAssignmentModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	assignment := tilgangsportalapi.TermsOfUseAssignment{
		RoleName:   data.RoleName.ValueString(),
		TermsOfUse: data.TermsOfUseIdentifier.ValueString(),
	}
	_, err := r.client.AssignTermsOfUseToRole(assignment)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to assign Terms of Use %s to Role %s, got error: %s", data.TermsOfUseIdentifier.ValueString(), data.RoleName.ValueString(), err))
		return
	}

	// Wait for assignment to be completed
	termsOfUseDetails, err := r.client.WaitForTermsOfUseAssignment(assignment.TermsOfUse, assignment.RoleName)
	if err != nil {
		resp.Diagnostics.AddError("Assignment Verification Failed", fmt.Sprintf("Terms of Use assignment failed: %s", err))
		return
	}

	// Setting assignment ID by combining RoleName and TermsOfUse with a pipe
	data.ID = types.StringValue(fmt.Sprintf("%s|%s", data.RoleName.ValueString(), data.TermsOfUseIdentifier.ValueString()))
	tflog.Trace(ctx, fmt.Sprintf("resource ID %s|%s added to resource", data.RoleName.ValueString(), data.TermsOfUseIdentifier.ValueString()))

	// Set description
	if termsOfUseDetails != nil {
		data.TermsOfUseDescription = types.StringValue(termsOfUseDetails.Description)
	} else {
		data.TermsOfUseDescription = types.StringNull()
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully created Terms of Use %s assignment to system Role %s. Assignment ID set to %s", data.TermsOfUseIdentifier.ValueString(), data.RoleName.ValueString(), data.ID.ValueString()))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read calls the API to get the latest data for the resource
func (r *NewTermsOfUseRoleAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TermsOfUseRoleAssignmentModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "Read Terms of Use Role Assignment to see if the terms of use is still assigned to the role")

	// Get the terms of use assigned to the role
	termsOfUse, err := r.client.GetTermsOfUseForRole(data.RoleName.ValueString())
	if err != nil {
		// Check if error is because no terms of use is assigned (error code 606)
		if strings.Contains(strings.ToLower(err.Error()), "no terms of use is assigned to") {
			// No terms of use assigned, remove from state
			tflog.Info(ctx, fmt.Sprintf("Terms of Use assignment for role %s not found. Removing from state.", data.RoleName.ValueString()))
			resp.State.RemoveResource(ctx)
			return
		}
		// Other errors (like role not found)
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get Terms of Use for System Role %s, got error: %s", data.RoleName.ValueString(), err))
		return
	}

	// Check if the assigned terms of use matches our state
	if termsOfUse == nil {
		// No terms of use assigned, remove from state
		tflog.Info(ctx, fmt.Sprintf("No Terms of Use is assigned to role %s. Removing from state.", data.RoleName.ValueString()))
		resp.State.RemoveResource(ctx)
		return
	}

	if termsOfUse.TermsOfUseIdentifier != data.TermsOfUseIdentifier.ValueString() {
		// Different terms of use assigned, update the resource in state
		tflog.Info(ctx, fmt.Sprintf("Terms of Use for role %s has changed from %s to %s. Updating state.", data.RoleName.ValueString(), data.TermsOfUseIdentifier.ValueString(), termsOfUse.TermsOfUseIdentifier))
		data.TermsOfUseIdentifier = types.StringValue(termsOfUse.TermsOfUseIdentifier)
	}

	tflog.Debug(ctx, fmt.Sprintf("Terms of Use %s is assigned to role %s", data.TermsOfUseIdentifier.ValueString(), data.RoleName.ValueString()))

	// Update the description from the current assignment
	data.TermsOfUseDescription = types.StringValue(termsOfUse.Description)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update handles updates to the terms of use role assignment resource
func (r *NewTermsOfUseRoleAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TermsOfUseRoleAssignmentModel
	var state TermsOfUseRoleAssignmentModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read current state data
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check if terms_of_use has changed
	if data.TermsOfUseIdentifier.ValueString() != state.TermsOfUseIdentifier.ValueString() {
		// Assign the new terms of use (this will overwrite any existing assignment)
		assignment := tilgangsportalapi.TermsOfUseAssignment{
			RoleName:   data.RoleName.ValueString(),
			TermsOfUse: data.TermsOfUseIdentifier.ValueString(),
		}
		_, err := r.client.AssignTermsOfUseToRole(assignment)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to assign new Terms of Use %s to Role %s during update, got error: %s", data.TermsOfUseIdentifier.ValueString(), data.RoleName.ValueString(), err))
			return
		}

		// Wait for assignment to be completed
		termsOfUseDetails, err := r.client.WaitForTermsOfUseAssignment(assignment.TermsOfUse, assignment.RoleName)
		if err != nil {
			resp.Diagnostics.AddError("Assignment Verification Failed", fmt.Sprintf("Terms of Use assignment failed: %s", err))
			return
		}

		// Set description
		if termsOfUseDetails != nil {
			data.TermsOfUseDescription = types.StringValue(termsOfUseDetails.Description)
		} else {
			data.TermsOfUseDescription = types.StringNull()
		}

		tflog.Debug(ctx, fmt.Sprintf("Successfully updated Terms of Use assignment for role %s from %s to %s", data.RoleName.ValueString(), state.TermsOfUseIdentifier.ValueString(), data.TermsOfUseIdentifier.ValueString()))
	} else {
		// Get the current terms of use details and update description
		termsOfUseDetails, err := r.client.GetTermsOfUseForRole(data.RoleName.ValueString())
		if err == nil && termsOfUseDetails != nil {
			data.TermsOfUseDescription = types.StringValue(termsOfUseDetails.Description)
		} else {
			data.TermsOfUseDescription = types.StringNull()
		}
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete removes the terms of use role assignment
func (r *NewTermsOfUseRoleAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TermsOfUseRoleAssignmentModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create removal request
	removal := tilgangsportalapi.TermsOfUseRemoval{
		RoleName: data.RoleName.ValueString(),
	}

	// Remove the terms of use from the role
	_, err := r.client.RemoveTermsOfUseFromRole(removal)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to remove Terms of Use from Role %s, got error: %s", data.RoleName.ValueString(), err))
		return
	}

	// Check if removal was successful
	err = r.client.WaitForTermsOfUseRemoval(data.RoleName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Removal Verification Failed", fmt.Sprintf("Terms of Use removal failed: %s", err))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully removed Terms of Use assignment from role %s", data.RoleName.ValueString()))

	// Removal from state is handled automatically by the plugin framework
}

// ImportState imports a terms of use role assignment to state
func (r *NewTermsOfUseRoleAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Expected format: "RoleName|TermsOfUseIdentifier"
	parts := strings.Split(req.ID, "|")
	if len(parts) != 2 {
		resp.Diagnostics.AddError("Invalid Import ID", "Expected format: RoleName|TermsOfUseIdentifier")
		return
	}

	roleName := parts[0]
	termsOfUseIdentifier := parts[1]

	// Create the data model
	data := TermsOfUseRoleAssignmentModel{
		ID:                   types.StringValue(req.ID),
		RoleName:             types.StringValue(roleName),
		TermsOfUseIdentifier: types.StringValue(termsOfUseIdentifier),
	}

	// Verify the assignment exists
	termsOfUse, err := r.client.GetTermsOfUseForRole(roleName)
	if err != nil {
		if strings.Contains(err.Error(), "no terms of use is assigned to role") {
			resp.Diagnostics.AddError("Assignment not found", fmt.Sprintf("No Terms of Use is assigned to role %s", roleName))
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to check if Terms of Use is assigned to Role %s, got error: %s", roleName, err))
		return
	}

	if termsOfUse == nil || termsOfUse.TermsOfUseIdentifier != termsOfUseIdentifier {
		resp.Diagnostics.AddError("Assignment not found", fmt.Sprintf("Terms of Use %s is not assigned to role %s (currently assigned: %s)", termsOfUseIdentifier, roleName, func() string {
			if termsOfUse != nil {
				return termsOfUse.TermsOfUseIdentifier
			}
			return "none"
		}()))
		return
	}

	// Set the description from the assignment
	if termsOfUse != nil {
		data.TermsOfUseDescription = types.StringValue(termsOfUse.Description)
	} else {
		data.TermsOfUseDescription = types.StringNull()
	}

	// Save imported data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
