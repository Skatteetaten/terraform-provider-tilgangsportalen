package provider

import (
	"context"
	"fmt"
	"regexp"
	"terraform-provider-tilgangsportalen/internal/tilgangsportalapi"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &NewSystemRoleResource{}
var _ resource.ResourceWithImportState = &NewSystemRoleResource{}

// CreateNewSystemRole is a helper function
func CreateNewSystemRole() resource.Resource {
	return &NewSystemRoleResource{}
}

// NewSystemRoleResource defines the resource implementation.
type NewSystemRoleResource struct {
	client *tilgangsportalapi.Client
}

// SystemRoleModel describes the resource data model.
type SystemRoleModel struct {
	ID                      types.String `tfsdk:"id"`
	Name                    types.String `tfsdk:"name"`
	SystemRoleOwner         types.String `tfsdk:"system_role_owner"`
	SystemRoleSecurityOwner types.String `tfsdk:"system_role_security_owner"`
	ApprovalLevel           types.String `tfsdk:"approval_level"`
	Description             types.String `tfsdk:"description"`
	ProductCategory         types.String `tfsdk:"product_category"`
	ItShopName              types.String `tfsdk:"it_shop_name"`
}

// Metadata returns the resource type name.
func (r *NewSystemRoleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_system_role"
}

// Schema defines the schema for the resource.
func (r *NewSystemRoleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This resource is used to create a new System Role in Tilgangsportalen",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Identifier for the System Role. Currently, as we do not get a unique ID we can use from the API, ID is set equal to Name",
				// Plan modifier to import id from previous state to avoid "know after apply" message
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the system role to be created. Please follow the standardized naming conventions for roles.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[æøåÆØÅa-zA-Z0-9 _-]+$`),
						"The name of the role may only contain alphanumeric characters, space ( ), underscore (_), and dash (-). The maximum length is 256 characters.",
					),
					stringvalidator.LengthAtMost(256),
				},
			},
			"system_role_owner": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The owner of the role, identified by their user ident",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[[a-zA-Z][0-9]{2}[a-zA-Z]{3}|[a-zA-Z][0-9]{5}$`),
						"Must be a valid user ident in the form of x00000 or x00xxx",
					),
				},
			},
			"system_role_security_owner": schema.StringAttribute{
				Optional:            true, // Note: is required if approval level is L3
				MarkdownDescription: "The security owner of the role. Required if the approval level is L3",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[[a-zA-Z][0-9]{2}[a-zA-Z]{3}|[a-zA-Z][0-9]{5}$`),
						"Must be a valid user ident in the form of x00000 or x00xxx",
					),
				},
			},
			"approval_level": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The approval level for the role",
				Validators: []validator.String{
					// Accepted approval levels are L0, L1, L2 and L3
					stringvalidator.OneOf([]string{"L0", "L1", "L2", "L3"}...),
				},
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "A description of the role",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[æøåÆØÅa-zA-Z0-9 .,?!_()/\-\[\]]*$`),
						"The description may only contain alphanumeric characters, punctuation (.,?!), space ( ), brackets (()), square brackets ([]),  forward slash (/), underscore (_), and dash (-).",
					),
				},
			},
			"product_category": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The product category (tjenestekategori) for this role. Should match an existing product catogory in tilgangsportalen",
			},
			"it_shop_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the IT shop",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *NewSystemRoleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create a new system role resource
func (r *NewSystemRoleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SystemRoleModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	role := tilgangsportalapi.SystemRole{
		Name:            data.Name.ValueString(),
		L2Ident:         data.SystemRoleOwner.ValueString(),
		L3Ident:         data.SystemRoleSecurityOwner.ValueString(),
		ApprovalLevel:   data.ApprovalLevel.ValueString(),
		Description:     data.Description.ValueString(),
		ProductCategory: data.ProductCategory.ValueString(),
		ItShopName:      data.ItShopName.ValueString(),
	}

	_, err := r.client.CreateAndPublishSystemRole(role)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create System Role %s, got error: %s", data.Name, err))
		return
	}

	// Setting role ID to be equal the role name
	data.ID = data.Name

	tflog.Debug(ctx, fmt.Sprintf("System Role %s created", data.Name))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read calls the API to get the latest data for the resource
func (r *NewSystemRoleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SystemRoleModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Check if the role exists
	exists, err := r.client.CheckIfRoleExists(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to check if System Role %s exists, got error: %s", data.Name.ValueString(), err))
		return
	}
	// remove from state if not exists
	if !exists {
		tflog.Info(ctx, fmt.Sprintf("System Role %s not found. Removing from state.", data.Name))
		resp.State.RemoveResource(ctx)
		return
	}

	// If the role exists, we get role and update state
	systemRole, err := r.client.GetSystemRole(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to import System Role %s, got error: %s", data.Name, err))
		return
	}

	// Map to SystemRoleModel and save updated data into Terraform state
	data.Name = types.StringValue(systemRole.Name)
	data.ApprovalLevel = types.StringValue(systemRole.ApprovalLevel)
	data.ProductCategory = types.StringValue(systemRole.ProductCategory)
	// If no description is set, GetSystemRole returns an empty string
	// We only want the plan to show change if description has actually changed
	if systemRole.Description == "" && data.Description != types.StringValue("") {
		data.Description = types.StringNull()
	} else {
		data.Description = types.StringValue(systemRole.Description)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// Update a system role resource
func (r *NewSystemRoleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var rolePlan SystemRoleModel
	var roleState SystemRoleModel
	var namePlan types.String
	var nameState types.String

	// Read Terraform plan into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &rolePlan)...)
	// Read Terraform state into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &roleState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Getting name to compare
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("name"), &namePlan)...)
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("name"), &nameState)...)

	// If the name of the role differs, we call RenameSystemRole
	if !namePlan.Equal(nameState) {

		renameRole := tilgangsportalapi.RenameSystemRole{
			OldName: nameState.ValueString(),
			NewName: namePlan.ValueString(),
		}

		_, err := r.client.RenameSystemRole(renameRole)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to rename System Role %s to %s, got error: %s", nameState, namePlan, err))
			return
		}
	}

	// If one or more of the fields Description, Approval Level, System Role
	// Owner, System Role Security Owner or Product Category differs, we call
	// UpdateRole.
	if !rolePlan.Description.Equal(roleState.Description) || !rolePlan.SystemRoleOwner.Equal(roleState.SystemRoleOwner) || !rolePlan.SystemRoleSecurityOwner.Equal(roleState.SystemRoleSecurityOwner) ||
		!rolePlan.ApprovalLevel.Equal(roleState.ApprovalLevel) || !rolePlan.ProductCategory.Equal(roleState.ProductCategory) {

		role := tilgangsportalapi.SystemRoleChange{
			RoleName:         namePlan.ValueString(), // identifier for the role, using plan in case the name was changed above
			L2Ident:          rolePlan.SystemRoleOwner.ValueString(),
			L3Ident:          rolePlan.SystemRoleSecurityOwner.ValueString(),
			NewApprovalLevel: rolePlan.ApprovalLevel.ValueString(),
			NewDescription:   rolePlan.Description.ValueString(),
			ProductCategory:  rolePlan.ProductCategory.ValueString(),
		}

		_, err := r.client.UpdateSystemRole(role)
		if err != nil {

			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update fields of System Role %s, got error: %s", namePlan, err))
			return
		}

	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &rolePlan)...)
}

// Delete a system role resource
func (r *NewSystemRoleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SystemRoleModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	role := tilgangsportalapi.DeleteSystemRole{
		Name:  data.Name.ValueString(),
		Force: "1",
	}

	_, err := r.client.DeleteSystemRole(role)
	if err != nil {
		// Known API error where deletion occasionally fails, handle this by checking if role still exists
		exists, _ := r.client.CheckIfRoleExists(data.Name.ValueString())
		if !exists {
			tflog.Warn(ctx, fmt.Sprintf("Deletion of System Role %s failed, but the role is no longer returned. You may be unable to create a new role with the same name as an internal representation of the role may still remain.", data.Name))
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete System Role %s, got error: %s", data.Name, err))
		return
	}
	// Removal from state is handled automatically by the plugin
}

// ImportState imports a system role to state
func (r *NewSystemRoleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	tflog.Debug(ctx, fmt.Sprintf("Importing System Role with name %s", req.ID))

	// Call the API to fetch role with id, if it exists
	response, err := r.client.GetSystemRole(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to import System Role %s, got error: %s", req.ID, err))
		return
	}

	// Save updated data into Terraform state
	role := SystemRoleModel{
		Name:                    types.StringValue(response.Name),
		Description:             types.StringValue(response.Description),
		ApprovalLevel:           types.StringValue(response.ApprovalLevel),
		SystemRoleOwner:         types.StringValue(response.L2Ident),
		SystemRoleSecurityOwner: types.StringValue(response.L3Ident),
		ProductCategory:         types.StringValue(response.ProductCategory),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &role)...)
}
