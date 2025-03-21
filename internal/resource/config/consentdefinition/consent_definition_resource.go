// Copyright © 2025 Ping Identity Corporation

package consentdefinition

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &consentDefinitionResource{}
	_ resource.ResourceWithConfigure   = &consentDefinitionResource{}
	_ resource.ResourceWithImportState = &consentDefinitionResource{}
	_ resource.Resource                = &defaultConsentDefinitionResource{}
	_ resource.ResourceWithConfigure   = &defaultConsentDefinitionResource{}
	_ resource.ResourceWithImportState = &defaultConsentDefinitionResource{}
)

// Create a Consent Definition resource
func NewConsentDefinitionResource() resource.Resource {
	return &consentDefinitionResource{}
}

func NewDefaultConsentDefinitionResource() resource.Resource {
	return &defaultConsentDefinitionResource{}
}

// consentDefinitionResource is the resource implementation.
type consentDefinitionResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultConsentDefinitionResource is the resource implementation.
type defaultConsentDefinitionResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *consentDefinitionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_consent_definition"
}

func (r *defaultConsentDefinitionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_consent_definition"
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

func (r *defaultConsentDefinitionResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type consentDefinitionResourceModel struct {
	Id              types.String `tfsdk:"id"`
	Notifications   types.Set    `tfsdk:"notifications"`
	RequiredActions types.Set    `tfsdk:"required_actions"`
	Type            types.String `tfsdk:"type"`
	UniqueID        types.String `tfsdk:"unique_id"`
	DisplayName     types.String `tfsdk:"display_name"`
	Parameter       types.Set    `tfsdk:"parameter"`
	Description     types.String `tfsdk:"description"`
}

// GetSchema defines the schema for the resource.
func (r *consentDefinitionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	consentDefinitionSchema(ctx, req, resp, false)
}

func (r *defaultConsentDefinitionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	consentDefinitionSchema(ctx, req, resp, true)
}

func consentDefinitionSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Consent Definition.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Consent Definition resource. Options are ['consent-definition']",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("consent-definition"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"consent-definition"}...),
				},
			},
			"unique_id": schema.StringAttribute{
				Description: "A version-independent unique identifier for this Consent Definition.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"display_name": schema.StringAttribute{
				Description: "A human-readable display name for this Consent Definition.",
				Optional:    true,
			},
			"parameter": schema.SetAttribute{
				Description: "Optional parameters for this Consent Definition.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Consent Definition",
				Optional:    true,
			},
		},
	}
	if isDefault {
		typeAttr := schemaDef.Attributes["type"].(schema.StringAttribute)
		typeAttr.Optional = false
		typeAttr.Required = false
		typeAttr.Computed = true
		typeAttr.PlanModifiers = []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		}
		schemaDef.Attributes["type"] = typeAttr
		// Add any default properties and set optional properties to computed where necessary
		config.SetAttributesToOptionalAndComputedAndRemoveDefaults(&schemaDef, []string{"type", "unique_id"})
	} else {
		// Add RequiresReplace modifier for read-only attributes
		uniqueIdAttr := schemaDef.Attributes["unique_id"].(schema.StringAttribute)
		uniqueIdAttr.PlanModifiers = append(uniqueIdAttr.PlanModifiers, stringplanmodifier.RequiresReplace())
		schemaDef.Attributes["unique_id"] = uniqueIdAttr
	}
	config.AddCommonResourceSchema(&schemaDef, false)
	resp.Schema = schemaDef
}

// Add optional fields to create request for consent-definition consent-definition
func addOptionalConsentDefinitionFields(ctx context.Context, addRequest *client.AddConsentDefinitionRequest, plan consentDefinitionResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.DisplayName) {
		addRequest.DisplayName = plan.DisplayName.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.Parameter) {
		var slice []string
		plan.Parameter.ElementsAs(ctx, &slice, false)
		addRequest.Parameter = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Populate any computed string values with empty strings, since that is equivalent to null to PD. This will reduce noise in plan output
func (model *consentDefinitionResourceModel) populateAllComputedStringAttributes() {
	if model.UniqueID.IsUnknown() || model.UniqueID.IsNull() {
		model.UniqueID = types.StringValue("")
	}
	if model.Description.IsUnknown() || model.Description.IsNull() {
		model.Description = types.StringValue("")
	}
	if model.DisplayName.IsUnknown() || model.DisplayName.IsNull() {
		model.DisplayName = types.StringValue("")
	}
}

// Read a ConsentDefinitionResponse object into the model struct
func readConsentDefinitionResponse(ctx context.Context, r *client.ConsentDefinitionResponse, state *consentDefinitionResourceModel, expectedValues *consentDefinitionResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("consent-definition")
	state.Id = types.StringValue(r.Id)
	state.UniqueID = types.StringValue(r.UniqueID)
	state.DisplayName = internaltypes.StringTypeOrNil(r.DisplayName, internaltypes.IsEmptyString(expectedValues.DisplayName))
	state.Parameter = internaltypes.GetStringSet(r.Parameter)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
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

// Create a consent-definition consent-definition
func (r *consentDefinitionResource) CreateConsentDefinition(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan consentDefinitionResourceModel) (*consentDefinitionResourceModel, error) {
	addRequest := client.NewAddConsentDefinitionRequest(plan.UniqueID.ValueString(),
		plan.UniqueID.ValueString())
	addOptionalConsentDefinitionFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ConsentDefinitionAPI.AddConsentDefinition(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddConsentDefinitionRequest(*addRequest)

	addResponse, httpResp, err := r.apiClient.ConsentDefinitionAPI.AddConsentDefinitionExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Consent Definition", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state consentDefinitionResourceModel
	readConsentDefinitionResponse(ctx, addResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
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

	state, err := r.CreateConsentDefinition(ctx, req, resp, plan)
	if err != nil {
		return
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, *state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create a new resource
// For edit only resources like this, create doesn't actually "create" anything - it "adopts" the existing
// config object into management by terraform. This method reads the existing config object
// and makes any changes needed to make it match the plan - similar to the Update method.
func (r *defaultConsentDefinitionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan consentDefinitionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ConsentDefinitionAPI.GetConsentDefinition(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.UniqueID.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Consent Definition", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state consentDefinitionResourceModel
	readConsentDefinitionResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.ConsentDefinitionAPI.UpdateConsentDefinition(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.UniqueID.ValueString())
	ops := createConsentDefinitionOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.ConsentDefinitionAPI.UpdateConsentDefinitionExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Consent Definition", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readConsentDefinitionResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
	}

	state.populateAllComputedStringAttributes()
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *consentDefinitionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readConsentDefinition(ctx, req, resp, r.apiClient, r.providerConfig, false)
}

func (r *defaultConsentDefinitionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readConsentDefinition(ctx, req, resp, r.apiClient, r.providerConfig, true)
}

func readConsentDefinition(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, isDefault bool) {
	// Get current state
	var state consentDefinitionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.ConsentDefinitionAPI.GetConsentDefinition(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.UniqueID.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 && !isDefault {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Consent Definition", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Consent Definition", err, httpResp)
		}
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readConsentDefinitionResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	if isDefault {
		state.populateAllComputedStringAttributes()
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *consentDefinitionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateConsentDefinition(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultConsentDefinitionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateConsentDefinition(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateConsentDefinition(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
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
	updateRequest := apiClient.ConsentDefinitionAPI.UpdateConsentDefinition(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.UniqueID.ValueString())

	// Determine what update operations are necessary
	ops := createConsentDefinitionOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.ConsentDefinitionAPI.UpdateConsentDefinitionExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Consent Definition", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readConsentDefinitionResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
// This config object is edit-only, so Terraform can't delete it.
// After running a delete, Terraform will just "forget" about this object and it can be managed elsewhere.
func (r *defaultConsentDefinitionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *consentDefinitionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state consentDefinitionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.ConsentDefinitionAPI.DeleteConsentDefinitionExecute(r.apiClient.ConsentDefinitionAPI.DeleteConsentDefinition(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.UniqueID.ValueString()))
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Consent Definition", err, httpResp)
		return
	}
}

func (r *consentDefinitionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importConsentDefinition(ctx, req, resp)
}

func (r *defaultConsentDefinitionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importConsentDefinition(ctx, req, resp)
}

func importConsentDefinition(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to unique_id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("unique_id"), req, resp)
}
