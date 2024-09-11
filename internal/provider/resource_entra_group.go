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
var _ resource.Resource = &NewEntraGroupResource{}
var _ resource.ResourceWithImportState = &NewEntraGroupResource{}

// CreateNewEntraGroup is a helper function
func CreateNewEntraGroup() resource.Resource {
	return &NewEntraGroupResource{}
}

// NewEntraGroupResource defines the resource implementation
type NewEntraGroupResource struct {
	client *tilgangsportalapi.Client
}

// EntraGroupModel is a mapping of the resource schema
type EntraGroupModel struct {
	Id               types.String `tfsdk:"id"`
	DisplayName      types.String `tfsdk:"name"`
	Description      types.String `tfsdk:"description"`
	InheritanceLevel types.String `tfsdk:"inheritance_level"`
}

// Metadata returns the resource type name.
func (r *NewEntraGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_entra_group"
}

// Schema defines the schema for the resource.
func (r *NewEntraGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This resource is used to create a new Entra Group using Tilgangsportalen",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Identifier for the Entra Group. Currently, as we do not get a unique ID we can use from the API, ID is set equal to DisplayName",
				// Plan modifier to import id from previous state to avoid "know after apply" message
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The display name of the Entra Group. Must be unique. Please follow the standardized naming conventions for Entra ID groups.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[\[][æøåÆØÅa-zA-Z0-9 _\-\[\]]+$`),
						"The name of the Entra group must start with a prefix enclosed in square brackets, and may only contain alphanumeric characters, space ( ), square brackets ([]), underscore (_), and dash (-). The maxiumum length is 256 characters.",
					),
					stringvalidator.LengthAtMost(256),
				},
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "A description of the Entra Group",
				Validators: []validator.String{
					// Tilgangsportalen API returns an error if the group description is longer than 1024 characters.
					stringvalidator.LengthAtMost(1024),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"inheritance_level": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The inheritance level of the Entra Group (User or Admin). Determines what type of AD account the group can be assigned to.",
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"User", "Admin"}...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *NewEntraGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create is used to create an Entra group resource
func (r *NewEntraGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data EntraGroupModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	entraGroup := tilgangsportalapi.EntraGroup{
		DisplayName:      data.DisplayName.ValueString(),
		Description:      data.Description.ValueString(),
		InheritanceLevel: data.InheritanceLevel.ValueString(),
	}

	_, err := r.client.CreateEntraGroup(entraGroup)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Entra Group %s, got error: %s", data.DisplayName, err))
		return
	}

	// Setting role ID to be equal the new role name
	data.Id = data.DisplayName

	tflog.Debug(ctx, fmt.Sprintf("Entra Group %s created", entraGroup.DisplayName))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read calls the API to get the latest data for the resource
func (r *NewEntraGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data EntraGroupModel

	// Read Terraform prior state data into the model
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// list entra groups belonging to API user and check if the group exists
	groupExists, err := r.client.CheckIfGroupExists(data.DisplayName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to check if Entra Group %s exists, got error: %s", data.DisplayName, err))
		return
	}

	if !groupExists {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update performs changes on Entra ID group names
// If there are changes to other fields than the group name the group will be
// deleted and re-created with the new name.
// This is due to a limitation in the underlying API which is missing a method to update the other fields.
func (r *NewEntraGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data EntraGroupModel
	var namePlan, nameState types.String

	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("name"), &namePlan)...)
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("name"), &nameState)...)

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if !namePlan.Equal(nameState) {

		group := tilgangsportalapi.RenameEntraGroup{
			OldName: nameState.ValueString(),
			NewName: namePlan.ValueString(),
		}

		_, err := r.client.RenameEntraGroup(group)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to rename Entra Group %s to %s, got error: %s", nameState, namePlan, err))
			return
		}

	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete an Entra group resource
func (r *NewEntraGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data EntraGroupModel

	// Read Terraform prior state data into the model
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	group := tilgangsportalapi.DeleteEntraGroup{
		Name:  data.DisplayName.ValueString(),
		Force: "1", // If the group has account assignments, these are also deleted
	}

	_, err := r.client.DeleteEntraGroup(group)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete Entra Group %s, got error: %s", data.DisplayName, err))
		return
	}
}

// ImportState imports an Entra group to state
func (r *NewEntraGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	tflog.Debug(ctx, fmt.Sprintf("Importing Entra Group with ID %s", req.ID))

	data := EntraGroupModel{
		Id:          types.StringValue(req.ID),
		DisplayName: types.StringValue(req.ID),
	}

	// TODO: Get other fields for entra group and add to Terraform state when we have a read group API method
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
