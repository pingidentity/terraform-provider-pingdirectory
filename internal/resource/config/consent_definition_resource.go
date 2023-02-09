package config

import (
	"context"
	"time"

	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9100"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &consentDefinitionResource{}
	_ resource.ResourceWithConfigure   = &consentDefinitionResource{}
	_ resource.ResourceWithImportState = &consentDefinitionResource{}
)

// Create a Consent Definition resource
func NewConsentDefinitionResource() resource.Resource {
	return &consentDefinitionResource{}
}

// consentDefinitionResource is the resource implementation.
type consentDefinitionResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *consentDefinitionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_consent_definition"
}

// Configure adds the provider configured client to the resource.
func (r *consentDefinitionResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type consentDefinitionResourceModel struct {
	Id              types.String `tfsdk:"id"`
	LastUpdated     types.String `tfsdk:"last_updated"`
	Notifications   types.Set    `tfsdk:"notifications"`
	RequiredActions types.Set    `tfsdk:"required_actions"`
	UniqueID        types.String `tfsdk:"unique_id"`
	DisplayName     types.String `tfsdk:"display_name"`
	Parameter       types.Set    `tfsdk:"parameter"`
	Description     types.String `tfsdk:"description"`
}

// GetSchema defines the schema for the resource.
func (r *consentDefinitionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages a Consent Definition.",
		Attributes: map[string]schema.Attribute{
			"unique_id": schema.StringAttribute{
				Description: "A version-independent unique identifier for this Consent Definition.",
				Required:    true,
			},
			"display_name": schema.StringAttribute{
				Description: "A human-readable display name for this Consent Definition.",
				Optional:    true,
			},
			"parameter": schema.SetAttribute{
				Description: "Optional parameters for this Consent Definition.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Consent Definition",
				Optional:    true,
			},
		},
	}
	AddCommonSchema(&schema, true)
	resp.Schema = schema
}

// Add optional fields to create request
func addOptionalConsentDefinitionFields(ctx context.Context, addRequest *client.AddConsentDefinitionRequest, plan consentDefinitionResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.DisplayName) {
		stringVal := plan.DisplayName.ValueString()
		addRequest.DisplayName = &stringVal
	}
	if internaltypes.IsDefined(plan.Parameter) {
		var slice []string
		plan.Parameter.ElementsAs(ctx, &slice, false)
		addRequest.Parameter = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		stringVal := plan.Description.ValueString()
		addRequest.Description = &stringVal
	}
}

// Read a ConsentDefinitionResponse object into the model struct
func readConsentDefinitionResponse(ctx context.Context, r *client.ConsentDefinitionResponse, state *consentDefinitionResourceModel, expectedValues *consentDefinitionResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.UniqueID = types.StringValue(r.UniqueID)
	state.DisplayName = internaltypes.StringTypeOrNil(r.DisplayName, internaltypes.IsEmptyString(expectedValues.DisplayName))
	state.Parameter = internaltypes.GetStringSet(r.Parameter)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createConsentDefinitionOperations(plan consentDefinitionResourceModel, state consentDefinitionResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.UniqueID, state.UniqueID, "unique-id")
	operations.AddStringOperationIfNecessary(&ops, plan.DisplayName, state.DisplayName, "display-name")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.Parameter, state.Parameter, "parameter")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	return ops
}

// Create a new resource
func (r *consentDefinitionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan consentDefinitionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddConsentDefinitionRequest(plan.Id.ValueString(),
		plan.UniqueID.ValueString())
	addOptionalConsentDefinitionFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ConsentDefinitionApi.AddConsentDefinition(
		ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddConsentDefinitionRequest(*addRequest)

	addResponse, httpResp, err := r.apiClient.ConsentDefinitionApi.AddConsentDefinitionExecute(apiAddRequest)
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Consent Definition", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state consentDefinitionResourceModel
	readConsentDefinitionResponse(ctx, addResponse, &state, &plan, &resp.Diagnostics)

	// Populate Computed attribute values
	state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *consentDefinitionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state consentDefinitionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ConsentDefinitionApi.GetConsentDefinition(
		ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Consent Definition", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readConsentDefinitionResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *consentDefinitionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan consentDefinitionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state consentDefinitionResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.ConsentDefinitionApi.UpdateConsentDefinition(
		ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createConsentDefinitionOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.ConsentDefinitionApi.UpdateConsentDefinitionExecute(updateRequest)
		if err != nil {
			ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Consent Definition", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readConsentDefinitionResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
		// Update computed values
		state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))
	} else {
		tflog.Warn(ctx, "No configuration API operations created for update")
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *consentDefinitionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state consentDefinitionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.ConsentDefinitionApi.DeleteConsentDefinitionExecute(r.apiClient.ConsentDefinitionApi.DeleteConsentDefinition(
		ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Consent Definition", err, httpResp)
		return
	}
}

func (r *consentDefinitionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
