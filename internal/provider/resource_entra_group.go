package provider

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"terraform-provider-tilgangsportalen/internal/tilgangsportalapi"
	"time"

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
	EntitlementUID   types.String `tfsdk:"entitlement_uid"`
	EntraIDOID       types.String `tfsdk:"object_id"`
	DisplayName      types.String `tfsdk:"name"`
	Alias            types.String `tfsdk:"alias"`
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
				MarkdownDescription: "Identifier for the Entra Group. Currently, as we do not get a unique ID we can use from the API, ID is set equal to `name`.",
				// Plan modifier to import id from previous state to avoid "know after apply" message
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"object_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Object identifier for the Entra Group.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"entitlement_uid": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique ID of the Entra group (entitlement) in Tilgangsportalen. Will be empty for resources created using provider version `0.12.x` or earlier and imported resources.",
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
			"alias": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Alias for the Entra Group. Deprecated and no longer in use.",
				DeprecationMessage:  "Alias is deprecated and not used in the API. Field will be removed in a future release.",
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
	ConfigureClientResource(ctx, req, resp, func(client *tilgangsportalapi.Client) {
		r.client = client
	})
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

	response, err := r.client.CreateEntraGroup(entraGroup)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Entra Group %s, got error: %s", data.DisplayName, err))
		return
	}

	// Setting role ID to be equal the new role name
	data.Id = data.DisplayName

	// Setting EntitlementUID to the RequestID from the API response
	data.EntitlementUID = types.StringValue(response.RequestID)

	// Check if we need to wait for the object_id to be set for this group
	waitForObjectId := checkIfGroupWillBeCreatedInEntra(r.client, entraGroup.DisplayName, entraGroup.Description)

	// Sleep before polling to reduce API traffic to Tilgangsportalen
	// This is located inside the resource creation function, so it will only run when a resource is created, thus avoiding adding a sleep for data sources
	sleepTime := 77 * time.Second
	if waitForObjectId {
		tflog.Debug(ctx, fmt.Sprintf("Sleeping %v before polling for Entra Group object_id", sleepTime))
		time.Sleep(sleepTime)
	}

	// Get EntraIDOID from the GetAzureADGroup API
	entraGroupRead, err := r.client.GetEntraGroupByID(data.EntitlementUID.ValueString(), waitForObjectId)
	if err != nil {
		resp.Diagnostics.AddError("Client error", fmt.Sprintf("Unable to import entra group %s, got error: %s", data.DisplayName, err))
		return
	}
	data.EntraIDOID = types.StringValue(entraGroupRead.EntraIDOID)

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

	// Check if the group exists
	groupExists, _, err := r.client.CheckIfGroupExists(data.DisplayName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to check if Entra Group %s exists, got error: %s", data.DisplayName, err))
		return
	}

	// Remove from state if group doesn't exists
	if !groupExists {
		resp.State.RemoveResource(ctx)
		return
	}

	// Prefer EntitlementUID lookup when available, otherwise fall back to display name.
	var entraGroup *tilgangsportalapi.EntraGroup
	if !data.EntitlementUID.IsNull() && data.EntitlementUID.ValueString() != "" {
		entraGroup, err = r.client.GetEntraGroupByID(data.EntitlementUID.ValueString(), false)
		if err != nil {
			resp.Diagnostics.AddError("Client error", fmt.Sprintf("Unable to get Entra Group with EntitlementUID %s, got error: %s", data.EntitlementUID.ValueString(), err))
			return
		}
	} else {
		entraGroup, err = r.client.GetEntraGroup(data.DisplayName.ValueString(), false)
		if err != nil {
			resp.Diagnostics.AddError("Client error", fmt.Sprintf("Unable to get Entra Group with display name %s, got error: %s", data.DisplayName.ValueString(), err))
			return
		}
	}

	// Map to EntraGroupModel and save updated data into Terraform state
	data.DisplayName = types.StringValue(entraGroup.DisplayName)
	data.EntitlementUID = types.StringValue(entraGroup.EntitlementUID)
	data.InheritanceLevel = types.StringValue(entraGroup.InheritanceLevel)
	data.EntraIDOID = types.StringValue(entraGroup.EntraIDOID)

	// If no description is set, GetEntraGroup returns an empty string.
	// We only want the plan to show change if description has actually changed
	if entraGroup.Description == "" && data.Description != types.StringValue("") {
		data.Description = types.StringNull()
	} else {
		data.Description = types.StringValue(entraGroup.Description)
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

	// Call the API to fetch group with name, if it exists
	response, err := r.client.GetEntraGroup(req.ID, false)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to import Entra Group %s, got error: %s", req.ID, err))
		return
	}

	group := EntraGroupModel{
		Id:               types.StringValue(response.DisplayName),
		DisplayName:      types.StringValue(response.DisplayName),
		InheritanceLevel: types.StringValue(response.InheritanceLevel),
		Description:      types.StringValue(response.Description),
		EntraIDOID:       types.StringValue(response.EntraIDOID),
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &group)...)
}

// Checks if the resource is created in tilgangsportalen test
func checkIfTestResource(client *tilgangsportalapi.Client) bool {
	testApi := "https://tilgang-test.sits.no/ApiServer"

	return client.GetBaseURL() == testApi
}

// Entra groups are only created in Entra via tilgangsportalen test if they meet certain requirements
func checkIfGroupWillBeCreatedInEntra(client *tilgangsportalapi.Client, name string, description string) bool {
	// assuming all groups created via tilgangsportalen prod to be created in Entra
	if !checkIfTestResource(client) {
		return true
	}

	// Name must start with "[APPTEST]" and description must be equal to "APPTEST" for entra group to be created in Entra via tilgangsportalen test
	// If name starts with "[TESTNOENTRA]", group will not be created in Entra, but the provider will still try to fetch the object_id.
	// This is a workaround for the test TestCreateNewEntraGroupThatIsNotCreatedInEntra
	if strings.HasPrefix(name, "[APPTEST]") && description == "APPTEST" {
		return true
	} else if strings.HasPrefix(name, "[TESTNOENTRA]") {
		return true
	}
	
	return false
}
